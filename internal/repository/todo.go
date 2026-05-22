package repository

import (
	"todo/model"

	"gorm.io/gorm"
)

type TodoRepository interface {
	FindAll() ([]model.Todo, error)
	FindByID(id uint) (*model.Todo, error)
	Create(todo *model.Todo) error
	Update(todo *model.Todo) error
	Delete(id uint) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll() ([]model.Todo, error) {
	var todos []model.Todo
	result := r.db.Find(&todos)
	return todos, result.Error
}

func (r *todoRepository) FindByID(id uint) (*model.Todo, error) {
	var todo model.Todo
	result := r.db.First(&todo, id)
	return &todo, result.Error
}

func (r *todoRepository) Create(todo *model.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) Update(todo *model.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	return r.db.Delete(&model.Todo{}, id).Error
}
