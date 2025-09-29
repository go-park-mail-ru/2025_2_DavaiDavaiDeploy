package authHandlers

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/hash"
	"kinopoisk/internal/pkg/auth/service"
	"kinopoisk/internal/pkg/auth/validation"
	"kinopoisk/internal/pkg/repo"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const (
	CookieName = "DDFilmsJWT"
)

type AuthHandler struct {
	JWTSecret string
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}

func (a *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	var req models.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, exists := repo.Users[req.Login]

	if exists {
		errorResp := models.Error{
			Message: "User already exists",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	if msg, passwordIsValid := validation.ValidatePassword(req.Password); !passwordIsValid {
		errorResp := models.Error{
			Message: msg,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	if msg, loginIsValid := validation.ValidateLogin(req.Login); !loginIsValid {
		errorResp := models.Error{
			Message: msg,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	passwordHash := hash.HashPass(req.Password)

	id := uuid.NewV4()

	user := models.User{
		ID:           id,
		Login:        req.Login,
		PasswordHash: passwordHash,
		Status:       "active",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	repo.Users[req.Login] = user
	authService := service.NewAuthService(a.JWTSecret)
	token, err := authService.GenerateToken(req.Login)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (a *AuthHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	var req models.SignInInput
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	enteredLogin := req.Login
	enteredPassword := req.Password

	var neededUser models.User
	for i, user := range repo.Users {
		if user.Login == enteredLogin {
			neededUser = repo.Users[i]
			break
		}
	}

	if neededUser.ID == uuid.Nil {
		errorResp := models.Error{
			Message: "Wrong login or password",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	if !hash.CheckPass(neededUser.PasswordHash, enteredPassword) {
		errorResp := models.Error{
			Message: "Wrong login or password",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)

		return
	}

	authService := service.NewAuthService(a.JWTSecret)
	token, err := authService.GenerateToken(req.Login)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(neededUser)
}

func (a *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	var neededUser models.User
	for i, user := range repo.Users {
		if user.ID == id {
			neededUser = repo.Users[i]
		}
	}

	if neededUser.ID == uuid.Nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(neededUser)
}

func (a *AuthHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authService := service.NewAuthService(a.JWTSecret)
		token := headerParts[1]
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

		next.ServeHTTP(w, r)
	})
}

func (a *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	var token string

	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 6 && authHeader[:6] == "Bearer" {
		token = authHeader[7:]
	}

	if token == "" {
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
		}
	}

	if token == "" {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	authService := service.NewAuthService(a.JWTSecret)

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

	user, err := authService.GetUser(login)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
