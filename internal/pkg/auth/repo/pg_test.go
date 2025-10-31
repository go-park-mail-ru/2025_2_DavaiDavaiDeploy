package repo

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func testContext() context.Context {
	testLogger := testLogger()
	return context.WithValue(context.Background(), logger.LoggerKey, testLogger)
}

func TestCheckUserExists_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewAuthRepository(mockPool)

	login := "testuser"

	rows := pgxpoolmock.NewRows([]string{"exists"}).
		AddRow(true).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), CheckUserExistsQuery, login).
		Return(rows)

	exists, err := repo.CheckUserExists(testContext(), login)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewAuthRepository(mockPool)

	userID := uuid.NewV4()
	avatar := "/static/default.jpg"
	user := models.User{
		ID:           userID,
		Login:        "testuser",
		PasswordHash: []byte("hash"),
		Avatar:       &avatar,
		Version:      1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockPool.EXPECT().
		Exec(gomock.Any(), CreateUserQuery, user.ID, user.Login, user.PasswordHash, user.CreatedAt, user.UpdatedAt).
		Return(nil, nil)

	err := repo.CreateUser(testContext(), user)

	assert.NoError(t, err)
}

func TestCheckUserLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewAuthRepository(mockPool)

	userID := uuid.NewV4()
	login := "testuser"
	avatar := "/static/default.jpg"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := pgxpoolmock.NewRows([]string{"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at"}).
		AddRow(userID, 1, login, []byte("hash"), &avatar, createdAt, updatedAt).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), CheckUserLoginQuery, login).
		Return(rows)

	user, err := repo.CheckUserLogin(testContext(), login)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, login, user.Login)
}

func TestIncrementUserVersion_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewAuthRepository(mockPool)

	userID := uuid.NewV4()

	mockPool.EXPECT().
		Exec(gomock.Any(), IncrementUserVersionQuery, userID).
		Return(nil, nil)

	err := repo.IncrementUserVersion(testContext(), userID)

	assert.NoError(t, err)
}

func TestGetUserByLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewAuthRepository(mockPool)

	userID := uuid.NewV4()
	login := "testuser"
	avatar := "/static/default.jpg"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := pgxpoolmock.NewRows([]string{"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at"}).
		AddRow(userID, 1, login, []byte("hash"), &avatar, createdAt, updatedAt).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetUserByLoginQuery, login).
		Return(rows)

	user, err := repo.GetUserByLogin(testContext(), login)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, login, user.Login)
}
