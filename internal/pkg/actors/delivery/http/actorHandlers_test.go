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
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	"kinopoisk/internal/pkg/films/mocks"
	"kinopoisk/internal/pkg/middleware/logger"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	birthDateStr := birthDate.String()

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
		mockSetup      func(mockClient *mocks.MockFilmsClient)
		expectedStatus int
		expectBody     bool
		expectedActor  models.ActorPage
	}{
		{
			name:   "Success",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetActor(gomock.Any(), &gen.GetActorRequest{ActorId: actorIDStr}).
					Return(&gen.GetActorResponse{
						Actor: &gen.ActorPage{
							Id:            actorIDStr,
							RussianName:   "Леонардо ДиКаприо",
							OriginalName:  &originalName,
							Photo:         "/photos/leo.jpg",
							Height:        183,
							BirthDate:     birthDateStr,
							Age:           49,
							ZodiacSign:    "Скорпион",
							BirthPlace:    "Лос-Анджелес, США",
							MaritalStatus: "Не женат",
							FilmsNumber:   45,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedActor:  expectedActor,
		},
		{
			name:           "Invalid ID - empty string",
			varsID:         "",
			mockSetup:      func(mockClient *mocks.MockFilmsClient) {},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:           "Invalid ID - not a uuid",
			varsID:         "not-a-uuid",
			mockSetup:      func(mockClient *mocks.MockFilmsClient) {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Actor Not Found",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetActor(gomock.Any(), &gen.GetActorRequest{ActorId: actorIDStr}).
					Return(nil, status.Error(codes.NotFound, "actor not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:   "Internal Server Error",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetActor(gomock.Any(), &gen.GetActorRequest{ActorId: actorIDStr}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Unknown gRPC Error",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetActor(gomock.Any(), &gen.GetActorRequest{ActorId: actorIDStr}).
					Return(nil, errors.New("unknown error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Success with nil original name",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetActor(gomock.Any(), &gen.GetActorRequest{ActorId: actorIDStr}).
					Return(&gen.GetActorResponse{
						Actor: &gen.ActorPage{
							Id:            actorIDStr,
							RussianName:   "Леонардо ДиКаприо",
							OriginalName:  nil,
							Photo:         "/photos/leo.jpg",
							Height:        183,
							BirthDate:     birthDateStr,
							Age:           49,
							ZodiacSign:    "Скорпион",
							BirthPlace:    "Лос-Анджелес, США",
							MaritalStatus: "Не женат",
							FilmsNumber:   45,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedActor: func() models.ActorPage {
				actor := expectedActor
				actor.OriginalName = nil
				return actor
			}(),
		},
		{
			name:   "Success with empty photo",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetActor(gomock.Any(), &gen.GetActorRequest{ActorId: actorIDStr}).
					Return(&gen.GetActorResponse{
						Actor: &gen.ActorPage{
							Id:            actorIDStr,
							RussianName:   "Леонардо ДиКаприо",
							OriginalName:  &originalName,
							Photo:         "",
							Height:        183,
							BirthDate:     birthDateStr,
							Age:           49,
							ZodiacSign:    "Скорпион",
							BirthPlace:    "Лос-Анджелес, США",
							MaritalStatus: "Не женат",
							FilmsNumber:   45,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedActor: func() models.ActorPage {
				actor := expectedActor
				actor.Photo = ""
				return actor
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockFilmsClient(ctrl)
			handler := NewActorHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
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
				assert.Equal(t, tt.expectedActor.Height, decoded.Height)
				assert.Equal(t, tt.expectedActor.Age, decoded.Age)
				assert.Equal(t, tt.expectedActor.FilmsNumber, decoded.FilmsNumber)

				if tt.expectedActor.OriginalName == nil {
					assert.Nil(t, decoded.OriginalName)
				} else {
					assert.Equal(t, *tt.expectedActor.OriginalName, *decoded.OriginalName)
				}
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
		mockSetup      func(mockClient *mocks.MockFilmsClient)
		expectedStatus int
		expectBody     bool
		expectedFilms  []models.MainPageFilm
	}{
		{
			name:   "Success with films",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				var grpcFilms []*gen.MainPageFilm
				for _, film := range expectedFilms {
					grpcFilms = append(grpcFilms, &gen.MainPageFilm{
						Id:     film.ID.String(),
						Cover:  film.Cover,
						Title:  film.Title,
						Rating: film.Rating,
						Year:   int32(film.Year),
						Genre:  film.Genre,
					})
				}
				mockClient.EXPECT().
					GetFilmsByActor(gomock.Any(), &gen.GetFilmsByActorRequest{
						ActorId: actorIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetFilmsByActorResponse{Films: grpcFilms}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:   "Success with empty films list",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetFilmsByActor(gomock.Any(), &gen.GetFilmsByActorRequest{
						ActorId: actorIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetFilmsByActorResponse{Films: []*gen.MainPageFilm{}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  emptyFilms,
		},
		{
			name:           "Invalid ID - not a uuid",
			url:            "/actors/not-a-uuid/films",
			varsID:         "not-a-uuid",
			mockSetup:      func(mockClient *mocks.MockFilmsClient) {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Actor Not Found",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetFilmsByActor(gomock.Any(), &gen.GetFilmsByActorRequest{
						ActorId: actorIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(nil, status.Error(codes.NotFound, "actor not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:   "Internal Server Error",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetFilmsByActor(gomock.Any(), &gen.GetFilmsByActorRequest{
						ActorId: actorIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Unknown gRPC Error",
			url:    "/actors/" + actorIDStr + "/films?count=10&offset=0",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetFilmsByActor(gomock.Any(), &gen.GetFilmsByActorRequest{
						ActorId: actorIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(nil, errors.New("unknown error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Success with default pager values",
			url:    "/actors/" + actorIDStr + "/films",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				var grpcFilms []*gen.MainPageFilm
				for _, film := range expectedFilms {
					grpcFilms = append(grpcFilms, &gen.MainPageFilm{
						Id:     film.ID.String(),
						Cover:  film.Cover,
						Title:  film.Title,
						Rating: film.Rating,
						Year:   int32(film.Year),
						Genre:  film.Genre,
					})
				}
				mockClient.EXPECT().
					GetFilmsByActor(gomock.Any(), &gen.GetFilmsByActorRequest{
						ActorId: actorIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetFilmsByActorResponse{Films: grpcFilms}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:   "Success with custom pager values",
			url:    "/actors/" + actorIDStr + "/films?count=5&offset=10",
			varsID: actorIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				var grpcFilms []*gen.MainPageFilm
				for _, film := range expectedFilms {
					grpcFilms = append(grpcFilms, &gen.MainPageFilm{
						Id:     film.ID.String(),
						Cover:  film.Cover,
						Title:  film.Title,
						Rating: film.Rating,
						Year:   int32(film.Year),
						Genre:  film.Genre,
					})
				}
				mockClient.EXPECT().
					GetFilmsByActor(gomock.Any(), &gen.GetFilmsByActorRequest{
						ActorId: actorIDStr,
						Pager: &gen.Pager{
							Count:  5,
							Offset: 10,
						},
					}).
					Return(&gen.GetFilmsByActorResponse{Films: grpcFilms}, nil)
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

			mockClient := mocks.NewMockFilmsClient(ctrl)
			handler := NewActorHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
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
					assert.Equal(t, tt.expectedFilms[0].Year, decoded[0].Year)
					assert.Equal(t, tt.expectedFilms[0].Genre, decoded[0].Genre)
				}
			}
		})
	}
}
