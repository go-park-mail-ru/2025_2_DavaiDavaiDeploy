package authHandlers

import (
	"context"
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/helpers"
	"net/http"
	"os"
	"time"
)

const (
	CookieName = "DDFilmsJWT"
)

type AuthHandler struct {
	JWTSecret      string
	CookieSecure   bool
	CookieSamesite http.SameSite
	uc             auth.AuthUsecase
}

func NewAuthHandler(uc auth.AuthUsecase) *AuthHandler {
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

	return &AuthHandler{
		JWTSecret:      os.Getenv("JWT_SECRET"),
		CookieSecure:   secure,
		CookieSamesite: samesite,
		uc:             uc,
	}
}

// SignupUser godoc
// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.SignUpInput true "User registration data"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /auth/signup [post]
func (a *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	user, token, err := a.uc.SignUpUser(r.Context(), req)

	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	//w.Header().Set("Authorization", "Bearer "+token)
	helpers.WriteJSON(w, user)
}

// SignInUser godoc
// @Summary User login
// @Description Authenticate user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.SignInInput true "User data"
// @Success 200 {object} models.User
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /auth/signin [post]
func (a *AuthHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var req models.SignInInput
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	user, token, err := a.uc.SignInUser(r.Context(), req)

	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	//w.Header().Set("Authorization", "Bearer "+token)
	helpers.WriteJSON(w, user)
}

func (a *AuthHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
		}

		user, err := a.uc.ValidateAndGetUser(r.Context(), token)
		if err != nil {
			helpers.WriteError(w, 401, err)
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckAuth godoc
// @Summary Check authentication status
// @Description Verify if user is authenticated and return user data
// @Tags auth
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /auth/check [get]
func (a *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	user, err := a.uc.CheckAuth(r.Context())
	if err != nil {
		helpers.WriteError(w, 401, err)
	}

	helpers.WriteJSON(w, user)
}

// LogOutUser godoc
// @Summary User logout
// @Description Clear authentication cookie and log out user
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} models.Error
// @Router /auth/logout [post]
func (a *AuthHandler) LogOutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := a.uc.LogOutUser(r.Context())
	if err != nil {
		helpers.WriteError(w, 401, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(-12 * time.Hour),
		Path:     "/",
	})

	helpers.WriteJSON(w, map[string]string{"message": "Successfully logged out"})
}
