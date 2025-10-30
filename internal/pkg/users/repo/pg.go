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

var (
	getUserByIDQuery = `
		SELECT id, version, login, password_hash, avatar, created_at, updated_at 
		FROM user_table 
		WHERE id = $1`

	getUserByLoginQuery = `
		SELECT id, version, login, password_hash, avatar, created_at, updated_at 
		FROM user_table 
		WHERE login = $1`

	updateUserPasswordQuery = `
		UPDATE user_table 
		SET password_hash = $1, version = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3`

	updateUserAvatarQuery = `
		UPDATE user_table 
		SET avatar = $1, version = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3`
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
		getUserByIDQuery,
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
		getUserByLoginQuery,
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

func (u *UserRepository) UpdateUserPassword(ctx context.Context, version int, userID uuid.UUID, passwordHash []byte) error {
	_, err := u.db.Exec(
		ctx,
		updateUserPasswordQuery,
		passwordHash, version, userID,
	)
	return err
}

func (u *UserRepository) UpdateUserAvatar(ctx context.Context, version int, userID uuid.UUID, avatarPath string) error {
	_, err := u.db.Exec(
		ctx,
		updateUserAvatarQuery,
		avatarPath, version, userID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL Error: %s, Code: %s, Detail: %s\n",
				pgErr.Message, pgErr.Code, pgErr.Detail)
		}
		return fmt.Errorf("failed to update avatar: %w", err)
	}
	return nil
}
