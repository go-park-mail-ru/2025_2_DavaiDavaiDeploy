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
		Login:        login,
		PasswordHash: hash.HashPass("password123"),
		Avatar:       "avatar1.jpg",
		Country:      "Russia",
		Status:       "active",
		SavedFilms: []models.Film{
			{
				ID:        uuid.NewV4(),
				Title:     "Интерстеллар",
				Year:      2014,
				Country:   "США, Великобритания",
				Rating:    8.6,
				Duration:  169,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.NewV4(),
				Title:     "Крестный отец",
				Year:      1972,
				Country:   "США",
				Rating:    9.2,
				Duration:  175,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		FavoriteGenres: []models.Genre{
			{
				ID:        uuid.Must(uuid.NewV4(), nil),
				Title:     "Фантастика",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.Must(uuid.NewV4(), nil),
				Title:     "Драма",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		FavoriteActors: []models.FilmProfessional{
			{
				ID:        uuid.NewV4(),
				Name:      "Леонардо",
				Surname:   "ДиКаприо",
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.NewV4(),
				Name:      "Ярослав",
				Surname:   "Михалёв",
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
