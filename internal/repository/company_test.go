package repository_test

import (
	"testing"

	"todo/internal/repository"
	"todo/model"
)

func truncateCompanies(t *testing.T) {
	t.Helper()
	testDB.Exec("TRUNCATE TABLE todos RESTART IDENTITY CASCADE")
	testDB.Exec("TRUNCATE TABLE team_members RESTART IDENTITY CASCADE")
	testDB.Exec("TRUNCATE TABLE teams RESTART IDENTITY CASCADE")
	testDB.Exec("TRUNCATE TABLE company_members RESTART IDENTITY CASCADE")
	testDB.Exec("TRUNCATE TABLE companies RESTART IDENTITY CASCADE")
}

func TestCompanyRepository_FindAll(t *testing.T) {
	truncateCompanies(t)
	repo := repository.NewCompanyRepository(testDB)

	testDB.Create(&model.Company{Name: "Alpha"})
	testDB.Create(&model.Company{Name: "Beta"})

	t.Run("returns all companies", func(t *testing.T) {
		companies, err := repo.FindAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(companies) != 2 {
			t.Errorf("len = %d, want 2", len(companies))
		}
	})
}

func TestCompanyRepository_FindByID(t *testing.T) {
	truncateCompanies(t)
	repo := repository.NewCompanyRepository(testDB)

	company := &model.Company{Name: "ACME"}
	testDB.Create(company)

	t.Run("found", func(t *testing.T) {
		got, err := repo.FindByID(company.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "ACME" {
			t.Errorf("name = %q, want ACME", got.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID(99999)
		if err != repository.ErrNotFound {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}

func TestCompanyRepository_CreateUpdateDelete(t *testing.T) {
	truncateCompanies(t)
	repo := repository.NewCompanyRepository(testDB)

	t.Run("create and find", func(t *testing.T) {
		company := &model.Company{Name: "New Co"}
		if err := repo.Create(company); err != nil {
			t.Fatalf("create: %v", err)
		}
		got, err := repo.FindByID(company.ID)
		if err != nil {
			t.Fatalf("find: %v", err)
		}
		if got.Name != "New Co" {
			t.Errorf("name = %q, want New Co", got.Name)
		}
	})

	t.Run("update", func(t *testing.T) {
		company := &model.Company{Name: "Old Name"}
		testDB.Create(company)
		company.Name = "New Name"
		if err := repo.Update(company); err != nil {
			t.Fatalf("update: %v", err)
		}
		got, _ := repo.FindByID(company.ID)
		if got.Name != "New Name" {
			t.Errorf("name = %q, want New Name", got.Name)
		}
	})

	t.Run("delete", func(t *testing.T) {
		company := &model.Company{Name: "To Delete"}
		testDB.Create(company)
		if err := repo.Delete(company.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}
		_, err := repo.FindByID(company.ID)
		if err != repository.ErrNotFound {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}
