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
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors"
	"kinopoisk/internal/pkg/actors/mocks"
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

func TestGetActor(t *testing.T) {
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
		varsID         string
		mockSetup      func(mockUsecase *mocks.MockActorUsecase)
		expectedStatus int
		expectBody     bool
		expectedActor  models.ActorPage
	}{
		{
			name:   "Success",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(expectedActor, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedActor:  expectedActor,
		},
		{
			name:           "Invalid ID - empty string",
			varsID:         "",
			mockSetup:      func(mockUsecase *mocks.MockActorUsecase) {},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:           "Invalid ID - not a uuid",
			varsID:         "not-a-uuid",
			mockSetup:      func(mockUsecase *mocks.MockActorUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Actor Not Found",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(models.ActorPage{}, actors.ErrorNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:   "Internal Server Error",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(models.ActorPage{}, actors.ErrorInternalServerError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Unknown Error",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(models.ActorPage{}, errors.New("unknown error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Success with nil original name",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				actorWithNilName := expectedActor
				actorWithNilName.OriginalName = nil
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(actorWithNilName, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedActor:  expectedActor,
		},
		{
			name:   "Success with empty photo",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				actorWithEmptyPhoto := expectedActor
				actorWithEmptyPhoto.Photo = ""
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(actorWithEmptyPhoto, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedActor:  expectedActor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockActorUsecase(ctrl)
			handler := NewActorHandler(mockUsecase)

			if tt.mockSetup != nil {
				tt.mockSetup(mockUsecase)
			}

			req := httptest.NewRequest(http.MethodGet, "/actors/"+tt.varsID, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/actors/{id}", handler.GetActor)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.ActorPage
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedActor.ID, decoded.ID)
				assert.Equal(t, tt.expectedActor.RussianName, decoded.RussianName)
			}
		})
	}
}

func TestGetFilmsByActor(t *testing.T) {
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

	emptyFilms := []models.MainPageFilm{}

	tests := []struct {
		name           string
		url            string
		varsID         string
		mockSetup      func(mockUsecase *mocks.MockActorUsecase)
		expectedStatus int
		expectBody     bool
		expectedFilms  []models.MainPageFilm
	}{
		{
			name:   "Success with films",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, models.Pager{Count: 10, Offset: 0}).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:   "Success with empty films list",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, models.Pager{Count: 10, Offset: 0}).
					Return(emptyFilms, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  emptyFilms,
		},
		{
			name:           "Invalid ID - not a uuid",
			url:            "/actors/not-a-uuid/films",
			varsID:         "not-a-uuid",
			mockSetup:      func(mockUsecase *mocks.MockActorUsecase) {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Actor Not Found",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{}, actors.ErrorNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:   "Internal Server Error",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{}, actors.ErrorInternalServerError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Unknown Error",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{}, errors.New("unknown error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Success with default pager values",
			url:    "/actors/" + actorIDStr + "/films",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, gomock.Any()).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:   "Success with custom pager values",
			url:    "/actors/" + actorIDStr + "/films?count=5&offset=10",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, models.Pager{Count: 5, Offset: 10}).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:   "Success with negative pager values",
			url:    "/actors/" + actorIDStr + "/films?count=-1&offset=-5",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, gomock.Any()).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:   "Success with large pager values",
			url:    "/actors/" + actorIDStr + "/films?count=1000&offset=500",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, models.Pager{Count: 1000, Offset: 500}).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockActorUsecase(ctrl)
			handler := NewActorHandler(mockUsecase)

			if tt.mockSetup != nil {
				tt.mockSetup(mockUsecase)
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/actors/{id}/films", handler.GetFilmsByActor)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.MainPageFilm
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedFilms), len(decoded))
				if len(decoded) > 0 {
					assert.Equal(t, tt.expectedFilms[0].Title, decoded[0].Title)
					assert.Equal(t, tt.expectedFilms[0].Rating, decoded[0].Rating)
				}
			}
		})
	}
}

func TestNewActorHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockActorUsecase(ctrl)

	t.Run("Success creation", func(t *testing.T) {
		handler := NewActorHandler(mockUsecase)
		assert.NotNil(t, handler)
		assert.Equal(t, mockUsecase, handler.uc)
	})

	t.Run("Creation with nil usecase", func(t *testing.T) {
		handler := NewActorHandler(nil)
		assert.NotNil(t, handler)
		assert.Nil(t, handler.uc)
	})
}

func TestGetActor_EdgeCases(t *testing.T) {
	actorID := uuid.NewV4()
	actorIDStr := actorID.String()

	tests := []struct {
		name           string
		varsID         string
		mockSetup      func(mockUsecase *mocks.MockActorUsecase)
		expectedStatus int
	}{
		{
			name:   "Context without logger",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(models.ActorPage{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Empty context",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetActor(gomock.Any(), actorID).
					Return(models.ActorPage{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockActorUsecase(ctrl)
			handler := NewActorHandler(mockUsecase)

			if tt.mockSetup != nil {
				tt.mockSetup(mockUsecase)
			}

			req := httptest.NewRequest(http.MethodGet, "/actors/"+tt.varsID, nil)
			if tt.name == "Empty context" {
				req = req.WithContext(context.Background())
			}
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/actors/{id}", handler.GetActor)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			} else {
				req = mux.SetURLVars(req, map[string]string{})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestGetFilmsByActor_EdgeCases(t *testing.T) {
	actorID := uuid.NewV4()
	actorIDStr := actorID.String()

	expectedFilms := []models.MainPageFilm{
		{
			ID:     uuid.NewV4(),
			Cover:  "/covers/test.jpg",
			Title:  "Test Film",
			Rating: 7.5,
			Year:   2020,
			Genre:  "Test Genre",
		},
	}

	tests := []struct {
		name           string
		url            string
		varsID         string
		mockSetup      func(mockUsecase *mocks.MockActorUsecase)
		expectedStatus int
	}{
		{
			name:   "Context without logger",
			url:    "/actors/" + actorIDStr + "/films",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, gomock.Any()).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Empty context",
			url:    "/actors/" + actorIDStr + "/films",
			varsID: actorIDStr,
			mockSetup: func(mockUsecase *mocks.MockActorUsecase) {
				mockUsecase.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, gomock.Any()).
					Return(expectedFilms, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocks.NewMockActorUsecase(ctrl)
			handler := NewActorHandler(mockUsecase)

			if tt.mockSetup != nil {
				tt.mockSetup(mockUsecase)
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			if tt.name == "Empty context" {
				req = req.WithContext(context.Background())
			}
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/actors/{id}/films", handler.GetFilmsByActor)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			} else {
				req = mux.SetURLVars(req, map[string]string{})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
