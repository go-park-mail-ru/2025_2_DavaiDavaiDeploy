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

func TestGetUserByID(t *testing.T) {
	userID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()
	avatar := "/static/avatar.png"

	tests := []struct {
		name       string
		userID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantUser   models.User
		wantErr    bool
	}{
		{
			name:   "Success",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at",
				}).
					AddRow(
						userID,
						1,
						"testuser",
						[]byte("hashedpassword"),
						avatar, // Убрать & - передавать строку, а не указатель
						createdAt,
						updatedAt,
					).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetUserByIDQuery, userID).
					Return(rows)
			},
			wantUser: models.User{
				ID:           userID,
				Version:      1,
				Login:        "testuser",
				PasswordHash: []byte("hashedpassword"),
				Avatar:       avatar,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewUserRepository(mockPool)
			user, err := repo.GetUserByID(testContext(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Version, user.Version)
				assert.Equal(t, tt.wantUser.Login, user.Login)
				assert.Equal(t, tt.wantUser.PasswordHash, user.PasswordHash)
				assert.Equal(t, tt.wantUser.Avatar, user.Avatar)
			}
		})
	}
}

func TestGetUserByLogin(t *testing.T) {
	userID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()
	avatar := "/static/avatar.png"

	tests := []struct {
		name       string
		login      string
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantUser   models.User
		wantErr    bool
	}{
		{
			name:  "Success",
			login: "testuser",
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at",
				}).
					AddRow(
						userID,
						1,
						"testuser",
						[]byte("hashedpassword"),
						avatar,
						createdAt,
						updatedAt,
					).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetUserByLoginQuery, "testuser").
					Return(rows)
			},
			wantUser: models.User{
				ID:           userID,
				Version:      1,
				Login:        "testuser",
				PasswordHash: []byte("hashedpassword"),
				Avatar:       avatar,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewUserRepository(mockPool)
			user, err := repo.GetUserByLogin(testContext(), tt.login)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Version, user.Version)
				assert.Equal(t, tt.wantUser.Login, user.Login)
				assert.Equal(t, tt.wantUser.PasswordHash, user.PasswordHash)
				assert.Equal(t, tt.wantUser.Avatar, user.Avatar)
			}
		})
	}
}

func TestUpdateUserPassword(t *testing.T) {
	userID := uuid.NewV4()
	version := 2
	newPasswordHash := []byte("newhashedpassword")

	tests := []struct {
		name         string
		version      int
		userID       uuid.UUID
		passwordHash []byte
		repoMocker   func(*pgxpoolmock.MockPgxPool)
		wantErr      bool
	}{
		{
			name:         "Success",
			version:      version,
			userID:       userID,
			passwordHash: newPasswordHash,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), UpdateUserPasswordQuery, newPasswordHash, version, userID).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name:         "Error",
			version:      version,
			userID:       userID,
			passwordHash: newPasswordHash,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), UpdateUserPasswordQuery, newPasswordHash, version, userID).
					Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewUserRepository(mockPool)
			err := repo.UpdateUserPassword(testContext(), tt.version, tt.userID, tt.passwordHash)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUserAvatar(t *testing.T) {
	userID := uuid.NewV4()
	version := 2
	avatarPath := "/static/new_avatar.png"

	tests := []struct {
		name       string
		version    int
		userID     uuid.UUID
		avatarPath string
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
	}{
		{
			name:       "Success",
			version:    version,
			userID:     userID,
			avatarPath: avatarPath,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), UpdateUserAvatarQuery, avatarPath, version, userID).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name:       "Error",
			version:    version,
			userID:     userID,
			avatarPath: avatarPath,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), UpdateUserAvatarQuery, avatarPath, version, userID).
					Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewUserRepository(mockPool)
			err := repo.UpdateUserAvatar(testContext(), tt.version, tt.userID, tt.avatarPath)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
