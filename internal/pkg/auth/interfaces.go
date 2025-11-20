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
	Enable2FA(ctx context.Context, userID uuid.UUID, has2FA bool) (models.EnableTwoFactorResponse, error)
	Disable2FA(ctx context.Context, userID uuid.UUID, has2FA bool) (models.DisableTwoFactorResponse, error)
	GenerateQRCode(login string) ([]byte, string, error)
	VerifyOTPCode(ctx context.Context, login, secretCode string) error
}

type AuthRepo interface {
	CheckUserExists(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, user models.User) error
	CheckUserLogin(ctx context.Context, login string) (models.User, error)
	IncrementUserVersion(ctx context.Context, userID uuid.UUID) error
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	Enable2FA(ctx context.Context, id uuid.UUID, secret string) (models.EnableTwoFactorResponse, error)
	Disable2FA(ctx context.Context, id uuid.UUID) (models.DisableTwoFactorResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetUserSecretCode(ctx context.Context, userID uuid.UUID) string
}
