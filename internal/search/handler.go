package search

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	searchtemplates "github.com/Suke2004/atlas-go/web/templates/search"
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

// GET /api/search
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("q")
	var searchResults []SearchResultItem
	if query != "" {
		res, err := h.service.Search(r.Context(), user.ID, query)
		if err != nil {
			h.logger.Error("failed to perform search", zap.Error(err))
		} else {
			searchResults = res
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = searchtemplates.Results(searchResults, query).Render(r.Context(), w)
}
