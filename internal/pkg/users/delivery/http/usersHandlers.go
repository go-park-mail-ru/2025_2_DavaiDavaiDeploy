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
	CookieName = "DDFilmsJWT"
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

		ctx := context.WithValue(r.Context(), "user_id", user.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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

	helpers.WriteJSON(w, neededUser)
}

func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
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

	user, token, err := u.uc.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	helpers.WriteJSON(w, user)
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
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		helpers.WriteError(w, 401, errors.New("user not authenticated"))
		return
	}

	const maxRequestBodySize = 10 * 1024 * 1024
	limitedReader := http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer limitedReader.Close()

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
	defer newReq.MultipartForm.RemoveAll()

	file, _, err := newReq.FormFile("avatar")
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}
	defer file.Close()

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

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	helpers.WriteJSON(w, user)
}
