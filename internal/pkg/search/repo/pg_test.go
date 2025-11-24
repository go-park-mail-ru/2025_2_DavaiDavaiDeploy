package repo

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/middleware/logger"
	"log/slog"
	"os"
	"testing"

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

type MockRow struct {
	err error
}

func (m MockRow) Scan(dest ...interface{}) error {
	return m.err
}

func TestGetFilmsFromSearch(t *testing.T) {
	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()

	filmColumns := []string{"id", "cover", "title", "rating", "year", "genre"}

	tests := []struct {
		name         string
		searchString string
		limit        int
		offset       int
		repoMocker   func(*pgxpoolmock.MockPgxPool)
		wantErr      bool
		wantFilms    []models.MainPageFilm
		errorType    error
	}{
		{
			name:         "Success",
			searchString: "матрица",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows(filmColumns).
					AddRow(filmID1, "matrix.jpg", "Матрица", 8.7, 1999, "Фантастика").
					AddRow(filmID2, "reload.jpg", "Матрица: Перезагрузка", 7.2, 2003, "Фантастика").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsFromSearchQuery, "матрица", 10, 0).
					Return(filmRows, nil)
			},
			wantErr: false,
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID1,
					Cover:  "matrix.jpg",
					Title:  "Матрица",
					Rating: 8.7,
					Year:   1999,
					Genre:  "Фантастика",
				},
				{
					ID:     filmID2,
					Cover:  "reload.jpg",
					Title:  "Матрица: Перезагрузка",
					Rating: 7.2,
					Year:   2003,
					Genre:  "Фантастика",
				},
			},
		},
		{
			name:         "EmptyResult",
			searchString: "фильм",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows(filmColumns).ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsFromSearchQuery, "фильм", 10, 0).
					Return(filmRows, nil)
			},
			wantErr:   false,
			wantFilms: []models.MainPageFilm{},
		},
		{
			name:         "QueryError",
			searchString: "матрица",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsFromSearchQuery, "матрица", 10, 0).
					Return(nil, errors.New("database error"))
			},
			wantErr:   true,
			wantFilms: nil,
			errorType: films.ErrorInternalServerError,
		},
		{
			name:         "EnglishSearch",
			searchString: "matrix",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows(filmColumns).
					AddRow(filmID1, "matrix.jpg", "Матрица", 8.7, 1999, "Фантастика").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsFromSearchQuery, "matrix", 10, 0).
					Return(filmRows, nil)
			},
			wantErr: false,
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID1,
					Cover:  "matrix.jpg",
					Title:  "Матрица",
					Rating: 8.7,
					Year:   1999,
					Genre:  "Фантастика",
				},
			},
		},
		{
			name:         "RatingRounding",
			searchString: "фильм",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows(filmColumns).
					AddRow(filmID1, "film1.jpg", "Фильм 1", 8.756, 2020, "Драма").
					AddRow(filmID2, "film2.jpg", "Фильм 2", 7.234, 2021, "Комедия").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsFromSearchQuery, "фильм", 10, 0).
					Return(filmRows, nil)
			},
			wantErr: false,
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID1,
					Cover:  "film1.jpg",
					Title:  "Фильм 1",
					Rating: 8.8,
					Year:   2020,
					Genre:  "Драма",
				},
				{
					ID:     filmID2,
					Cover:  "film2.jpg",
					Title:  "Фильм 2",
					Rating: 7.2,
					Year:   2021,
					Genre:  "Комедия",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewSearchRepository(mockPool)
			films, err := repo.GetFilmsFromSearch(testContext(), tt.searchString, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.True(t, errors.Is(err, tt.errorType),
						"Expected error type: %v, got: %v", tt.errorType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, films, len(tt.wantFilms))
				if len(tt.wantFilms) > 0 {
					for i := range tt.wantFilms {
						assert.Equal(t, tt.wantFilms[i].ID, films[i].ID)
						assert.Equal(t, tt.wantFilms[i].Title, films[i].Title)
						assert.Equal(t, tt.wantFilms[i].Cover, films[i].Cover)
						assert.Equal(t, tt.wantFilms[i].Year, films[i].Year)
						assert.Equal(t, tt.wantFilms[i].Genre, films[i].Genre)
						assert.Equal(t, tt.wantFilms[i].Rating, films[i].Rating)
					}
				}
			}
		})
	}
}

func TestGetActorsFromSearch(t *testing.T) {
	actorID1 := uuid.NewV4()
	actorID2 := uuid.NewV4()

	actorColumns := []string{"id", "russian_name", "photo"}

	tests := []struct {
		name         string
		searchString string
		limit        int
		offset       int
		repoMocker   func(*pgxpoolmock.MockPgxPool)
		wantErr      bool
		wantActors   []models.MainPageActor
		errorType    error
	}{
		{
			name:         "Success",
			searchString: "том хэнкс",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				actorRows := pgxpoolmock.NewRows(actorColumns).
					AddRow(actorID1, "Том Хэнкс", "tom.jpg").
					AddRow(actorID2, "Том Круз", "tom_cruise.jpg").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetActorsFromSearchQuery, "том хэнкс", 10, 0).
					Return(actorRows, nil)
			},
			wantErr: false,
			wantActors: []models.MainPageActor{
				{
					ID:          actorID1,
					RussianName: "Том Хэнкс",
					Photo:       "tom.jpg",
				},
				{
					ID:          actorID2,
					RussianName: "Том Круз",
					Photo:       "tom_cruise.jpg",
				},
			},
		},
		{
			name:         "EmptyResult",
			searchString: "актер",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				actorRows := pgxpoolmock.NewRows(actorColumns).ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetActorsFromSearchQuery, "актер", 10, 0).
					Return(actorRows, nil)
			},
			wantErr:    false,
			wantActors: []models.MainPageActor{},
		},
		{
			name:         "QueryError",
			searchString: "том",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetActorsFromSearchQuery, "том", 10, 0).
					Return(nil, errors.New("database error"))
			},
			wantErr:    true,
			wantActors: nil,
			errorType:  actors.ErrorInternalServerError,
		},
		{
			name:         "EnglishSearch",
			searchString: "tom hanks",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				actorRows := pgxpoolmock.NewRows(actorColumns).
					AddRow(actorID1, "Том Хэнкс", "tom.jpg").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetActorsFromSearchQuery, "tom hanks", 10, 0).
					Return(actorRows, nil)
			},
			wantErr: false,
			wantActors: []models.MainPageActor{
				{
					ID:          actorID1,
					RussianName: "Том Хэнкс",
					Photo:       "tom.jpg",
				},
			},
		},
		{
			name:         "PartialName",
			searchString: "хэнкс",
			limit:        10,
			offset:       0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				actorRows := pgxpoolmock.NewRows(actorColumns).
					AddRow(actorID1, "Том Хэнкс", "tom.jpg").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetActorsFromSearchQuery, "хэнкс", 10, 0).
					Return(actorRows, nil)
			},
			wantErr: false,
			wantActors: []models.MainPageActor{
				{
					ID:          actorID1,
					RussianName: "Том Хэнкс",
					Photo:       "tom.jpg",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewSearchRepository(mockPool)
			actors, err := repo.GetActorsFromSearch(testContext(), tt.searchString, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.True(t, errors.Is(err, tt.errorType),
						"Expected error type: %v, got: %v", tt.errorType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, actors, len(tt.wantActors))
				if len(tt.wantActors) > 0 {
					for i := range tt.wantActors {
						assert.Equal(t, tt.wantActors[i].ID, actors[i].ID)
						assert.Equal(t, tt.wantActors[i].RussianName, actors[i].RussianName)
						assert.Equal(t, tt.wantActors[i].Photo, actors[i].Photo)
					}
				}
			}
		})
	}
}
