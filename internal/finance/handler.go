package finance

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
	financetemplates "github.com/Suke2004/atlas-go/web/templates/finance"
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

// GET /finance
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	txns, err := h.service.ListTransactions(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to list transactions", zap.Error(err))
	}

	summary, _ := h.service.GetFinanceSummary(r.Context(), user.ID)

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return financetemplates.Finance(txns, summary).Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("Finance & Cost Attribution", "/finance", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /finance
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	input := TransactionInput{
		Amount:          amount,
		Type:            r.FormValue("type"),
		Category:        r.FormValue("category"),
		Description:     r.FormValue("description"),
		TransactionDate: r.FormValue("transaction_date"),
	}

	_, err := h.service.CreateTransaction(r.Context(), user.ID, input)
	if err != nil {
		h.logger.Error("failed to create transaction", zap.Error(err))
	}

	http.Redirect(w, r, "/finance", http.StatusSeeOther)
}

// POST /finance/{id}/delete
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	txnID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	_ = h.service.DeleteTransaction(r.Context(), user.ID, txnID)
	http.Redirect(w, r, "/finance", http.StatusSeeOther)
}
