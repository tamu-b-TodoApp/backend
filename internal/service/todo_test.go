package service_test

import (
	"errors"
	"testing"

	"todo/internal/repository"
	"todo/internal/service"
	"todo/model"
)

type mockTodoRepo struct {
	findByTeamIDFn  func(teamID uint) ([]model.Todo, error)
	findByIDFn      func(id uint) (*model.Todo, error)
	createFn        func(todo *model.Todo) error
	updateFn        func(todo *model.Todo) error
	deleteFn        func(id uint) error
	isTeamMemberFn  func(teamID, companyMemberID uint) (bool, error)
}

func (m *mockTodoRepo) FindByTeamID(teamID uint) ([]model.Todo, error) {
	return m.findByTeamIDFn(teamID)
}

func (m *mockTodoRepo) FindByID(id uint) (*model.Todo, error) {
	return m.findByIDFn(id)
}

func (m *mockTodoRepo) Create(todo *model.Todo) error {
	return m.createFn(todo)
}

func (m *mockTodoRepo) Update(todo *model.Todo) error {
	return m.updateFn(todo)
}

func (m *mockTodoRepo) Delete(id uint) error {
	return m.deleteFn(id)
}

func (m *mockTodoRepo) IsTeamMember(teamID, companyMemberID uint) (bool, error) {
	return m.isTeamMemberFn(teamID, companyMemberID)
}

func TestTodoService_GetByTeamID(t *testing.T) {
	t.Run("returns todos", func(t *testing.T) {
		svc := service.NewTodoService(&mockTodoRepo{
			findByTeamIDFn: func(teamID uint) ([]model.Todo, error) {
				return []model.Todo{{ID: 1, TeamID: teamID, Title: "Task"}}, nil
			},
		})
		todos, err := svc.GetByTeamID(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 1 {
			t.Errorf("len = %d, want 1", len(todos))
		}
	})
}

func TestTodoService_GetByID(t *testing.T) {
	t.Run("found returns todo", func(t *testing.T) {
		svc := service.NewTodoService(&mockTodoRepo{
			findByIDFn: func(id uint) (*model.Todo, error) {
				return &model.Todo{ID: id, Title: "Task"}, nil
			},
		})
		todo, err := svc.GetByID(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if todo.ID != 1 {
			t.Errorf("id = %d, want 1", todo.ID)
		}
	})

	t.Run("not found returns ErrTodoNotFound", func(t *testing.T) {
		svc := service.NewTodoService(&mockTodoRepo{
			findByIDFn: func(id uint) (*model.Todo, error) {
				return nil, repository.ErrNotFound
			},
		})
		_, err := svc.GetByID(99)
		if !errors.Is(err, service.ErrTodoNotFound) {
			t.Errorf("err = %v, want ErrTodoNotFound", err)
		}
	})
}

func TestTodoService_Create(t *testing.T) {
	t.Run("success without assignee", func(t *testing.T) {
		svc := service.NewTodoService(&mockTodoRepo{
			createFn: func(todo *model.Todo) error { return nil },
		})
		err := svc.Create(&model.Todo{TeamID: 1, Title: "Task"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("success with valid assignee", func(t *testing.T) {
		assigneeID := uint(10)
		svc := service.NewTodoService(&mockTodoRepo{
			isTeamMemberFn: func(teamID, companyMemberID uint) (bool, error) { return true, nil },
			createFn:       func(todo *model.Todo) error { return nil },
		})
		err := svc.Create(&model.Todo{TeamID: 1, Title: "Task", AssigneeID: &assigneeID})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("assignee not team member returns ErrAssigneeNotTeamMember", func(t *testing.T) {
		assigneeID := uint(99)
		svc := service.NewTodoService(&mockTodoRepo{
			isTeamMemberFn: func(teamID, companyMemberID uint) (bool, error) { return false, nil },
		})
		err := svc.Create(&model.Todo{TeamID: 1, Title: "Task", AssigneeID: &assigneeID})
		if !errors.Is(err, service.ErrAssigneeNotTeamMember) {
			t.Errorf("err = %v, want ErrAssigneeNotTeamMember", err)
		}
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		dbErr := errors.New("db error")
		svc := service.NewTodoService(&mockTodoRepo{
			createFn: func(todo *model.Todo) error { return dbErr },
		})
		err := svc.Create(&model.Todo{TeamID: 1, Title: "Task"})
		if !errors.Is(err, dbErr) {
			t.Errorf("err = %v, want %v", err, dbErr)
		}
	})
}

func TestTodoService_Update(t *testing.T) {
	t.Run("success without assignee", func(t *testing.T) {
		svc := service.NewTodoService(&mockTodoRepo{
			updateFn: func(todo *model.Todo) error { return nil },
		})
		err := svc.Update(&model.Todo{ID: 1, TeamID: 1, Title: "Task"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("assignee not team member returns ErrAssigneeNotTeamMember", func(t *testing.T) {
		assigneeID := uint(99)
		svc := service.NewTodoService(&mockTodoRepo{
			isTeamMemberFn: func(teamID, companyMemberID uint) (bool, error) { return false, nil },
		})
		err := svc.Update(&model.Todo{ID: 1, TeamID: 1, Title: "Task", AssigneeID: &assigneeID})
		if !errors.Is(err, service.ErrAssigneeNotTeamMember) {
			t.Errorf("err = %v, want ErrAssigneeNotTeamMember", err)
		}
	})
}

func TestTodoService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := service.NewTodoService(&mockTodoRepo{
			deleteFn: func(id uint) error { return nil },
		})
		if err := svc.Delete(1); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		dbErr := errors.New("db error")
		svc := service.NewTodoService(&mockTodoRepo{
			deleteFn: func(id uint) error { return dbErr },
		})
		err := svc.Delete(1)
		if !errors.Is(err, dbErr) {
			t.Errorf("err = %v, want %v", err, dbErr)
		}
	})
}
