package auth

import (
	"crypto/rand"
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/hash"
	"kinopoisk/internal/pkg/auth/validation"
	"kinopoisk/internal/repo"
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
	id := uuid.NewV4()

	password := "password456"
	err := validation.ValidatePassword(password)
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

	login := "ivanova"

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


	user := models.User{
		ID:           id,
		Login:        login,
		PasswordHash: passwordHash,
		Avatar:       "avatar1.jpg",
		Country:      "Russia",
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	repo.Users[id] = user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (c *AuthHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	enteredLogin := "ivanov"
	enteredPassword := "password123"

	var neededUser *models.User
	for _, user := range repo.Users {
		if user.Login == enteredLogin {
			neededUser = &user
			break
		}
	}

	if neededUser == nil {
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
