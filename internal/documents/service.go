package documents

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/Suke2004/atlas-go/internal/ai"
	"github.com/Suke2004/atlas-go/internal/db"
)

const (
	// MaxUploadSize is the maximum allowed file size per upload (100 MB).
	MaxUploadSize = 100 << 20

	// UploadsDir is the root directory for document storage, relative to the working directory.
	UploadsDir = "data/uploads"
)

// Service defines the business logic interface for the documents module.
type Service interface {
	Upload(ctx context.Context, userID int64, file multipart.File, header *multipart.FileHeader) (db.Document, error)
	Get(ctx context.Context, userID, id int64) (db.Document, error)
	List(ctx context.Context, userID int64) ([]db.Document, error)
	UpdateMeta(ctx context.Context, userID, id int64, title string, tags []string) (db.Document, error)
	Summarise(ctx context.Context, userID, id int64) (string, error)
	Delete(ctx context.Context, userID, id int64) error
	Count(ctx context.Context, userID int64) (int64, error)
}

type service struct {
	repo       Repository
	aiProvider ai.Provider // may be nil if AI not configured
}

// NewService creates a documents service.
// aiProvider may be nil — summaries will return an error in that case.
func NewService(repo Repository, aiProvider ai.Provider) Service {
	return &service{repo: repo, aiProvider: aiProvider}
}

// Upload stores the file to disk and creates a DB record.
func (s *service) Upload(ctx context.Context, userID int64, file multipart.File, header *multipart.FileHeader) (db.Document, error) {
	if header.Size > MaxUploadSize {
		return db.Document{}, fmt.Errorf("file too large: maximum upload size is 100 MB")
	}

	// Derive a safe storage filename: uuid + original extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	storedName := uuid.New().String() + ext
	storagePath := filepath.Join(UploadsDir, storedName)

	// Ensure the uploads directory exists
	if err := os.MkdirAll(UploadsDir, 0755); err != nil {
		return db.Document{}, fmt.Errorf("failed to create uploads dir: %w", err)
	}

	// Write file to disk
	out, err := os.Create(storagePath)
	if err != nil {
		return db.Document{}, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		os.Remove(storagePath) // cleanup on failure
		return db.Document{}, fmt.Errorf("failed to write file: %w", err)
	}

	// Detect MIME type from extension
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Use the original filename (without ext) as the default title
	title := strings.TrimSuffix(header.Filename, ext)
	if title == "" {
		title = header.Filename
	}

	doc, err := s.repo.Create(ctx, db.CreateDocumentParams{
		UserID:       userID,
		Filename:     storedName,
		OriginalName: header.Filename,
		MimeType:     mimeType,
		FileSize:     header.Size,
		StoragePath:  storagePath,
		Title:        title,
		Tags:         "[]",
	})
	if err != nil {
		os.Remove(storagePath) // rollback file on DB error
		return db.Document{}, fmt.Errorf("failed to save document record: %w", err)
	}

	// For text files, extract content for FTS5 search (non-fatal)
	if isTextMime(mimeType) || ext == ".md" || ext == ".txt" {
		_ = s.extractTextContent(ctx, doc, storagePath)
	}

	return doc, nil
}

func (s *service) extractTextContent(ctx context.Context, doc db.Document, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return s.repo.UpdateContent(ctx, doc.ID, doc.UserID, string(data))
}

func (s *service) Get(ctx context.Context, userID, id int64) (db.Document, error) {
	doc, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return db.Document{}, fmt.Errorf("document not found: %w", err)
	}
	return doc, nil
}

func (s *service) List(ctx context.Context, userID int64) ([]db.Document, error) {
	return s.repo.List(ctx, userID)
}

// UpdateMeta saves the document title and tags.
func (s *service) UpdateMeta(ctx context.Context, userID, id int64, title string, tags []string) (db.Document, error) {
	if title == "" {
		return db.Document{}, fmt.Errorf("title cannot be empty")
	}
	// Build JSON tags string
	tagsJSON := "[]"
	if len(tags) > 0 {
		parts := make([]string, len(tags))
		for i, t := range tags {
			parts[i] = `"` + strings.ReplaceAll(t, `"`, `\"`) + `"`
		}
		tagsJSON = "[" + strings.Join(parts, ",") + "]"
	}
	return s.repo.UpdateMeta(ctx, db.UpdateDocumentMetaParams{
		Title:  title,
		Tags:   tagsJSON,
		ID:     id,
		UserID: userID,
	})
}

// Summarise calls the AI provider to generate a summary.
func (s *service) Summarise(ctx context.Context, userID, id int64) (string, error) {
	if s.aiProvider == nil {
		return "", fmt.Errorf("AI provider not configured — set one in Settings → AI Provider")
	}

	doc, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return "", fmt.Errorf("document not found: %w", err)
	}

	if !doc.ContentText.Valid || strings.TrimSpace(doc.ContentText.String) == "" {
		return "", fmt.Errorf("document has no extractable text to summarise (only plain text files are supported; OCR coming in v0.7.0)")
	}

	// Truncate to ~8000 chars to stay within token limits
	content := doc.ContentText.String
	if len(content) > 8000 {
		content = content[:8000] + "\n...[content truncated]"
	}

	messages := []ai.Message{
		{Role: ai.RoleSystem, Content: ai.SystemPromptFor("summarise")},
		{Role: ai.RoleUser, Content: fmt.Sprintf("Please summarise the following document titled %q:\n\n%s", doc.OriginalName, content)},
	}

	summary, err := s.aiProvider.Complete(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("AI summarisation failed: %w", err)
	}

	// Persist the summary (non-fatal)
	_ = s.repo.UpdateSummary(ctx, id, userID, summary)

	return summary, nil
}

// Delete removes the document record and its file from disk.
func (s *service) Delete(ctx context.Context, userID, id int64) error {
	doc, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	if err := os.Remove(doc.StoragePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file from disk: %w", err)
	}

	return s.repo.Delete(ctx, userID, id)
}

func (s *service) Count(ctx context.Context, userID int64) (int64, error) {
	return s.repo.Count(ctx, userID)
}

func isTextMime(mimeType string) bool {
	return strings.HasPrefix(mimeType, "text/") ||
		mimeType == "application/json" ||
		mimeType == "application/xml"
}
