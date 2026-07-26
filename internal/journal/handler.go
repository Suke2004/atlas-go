package journal

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
	journaltemplates "github.com/Suke2004/atlas-go/web/templates/journal"
)

type Handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// GET /journal
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	details, err := h.service.GetDailyJournal(r.Context(), user.ID, dateStr)
	if err != nil {
		h.logger.Error("failed to load daily journal", zap.Error(err))
	}

	summary, _ := h.service.GetJournalSummary(r.Context(), user.ID)

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return journaltemplates.Journal(*details, summary, dateStr).Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("Daily Journal & Mind-Sync", "/journal", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /journal/save (On-blur / 30s auto-save)
func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dateStr := r.FormValue("entry_date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	mood, _ := strconv.Atoi(r.FormValue("mood_rating"))
	energy, _ := strconv.Atoi(r.FormValue("energy_rating"))
	sleep, _ := strconv.ParseFloat(r.FormValue("sleep_hours"), 64)

	input := JournalInput{
		EntryDate:    dateStr,
		MoodRating:   mood,
		EnergyRating: energy,
		SleepHours:  sleep,
		Summary:      r.FormValue("summary"),
	}

	_, err := h.service.SaveDailyJournal(r.Context(), user.ID, input)
	if err != nil {
		h.logger.Error("failed to save journal entry", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<span class="text-xs text-emerald-400 font-mono-code flex items-center gap-1"><i data-lucide="check-circle-2" class="w-3.5 h-3.5"></i>Synced</span>`))
}

// POST /journal/items
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	entryID, _ := strconv.ParseInt(r.FormValue("entry_id"), 10, 64)
	category := r.FormValue("category")
	content := r.FormValue("content")

	if content != "" && entryID > 0 {
		_, err := h.service.AddJournalItem(r.Context(), user.ID, entryID, category, content)
		if err != nil {
			h.logger.Error("failed to add journal item", zap.Error(err))
		}
	}

	dateStr := r.FormValue("entry_date")
	http.Redirect(w, r, "/journal?date="+dateStr, http.StatusSeeOther)
}

// POST /journal/items/{id}/delete
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	itemID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	_ = h.service.DeleteJournalItem(r.Context(), itemID)
	dateStr := r.URL.Query().Get("date")
	http.Redirect(w, r, "/journal?date="+dateStr, http.StatusSeeOther)
}
