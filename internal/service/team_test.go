package service_test

import (
	"errors"
	"testing"

	"todo/internal/repository"
	"todo/internal/service"
	"todo/model"
)

type mockTeamRepo struct {
	findByCompanyIDFn func(companyID uint) ([]model.Team, error)
	findByIDFn        func(id uint) (*model.Team, error)
}

func (m *mockTeamRepo) FindByCompanyID(companyID uint) ([]model.Team, error) {
	return m.findByCompanyIDFn(companyID)
}

func (m *mockTeamRepo) FindByID(id uint) (*model.Team, error) {
	return m.findByIDFn(id)
}

func TestTeamService_GetByCompanyID(t *testing.T) {
	t.Run("returns teams", func(t *testing.T) {
		svc := service.NewTeamService(&mockTeamRepo{
			findByCompanyIDFn: func(companyID uint) ([]model.Team, error) {
				return []model.Team{{ID: 1, CompanyID: companyID, Name: "Dev"}}, nil
			},
		})
		teams, err := svc.GetByCompanyID(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(teams) != 1 {
			t.Errorf("len = %d, want 1", len(teams))
		}
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		dbErr := errors.New("db error")
		svc := service.NewTeamService(&mockTeamRepo{
			findByCompanyIDFn: func(companyID uint) ([]model.Team, error) {
				return nil, dbErr
			},
		})
		_, err := svc.GetByCompanyID(1)
		if !errors.Is(err, dbErr) {
			t.Errorf("err = %v, want %v", err, dbErr)
		}
	})
}

func TestTeamService_GetByID(t *testing.T) {
	t.Run("found returns team", func(t *testing.T) {
		svc := service.NewTeamService(&mockTeamRepo{
			findByIDFn: func(id uint) (*model.Team, error) {
				return &model.Team{ID: id, Name: "Dev"}, nil
			},
		})
		team, err := svc.GetByID(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if team.ID != 1 {
			t.Errorf("id = %d, want 1", team.ID)
		}
	})

	t.Run("not found returns ErrTeamNotFound", func(t *testing.T) {
		svc := service.NewTeamService(&mockTeamRepo{
			findByIDFn: func(id uint) (*model.Team, error) {
				return nil, repository.ErrNotFound
			},
		})
		_, err := svc.GetByID(99)
		if !errors.Is(err, service.ErrTeamNotFound) {
			t.Errorf("err = %v, want ErrTeamNotFound", err)
		}
	})
}
