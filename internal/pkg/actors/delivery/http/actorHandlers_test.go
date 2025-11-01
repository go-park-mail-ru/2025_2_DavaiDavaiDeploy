package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors/mocks"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockActorUsecase(ctrl)
	handler := NewActorHandler(mockUsecase)

	actorID := uuid.NewV4()
	actorIDStr := actorID.String()
	originalName := "Leonardo DiCaprio"
	birthDate := time.Date(1974, 11, 11, 0, 0, 0, 0, time.UTC)

	expectedActor := models.ActorPage{
		ID:            actorID,
		RussianName:   "Леонардо ДиКаприо",
		OriginalName:  &originalName,
		Photo:         "/photos/leo.jpg",
		Height:        183,
		BirthDate:     birthDate,
		Age:           49,
		ZodiacSign:    "Скорпион",
		BirthPlace:    "Лос-Анджелес, США",
		MaritalStatus: "Не женат",
		FilmsNumber:   45,
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
			url:    "/actors/" + actorIDStr,
			varsID: actorIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(expectedActor, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/actors/not-a-uuid",
			varsID:         "not-a-uuid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Usecase error",
			url:    "/actors/" + actorIDStr,
			varsID: actorIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(models.ActorPage{}, errors.New("actor not exists"))
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

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/actors/{id}", handler.GetActor)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.ActorPage
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedActor, decoded)
			}
		})
	}
}

func TestGetFilmsByActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockActorUsecase(ctrl)
	handler := NewActorHandler(mockUsecase)

	actorID := uuid.NewV4()
	actorIDStr := actorID.String()

	expectedFilms := []models.MainPageFilm{
		{
			ID:     uuid.NewV4(),
			Cover:  "/covers/titanic.jpg",
			Title:  "Титаник",
			Rating: 8.5,
			Year:   1997,
			Genre:  "Драма",
		},
		{
			ID:     uuid.NewV4(),
			Cover:  "/covers/inception.jpg",
			Title:  "Начало",
			Rating: 8.8,
			Year:   2010,
			Genre:  "Фантастика",
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
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, models.Pager{Count: 10, Offset: 0}).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/actors/not-a-uuid/films",
			varsID:         "not-a-uuid",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Usecase error",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func() {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, models.Pager{Count: 10, Offset: 0}).
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

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/actors/{id}/films", handler.GetFilmsByActor)
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
