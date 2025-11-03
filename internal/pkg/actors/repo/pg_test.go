package repo

import (
	"context"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"
	"log/slog"
	"os"
	"testing"
	"time"

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
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantActor.ID, actor.ID)
				assert.Equal(t, tt.wantActor.RussianName, actor.RussianName)
				assert.Equal(t, tt.wantActor.Photo, actor.Photo)
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
					AddRow(8.5).
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
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRating, rating)
			}
		})
	}
}

func TestGetFilmsByActor(t *testing.T) {
	actorID := uuid.NewV4()
	filmID := uuid.NewV4()
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	originalName := "Tom Hanks"

	actorColumns := []string{
		"id", "russian_name", "original_name", "photo", "height",
		"birth_date", "death_date", "zodiac_sign", "birth_place", "marital_status",
	}

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

				filmRows := pgxpoolmock.NewRows(filmColumns).
					AddRow(filmID, "film1.jpg", "Форрест Гамп", 1994, "Драма").
					ToPgxRows()

				ratingRows := pgxpoolmock.NewRows([]string{"avg"}).AddRow(8.8).ToPgxRows()
				ratingRows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetActorByID, actorID).
					Return(actorRows)

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
					Return(filmRows, nil)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRating, filmID).
					Return(ratingRows)
			},
			wantErr: false,
			wantFilms: []models.MainPageFilm{
				{
					ID:     filmID,
					Cover:  "film1.jpg",
					Title:  "Форрест Гамп",
					Year:   1994,
					Genre:  "Драма",
					Rating: 8.8,
				},
			},
		},
		{
			name:    "EmptyResult",
			actorID: actorID,
			limit:   10,
			offset:  0,
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

				filmRows := pgxpoolmock.NewRows(filmColumns).ToPgxRows()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetActorByID, actorID).
					Return(actorRows)

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
					Return(filmRows, nil)
			},
			wantErr:   false,
			wantFilms: []models.MainPageFilm{},
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
					assert.Equal(t, tt.wantFilms[0].ID, films[0].ID)
					assert.Equal(t, tt.wantFilms[0].Title, films[0].Title)
					assert.Equal(t, tt.wantFilms[0].Rating, films[0].Rating)
				}
			}
		})
	}
}
