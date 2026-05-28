package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"todo/internal/repository"
	"todo/model"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
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
	if err != nil {
		return "", "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := generateToken(user.ID, "access", 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateToken(user.ID, "refresh", 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Refresh(refreshToken string) (string, error) {
	claims, err := parseToken(refreshToken)
	if err != nil {
		return "", ErrInvalidToken
	}
	if claims["type"] != "refresh" {
		return "", ErrInvalidToken
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return "", ErrInvalidToken
	}

	return generateToken(uint(sub), "access", 15*time.Minute)
}

func (s *authService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func generateToken(userID uint, tokenType string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"type": tokenType,
		"exp":  time.Now().Add(duration).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func parseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
