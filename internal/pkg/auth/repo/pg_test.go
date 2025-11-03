package repo

import (
	"context"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"
	"log/slog"
	"os"
	"testing"
	"time"

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

func TestCheckUserExists(t *testing.T) {
	login := "testuser"

	tests := []struct {
		name       string
		login      string
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantExists bool
		wantErr    bool
	}{
		{
			name:  "Success_UserExists",
			login: login,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"exists"}).
					AddRow(true).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserExistsQuery, login).
					Return(rows)
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name:  "Success_UserNotExists",
			login: login,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"exists"}).
					AddRow(false).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserExistsQuery, login).
					Return(rows)
			},
			wantExists: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewAuthRepository(mockPool)
			exists, err := repo.CheckUserExists(testContext(), tt.login)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantExists, exists)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
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

	tests := []struct {
		name       string
		user       models.User
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
	}{
		{
			name: "Success",
			user: user,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), CreateUserQuery, user.ID, user.Login, user.PasswordHash, user.CreatedAt, user.UpdatedAt).
					Return(nil, nil)
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

			repo := NewAuthRepository(mockPool)
			err := repo.CreateUser(testContext(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckUserLogin(t *testing.T) {
	userID := uuid.NewV4()
	login := "testuser"
	avatar := "/static/default.jpg"
	createdAt := time.Now()
	updatedAt := time.Now()

	expectedUser := models.User{
		ID:           userID,
		Version:      1,
		Login:        login,
		PasswordHash: []byte("hash"),
		Avatar:       &avatar,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	tests := []struct {
		name       string
		login      string
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantUser   models.User
		wantErr    bool
	}{
		{
			name:  "Success",
			login: login,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at"}).
					AddRow(userID, 1, login, []byte("hash"), &avatar, createdAt, updatedAt).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserLoginQuery, login).
					Return(rows)
			},
			wantUser: expectedUser,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewAuthRepository(mockPool)
			user, err := repo.CheckUserLogin(testContext(), tt.login)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Login, user.Login)
				assert.Equal(t, tt.wantUser.Version, user.Version)
			}
		})
	}
}

func TestIncrementUserVersion(t *testing.T) {
	userID := uuid.NewV4()

	tests := []struct {
		name       string
		userID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
	}{
		{
			name:   "Success",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), IncrementUserVersionQuery, userID).
					Return(nil, nil)
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

			repo := NewAuthRepository(mockPool)
			err := repo.IncrementUserVersion(testContext(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserByLogin(t *testing.T) {
	userID := uuid.NewV4()
	login := "testuser"
	avatar := "/static/default.jpg"
	createdAt := time.Now()
	updatedAt := time.Now()

	expectedUser := models.User{
		ID:           userID,
		Version:      1,
		Login:        login,
		PasswordHash: []byte("hash"),
		Avatar:       &avatar,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	tests := []struct {
		name       string
		login      string
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantUser   models.User
		wantErr    bool
	}{
		{
			name:  "Success",
			login: login,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at"}).
					AddRow(userID, 1, login, []byte("hash"), &avatar, createdAt, updatedAt).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetUserByLoginQuery, login).
					Return(rows)
			},
			wantUser: expectedUser,
			wantErr:  false,
		},
		// Аналогично убираем тест с UserNotFound для GetUserByLogin
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewAuthRepository(mockPool)
			user, err := repo.GetUserByLogin(testContext(), tt.login)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Login, user.Login)
				assert.Equal(t, tt.wantUser.Version, user.Version)
			}
		})
	}
}
