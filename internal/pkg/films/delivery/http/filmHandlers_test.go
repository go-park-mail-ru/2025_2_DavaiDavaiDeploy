package filmHandlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/films/mocks"
	"kinopoisk/internal/pkg/middleware/logger"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func testContext() context.Context {
	testLogger := testLogger()
	return context.WithValue(context.Background(), logger.LoggerKey, testLogger)
}

func TestGetPromoFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockFilmUsecase(ctrl)
	handler := NewFilmHandler(mockUsecase)

	promoFilmID := uuid.NewV4()
	expectedPromoFilm := models.PromoFilm{
		ID:               promoFilmID,
		Image:            "/images/promo.jpg",
		Title:            "Промо фильм",
		Rating:           8.7,
		ShortDescription: "Краткое описание промо фильма",
		Year:             2024,
		Genre:            "Драма",
		Duration:         120,
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
			url:  "/films/promo",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetPromoFilm(gomock.Any()).
					Return(expectedPromoFilm, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name: "Usecase error",
			url:  "/films/promo",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetPromoFilm(gomock.Any()).
					Return(models.PromoFilm{}, errors.New("promo film not found"))
			},
			expectedStatus: http.StatusNotFound,
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
			router.HandleFunc("/films/promo", handler.GetPromoFilm)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.PromoFilm
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedPromoFilm, decoded)
			}
		})
	}
}

func TestGetFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockFilmUsecase(ctrl)
	handler := NewFilmHandler(mockUsecase)

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
			Genre:  "Комедия",
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
			url:  "/films?count=10&offset=0",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilms(gomock.Any(), models.Pager{Count: 10, Offset: 0}).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name: "Usecase error",
			url:  "/films?count=10&offset=0",
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilms(gomock.Any(), models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{}, errors.New("no films"))
			},
			expectedStatus: http.StatusBadRequest,
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
			router.HandleFunc("/films", handler.GetFilms)
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

func TestGetFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockFilmUsecase(ctrl)
	handler := NewFilmHandler(mockUsecase)

	filmID := uuid.NewV4()
	filmIDStr := filmID.String()
	originalTitle := "Original Film Title"
	slogan := "Great film slogan"
	image1 := "/images/image1.jpg"
	image2 := "/images/image2.jpg"
	image3 := "/images/image3.jpg"

	expectedFilm := models.FilmPage{
		ID:               filmID,
		Title:            "Название фильма",
		OriginalTitle:    &originalTitle,
		Cover:            "/covers/cover.jpg",
		Poster:           "/posters/poster.jpg",
		Genre:            "Драма",
		ShortDescription: "Краткое описание",
		Description:      "Полное описание фильма",
		AgeCategory:      "16+",
		Budget:           10000000,
		WorldwideFees:    50000000,
		TrailerURL:       nil,
		NumberOfRatings:  1500,
		Year:             2024,
		Rating:           8.5,
		Country:          "США",
		Slogan:           &slogan,
		Duration:         120,
		Image1:           &image1,
		Image2:           &image2,
		Image3:           &image3,
		Actors:           []models.Actor{},
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
			url:    "/films/" + filmIDStr,
			varsID: filmIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilm(gomock.Any(), filmID).
					Return(expectedFilm, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/films/not-a-uuid",
			varsID:         "not-a-uuid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Usecase error",
			url:    "/films/" + filmIDStr,
			varsID: filmIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilm(gomock.Any(), filmID).
					Return(models.FilmPage{}, errors.New("film not exists"))
			},
			expectedStatus: http.StatusBadRequest,
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
			router.HandleFunc("/films/{id}", handler.GetFilm)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.FilmPage
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedFilm, decoded)
			}
		})
	}
}

func TestGetFilmFeedbacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockFilmUsecase(ctrl)
	handler := NewFilmHandler(mockUsecase)

	filmID := uuid.NewV4()
	filmIDStr := filmID.String()
	title1 := "Отличный фильм!"
	text1 := "Потрясающая актерская игра и сюжет"
	title2 := "Хороший фильм"
	text2 := "Приятно посмотреть"

	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	expectedFeedbacks := []models.FilmFeedback{
		{
			ID:         uuid.NewV4(),
			UserID:     uuid.NewV4(),
			FilmID:     filmID,
			Title:      &title1,
			Text:       &text1,
			Rating:     9,
			CreatedAt:  fixedTime,
			UpdatedAt:  fixedTime,
			UserLogin:  "user1",
			UserAvatar: "/avatars/user1.jpg",
			IsMine:     false,
		},
		{
			ID:         uuid.NewV4(),
			UserID:     uuid.NewV4(),
			FilmID:     filmID,
			Title:      &title2,
			Text:       &text2,
			Rating:     8,
			CreatedAt:  fixedTime,
			UpdatedAt:  fixedTime,
			UserLogin:  "user2",
			UserAvatar: "/avatars/user2.jpg",
			IsMine:     false,
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
			url:    "/films/" + filmIDStr + "/feedbacks?count=10&offset=0",
			varsID: filmIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, models.Pager{Count: 10, Offset: 0}).
					Return(expectedFeedbacks, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/films/not-a-uuid/feedbacks",
			varsID:         "not-a-uuid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Usecase error",
			url:    "/films/" + filmIDStr + "/feedbacks?count=10&offset=0",
			varsID: filmIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, models.Pager{Count: 10, Offset: 0}).
					Return([]models.FilmFeedback{}, errors.New("no feedbacks"))
			},
			expectedStatus: http.StatusBadRequest,
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
			router.HandleFunc("/films/{id}/feedbacks", handler.GetFilmFeedbacks)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.FilmFeedback
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)

				assert.Equal(t, len(expectedFeedbacks), len(decoded))
				for i := range decoded {
					assert.Equal(t, expectedFeedbacks[i].ID, decoded[i].ID)
					assert.Equal(t, expectedFeedbacks[i].UserID, decoded[i].UserID)
					assert.Equal(t, expectedFeedbacks[i].FilmID, decoded[i].FilmID)
					assert.Equal(t, expectedFeedbacks[i].Title, decoded[i].Title)
					assert.Equal(t, expectedFeedbacks[i].Text, decoded[i].Text)
					assert.Equal(t, expectedFeedbacks[i].Rating, decoded[i].Rating)
					assert.Equal(t, expectedFeedbacks[i].UserLogin, decoded[i].UserLogin)
					assert.Equal(t, expectedFeedbacks[i].UserAvatar, decoded[i].UserAvatar)
					assert.Equal(t, expectedFeedbacks[i].IsMine, decoded[i].IsMine)
				}
			}
		})
	}
}
