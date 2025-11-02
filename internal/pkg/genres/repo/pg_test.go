package repo

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

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

func TestGetGenreByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewGenreRepository(mockPool)

	genreID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

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

	genre, err := repo.GetGenreByID(testContext(), genreID)

	assert.NoError(t, err)
	assert.Equal(t, genreID, genre.ID)
	assert.Equal(t, "Drama", genre.Title)
	assert.Equal(t, "Drama films description", genre.Description)
	assert.Equal(t, "/static/drama.png", genre.Icon)
}

func TestGetGenresWithPagination_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewGenreRepository(mockPool)

	genreID1 := uuid.NewV4()
	genreID2 := uuid.NewV4()
	limit := 10
	offset := 0
	createdAt := time.Now()
	updatedAt := time.Now()

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

	genres, err := repo.GetGenresWithPagination(testContext(), limit, offset)

	assert.NoError(t, err)
	assert.Len(t, genres, 2)
	assert.Equal(t, genreID1, genres[0].ID)
	assert.Equal(t, "Action", genres[0].Title)
	assert.Equal(t, genreID2, genres[1].ID)
	assert.Equal(t, "Comedy", genres[1].Title)
}

func TestGetFilmAvgRating_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewGenreRepository(mockPool)

	filmID := uuid.NewV4()
	avgRating := 8.5

	rows := pgxpoolmock.NewRows([]string{"coalesce"}).
		AddRow(avgRating).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID).
		Return(rows)

	rating, err := repo.GetFilmAvgRating(testContext(), filmID)

	assert.NoError(t, err)
	assert.Equal(t, 8.5, rating)
}

func TestGetFilmsByGenre_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewGenreRepository(mockPool)

	genreID := uuid.NewV4()
	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()
	limit := 10
	offset := 0

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

	films, err := repo.GetFilmsByGenre(testContext(), genreID, limit, offset)

	assert.NoError(t, err)
	assert.Len(t, films, 2)
	assert.Equal(t, filmID1, films[0].ID)
	assert.Equal(t, "Film 1", films[0].Title)
	assert.Equal(t, 8.5, films[0].Rating)
	assert.Equal(t, filmID2, films[1].ID)
	assert.Equal(t, "Film 2", films[1].Title)
	assert.Equal(t, 7.8, films[1].Rating)
}

func TestGetGenresWithPagination_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewGenreRepository(mockPool)

	limit := 10
	offset := 0

	mockPool.EXPECT().
		Query(gomock.Any(), GetGenresWithPaginationQuery, limit, offset).
		Return(nil, assert.AnError)

	genres, err := repo.GetGenresWithPagination(testContext(), limit, offset)

	assert.Error(t, err)
	assert.Nil(t, genres)
}

func TestGetFilmsByGenre_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewGenreRepository(mockPool)

	genreID := uuid.NewV4()
	limit := 10
	offset := 0

	mockPool.EXPECT().
		Query(gomock.Any(), GetFilmsByGenreQuery, genreID, limit, offset).
		Return(nil, assert.AnError)

	films, err := repo.GetFilmsByGenre(testContext(), genreID, limit, offset)

	assert.Error(t, err)
	assert.Nil(t, films)
}

func TestGetFilmsByGenre_ScanError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewGenreRepository(mockPool)

	genreID := uuid.NewV4()
	filmID := uuid.NewV4()
	limit := 10
	offset := 0

	mainRows := pgxpoolmock.NewRows([]string{
		"id", "cover", "title", "year", "genre_title",
	}).
		AddRow(filmID, "", "", 0, "").
		ToPgxRows()

	mockPool.EXPECT().
		Query(gomock.Any(), GetFilmsByGenreQuery, genreID, limit, offset).
		Return(mainRows, nil)

	ratingRows := pgxpoolmock.NewRows([]string{"coalesce"}).
		AddRow(0.0).
		ToPgxRows()
	ratingRows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID).
		Return(ratingRows)

	films, err := repo.GetFilmsByGenre(testContext(), genreID, limit, offset)

	assert.NoError(t, err)
	assert.Len(t, films, 1)
	assert.Equal(t, filmID, films[0].ID)
	assert.Equal(t, 0.0, films[0].Rating)
}
