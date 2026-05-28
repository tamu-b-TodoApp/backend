package repository_test

import (
	"testing"

	"todo/internal/repository"
	"todo/model"
)

func truncateTeams(t *testing.T) {
	t.Helper()
	testDB.Exec("TRUNCATE TABLE teams RESTART IDENTITY CASCADE")
	testDB.Exec("TRUNCATE TABLE companies RESTART IDENTITY CASCADE")
}

func TestTeamRepository_FindByCompanyID(t *testing.T) {
	truncateTeams(t)
	repo := repository.NewTeamRepository(testDB)

	company := &model.Company{Name: "ACME"}
	testDB.Create(company)
	testDB.Create(&model.Team{CompanyID: company.ID, Name: "Alpha"})
	testDB.Create(&model.Team{CompanyID: company.ID, Name: "Beta"})

	t.Run("returns teams for company", func(t *testing.T) {
		teams, err := repo.FindByCompanyID(company.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(teams) != 2 {
			t.Errorf("len = %d, want 2", len(teams))
		}
	})

	t.Run("returns empty slice for unknown company", func(t *testing.T) {
		teams, err := repo.FindByCompanyID(99999)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(teams) != 0 {
			t.Errorf("len = %d, want 0", len(teams))
		}
	})
}

func TestTeamRepository_FindByID(t *testing.T) {
	truncateTeams(t)
	repo := repository.NewTeamRepository(testDB)

	company := &model.Company{Name: "ACME"}
	testDB.Create(company)
	team := &model.Team{CompanyID: company.ID, Name: "Alpha"}
	testDB.Create(team)

	t.Run("found", func(t *testing.T) {
		got, err := repo.FindByID(team.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != team.ID {
			t.Errorf("id = %d, want %d", got.ID, team.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID(99999)
		if err != repository.ErrNotFound {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}
