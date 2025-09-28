package authHandlers

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/hash"
	"kinopoisk/internal/pkg/auth/service"
	"kinopoisk/internal/pkg/auth/validation"
	"kinopoisk/internal/pkg/repo"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (c *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
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
			Message: "user already exists",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	err = validation.ValidatePassword(req.Password)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	err = validation.ValidateLogin(req.Login)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
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
	authService := service.NewAuthService("sus")
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

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

	http.SetCookie(w, &http.Cookie{
		Name:     "AdminJWT",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
}

func (c *AuthHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
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
			Message: "wrong login or password",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	if !hash.CheckPass(neededUser.PasswordHash, enteredPassword) {
		errorResp := models.Error{
			Message: "wrong login or password",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)

		return
	}
	authService := service.NewAuthService("sus")
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

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(neededUser)

	http.SetCookie(w, &http.Cookie{
		Name:     "AdminJWT",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
}

func (c *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
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
		errorResp := models.Error{
			Message: "user not found",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(neededUser)
}
