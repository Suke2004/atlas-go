package documents

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	doclayouts "github.com/Suke2004/atlas-go/web/templates/documents"
	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
)

// Handler serves all HTTP endpoints for the documents module.
type Handler struct {
	service Service
	logger  *zap.Logger
}

// NewHandler constructs a documents handler.
func NewHandler(svc Service, logger *zap.Logger) *Handler {
	return &Handler{service: svc, logger: logger}
}

// GET /documents — library page.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rawDocs, err := h.service.List(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to list documents", zap.Error(err))
		http.Error(w, "Failed to load documents", http.StatusInternalServerError)
		return
	}

	docs := doclayouts.EnrichDocuments(rawDocs)
	count, _ := h.service.Count(r.Context(), user.ID)
	errorMsg := r.URL.Query().Get("error")

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return doclayouts.List(docs, count, errorMsg).Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("Documents", "/documents", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /documents — upload a new document.
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Limit total body size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize+1<<20)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Warn("failed to parse multipart form", zap.Error(err))
		http.Redirect(w, r, "/documents?error=invalid_form", http.StatusSeeOther)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Warn("no file in form", zap.Error(err))
		http.Redirect(w, r, "/documents?error=no_file", http.StatusSeeOther)
		return
	}
	defer file.Close()

	doc, err := h.service.Upload(r.Context(), user.ID, file, header)
	if err != nil {
		h.logger.Error("upload failed", zap.Error(err))
		http.Redirect(w, r, "/documents?error=upload_failed", http.StatusSeeOther)
		return
	}

	h.logger.Info("document uploaded",
		zap.Int64("doc_id", doc.ID),
		zap.String("original_name", doc.OriginalName),
		zap.Int64("size", doc.FileSize),
	)
	http.Redirect(w, r, "/documents/"+strconv.FormatInt(doc.ID, 10), http.StatusSeeOther)
}

// GET /documents/{id} — document detail page.
func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	rawDoc, err := h.service.Get(r.Context(), user.ID, id)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	doc := doclayouts.EnrichDocument(rawDoc)

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return doclayouts.Detail(doc).Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base(doc.OriginalName+" — Documents", "/documents", user.DisplayName, pageContent).Render(r.Context(), w)
}

// GET /documents/{id}/raw — serve the raw file for inline preview.
func (h *Handler) ServeRaw(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	doc, err := h.service.Get(r.Context(), user.ID, id)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", doc.MimeType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+doc.OriginalName+"\"")

	f, err := os.Open(doc.StoragePath)
	if err != nil {
		http.Error(w, "File not found on disk", http.StatusNotFound)
		return
	}
	defer f.Close()

	http.ServeContent(w, r, doc.OriginalName, doc.CreatedAt, f)
}

// GET /documents/{id}/download — force-download the file.
func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	doc, err := h.service.Get(r.Context(), user.ID, id)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+doc.OriginalName+"\"")
	w.Header().Set("Content-Type", doc.MimeType)

	http.ServeFile(w, r, doc.StoragePath)
}

// POST /documents/{id}/meta — update title and tags.
func (h *Handler) UpdateMeta(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	tagsRaw := strings.TrimSpace(r.FormValue("tags"))

	var tags []string
	for _, t := range strings.Split(tagsRaw, ",") {
		if t = strings.TrimSpace(t); t != "" {
			tags = append(tags, t)
		}
	}

	if _, err := h.service.UpdateMeta(r.Context(), user.ID, id, title, tags); err != nil {
		h.logger.Error("failed to update document meta", zap.Error(err))
		http.Redirect(w, r, "/documents/"+strconv.FormatInt(id, 10)+"?error=save_failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/documents/"+strconv.FormatInt(id, 10)+"?saved=1", http.StatusSeeOther)
}

// POST /documents/{id}/summarise — trigger AI summarisation (JSON response for fetch()).
func (h *Handler) Summarise(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	summary, err := h.service.Summarise(r.Context(), user.ID, id)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}

// POST /documents/{id}/delete — delete a document.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), user.ID, id); err != nil {
		h.logger.Error("failed to delete document", zap.Error(err))
		http.Redirect(w, r, "/documents?error=delete_failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/documents", http.StatusSeeOther)
}
