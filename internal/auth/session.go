package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gorilla/sessions"
)

const SessionCookieName = "atlas_session"

// Store wraps gorilla/sessions CookieStore.
type Store struct {
	cookieStore *sessions.CookieStore
}

// NewStore creates a session store using the provided 32-byte secret key.
func NewStore(secret string, maxAge int) *Store {
	if secret == "" {
		// Dev fallback secret if empty
		secret = "atlas-dev-secret-key-32bytes-long!"
	}

	cs := sessions.NewCookieStore([]byte(secret))
	cs.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Will be set to true in production behind HTTPS
	}

	return &Store{cookieStore: cs}
}

// Get returns the session for the HTTP request.
func (s *Store) Get(r *http.Request) (*sessions.Session, error) {
	return s.cookieStore.Get(r, SessionCookieName)
}

// GenerateToken generates a cryptographically random 32-byte hex session token string.
func GenerateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
