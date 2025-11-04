package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"
	"kinopoisk/internal/pkg/users"
	"kinopoisk/internal/pkg/users/mocks"

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

func TestUserUsecase_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUsersRepo(ctrl)
	mockS3Repo := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	userID := uuid.NewV4()
	expectedUser := models.User{
		ID:        userID,
		Version:   1,
		Login:     "testuser",
		Avatar:    nil,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(expectedUser, nil)
		result, err := usecase.GetUser(testContext(), userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(models.User{}, errors.New("not found"))
		result, err := usecase.GetUser(testContext(), userID)
		assert.Error(t, err)
		assert.Equal(t, models.User{}, result)
	})
}

func TestUserUsecase_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUsersRepo(ctrl)
	mockS3Repo := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	userID := uuid.NewV4()
	oldPassword := "oldPassword123"
	newPassword := "newPassword123"

	existingUser := models.User{
		ID:           userID,
		Version:      1,
		Login:        "testuser",
		PasswordHash: HashPass(oldPassword),
		Avatar:       nil,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
		mockRepo.EXPECT().UpdateUserPassword(gomock.Any(), 2, userID, gomock.Any()).Return(nil)
		result, token, err := usecase.ChangePassword(testContext(), userID, oldPassword, newPassword)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, userID, result.ID)
	})

	t.Run("Wrong password", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
		_, _, err := usecase.ChangePassword(testContext(), userID, "wrong", newPassword)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, users.ErrorBadRequest))
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(models.User{}, errors.New("not found"))
		_, _, err := usecase.ChangePassword(testContext(), userID, oldPassword, newPassword)
		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
	})

	t.Run("Same passwords", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
		_, _, err := usecase.ChangePassword(testContext(), userID, oldPassword, oldPassword)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, users.ErrorBadRequest))
	})

	t.Run("Invalid new password", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
		_, _, err := usecase.ChangePassword(testContext(), userID, oldPassword, "short")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, users.ErrorBadRequest))
	})
}

func TestUserUsecase_GenerateAndParseToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUsersRepo(ctrl)
	os.Setenv("JWT_SECRET", "test-secret-key")
	mockS3Repo := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	userID := uuid.NewV4()
	login := "testuser"

	t.Run("Generate and parse success", func(t *testing.T) {
		token, err := usecase.GenerateToken(userID, login)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := usecase.ParseToken(token)
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)
	})

	t.Run("Parse invalid token", func(t *testing.T) {
		_, err := usecase.ParseToken("invalid")
		assert.Error(t, err)
	})
}

func TestUserUsecase_ValidateAndGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUsersRepo(ctrl)
	os.Setenv("JWT_SECRET", "test-secret-key")
	mockS3Repo := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	userID := uuid.NewV4()
	login := "testuser"
	expectedUser := models.User{
		ID:        userID,
		Version:   1,
		Login:     login,
		Avatar:    nil,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	validToken, _ := usecase.GenerateToken(userID, login)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByLogin(gomock.Any(), login).Return(expectedUser, nil)
		result, err := usecase.ValidateAndGetUser(testContext(), validToken)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("Empty token", func(t *testing.T) {
		_, err := usecase.ValidateAndGetUser(testContext(), "")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, users.ErrorUnauthorized))
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByLogin(gomock.Any(), login).Return(models.User{}, errors.New("not found"))
		_, err := usecase.ValidateAndGetUser(testContext(), validToken)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, users.ErrorUnauthorized))
	})

	t.Run("Invalid token", func(t *testing.T) {
		_, err := usecase.ValidateAndGetUser(testContext(), "invalid.token.here")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, users.ErrorUnauthorized))
	})
}

func TestUserUsecase_ChangeUserAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUsersRepo(ctrl)
	os.Setenv("AVATARS_DIR", t.TempDir())
	os.Setenv("JWT_SECRET", "test-secret-key")
	mockS3Repo := mocks.NewMockStorageRepo(ctrl)
	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	// Корректная JPEG сигнатура
	jpegBuffer := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, // JPEG signature
		0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, // JFIF header
		0x01, 0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, // Some image data
	}

	t.Run("Success JPEG", func(t *testing.T) {
		userID := uuid.NewV4()
		existingUser := models.User{
			ID:        userID,
			Version:   1,
			Login:     "testuser",
			Avatar:    nil,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
		mockRepo.EXPECT().UpdateUserAvatar(gomock.Any(), 2, userID, gomock.Any()).Return(nil)
		result, token, err := usecase.ChangeUserAvatar(testContext(), userID, jpegBuffer)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, 2, result.Version)
	})

	t.Run("User not found", func(t *testing.T) {
		userID := uuid.NewV4()
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(models.User{}, errors.New("not found"))
		_, _, err := usecase.ChangeUserAvatar(testContext(), userID, jpegBuffer)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, users.ErrorUnauthorized))
	})

	t.Run("Invalid format", func(t *testing.T) {
		userID := uuid.NewV4()
		existingUser := models.User{
			ID:        userID,
			Version:   1,
			Login:     "testuser",
			Avatar:    nil,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		mockRepo.EXPECT().GetUserByID(gomock.Any(), userID).Return(existingUser, nil)
		_, _, err := usecase.ChangeUserAvatar(testContext(), userID, []byte{0x00, 0x01, 0x02, 0x03})
		assert.Error(t, err)
		assert.True(t, errors.Is(err, users.ErrorBadRequest))
	})
}
