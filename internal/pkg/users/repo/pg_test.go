package repo

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

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

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mockPool)

	userID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()
	avatar := "/static/avatar.png"

	rows := pgxpoolmock.NewRows([]string{
		"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at",
	}).
		AddRow(
			userID,
			1,
			"testuser",
			[]byte("hashedpassword"),
			&avatar,
			createdAt,
			updatedAt,
		).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetUserByIDQuery, userID).
		Return(rows)

	user, err := repo.GetUserByID(testContext(), userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, 1, user.Version)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, []byte("hashedpassword"), user.PasswordHash)
	assert.Equal(t, &avatar, user.Avatar)
}

func TestGetUserByID_WithNullAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mockPool)

	userID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := pgxpoolmock.NewRows([]string{
		"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at",
	}).
		AddRow(
			userID,
			1,
			"testuser",
			[]byte("hashedpassword"),
			nil,
			createdAt,
			updatedAt,
		).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetUserByIDQuery, userID).
		Return(rows)

	user, err := repo.GetUserByID(testContext(), userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, 1, user.Version)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, []byte("hashedpassword"), user.PasswordHash)
	assert.Nil(t, user.Avatar)
}

func TestGetUserByLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mockPool)

	userID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()
	avatar := "/static/avatar.png"

	rows := pgxpoolmock.NewRows([]string{
		"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at",
	}).
		AddRow(
			userID,
			1,
			"testuser",
			[]byte("hashedpassword"),
			&avatar,
			createdAt,
			updatedAt,
		).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetUserByLoginQuery, "testuser").
		Return(rows)

	user, err := repo.GetUserByLogin(testContext(), "testuser")

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, 1, user.Version)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, []byte("hashedpassword"), user.PasswordHash)
	assert.Equal(t, &avatar, user.Avatar)
}

func TestGetUserByLogin_WithNullAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mockPool)

	userID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := pgxpoolmock.NewRows([]string{
		"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at",
	}).
		AddRow(
			userID,
			1,
			"testuser",
			[]byte("hashedpassword"),
			nil,
			createdAt,
			updatedAt,
		).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetUserByLoginQuery, "testuser").
		Return(rows)

	user, err := repo.GetUserByLogin(testContext(), "testuser")

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, 1, user.Version)
	assert.Equal(t, "testuser", user.Login)
	assert.Equal(t, []byte("hashedpassword"), user.PasswordHash)
	assert.Nil(t, user.Avatar)
}

func TestUpdateUserPassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mockPool)

	userID := uuid.NewV4()
	version := 2
	newPasswordHash := []byte("newhashedpassword")

	mockPool.EXPECT().
		Exec(gomock.Any(), UpdateUserPasswordQuery, newPasswordHash, version, userID).
		Return(nil, nil)

	err := repo.UpdateUserPassword(testContext(), version, userID, newPasswordHash)

	assert.NoError(t, err)
}

func TestUpdateUserPassword_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mockPool)

	userID := uuid.NewV4()
	version := 2
	newPasswordHash := []byte("newhashedpassword")

	mockPool.EXPECT().
		Exec(gomock.Any(), UpdateUserPasswordQuery, newPasswordHash, version, userID).
		Return(nil, assert.AnError)

	err := repo.UpdateUserPassword(testContext(), version, userID, newPasswordHash)

	assert.Error(t, err)
}

func TestUpdateUserAvatar_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mockPool)

	userID := uuid.NewV4()
	version := 2
	avatarPath := "/static/new_avatar.png"

	mockPool.EXPECT().
		Exec(gomock.Any(), UpdateUserAvatarQuery, avatarPath, version, userID).
		Return(nil, nil)

	err := repo.UpdateUserAvatar(testContext(), version, userID, avatarPath)

	assert.NoError(t, err)
}

func TestUpdateUserAvatar_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewUserRepository(mockPool)

	userID := uuid.NewV4()
	version := 2
	avatarPath := "/static/new_avatar.png"

	mockPool.EXPECT().
		Exec(gomock.Any(), UpdateUserAvatarQuery, avatarPath, version, userID).
		Return(nil, assert.AnError)

	err := repo.UpdateUserAvatar(testContext(), version, userID, avatarPath)

	assert.Error(t, err)
}
