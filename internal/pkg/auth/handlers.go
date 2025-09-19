package auth

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

type FilmController struct {
}

func NewFilmController() (*FilmController, error) {
	return &FilmController{}, nil
}

func (c *FilmController) SignupUser(w http.ResponseWriter, r *http.Request) {
	users := []models.User{
		{
			ID:               uuid.NewV4(),
			Login:            "ivanov",
			Password:         "password123",
			Avatar:           "avatar1.jpg",
			Country:          "Russia",
			Status:           "active",
			SavedFilmsID:     []uuid.UUID{},
			FavoriteGenresID: []uuid.UUID{},
			FavoriteActorsID: []uuid.UUID{},
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *FilmController) LoginUser(w http.ResponseWriter, r *http.Request) {
	users := []models.User{
		{
			ID:       uuid.NewV4(),
			Login:    "ivanov",
			Password: "password123",
			Status:   "active",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *FilmController) GetUser(w http.ResponseWriter, r *http.Request) {
	users := []models.User{
		{
			ID:               uuid.NewV4(),
			Login:            "ivanov",
			Password:         "password123",
			Avatar:           "avatar1.jpg",
			Country:          "Russia",
			Status:           "active",
			SavedFilmsID:     []uuid.UUID{uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4()},
			FavoriteGenresID: []uuid.UUID{uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4()},
			FavoriteActorsID: []uuid.UUID{uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4()},
			CreatedAt:        time.Now().Add(-30 * 24 * time.Hour), // 30 дней назад
			UpdatedAt:        time.Now(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *FilmController) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		ID:               uuid.NewV4(),
		Login:            "ivanov",
		Password:         "password123",
		Avatar:           "avatar1.jpg",
		Country:          "Russia",
		Status:           "active",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	users := []models.User {}
	users = append(users, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *FilmController) GetFilms(w http.ResponseWriter, r *http.Request) {
	films := []models.Film{
		{
			ID:         	uuid.NewV4(),
			Title:			"film1",
			Genres:			"genre1",
			Year:			2025,
			Country:    	"Russia",
			Rating:			10,
			Budget:			1000000,
			Fees:			10000000,
			PremierDate: 	time.Now().Add(-30 * 24 * time.Hour),
			Duration:		120,
			CreatedAt:		time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt: 		time.Now(),
		},
		{
			ID:         	uuid.NewV4(),
			Title:			"film2",
			Genres:			"genre2",
			Year:			2023,
			Country:    	"Russia",
			Rating:			9,
			Budget:			100000,
			Fees:			10000000,
			PremierDate: 	time.Now().Add(-30 * 24 * time.Hour),
			Duration:		100,
			CreatedAt:		time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt: 		time.Now(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(films)
}

func (c *FilmController) GetFilm(w http.ResponseWriter, r *http.Request) {
	films := []models.Film{
		{
			ID:         	uuid.NewV4(),
			Title:			"film1",
			Genres:			"genre1",
			Year:			2025,
			Country:    	"Russia",
			Rating:			10,
			Budget:			1000000,
			Fees:			10000000,
			PremierDate: 	time.Now().Add(-30 * 24 * time.Hour),
			Duration:		120,
			CreatedAt:		time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt: 		time.Now(),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(films)
}
