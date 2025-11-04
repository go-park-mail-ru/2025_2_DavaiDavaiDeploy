package repo

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors"
	"kinopoisk/internal/pkg/middleware/logger"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
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

func TestGetActorByID(t *testing.T) {
	actorID := uuid.NewV4()
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	originalName := "Tom Hanks"

	columns := []string{
		"id", "russian_name", "original_name", "photo", "height",
		"birth_date", "death_date", "zodiac_sign", "birth_place", "marital_status",
	}

	tests := []struct {
		name       string
		actorID    uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
		wantActor  models.Actor
	}{
		{
			name:    "Success",
			actorID: actorID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows(columns).
					AddRow(
						actorID,
						"Том Хэнкс",
						&originalName,
						"tom.jpg",
						183,
						birthDate,
						nil,
						"Козерог",
						"Конкорд, Калифорния, США",
						"Женат",
					).ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetActorByID, actorID).
					Return(rows)
			},
			wantErr: false,
			wantActor: models.Actor{
				ID:            actorID,
				RussianName:   "Том Хэнкс",
				OriginalName:  &originalName,
				Photo:         "tom.jpg",
				Height:        183,
				BirthDate:     birthDate,
				DeathDate:     nil,
				ZodiacSign:    "Козерог",
				BirthPlace:    "Конкорд, Калифорния, США",
				MaritalStatus: "Женат",
			},
		},
		{
			name:    "NotFound",
			actorID: actorID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetActorByID, actorID).
					Return(MockRow{err: pgx.ErrNoRows})
			},
			wantErr:   true,
			wantActor: models.Actor{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewActorRepository(mockPool)
			actor, err := repo.GetActorByID(testContext(), tt.actorID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "NotFound" {
					assert.True(t, errors.Is(err, actors.ErrorNotFound))
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantActor.ID, actor.ID)
				assert.Equal(t, tt.wantActor.RussianName, actor.RussianName)
				assert.Equal(t, tt.wantActor.Photo, actor.Photo)
				assert.Equal(t, tt.wantActor.Height, actor.Height)
				assert.Equal(t, tt.wantActor.BirthDate, actor.BirthDate)
				assert.Equal(t, tt.wantActor.ZodiacSign, actor.ZodiacSign)
				assert.Equal(t, tt.wantActor.BirthPlace, actor.BirthPlace)
				assert.Equal(t, tt.wantActor.MaritalStatus, actor.MaritalStatus)
				if tt.wantActor.OriginalName != nil {
					assert.Equal(t, *tt.wantActor.OriginalName, *actor.OriginalName)
				} else {
					assert.Nil(t, actor.OriginalName)
				}
				if tt.wantActor.DeathDate != nil {
					assert.Equal(t, *tt.wantActor.DeathDate, *actor.DeathDate)
				} else {
					assert.Nil(t, actor.DeathDate)
				}
			}
		})
	}
}

func TestGetActorFilmsCount(t *testing.T) {
	actorID := uuid.NewV4()
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	originalName := "Tom Hanks"

	actorColumns := []string{
		"id", "russian_name", "original_name", "photo", "height",
		"birth_date", "death_date", "zodiac_sign", "birth_place", "marital_status",
	}

	tests := []struct {
		name       string
		actorID    uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
		wantCount  int
	}{
		{
			name:    "Success",
			actorID: actorID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				actorRows := pgxpoolmock.NewRows(actorColumns).
					AddRow(
						actorID,
						"Том Хэнкс",
						&originalName,
						"tom.jpg",
						183,
						birthDate,
						nil,
						"Козерог",
						"Конкорд, Калифорния, США",
						"Женат",
					).ToPgxRows()
				actorRows.Next()

				countRows := pgxpoolmock.NewRows([]string{"count"}).
					AddRow(15).
					ToPgxRows()
				countRows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetActorByID, actorID).
					Return(actorRows)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetActorFilmsCount, actorID).
					Return(countRows)
			},
			wantErr:   false,
			wantCount: 15,
		},
		{
			name:    "ActorNotFound",
			actorID: actorID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetActorByID, actorID).
					Return(MockRow{err: pgx.ErrNoRows})
			},
			wantErr:   true,
			wantCount: 0,
		},
		{
			name:    "NoFilmsFound",
			actorID: actorID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				actorRows := pgxpoolmock.NewRows(actorColumns).
					AddRow(
						actorID,
						"Том Хэнкс",
						&originalName,
						"tom.jpg",
						183,
						birthDate,
						nil,
						"Козерог",
						"Конкорд, Калифорния, США",
						"Женат",
					).ToPgxRows()
				actorRows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetActorByID, actorID).
					Return(actorRows)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetActorFilmsCount, actorID).
					Return(MockRow{err: pgx.ErrNoRows})
			},
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewActorRepository(mockPool)
			count, err := repo.GetActorFilmsCount(testContext(), tt.actorID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, count)
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
		wantErr    bool
		wantRating float64
	}{
		{
			name:   "Success",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"avg"}).
					AddRow(8.54321).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID).
					Return(rows)
			},
			wantErr:    false,
			wantRating: 8.5,
		},
		{
			name:   "ZeroRating",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"avg"}).
					AddRow(0.0).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID).
					Return(rows)
			},
			wantErr:    false,
			wantRating: 0.0,
		},
		{
			name:   "HighPrecisionRating",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"avg"}).
					AddRow(9.98765).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID).
					Return(rows)
			},
			wantErr:    false,
			wantRating: 10.0,
		},
		{
			name:   "NotFound",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID).
					Return(MockRow{err: pgx.ErrNoRows})
			},
			wantErr:    true,
			wantRating: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewActorRepository(mockPool)
			rating, err := repo.GetFilmAvgRating(testContext(), tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "NotFound" {
					assert.True(t, errors.Is(err, actors.ErrorNotFound))
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRating, rating)
			}
		})
	}
}

func TestGetFilmsByActor(t *testing.T) {
	actorID := uuid.NewV4()
	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()

	filmColumns := []string{"id", "cover", "title", "year", "genre"}

	tests := []struct {
		name       string
		actorID    uuid.UUID
		limit      int
		offset     int
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
		wantFilms  []models.MainPageFilm
	}{
		{
			name:    "Success",
			actorID: actorID,
			limit:   10,
			offset:  0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows(filmColumns).
					AddRow(filmID1, "film1.jpg", "Форрест Гамп", 1994, "Драма").
					AddRow(filmID2, "film2.jpg", "Зеленая миля", 1999, "Драма").
					ToPgxRows()

				ratingRows1 := pgxpoolmock.NewRows([]string{"avg"}).AddRow(8.8).ToPgxRows()
				ratingRows1.Next()
				ratingRows2 := pgxpoolmock.NewRows([]string{"avg"}).AddRow(8.6).ToPgxRows()
				ratingRows2.Next()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
					Return(filmRows, nil)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID1).
					Return(ratingRows1)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID2).
					Return(ratingRows2)
			},
			wantErr: false,
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID1,
					Cover:  "film1.jpg",
					Title:  "Форрест Гамп",
					Year:   1994,
					Genre:  "Драма",
					Rating: 8.8,
				},
				{
					ID:     filmID2,
					Cover:  "film2.jpg",
					Title:  "Зеленая миля",
					Year:   1999,
					Genre:  "Драма",
					Rating: 8.6,
				},
			},
		},
		{
			name:    "EmptyResult",
			actorID: actorID,
			limit:   10,
			offset:  0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows(filmColumns).ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
					Return(filmRows, nil)
			},
			wantErr:   false,
			wantFilms: []models.MainPageFilm{},
		},
		{
			name:    "QueryError",
			actorID: actorID,
			limit:   10,
			offset:  0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
					Return(nil, assert.AnError)
			},
			wantErr:   true,
			wantFilms: nil,
		},
		{
			name:    "NoRowsError",
			actorID: actorID,
			limit:   10,
			offset:  0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
					Return(nil, pgx.ErrNoRows)
			},
			wantErr:   true,
			wantFilms: nil,
		},
		{
			name:    "ScanError",
			actorID: actorID,
			limit:   10,
			offset:  0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows([]string{"id", "cover"}).
					AddRow(filmID1, "film1.jpg").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
					Return(filmRows, nil)
			},
			wantErr:   true,
			wantFilms: nil,
		},
		{
			name:    "RatingQueryError",
			actorID: actorID,
			limit:   10,
			offset:  0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows(filmColumns).
					AddRow(filmID1, "film1.jpg", "Форрест Гамп", 1994, "Драма").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
					Return(filmRows, nil)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID1).
					Return(MockRow{err: pgx.ErrNoRows})
			},
			wantErr: false,
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID1,
					Cover:  "film1.jpg",
					Title:  "Форрест Гамп",
					Year:   1994,
					Genre:  "Драма",
					Rating: 0.0,
				},
			},
		},
		{
			name:    "RatingInternalError",
			actorID: actorID,
			limit:   10,
			offset:  0,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows(filmColumns).
					AddRow(filmID1, "film1.jpg", "Форрест Гамп", 1994, "Драма").
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
					Return(filmRows, nil)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID1).
					Return(MockRow{err: assert.AnError})
			},
			wantErr: false,
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID1,
					Cover:  "film1.jpg",
					Title:  "Форрест Гамп",
					Year:   1994,
					Genre:  "Драма",
					Rating: 0.0,
				},
			},
		},
		{
			name:    "DifferentLimitOffset",
			actorID: actorID,
			limit:   5,
			offset:  10,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows(filmColumns).
					AddRow(filmID1, "film1.jpg", "Форрест Гамп", 1994, "Драма").
					ToPgxRows()

				ratingRows1 := pgxpoolmock.NewRows([]string{"avg"}).AddRow(8.8).ToPgxRows()
				ratingRows1.Next()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 5, 10).
					Return(filmRows, nil)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID1).
					Return(ratingRows1)
			},
			wantErr: false,
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID1,
					Cover:  "film1.jpg",
					Title:  "Форрест Гамп",
					Year:   1994,
					Genre:  "Драма",
					Rating: 8.8,
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

			repo := NewActorRepository(mockPool)
			films, err := repo.GetFilmsByActor(testContext(), tt.actorID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
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
