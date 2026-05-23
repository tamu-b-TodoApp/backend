package service

import (
	"crypto/rand"
	"encoding/hex"
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
	Logout(refreshToken string) error
	GetUserByID(id uint) (*model.User, error)
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewAuthService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	return &authService{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	rt, err := s.issueRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	accessToken, err := generateAccessToken(user.ID, rt.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, rt.Token, nil
}

func (s *authService) Refresh(refreshToken string) (string, error) {
	rt, err := s.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		return "", ErrInvalidToken
	}
	if time.Now().After(rt.ExpiresAt) {
		_ = s.refreshTokenRepo.DeleteByToken(refreshToken)
		return "", ErrInvalidToken
	}

	return generateAccessToken(rt.UserID, rt.ID)
}

func (s *authService) Logout(refreshToken string) error {
	return s.refreshTokenRepo.DeleteByToken(refreshToken)
}

func (s *authService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *authService) issueRefreshToken(userID uint) (*model.RefreshToken, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	rt := &model.RefreshToken{
		UserID:    userID,
		Token:     hex.EncodeToString(b),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.refreshTokenRepo.Create(rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func generateAccessToken(userID uint, refreshTokenID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"rti": refreshTokenID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
