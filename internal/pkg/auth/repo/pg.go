package repo

import (
	"context"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"

	"github.com/jackc/pgtype/pgxtype"
	uuid "github.com/satori/go.uuid"
)

type AuthRepository struct {
	db pgxtype.Querier
}

func NewAuthRepository(db pgxtype.Querier) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CheckUserExists(ctx context.Context, login string) (bool, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var exists bool
	err := r.db.QueryRow(
		ctx,
		CheckUserExistsQuery,
		login,
	).Scan(&exists)
	if err != nil {
		logger.Error("failed to scan user: " + err.Error())
		return false, err
	}
	return exists, err
}

func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	_, err := r.db.Exec(
		ctx,
		CreateUserQuery,
		user.ID, user.Login, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		logger.Error("failed to create user: " + err.Error())
	}
	return err
}

func (r *AuthRepository) CheckUserLogin(ctx context.Context, login string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var user models.User
	err := r.db.QueryRow(ctx,
		CheckUserLoginQuery,
		login,
	).Scan(&user.ID,
		&user.Version,
		&user.Login,
		&user.PasswordHash,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		logger.Error("failed to scan user: " + err.Error())
		return models.User{}, err
	}
	return user, nil
}

func (r *AuthRepository) IncrementUserVersion(ctx context.Context, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	_, err := r.db.Exec(
		ctx,
		IncrementUserVersionQuery,
		userID,
	)
	if err != nil {
		logger.Error("failed to increment version: " + err.Error())
	}
	return err
}

func (r *AuthRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var user models.User
	err := r.db.QueryRow(
		ctx,
		GetUserByLoginQuery,
		login,
	).Scan(
		&user.ID, &user.Version, &user.Login,
		&user.PasswordHash, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		logger.Error("failed to scan user: " + err.Error())
		return models.User{}, err
	}
	return user, nil
}
