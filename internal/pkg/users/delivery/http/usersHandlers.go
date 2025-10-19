package http

import (
	"encoding/json"
	"errors"
	"io"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/users/repo"
	"kinopoisk/internal/pkg/users/usecase"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

const (
	CookieName = "DDFilmsJWT"
)

type UserHandler struct {
	JWTSecret      string
	CookieSecure   bool
	CookieSamesite http.SameSite
	UserRepo       *repo.UserRepository
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	secure := false
	cookieValue := os.Getenv("COOKIE_SECURE")
	if cookieValue == "true" {
		secure = true
	}

	samesite := http.SameSiteLaxMode
	samesiteValue := os.Getenv("COOKIE_SAMESITE")
	if samesiteValue == "Strict" {
		samesite = http.SameSiteStrictMode
	}

	userRepo := repo.NewUserRepository(db)

	return &UserHandler{
		JWTSecret:      os.Getenv("JWT_SECRET"),
		CookieSecure:   secure,
		CookieSamesite: samesite,
		UserRepo:       userRepo,
	}
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	neededUser, err := u.UserRepo.GetUserByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(neededUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var token string
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie(CookieName)
	{
		if err == nil {
			token = cookie.Value
		}
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	authService := usecase.NewUserService(u.JWTSecret, u.UserRepo)
	parsedToken, err := authService.ParseToken(token)

	if err != nil || !parsedToken.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	login, ok := claims["login"].(string)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	neededUser, err := u.UserRepo.GetUserByLogin(r.Context(), login)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req models.ChangePasswordInput
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if !usecase.CheckPass(neededUser.PasswordHash, req.OldPassword) {
		errorResp := models.Error{
			Message: "Wrong password",
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	msg, passwordIsValid := usecase.ValidatePassword(req.NewPassword)
	if !passwordIsValid {
		errorResp := models.Error{
			Message: msg,
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if req.NewPassword == req.OldPassword {
		errorResp := models.Error{
			Message: "The passwords should be different",
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	neededUser.PasswordHash = usecase.HashPass(req.NewPassword)
	neededUser.UpdatedAt = time.Now().UTC()

	token, err = authService.GenerateToken(login)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = u.UserRepo.UpdateUserPassword(r.Context(), neededUser.ID, neededUser.PasswordHash)

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.CookieSecure,
		SameSite: u.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(neededUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// ChangeAvatar  godoc
// @Summary      Changing the avatar
// @Description  Changing the avatar by user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Error
// @Router       /auth/change/avatar [put]
func (u *UserHandler) ChangeAvatar(w http.ResponseWriter, r *http.Request) {
	var token string
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie(CookieName)
	{
		if err == nil {
			token = cookie.Value
		}
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	authService := usecase.NewUserService(u.JWTSecret, u.UserRepo)
	parsedToken, err := authService.ParseToken(token)

	if err != nil || !parsedToken.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	login, ok := claims["login"].(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	neededUser, err := u.UserRepo.GetUserByLogin(r.Context(), login)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	const maxRequestBodySize = 10 * 1024 * 1024
	limitedReader := http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer limitedReader.Close()

	newReq := *r
	newReq.Body = limitedReader

	err = newReq.ParseMultipartForm(maxRequestBodySize)
	if err != nil {
		if errors.As(err, new(*http.MaxBytesError)) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, _, err := newReq.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileFormat := http.DetectContentType(buffer)
	if fileFormat != "image/jpeg" && fileFormat != "image/png" && fileFormat != "image/webp" {
		errResp := models.Error{
			Message: "Wrong format of file",
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	avatarExtension := ""
	if fileFormat == "image/jpeg" {
		avatarExtension = ".jpg"
	} else if fileFormat == "image/png" {
		avatarExtension = ".png"
	} else {
		avatarExtension = ".webp"
	}

	avatarPath := neededUser.ID.String() + avatarExtension
	neededUser.Avatar = &avatarPath

	avatarsDir := "/opt/static/avatars"

	filePath := avatarsDir + "/" + avatarPath

	err = os.WriteFile(filePath, buffer, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = u.UserRepo.UpdateUserAvatar(r.Context(), neededUser.ID, filePath)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.CookieSecure,
		SameSite: u.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(neededUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
