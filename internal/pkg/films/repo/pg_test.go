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

func TestGetPromoFilmByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	filmID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := pgxpoolmock.NewRows([]string{
		"id", "image", "title", "short_description", "year", "genre", "duration", "created_at", "updated_at",
	}).
		AddRow(
			filmID,
			"/static/poster.jpg",
			"Test Film",
			"Short description",
			2023,
			"Drama",
			120,
			createdAt,
			updatedAt,
		).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetPromoFilmByIDQuery, filmID).
		Return(rows)

	film, err := repo.GetPromoFilmByID(testContext(), filmID)

	assert.NoError(t, err)
	assert.Equal(t, filmID, film.ID)
	assert.Equal(t, "Test Film", film.Title)
	assert.Equal(t, "Drama", film.Genre)
}

func TestGetFilmByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	filmID := uuid.NewV4()
	countryID := uuid.NewV4()
	genreID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	originalTitle := "Test Film Original"
	trailerURL := ""
	slogan := "Great film slogan"
	image1 := "/static/image1.jpg"
	image2 := "/static/image2.jpg"
	image3 := "/static/image3.jpg"

	rows := pgxpoolmock.NewRows([]string{
		"id", "title", "original_title", "cover", "poster",
		"short_description", "description", "age_category", "budget",
		"worldwide_fees", "trailer_url", "year", "country_id",
		"genre_id", "slogan", "duration", "image1", "image2",
		"image3", "created_at", "updated_at",
	}).
		AddRow(
			filmID,
			"Test Film",
			&originalTitle,
			"/static/cover.jpg",
			"/static/poster.jpg",
			"Short description",
			"Full description",
			"18+",
			1000000,
			5000000,
			&trailerURL,
			2023,
			countryID,
			genreID,
			&slogan,
			120,
			&image1,
			&image2,
			&image3,
			createdAt,
			updatedAt,
		).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetFilmByIDQuery, filmID).
		Return(rows)

	film, err := repo.GetFilmByID(testContext(), filmID)

	assert.NoError(t, err)
	assert.Equal(t, filmID, film.ID)
	assert.Equal(t, "Test Film", film.Title)
	assert.Equal(t, &originalTitle, film.OriginalTitle)
	assert.Equal(t, "18+", film.AgeCategory)
}

func TestGetGenreTitle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	genreID := uuid.NewV4()
	genreTitle := "Drama"

	rows := pgxpoolmock.NewRows([]string{"title"}).
		AddRow(genreTitle).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetGenreTitleQuery, genreID).
		Return(rows)

	title, err := repo.GetGenreTitle(testContext(), genreID)

	assert.NoError(t, err)
	assert.Equal(t, genreTitle, title)
}

func TestGetFilmAvgRating_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

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

func TestGetFilmsWithPagination_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()
	limit := 10
	offset := 0

	mainRows := pgxpoolmock.NewRows([]string{
		"id", "cover", "title", "year", "genre_title",
	}).
		AddRow(filmID1, "/static/cover1.jpg", "Film 1", 2023, "Drama").
		AddRow(filmID2, "/static/cover2.jpg", "Film 2", 2022, "Comedy").
		ToPgxRows()

	mockPool.EXPECT().
		Query(gomock.Any(), GetFilmsWithPaginationQuery, limit, offset).
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

	films, err := repo.GetFilmsWithPagination(testContext(), limit, offset)

	assert.NoError(t, err)
	assert.Len(t, films, 2)
	assert.Equal(t, filmID1, films[0].ID)
	assert.Equal(t, "Film 1", films[0].Title)
	assert.Equal(t, 8.5, films[0].Rating)
	assert.Equal(t, filmID2, films[1].ID)
	assert.Equal(t, "Film 2", films[1].Title)
	assert.Equal(t, 7.8, films[1].Rating)
}

func TestGetFilmFeedbacks_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	filmID := uuid.NewV4()
	userID := uuid.NewV4()
	feedbackID := uuid.NewV4()
	limit := 10
	offset := 0
	createdAt := time.Now()
	updatedAt := time.Now()

	title := "Great film!"
	text := "Amazing storyline and acting"

	rows := pgxpoolmock.NewRows([]string{
		"id", "user_id", "film_id", "title", "text", "rating",
		"created_at", "updated_at", "user_login", "user_avatar",
	}).
		AddRow(
			feedbackID,
			userID,
			filmID,
			&title,
			&text,
			9,
			createdAt,
			updatedAt,
			"testuser",
			"/static/avatar.jpg",
		).
		ToPgxRows()

	mockPool.EXPECT().
		Query(gomock.Any(), GetFilmFeedbacksQuery, filmID, limit, offset).
		Return(rows, nil)

	feedbacks, err := repo.GetFilmFeedbacks(testContext(), filmID, limit, offset)

	assert.NoError(t, err)
	assert.Len(t, feedbacks, 1)
	assert.Equal(t, feedbackID, feedbacks[0].ID)
	assert.Equal(t, userID, feedbacks[0].UserID)
	assert.Equal(t, filmID, feedbacks[0].FilmID)
	assert.Equal(t, &title, feedbacks[0].Title)
	assert.Equal(t, "testuser", feedbacks[0].UserLogin)
}

func TestCheckUserFeedbackExists_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	feedbackID := uuid.NewV4()
	userID := uuid.NewV4()
	filmID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	title := "Good film"
	text := "Nice cinematography"

	rows := pgxpoolmock.NewRows([]string{
		"id", "user_id", "film_id", "title", "text", "rating",
		"created_at", "updated_at", "user_login", "user_avatar",
	}).
		AddRow(
			feedbackID,
			userID,
			filmID,
			&title,
			&text,
			8,
			createdAt,
			updatedAt,
			"testuser",
			"/static/avatar.jpg",
		).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), CheckUserFeedbackExistsQuery, userID, filmID).
		Return(rows)

	feedback, err := repo.CheckUserFeedbackExists(testContext(), userID, filmID)

	assert.NoError(t, err)
	assert.Equal(t, feedbackID, feedback.ID)
	assert.Equal(t, userID, feedback.UserID)
	assert.Equal(t, filmID, feedback.FilmID)
	assert.Equal(t, &title, feedback.Title)
	assert.Equal(t, "testuser", feedback.UserLogin)
}

func TestUpdateFeedback_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	feedbackID := uuid.NewV4()
	userID := uuid.NewV4()
	filmID := uuid.NewV4()
	title := "Updated Title"
	text := "Updated text"

	feedback := models.FilmFeedback{
		ID:     feedbackID,
		UserID: userID,
		FilmID: filmID,
		Title:  &title,
		Text:   &text,
		Rating: 9,
	}

	mockPool.EXPECT().
		Exec(gomock.Any(), UpdateFeedbackQuery, feedback.Title, feedback.Text, feedback.Rating, feedback.ID).
		Return(nil, nil)

	err := repo.UpdateFeedback(testContext(), feedback)

	assert.NoError(t, err)
}

func TestCreateFeedback_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	feedbackID := uuid.NewV4()
	userID := uuid.NewV4()
	filmID := uuid.NewV4()
	title := "New Title"
	text := "Great film"

	feedback := models.FilmFeedback{
		ID:     feedbackID,
		UserID: userID,
		FilmID: filmID,
		Title:  &title,
		Text:   &text,
		Rating: 9,
	}

	mockPool.EXPECT().
		Exec(gomock.Any(), CreateFeedbackQuery,
			feedback.ID, feedback.UserID, feedback.FilmID,
			feedback.Title, feedback.Text, feedback.Rating).
		Return(nil, nil)

	err := repo.CreateFeedback(testContext(), feedback)

	assert.NoError(t, err)
}

func TestSetRating_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	feedbackID := uuid.NewV4()
	userID := uuid.NewV4()
	filmID := uuid.NewV4()

	feedback := models.FilmFeedback{
		ID:     feedbackID,
		UserID: userID,
		FilmID: filmID,
		Rating: 8,
	}

	mockPool.EXPECT().
		Exec(gomock.Any(), SetRatingQuery,
			feedback.ID, feedback.UserID, feedback.FilmID, feedback.Rating).
		Return(nil, nil)

	err := repo.SetRating(testContext(), feedback)

	assert.NoError(t, err)
}

func TestGetUserByLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	userID := uuid.NewV4()
	login := "testuser"
	avatar := "/static/default.jpg"
	createdAt := time.Now()
	updatedAt := time.Now()

	rows := pgxpoolmock.NewRows([]string{
		"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at",
	}).
		AddRow(userID, 1, login, []byte("hash"), &avatar, createdAt, updatedAt).
		ToPgxRows()
	rows.Next()

	mockPool.EXPECT().
		QueryRow(gomock.Any(), GetUserByLoginQuery, login).
		Return(rows)

	user, err := repo.GetUserByLogin(testContext(), login)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, login, user.Login)
	assert.Equal(t, &avatar, user.Avatar)
}

func TestGetFilmsWithPagination_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	limit := 10
	offset := 0

	mockPool.EXPECT().
		Query(gomock.Any(), GetFilmsWithPaginationQuery, limit, offset).
		Return(nil, assert.AnError)

	films, err := repo.GetFilmsWithPagination(testContext(), limit, offset)

	assert.Error(t, err)
	assert.Nil(t, films)
}

func TestGetFilmFeedbacks_QueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	repo := NewFilmRepository(mockPool)

	filmID := uuid.NewV4()
	limit := 10
	offset := 0

	mockPool.EXPECT().
		Query(gomock.Any(), GetFilmFeedbacksQuery, filmID, limit, offset).
		Return(nil, assert.AnError)

	feedbacks, err := repo.GetFilmFeedbacks(testContext(), filmID, limit, offset)

	assert.Error(t, err)
	assert.Nil(t, feedbacks)
}
