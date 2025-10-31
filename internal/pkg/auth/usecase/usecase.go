package usecase

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/argon2"
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

const (
	ValidChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?`~"
)

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
	//salt := passHash[:8] - раньше было так
	salt := make([]byte, 8)
	copy(salt, passHash[:8])
	userHash := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	userHashedPassword := append(salt, userHash...)
	return bytes.Equal(userHashedPassword, passHash)
}

type AuthUsecase struct {
	secret   string
	authRepo auth.AuthRepo
}

func NewAuthUsecase(repo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		authRepo: repo,
		secret:   os.Getenv("JWT_SECRET"),
	}
}

func (uc *AuthUsecase) GenerateToken(id uuid.UUID, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"login": login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(uc.secret))
}

func (uc *AuthUsecase) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.secret), nil
	})
}

func (uc *AuthUsecase) SignUpUser(ctx context.Context, req models.SignUpInput) (models.User, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	if msg, passwordIsValid := ValidatePassword(req.Password); !passwordIsValid {
		logger.Error(msg)
		return models.User{}, "", errors.New(msg)
	}

	if msg, loginIsValid := ValidateLogin(req.Login); !loginIsValid {
		logger.Error(msg)
		return models.User{}, "", errors.New(msg)
	}

	exists, err := uc.authRepo.CheckUserExists(ctx, req.Login)
	if err != nil {
		return models.User{}, "", err
	}
	if exists {
		logger.Error("user already exists")
		return models.User{}, "", errors.New("user already exists")
	}

	passwordHash := HashPass(req.Password)

	id := uuid.NewV4()
	defaultAvatar := "/static/default.jpg"

	user := models.User{
		ID:           id,
		Login:        req.Login,
		PasswordHash: passwordHash,
		Avatar:       &defaultAvatar,
		Version:      1,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	err = uc.authRepo.CreateUser(ctx, user)
	if err != nil {
		return models.User{}, "", err
	}

	token, err := uc.GenerateToken(id, req.Login)
	if err != nil {
		logger.Error("cannot generate token")
		return models.User{}, "", err
	}

	return user, token, nil
}

func (uc *AuthUsecase) SignInUser(ctx context.Context, req models.SignInInput) (models.User, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	neededUser, err := uc.authRepo.CheckUserLogin(ctx, req.Login)

	if err != nil {
		return models.User{}, "", err
	}

	if neededUser.ID == uuid.Nil {
		logger.Error("no such user")
		return models.User{}, "", errors.New("wrong login or password")
	}

	if !CheckPass(neededUser.PasswordHash, req.Password) {
		logger.Error("no such user")
		return models.User{}, "", errors.New("wrong login or password")
	}

	token, err := uc.GenerateToken(neededUser.ID, req.Login)
	if err != nil {
		logger.Error("cannot generate token")
		return models.User{}, "", err
	}

	return neededUser, token, nil
}

func (uc *AuthUsecase) CheckAuth(ctx context.Context) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	user, ok := ctx.Value(auth.UserKey).(models.User)
	if !ok {
		logger.Info("no such user in context")
		return models.User{}, errors.New("user is not authorized")
	}

	return user, nil
}

func (uc *AuthUsecase) LogOutUser(ctx context.Context) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	user, ok := ctx.Value(auth.UserKey).(models.User)
	if !ok {
		logger.Error("no such user in context")
		return errors.New("user is not authorized")
	}
	err := uc.authRepo.IncrementUserVersion(ctx, user.ID)
	if err != nil {
		logger.Error("cannot increment version")
		return errors.New("user is not authorized")
	}

	return nil
}

func (uc *AuthUsecase) ValidateAndGetUser(ctx context.Context, token string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if token == "" {
		logger.Error("user is not authorized")
		return models.User{}, errors.New("user is not authorized")
	}

	parsedToken, err := uc.ParseToken(token)
	if err != nil || !parsedToken.Valid {
		logger.Error("user is not authorized or invalid token")
		return models.User{}, errors.New("user is not authorized")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("invalid claims")
		return models.User{}, errors.New("user not authenticated")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		logger.Error("invalid exp claim")
		return models.User{}, errors.New("user not authenticated")
	}

	login, ok := claims["login"].(string)
	if !ok || login == "" {
		logger.Error("invalid login claim")
		return models.User{}, errors.New("user not authenticated")
	}

	user, err := uc.authRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, errors.New("user not authenticated")
	}
	return user, nil
}
