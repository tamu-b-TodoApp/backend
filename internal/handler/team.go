package handler

import (
	"errors"
	"net/http"
	"strconv"

	"todo/internal/service"
)

type TeamHandler struct {
	service        service.TeamService
	authMiddleware func(http.Handler) http.Handler
}

func NewTeamHandler(svc service.TeamService, authMiddleware func(http.Handler) http.Handler) *TeamHandler {
	return &TeamHandler{service: svc, authMiddleware: authMiddleware}
}

func (h *TeamHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /companies/{id}/teams", h.authMiddleware(http.HandlerFunc(h.ListByCompany)))
}

func (h *TeamHandler) ListByCompany(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	teams, err := h.service.GetByCompanyID(uint(id))
	if errors.Is(err, service.ErrTeamNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, teams)
}
