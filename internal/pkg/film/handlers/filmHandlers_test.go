package filmHandlers

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/repo"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRepo() {
	repo.InitRepo()

	genre1 := models.Genre{ID: uuid.NewV4(), Title: "Action"}
	genre2 := models.Genre{ID: uuid.NewV4(), Title: "Drama"}
	repo.Genres = []models.Genre{genre1, genre2}

	film1 := models.Film{ID: uuid.NewV4(), Title: "Film One", GenreID: genre1.ID}
	film2 := models.Film{ID: uuid.NewV4(), Title: "Film Two", GenreID: genre2.ID}
	repo.Films = []models.Film{film1, film2}
}

func TestGetFilms(t *testing.T) {
	setupRepo()
	handler := NewFilmHandler()

	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedLen  int
	}{
		{
			name:         "Default pagination",
			url:          "/api/films",
			expectedCode: http.StatusOK,
			expectedLen:  2,
		},
		{
			name:         "Pagination with count=1",
			url:          "/api/films?count=1&offset=0",
			expectedCode: http.StatusOK,
			expectedLen:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", tt.url, nil)

			handler.GetFilms(w, r)
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var films []models.Film
				err := json.NewDecoder(w.Body).Decode(&films)
				require.NoError(t, err)
				assert.Len(t, films, tt.expectedLen)
			}
		})
	}
}

func TestGetFilm(t *testing.T) {
	setupRepo()
	handler := NewFilmHandler()
	existingFilm := repo.Films[0]

	tests := []struct {
		name         string
		id           string
		expectedCode int
		expectFound  bool
	}{
		{
			name:         "Existing film",
			id:           existingFilm.ID.String(),
			expectedCode: http.StatusOK,
			expectFound:  true,
		},
		{
			name:         "Invalid UUID",
			id:           "not-a-uuid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Not found film",
			id:           uuid.NewV4().String(),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/films/"+tt.id, nil)

			router := mux.NewRouter()
			router.HandleFunc("/api/films/{id}", handler.GetFilm)
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectFound {
				var film models.Film
				err := json.NewDecoder(w.Body).Decode(&film)
				require.NoError(t, err)
				assert.Equal(t, existingFilm.ID, film.ID)
				assert.Equal(t, existingFilm.Title, film.Title)
			}
		})
	}
}

func TestGetGenre(t *testing.T) {
	setupRepo()
	handler := NewFilmHandler()
	existingGenre := repo.Genres[0]

	tests := []struct {
		name         string
		id           string
		expectedCode int
		expectFound  bool
	}{
		{
			name:         "Existing genre",
			id:           existingGenre.ID.String(),
			expectedCode: http.StatusOK,
			expectFound:  true,
		},
		{
			name:         "Invalid UUID",
			id:           "bad-uuid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Not found genre",
			id:           uuid.NewV4().String(),
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/genres/"+tt.id, nil)

			router := mux.NewRouter()
			router.HandleFunc("/api/genres/{id}", handler.GetGenre)
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectFound {
				var genre models.Genre
				err := json.NewDecoder(w.Body).Decode(&genre)
				require.NoError(t, err)
				assert.Equal(t, existingGenre.ID, genre.ID)
				assert.Equal(t, existingGenre.Title, genre.Title)
			}
		})
	}
}

func TestGetGenres(t *testing.T) {
	setupRepo()
	handler := NewFilmHandler()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/genres", nil)

	handler.GetGenres(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	var genres []models.Genre
	err := json.NewDecoder(w.Body).Decode(&genres)
	require.NoError(t, err)
	assert.Len(t, genres, len(repo.Genres))
}

func TestGetFilmsByGenre(t *testing.T) {
	setupRepo()
	handler := NewFilmHandler()
	existingGenre := repo.Genres[0]

	tests := []struct {
		name         string
		id           string
		expectedCode int
		expectedLen  int
	}{
		{
			name:         "Films by valid genre",
			id:           existingGenre.ID.String(),
			expectedCode: http.StatusOK,
			expectedLen:  1,
		},
		{
			name:         "Invalid UUID",
			id:           "wrong-id",
			expectedCode: http.StatusBadRequest,
			expectedLen:  0,
		},
		{
			name:         "Genre with no films",
			id:           uuid.NewV4().String(),
			expectedCode: http.StatusOK,
			expectedLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/genres/"+tt.id+"/films", nil)

			router := mux.NewRouter()
			router.HandleFunc("/api/genres/{id}/films", handler.GetFilmsByGenre)
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var films []models.Film
				err := json.NewDecoder(w.Body).Decode(&films)
				require.NoError(t, err)
				assert.Len(t, films, tt.expectedLen)
			}
		})
	}
}
