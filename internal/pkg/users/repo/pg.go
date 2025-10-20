package repo

import (
	"context"
	"errors"
	"fmt"
	"kinopoisk/internal/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	var user models.User
	err := u.db.QueryRow(
		ctx,
		"SELECT id, version, login, password_hash, avatar, created_at, updated_at FROM user_table WHERE id = $1",
		id,
	).Scan(
		&user.ID, &user.Version, &user.Login,
		&user.PasswordHash, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL Error: %s, Code: %s, Detail: %s\n",
				pgErr.Message, pgErr.Code, pgErr.Detail)
		}

		fmt.Printf("Error %s: %v\n", id, err)
		return models.User{}, fmt.Errorf("failed to: %w", err)
	}
	return user, nil
}

func (u *UserRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User
	err := u.db.QueryRow(
		ctx,
		"SELECT id, version, login, password_hash, avatar, created_at, updated_at FROM user_table WHERE login = $1",
		login,
	).Scan(
		&user.ID, &user.Version, &user.Login,
		&user.PasswordHash, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u *UserRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error {
	_, err := u.db.Exec(
		ctx,
		"UPDATE user_table SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		passwordHash, userID,
	)
	return err
}

func (u *UserRepository) UpdateUserAvatar(ctx context.Context, userID uuid.UUID, avatarPath string) error {
	_, err := u.db.Exec(
		ctx,
		"UPDATE user_table SET avatar = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		avatarPath, userID,
	)
	return err
}
