package authHandlers

import (
	"context"
	"encoding/json"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/auth/delivery/grpc/gen"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CookieName     = "DDFilmsJWT"
	CSRFCookieName = "DDFilmsCSRF"
)

type AuthHandler struct {
	JWTSecret      string
	CookieSecure   bool
	CookieSamesite http.SameSite
	client         gen.AuthClient
}

func NewAuthHandler(client gen.AuthClient) *AuthHandler {
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
		client:         client,
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
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /auth/signup [post]
func (a *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	var req models.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid input"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	req.Sanitize()

	user, err := a.client.SignupUser(r.Context(), &gen.SignupRequest{
		Login:    req.Login,
		Password: req.Password})

	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.InvalidArgument:
			helpers.WriteError(w, http.StatusBadRequest)
		case codes.AlreadyExists:
			helpers.WriteError(w, http.StatusConflict)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    user.CSRFToken,
		HttpOnly: false,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    user.JWTToken,
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	w.Header().Set("X-CSRF-Token", user.CSRFToken)
	helpers.WriteJSON(w, user)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// SignInUser godoc
// @Summary User login
// @Description Authenticate user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.SignInInput true "User data"
// @Success 200 {object} models.User
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /auth/signin [post]
func (a *AuthHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	var req models.SignInInput
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid input"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	req.Sanitize()

	user, err := a.client.SignInUser(r.Context(), &gen.SignInRequest{
		Login:    req.Login,
		Password: req.Password})

	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.InvalidArgument:
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    user.CSRFToken,
		HttpOnly: false,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    user.JWTToken,
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	w.Header().Set("X-CSRF-Token", user.CSRFToken)
	helpers.WriteJSON(w, user)

	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (a *AuthHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
		csrfCookie, err := r.Cookie(CSRFCookieName)
		if err != nil {
			log.LogHandlerError(logger, errors.New("invalid csrf token"), http.StatusUnauthorized)
			helpers.WriteError(w, http.StatusUnauthorized)
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
				helpers.WriteError(w, http.StatusUnauthorized)
				return
			}
		}
		if csrfCookie.Value != csrfToken {
			log.LogHandlerError(logger, errors.New("invalid csrf-token"), http.StatusUnauthorized)
			helpers.WriteError(w, http.StatusUnauthorized)
			return
		}

		var token string
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
		}

		user, err := a.client.ValidateAndGetUser(r.Context(), &gen.ValidateAndGetUserRequest{Token: token})
		if err != nil {
			st, _ := status.FromError(err)
			switch st.Code() {
			case codes.Unauthenticated:
				helpers.WriteError(w, http.StatusUnauthorized)
			default:
				helpers.WriteError(w, http.StatusInternalServerError)
			}
		}
		ctx := context.WithValue(r.Context(), auth.UserKey, user)

		log.LogHandlerInfo(logger, "success", http.StatusOK)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckAuth godoc
// @Summary Check authentication status
// @Description Verify if user is authenticated and return user data
// @Tags auth
// @Produce json
// @Success 200 {object} models.User
// @Failure 401
// @Failure 500
// @Router /auth/check [get]
func (a *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	user, err := a.client.CheckAuth(r.Context(), &gen.CheckAuthRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.Unauthenticated:
			helpers.WriteError(w, http.StatusUnauthorized)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
	}
	helpers.WriteJSON(w, user)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// LogOutUser godoc
// @Summary User logout
// @Description Clear authentication cookie and log out user
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401
// @Failure 500
// @Router /auth/logout [post]
func (a *AuthHandler) LogOutUser(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	_, err := a.client.LogOutUser(r.Context(), &gen.LogOutUserRequest{})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.Unauthenticated:
			helpers.WriteError(w, http.StatusUnauthorized)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
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

	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
