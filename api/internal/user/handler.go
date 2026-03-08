package user

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

// List GET /api/users?page=1&limit=20&role=admin&status=active
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	params := pagination.Parse(r, h.config.Pagination.DefaultLimit, h.config.Pagination.MaxLimit)

	role := r.URL.Query().Get("role")
	status := r.URL.Query().Get("status")
	sort := pagination.ParseSort(r, []string{
		"created_at", "first_name", "last_name", "role", "status",
	}, "created_at")

	users, total, err := h.repo.List(r.Context(), params.Limit, params.Offset, role, status, sort.Column, sort.Order)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	response.Paginated(w, users, params.Page, params.Limit, total)
}

// GetByID GET /api/users/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	if user == nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, user)
}
