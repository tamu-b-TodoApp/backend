package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"todo/internal/middleware"
	"todo/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/logout", h.Logout)
	mux.Handle("GET /auth/me", middleware.Auth(http.HandlerFunc(h.Me)))
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	user.Password = ""
	writeJSON(w, http.StatusOK, user)
}
