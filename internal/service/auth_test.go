package service_test

import (
	"errors"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"todo/internal/repository"
	"todo/internal/service"
	"todo/model"
)

type mockUserRepo struct {
	findByEmailFn func(email string) (*model.User, error)
	findByIDFn    func(id uint) (*model.User, error)
}

func (m *mockUserRepo) FindByEmail(email string) (*model.User, error) {
	return m.findByEmailFn(email)
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
	return m.findByIDFn(id)
}

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests")
	os.Exit(m.Run())
}

func hashPass(t *testing.T, pass string) string {
	t.Helper()
	b, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestAuthService_Login(t *testing.T) {
	t.Run("success returns access and refresh tokens", func(t *testing.T) {
		hash := hashPass(t, "pass123")
		svc := service.NewAuthService(&mockUserRepo{
			findByEmailFn: func(email string) (*model.User, error) {
				return &model.User{ID: 1, Email: email, Password: hash}, nil
			},
		})
		access, refresh, err := svc.Login("u@e.com", "pass123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if access == "" || refresh == "" {
			t.Error("expected non-empty tokens")
		}
	})

	t.Run("user not found returns ErrInvalidCredentials", func(t *testing.T) {
		svc := service.NewAuthService(&mockUserRepo{
			findByEmailFn: func(email string) (*model.User, error) {
				return nil, repository.ErrNotFound
			},
		})
		_, _, err := svc.Login("nobody@e.com", "pass")
		if !errors.Is(err, service.ErrInvalidCredentials) {
			t.Errorf("err = %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("wrong password returns ErrInvalidCredentials", func(t *testing.T) {
		hash := hashPass(t, "correct")
		svc := service.NewAuthService(&mockUserRepo{
			findByEmailFn: func(email string) (*model.User, error) {
				return &model.User{ID: 1, Email: email, Password: hash}, nil
			},
		})
		_, _, err := svc.Login("u@e.com", "wrong")
		if !errors.Is(err, service.ErrInvalidCredentials) {
			t.Errorf("err = %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		dbErr := errors.New("db connection error")
		svc := service.NewAuthService(&mockUserRepo{
			findByEmailFn: func(email string) (*model.User, error) {
				return nil, dbErr
			},
		})
		_, _, err := svc.Login("u@e.com", "pass")
		if !errors.Is(err, dbErr) {
			t.Errorf("err = %v, want %v", err, dbErr)
		}
	})
}

func TestAuthService_Refresh(t *testing.T) {
	hash := hashPass(t, "pass")
	svc := service.NewAuthService(&mockUserRepo{
		findByEmailFn: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, Password: hash}, nil
		},
	})
	_, refreshToken, err := svc.Login("u@e.com", "pass")
	if err != nil {
		t.Fatalf("setup login failed: %v", err)
	}

	t.Run("success returns new access token", func(t *testing.T) {
		access, err := svc.Refresh(refreshToken)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if access == "" {
			t.Error("expected non-empty access token")
		}
	})

	t.Run("invalid token string returns ErrInvalidToken", func(t *testing.T) {
		_, err := svc.Refresh("not-a-valid-token")
		if !errors.Is(err, service.ErrInvalidToken) {
			t.Errorf("err = %v, want ErrInvalidToken", err)
		}
	})

	t.Run("access token used as refresh token returns ErrInvalidToken", func(t *testing.T) {
		accessToken, _, err := svc.Login("u@e.com", "pass")
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}
		_, err = svc.Refresh(accessToken)
		if !errors.Is(err, service.ErrInvalidToken) {
			t.Errorf("err = %v, want ErrInvalidToken", err)
		}
	})
}

func TestAuthService_GetUserByID(t *testing.T) {
	t.Run("found returns user", func(t *testing.T) {
		svc := service.NewAuthService(&mockUserRepo{
			findByIDFn: func(id uint) (*model.User, error) {
				return &model.User{ID: id, Email: "u@e.com"}, nil
			},
		})
		user, err := svc.GetUserByID(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != 1 {
			t.Errorf("id = %d, want 1", user.ID)
		}
	})

	t.Run("not found returns ErrUserNotFound", func(t *testing.T) {
		svc := service.NewAuthService(&mockUserRepo{
			findByIDFn: func(id uint) (*model.User, error) {
				return nil, repository.ErrNotFound
			},
		})
		_, err := svc.GetUserByID(99)
		if !errors.Is(err, service.ErrUserNotFound) {
			t.Errorf("err = %v, want ErrUserNotFound", err)
		}
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		dbErr := errors.New("db connection error")
		svc := service.NewAuthService(&mockUserRepo{
			findByIDFn: func(id uint) (*model.User, error) {
				return nil, dbErr
			},
		})
		_, err := svc.GetUserByID(1)
		if !errors.Is(err, dbErr) {
			t.Errorf("err = %v, want %v", err, dbErr)
		}
	})
}
