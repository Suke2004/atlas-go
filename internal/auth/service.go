// Package auth provides authentication, session management, password verification, and security middleware.
package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Suke2004/atlas-go/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("auth: invalid username or password")
	ErrSessionExpired     = errors.New("auth: session expired or invalid")
)

// Service provides business logic for user authentication and session validation.
type Service struct {
	db          *db.DB
	sessionStore *Store
}

// NewService creates a new auth Service.
func NewService(database *db.DB, store *Store) *Service {
	return &Service{
		db:           database,
		sessionStore: store,
	}
}

// Authenticate verifies the username and password, returns the User if valid.
func (s *Service) Authenticate(ctx context.Context, username, password string) (*db.User, error) {
	user, err := s.db.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

// CreateSession creates a new session in SQLite for the user.
func (s *Service) CreateSession(ctx context.Context, userID int64, duration time.Duration) (*db.Session, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("auth: failed to generate session token: %w", err)
	}

	expiresAt := time.Now().Add(duration)
	sess, err := s.db.CreateSession(ctx, db.CreateSessionParams{
		ID:        token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("auth: failed to create session in db: %w", err)
	}

	return &sess, nil
}

// ValidateSession retrieves and validates an active session token.
func (s *Service) ValidateSession(ctx context.Context, sessionID string) (*db.User, error) {
	sess, err := s.db.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, ErrSessionExpired
	}

	user, err := s.db.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, ErrSessionExpired
	}

	return &user, nil
}

// DestroySession deletes a session by token ID.
func (s *Service) DestroySession(ctx context.Context, sessionID string) error {
	return s.db.DeleteSession(ctx, sessionID)
}

// Store returns the underlying session store.
func (s *Service) Store() *Store {
	return s.sessionStore
}
