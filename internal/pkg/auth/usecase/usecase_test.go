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

	jwt "github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
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
			assert.GreaterOrEqual(t, len(hash), 40)
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
	version := 1

	t.Run("Success", func(t *testing.T) {
		token, err := usecase.GenerateToken(userID, login, version)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := usecase.ParseToken(token)
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, login, claims["login"])
		assert.Equal(t, userID.String(), claims["id"])
		assert.Equal(t, float64(version), claims["version"])
	})

	t.Run("Empty login", func(t *testing.T) {
		token, err := usecase.GenerateToken(userID, "", version)
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
	version := 1
	validToken, _ := usecase.GenerateToken(userID, login, version)

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
				assert.Equal(t, "avatars/default.png", user.Avatar)
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
	version := 1

	existingUser := models.User{
		ID:           userID,
		Login:        login,
		PasswordHash: HashPass(password),
		Version:      version,
	}

	tests := []struct {
		name        string
		setupMock   func()
		req         models.SignInInput
		expectError bool
		errorType   error
	}{
		{
			name: "Success without 2FA",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserLogin(gomock.Any(), login).
					Return(existingUser, nil)
				mockRepo.EXPECT().
					GetUserSecretCode(gomock.Any(), userID).
					Return("")
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
				mockRepo.EXPECT().
					GetUserSecretCode(gomock.Any(), userID).
					Return("")
			},
			req: models.SignInInput{
				Login:    login,
				Password: "wrongpass",
			},
			expectError: true,
			errorType:   auth.ErrorBadRequest,
		},
		{
			name: "Error - 2FA enabled but no code provided",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserLogin(gomock.Any(), login).
					Return(existingUser, nil)
				mockRepo.EXPECT().
					GetUserSecretCode(gomock.Any(), userID).
					Return("SECRETCODE123")
			},
			req: models.SignInInput{
				Login:    login,
				Password: password,
			},
			expectError: true,
			errorType:   auth.ErrorPreconditionFailed,
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

func TestAuthUsecase_LogOutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	userID := uuid.NewV4()

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					IncrementUserVersion(gomock.Any(), userID).
					Return(nil)
			},
			expectError: false,
		},
		{
			name: "Error - IncrementUserVersion fails",
			setupMock: func() {
				mockRepo.EXPECT().
					IncrementUserVersion(gomock.Any(), userID).
					Return(errors.New("db error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := usecase.LogOutUser(testContext(), userID)

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
	version := 1
	user := models.User{
		ID:      userID,
		Login:   login,
		Version: version,
	}

	validToken, _ := usecase.GenerateToken(userID, login, version)

	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      userID,
		"login":   login,
		"version": version,
		"exp":     time.Now().Add(-time.Hour).Unix(),
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
			token:       "invalid.token.123",
			setupMock:   func() {},
			expectError: true,
		},
		{
			name:        "Error - expired token",
			token:       expiredTokenString,
			setupMock:   func() {},
			expectError: true,
		},
		{
			name:  "Error - GetUserByLogin fails",
			token: validToken,
			setupMock: func() {
				mockRepo.EXPECT().
					GetUserByLogin(gomock.Any(), login).
					Return(models.User{}, errors.New("user not found"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.ValidateAndGetUser(testContext(), tt.token)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, user, result)
			}
		})
	}
}

func TestAuthUsecase_Enable2FA(t *testing.T) {
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
		has2FA      bool
		setupMock   func()
		expectError bool
		errorType   error
	}{
		{
			name:   "Success",
			has2FA: false,
			setupMock: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(user, nil)
				mockRepo.EXPECT().
					Enable2FA(gomock.Any(), userID, gomock.Any()).
					Return(models.EnableTwoFactorResponse{Has2FA: true}, nil)
			},
			expectError: false,
		},
		{
			name:        "Error - 2FA already enabled",
			has2FA:      true,
			setupMock:   func() {},
			expectError: true,
			errorType:   auth.ErrorBadRequest,
		},
		{
			name:   "Error - GetUserByID fails",
			has2FA: false,
			setupMock: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(models.User{}, errors.New("user not found"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.Enable2FA(testContext(), userID, tt.has2FA)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
			} else {
				assert.NoError(t, err)
				assert.True(t, result.Has2FA)
				assert.NotEmpty(t, result.QrCode)
			}
		})
	}
}

func TestAuthUsecase_VerifyOTPCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	login := "testuser"
	secretCode := "JBSWY3DPEHPK3PXP"

	tests := []struct {
		name        string
		userCode    string
		expectError bool
	}{
		{
			name:        "Error - invalid OTP code",
			userCode:    "000000",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.VerifyOTPCode(testContext(), login, secretCode, tt.userCode)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthUsecase_GenerateQRCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	login := "testuser"

	t.Run("Success", func(t *testing.T) {
		qrCode, secret, err := usecase.GenerateQRCode(login)
		assert.NoError(t, err)
		assert.NotEmpty(t, qrCode)
		assert.NotEmpty(t, secret)
		assert.Greater(t, len(qrCode), 0)
		assert.Greater(t, len(secret), 0)
	})
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
			_, dataIsValid := auth.Validation(tt.login, tt.password)
			assert.Equal(t, tt.expectedValid, dataIsValid)
		})
	}
}

// Новые тесты для методов VK
func TestAuthUsecase_SignInVKUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	vkid := "123456789"
	userID := uuid.NewV4()
	login := "vkuser"
	version := 1

	vkUser := models.User{
		ID:      userID,
		Login:   login,
		Version: version,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
		errorType   error
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetVKUser(gomock.Any(), vkid).
					Return(vkUser, nil)
			},
			expectError: false,
		},
		{
			name: "Error - user not found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetVKUser(gomock.Any(), vkid).
					Return(models.User{}, errors.New("user not found"))
			},
			expectError: true,
			errorType:   auth.ErrorPreconditionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			user, token, err := usecase.SignInVKUser(testContext(), vkid)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Equal(t, vkUser.ID, user.ID)
			}
		})
	}
}

func TestAuthUsecase_SignUpVKUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepo(ctrl)
	usecase := NewAuthUsecase(mockRepo)

	vkid := "123456789"
	login := "vkuser"

	tests := []struct {
		name        string
		setupMock   func()
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
					CreateVKUser(gomock.Any(), gomock.Any(), vkid).
					Return(nil)
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
			expectError: true,
			errorType:   auth.ErrorBadRequest,
		},
		{
			name: "Error - CreateVKUser fails",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserExists(gomock.Any(), login).
					Return(false, nil)
				mockRepo.EXPECT().
					CreateVKUser(gomock.Any(), gomock.Any(), vkid).
					Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			user, token, err := usecase.SignUpVKUser(testContext(), vkid, login)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Equal(t, login, user.Login)
				assert.NotNil(t, user.ID)
				assert.Equal(t, "avatars/default.png", user.Avatar)
			}
		})
	}
}
