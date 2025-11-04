package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/genres"
	"kinopoisk/internal/pkg/genres/mocks"
	"kinopoisk/internal/pkg/middleware/logger"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func testContext() context.Context {
	testLogger := testLogger()
	return context.WithValue(context.Background(), logger.LoggerKey, testLogger)
}

func TestGetGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockGenreUsecase(ctrl)
	handler := NewGenreHandler(mockUsecase)

	genreID := uuid.NewV4()
	genreIDStr := genreID.String()

	expectedGenre := models.Genre{
		ID:          genreID,
		Title:       "Драма",
		Description: "Драматические фильмы",
		Icon:        "/icons/drama.png",
	}

	tests := []struct {
		name           string
		url            string
		varsID         string
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name:   "Success",
			url:    "/genres/" + genreIDStr,
			varsID: genreIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetGenre(gomock.Any(), genreID).
					Return(expectedGenre, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/genres/not-a-uuid",
			varsID:         "not-a-uuid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Usecase not found error",
			url:    "/genres/" + genreIDStr,
			varsID: genreIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetGenre(gomock.Any(), genreID).
					Return(models.Genre{}, genres.ErrorNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:   "Usecase internal error",
			url:    "/genres/" + genreIDStr,
			varsID: genreIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetGenre(gomock.Any(), genreID).
					Return(models.Genre{}, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/genres/{id}", handler.GetGenre)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.Genre
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedGenre, decoded)
			}
		})
	}
}

func TestGetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockGenreUsecase(ctrl)
	handler := NewGenreHandler(mockUsecase)

	expectedGenres := []models.Genre{
		{
			ID:          uuid.NewV4(),
			Title:       "Драма",
			Description: "Драматические фильмы",
			Icon:        "/icons/drama.png",
		},
		{
			ID:          uuid.NewV4(),
			Title:       "Комедия",
			Description: "Комедийные фильмы",
			Icon:        "/icons/comedy.png",
		},
		{
			ID:          uuid.NewV4(),
			Title:       "Боевик",
			Description: "Экшн-фильмы",
			Icon:        "/icons/action.png",
		},
	}

	tests := []struct {
		name           string
		url            string
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name: "Success",
			url:  "/genres?count=10&offset=0",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetGenres(gomock.Any(), models.Pager{Count: 10, Offset: 0}).
					Return(expectedGenres, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name: "Usecase not found error",
			url:  "/genres?count=10&offset=0",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetGenres(gomock.Any(), models.Pager{Count: 10, Offset: 0}).
					Return([]models.Genre{}, genres.ErrorNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name: "Usecase internal error",
			url:  "/genres?count=10&offset=0",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetGenres(gomock.Any(), models.Pager{Count: 10, Offset: 0}).
					Return([]models.Genre{}, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/genres", handler.GetGenres)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.Genre
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedGenres, decoded)
			}
		})
	}
}

func TestGetFilmsByGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockGenreUsecase(ctrl)
	handler := NewGenreHandler(mockUsecase)

	genreID := uuid.NewV4()
	genreIDStr := genreID.String()

	expectedFilms := []models.MainPageFilm{
		{
			ID:     uuid.NewV4(),
			Cover:  "/covers/film1.jpg",
			Title:  "Фильм 1",
			Rating: 8.5,
			Year:   2024,
			Genre:  "Драма",
		},
		{
			ID:     uuid.NewV4(),
			Cover:  "/covers/film2.jpg",
			Title:  "Фильм 2",
			Rating: 7.9,
			Year:   2023,
			Genre:  "Драма",
		},
	}

	tests := []struct {
		name           string
		url            string
		varsID         string
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name:   "Success",
			url:    "/genres/" + genreIDStr + "/films?count=10&offset=0",
			varsID: genreIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilmsByGenre(gomock.Any(), genreID, models.Pager{Count: 10, Offset: 0}).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/genres/not-a-uuid/films",
			varsID:         "not-a-uuid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Usecase not found error",
			url:    "/genres/" + genreIDStr + "/films?count=10&offset=0",
			varsID: genreIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilmsByGenre(gomock.Any(), genreID, models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{}, genres.ErrorNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:   "Usecase internal error",
			url:    "/genres/" + genreIDStr + "/films?count=10&offset=0",
			varsID: genreIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilmsByGenre(gomock.Any(), genreID, models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{}, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/genres/{id}/films", handler.GetFilmsByGenre)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.MainPageFilm
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedFilms, decoded)
			}
		})
	}
}
