package repo

import (
	"crypto/rand"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/hash"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	Users map[uuid.UUID]models.User
)

func init() {
	Users = make(map[uuid.UUID]models.User)
	id1 := uuid.NewV4()
	salt := make([]byte, 8)
	rand.Read(salt)
	Users[id1] = models.User{
		ID:           id1,
		Login:        "ivanov",
		PasswordHash: hash.HashPass(salt, "password123"),
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
				ID:        uuid.NewV4(),
				Title:     "Фантастика",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.NewV4(),
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
