package repository

import (
	"errors"

	"gorm.io/gorm"
	"todo/model"
)

type TodoRepository interface {
	FindByTeamID(teamID uint) ([]model.Todo, error)
	FindByID(id uint) (*model.Todo, error)
	Create(todo *model.Todo) error
	Update(todo *model.Todo) error
	Delete(id uint) error
	IsTeamMember(teamID, companyMemberID uint) (bool, error)
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindByTeamID(teamID uint) ([]model.Todo, error) {
	var todos []model.Todo
	result := r.db.Where("team_id = ?", teamID).Find(&todos)
	return todos, result.Error
}

func (r *todoRepository) FindByID(id uint) (*model.Todo, error) {
	var todo model.Todo
	result := r.db.First(&todo, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
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

func (r *todoRepository) IsTeamMember(teamID, companyMemberID uint) (bool, error) {
	var count int64
	result := r.db.Model(&model.TeamMember{}).
		Where("team_id = ? AND company_member_id = ?", teamID, companyMemberID).
		Count(&count)
	return count > 0, result.Error
}
