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

func TestGetGenre(t *testing.T) {
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
		varsID         string
		mockSetup      func(mockClient *mocks.MockFilmsClient)
		expectedStatus int
		expectBody     bool
		expectedGenre  models.Genre
	}{
		{
			name:   "Success",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetGenre(gomock.Any(), &gen.GetGenreRequest{GenreId: genreIDStr}).
					Return(&gen.GetGenreResponse{
						Genre: &gen.Genre{
							Id:          genreIDStr,
							Name:        "Драма",
							Description: "Драматические фильмы",
							Icon:        "/icons/drama.png",
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedGenre:  expectedGenre,
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
			name:   "Genre Not Found",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetGenre(gomock.Any(), &gen.GetGenreRequest{GenreId: genreIDStr}).
					Return(nil, status.Error(codes.NotFound, "genre not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:   "Invalid Argument",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetGenre(gomock.Any(), &gen.GetGenreRequest{GenreId: genreIDStr}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid argument"))
			},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Internal Server Error",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetGenre(gomock.Any(), &gen.GetGenreRequest{GenreId: genreIDStr}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:   "Unknown gRPC Error",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetGenre(gomock.Any(), &gen.GetGenreRequest{GenreId: genreIDStr}).
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
			handler := NewGenreHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			req := httptest.NewRequest(http.MethodGet, "/genres/"+tt.varsID, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/genres/{id}", handler.GetGenre)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.Genre
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedGenre, decoded)
			}
		})
	}
}

func TestGetGenres(t *testing.T) {
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
		mockSetup      func(mockClient *mocks.MockFilmsClient)
		expectedStatus int
		expectBody     bool
		expectedGenres []models.Genre
	}{
		{
			name: "Success with genres",
			url:  "/genres?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				var grpcGenres []*gen.Genre
				for _, genre := range expectedGenres {
					grpcGenres = append(grpcGenres, &gen.Genre{
						Id:          genre.ID.String(),
						Name:        genre.Title,
						Description: genre.Description,
						Icon:        genre.Icon,
					})
				}
				mockClient.EXPECT().
					GetGenres(gomock.Any(), &gen.GetGenresRequest{
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetGenresResponse{Genres: grpcGenres}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedGenres: expectedGenres,
		},
		{
			name: "Success with empty genres list",
			url:  "/genres?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetGenres(gomock.Any(), &gen.GetGenresRequest{
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetGenresResponse{Genres: []*gen.Genre{}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedGenres: []models.Genre{},
		},
		{
			name: "Genres Not Found",
			url:  "/genres?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetGenres(gomock.Any(), &gen.GetGenresRequest{
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(nil, status.Error(codes.NotFound, "genres not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name: "Internal Server Error",
			url:  "/genres?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetGenres(gomock.Any(), &gen.GetGenresRequest{
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
			url:  "/genres?count=10&offset=0",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetGenres(gomock.Any(), &gen.GetGenresRequest{
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
			url:  "/genres",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				var grpcGenres []*gen.Genre
				for _, genre := range expectedGenres {
					grpcGenres = append(grpcGenres, &gen.Genre{
						Id:          genre.ID.String(),
						Name:        genre.Title,
						Description: genre.Description,
						Icon:        genre.Icon,
					})
				}
				mockClient.EXPECT().
					GetGenres(gomock.Any(), &gen.GetGenresRequest{
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetGenresResponse{Genres: grpcGenres}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedGenres: expectedGenres,
		},
		{
			name: "Success with custom pager values",
			url:  "/genres?count=5&offset=10",
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				var grpcGenres []*gen.Genre
				for _, genre := range expectedGenres {
					grpcGenres = append(grpcGenres, &gen.Genre{
						Id:          genre.ID.String(),
						Name:        genre.Title,
						Description: genre.Description,
						Icon:        genre.Icon,
					})
				}
				mockClient.EXPECT().
					GetGenres(gomock.Any(), &gen.GetGenresRequest{
						Pager: &gen.Pager{
							Count:  5,
							Offset: 10,
						},
					}).
					Return(&gen.GetGenresResponse{Genres: grpcGenres}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedGenres: expectedGenres,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockFilmsClient(ctrl)
			handler := NewGenreHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
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
				assert.Equal(t, len(tt.expectedGenres), len(decoded))
				if len(decoded) > 0 {
					assert.Equal(t, tt.expectedGenres[0].Title, decoded[0].Title)
					assert.Equal(t, tt.expectedGenres[0].Description, decoded[0].Description)
					assert.Equal(t, tt.expectedGenres[0].Icon, decoded[0].Icon)
				}
			}
		})
	}
}

func TestGetFilmsByGenre(t *testing.T) {
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
			url:    "/genres/" + genreIDStr + "/films?count=10&offset=0",
			varsID: genreIDStr,
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
					GetFilmsByGenre(gomock.Any(), &gen.GetFilmsByGenreRequest{
						GenreId: genreIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetFilmsByGenreResponse{Films: grpcFilms}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:   "Success with empty films list",
			url:    "/genres/" + genreIDStr + "/films?count=10&offset=0",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetFilmsByGenre(gomock.Any(), &gen.GetFilmsByGenreRequest{
						GenreId: genreIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetFilmsByGenreResponse{Films: []*gen.MainPageFilm{}}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  emptyFilms,
		},
		{
			name:           "Invalid ID - not a uuid",
			url:            "/genres/not-a-uuid/films",
			varsID:         "not-a-uuid",
			mockSetup:      func(mockClient *mocks.MockFilmsClient) {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Genre Not Found",
			url:    "/genres/" + genreIDStr + "/films?count=10&offset=0",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetFilmsByGenre(gomock.Any(), &gen.GetFilmsByGenreRequest{
						GenreId: genreIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(nil, status.Error(codes.NotFound, "genre not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name:   "Invalid Argument",
			url:    "/genres/" + genreIDStr + "/films?count=10&offset=0",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetFilmsByGenre(gomock.Any(), &gen.GetFilmsByGenreRequest{
						GenreId: genreIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid argument"))
			},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:   "Internal Server Error",
			url:    "/genres/" + genreIDStr + "/films?count=10&offset=0",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetFilmsByGenre(gomock.Any(), &gen.GetFilmsByGenreRequest{
						GenreId: genreIDStr,
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
			url:    "/genres/" + genreIDStr + "/films?count=10&offset=0",
			varsID: genreIDStr,
			mockSetup: func(mockClient *mocks.MockFilmsClient) {
				mockClient.EXPECT().
					GetFilmsByGenre(gomock.Any(), &gen.GetFilmsByGenreRequest{
						GenreId: genreIDStr,
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
			url:    "/genres/" + genreIDStr + "/films",
			varsID: genreIDStr,
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
					GetFilmsByGenre(gomock.Any(), &gen.GetFilmsByGenreRequest{
						GenreId: genreIDStr,
						Pager: &gen.Pager{
							Count:  10,
							Offset: 0,
						},
					}).
					Return(&gen.GetFilmsByGenreResponse{Films: grpcFilms}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedFilms:  expectedFilms,
		},
		{
			name:   "Success with custom pager values",
			url:    "/genres/" + genreIDStr + "/films?count=5&offset=10",
			varsID: genreIDStr,
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
					GetFilmsByGenre(gomock.Any(), &gen.GetFilmsByGenreRequest{
						GenreId: genreIDStr,
						Pager: &gen.Pager{
							Count:  5,
							Offset: 10,
						},
					}).
					Return(&gen.GetFilmsByGenreResponse{Films: grpcFilms}, nil)
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
			handler := NewGenreHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/genres/{id}/films", handler.GetFilmsByGenre)

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
