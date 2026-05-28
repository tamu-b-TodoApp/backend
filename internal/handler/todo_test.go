package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo/internal/handler"
	"todo/internal/service"
	"todo/model"
)

type mockTodoService struct {
	getByTeamIDFn func(teamID uint) ([]model.Todo, error)
	getByIDFn     func(id uint) (*model.Todo, error)
	createFn      func(todo *model.Todo) error
	updateFn      func(todo *model.Todo) error
	deleteFn      func(id uint) error
}

func (m *mockTodoService) GetByTeamID(teamID uint) ([]model.Todo, error) {
	return m.getByTeamIDFn(teamID)
}

func (m *mockTodoService) GetByID(id uint) (*model.Todo, error) {
	return m.getByIDFn(id)
}

func (m *mockTodoService) Create(todo *model.Todo) error {
	return m.createFn(todo)
}

func (m *mockTodoService) Update(todo *model.Todo) error {
	return m.updateFn(todo)
}

func (m *mockTodoService) Delete(id uint) error {
	return m.deleteFn(id)
}

func teamInCompany(companyID uint) *mockTeamService {
	return &mockTeamService{
		getByIDFn: func(id uint) (*model.Team, error) {
			return &model.Team{ID: id, CompanyID: companyID}, nil
		},
	}
}

func newTodoHandler(todoSvc service.TodoService, teamSvc service.TeamService) *handler.TodoHandler {
	return handler.NewTodoHandler(todoSvc, teamSvc, noopMiddleware)
}

func TestTodoHandler_ListByTeam(t *testing.T) {
	t.Run("success returns 200 with todos", func(t *testing.T) {
		svc := &mockTodoService{
			getByTeamIDFn: func(teamID uint) ([]model.Todo, error) {
				return []model.Todo{{ID: 1, TeamID: teamID, Title: "Task"}}, nil
			},
		}
		h := newTodoHandler(svc, teamInCompany(1))
		req := httptest.NewRequest(http.MethodGet, "/companies/1/teams/10/todos", nil)
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "10")
		w := httptest.NewRecorder()
		h.ListByTeam(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})

	t.Run("team belongs to different company returns 404", func(t *testing.T) {
		teamSvc := &mockTeamService{
			getByIDFn: func(id uint) (*model.Team, error) {
				return &model.Team{ID: id, CompanyID: 99}, nil
			},
		}
		h := newTodoHandler(&mockTodoService{}, teamSvc)
		req := httptest.NewRequest(http.MethodGet, "/companies/1/teams/10/todos", nil)
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "10")
		w := httptest.NewRecorder()
		h.ListByTeam(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("team not found returns 404", func(t *testing.T) {
		teamSvc := &mockTeamService{
			getByIDFn: func(id uint) (*model.Team, error) {
				return nil, service.ErrTeamNotFound
			},
		}
		h := newTodoHandler(&mockTodoService{}, teamSvc)
		req := httptest.NewRequest(http.MethodGet, "/companies/1/teams/99/todos", nil)
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "99")
		w := httptest.NewRecorder()
		h.ListByTeam(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("invalid companyID returns 400", func(t *testing.T) {
		h := newTodoHandler(&mockTodoService{}, &mockTeamService{})
		req := httptest.NewRequest(http.MethodGet, "/companies/abc/teams/1/todos", nil)
		req.SetPathValue("companyID", "abc")
		req.SetPathValue("teamID", "1")
		w := httptest.NewRecorder()
		h.ListByTeam(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockTodoService{
			getByTeamIDFn: func(teamID uint) ([]model.Todo, error) {
				return nil, errors.New("db error")
			},
		}
		h := newTodoHandler(svc, teamInCompany(1))
		req := httptest.NewRequest(http.MethodGet, "/companies/1/teams/10/todos", nil)
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "10")
		w := httptest.NewRecorder()
		h.ListByTeam(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}

func TestTodoHandler_Create(t *testing.T) {
	t.Run("success returns 201", func(t *testing.T) {
		svc := &mockTodoService{
			createFn: func(todo *model.Todo) error { return nil },
		}
		h := newTodoHandler(svc, teamInCompany(1))
		req := httptest.NewRequest(http.MethodPost, "/companies/1/teams/10/todos", jsonBody(map[string]any{
			"title": "New Task", "description": "desc",
		}))
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "10")
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", w.Code)
		}
	})

	t.Run("team belongs to different company returns 404", func(t *testing.T) {
		teamSvc := &mockTeamService{
			getByIDFn: func(id uint) (*model.Team, error) {
				return &model.Team{ID: id, CompanyID: 99}, nil
			},
		}
		h := newTodoHandler(&mockTodoService{}, teamSvc)
		req := httptest.NewRequest(http.MethodPost, "/companies/1/teams/10/todos", jsonBody(map[string]any{"title": "T"}))
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "10")
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		h := newTodoHandler(&mockTodoService{}, teamInCompany(1))
		req := httptest.NewRequest(http.MethodPost, "/companies/1/teams/10/todos", bytes.NewReader([]byte("bad json")))
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "10")
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("missing title returns 400", func(t *testing.T) {
		h := newTodoHandler(&mockTodoService{}, teamInCompany(1))
		req := httptest.NewRequest(http.MethodPost, "/companies/1/teams/10/todos", jsonBody(map[string]any{"description": "d"}))
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "10")
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("assignee not team member returns 400", func(t *testing.T) {
		svc := &mockTodoService{
			createFn: func(todo *model.Todo) error { return service.ErrAssigneeNotTeamMember },
		}
		h := newTodoHandler(svc, teamInCompany(1))
		req := httptest.NewRequest(http.MethodPost, "/companies/1/teams/10/todos", jsonBody(map[string]any{
			"title": "Task", "assignee_id": 99,
		}))
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "10")
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockTodoService{
			createFn: func(todo *model.Todo) error { return errors.New("db error") },
		}
		h := newTodoHandler(svc, teamInCompany(1))
		req := httptest.NewRequest(http.MethodPost, "/companies/1/teams/10/todos", jsonBody(map[string]any{"title": "Task"}))
		req.SetPathValue("companyID", "1")
		req.SetPathValue("teamID", "10")
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}

func TestTodoHandler_GetByID(t *testing.T) {
	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockTodoService{
			getByIDFn: func(id uint) (*model.Todo, error) {
				return &model.Todo{ID: id, Title: "Task"}, nil
			},
		}
		h := newTodoHandler(svc, &mockTeamService{})
		req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.GetByID(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockTodoService{
			getByIDFn: func(id uint) (*model.Todo, error) { return nil, service.ErrTodoNotFound },
		}
		h := newTodoHandler(svc, &mockTeamService{})
		req := httptest.NewRequest(http.MethodGet, "/todos/99", nil)
		req.SetPathValue("id", "99")
		w := httptest.NewRecorder()
		h.GetByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		h := newTodoHandler(&mockTodoService{}, &mockTeamService{})
		req := httptest.NewRequest(http.MethodGet, "/todos/abc", nil)
		req.SetPathValue("id", "abc")
		w := httptest.NewRecorder()
		h.GetByID(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})
}

func TestTodoHandler_Update(t *testing.T) {
	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockTodoService{
			getByIDFn: func(id uint) (*model.Todo, error) {
				return &model.Todo{ID: id, TeamID: 1, Title: "Old"}, nil
			},
			updateFn: func(todo *model.Todo) error { return nil },
		}
		h := newTodoHandler(svc, &mockTeamService{})
		req := httptest.NewRequest(http.MethodPut, "/todos/1", jsonBody(map[string]any{
			"title": "Updated", "description": "",
		}))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.Update(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})

	t.Run("todo not found returns 404", func(t *testing.T) {
		svc := &mockTodoService{
			getByIDFn: func(id uint) (*model.Todo, error) { return nil, service.ErrTodoNotFound },
		}
		h := newTodoHandler(svc, &mockTeamService{})
		req := httptest.NewRequest(http.MethodPut, "/todos/99", jsonBody(map[string]any{"title": "T"}))
		req.SetPathValue("id", "99")
		w := httptest.NewRecorder()
		h.Update(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("assignee not team member returns 400", func(t *testing.T) {
		svc := &mockTodoService{
			getByIDFn: func(id uint) (*model.Todo, error) {
				return &model.Todo{ID: id, TeamID: 1, Title: "Task"}, nil
			},
			updateFn: func(todo *model.Todo) error { return service.ErrAssigneeNotTeamMember },
		}
		h := newTodoHandler(svc, &mockTeamService{})
		req := httptest.NewRequest(http.MethodPut, "/todos/1", jsonBody(map[string]any{
			"title": "Task", "assignee_id": 99,
		}))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.Update(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})
}

func TestTodoHandler_Delete(t *testing.T) {
	t.Run("success returns 204", func(t *testing.T) {
		svc := &mockTodoService{
			deleteFn: func(id uint) error { return nil },
		}
		h := newTodoHandler(svc, &mockTeamService{})
		req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.Delete(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204", w.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		h := newTodoHandler(&mockTodoService{}, &mockTeamService{})
		req := httptest.NewRequest(http.MethodDelete, "/todos/abc", nil)
		req.SetPathValue("id", "abc")
		w := httptest.NewRecorder()
		h.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockTodoService{
			deleteFn: func(id uint) error { return errors.New("db error") },
		}
		h := newTodoHandler(svc, &mockTeamService{})
		req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.Delete(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}
