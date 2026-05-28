package service

import (
	"errors"

	"todo/internal/repository"
	"todo/model"
)

var (
	ErrTodoNotFound          = errors.New("todo not found")
	ErrAssigneeNotTeamMember = errors.New("assignee is not a team member")
)

type TodoService interface {
	GetByTeamID(teamID uint) ([]model.Todo, error)
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

func (s *todoService) GetByTeamID(teamID uint) ([]model.Todo, error) {
	return s.repo.FindByTeamID(teamID)
}

func (s *todoService) GetByID(id uint) (*model.Todo, error) {
	todo, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTodoNotFound
	}
	return todo, err
}

func (s *todoService) Create(todo *model.Todo) error {
	if todo.AssigneeID != nil {
		ok, err := s.repo.IsTeamMember(todo.TeamID, *todo.AssigneeID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrAssigneeNotTeamMember
		}
	}
	return s.repo.Create(todo)
}

func (s *todoService) Update(todo *model.Todo) error {
	if todo.AssigneeID != nil {
		ok, err := s.repo.IsTeamMember(todo.TeamID, *todo.AssigneeID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrAssigneeNotTeamMember
		}
	}
	return s.repo.Update(todo)
}

func (s *todoService) Delete(id uint) error {
	return s.repo.Delete(id)
}
