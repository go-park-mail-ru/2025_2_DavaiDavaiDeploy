package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"kinopoisk/internal/models"
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
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

	neededUser, err := u.uc.GetUser(r.Context(), id)
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
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		errorResp := models.Error{
			Message: "User not authenticated",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	var req models.ChangePasswordInput
	err := json.NewDecoder(r.Body).Decode(&req)
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

	user, token, err := u.uc.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
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
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		errorResp := models.Error{
			Message: "User not authenticated",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	var req models.ChangePasswordInput
	err := json.NewDecoder(r.Body).Decode(&req)
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

	user, token, err := u.uc.ChangeUserAvatar(r.Context(), userID, buffer)

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
