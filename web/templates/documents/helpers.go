package documents

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Suke2004/atlas-go/internal/db"
)

// DocumentWithMeta is a template-layer enrichment of db.Document.
// It lives here (not in internal/documents) to avoid import cycles between
// the service layer and the template package.
type DocumentWithMeta struct {
	db.Document
	HumanSize   string
	FileIcon    string
	PreviewType string
	TagList     []string
}

// EnrichDocument converts a db.Document into a DocumentWithMeta.
func EnrichDocument(doc db.Document) DocumentWithMeta {
	dm := DocumentWithMeta{Document: doc}
	dm.HumanSize = HumanFileSize(doc.FileSize)
	dm.FileIcon = FileIcon(doc.MimeType, doc.Filename)
	dm.PreviewType = PreviewType(doc.MimeType, doc.Filename)
	var tags []string
	if err := json.Unmarshal([]byte(doc.Tags), &tags); err == nil {
		dm.TagList = tags
	}
	return dm
}

// EnrichDocuments converts a slice of db.Document to DocumentWithMeta.
func EnrichDocuments(docs []db.Document) []DocumentWithMeta {
	out := make([]DocumentWithMeta, len(docs))
	for i, d := range docs {
		out[i] = EnrichDocument(d)
	}
	return out
}

// UploadFormattedDate returns a human-readable date.
func UploadFormattedDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}

// NullString safely extracts a sql.NullString value.
func NullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// joinTags renders a comma-separated tag string for the edit form.
func joinTags(tags []string) string {
	return strings.Join(tags, ", ")
}

// HumanFileSize formats bytes as a human-readable size string.
func HumanFileSize(bytes int64) string {
	switch {
	case bytes >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(bytes)/(1<<30))
	case bytes >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(bytes)/(1<<10))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// FileIcon returns the Lucide icon name for a file.
func FileIcon(mimeType, filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case mimeType == "application/pdf":
		return "file-text"
	case ext == ".md" || ext == ".txt":
		return "file-text"
	case ext == ".go", ext == ".py", ext == ".js", ext == ".ts",
		ext == ".java", ext == ".c", ext == ".cpp", ext == ".rs":
		return "code-2"
	case ext == ".xlsx" || ext == ".csv":
		return "table"
	case ext == ".docx" || ext == ".doc":
		return "file-type"
	case ext == ".zip" || ext == ".tar" || ext == ".gz":
		return "archive"
	case ext == ".mp4" || ext == ".mov" || ext == ".avi":
		return "film"
	case ext == ".mp3" || ext == ".wav" || ext == ".flac":
		return "music"
	default:
		return "file"
	}
}

// PreviewType classifies a file for inline preview rendering.
func PreviewType(mimeType, filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch {
	case mimeType == "application/pdf":
		return "pdf"
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case IsTextMime(mimeType) || ext == ".md" || ext == ".txt" || ext == ".csv":
		return "text"
	default:
		return "binary"
	}
}

// IsTextMime returns true for MIME types that can be rendered as text.
func IsTextMime(mimeType string) bool {
	return strings.HasPrefix(mimeType, "text/") ||
		mimeType == "application/json" ||
		mimeType == "application/xml"
}

// previewBadgeClass returns TailwindCSS colour classes for the preview type badge.
func previewBadgeClass(pt string) string {
	switch pt {
	case "pdf":
		return "bg-rose-500/10 text-rose-400 border-rose-500/20"
	case "image":
		return "bg-emerald-500/10 text-emerald-400 border-emerald-500/20"
	case "text":
		return "bg-sky-500/10 text-sky-400 border-sky-500/20"
	default:
		return "bg-slate-700/40 text-slate-400 border-slate-600/20"
	}
}

// previewLabel returns a display label for the preview type.
func previewLabel(pt string) string {
	switch pt {
	case "pdf":
		return "PDF"
	case "image":
		return "Image"
	case "text":
		return "Text"
	default:
		return "Binary"
	}
}
