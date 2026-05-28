package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"todo/internal/service"
	"todo/model"
)

type TodoHandler struct {
	service        service.TodoService
	teamService    service.TeamService
	authMiddleware func(http.Handler) http.Handler
}

func NewTodoHandler(svc service.TodoService, teamSvc service.TeamService, authMiddleware func(http.Handler) http.Handler) *TodoHandler {
	return &TodoHandler{service: svc, teamService: teamSvc, authMiddleware: authMiddleware}
}

func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /companies/{companyID}/teams/{teamID}/todos", h.authMiddleware(http.HandlerFunc(h.ListByTeam)))
	mux.Handle("POST /companies/{companyID}/teams/{teamID}/todos", h.authMiddleware(http.HandlerFunc(h.Create)))
	mux.Handle("GET /todos/{id}", h.authMiddleware(http.HandlerFunc(h.GetByID)))
	mux.Handle("PUT /todos/{id}", h.authMiddleware(http.HandlerFunc(h.Update)))
	mux.Handle("DELETE /todos/{id}", h.authMiddleware(http.HandlerFunc(h.Delete)))
}

type todoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ParentID    *uint      `json:"parent_id"`
	AssigneeID  *uint      `json:"assignee_id"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	StoryPoints *int       `json:"story_points"`
}

func (r todoRequest) validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *TodoHandler) resolveTeam(w http.ResponseWriter, r *http.Request) (companyID, teamID uint, ok bool) {
	cid, err := strconv.ParseUint(r.PathValue("companyID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid company id", http.StatusBadRequest)
		return 0, 0, false
	}
	tid, err := strconv.ParseUint(r.PathValue("teamID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid team id", http.StatusBadRequest)
		return 0, 0, false
	}

	team, err := h.teamService.GetByID(uint(tid))
	if errors.Is(err, service.ErrTeamNotFound) || (err == nil && team.CompanyID != uint(cid)) {
		http.Error(w, "not found", http.StatusNotFound)
		return 0, 0, false
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return 0, 0, false
	}

	return uint(cid), uint(tid), true
}

func (h *TodoHandler) ListByTeam(w http.ResponseWriter, r *http.Request) {
	_, teamID, ok := h.resolveTeam(w, r)
	if !ok {
		return
	}

	todos, err := h.service.GetByTeamID(teamID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	_, teamID, ok := h.resolveTeam(w, r)
	if !ok {
		return
	}

	var req todoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := req.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := req.Status
	if status == "" {
		status = "not_started"
	}

	todo := &model.Todo{
		TeamID:      teamID,
		ParentID:    req.ParentID,
		AssigneeID:  req.AssigneeID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		DueDate:     req.DueDate,
		StoryPoints: req.StoryPoints,
	}
	if err := h.service.Create(todo); errors.Is(err, service.ErrAssigneeNotTeamMember) {
		http.Error(w, "assignee is not a team member", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo, err := h.service.GetByID(uint(id))
	if errors.Is(err, service.ErrTodoNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo, err := h.service.GetByID(uint(id))
	if errors.Is(err, service.ErrTodoNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var req todoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := req.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.Title = req.Title
	todo.Description = req.Description
	todo.ParentID = req.ParentID
	todo.AssigneeID = req.AssigneeID
	todo.Status = req.Status
	todo.DueDate = req.DueDate
	todo.StoryPoints = req.StoryPoints

	if err := h.service.Update(todo); errors.Is(err, service.ErrAssigneeNotTeamMember) {
		http.Error(w, "assignee is not a team member", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
