package repo

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
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

func TestGetGenreByID(t *testing.T) {
	genreID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	tests := []struct {
		name       string
		genreID    uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantGenre  models.Genre
		wantErr    bool
	}{
		{
			name:    "Success",
			genreID: genreID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "title", "description", "icon", "created_at", "updated_at",
				}).
					AddRow(
						genreID,
						"Drama",
						"Drama films description",
						"/static/drama.png",
						createdAt,
						updatedAt,
					).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetGenreByIDQuery, genreID).
					Return(rows)
			},
			wantGenre: models.Genre{
				ID:          genreID,
				Title:       "Drama",
				Description: "Drama films description",
				Icon:        "/static/drama.png",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewGenreRepository(mockPool)
			genre, err := repo.GetGenreByID(testContext(), tt.genreID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantGenre.ID, genre.ID)
				assert.Equal(t, tt.wantGenre.Title, genre.Title)
				assert.Equal(t, tt.wantGenre.Description, genre.Description)
				assert.Equal(t, tt.wantGenre.Icon, genre.Icon)
			}
		})
	}
}

func TestGetGenresWithPagination(t *testing.T) {
	genreID1 := uuid.NewV4()
	genreID2 := uuid.NewV4()
	limit := 10
	offset := 0
	createdAt := time.Now()
	updatedAt := time.Now()

	tests := []struct {
		name       string
		limit      int
		offset     int
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantGenres []models.Genre
		wantErr    bool
	}{
		{
			name:   "Success",
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "title", "description", "icon", "created_at", "updated_at",
				}).
					AddRow(
						genreID1,
						"Action",
						"Action films description",
						"/static/action.png",
						createdAt,
						updatedAt,
					).
					AddRow(
						genreID2,
						"Comedy",
						"Comedy films description",
						"/static/comedy.png",
						createdAt,
						updatedAt,
					).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetGenresWithPaginationQuery, limit, offset).
					Return(rows, nil)
			},
			wantGenres: []models.Genre{
				{
					ID:          genreID1,
					Title:       "Action",
					Description: "Action films description",
					Icon:        "/static/action.png",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
				{
					ID:          genreID2,
					Title:       "Comedy",
					Description: "Comedy films description",
					Icon:        "/static/comedy.png",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
			},
			wantErr: false,
		},
		{
			name:   "QueryError",
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetGenresWithPaginationQuery, limit, offset).
					Return(nil, assert.AnError)
			},
			wantGenres: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewGenreRepository(mockPool)
			genres, err := repo.GetGenresWithPagination(testContext(), tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, genres)
			} else {
				assert.NoError(t, err)
				assert.Len(t, genres, len(tt.wantGenres))
				if len(genres) > 0 {
					assert.Equal(t, tt.wantGenres[0].ID, genres[0].ID)
					assert.Equal(t, tt.wantGenres[0].Title, genres[0].Title)
					assert.Equal(t, tt.wantGenres[1].ID, genres[1].ID)
					assert.Equal(t, tt.wantGenres[1].Title, genres[1].Title)
				}
			}
		})
	}
}

func TestGetFilmAvgRating(t *testing.T) {
	filmID := uuid.NewV4()

	tests := []struct {
		name       string
		filmID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantRating float64
		wantErr    bool
	}{
		{
			name:   "Success",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(8.5).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID).
					Return(rows)
			},
			wantRating: 8.5,
			wantErr:    false,
		},
		{
			name:   "ZeroRating",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(0.0).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID).
					Return(rows)
			},
			wantRating: 0.0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewGenreRepository(mockPool)
			rating, err := repo.GetFilmAvgRating(testContext(), tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRating, rating)
			}
		})
	}
}

func TestGetFilmsByGenre(t *testing.T) {
	genreID := uuid.NewV4()
	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()
	limit := 10
	offset := 0

	tests := []struct {
		name       string
		genreID    uuid.UUID
		limit      int
		offset     int
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilms  []models.MainPageFilm
		wantErr    bool
	}{
		{
			name:    "Success",
			genreID: genreID,
			limit:   limit,
			offset:  offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mainRows := pgxpoolmock.NewRows([]string{
					"id", "cover", "title", "year", "genre_title",
				}).
					AddRow(filmID1, "/static/cover1.jpg", "Film 1", 2023, "Drama").
					AddRow(filmID2, "/static/cover2.jpg", "Film 2", 2022, "Drama").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByGenreQuery, genreID, limit, offset).
					Return(mainRows, nil)

				ratingRows1 := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(8.5).
					ToPgxRows()
				ratingRows1.Next()

				ratingRows2 := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(7.8).
					ToPgxRows()
				ratingRows2.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID1).
					Return(ratingRows1)
				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID2).
					Return(ratingRows2)
			},
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID1,
					Cover:  "/static/cover1.jpg",
					Title:  "Film 1",
					Year:   2023,
					Genre:  "Drama",
					Rating: 8.5,
				},
				{
					ID:     filmID2,
					Cover:  "/static/cover2.jpg",
					Title:  "Film 2",
					Year:   2022,
					Genre:  "Drama",
					Rating: 7.8,
				},
			},
			wantErr: false,
		},
		{
			name:    "QueryError",
			genreID: genreID,
			limit:   limit,
			offset:  offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByGenreQuery, genreID, limit, offset).
					Return(nil, assert.AnError)
			},
			wantFilms: nil,
			wantErr:   true,
		},
		{
			name:    "ScanError",
			genreID: genreID,
			limit:   limit,
			offset:  offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mainRows := pgxpoolmock.NewRows([]string{
					"id", "cover", "title", "year", "genre_title",
				}).
					AddRow(filmID1, "", "", 0, "").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByGenreQuery, genreID, limit, offset).
					Return(mainRows, nil)

				ratingRows := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(0.0).
					ToPgxRows()
				ratingRows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID1).
					Return(ratingRows)
			},
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID1,
					Cover:  "",
					Title:  "",
					Year:   0,
					Genre:  "",
					Rating: 0.0,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewGenreRepository(mockPool)
			films, err := repo.GetFilmsByGenre(testContext(), tt.genreID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, films)
			} else {
				assert.NoError(t, err)
				assert.Len(t, films, len(tt.wantFilms))
				if len(films) > 0 {
					assert.Equal(t, tt.wantFilms[0].ID, films[0].ID)
					assert.Equal(t, tt.wantFilms[0].Title, films[0].Title)
					assert.Equal(t, tt.wantFilms[0].Rating, films[0].Rating)
				}
			}
		})
	}
}
