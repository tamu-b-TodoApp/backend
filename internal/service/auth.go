package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"todo/internal/repository"
	"todo/internal/token"
	"todo/model"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	Login(email, password string) (accessToken, refreshToken string, err error)
	Refresh(refreshToken string) (accessToken string, err error)
	GetUserByID(id uint) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if errors.Is(err, repository.ErrNotFound) {
		return "", "", ErrInvalidCredentials
	}
	if err != nil {
		return "", "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := generateToken(user.ID, "access", accessTokenDuration)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateToken(user.ID, "refresh", refreshTokenDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Refresh(refreshToken string) (string, error) {
	claims, err := token.Parse(refreshToken)
	if err != nil {
		return "", ErrInvalidToken
	}
	if claims.Type != "refresh" {
		return "", ErrInvalidToken
	}

	return generateToken(claims.UserID, "access", accessTokenDuration)
}

func (s *authService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	return user, err
}

func generateToken(userID uint, tokenType string, duration time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &token.Claims{
		UserID: userID,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	})
	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
