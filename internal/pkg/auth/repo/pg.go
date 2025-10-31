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
		CheckUserExistsQuery,
		login,
	).Scan(&exists)
	return exists, err
}

func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.db.Exec(
		ctx,
		CreateUserQuery,
		user.ID, user.Login, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *AuthRepository) CheckUserLogin(ctx context.Context, login string) (models.User, error) {
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
		return models.User{}, err
	}
	return user, nil
}

func (r *AuthRepository) IncrementUserVersion(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(
		ctx,
		IncrementUserVersionQuery,
		userID,
	)
	return err
}

func (r *AuthRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL Error: %s, Code: %s, Detail: %s\n",
				pgErr.Message, pgErr.Code, pgErr.Detail)
		}

		fmt.Printf("Error getting user by ID %v: %v\n", user, err)
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
