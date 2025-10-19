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

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CheckUserExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM user_table WHERE login = $1)",
		login,
	).Scan(&exists)
	return exists, err
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO user_table (id, login, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		user.ID, user.Login, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *AuthRepository) CheckUserLogin(ctx context.Context, login string) (*models.User, error) {
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
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) IncrementUserVersion(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(
		ctx,
		"UPDATE user_table SET version = version + 1 WHERE id = $1",
		userID,
	)
	return err
}

func (r *AuthRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL Error: %s, Code: %s, Detail: %s\n",
				pgErr.Message, pgErr.Code, pgErr.Detail)
		}

		fmt.Printf("Error getting film by ID %s: %v\n", user, err)
		return nil, fmt.Errorf("failed to get film: %w", err)
	}
	return &user, nil
}
