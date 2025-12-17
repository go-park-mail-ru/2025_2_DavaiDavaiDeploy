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
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	"kinopoisk/internal/pkg/films/mocks"
	"kinopoisk/internal/pkg/middleware/logger"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

func testContextWithUser(userID uuid.UUID) context.Context {
	ctx := testContext()
	user := models.User{
		ID: userID,
	}
	return context.WithValue(ctx, auth.UserKey, user)
}

func TestGetCompilation(t *testing.T) {
	compilationID := uuid.NewV4()
	compilationIDStr := compilationID.String()

	expectedCompilation := models.Compilation{
		ID:          compilationID,
		Title:       "Лучшие фильмы 2024",
		Description: "Подборка самых популярных фильмов 2024 года",
		Icon:        "/icons/best2024.png",
	}

	tests := []struct {
		name                string
		varsID              string
		mockSetup           func(mockClient *mocks.MockFilmsClient)
		expectedStatus      int
		expectBody          bool
		expectedCompilation models.Compilation
	}{
		{
			name:   "Success",
			varsID: compilationIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetCompilation(gomock.Any(), &gen.GetCompilationRequest{CompilationId: compilationIDStr}).
					Return(&gen.GetCompilationResponse{
						Compilation: &gen.Compilation{
							Id:          compilationIDStr,
							Name:        "Лучшие фильмы 2024",
							Description: "Подборка самых популярных фильмов 2024 года",
							Icon:        "/icons/best2024.png",
						},
					}, nil)
			},
			expectedStatus:      http.StatusOK,
			expectBody:          true,
			expectedCompilation: expectedCompilation,
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
			name:   "Compilation Not Found",
			varsID: compilationIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetCompilation(gomock.Any(), &gen.GetCompilationRequest{CompilationId: compilationIDStr}).
					Return(nil, status.Error(codes.NotFound, "compilation not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:   "Internal Server Error",
			varsID: compilationIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetCompilation(gomock.Any(), &gen.GetCompilationRequest{CompilationId: compilationIDStr}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Unknown gRPC Error",
			varsID: compilationIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetCompilation(gomock.Any(), &gen.GetCompilationRequest{CompilationId: compilationIDStr}).
					Return(nil, errors.New("unknown error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockFilmsClient(ctrl)
			handler := NewCompilationHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			req := httptest.NewRequest(http.MethodGet, "/compilations/"+tt.varsID, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/compilations/{id}", handler.GetCompilation)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.Compilation
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCompilation, decoded)
			}
		})
	}
}

func TestGetCompilations(t *testing.T) {
	expectedCompilations := []models.Compilation{
		{
			ID:          uuid.NewV4(),
			Title:       "Лучшие фильмы 2024",
			Description: "Подборка самых популярных фильмов 2024 года",
			Icon:        "/icons/best2024.png",
		},
		{
			ID:          uuid.NewV4(),
			Title:       "Классика кино",
			Description: "Вечная классика мирового кинематографа",
			Icon:        "/icons/classic.png",
		},
		{
			ID:          uuid.NewV4(),
			Title:       "Новинки",
			Description: "Самые свежие премьеры",
			Icon:        "/icons/new.png",
		},
	}

	tests := []struct {
		name                 string
		url                  string
		mockSetup            func(mockClient *mocks.MockFilmsClient)
		expectedStatus       int
		expectBody           bool
		expectedCompilations []models.Compilation
	}{
		{
			name: "Success with compilations",
			url:  "/compilations?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				var grpcCompilations []*gen.Compilation
				for _, compilation := range expectedCompilations {
					grpcCompilations = append(grpcCompilations, &gen.Compilation{
						Id:          compilation.ID.String(),
						Name:        compilation.Title,
						Description: compilation.Description,
						Icon:        compilation.Icon,
					})
				}
				mockClient.EXPECT().
					GetCompilations(gomock.Any(), &gen.GetCompilationsRequest{
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetCompilationsResponse{Compilations: grpcCompilations}, nil)
			},
			expectedStatus:       http.StatusOK,
			expectBody:           true,
			expectedCompilations: expectedCompilations,
		},
		{
			name: "Success with empty compilations list",
			url:  "/compilations?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetCompilations(gomock.Any(), &gen.GetCompilationsRequest{
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetCompilationsResponse{Compilations: []*gen.Compilation{}}, nil)
			},
			expectedStatus:       http.StatusOK,
			expectBody:           true,
			expectedCompilations: []models.Compilation{},
		},
		{
			name: "Compilations Not Found",
			url:  "/compilations?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetCompilations(gomock.Any(), &gen.GetCompilationsRequest{
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(nil, status.Error(codes.NotFound, "compilations not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name: "Internal Server Error",
			url:  "/compilations?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetCompilations(gomock.Any(), &gen.GetCompilationsRequest{
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
			name: "Unknown gRPC Error",
			url:  "/compilations?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetCompilations(gomock.Any(), &gen.GetCompilationsRequest{
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
			name: "Success with default pager values",
			url:  "/compilations",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				var grpcCompilations []*gen.Compilation
				for _, compilation := range expectedCompilations {
					grpcCompilations = append(grpcCompilations, &gen.Compilation{
						Id:          compilation.ID.String(),
						Name:        compilation.Title,
						Description: compilation.Description,
						Icon:        compilation.Icon,
					})
				}
				mockClient.EXPECT().
					GetCompilations(gomock.Any(), &gen.GetCompilationsRequest{
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetCompilationsResponse{Compilations: grpcCompilations}, nil)
			},
			expectedStatus:       http.StatusOK,
			expectBody:           true,
			expectedCompilations: expectedCompilations,
		},
		{
			name: "Success with custom pager values",
			url:  "/compilations?count=5&offset=10",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				var grpcCompilations []*gen.Compilation
				for _, compilation := range expectedCompilations {
					grpcCompilations = append(grpcCompilations, &gen.Compilation{
						Id:          compilation.ID.String(),
						Name:        compilation.Title,
						Description: compilation.Description,
						Icon:        compilation.Icon,
					})
				}
				mockClient.EXPECT().
					GetCompilations(gomock.Any(), &gen.GetCompilationsRequest{
						Pager: &gen.Pager{
							Count:  5,
							Offset: 10,
						},
					}).
					Return(&gen.GetCompilationsResponse{Compilations: grpcCompilations}, nil)
			},
			expectedStatus:       http.StatusOK,
			expectBody:           true,
			expectedCompilations: expectedCompilations,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockFilmsClient(ctrl)
			handler := NewCompilationHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/compilations", handler.GetCompilations)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.Compilation
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedCompilations), len(decoded))
				if len(decoded) > 0 {
					assert.Equal(t, tt.expectedCompilations[0].Title, decoded[0].Title)
					assert.Equal(t, tt.expectedCompilations[0].Description, decoded[0].Description)
					assert.Equal(t, tt.expectedCompilations[0].Icon, decoded[0].Icon)
				}
			}
		})
	}
}

func TestGetFilmsByCompilation(t *testing.T) {
	compilationID := uuid.NewV4()
	compilationIDStr := compilationID.String()
	userID := uuid.NewV4()

	expectedFilms := []models.CompFilm{
		{
			ID:               uuid.NewV4(),
			Image:            "/covers/film1.jpg",
			Title:            "Фильм 1",
			Rating:           8.5,
			Year:             2024,
			Genre:            "Драма",
			ShortDescription: "Описание фильма 1",
			Duration:         120,
			IsLiked:          true,
		},
		{
			ID:               uuid.NewV4(),
			Image:            "/covers/film2.jpg",
			Title:            "Фильм 2",
			Rating:           7.9,
			Year:             2023,
			Genre:            "Комедия",
			ShortDescription: "Описание фильма 2",
			Duration:         110,
			IsLiked:          false,
		},
	}

	emptyFilms := []models.CompFilm{}

	tests := []struct {
		name           string
		url            string
		varsID         string
		context        context.Context
		mockSetup      func(mockClient *mocks.MockFilmsClient, userID string)
		expectedStatus int
		expectBody     bool
		expectedFilms  []models.CompFilm
	}{
		{
			name:    "Success with films",
			url:     "/compilations/" + compilationIDStr + "/films?count=10&offset=0",
			varsID:  compilationIDStr,
			context: testContextWithUser(userID),
			mockSetup: func(mockClient *mocks.MockFilmsClient, userID string) {
				var grpcFilms []*gen.CompFilm
				for _, film := range expectedFilms {
					grpcFilms = append(grpcFilms, &gen.CompFilm{
						Id:               film.ID.String(),
						Image:            film.Image,
						Title:            film.Title,
						Rating:           film.Rating,
						Year:             int32(film.Year),
						Genre:            film.Genre,
						ShortDescription: film.ShortDescription,
						Duration:         int32(film.Duration),
						IsLiked:          film.IsLiked,
					})
				}
				mockClient.EXPECT().
					GetFilmsByCompilation(gomock.Any(), &gen.GetFilmsByCompilationRequest{
						CompilationId: compilationIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
						UserId: userID,
					}).
					Return(&gen.GetFilmsByCompilationResponse{Films: grpcFilms}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:    "Success with empty films list",
			url:     "/compilations/" + compilationIDStr + "/films?count=10&offset=0",
			varsID:  compilationIDStr,
			context: testContextWithUser(userID),
			mockSetup: func(mockClient *mocks.MockFilmsClient, userID string) {
				mockClient.EXPECT().
					GetFilmsByCompilation(gomock.Any(), &gen.GetFilmsByCompilationRequest{
						CompilationId: compilationIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
						UserId: userID,
					}).
					Return(&gen.GetFilmsByCompilationResponse{Films: []*gen.CompFilm{}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  emptyFilms,
		},
		{
			name:           "Invalid ID - not a uuid",
			url:            "/compilations/not-a-uuid/films",
			varsID:         "not-a-uuid",
			context:        testContextWithUser(userID),
			mockSetup:      func(mockClient *mocks.MockFilmsClient, userID string) {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:    "Compilation Not Found",
			url:     "/compilations/" + compilationIDStr + "/films?count=10&offset=0",
			varsID:  compilationIDStr,
			context: testContextWithUser(userID),
			mockSetup: func(mockClient *mocks.MockFilmsClient, userID string) {
				mockClient.EXPECT().
					GetFilmsByCompilation(gomock.Any(), &gen.GetFilmsByCompilationRequest{
						CompilationId: compilationIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
						UserId: userID,
					}).
					Return(nil, status.Error(codes.NotFound, "compilation not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:    "Internal Server Error",
			url:     "/compilations/" + compilationIDStr + "/films?count=10&offset=0",
			varsID:  compilationIDStr,
			context: testContextWithUser(userID),
			mockSetup: func(mockClient *mocks.MockFilmsClient, userID string) {
				mockClient.EXPECT().
					GetFilmsByCompilation(gomock.Any(), &gen.GetFilmsByCompilationRequest{
						CompilationId: compilationIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
						UserId: userID,
					}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:    "Unknown gRPC Error",
			url:     "/compilations/" + compilationIDStr + "/films?count=10&offset=0",
			varsID:  compilationIDStr,
			context: testContextWithUser(userID),
			mockSetup: func(mockClient *mocks.MockFilmsClient, userID string) {
				mockClient.EXPECT().
					GetFilmsByCompilation(gomock.Any(), &gen.GetFilmsByCompilationRequest{
						CompilationId: compilationIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
						UserId: userID,
					}).
					Return(nil, errors.New("unknown error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:    "Success with default pager values",
			url:     "/compilations/" + compilationIDStr + "/films",
			varsID:  compilationIDStr,
			context: testContextWithUser(userID),
			mockSetup: func(mockClient *mocks.MockFilmsClient, userID string) {
				var grpcFilms []*gen.CompFilm
				for _, film := range expectedFilms {
					grpcFilms = append(grpcFilms, &gen.CompFilm{
						Id:               film.ID.String(),
						Image:            film.Image,
						Title:            film.Title,
						Rating:           film.Rating,
						Year:             int32(film.Year),
						Genre:            film.Genre,
						ShortDescription: film.ShortDescription,
						Duration:         int32(film.Duration),
						IsLiked:          film.IsLiked,
					})
				}
				mockClient.EXPECT().
					GetFilmsByCompilation(gomock.Any(), &gen.GetFilmsByCompilationRequest{
						CompilationId: compilationIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
						UserId: userID,
					}).
					Return(&gen.GetFilmsByCompilationResponse{Films: grpcFilms}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:    "Success with custom pager values",
			url:     "/compilations/" + compilationIDStr + "/films?count=5&offset=10",
			varsID:  compilationIDStr,
			context: testContextWithUser(userID),
			mockSetup: func(mockClient *mocks.MockFilmsClient, userID string) {
				var grpcFilms []*gen.CompFilm
				for _, film := range expectedFilms {
					grpcFilms = append(grpcFilms, &gen.CompFilm{
						Id:               film.ID.String(),
						Image:            film.Image,
						Title:            film.Title,
						Rating:           film.Rating,
						Year:             int32(film.Year),
						Genre:            film.Genre,
						ShortDescription: film.ShortDescription,
						Duration:         int32(film.Duration),
						IsLiked:          film.IsLiked,
					})
				}
				mockClient.EXPECT().
					GetFilmsByCompilation(gomock.Any(), &gen.GetFilmsByCompilationRequest{
						CompilationId: compilationIDStr,
						Pager: &gen.Pager{
							Count:  5,
							Offset: 10,
						},
						UserId: userID,
					}).
					Return(&gen.GetFilmsByCompilationResponse{Films: grpcFilms}, nil)
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
			handler := NewCompilationHandler(mockClient)

			var userIDStr string
			if user, ok := tt.context.Value(auth.UserKey).(models.User); ok {
				userIDStr = user.ID.String()
			}

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient, userIDStr)
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(tt.context)
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/compilations/{id}/films", handler.GetFilmsByCompilation)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.CompFilm
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedFilms), len(decoded))
				if len(decoded) > 0 {
					assert.Equal(t, tt.expectedFilms[0].Title, decoded[0].Title)
					assert.Equal(t, tt.expectedFilms[0].Rating, decoded[0].Rating)
					assert.Equal(t, tt.expectedFilms[0].Year, decoded[0].Year)
					assert.Equal(t, tt.expectedFilms[0].Genre, decoded[0].Genre)
					assert.Equal(t, tt.expectedFilms[0].ShortDescription, decoded[0].ShortDescription)
					assert.Equal(t, tt.expectedFilms[0].Duration, decoded[0].Duration)
					assert.Equal(t, tt.expectedFilms[0].IsLiked, decoded[0].IsLiked)
				}
			}
		})
	}
}
