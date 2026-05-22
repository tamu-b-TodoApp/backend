package service

import (
	"todo/internal/repository"
	"todo/model"
)

type TodoService interface {
	GetAll() ([]model.Todo, error)
	GetByID(id uint) (*model.Todo, error)
	Create(todo *model.Todo) error
	Update(todo *model.Todo) error
	Delete(id uint) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) GetAll() ([]model.Todo, error) {
	return s.repo.FindAll()
}

func (s *todoService) GetByID(id uint) (*model.Todo, error) {
	return s.repo.FindByID(id)
}

func (s *todoService) Create(todo *model.Todo) error {
	return s.repo.Create(todo)
}

func (s *todoService) Update(todo *model.Todo) error {
	return s.repo.Update(todo)
}

func (s *todoService) Delete(id uint) error {
	return s.repo.Delete(id)
}
