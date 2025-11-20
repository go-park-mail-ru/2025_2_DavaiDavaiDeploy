package usecase

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/dgryski/dgoogauth"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"github.com/skip2/go-qrcode"
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

	msg, dataIsValid := auth.Validaton(req.Login, req.Password)
	if !dataIsValid {
		logger.Error(msg)
		return models.User{}, "", auth.ErrorBadRequest
	}

	exists, err := uc.authRepo.CheckUserExists(ctx, req.Login)
	if err != nil {
		return models.User{}, "", err
	}
	if exists {
		logger.Error("user already exists")
		return models.User{}, "", auth.ErrorConflict
	}

	passwordHash := HashPass(req.Password)

	id := uuid.NewV4()
	defaultAvatar := "avatars/default.png"

	user := models.User{
		ID:           id,
		Login:        req.Login,
		PasswordHash: passwordHash,
		Avatar:       defaultAvatar,
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
		return models.User{}, "", auth.ErrorInternalServerError
	}

	return user, token, nil
}

func (uc *AuthUsecase) VerifyOTPCode(ctx context.Context, login, secretCode string, userCode string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	// Создаем конфигурацию OTP
	otpConfig := &dgoogauth.OTPConfig{
		Secret:      base32.StdEncoding.EncodeToString([]byte(secretCode)),
		WindowSize:  30, //тут должно быть маленькое число
		HotpCounter: 0,
	}
	// Проверяем код
	isValid, err := otpConfig.Authenticate(userCode)
	if err != nil || !isValid {
		logger.Error("OTP authentication error: " + err.Error())
		return auth.ErrorBadRequest
	}

	logger.Info("OTP code verified successfully", slog.String("login", login))
	return nil
}

func (uc *AuthUsecase) SignInUser(ctx context.Context, req models.SignInInput) (models.User, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	neededUser, err := uc.authRepo.CheckUserLogin(ctx, req.Login)
	if err != nil {
		return models.User{}, "", err
	}

	secretCode := uc.authRepo.GetUserSecretCode(ctx, neededUser.ID)
	if secretCode == "" {
		if !CheckPass(neededUser.PasswordHash, req.Password) {
			logger.Error("wrong password")
			return models.User{}, "", auth.ErrorBadRequest
		}

		token, err := uc.GenerateToken(neededUser.ID, req.Login)
		if err != nil {
			logger.Error("cannot generate token")
			return models.User{}, "", auth.ErrorInternalServerError
		}

		return neededUser, token, nil
	}
	emptyCode := ""

	if req.Code == &emptyCode {
		logger.Warn("no code given")
		return models.User{}, "", auth.ErrorPreconditionFailed
	}

	if !CheckPass(neededUser.PasswordHash, req.Password) {
		logger.Error("wrong password")
		return models.User{}, "", auth.ErrorBadRequest
	}

	err = uc.VerifyOTPCode(ctx, neededUser.Login, secretCode, *req.Code)
	if err != nil {
		logger.Error("OTP authentication error: " + err.Error())
		return models.User{}, "", auth.ErrorBadRequest
	}

	token, err := uc.GenerateToken(neededUser.ID, req.Login)
	if err != nil {
		logger.Error("cannot generate token")
		return models.User{}, "", auth.ErrorInternalServerError
	}

	return neededUser, token, nil
}

func (uc *AuthUsecase) LogOutUser(ctx context.Context, userID uuid.UUID) error {
	err := uc.authRepo.IncrementUserVersion(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AuthUsecase) GenerateQRCode(login string) ([]byte, string, error) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return []byte{}, "", err
	}

	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	issuer := "kinopoisk"
	otpURL := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		url.PathEscape(issuer),
		url.PathEscape(login),
		secretBase32,
		url.PathEscape(issuer))

	qrCode, err := qrcode.Encode(otpURL, qrcode.Medium, 256)
	if err != nil {
		return []byte{}, "", err
	}

	return qrCode, secretBase32, nil
}

func (uc *AuthUsecase) Enable2FA(ctx context.Context, userID uuid.UUID, has2FA bool) (models.EnableTwoFactorResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	if has2FA == true {
		logger.Error("user already enabled the 2fa")
		return models.EnableTwoFactorResponse{}, auth.ErrorBadRequest
	}

	user, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error("failed to get user by ID: " + err.Error())
		return models.EnableTwoFactorResponse{}, err
	}

	qrCode, secret, err := uc.GenerateQRCode(user.Login)
	if err != nil {
		logger.Error("failed to generate QR code: " + err.Error())
		return models.EnableTwoFactorResponse{}, auth.ErrorInternalServerError
	}

	response, err := uc.authRepo.Enable2FA(ctx, userID, secret)
	if err != nil {
		logger.Error("failed to enable 2FA: " + err.Error())
		return models.EnableTwoFactorResponse{}, err
	}

	response.QrCode = qrCode
	return response, nil
}

func (uc *AuthUsecase) Disable2FA(ctx context.Context, userID uuid.UUID, has2FA bool) (models.DisableTwoFactorResponse, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	if has2FA == false {
		logger.Error("user already disabled the 2fa")
		return models.DisableTwoFactorResponse{}, auth.ErrorBadRequest
	}
	return uc.authRepo.Disable2FA(ctx, userID)
}

func (uc *AuthUsecase) ValidateAndGetUser(ctx context.Context, token string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if token == "" {
		logger.Error("user is not authorized")
		return models.User{}, auth.ErrorUnauthorized
	}

	parsedToken, err := uc.ParseToken(token)
	if err != nil || !parsedToken.Valid {
		logger.Error("user is not authorized or invalid token")
		return models.User{}, auth.ErrorUnauthorized
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("invalid claims")
		return models.User{}, auth.ErrorUnauthorized
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		logger.Error("invalid exp claim")
		return models.User{}, auth.ErrorUnauthorized
	}

	login, ok := claims["login"].(string)
	if !ok || login == "" {
		logger.Error("invalid login claim")
		return models.User{}, auth.ErrorUnauthorized
	}

	user, err := uc.authRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
