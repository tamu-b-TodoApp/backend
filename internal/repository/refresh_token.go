package repository

import (
	"todo/model"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *model.RefreshToken) error
	FindByID(id uint) (*model.RefreshToken, error)
	FindByToken(token string) (*model.RefreshToken, error)
	DeleteByToken(token string) error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) FindByID(id uint) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	result := r.db.First(&rt, id)
	return &rt, result.Error
}

func (r *refreshTokenRepository) FindByToken(token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	result := r.db.Where("token = ?", token).First(&rt)
	return &rt, result.Error
}

func (r *refreshTokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&model.RefreshToken{}).Error
}
