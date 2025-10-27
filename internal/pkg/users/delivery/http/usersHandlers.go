package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/users"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const (
	CookieName     = "DDFilmsJWT"
	CSRFCookieName = "DDFilmsCSRF"
)

type UserHandler struct {
	uc             users.UsersUsecase
	cookieSecure   bool
	cookieSamesite http.SameSite
}

func NewUserHandler(uc users.UsersUsecase) *UserHandler {
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
	return &UserHandler{
		uc:             uc,
		cookieSecure:   secure,
		cookieSamesite: samesite,
	}
}

func (u *UserHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfCookie, err := r.Cookie(CSRFCookieName)
		if err != nil {
			return
		}
		var csrfToken string

		tokenFromHeader := r.Header.Get("X-CSRF-Token")
		if tokenFromHeader != "" {
			csrfToken = tokenFromHeader
		} else {
			tokenFromForm := r.FormValue("csrftoken")
			if tokenFromForm != "" {
				csrfToken = tokenFromForm
			} else {
				return
			}
		}

		if csrfCookie.Value != csrfToken {
			return
		}
		var token string
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
		}

		user, err := u.uc.ValidateAndGetUser(r.Context(), token)
		if err != nil {
			helpers.WriteError(w, 401, err)
			return
		}
		user.Sanitize()
		ctx := context.WithValue(r.Context(), users.UserKey, user.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param        id   path      string  true  "Genre ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/{id} [get]
func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	neededUser, err := u.uc.GetUser(r.Context(), id)
	if err != nil {
		helpers.WriteError(w, 500, err)
		return
	}
	neededUser.Sanitize()
	helpers.WriteJSON(w, neededUser)
}

// ChangePassword godoc
// @Summary Change user password
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.ChangePasswordInput true "Password data (old_password and new_password are required)"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error "User not authenticated"
// @Failure 500 {object} models.Error
// @Router /users/password [put]
func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		helpers.WriteError(w, 401, errors.New("user not authenticated"))
		return
	}

	var req models.ChangePasswordInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}
	req.Sanitize()

	user, token, err := u.uc.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	csrfToken := uuid.NewV4().String()

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    csrfToken,
		HttpOnly: false,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
	user.Sanitize()
	helpers.WriteJSON(w, user)
}

// ChangeAvatar godoc
// @Summary Change user avatar
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar image file (required, max 10MB, formats: jpg, png, webp)"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 413 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /users/avatar [put]
func (u *UserHandler) ChangeAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		helpers.WriteError(w, 401, errors.New("user not authenticated"))
		return
	}

	const maxRequestBodySize = 10 * 1024 * 1024
	limitedReader := http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer func() {
		if limitedReader.Close() != nil {
			_ = limitedReader.Close()

		}
	}()
	newReq := *r
	newReq.Body = limitedReader

	err := newReq.ParseMultipartForm(maxRequestBodySize)
	if err != nil {
		if errors.As(err, new(*http.MaxBytesError)) {
			helpers.WriteError(w, 413, err)
			return
		}
		helpers.WriteError(w, 400, err)
		return
	}
	defer func() {
		if newReq.MultipartForm != nil {
			_ = newReq.MultipartForm.RemoveAll()
		}
	}()

	file, _, err := newReq.FormFile("avatar")
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}
	defer func() {
		if file.Close() != nil {
			_ = file.Close()

		}
	}()

	buffer, err := io.ReadAll(file)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	user, token, err := u.uc.ChangeUserAvatar(r.Context(), userID, buffer)
	if err != nil {
		helpers.WriteError(w, 500, err)
		return
	}

	csrfToken := uuid.NewV4().String()

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    csrfToken,
		HttpOnly: false,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
	user.Sanitize()
	helpers.WriteJSON(w, user)
}
