package post

import (
	"go-sandbox/api/internal/config"
	"go-sandbox/api/internal/pagination"
	"go-sandbox/api/internal/response"
	"net/http"
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

// List GET /api/posts?page=1&limit=20&status=active
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params := pagination.Parse(r, h.config.Pagination.DefaultLimit, h.config.Pagination.MaxLimit)

	status := r.URL.Query().Get("status")
	posts, total, err := h.repo.List(r.Context(), params.Limit, params.Offset, status)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	response.Paginated(w, posts, params.Page, params.Limit, total)
}
