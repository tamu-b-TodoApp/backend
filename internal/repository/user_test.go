package repository_test

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"todo/internal/repository"
	"todo/model"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect test db: " + err.Error())
	}
	if err := db.AutoMigrate(
		&model.User{},
		&model.Company{},
		&model.CompanyMember{},
		&model.Team{},
		&model.TeamMember{},
		&model.Todo{},
	); err != nil {
		panic("failed to migrate: " + err.Error())
	}
	testDB = db
	os.Exit(m.Run())
}

func truncateUsers(t *testing.T) {
	t.Helper()
	testDB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
}

func TestUserRepository_FindByEmail(t *testing.T) {
	truncateUsers(t)
	repo := repository.NewUserRepository(testDB)

	testDB.Create(&model.User{Email: "test@example.com", Password: "hashed"})

	t.Run("found", func(t *testing.T) {
		got, err := repo.FindByEmail("test@example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Email != "test@example.com" {
			t.Errorf("email = %q, want %q", got.Email, "test@example.com")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByEmail("nobody@example.com")
		if err != repository.ErrNotFound {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	truncateUsers(t)
	repo := repository.NewUserRepository(testDB)

	user := &model.User{Email: "test@example.com", Password: "hashed"}
	testDB.Create(user)

	t.Run("found", func(t *testing.T) {
		got, err := repo.FindByID(user.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != user.ID {
			t.Errorf("id = %d, want %d", got.ID, user.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID(99999)
		if err != repository.ErrNotFound {
			t.Errorf("err = %v, want ErrNotFound", err)
		}
	})
}
