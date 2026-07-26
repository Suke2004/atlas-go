package unit

import (
	"bytes"
	"context"
	"mime/multipart"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/documents"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestDocuments_UploadAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "docs_test.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	if err := db.MigrateUp(database.Raw); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	ctx := context.Background()
	setupSvc := setup.NewService(database)
	user, err := setupSvc.CreateFirstUser(ctx, setup.CreateFirstUserInput{
		Username:    "docuser",
		DisplayName: "Doc User",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	repo := documents.NewRepository(database)
	svc := documents.NewService(repo, nil)

	// Create test file in multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test_doc.txt")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	content := "Hello Atlas Document Engine"
	part.Write([]byte(content))
	writer.Close()

	// Parse header for service Upload
	req := multipart.NewReader(body, writer.Boundary())
	form, err := req.ReadForm(1000)
	if err != nil {
		t.Fatalf("failed to read form: %v", err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("no file in parsed form")
	}

	fileHeader := files[0]
	file, err := fileHeader.Open()
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	doc, err := svc.Upload(ctx, user.ID, file, fileHeader)
	if err != nil {
		t.Fatalf("failed to upload doc: %v", err)
	}

	if doc.OriginalName != "test_doc.txt" {
		t.Errorf("expected original_name test_doc.txt, got %s", doc.OriginalName)
	}

	fetched, err := svc.Get(ctx, user.ID, doc.ID)
	if err != nil {
		t.Fatalf("failed to get doc: %v", err)
	}

	if fetched.Title != "test_doc" {
		t.Errorf("expected title test_doc, got %s", fetched.Title)
	}

	// Clean up
	_ = svc.Delete(ctx, user.ID, doc.ID)
}
