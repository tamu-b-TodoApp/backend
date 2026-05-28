package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo/internal/handler"
	"todo/internal/middleware"
	"todo/internal/service"
	"todo/model"
)

type mockAuthService struct {
	loginFn       func(email, password string) (string, string, error)
	refreshFn     func(refreshToken string) (string, error)
	getUserByIDFn func(id uint) (*model.User, error)
}

func (m *mockAuthService) Login(email, password string) (string, string, error) {
	return m.loginFn(email, password)
}

func (m *mockAuthService) Refresh(refreshToken string) (string, error) {
	return m.refreshFn(refreshToken)
}

func (m *mockAuthService) GetUserByID(id uint) (*model.User, error) {
	return m.getUserByIDFn(id)
}

func noopMiddleware(next http.Handler) http.Handler { return next }

func jsonBody(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("success returns 200 with tokens", func(t *testing.T) {
		svc := &mockAuthService{
			loginFn: func(email, password string) (string, string, error) {
				return "access_tok", "refresh_tok", nil
			},
		}
		h := handler.NewAuthHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(map[string]string{
			"email": "u@e.com", "password": "pass",
		}))
		w := httptest.NewRecorder()
		h.Login(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.AccessToken != "access_tok" || resp.RefreshToken != "refresh_tok" {
			t.Errorf("tokens = %q/%q, want access_tok/refresh_tok", resp.AccessToken, resp.RefreshToken)
		}
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		h := handler.NewAuthHandler(&mockAuthService{}, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte("bad json")))
		w := httptest.NewRecorder()
		h.Login(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("missing email returns 400", func(t *testing.T) {
		h := handler.NewAuthHandler(&mockAuthService{}, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(map[string]string{"password": "pass"}))
		w := httptest.NewRecorder()
		h.Login(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("missing password returns 400", func(t *testing.T) {
		h := handler.NewAuthHandler(&mockAuthService{}, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(map[string]string{"email": "u@e.com"}))
		w := httptest.NewRecorder()
		h.Login(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("invalid credentials returns 401", func(t *testing.T) {
		svc := &mockAuthService{
			loginFn: func(email, password string) (string, string, error) {
				return "", "", service.ErrInvalidCredentials
			},
		}
		h := handler.NewAuthHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(map[string]string{
			"email": "u@e.com", "password": "wrong",
		}))
		w := httptest.NewRecorder()
		h.Login(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want 401", w.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockAuthService{
			loginFn: func(email, password string) (string, string, error) {
				return "", "", errors.New("unexpected error")
			},
		}
		h := handler.NewAuthHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(map[string]string{
			"email": "u@e.com", "password": "pass",
		}))
		w := httptest.NewRecorder()
		h.Login(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}

func TestAuthHandler_Refresh(t *testing.T) {
	t.Run("success returns 200 with access token", func(t *testing.T) {
		svc := &mockAuthService{
			refreshFn: func(refreshToken string) (string, error) {
				return "new_access_tok", nil
			},
		}
		h := handler.NewAuthHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", jsonBody(map[string]string{
			"refresh_token": "some_refresh_token",
		}))
		w := httptest.NewRecorder()
		h.Refresh(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp struct {
			AccessToken string `json:"access_token"`
		}
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.AccessToken != "new_access_tok" {
			t.Errorf("access_token = %q, want new_access_tok", resp.AccessToken)
		}
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		h := handler.NewAuthHandler(&mockAuthService{}, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader([]byte("bad json")))
		w := httptest.NewRecorder()
		h.Refresh(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("missing refresh_token returns 400", func(t *testing.T) {
		h := handler.NewAuthHandler(&mockAuthService{}, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", jsonBody(map[string]string{}))
		w := httptest.NewRecorder()
		h.Refresh(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		svc := &mockAuthService{
			refreshFn: func(refreshToken string) (string, error) {
				return "", service.ErrInvalidToken
			},
		}
		h := handler.NewAuthHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", jsonBody(map[string]string{
			"refresh_token": "expired_token",
		}))
		w := httptest.NewRecorder()
		h.Refresh(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want 401", w.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockAuthService{
			refreshFn: func(refreshToken string) (string, error) {
				return "", errors.New("unexpected error")
			},
		}
		h := handler.NewAuthHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodPost, "/auth/refresh", jsonBody(map[string]string{
			"refresh_token": "some_token",
		}))
		w := httptest.NewRecorder()
		h.Refresh(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}

func TestAuthHandler_Me(t *testing.T) {
	t.Run("success returns 200 with user", func(t *testing.T) {
		svc := &mockAuthService{
			getUserByIDFn: func(id uint) (*model.User, error) {
				return &model.User{ID: id, Email: "u@e.com"}, nil
			},
		}
		h := handler.NewAuthHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()
		h.Me(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}
		var resp model.User
		json.NewDecoder(w.Body).Decode(&resp)
		if resp.Email != "u@e.com" {
			t.Errorf("email = %q, want u@e.com", resp.Email)
		}
	})

	t.Run("no userID in context returns 401", func(t *testing.T) {
		h := handler.NewAuthHandler(&mockAuthService{}, noopMiddleware)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		w := httptest.NewRecorder()
		h.Me(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want 401", w.Code)
		}
	})

	t.Run("user not found returns 404", func(t *testing.T) {
		svc := &mockAuthService{
			getUserByIDFn: func(id uint) (*model.User, error) {
				return nil, service.ErrUserNotFound
			},
		}
		h := handler.NewAuthHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(99))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()
		h.Me(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &mockAuthService{
			getUserByIDFn: func(id uint) (*model.User, error) {
				return nil, errors.New("unexpected error")
			},
		}
		h := handler.NewAuthHandler(svc, noopMiddleware)
		req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, uint(1))
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()
		h.Me(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500", w.Code)
		}
	})
}
