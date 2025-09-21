package auth

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/source"
	"kinopoisk/internal/storage"
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
	if !source.ValidatePassword(password) {
		http.Error(w, `{"error": "password is invalid"}`, http.StatusBadRequest)
		return
	}

	login := "ivanova"
	if !source.ValidateLogin(login) {
		http.Error(w, `{"error": "login is invalid"}`, http.StatusBadRequest)
		return
	}

	passwordHash := source.HashPassword(password)
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

	storage.Users[id.String()] = user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (c *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	enteredLogin := "ivanov"
	enteredPassword := "password123"

	var necessaryUser *models.User
	for _, user := range storage.Users {
		if user.Login == enteredLogin {
			necessaryUser = &user
		}
	}

	if necessaryUser == nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusUnauthorized)
		return
	}

	if necessaryUser.PasswordHash != source.HashPassword(enteredPassword) {
		http.Error(w, `{"error": "password is wrong"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(necessaryUser)
}
