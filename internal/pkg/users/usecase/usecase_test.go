package usecase

import (
	"context"
	"encoding/pem"
	"errors"
	"log/slog"
	"math/big"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"
	"kinopoisk/internal/pkg/users"
	"kinopoisk/internal/pkg/users/mocks"

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

// Тестовый RSA приватный ключ для тестов с неправильным алгоритмом подписи
var testRSAPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwV2WlP7UYVK6oA8H6HkKxmLnD4i4h7G0J2vKj7V9q6Q0K9XQ
3Lt8Z3P6D5g8U0bJ8X5fWm7q2K0Q3K8L9Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9
Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P
3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6
Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X
0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6
G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H
8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9
L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T
3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T3U0K6Q8L9Y1P3L6Y8X0V6G1H8K9L4T3U0
-----END RSA PRIVATE KEY-----`

func TestHashPass(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantLen  int
	}{
		{
			name:     "valid password",
			password: "Password123!",
			wantLen:  40, // 8 bytes salt + 32 bytes hash
		},
		{
			name:     "empty password",
			password: "",
			wantLen:  40,
		},
		{
			name:     "long password",
			password: "VeryLongPassword1234567890!@#$%^&*()",
			wantLen:  40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashPass(tt.password)
			assert.Len(t, got, tt.wantLen)
			// Проверяем что хэши разные для одного пароля
			got2 := HashPass(tt.password)
			assert.NotEqual(t, got, got2) // разные соли должны давать разные хэши
		})
	}
}

func TestCheckPass(t *testing.T) {
	password := "Password123!"
	hash := HashPass(password)

	tests := []struct {
		name          string
		passwordHash  []byte
		plainPassword string
		want          bool
		description   string
	}{
		{
			name:          "correct password",
			passwordHash:  hash,
			plainPassword: password,
			want:          true,
			description:   "правильный пароль должен проходить проверку",
		},
		{
			name:          "wrong password",
			passwordHash:  hash,
			plainPassword: "WrongPassword123!",
			want:          false,
			description:   "неправильный пароль не должен проходить проверку",
		},
		{
			name:          "different hash",
			passwordHash:  HashPass("DifferentPassword123!"),
			plainPassword: password,
			want:          false,
			description:   "хэш от другого пароля не должен проходить проверку",
		},
		{
			name:          "empty password with hash",
			passwordHash:  HashPass(""),
			plainPassword: "",
			want:          true,
			description:   "пустой пароль с его хэшем должен проходить проверку",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPass(tt.passwordHash, tt.plainPassword)
			assert.Equal(t, tt.want, got, tt.description)
		})
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUsersRepo(ctrl)
	mockStorage := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockStorage)

	userID := uuid.NewV4()
	expectedUser := models.User{
		ID:        userID,
		Version:   1,
		Login:     "testuser",
		Avatar:    "avatars/default.png",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	tests := []struct {
		name        string
		mockSetup   func()
		userID      uuid.UUID
		wantUser    models.User
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(expectedUser, nil)
			},
			userID:   userID,
			wantUser: expectedUser,
			wantErr:  false,
		},
		{
			name: "user not found",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(models.User{}, errors.New("not found"))
			},
			userID:      userID,
			wantUser:    models.User{},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := usecase.GetUser(testContext(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUser, got)
		})
	}
}

func TestUserUsecase_GenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("JWT_SECRET", "test_secret_key_for_jwt_tokens_12345")
	mockRepo := mocks.NewMockUsersRepo(ctrl)
	mockStorage := mocks.NewMockStorageRepo(ctrl)

	userID := uuid.NewV4()
	login := "testuser"
	version := 1

	tests := []struct {
		name        string
		envSetup    func()
		id          uuid.UUID
		login       string
		version     int
		wantErr     bool
		errContains string
		checkToken  func(t *testing.T, token string, usecase *UserUsecase) bool
	}{
		{
			name:    "success",
			id:      userID,
			login:   login,
			version: version,
			wantErr: false,
			checkToken: func(t *testing.T, token string, usecase *UserUsecase) bool {
				assert.NotEmpty(t, token)
				parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				if !assert.NoError(t, err) {
					return false
				}
				if !assert.True(t, parsed.Valid) {
					return false
				}

				claims, ok := parsed.Claims.(jwt.MapClaims)
				if !assert.True(t, ok) {
					return false
				}

				assert.Equal(t, login, claims["login"])
				assert.Equal(t, userID.String(), claims["id"])
				assert.Equal(t, float64(version), claims["version"])
				assert.NotNil(t, claims["exp"])

				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldSecret := os.Getenv("JWT_SECRET")
			if tt.envSetup != nil {
				tt.envSetup()
			}
			defer os.Setenv("JWT_SECRET", oldSecret)

			usecase := NewUserUsecase(mockRepo, mockStorage)

			got, err := usecase.GenerateToken(tt.id, tt.login, tt.version)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.checkToken != nil {
					tt.checkToken(t, got, usecase)
				}
			}
		})
	}
}

func TestUserUsecase_ParseToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("JWT_SECRET", "test_secret_key_for_jwt_tokens")
	mockRepo := mocks.NewMockUsersRepo(ctrl)
	mockStorage := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockStorage)

	// Создаем валидный токен для тестов
	userID := uuid.NewV4()
	validToken, _ := usecase.GenerateToken(userID, "testuser", 1)

	tests := []struct {
		name        string
		token       string
		wantErr     bool
		errContains string
		checkToken  func(t *testing.T, token *jwt.Token) bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
			checkToken: func(t *testing.T, token *jwt.Token) bool {
				return assert.True(t, token.Valid)
			},
		},
		{
			name:        "empty token",
			token:       "",
			wantErr:     true,
			errContains: "token contains an invalid number of segments",
		},
		{
			name:        "malformed token",
			token:       "malformed.token",
			wantErr:     true,
			errContains: "token contains an invalid number of segments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecase.ParseToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				if tt.checkToken != nil {
					tt.checkToken(t, got)
				}
			}
		})
	}
}

// Тест для проверки неправильного алгоритма подписи (отдельный тест чтобы избежать паники)
func TestUserUsecase_ParseToken_WrongSigningMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("JWT_SECRET", "test_secret_key_for_jwt_tokens")
	mockRepo := mocks.NewMockUsersRepo(ctrl)
	mockStorage := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockStorage)

	// Создаем минимально валидный RSA ключ для теста
	block, _ := pem.Decode([]byte(testRSAPrivateKey))
	if block == nil {
		t.Skip("Не удалось декодировать тестовый RSA ключ")
	}

	t.Run("wrong signing method", func(t *testing.T) {
		tokenStr := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEiLCJsb2dpbiI6InRlc3R1c2VyIiwidmVyc2lvbiI6MX0.fake_signature_for_rsa"

		parsedToken, err := usecase.ParseToken(tokenStr)
		assert.Error(t, err)
		assert.Nil(t, parsedToken)
		assert.Contains(t, err.Error(), "unexpected signing method")
	})
}

// Вспомогательная функция для создания big.Int из hex
func bigIntFromHex(hex string) *big.Int {
	n := new(big.Int)
	n.SetString(hex, 16)
	return n
}

func TestUserUsecase_ValidateAndGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("JWT_SECRET", "test_secret_key_for_jwt_tokens")
	mockRepo := mocks.NewMockUsersRepo(ctrl)
	mockStorage := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockStorage)

	userID := uuid.NewV4()
	login := "testuser"
	version := 1
	validToken, _ := usecase.GenerateToken(userID, login, version)

	expectedUser := models.User{
		ID:        userID,
		Version:   version,
		Login:     login,
		Avatar:    "avatars/default.png",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	tests := []struct {
		name        string
		token       string
		mockSetup   func()
		wantUser    models.User
		wantErr     bool
		errIs       error
		errContains string
	}{
		{
			name:  "success",
			token: validToken,
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByLogin(gomock.Any(), login).Return(expectedUser, nil)
			},
			wantUser: expectedUser,
			wantErr:  false,
		},
		{
			name:      "empty token",
			token:     "",
			mockSetup: func() {},
			wantUser:  models.User{},
			wantErr:   true,
			errIs:     users.ErrorUnauthorized,
		},
		{
			name:      "invalid token",
			token:     "invalid.token",
			mockSetup: func() {},
			wantUser:  models.User{},
			wantErr:   true,
			errIs:     users.ErrorUnauthorized,
		},
		{
			name: "expired token",
			token: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"id":      userID,
					"login":   login,
					"version": version,
					"exp":     time.Now().Add(-time.Hour).Unix(),
				})
				tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
				return tokenStr
			}(),
			mockSetup: func() {},
			wantUser:  models.User{},
			wantErr:   true,
			errIs:     users.ErrorUnauthorized,
		},
		{
			name:  "user not found",
			token: validToken,
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByLogin(gomock.Any(), login).Return(models.User{}, errors.New("not found"))
			},
			wantUser: models.User{},
			wantErr:  true,
			errIs:    users.ErrorUnauthorized,
		},
		{
			name: "invalid exp claim type",
			token: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"id":      userID,
					"login":   login,
					"version": version,
					"exp":     "not a number",
				})
				tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
				return tokenStr
			}(),
			mockSetup: func() {},
			wantUser:  models.User{},
			wantErr:   true,
			errIs:     users.ErrorUnauthorized,
		},
		{
			name: "invalid login claim",
			token: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"id":      userID,
					"login":   123, // Неправильный тип
					"version": version,
					"exp":     time.Now().Add(time.Hour).Unix(),
				})
				tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
				return tokenStr
			}(),
			mockSetup: func() {},
			wantUser:  models.User{},
			wantErr:   true,
			errIs:     users.ErrorUnauthorized,
		},
		{
			name: "empty login claim",
			token: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"id":      userID,
					"login":   "", // Пустой логин
					"version": version,
					"exp":     time.Now().Add(time.Hour).Unix(),
				})
				tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
				return tokenStr
			}(),
			mockSetup: func() {},
			wantUser:  models.User{},
			wantErr:   true,
			errIs:     users.ErrorUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := usecase.ValidateAndGetUser(testContext(), tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errIs != nil {
					assert.ErrorIs(t, err, tt.errIs)
				}
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUser, got)
		})
	}
}

func TestUserUsecase_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUsersRepo(ctrl)
	mockStorage := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockStorage)

	userID := uuid.NewV4()
	oldPassword := "OldPassword123!"
	newPassword := "NewPassword123!"

	existingUser := models.User{
		ID:           userID,
		Version:      1,
		Login:        "testuser",
		PasswordHash: HashPass(oldPassword),
		Avatar:       "avatars/default.png",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	tests := []struct {
		name        string
		mockSetup   func()
		userID      uuid.UUID
		oldPassword string
		newPassword string
		wantErr     bool
		errIs       error
		errContains string
		checkResult func(t *testing.T, user models.User, token string, err error)
	}{
		{
			name: "success",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
				mockRepo.EXPECT().UpdateUserPassword(gomock.Any(), 2, userID, gomock.Any()).DoAndReturn(
					func(ctx context.Context, version int, id uuid.UUID, hash []byte) error {
						assert.Len(t, hash, 40)                             // Проверяем что хэш правильной длины
						assert.NotEqual(t, existingUser.PasswordHash, hash) // Новый хэш должен отличаться
						return nil
					},
				)
			},
			userID:      userID,
			oldPassword: oldPassword,
			newPassword: newPassword,
			wantErr:     false,
			checkResult: func(t *testing.T, user models.User, token string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Equal(t, userID, user.ID)
				assert.Equal(t, 2, user.Version) // Версия должна увеличиться
			},
		},
		{
			name: "wrong old password",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
			},
			userID:      userID,
			oldPassword: "WrongPassword123!",
			newPassword: newPassword,
			wantErr:     true,
			errIs:       users.ErrorBadRequest,
		},
		{
			name: "user not found",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(models.User{}, errors.New("not found"))
			},
			userID:      userID,
			oldPassword: oldPassword,
			newPassword: newPassword,
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "same passwords",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
			},
			userID:      userID,
			oldPassword: oldPassword,
			newPassword: oldPassword,
			wantErr:     true,
			errIs:       users.ErrorBadRequest,
		},
		{
			name: "invalid new password - too short",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
			},
			userID:      userID,
			oldPassword: oldPassword,
			newPassword: "short",
			wantErr:     true,
			errIs:       users.ErrorBadRequest,
		},
		{
			name: "update password fails",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
				mockRepo.EXPECT().UpdateUserPassword(gomock.Any(), 2, userID, gomock.Any()).Return(errors.New("db error"))
			},
			userID:      userID,
			oldPassword: oldPassword,
			newPassword: newPassword,
			wantErr:     true,
			errContains: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, token, err := usecase.ChangePassword(testContext(), tt.userID, tt.oldPassword, tt.newPassword)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errIs != nil {
					assert.ErrorIs(t, err, tt.errIs)
				}
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, user, token, err)
				}
			}
		})
	}
}

func TestUserUsecase_ChangeUserAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUsersRepo(ctrl)
	mockStorage := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockStorage)

	userID := uuid.NewV4()
	existingUser := models.User{
		ID:        userID,
		Version:   1,
		Login:     "testuser",
		Avatar:    "avatars/default.png",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	userWithCustomAvatar := existingUser
	userWithCustomAvatar.Avatar = "avatars/custom.jpg"

	buffer := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG signature

	tests := []struct {
		name        string
		mockSetup   func()
		userID      uuid.UUID
		buffer      []byte
		fileFormat  string
		wantErr     bool
		errIs       error
		errContains string
		checkResult func(t *testing.T, user models.User, token string, err error)
	}{
		{
			name: "success with default avatar",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
				// Не должно вызывать DeleteAvatar для default.png
				mockStorage.EXPECT().DeleteAvatar(gomock.Any(), "avatars/default.png").Times(0)
				mockStorage.EXPECT().UploadAvatar(gomock.Any(), userID.String(), buffer, "image/png", ".png").Return("avatars/new.png", nil)
				mockRepo.EXPECT().UpdateUserAvatar(gomock.Any(), 1, userID, "avatars/new.png").Return(nil)
			},
			userID:     userID,
			buffer:     buffer,
			fileFormat: "image/png",
			wantErr:    false,
			checkResult: func(t *testing.T, user models.User, token string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Equal(t, "avatars/new.png", user.Avatar)
			},
		},
		{
			name: "success with custom avatar - delete old",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(userWithCustomAvatar, nil)
				mockStorage.EXPECT().DeleteAvatar(gomock.Any(), "avatars/custom.jpg").Return(nil)
				mockStorage.EXPECT().UploadAvatar(gomock.Any(), userID.String(), buffer, "image/jpeg", ".jpg").Return("avatars/new.jpg", nil)
				mockRepo.EXPECT().UpdateUserAvatar(gomock.Any(), 1, userID, "avatars/new.jpg").Return(nil)
			},
			userID:     userID,
			buffer:     buffer,
			fileFormat: "image/jpeg",
			wantErr:    false,
			checkResult: func(t *testing.T, user models.User, token string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.Equal(t, "avatars/new.jpg", user.Avatar)
			},
		},
		{
			name: "user not found",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(models.User{}, errors.New("not found"))
			},
			userID:      userID,
			buffer:      buffer,
			fileFormat:  "image/png",
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "unsupported file format",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
			},
			userID:     userID,
			buffer:     buffer,
			fileFormat: "image/gif",
			wantErr:    true,
			errIs:      users.ErrorBadRequest,
		},
		{
			name: "delete old avatar fails but continues",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(userWithCustomAvatar, nil)
				mockStorage.EXPECT().DeleteAvatar(gomock.Any(), "avatars/custom.jpg").Return(errors.New("delete failed"))
				mockStorage.EXPECT().UploadAvatar(gomock.Any(), userID.String(), buffer, "image/png", ".png").Return("avatars/new.png", nil)
				mockRepo.EXPECT().UpdateUserAvatar(gomock.Any(), 1, userID, "avatars/new.png").Return(nil)
			},
			userID:     userID,
			buffer:     buffer,
			fileFormat: "image/png",
			wantErr:    false,
			checkResult: func(t *testing.T, user models.User, token string, err error) {
				assert.NoError(t, err) // Ошибка удаления старого аватара не должна прерывать процесс
				assert.NotEmpty(t, token)
			},
		},
		{
			name: "upload avatar fails",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
				mockStorage.EXPECT().UploadAvatar(gomock.Any(), userID.String(), buffer, "image/png", ".png").Return("", errors.New("upload failed"))
			},
			userID:     userID,
			buffer:     buffer,
			fileFormat: "image/png",
			wantErr:    true,
			errIs:      users.ErrorInternalServerError,
		},
		{
			name: "update user avatar fails",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
				mockStorage.EXPECT().UploadAvatar(gomock.Any(), userID.String(), buffer, "image/png", ".png").Return("avatars/new.png", nil)
				mockRepo.EXPECT().UpdateUserAvatar(gomock.Any(), 1, userID, "avatars/new.png").Return(errors.New("db error"))
			},
			userID:      userID,
			buffer:      buffer,
			fileFormat:  "image/png",
			wantErr:     true,
			errContains: "db error",
		},
		{
			name: "webp format",
			mockSetup: func() {
				mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
				mockStorage.EXPECT().UploadAvatar(gomock.Any(), userID.String(), buffer, "image/webp", ".webp").Return("avatars/new.webp", nil)
				mockRepo.EXPECT().UpdateUserAvatar(gomock.Any(), 1, userID, "avatars/new.webp").Return(nil)
			},
			userID:     userID,
			buffer:     buffer,
			fileFormat: "image/webp",
			wantErr:    false,
			checkResult: func(t *testing.T, user models.User, token string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "avatars/new.webp", user.Avatar)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			user, token, err := usecase.ChangeUserAvatar(testContext(), tt.userID, tt.buffer, tt.fileFormat)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errIs != nil {
					assert.ErrorIs(t, err, tt.errIs)
				}
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, user, token, err)
				}
			}
		})
	}
}
