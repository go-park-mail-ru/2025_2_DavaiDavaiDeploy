package service

import (
	"errors"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/repo"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type AuthService struct {
	repository map[string]models.User
	secret     string
}

func NewAuthService(secret string) *AuthService {

	return &AuthService{
		repository: repo.Users,
		secret:     secret,
	}
}

func (s *AuthService) GetUser(login string) (models.User, error) {
	var neededUser models.User
	for i, user := range repo.Users {
		if user.Login == login {
			neededUser = repo.Users[i]
			break
		}
	}
	if neededUser.ID == uuid.Nil {
		return models.User{}, errors.New("no such user")
	}
	return neededUser, nil
}

func (s *AuthService) GenerateToken(login string) (string, error) {
	user, err := s.GetUser(login)
	if err != nil {
		return "", err
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
