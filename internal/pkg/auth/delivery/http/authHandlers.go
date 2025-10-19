package authHandlers

import (
	"context"
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/repo"
	"kinopoisk/internal/pkg/auth/usecase"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

const (
	CookieName = "DDFilmsJWT"
)

type AuthHandler struct {
	JWTSecret      string
	CookieSecure   bool
	CookieSamesite http.SameSite
	authRepo       *repo.AuthRepository
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
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

	authRepo := repo.NewAuthRepository(db)

	return &AuthHandler{
		JWTSecret:      os.Getenv("JWT_SECRET"),
		CookieSecure:   secure,
		CookieSamesite: samesite,
		authRepo:       authRepo,
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	exists, _ := a.authRepo.CheckUserExists(r.Context(), req.Login)

	if exists {
		errorResp := models.Error{
			Message: "User already exists",
		}
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if msg, passwordIsValid := usecase.ValidatePassword(req.Password); !passwordIsValid {
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

	if msg, loginIsValid := usecase.ValidateLogin(req.Login); !loginIsValid {
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

	passwordHash := usecase.HashPass(req.Password)

	id := uuid.NewV4()

	user := models.User{
		ID:           id,
		Login:        req.Login,
		Version:      1,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	err = a.authRepo.CreateUser(r.Context(), &user)

	authService := usecase.NewAuthService(a.JWTSecret, a.authRepo)
	token, err := authService.GenerateToken(req.Login)
	if err != nil {
		errorResp := models.Error{
			Message: "токен не вышел :/",
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

	neededUser, err := a.authRepo.CheckUserLogin(r.Context(), req.Login)

	if neededUser.ID == uuid.Nil {
		errorResp := models.Error{
			Message: "Wrong login or password",
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if !usecase.CheckPass(neededUser.PasswordHash, req.Login) {
		errorResp := models.Error{
			Message: "Wrong login or password",
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	authService := usecase.NewAuthService(a.JWTSecret, a.authRepo)
	token, err := authService.GenerateToken(req.Login)
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
	err = json.NewEncoder(w).Encode(neededUser)
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
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authService := usecase.NewAuthService(a.JWTSecret, a.authRepo)
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

		if int64(claims["exp"].(float64)) < time.Now().Unix() {
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

		user, err := authService.GetUser(login)
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
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
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
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *AuthHandler) LogOutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		errorResp := models.Error{
			Message: "User not authenticated",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}
	err := a.authRepo.IncrementUserVersion(r.Context(), user.ID)
	if err != nil {
		errorResp := models.Error{
			Message: "Failed to update user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
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
