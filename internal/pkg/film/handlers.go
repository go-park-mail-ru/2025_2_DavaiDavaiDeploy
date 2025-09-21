package film

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

type FilmHandler struct {
}

func NewFilmHandler() *FilmHandler {
	return &FilmHandler{}
}

func (c *FilmHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	users := []models.User{
		{
			ID:           uuid.NewV4(),
			Login:        "ivanov",
			PasswordHash: "password123",
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
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *FilmHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	users := []models.User{
		{
			ID:           uuid.NewV4(),
			Login:        "ivanov",
			PasswordHash: "password123",
			Status:       "active",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *FilmHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	users := []models.User{
		{
			ID:           uuid.NewV4(),
			Login:        "ivanov",
			PasswordHash: "password123",
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
					Surname:   "Ди Каприо",
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
			CreatedAt: time.Now().Add(-30 * 24 * time.Hour), // 30 дней назад
			UpdatedAt: time.Now(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *FilmHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		ID:           uuid.NewV4(),
		Login:        "ivanov",
		PasswordHash: "password123",
		Avatar:       "avatar1.jpg",
		Country:      "Russia",
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	users := []models.User{}
	users = append(users, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	films := []models.Film{
		{
			ID:    uuid.NewV4(),
			Title: "film1",
			Genres: []models.Genre{
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
			Year:        2025,
			Country:     "Russia",
			Rating:      10,
			Budget:      1000000,
			Fees:        10000000,
			PremierDate: time.Now().Add(-30 * 24 * time.Hour),
			Duration:    120,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:    uuid.NewV4(),
			Title: "film2",
			Genres: []models.Genre{
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
			Year:        2023,
			Country:     "Russia",
			Rating:      9,
			Budget:      100000,
			Fees:        10000000,
			PremierDate: time.Now().Add(-30 * 24 * time.Hour),
			Duration:    100,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(films)
}

func (c *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	films := []models.Film{
		{
			ID:    uuid.NewV4(),
			Title: "film1",
			Genres: []models.Genre{
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
			Year:        2025,
			Country:     "Russia",
			Rating:      10,
			Budget:      1000000,
			Fees:        10000000,
			PremierDate: time.Now().Add(-30 * 24 * time.Hour),
			Duration:    120,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(films)
}
