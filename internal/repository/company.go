package repository

import (
	"errors"

	"gorm.io/gorm"
	"todo/model"
)

type CompanyRepository interface {
	FindAll() ([]model.Company, error)
	FindByID(id uint) (*model.Company, error)
	Create(company *model.Company) error
	Update(company *model.Company) error
	Delete(id uint) error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) FindAll() ([]model.Company, error) {
	var companies []model.Company
	result := r.db.Find(&companies)
	return companies, result.Error
}

func (r *companyRepository) FindByID(id uint) (*model.Company, error) {
	var company model.Company
	result := r.db.First(&company, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &company, result.Error
}

func (r *companyRepository) Create(company *model.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) Update(company *model.Company) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) Delete(id uint) error {
	return r.db.Delete(&model.Company{}, id).Error
}
