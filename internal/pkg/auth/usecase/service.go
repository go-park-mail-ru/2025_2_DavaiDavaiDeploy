package usecase

import (
	"context"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/repo"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	secret   string
	authRepo *repo.AuthRepository
}

func NewAuthService(secret string, authRepo *repo.AuthRepository) *AuthService {
	return &AuthService{
		secret:   secret,
		authRepo: authRepo,
	}
}

func (s *AuthService) GetUser(login string) (models.User, error) {
	user, err := s.authRepo.GetUserByLogin(context.Background(), login)
	if err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (s *AuthService) GenerateToken(login string) (string, error) {
	user, err := s.GetUser(login)
	if err != nil {
		return "suslik", nil
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"login": user.Login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(s.secret))
}

func (s *AuthService) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
}
