package auth

import (
	"context"
	"kinopoisk/internal/models"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type AuthUsecase interface {
	GenerateToken(id uuid.UUID, login string) (string, error)
	ParseToken(token string) (*jwt.Token, error)
	SignUpUser(ctx context.Context, req models.SignUpInput) (models.User, string, error)
	SignInUser(ctx context.Context, req models.SignInInput) (models.User, string, error)
	LogOutUser(ctx context.Context, userID uuid.UUID) error
	ValidateAndGetUser(ctx context.Context, token string) (models.User, error)
}

type AuthRepo interface {
	CheckUserExists(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, user models.User) error
	CheckUserLogin(ctx context.Context, login string) (models.User, error)
	IncrementUserVersion(ctx context.Context, userID uuid.UUID) error
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	Enable2FA(ctx context.Context, id uuid.UUID) (models.EnableTwoFactorResponse, error)
	Disable2FA(ctx context.Context, id uuid.UUID) (models.DisableTwoFactorResponse, error)
}
