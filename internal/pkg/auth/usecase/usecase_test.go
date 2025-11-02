package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/auth/mocks"
	"kinopoisk/internal/pkg/middleware/logger"

	"github.com/golang-jwt/jwt"
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

func TestAuthUsecase_SignUpUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	login := "testuser"
	password := "testpass123"

	tests := []struct {
		name        string
		setupMock   func()
		req         models.SignUpInput
		expectError bool
		errorType   error
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserExists(gomock.Any(), login).
					Return(false, nil)
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			req: models.SignUpInput{
				Login:    login,
				Password: password,
			},
			expectError: false,
		},
		{
			name: "Error - user already exists",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserExists(gomock.Any(), login).
					Return(true, nil)
			},
			req: models.SignUpInput{
				Login:    login,
				Password: password,
			},
			expectError: true,
			errorType:   auth.ErrorConflict,
		},
		{
			name:      "Error - invalid login",
			setupMock: func() {},
			req: models.SignUpInput{
				Login:    "usr",
				Password: password,
			},
			expectError: true,
			errorType:   auth.ErrorBadRequest,
		},
		{
			name:      "Error - invalid password",
			setupMock: func() {},
			req: models.SignUpInput{
				Login:    login,
				Password: "pass",
			},
			expectError: true,
			errorType:   auth.ErrorBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			user, token, err := usecase.SignUpUser(testContext(), tt.req)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Equal(t, tt.req.Login, user.Login)
			}
		})
	}
}

func TestAuthUsecase_SignInUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	userID := uuid.NewV4()
	login := "testuser"
	password := "testpass123"

	existingUser := models.User{
		ID:           userID,
		Login:        login,
		PasswordHash: HashPass(password),
		Version:      1,
	}

	notFoundError := errors.New("user not found")

	tests := []struct {
		name        string
		setupMock   func()
		req         models.SignInInput
		expectError bool
		errorType   error
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserLogin(gomock.Any(), login).
					Return(existingUser, nil)
			},
			req: models.SignInInput{
				Login:    login,
				Password: password,
			},
			expectError: false,
		},
		{
			name: "Error - user not found",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserLogin(gomock.Any(), login).
					Return(models.User{}, notFoundError)
			},
			req: models.SignInInput{
				Login:    login,
				Password: password,
			},
			expectError: true,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserLogin(gomock.Any(), login).
					Return(models.User{}, errors.New("database error"))
			},
			req: models.SignInInput{
				Login:    login,
				Password: password,
			},
			expectError: true,
		},
		{
			name: "Error - wrong password",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserLogin(gomock.Any(), login).
					Return(existingUser, nil)
			},
			req: models.SignInInput{
				Login:    login,
				Password: "wrongpass",
			},
			expectError: true,
			errorType:   auth.ErrorBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			user, token, err := usecase.SignInUser(testContext(), tt.req)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Equal(t, existingUser.ID, user.ID)
			}
		})
	}
}

func TestAuthUsecase_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	userID := uuid.NewV4()
	user := models.User{
		ID:    userID,
		Login: "testuser",
	}

	tests := []struct {
		name        string
		ctx         context.Context
		expectError bool
	}{
		{
			name:        "Success",
			ctx:         context.WithValue(testContext(), auth.UserKey, user),
			expectError: false,
		},
		{
			name:        "Error - no user in context",
			ctx:         testContext(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := usecase.CheckAuth(tt.ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, auth.ErrorUnauthorized)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, user, result)
			}
		})
	}
}

func TestAuthUsecase_LogOutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	userID := uuid.NewV4()
	user := models.User{
		ID:    userID,
		Login: "testuser",
	}

	tests := []struct {
		name        string
		ctx         context.Context
		setupMock   func()
		expectError bool
	}{
		{
			name: "Success",
			ctx:  context.WithValue(testContext(), auth.UserKey, user),
			setupMock: func() {
				mockRepo.EXPECT().
					IncrementUserVersion(gomock.Any(), userID).
					Return(nil)
			},
			expectError: false,
		},
		{
			name:        "Error - no user in context",
			ctx:         testContext(),
			setupMock:   func() {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := usecase.LogOutUser(tt.ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, auth.ErrorUnauthorized)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthUsecase_ValidateAndGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	userID := uuid.NewV4()
	login := "testuser"
	user := models.User{
		ID:    userID,
		Login: login,
	}

	validToken, _ := usecase.GenerateToken(userID, login)

	tests := []struct {
		name        string
		token       string
		setupMock   func()
		expectError bool
	}{
		{
			name:  "Success",
			token: validToken,
			setupMock: func() {
				mockRepo.EXPECT().
					GetUserByLogin(gomock.Any(), login).
					Return(user, nil)
			},
			expectError: false,
		},
		{
			name:        "Error - empty token",
			token:       "",
			setupMock:   func() {},
			expectError: true,
		},
		{
			name:        "Error - invalid token",
			token:       "invalid.token.here",
			setupMock:   func() {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.ValidateAndGetUser(testContext(), tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, auth.ErrorUnauthorized)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, user, result)
			}
		})
	}
}

func TestAuthUsecase_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	userID := uuid.NewV4()
	login := "testuser"

	token, err := usecase.GenerateToken(userID, login)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := usecase.ParseToken(token)
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, login, claims["login"])
}

func TestValidateFunctions(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		loginValid    bool
		passwordValid bool
	}{
		{
			name:          "Valid credentials",
			login:         "user123",
			password:      "pass123",
			loginValid:    true,
			passwordValid: true,
		},
		{
			name:          "Invalid login length",
			login:         "usr",
			password:      "pass123",
			loginValid:    false,
			passwordValid: true,
		},
		{
			name:          "Invalid password length",
			login:         "user123",
			password:      "pass",
			loginValid:    true,
			passwordValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, loginValid := ValidateLogin(tt.login)
			_, passValid := ValidatePassword(tt.password)

			assert.Equal(t, tt.loginValid, loginValid)
			assert.Equal(t, tt.passwordValid, passValid)
		})
	}
}
