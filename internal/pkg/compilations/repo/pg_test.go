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

func TestGetCompilationByID(t *testing.T) {
	compilationID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	tests := []struct {
		name            string
		compilationID   uuid.UUID
		repoMocker      func(*pgxpoolmock.MockPgxPool)
		wantCompilation models.Compilation
		wantErr         bool
	}{
		{
			name:          "Success",
			compilationID: compilationID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "title", "description", "icon", "created_at", "updated_at",
				}).
					AddRow(
						compilationID,
						"Best Dramas",
						"Collection of the best drama films",
						"/static/drama-compilation.png",
						createdAt,
						updatedAt,
					).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetCompilationByIDQuery, compilationID).
					Return(rows)
			},
			wantCompilation: models.Compilation{
				ID:          compilationID,
				Title:       "Best Dramas",
				Description: "Collection of the best drama films",
				Icon:        "/static/drama-compilation.png",
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

			repo := NewCompilationRepository(mockPool)
			compilation, err := repo.GetCompilationByID(testContext(), tt.compilationID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCompilation.ID, compilation.ID)
				assert.Equal(t, tt.wantCompilation.Title, compilation.Title)
				assert.Equal(t, tt.wantCompilation.Description, compilation.Description)
				assert.Equal(t, tt.wantCompilation.Icon, compilation.Icon)
			}
		})
	}
}

func TestGetCompilationsWithPagination(t *testing.T) {
	compilationID1 := uuid.NewV4()
	compilationID2 := uuid.NewV4()
	limit := 10
	offset := 0
	createdAt := time.Now()
	updatedAt := time.Now()

	tests := []struct {
		name             string
		limit            int
		offset           int
		repoMocker       func(*pgxpoolmock.MockPgxPool)
		wantCompilations []models.Compilation
		wantErr          bool
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
						compilationID1,
						"Action Collection",
						"Best action films",
						"/static/action-collection.png",
						createdAt,
						updatedAt,
					).
					AddRow(
						compilationID2,
						"Comedy Collection",
						"Funniest comedy films",
						"/static/comedy-collection.png",
						createdAt,
						updatedAt,
					).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetCompilationsWithPaginationQuery, limit, offset).
					Return(rows, nil)
			},
			wantCompilations: []models.Compilation{
				{
					ID:          compilationID1,
					Title:       "Action Collection",
					Description: "Best action films",
					Icon:        "/static/action-collection.png",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
				{
					ID:          compilationID2,
					Title:       "Comedy Collection",
					Description: "Funniest comedy films",
					Icon:        "/static/comedy-collection.png",
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
					Query(gomock.Any(), GetCompilationsWithPaginationQuery, limit, offset).
					Return(nil, assert.AnError)
			},
			wantCompilations: nil,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewCompilationRepository(mockPool)
			compilations, err := repo.GetCompilationsWithPagination(testContext(), tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, compilations)
			} else {
				assert.NoError(t, err)
				assert.Len(t, compilations, len(tt.wantCompilations))
				if len(compilations) > 0 {
					assert.Equal(t, tt.wantCompilations[0].ID, compilations[0].ID)
					assert.Equal(t, tt.wantCompilations[0].Title, compilations[0].Title)
					assert.Equal(t, tt.wantCompilations[1].ID, compilations[1].ID)
					assert.Equal(t, tt.wantCompilations[1].Title, compilations[1].Title)
				}
			}
		})
	}
}

func TestGetFilmsByCompilation(t *testing.T) {
	compilationID := uuid.NewV4()
	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()
	limit := 10
	offset := 0

	tests := []struct {
		name          string
		compilationID uuid.UUID
		limit         int
		offset        int
		repoMocker    func(*pgxpoolmock.MockPgxPool)
		wantFilms     []models.FavFilm
		wantErr       bool
	}{
		{
			name:          "Success",
			compilationID: compilationID,
			limit:         limit,
			offset:        offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mainRows := pgxpoolmock.NewRows([]string{
					"id", "title", "genre_title", "year", "duration", "cover", "short_description", "rating",
				}).
					AddRow(filmID1, "Film 1", "Drama", 2023, 120, "/static/cover1.jpg", "Short description 1", 8.5).
					AddRow(filmID2, "Film 2", "Comedy", 2022, 110, "/static/cover2.jpg", "Short description 2", 7.8).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByCompilationQuery, compilationID, limit, offset).
					Return(mainRows, nil)
			},
			wantFilms: []models.FavFilm{
				{
					ID:               filmID1,
					Title:            "Film 1",
					Genre:            "Drama",
					Year:             2023,
					Duration:         120,
					Image:            "/static/cover1.jpg",
					ShortDescription: "Short description 1",
					Rating:           8.5,
				},
				{
					ID:               filmID2,
					Title:            "Film 2",
					Genre:            "Comedy",
					Year:             2022,
					Duration:         110,
					Image:            "/static/cover2.jpg",
					ShortDescription: "Short description 2",
					Rating:           7.8,
				},
			},
			wantErr: false,
		},
		{
			name:          "QueryError",
			compilationID: compilationID,
			limit:         limit,
			offset:        offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByCompilationQuery, compilationID, limit, offset).
					Return(nil, assert.AnError)
			},
			wantFilms: nil,
			wantErr:   true,
		},
		{
			name:          "ScanError",
			compilationID: compilationID,
			limit:         limit,
			offset:        offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mainRows := pgxpoolmock.NewRows([]string{
					"id", "title", "genre_title", "year", "duration", "cover", "short_description", "rating",
				}).
					AddRow(filmID1, "", "", 0, 0, "", "", 0.0).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByCompilationQuery, compilationID, limit, offset).
					Return(mainRows, nil)
			},
			wantFilms: []models.FavFilm{
				{
					ID:               filmID1,
					Title:            "",
					Genre:            "",
					Year:             0,
					Duration:         0,
					Image:            "",
					ShortDescription: "",
					Rating:           0.0,
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

			repo := NewCompilationRepository(mockPool)
			films, err := repo.GetFilmsByCompilation(testContext(), tt.compilationID, tt.limit, tt.offset)

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
