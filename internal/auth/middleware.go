package auth

import (
	"context"
	"net/http"

	"github.com/Suke2004/atlas-go/internal/db"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthRequired middleware validates the session cookie and attaches the user to request context.
// Unauthenticated requests are redirected to /login.
func AuthRequired(svc *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := svc.Store().Get(r)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			token, ok := sess.Values["token"].(string)
			if !ok || token == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			user, err := svc.ValidateSession(r.Context(), token)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the authenticated *db.User from the request context.
func GetUserFromContext(ctx context.Context) *db.User {
	user, ok := ctx.Value(UserContextKey).(*db.User)
	if !ok {
		return nil
	}
	return user
}
