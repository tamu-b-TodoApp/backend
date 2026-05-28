package repository

import (
	"errors"

	"todo/model"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &user, result.Error
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &user, result.Error
}
