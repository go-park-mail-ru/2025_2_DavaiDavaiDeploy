package authHandlers

import (
	"context"
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
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
// @Summary      User signup
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.SignUpInput  true  "User credentials"
// @Success      200    {object}  models.User
// @Failure      400    {object}  models.Error
// @Failure      409    {object}  models.Error
// @Router       /auth/signup [post]
func (a *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var req models.SignUpInput
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

	user, token, err := a.uc.SignUpUser(r.Context(), req)

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
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// SignInUser godoc
// @Summary      User login
// @Description  Authenticate existing user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.SignInInput  true  "User credentials"
// @Success      200    {object}  models.User
// @Failure      400    {object}  models.Error
// @Failure      401    {object}  models.Error
// @Router       /auth/signin [post]
func (a *AuthHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var req models.SignInInput
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

	user, token, err := a.uc.SignInUser(r.Context(), req)

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
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckAuth godoc
// @Summary      Check authentication
// @Description  Verify JWT token in cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      401  {object}  models.Error
// @Router       /auth/check [get]
func (a *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	user, err := a.uc.CheckAuth(r.Context())
	if err != nil {
		errorResp := models.Error{
			Message: "User not authenticated",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *AuthHandler) LogOutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := a.uc.LogOutUser(r.Context())
	if err != nil {
		errorResp := models.Error{
			Message: "User not authenticated",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
}
