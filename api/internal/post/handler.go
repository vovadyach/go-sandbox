package post

import (
	"go-sandbox/api/internal/config"
	"go-sandbox/api/internal/pagination"
	"go-sandbox/api/internal/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	repo   *Repository
	config *config.Config
}

func NewHandler(repo *Repository, config *config.Config) *Handler {
	return &Handler{
		repo,
		config,
	}
}

// List GET /api/posts?page=1&limit=20&status=published
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params := pagination.Parse(r, h.config.Pagination.DefaultLimit, h.config.Pagination.MaxLimit)

	status := r.URL.Query().Get("status")
	sort := pagination.ParseSort(r, []string{
		"created_at", "title", "status",
	}, "created_at")

	posts, total, err := h.repo.List(r.Context(), params.Limit, params.Offset, status, sort.Column, sort.Order)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	response.Paginated(w, posts, params.Page, params.Limit, total)
}

// ListByUserID GET /api/users/{id}/posts?page=1&limit=20&status=published
func (h *Handler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	params := pagination.Parse(r, h.config.Pagination.DefaultLimit, h.config.Pagination.MaxLimit)

	userID := chi.URLParam(r, "id")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, "user id is required")
		return
	}
	if _, err := uuid.Parse(userID); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	posts, total, err := h.repo.ListByUserID(r.Context(), userID, params.Limit, params.Offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	response.Paginated(w, posts, params.Page, params.Limit, total)
}
