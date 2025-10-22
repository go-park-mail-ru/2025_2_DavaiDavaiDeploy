package usecase

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/users"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/argon2"
)

const (
	ValidChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?`~"
)

func ValidateLogin(login string) (string, bool) {
	if len(login) < 6 || len(login) > 20 {
		return "Invalid login length", false
	}

	for _, char := range login {
		if !strings.ContainsRune(ValidChars, char) {
			return "Login contains invalid characters", false
		}
	}
	return "Ok", true
}

func ValidatePassword(password string) (string, bool) {
	if len(password) < 6 || len(password) > 20 {
		return "Invalid password length", false
	}

	for _, char := range password {
		if !strings.ContainsRune(ValidChars, char) {
			return "Password contains invalid characters", false
		}
	}
	return "Ok", true

}

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
	salt := passHash[:8]
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
	if token == "" {
		return models.User{}, errors.New("user not authenticated")
	}

	parsedToken, err := uc.ParseToken(token)
	if err != nil || !parsedToken.Valid {
		return models.User{}, errors.New("user not authenticated")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return models.User{}, errors.New("user not authenticated")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return models.User{}, errors.New("user not authenticated")
	}

	login, ok := claims["login"].(string)
	if !ok || login == "" {
		return models.User{}, errors.New("user not authenticated")
	}

	user, err := uc.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, errors.New("user not authenticated")
	}

	return user, nil
}

func (uc *UserUsecase) GetUser(ctx context.Context, id uuid.UUID) (models.User, error) {
	user, err := uc.userRepo.GetUserByID(context.Background(), id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string) (models.User, string, error) {
	neededUser, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, "", errors.New("User not authenticated")
	}

	if !CheckPass(neededUser.PasswordHash, oldPassword) {
		return models.User{}, "", errors.New("Wrong password")
	}

	msg, passwordIsValid := ValidatePassword(newPassword)
	if !passwordIsValid {
		return models.User{}, "", errors.New(msg)
	}

	if newPassword == oldPassword {
		return models.User{}, "", errors.New("The passwords should be different")
	}

	neededUser.Version += 1

	err = uc.userRepo.UpdateUserPassword(ctx, neededUser.Version, neededUser.ID, HashPass(newPassword))
	if err != nil {
		return models.User{}, "", errors.New("Failed to update the password")
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
	neededUser, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return models.User{}, "", errors.New("User not authenticated")
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
		return models.User{}, "", errors.New("unsupported image format")
	}

	avatarPath := neededUser.ID.String() + avatarExtension
	neededUser.Avatar = &avatarPath

	avatarsDir := os.Getenv("AVATARS_DIR")

	filePath := filepath.Join(avatarsDir, avatarPath)
	filePath = filepath.ToSlash(filePath)

	err = os.WriteFile(filePath, buffer, 0644)
	if err != nil {
		return models.User{}, "", err
	}

	err = uc.userRepo.UpdateUserAvatar(ctx, neededUser.ID, filePath)
	if err != nil {
		os.Remove(filePath)
		return models.User{}, "", err
	}

	token, err := uc.GenerateToken(neededUser.ID, neededUser.Login)
	if err != nil {
		os.Remove(filePath)
		return models.User{}, "", err
	}

	return neededUser, token, nil
}
