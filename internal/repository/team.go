package repository

import (
	"errors"

	"gorm.io/gorm"
	"todo/model"
)

type TeamRepository interface {
	FindByCompanyID(companyID uint) ([]model.Team, error)
	FindByID(id uint) (*model.Team, error)
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) FindByCompanyID(companyID uint) ([]model.Team, error) {
	var teams []model.Team
	result := r.db.Where("company_id = ?", companyID).Find(&teams)
	return teams, result.Error
}

func (r *teamRepository) FindByID(id uint) (*model.Team, error) {
	var team model.Team
	result := r.db.First(&team, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &team, result.Error
}
