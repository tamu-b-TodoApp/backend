package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo/internal/handler"
	"todo/internal/service"
	"todo/model"
)

type mockCompanyService struct {
	getAllFn    func() ([]model.Company, error)
	getByIDFn  func(id uint) (*model.Company, error)
	createFn   func(company *model.Company) error
	updateFn   func(company *model.Company) error
	deleteFn   func(id uint) error
}

func (m *mockCompanyService) GetAll() ([]model.Company, error)          { return m.getAllFn() }
func (m *mockCompanyService) GetByID(id uint) (*model.Company, error)   { return m.getByIDFn(id) }
func (m *mockCompanyService) Create(company *model.Company) error         { return m.createFn(company) }
func (m *mockCompanyService) Update(company *model.Company) error         { return m.updateFn(company) }
func (m *mockCompanyService) Delete(id uint) error                        { return m.deleteFn(id) }

func newCompanyHandler(svc service.CompanyService) *handler.CompanyHandler {
	return handler.NewCompanyHandler(svc, noopMiddleware)
}

func TestCompanyHandler_List(t *testing.T) {
	t.Run("success returns 200 with companies", func(t *testing.T) {
		svc := &mockCompanyService{
			getAllFn: func() ([]model.Company, error) {
				return []model.Company{{ID: 1, Name: "ACME"}}, nil
			},
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodGet, "/companies", nil)
		w := httptest.NewRecorder()
		h.List(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp []model.Company
		json.NewDecoder(w.Body).Decode(&resp)
		if len(resp) != 1 {
			t.Errorf("len = %d, want 1", len(resp))
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockCompanyService{
			getAllFn: func() ([]model.Company, error) { return nil, errors.New("db error") },
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodGet, "/companies", nil)
		w := httptest.NewRecorder()
		h.List(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}

func TestCompanyHandler_Create(t *testing.T) {
	t.Run("success returns 201 with company", func(t *testing.T) {
		svc := &mockCompanyService{
			createFn: func(company *model.Company) error { return nil },
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodPost, "/companies", jsonBody(map[string]string{"name": "ACME"}))
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", w.Code)
		}
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		h := newCompanyHandler(&mockCompanyService{})
		req := httptest.NewRequest(http.MethodPost, "/companies", bytes.NewReader([]byte("bad json")))
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("missing name returns 400", func(t *testing.T) {
		h := newCompanyHandler(&mockCompanyService{})
		req := httptest.NewRequest(http.MethodPost, "/companies", jsonBody(map[string]string{}))
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockCompanyService{
			createFn: func(company *model.Company) error { return errors.New("db error") },
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodPost, "/companies", jsonBody(map[string]string{"name": "ACME"}))
		w := httptest.NewRecorder()
		h.Create(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}

func TestCompanyHandler_GetByID(t *testing.T) {
	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockCompanyService{
			getByIDFn: func(id uint) (*model.Company, error) {
				return &model.Company{ID: id, Name: "ACME"}, nil
			},
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodGet, "/companies/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.GetByID(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockCompanyService{
			getByIDFn: func(id uint) (*model.Company, error) { return nil, service.ErrCompanyNotFound },
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodGet, "/companies/99", nil)
		req.SetPathValue("id", "99")
		w := httptest.NewRecorder()
		h.GetByID(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		h := newCompanyHandler(&mockCompanyService{})
		req := httptest.NewRequest(http.MethodGet, "/companies/abc", nil)
		req.SetPathValue("id", "abc")
		w := httptest.NewRecorder()
		h.GetByID(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})
}

func TestCompanyHandler_Update(t *testing.T) {
	t.Run("success returns 200", func(t *testing.T) {
		svc := &mockCompanyService{
			getByIDFn: func(id uint) (*model.Company, error) {
				return &model.Company{ID: id, Name: "Old"}, nil
			},
			updateFn: func(company *model.Company) error { return nil },
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodPut, "/companies/1", jsonBody(map[string]string{"name": "New"}))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.Update(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp model.Company
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.Name != "New" {
			t.Errorf("name = %q, want New", resp.Name)
		}
	})

	t.Run("not found returns 404", func(t *testing.T) {
		svc := &mockCompanyService{
			getByIDFn: func(id uint) (*model.Company, error) { return nil, service.ErrCompanyNotFound },
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodPut, "/companies/99", jsonBody(map[string]string{"name": "X"}))
		req.SetPathValue("id", "99")
		w := httptest.NewRecorder()
		h.Update(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("missing name returns 400", func(t *testing.T) {
		svc := &mockCompanyService{
			getByIDFn: func(id uint) (*model.Company, error) {
				return &model.Company{ID: id, Name: "Old"}, nil
			},
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodPut, "/companies/1", jsonBody(map[string]string{}))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.Update(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})
}

func TestCompanyHandler_Delete(t *testing.T) {
	t.Run("success returns 204", func(t *testing.T) {
		svc := &mockCompanyService{
			deleteFn: func(id uint) error { return nil },
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodDelete, "/companies/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.Delete(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204", w.Code)
		}
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		h := newCompanyHandler(&mockCompanyService{})
		req := httptest.NewRequest(http.MethodDelete, "/companies/abc", nil)
		req.SetPathValue("id", "abc")
		w := httptest.NewRecorder()
		h.Delete(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockCompanyService{
			deleteFn: func(id uint) error { return errors.New("db error") },
		}
		h := newCompanyHandler(svc)
		req := httptest.NewRequest(http.MethodDelete, "/companies/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		h.Delete(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}
