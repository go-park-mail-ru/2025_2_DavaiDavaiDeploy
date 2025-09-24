package film

import (
	"encoding/json"
	"kinopoisk/internal/models"
	storage "kinopoisk/internal/repo"
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

func (c *FilmHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewV4()
	var neededUser *models.User
	for _, user := range storage.Users {
		if user.ID == id {
			neededUser = &user
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(neededUser)
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
	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)
	filmTotal := len(films)

	endingIndex := min(offset+count, filmTotal)

	paginatedFilms := films[offset:endingIndex]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedFilms)
}

func (c *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewV4()

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

	var result *models.Film
	for _, film := range films {
		if film.ID == id {
			result = &film
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (c *FilmHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	neededGenre := "template"

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

	var result []models.Film
	for _, film := range films {
		for _, genre := range film.Genres {
			if neededGenre == genre.Title {
				result = append(result, film)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
