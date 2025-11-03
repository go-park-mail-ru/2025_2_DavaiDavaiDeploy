package usecase

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/users"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/argon2"
)

func HashPass(plainPassword string) []byte {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return []byte{}
	}
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func CheckPass(passHash []byte, plainPassword string) bool {
	salt := make([]byte, 8)
	copy(salt, passHash[:8])
	userHash := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	userHashedPassword := append(salt, userHash...)
	return bytes.Equal(userHashedPassword, passHash)
}

type UserUsecase struct {
	secret   string
	userRepo users.UsersRepo
}

func NewUserUsecase(userRepo users.UsersRepo) *UserUsecase {
	return &UserUsecase{
		secret:   os.Getenv("JWT_SECRET"),
		userRepo: userRepo,
	}
}

func (uc *UserUsecase) GenerateToken(id uuid.UUID, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"login": login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(uc.secret))
}

func (uc *UserUsecase) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.secret), nil
	})
}

func (uc *UserUsecase) ValidateAndGetUser(ctx context.Context, token string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if token == "" {
		logger.Error("no token")
		return models.User{}, users.ErrorUnauthorized
	}

	parsedToken, err := uc.ParseToken(token)
	if err != nil || !parsedToken.Valid {
		return models.User{}, users.ErrorUnauthorized
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("invalid claims")
		return models.User{}, users.ErrorUnauthorized
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		logger.Error("invalid exp claim")
		return models.User{}, users.ErrorUnauthorized
	}

	login, ok := claims["login"].(string)
	if !ok || login == "" {
		logger.Error("invalid login claim")
		return models.User{}, users.ErrorUnauthorized
	}

	user, err := uc.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, users.ErrorUnauthorized
	}

	return user, nil
}

func (uc *UserUsecase) GetUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) (models.User, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	neededUser, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, "", err
	}

	if !CheckPass(neededUser.PasswordHash, oldPassword) {
		logger.Error("wrong old password")
		return models.User{}, "", users.ErrorBadRequest
	}

	msg, passwordIsValid := users.Validaton(neededUser.Login, newPassword)
	if !passwordIsValid {
		logger.Error(msg)
		return models.User{}, "", users.ErrorBadRequest
	}

	if newPassword == oldPassword {
		logger.Error("passwords are equal")
		return models.User{}, "", users.ErrorBadRequest
	}

	neededUser.Version += 1

	err = uc.userRepo.UpdateUserPassword(ctx, neededUser.Version, neededUser.ID, HashPass(newPassword))
	if err != nil {
		return models.User{}, "", err
	}

	neededUser.PasswordHash = HashPass(newPassword)
	neededUser.UpdatedAt = time.Now().UTC()

	token, err := uc.GenerateToken(neededUser.ID, neededUser.Login)
	if err != nil {
		return models.User{}, "", err
	}

	return neededUser, token, nil
}

func (uc *UserUsecase) ChangeUserAvatar(ctx context.Context, id uuid.UUID, buffer []byte) (models.User, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	neededUser, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, "", users.ErrorUnauthorized
	}

	fileFormat := http.DetectContentType(buffer)
	var avatarExtension string
	switch fileFormat {
	case "image/jpeg":
		avatarExtension = ".jpg"
	case "image/png":
		avatarExtension = ".png"
	case "image/webp":
		avatarExtension = ".webp"
	default:
		logger.Error("invalid format of file")
		return models.User{}, "", users.ErrorBadRequest
	}

	avatarPath := neededUser.ID.String() + avatarExtension
	neededUser.Avatar = &avatarPath

	avatarsDir := os.Getenv("AVATARS_DIR")

	filePath := filepath.Join(avatarsDir, avatarPath)
	filePath = filepath.ToSlash(filePath)

	err = os.WriteFile(filePath, buffer, 0644)
	if err != nil {
		logger.Error("failed to write file")
		return models.User{}, "", users.ErrorInternalServerError
	}

	neededUser.Version += 1

	err = uc.userRepo.UpdateUserAvatar(ctx, neededUser.Version, neededUser.ID, filePath)
	if err != nil {
		_ = os.Remove(filePath)
		return models.User{}, "", err
	}

	token, err := uc.GenerateToken(neededUser.ID, neededUser.Login)
	if err != nil {
		_ = os.Remove(filePath)
		return models.User{}, "", err
	}

	return neededUser, token, nil
}
