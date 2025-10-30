package authHandlers

import (
	"context"
	"encoding/json"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	CookieName     = "DDFilmsJWT"
	CSRFCookieName = "DDFilmsCSRF"
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
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	var req models.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid input"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	req.Sanitize()

	user, token, err := a.uc.SignUpUser(r.Context(), req)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	csrfToken := uuid.NewV4().String()

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    csrfToken,
		HttpOnly: false,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
	user.Sanitize()

	log.LogHandlerInfo(logger, "Success", http.StatusOK)
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
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	var req models.SignInInput
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid input"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	req.Sanitize()

	user, token, err := a.uc.SignInUser(r.Context(), req)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	csrfToken := uuid.NewV4().String()

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    csrfToken,
		HttpOnly: false,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
	user.Sanitize()
	//w.Header().Set("Authorization", "Bearer "+token)
	helpers.WriteJSON(w, user)

	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}

func (a *AuthHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
		csrfCookie, err := r.Cookie(CSRFCookieName)
		if err != nil {
			log.LogHandlerError(logger, errors.New("invalid csrf-token"), http.StatusUnauthorized)
			helpers.WriteError(w, http.StatusUnauthorized, errors.New("user is not authorized"))
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
				log.LogHandlerError(logger, errors.New("csrf-token is empty"), http.StatusUnauthorized)
				helpers.WriteError(w, http.StatusUnauthorized, errors.New("user is not authorized"))
				return
			}
		}
		if csrfCookie.Value != csrfToken {
			log.LogHandlerError(logger, errors.New("invalid csrf-token"), http.StatusUnauthorized)
			helpers.WriteError(w, http.StatusUnauthorized, errors.New("user is not authorized"))
			return
		}

		var token string
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
		}

		user, err := a.uc.ValidateAndGetUser(r.Context(), token)
		if err != nil {
			helpers.WriteError(w, http.StatusUnauthorized, errors.New("user is not authorized"))
			return
		}
		user.Sanitize()
		ctx := context.WithValue(r.Context(), auth.UserKey, user)

		log.LogHandlerInfo(logger, "Success", http.StatusOK)
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
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	user, err := a.uc.CheckAuth(r.Context())
	if err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err)
	}
	user.Sanitize()
	helpers.WriteJSON(w, user)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
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
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	err := a.uc.LogOutUser(r.Context())
	if err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    "",
		HttpOnly: false,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(-12 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(-12 * time.Hour),
		Path:     "/",
	})

	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}
