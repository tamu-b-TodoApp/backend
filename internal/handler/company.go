package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"todo/internal/service"
	"todo/model"
)

type CompanyHandler struct {
	service        service.CompanyService
	authMiddleware func(http.Handler) http.Handler
}

func NewCompanyHandler(svc service.CompanyService, authMiddleware func(http.Handler) http.Handler) *CompanyHandler {
	return &CompanyHandler{service: svc, authMiddleware: authMiddleware}
}

func (h *CompanyHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /companies", h.authMiddleware(http.HandlerFunc(h.List)))
	mux.Handle("POST /companies", h.authMiddleware(http.HandlerFunc(h.Create)))
	mux.Handle("GET /companies/{id}", h.authMiddleware(http.HandlerFunc(h.GetByID)))
	mux.Handle("PUT /companies/{id}", h.authMiddleware(http.HandlerFunc(h.Update)))
	mux.Handle("DELETE /companies/{id}", h.authMiddleware(http.HandlerFunc(h.Delete)))
}

type companyRequest struct {
	Name string `json:"name"`
}

func (r companyRequest) validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	companies, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, companies)
}

func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req companyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := req.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	company := &model.Company{Name: req.Name}
	if err := h.service.Create(company); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, company)
}

func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	company, err := h.service.GetByID(uint(id))
	if errors.Is(err, service.ErrCompanyNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, company)
}

func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	company, err := h.service.GetByID(uint(id))
	if errors.Is(err, service.ErrCompanyNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	var req companyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := req.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	company.Name = req.Name
	if err := h.service.Update(company); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, company)
}

func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
