package unit

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/Suke2004/atlas-go/internal/auth"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestAuth_AuthenticateAndSession(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "auth_test.db")

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
	authSvc := auth.NewService(database, auth.NewStore("test-secret-key-32-bytes-long!", 3600))

	// 1. Create first user via setup service
	user, err := setupSvc.CreateFirstUser(ctx, setup.CreateFirstUserInput{
		Username:    "testuser",
		DisplayName: "Test User",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	// 2. Test successful authentication
	authUser, err := authSvc.Authenticate(ctx, "testuser", "password123")
	if err != nil {
		t.Fatalf("expected successful auth, got: %v", err)
	}

	if authUser.ID != user.ID {
		t.Errorf("expected user ID %d, got %d", user.ID, authUser.ID)
	}

	// 3. Test invalid password
	_, err = authSvc.Authenticate(ctx, "testuser", "wrongpassword")
	if err != auth.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}

	// 4. Test session creation and validation
	sess, err := authSvc.CreateSession(ctx, user.ID, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	valUser, err := authSvc.ValidateSession(ctx, sess.ID)
	if err != nil {
		t.Fatalf("failed to validate active session: %v", err)
	}

	if valUser.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", valUser.Username)
	}

	// 5. Test session destruction
	if err := authSvc.DestroySession(ctx, sess.ID); err != nil {
		t.Fatalf("failed to destroy session: %v", err)
	}

	_, err = authSvc.ValidateSession(ctx, sess.ID)
	if err != auth.ErrSessionExpired {
		t.Errorf("expected ErrSessionExpired after destruction, got: %v", err)
	}
}
