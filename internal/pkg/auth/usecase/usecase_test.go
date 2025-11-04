package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

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

func TestHashPass(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Success",
			password: "testpassword123",
		},
		{
			name:     "Empty password",
			password: "",
		},
		{
			name:     "Long password",
			password: "verylongpasswordwithspecialchars!@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := HashPass(tt.password)
			assert.NotEmpty(t, hash)
			assert.GreaterOrEqual(t, len(hash), 40) // salt (8) + hash (32)
		})
	}
}

func TestCheckPass(t *testing.T) {
	password := "testpassword123"
	correctHash := HashPass(password)
	wrongHash := HashPass("wrongpassword")

	tests := []struct {
		name     string
		hash     []byte
		password string
		expected bool
	}{
		{
			name:     "Correct password",
			hash:     correctHash,
			password: password,
			expected: true,
		},
		{
			name:     "Wrong password",
			hash:     correctHash,
			password: "wrongpassword",
			expected: false,
		},
		{
			name:     "Different hash",
			hash:     wrongHash,
			password: password,
			expected: false,
		},
		{
			name:     "Empty password",
			hash:     correctHash,
			password: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPass(tt.hash, tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewAuthUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)

	t.Run("Success creation", func(t *testing.T) {
		usecase := NewAuthUsecase(mockRepo)
		assert.NotNil(t, usecase)
		assert.Equal(t, mockRepo, usecase.authRepo)
	})

	t.Run("Creation with nil repo", func(t *testing.T) {
		usecase := NewAuthUsecase(nil)
		assert.NotNil(t, usecase)
		assert.Nil(t, usecase.authRepo)
	})
}

func TestAuthUsecase_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	userID := uuid.NewV4()
	login := "testuser"

	t.Run("Success", func(t *testing.T) {
		token, err := usecase.GenerateToken(userID, login)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token can be parsed
		parsedToken, err := usecase.ParseToken(token)
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, login, claims["login"])
		assert.Equal(t, userID.String(), claims["id"])
	})

	t.Run("Empty login", func(t *testing.T) {
		token, err := usecase.GenerateToken(userID, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestAuthUsecase_ParseToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	userID := uuid.NewV4()
	login := "testuser"
	validToken, _ := usecase.GenerateToken(userID, login)

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "Valid token",
			token:       validToken,
			expectError: false,
		},
		{
			name:        "Empty token",
			token:       "",
			expectError: true,
		},
		{
			name:        "Malformed token",
			token:       "malformed.token",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedToken, err := usecase.ParseToken(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, parsedToken)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, parsedToken)
				assert.True(t, parsedToken.Valid)
			}
		})
	}
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
			name: "Error - CheckUserExists fails",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserExists(gomock.Any(), login).
					Return(false, errors.New("db error"))
			},
			req: models.SignUpInput{
				Login:    login,
				Password: password,
			},
			expectError: true,
		},
		{
			name: "Error - CreateUser fails",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserExists(gomock.Any(), login).
					Return(false, nil)
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(errors.New("create error"))
			},
			req: models.SignUpInput{
				Login:    login,
				Password: password,
			},
			expectError: true,
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
				assert.NotNil(t, user.ID)
				assert.NotNil(t, user.Avatar)
				assert.Equal(t, "avatars/default.png", *user.Avatar)
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
		{
			name: "Error - empty password",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserLogin(gomock.Any(), login).
					Return(existingUser, nil)
			},
			req: models.SignInInput{
				Login:    login,
				Password: "",
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
		{
			name:        "Error - wrong type in context",
			ctx:         context.WithValue(testContext(), auth.UserKey, "not-a-user"),
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
			name: "Error - IncrementUserVersion fails",
			ctx:  context.WithValue(testContext(), auth.UserKey, user),
			setupMock: func() {
				mockRepo.EXPECT().
					IncrementUserVersion(gomock.Any(), userID).
					Return(errors.New("db error"))
			},
			expectError: true,
		},
		{
			name:        "Error - no user in context",
			ctx:         testContext(),
			setupMock:   func() {},
			expectError: true,
		},
		{
			name:        "Error - wrong type in context",
			ctx:         context.WithValue(testContext(), auth.UserKey, "not-a-user"),
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

	// Create expired token
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userID,
		"login": login,
		"exp":   time.Now().Add(-time.Hour).Unix(), // expired
	})
	expiredTokenString, _ := expiredToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

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
		{
			name:        "Error - expired token",
			token:       expiredTokenString,
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

func TestValidateFunctions(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		expectedValid bool
	}{
		{
			name:          "Valid credentials",
			login:         "user123",
			password:      "pass123",
			expectedValid: true,
		},
		{
			name:          "Invalid login length",
			login:         "usr",
			password:      "pass123",
			expectedValid: false,
		},
		{
			name:          "Invalid password length",
			login:         "user123",
			password:      "pass",
			expectedValid: false,
		},
		{
			name:          "Both invalid",
			login:         "usr",
			password:      "pass",
			expectedValid: false,
		},
		{
			name:          "Empty credentials",
			login:         "",
			password:      "",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, dataIsValid := auth.Validaton(tt.login, tt.password)
			assert.Equal(t, tt.expectedValid, dataIsValid)
		})
	}
}
