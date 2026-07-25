package auth

import (
	"net/http"
	"time"

	templates "github.com/Suke2004/atlas-go/web/templates/auth"
	"go.uber.org/zap"
)

// Handler manages HTTP routes for authentication (login & logout).
type Handler struct {
	svc *Service
	log *zap.Logger
}

// NewHandler constructs an auth Handler.
func NewHandler(svc *Service, log *zap.Logger) *Handler {
	return &Handler{
		svc: svc,
		log: log,
	}
}

// ShowLogin handles GET /login.
func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.Login("").Render(r.Context(), w)
}

// ProcessLogin handles POST /login.
func (h *Handler) ProcessLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_ = templates.Login("Invalid credentials").Render(r.Context(), w)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.svc.Authenticate(r.Context(), username, password)
	if err != nil {
		h.log.Warn("auth: failed login attempt", zap.String("username", username))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		_ = templates.Login("Invalid username or password").Render(r.Context(), w)
		return
	}

	sess, err := h.svc.CreateSession(r.Context(), user.ID, 24*30*time.Hour) // 30-day session
	if err != nil {
		h.log.Error("auth: failed to create session", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cookieSess, _ := h.svc.Store().Get(r)
	cookieSess.Values["token"] = sess.ID
	cookieSess.Values["user_id"] = user.ID
	_ = cookieSess.Save(r, w)

	h.log.Info("auth: user logged in successfully", zap.String("username", user.Username))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout handles POST /logout.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookieSess, err := h.svc.Store().Get(r)
	if err == nil {
		if token, ok := cookieSess.Values["token"].(string); ok {
			_ = h.svc.DestroySession(r.Context(), token)
		}
		cookieSess.Options.MaxAge = -1
		_ = cookieSess.Save(r, w)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
