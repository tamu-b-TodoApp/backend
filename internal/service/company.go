package service

import (
	"errors"

	"todo/internal/repository"
	"todo/model"
)

var ErrCompanyNotFound = errors.New("company not found")

type CompanyService interface {
	GetAll() ([]model.Company, error)
	GetByID(id uint) (*model.Company, error)
	Create(company *model.Company) error
	Update(company *model.Company) error
	Delete(id uint) error
}

type companyService struct {
	repo repository.CompanyRepository
}

func NewCompanyService(repo repository.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) GetAll() ([]model.Company, error) {
	return s.repo.FindAll()
}

func (s *companyService) GetByID(id uint) (*model.Company, error) {
	company, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrCompanyNotFound
	}
	return company, err
}

func (s *companyService) Create(company *model.Company) error {
	return s.repo.Create(company)
}

func (s *companyService) Update(company *model.Company) error {
	return s.repo.Update(company)
}

func (s *companyService) Delete(id uint) error {
	return s.repo.Delete(id)
}
