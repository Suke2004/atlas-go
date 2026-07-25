package setup

import (
	"net/http"

	"github.com/Suke2004/atlas-go/internal/auth"
	templates "github.com/Suke2004/atlas-go/web/templates/setup"
	"go.uber.org/zap"
)

// Handler manages HTTP routes for the setup wizard.
type Handler struct {
	svc     *Service
	authSvc *auth.Service
	log     *zap.Logger
}

// NewHandler constructs a setup Handler.
func NewHandler(svc *Service, authSvc *auth.Service, log *zap.Logger) *Handler {
	return &Handler{
		svc:     svc,
		authSvc: authSvc,
		log:     log,
	}
}

// ShowWizard handles GET /setup.
func (h *Handler) ShowWizard(w http.ResponseWriter, r *http.Request) {
	isFirst, err := h.svc.IsFirstRun(r.Context())
	if err != nil {
		h.log.Error("setup: check first run failed", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !isFirst {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.Wizard("").Render(r.Context(), w)
}

// ProcessSetup handles POST /setup.
func (h *Handler) ProcessSetup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_ = templates.Wizard("Invalid form parameters").Render(r.Context(), w)
		return
	}

	input := CreateFirstUserInput{
		Username:    r.FormValue("username"),
		DisplayName: r.FormValue("display_name"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		Timezone:    r.FormValue("timezone"),
		Theme:       r.FormValue("theme"),
	}

	user, err := h.svc.CreateFirstUser(r.Context(), input)
	if err != nil {
		h.log.Warn("setup: create user failed", zap.Error(err))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = templates.Wizard(err.Error()).Render(r.Context(), w)
		return
	}

	// Create session and set session cookie
	sess, err := h.authSvc.CreateSession(r.Context(), user.ID, 24*365) // 1 year session for owner
	if err == nil {
		cookieSess, _ := h.authSvc.Store().Get(r)
		cookieSess.Values["token"] = sess.ID
		cookieSess.Values["user_id"] = user.ID
		_ = cookieSess.Save(r, w)
	}

	h.log.Info("setup: first run owner created successfully", zap.String("username", user.Username))
	http.Redirect(w, r, "/setup/demo-choice", http.StatusSeeOther)
}

// ShowDemoChoice handles GET /setup/demo-choice.
func (h *Handler) ShowDemoChoice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.DemoChoice().Render(r.Context(), w)
}

// ProcessDemoSeed handles POST /setup/seed.
func (h *Handler) ProcessDemoSeed(w http.ResponseWriter, r *http.Request) {
	cookieSess, err := h.authSvc.Store().Get(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, ok := cookieSess.Values["user_id"].(int64)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := h.svc.SeedDemoData(r.Context(), userID); err != nil {
		h.log.Error("setup: seed demo data failed", zap.Error(err))
	} else {
		h.log.Info("setup: demo data seeded successfully", zap.Int64("user_id", userID))
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
