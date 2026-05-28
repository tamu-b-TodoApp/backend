package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo/internal/handler"
	"todo/internal/service"
	"todo/model"
)

type mockTeamService struct {
	getByCompanyIDFn func(companyID uint) ([]model.Team, error)
	getByIDFn        func(id uint) (*model.Team, error)
}

func (m *mockTeamService) GetByCompanyID(companyID uint) ([]model.Team, error) {
	return m.getByCompanyIDFn(companyID)
}

func (m *mockTeamService) GetByID(id uint) (*model.Team, error) {
	return m.getByIDFn(id)
}

func TestTeamHandler_ListByCompany(t *testing.T) {
	t.Run("success returns 200 with teams", func(t *testing.T) {
		svc := &mockTeamService{
			getByCompanyIDFn: func(companyID uint) ([]model.Team, error) {
				return []model.Team{{ID: 1, CompanyID: companyID, Name: "Dev"}}, nil
			},
		}
		h := handler.NewTeamHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodGet, "/companies/1/teams", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.ListByCompany(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		h := handler.NewTeamHandler(&mockTeamService{}, noopMiddleware)
		req := httptest.NewRequest(http.MethodGet, "/companies/abc/teams", nil)
		req.SetPathValue("id", "abc")
		w := httptest.NewRecorder()
		h.ListByCompany(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockTeamService{
			getByCompanyIDFn: func(companyID uint) ([]model.Team, error) {
				return nil, errors.New("db error")
			},
		}
		h := handler.NewTeamHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodGet, "/companies/1/teams", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.ListByCompany(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})

	t.Run("empty result returns 200 with empty array", func(t *testing.T) {
		svc := &mockTeamService{
			getByCompanyIDFn: func(companyID uint) ([]model.Team, error) {
				return []model.Team{}, nil
			},
		}
		h := handler.NewTeamHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodGet, "/companies/99/teams", nil)
		req.SetPathValue("id", "99")
		w := httptest.NewRecorder()
		h.ListByCompany(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", w.Code)
		}
	})

	t.Run("ErrTeamNotFound returns 404", func(t *testing.T) {
		svc := &mockTeamService{
			getByCompanyIDFn: func(companyID uint) ([]model.Team, error) {
				return nil, service.ErrTeamNotFound
			},
		}
		h := handler.NewTeamHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodGet, "/companies/1/teams", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.ListByCompany(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})
}
