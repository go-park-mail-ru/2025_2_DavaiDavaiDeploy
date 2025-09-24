package repo

import (
	"crypto/rand"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/hash"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	Users map[string]models.User
)

func init() {
	Users = make(map[string]models.User)
	login := "ivanov"
	salt := make([]byte, 8)
	rand.Read(salt)
	Users[login] = models.User{
		ID:           uuid.NewV4(),
		Login:        login,
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
				ID:        uuid.Must(uuid.FromString("w3x4y5z6-a7b8-9012-wxyz-345678901234")),
				Title:     "Фантастика",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.Must(uuid.FromString("g7h8i9j0-k1l2-3456-ghij-789012345678")),
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
