package film

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"math"
	"net/http"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

type FilmHandler struct {
}

func NewFilmHandler() *FilmHandler {
	return &FilmHandler{}
}

func GetParameter(r *http.Request, s string, defaultValue int) int {
	strValue := r.URL.Query().Get(s)
	if strValue == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(strValue)
	if err != nil || result <= 0 {
		return defaultValue
	}
	return result
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
			Title: "Интерстеллар",
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
				{
					ID:        uuid.NewV4(),
					Title:     "Приключения",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			Year:        2014,
			Country:     "США",
			Rating:      8.6,
			Budget:      165000000,
			Fees:        677000000,
			PremierDate: time.Date(2014, 10, 26, 0, 0, 0, 0, time.UTC),
			Duration:    169,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:    uuid.NewV4(),
			Title: "Крестный отец",
			Genres: []models.Genre{
				{
					ID:        uuid.NewV4(),
					Title:     "Криминал",
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
			Year:        1972,
			Country:     "США",
			Rating:      9.2,
			Budget:      6000000,
			Fees:        245000000,
			PremierDate: time.Date(1972, 3, 15, 0, 0, 0, 0, time.UTC),
			Duration:    175,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:    uuid.NewV4(),
			Title: "Темный рыцарь",
			Genres: []models.Genre{
				{
					ID:        uuid.NewV4(),
					Title:     "Боевик",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        uuid.NewV4(),
					Title:     "Криминал",
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
			Year:        2008,
			Country:     "США",
			Rating:      9.0,
			Budget:      185000000,
			Fees:        1005000000,
			PremierDate: time.Date(2008, 7, 18, 0, 0, 0, 0, time.UTC),
			Duration:    152,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:    uuid.NewV4(),
			Title: "Брат",
			Genres: []models.Genre{
				{
					ID:        uuid.NewV4(),
					Title:     "Криминал",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        uuid.NewV4(),
					Title:     "Драма",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        uuid.NewV4(),
					Title:     "Боевик",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			Year:        1997,
			Country:     "Россия",
			Rating:      8.3,
			Budget:      10000,
			Fees:        1000000,
			PremierDate: time.Date(1997, 12, 12, 0, 0, 0, 0, time.UTC),
			Duration:    100,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:    uuid.NewV4(),
			Title: "Назад в будущее",
			Genres: []models.Genre{
				{
					ID:        uuid.NewV4(),
					Title:     "Фантастика",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        uuid.NewV4(),
					Title:     "Комедия",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        uuid.NewV4(),
					Title:     "Приключения",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			Year:        1985,
			Country:     "США",
			Rating:      8.5,
			Budget:      19000000,
			Fees:        381000000,
			PremierDate: time.Date(1985, 7, 3, 0, 0, 0, 0, time.UTC),
			Duration:    116,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:    uuid.NewV4(),
			Title: "Леон",
			Genres: []models.Genre{
				{
					ID:        uuid.NewV4(),
					Title:     "Боевик",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        uuid.NewV4(),
					Title:     "Криминал",
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
			Year:        1994,
			Country:     "Франция",
			Rating:      8.5,
			Budget:      16000000,
			Fees:        45000000,
			PremierDate: time.Date(1994, 9, 14, 0, 0, 0, 0, time.UTC),
			Duration:    110,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
		{
			ID:    uuid.NewV4(),
			Title: "Джентльмены",
			Genres: []models.Genre{
				{
					ID:        uuid.NewV4(),
					Title:     "Криминал",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        uuid.NewV4(),
					Title:     "Комедия",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        uuid.NewV4(),
					Title:     "Боевик",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			Year:        2019,
			Country:     "Великобритания",
			Rating:      8.5,
			Budget:      22000000,
			Fees:        115000000,
			PremierDate: time.Date(2019, 12, 3, 0, 0, 0, 0, time.UTC),
			Duration:    113,
			CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		},
	}
	page := GetParameter(r, "page", 1)
	filmLimit := GetParameter(r, "limit", 10)
	filmTotal := len(films)
	totalPages := int(math.Ceil(float64(filmTotal) / float64(filmLimit)))

	startingIndex := (page - 1) * filmLimit
	endingIndex := min(startingIndex+filmLimit, filmTotal)

	paginatedFilms := films[startingIndex:endingIndex]

	result := map[string]interface{}{
		"data": paginatedFilms,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       filmLimit,
			"total":       filmTotal,
			"totalPages":  totalPages,
			"hasNext":     page < totalPages,
			"hasPrevious": page > 1,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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
