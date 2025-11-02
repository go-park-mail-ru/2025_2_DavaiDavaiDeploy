package repo

import (
	"context"
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

func TestGetActorByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewActorRepository(mockPool)

	actorID := uuid.NewV4()
	birthDate := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	originalName := "Tom Hanks"

	columns := []string{
		"id", "russian_name", "original_name", "photo", "height",
		"birth_date", "death_date", "zodiac_sign", "birth_place", "marital_status",
	}
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

	actor, err := repo.GetActorByID(testContext(), actorID)

	assert.NoError(t, err)
	assert.Equal(t, actorID, actor.ID)
	assert.Equal(t, "Том Хэнкс", actor.RussianName)
	assert.Equal(t, "tom.jpg", actor.Photo)
}

func TestGetActorFilmsCount_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewActorRepository(mockPool)

	actorID := uuid.NewV4()

	rows := pgxpoolmock.NewRows([]string{"count"}).
		AddRow(15).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetActorFilmsCount, actorID).
		Return(rows)

	count, err := repo.GetActorFilmsCount(testContext(), actorID)

	assert.NoError(t, err)
	assert.Equal(t, 15, count)
}

func TestGetFilmAvgRating_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewActorRepository(mockPool)

	filmID := uuid.NewV4()

	rows := pgxpoolmock.NewRows([]string{"avg"}).
		AddRow(8.5).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetFilmAvgRating, filmID).
		Return(rows)

	rating, err := repo.GetFilmAvgRating(testContext(), filmID)

	assert.NoError(t, err)
	assert.Equal(t, 8.5, rating)
}

func TestGetFilmsByActor_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewActorRepository(mockPool)

	actorID := uuid.NewV4()
	filmID := uuid.NewV4()

	columns := []string{"id", "cover", "title", "year", "genre"}
	filmRows := pgxpoolmock.NewRows(columns).
		AddRow(filmID, "film1.jpg", "Форрест Гамп", 1994, "Драма").
		ToPgxRows()

	mockPool.EXPECT().
		Query(gomock.Any(), GetFilmsByActor, actorID, 10, 0).
		Return(filmRows, nil)

	ratingRows := pgxpoolmock.NewRows([]string{"avg"}).AddRow(8.8).ToPgxRows()
	ratingRows.Next()
	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetFilmAvgRating, filmID).
		Return(ratingRows)

	films, err := repo.GetFilmsByActor(testContext(), actorID, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, films, 1)
	assert.Equal(t, filmID, films[0].ID)
	assert.Equal(t, "Форрест Гамп", films[0].Title)
	assert.Equal(t, 8.8, films[0].Rating)
}
