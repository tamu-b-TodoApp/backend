package service

import (
	"errors"

	"todo/internal/repository"
	"todo/model"
)

var ErrTeamNotFound = errors.New("team not found")

type TeamService interface {
	GetByCompanyID(companyID uint) ([]model.Team, error)
	GetByID(id uint) (*model.Team, error)
}

type teamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{repo: repo}
}

func (s *teamService) GetByCompanyID(companyID uint) ([]model.Team, error) {
	return s.repo.FindByCompanyID(companyID)
}

func (s *teamService) GetByID(id uint) (*model.Team, error) {
	team, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrTeamNotFound
	}
	return team, err
}
