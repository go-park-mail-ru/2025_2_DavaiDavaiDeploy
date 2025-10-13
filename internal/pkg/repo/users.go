package repo

import (
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/hash"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	Users map[string]models.User
	Mutex sync.RWMutex
)

func InitUsers() {
	Users = make(map[string]models.User)
	login := "ivanov"
	Users[login] = models.User{
		ID:           uuid.NewV4(),
		Version:      1,
		Login:        login,
		PasswordHash: hash.HashPass("password123"),
		Avatar:       "avatar1.jpg",
		Country:      "Russia",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
