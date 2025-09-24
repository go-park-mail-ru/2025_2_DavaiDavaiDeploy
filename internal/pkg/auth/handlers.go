package auth

import (
	"crypto/rand"
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/hash"
	"kinopoisk/internal/pkg/auth/validation"
	"kinopoisk/internal/pkg/repo"
	"net/http"
	"time"

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
		errorResp := models.Error{
			Type:    "BAD_REQUEST",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	password := req.Password
	err = validation.ValidatePassword(password)
	if err != nil {
		errorResp := models.Error{
			Type:    "VALIDATION_ERROR",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	login := req.Login

	err = validation.ValidateLogin(login)
	if err != nil {
		errorResp := models.Error{
			Type:    "VALIDATION_ERROR",
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	salt := make([]byte, 8)
	rand.Read(salt)
	passwordHash := hash.HashPass(salt, password)
	if !hash.CheckPass(passwordHash, password) {
		errorResp := models.Error{
			Type:    "VALIDATION_ERROR",
			Message: "hashing error",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	avatar := req.Avatar
	if avatar == "" {
		avatar = "avatar1.jpg"
	}

	country := req.Country
	if country == "" {
		country = "Russia"
	}

	id := uuid.NewV4()

	user := models.User{
		ID:           id,
		Login:        login,
		PasswordHash: passwordHash,
		Avatar:       avatar,
		Country:      country,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	repo.Users[login] = user

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *AuthHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	var req models.SignInInput
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		errorResp := models.Error{
			Type:    "BAD_REQUEST",
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
			Type:    "NOT_FOUND",
			Message: "user not found",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	if !hash.CheckPass(neededUser.PasswordHash, enteredPassword) {
		errorResp := models.Error{
			Type:    "VALIDATION_ERROR",
			Message: "wrong password",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(neededUser)
}
