package repository

import (
	"todo/model"

	"gorm.io/gorm"
)

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
	return &user, result.Error
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, id)
	return &user, result.Error
}
