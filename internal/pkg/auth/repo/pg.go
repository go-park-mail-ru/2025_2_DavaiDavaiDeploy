package repo

import (
	"context"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CheckUserExists(ctx context.Context, login string) (bool, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var exists bool
	err := r.db.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM user_table WHERE login = $1)",
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
		"INSERT INTO user_table (id, login, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.Login, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		logger.Error("failed to insert into user_table: " + err.Error())
		return err
	}
	return err
}

func (r *AuthRepository) CheckUserLogin(ctx context.Context, login string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	var user models.User
	err := r.db.QueryRow(ctx,
		"SELECT id, version, login, password_hash, avatar, created_at, updated_at FROM user_table WHERE login = $1",
		login,
	).Scan(&user.ID,
		&user.Version,
		&user.Login,
		&user.PasswordHash,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		logger.Error("failed to get user: " + err.Error())
		return models.User{}, err
	}
	return user, nil
}

func (r *AuthRepository) IncrementUserVersion(ctx context.Context, userID uuid.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	_, err := r.db.Exec(
		ctx,
		"UPDATE user_table SET version = version + 1 WHERE id = $1",
		userID,
	)
	if err != nil {
		logger.Error("failed to update version: " + err.Error())
		return err
	}
	return err
}

func (r *AuthRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var user models.User
	err := r.db.QueryRow(
		ctx,
		"SELECT id, version, login, password_hash, avatar, created_at, updated_at FROM user_table WHERE login = $1",
		login,
	).Scan(
		&user.ID, &user.Version, &user.Login,
		&user.PasswordHash, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		logger.Error("failed to get user: " + err.Error())
		return models.User{}, err
	}
	return user, nil
}
