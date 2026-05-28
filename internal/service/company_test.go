package service_test

import (
	"errors"
	"testing"

	"todo/internal/repository"
	"todo/internal/service"
	"todo/model"
)

type mockCompanyRepo struct {
	findAllFn  func() ([]model.Company, error)
	findByIDFn func(id uint) (*model.Company, error)
	createFn   func(company *model.Company) error
	updateFn   func(company *model.Company) error
	deleteFn   func(id uint) error
}

func (m *mockCompanyRepo) FindAll() ([]model.Company, error)          { return m.findAllFn() }
func (m *mockCompanyRepo) FindByID(id uint) (*model.Company, error)   { return m.findByIDFn(id) }
func (m *mockCompanyRepo) Create(company *model.Company) error         { return m.createFn(company) }
func (m *mockCompanyRepo) Update(company *model.Company) error         { return m.updateFn(company) }
func (m *mockCompanyRepo) Delete(id uint) error                        { return m.deleteFn(id) }

func TestCompanyService_GetAll(t *testing.T) {
	t.Run("returns companies", func(t *testing.T) {
		svc := service.NewCompanyService(&mockCompanyRepo{
			findAllFn: func() ([]model.Company, error) {
				return []model.Company{{ID: 1, Name: "ACME"}}, nil
			},
		})
		companies, err := svc.GetAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(companies) != 1 {
			t.Errorf("len = %d, want 1", len(companies))
		}
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		dbErr := errors.New("db error")
		svc := service.NewCompanyService(&mockCompanyRepo{
			findAllFn: func() ([]model.Company, error) { return nil, dbErr },
		})
		_, err := svc.GetAll()
		if !errors.Is(err, dbErr) {
			t.Errorf("err = %v, want %v", err, dbErr)
		}
	})
}

func TestCompanyService_GetByID(t *testing.T) {
	t.Run("found returns company", func(t *testing.T) {
		svc := service.NewCompanyService(&mockCompanyRepo{
			findByIDFn: func(id uint) (*model.Company, error) {
				return &model.Company{ID: id, Name: "ACME"}, nil
			},
		})
		company, err := svc.GetByID(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if company.ID != 1 {
			t.Errorf("id = %d, want 1", company.ID)
		}
	})

	t.Run("not found returns ErrCompanyNotFound", func(t *testing.T) {
		svc := service.NewCompanyService(&mockCompanyRepo{
			findByIDFn: func(id uint) (*model.Company, error) {
				return nil, repository.ErrNotFound
			},
		})
		_, err := svc.GetByID(99)
		if !errors.Is(err, service.ErrCompanyNotFound) {
			t.Errorf("err = %v, want ErrCompanyNotFound", err)
		}
	})
}

func TestCompanyService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := service.NewCompanyService(&mockCompanyRepo{
			createFn: func(company *model.Company) error { return nil },
		})
		if err := svc.Create(&model.Company{Name: "ACME"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		dbErr := errors.New("db error")
		svc := service.NewCompanyService(&mockCompanyRepo{
			createFn: func(company *model.Company) error { return dbErr },
		})
		if err := svc.Create(&model.Company{Name: "ACME"}); !errors.Is(err, dbErr) {
			t.Errorf("err = %v, want %v", err, dbErr)
		}
	})
}

func TestCompanyService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := service.NewCompanyService(&mockCompanyRepo{
			updateFn: func(company *model.Company) error { return nil },
		})
		if err := svc.Update(&model.Company{ID: 1, Name: "Updated"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCompanyService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := service.NewCompanyService(&mockCompanyRepo{
			deleteFn: func(id uint) error { return nil },
		})
		if err := svc.Delete(1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		dbErr := errors.New("db error")
		svc := service.NewCompanyService(&mockCompanyRepo{
			deleteFn: func(id uint) error { return dbErr },
		})
		if err := svc.Delete(1); !errors.Is(err, dbErr) {
			t.Errorf("err = %v, want %v", err, dbErr)
		}
	})
}
