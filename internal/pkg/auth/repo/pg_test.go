package repo

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
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

type errorRow struct {
	err error
}

func (r errorRow) Scan(dest ...interface{}) error {
	return r.err
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
				rows := pgxpoolmock.NewRows([]string{"exists"}).AddRow(true).ToPgxRows()
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
				rows := pgxpoolmock.NewRows([]string{"exists"}).AddRow(false).ToPgxRows()
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
		Avatar:       avatar,
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
		{
			name: "Error_DatabaseError",
			user: user,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), CreateUserQuery, user.ID, user.Login, user.PasswordHash, user.CreatedAt, user.UpdatedAt).
					Return(nil, errors.New("database error"))
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
		{
			name:   "Error_DatabaseError",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), IncrementUserVersionQuery, userID).
					Return(nil, errors.New("database error"))
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
		Avatar:       avatar,
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
					AddRow(userID, 1, login, []byte("hash"), avatar, createdAt, updatedAt).
					ToPgxRows()
				rows.Next()
				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserLoginQuery, login).
					Return(rows)
			},
			wantUser: expectedUser,
			wantErr:  false,
		},
		{
			name:  "Error_UserNotFound",
			login: login,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserLoginQuery, login).
					Return(errorRow{err: pgx.ErrNoRows})
			},
			wantUser: models.User{},
			wantErr:  true,
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
				assert.Equal(t, tt.wantUser.Avatar, user.Avatar)
			}
		})
	}
}

func TestGetUserByLogin(t *testing.T) {
	userID := uuid.NewV4()
	login := "testuser"
	avatar := "/static/default.jpg"
	has2FA := false
	createdAt := time.Now()
	updatedAt := time.Now()

	expectedUser := models.User{
		ID:           userID,
		Version:      1,
		Login:        login,
		PasswordHash: []byte("hash"),
		Avatar:       avatar,
		Has2FA:       has2FA,
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
				rows := pgxpoolmock.NewRows([]string{"id", "version", "login", "password_hash", "avatar", "has_2fa", "created_at", "updated_at"}).
					AddRow(userID, 1, login, []byte("hash"), avatar, has2FA, createdAt, updatedAt).
					ToPgxRows()
				rows.Next()
				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetUserByLoginQuery, login).
					Return(rows)
			},
			wantUser: expectedUser,
			wantErr:  false,
		},
		{
			name:  "Error_UserNotFound",
			login: login,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetUserByLoginQuery, login).
					Return(errorRow{err: pgx.ErrNoRows})
			},
			wantUser: models.User{},
			wantErr:  true,
		},
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
				assert.Equal(t, tt.wantUser.Avatar, user.Avatar)
				assert.Equal(t, tt.wantUser.Has2FA, user.Has2FA)
			}
		})
	}
}

func TestEnable2FA(t *testing.T) {
	userID := uuid.NewV4()
	secret := "QWERTYASDFGZ"

	tests := []struct {
		name       string
		userID     uuid.UUID
		secret     string
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantResult models.EnableTwoFactorResponse
		wantErr    bool
	}{
		{
			name:   "Success",
			userID: userID,
			secret: secret,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"has_2fa"}).
					AddRow(true).
					ToPgxRows()
				rows.Next()
				mockPool.EXPECT().
					QueryRow(gomock.Any(), Enable2FaQuery, userID, secret).
					Return(rows)
			},
			wantResult: models.EnableTwoFactorResponse{Has2FA: true},
			wantErr:    false,
		},
		{
			name:   "Error_UserNotFound",
			userID: userID,
			secret: secret,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					QueryRow(gomock.Any(), Enable2FaQuery, userID, secret).
					Return(errorRow{err: pgx.ErrNoRows})
			},
			wantResult: models.EnableTwoFactorResponse{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewAuthRepository(mockPool)
			result, err := repo.Enable2FA(testContext(), tt.userID, tt.secret)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult.Has2FA, result.Has2FA)
			}
		})
	}
}

func TestDisable2FA(t *testing.T) {
	userID := uuid.NewV4()

	tests := []struct {
		name       string
		userID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantResult models.DisableTwoFactorResponse
		wantErr    bool
	}{
		{
			name:   "Success",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"has_2fa"}).
					AddRow(false).
					ToPgxRows()
				rows.Next()
				mockPool.EXPECT().
					QueryRow(gomock.Any(), Disable2FaQuery, userID).
					Return(rows)
			},
			wantResult: models.DisableTwoFactorResponse{Has2FA: false},
			wantErr:    false,
		},
		{
			name:   "Error_UserNotFound",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					QueryRow(gomock.Any(), Disable2FaQuery, userID).
					Return(errorRow{err: pgx.ErrNoRows})
			},
			wantResult: models.DisableTwoFactorResponse{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewAuthRepository(mockPool)
			result, err := repo.Disable2FA(testContext(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult.Has2FA, result.Has2FA)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	userID := uuid.NewV4()
	avatar := "/static/default.jpg"
	createdAt := time.Now()
	updatedAt := time.Now()

	expectedUser := models.User{
		ID:           userID,
		Version:      1,
		Login:        "testuser",
		PasswordHash: []byte("hash"),
		Avatar:       avatar,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

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
				rows := pgxpoolmock.NewRows([]string{"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at"}).
					AddRow(userID, 1, "testuser", []byte("hash"), avatar, createdAt, updatedAt).
					ToPgxRows()
				rows.Next()
				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetUserByIDQuery, userID).
					Return(rows)
			},
			wantUser: expectedUser,
			wantErr:  false,
		},
		{
			name:   "Error_UserNotFound",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetUserByIDQuery, userID).
					Return(errorRow{err: pgx.ErrNoRows})
			},
			wantUser: models.User{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewAuthRepository(mockPool)
			user, err := repo.GetUserByID(testContext(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Login, user.Login)
				assert.Equal(t, tt.wantUser.Version, user.Version)
				assert.Equal(t, tt.wantUser.Avatar, user.Avatar)
			}
		})
	}
}

func TestCheckUserTwoFactor(t *testing.T) {
	userID := uuid.NewV4()

	tests := []struct {
		name       string
		userID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantHas2FA bool
		wantErr    bool
	}{
		{
			name:   "Success_2FAEnabled",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"has_2fa"}).
					AddRow(true).
					ToPgxRows()
				rows.Next()
				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserTwoFactorQuery, userID).
					Return(rows)
			},
			wantHas2FA: true,
			wantErr:    false,
		},
		{
			name:   "Success_2FADisabled",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"has_2fa"}).
					AddRow(false).
					ToPgxRows()
				rows.Next()
				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserTwoFactorQuery, userID).
					Return(rows)
			},
			wantHas2FA: false,
			wantErr:    false,
		},
		{
			name:   "Error_UserNotFound",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserTwoFactorQuery, userID).
					Return(errorRow{err: pgx.ErrNoRows})
			},
			wantHas2FA: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewAuthRepository(mockPool)
			has2FA, err := repo.CheckUserTwoFactor(testContext(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantHas2FA, has2FA)
			}
		})
	}
}

func TestGetUserSecretCode(t *testing.T) {
	userID := uuid.NewV4()
	secretCode := "QWERTYASDFGZ"

	tests := []struct {
		name       string
		userID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantSecret string
	}{
		{
			name:   "Success",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"secret_code"}).
					AddRow(secretCode).
					ToPgxRows()
				rows.Next()
				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserSecretCodeQuery, userID).
					Return(rows)
			},
			wantSecret: secretCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewAuthRepository(mockPool)
			secret := repo.GetUserSecretCode(testContext(), tt.userID)

			assert.Equal(t, tt.wantSecret, secret)
		})
	}
}
