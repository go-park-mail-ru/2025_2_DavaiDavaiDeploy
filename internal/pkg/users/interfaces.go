package users

import (
	"context"
	"kinopoisk/internal/models"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type UsersUsecase interface {
	GenerateToken(id uuid.UUID, login string) (string, error)
	ParseToken(token string) (*jwt.Token, error)
	GetUser(ctx context.Context, id uuid.UUID) (models.User, error)
	ValidateAndGetUser(ctx context.Context, token string) (models.User, error)
	ChangePassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) (models.User, string, error)
	ChangeUserAvatar(ctx context.Context, userID uuid.UUID, fileBytes []byte, fileFormat string) (models.User, string, error)
	CreateFeedback(ctx context.Context, feedback *models.SupportFeedback) error
	GetFeedbackByID(ctx context.Context, id uuid.UUID) (models.SupportFeedback, error)
	GetFeedbacksByUserID(ctx context.Context, userID uuid.UUID) ([]models.SupportFeedback, error)
	UpdateFeedback(ctx context.Context, feedback *models.SupportFeedback) error
	GetFeedbackStats(ctx context.Context) (models.FeedbackStats, error)
	GetUserFeedbackStats(ctx context.Context, userID uuid.UUID) (models.FeedbackStats, error)
	GetAllFeedbacks(ctx context.Context) ([]models.SupportFeedback, error)
}

type UsersRepo interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	UpdateUserPassword(ctx context.Context, version int, userID uuid.UUID, passwordHash []byte) error
	UpdateUserAvatar(ctx context.Context, version int, userID uuid.UUID, avatarPath string) error
	CreateFeedback(ctx context.Context, feedback *models.SupportFeedback) error
	GetFeedbackByID(ctx context.Context, id uuid.UUID) (models.SupportFeedback, error)
	GetFeedbacksByUserID(ctx context.Context, userID uuid.UUID) ([]models.SupportFeedback, error)
	UpdateFeedback(ctx context.Context, feedback *models.SupportFeedback) error
	GetFeedbackStats(ctx context.Context) (models.FeedbackStats, error)
	GetUserFeedbackStats(ctx context.Context, userID uuid.UUID) (models.FeedbackStats, error)
	GetAllFeedbacks(ctx context.Context) ([]models.SupportFeedback, error)
}

type StorageRepo interface {
	DeleteAvatar(ctx context.Context, avatarPath string) error
	UploadAvatar(ctx context.Context, userID string, buffer []byte, fileFormat string, avatarExtension string) (string, error)
}
