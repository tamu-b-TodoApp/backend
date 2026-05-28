package repository_test

import (
	"testing"

	"todo/internal/repository"
	"todo/model"
)

func truncateTodos(t *testing.T) {
	t.Helper()
	testDB.Exec("TRUNCATE TABLE todos RESTART IDENTITY CASCADE")
	testDB.Exec("TRUNCATE TABLE team_members RESTART IDENTITY CASCADE")
	testDB.Exec("TRUNCATE TABLE company_members RESTART IDENTITY CASCADE")
	testDB.Exec("TRUNCATE TABLE teams RESTART IDENTITY CASCADE")
	testDB.Exec("TRUNCATE TABLE companies RESTART IDENTITY CASCADE")
}

func setupTodoFixture(t *testing.T) (team *model.Team, member *model.CompanyMember) {
	t.Helper()
	company := &model.Company{Name: "ACME"}
	testDB.Create(company)
	team = &model.Team{CompanyID: company.ID, Name: "Dev"}
	testDB.Create(team)
	member = &model.CompanyMember{CompanyID: company.ID, UserID: 1}
	testDB.Create(member)
	testDB.Create(&model.TeamMember{TeamID: team.ID, CompanyMemberID: member.ID})
	return team, member
}

func TestTodoRepository_FindByTeamID(t *testing.T) {
	truncateTodos(t)
	repo := repository.NewTodoRepository(testDB)
	team, _ := setupTodoFixture(t)

	testDB.Create(&model.Todo{TeamID: team.ID, Title: "Task 1", Description: "", Status: "not_started"})
	testDB.Create(&model.Todo{TeamID: team.ID, Title: "Task 2", Description: "", Status: "not_started"})

	t.Run("returns todos for team", func(t *testing.T) {
		todos, err := repo.FindByTeamID(team.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 2 {
			t.Errorf("len = %d, want 2", len(todos))
		}
	})

	t.Run("returns empty slice for unknown team", func(t *testing.T) {
		todos, err := repo.FindByTeamID(99999)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 0 {
			t.Errorf("len = %d, want 0", len(todos))
		}
	})
}

func TestTodoRepository_FindByID(t *testing.T) {
	truncateTodos(t)
	repo := repository.NewTodoRepository(testDB)
	team, _ := setupTodoFixture(t)

	todo := &model.Todo{TeamID: team.ID, Title: "Task", Description: "", Status: "not_started"}
	testDB.Create(todo)

	t.Run("found", func(t *testing.T) {
		got, err := repo.FindByID(todo.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != todo.ID {
			t.Errorf("id = %d, want %d", got.ID, todo.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID(99999)
		if err != repository.ErrNotFound {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}

func TestTodoRepository_IsTeamMember(t *testing.T) {
	truncateTodos(t)
	repo := repository.NewTodoRepository(testDB)
	team, member := setupTodoFixture(t)

	t.Run("is member", func(t *testing.T) {
		ok, err := repo.IsTeamMember(team.ID, member.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Error("expected true, got false")
		}
	})

	t.Run("is not member", func(t *testing.T) {
		ok, err := repo.IsTeamMember(team.ID, 99999)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Error("expected false, got true")
		}
	})
}

func TestTodoRepository_CreateUpdateDelete(t *testing.T) {
	truncateTodos(t)
	repo := repository.NewTodoRepository(testDB)
	team, _ := setupTodoFixture(t)

	t.Run("create and find", func(t *testing.T) {
		todo := &model.Todo{TeamID: team.ID, Title: "New Task", Description: "desc", Status: "not_started"}
		if err := repo.Create(todo); err != nil {
			t.Fatalf("create: %v", err)
		}
		got, err := repo.FindByID(todo.ID)
		if err != nil {
			t.Fatalf("find: %v", err)
		}
		if got.Title != "New Task" {
			t.Errorf("title = %q, want New Task", got.Title)
		}
	})

	t.Run("update", func(t *testing.T) {
		todo := &model.Todo{TeamID: team.ID, Title: "Old", Description: "", Status: "not_started"}
		testDB.Create(todo)
		todo.Title = "Updated"
		if err := repo.Update(todo); err != nil {
			t.Fatalf("update: %v", err)
		}
		got, _ := repo.FindByID(todo.ID)
		if got.Title != "Updated" {
			t.Errorf("title = %q, want Updated", got.Title)
		}
	})

	t.Run("delete", func(t *testing.T) {
		todo := &model.Todo{TeamID: team.ID, Title: "To Delete", Description: "", Status: "not_started"}
		testDB.Create(todo)
		if err := repo.Delete(todo.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}
		_, err := repo.FindByID(todo.ID)
		if err != repository.ErrNotFound {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}
