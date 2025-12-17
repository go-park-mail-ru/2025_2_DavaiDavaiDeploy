package repo

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/middleware/logger"

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

func TestGetPromoFilmByID(t *testing.T) {
	filmID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	tests := []struct {
		name       string
		filmID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilm   models.PromoFilm
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
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
			},
			wantFilm: models.PromoFilm{
				ID:               filmID,
				Image:            "/static/poster.jpg",
				Title:            "Test Film",
				ShortDescription: "Short description",
				Year:             2023,
				Genre:            "Drama",
				Duration:         120,
				CreatedAt:        createdAt,
				UpdatedAt:        updatedAt,
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

			repo := NewFilmRepository(mockPool)
			film, err := repo.GetPromoFilmByID(testContext(), tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFilm.ID, film.ID)
				assert.Equal(t, tt.wantFilm.Title, film.Title)
				assert.Equal(t, tt.wantFilm.Genre, film.Genre)
			}
		})
	}
}

func TestGetFilmByID(t *testing.T) {
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

	tests := []struct {
		name       string
		filmID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilm   models.Film
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
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
			},
			wantFilm: models.Film{
				ID:               filmID,
				Title:            "Test Film",
				OriginalTitle:    &originalTitle,
				Cover:            "/static/cover.jpg",
				Poster:           "/static/poster.jpg",
				ShortDescription: "Short description",
				Description:      "Full description",
				AgeCategory:      "18+",
				Budget:           1000000,
				WorldwideFees:    5000000,
				TrailerURL:       &trailerURL,
				Year:             2023,
				CountryID:        countryID,
				GenreID:          genreID,
				Slogan:           &slogan,
				Duration:         120,
				Image1:           &image1,
				Image2:           &image2,
				Image3:           &image3,
				CreatedAt:        createdAt,
				UpdatedAt:        updatedAt,
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

			repo := NewFilmRepository(mockPool)
			film, err := repo.GetFilmByID(testContext(), tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFilm.ID, film.ID)
				assert.Equal(t, tt.wantFilm.Title, film.Title)
				assert.Equal(t, tt.wantFilm.OriginalTitle, film.OriginalTitle)
				assert.Equal(t, tt.wantFilm.AgeCategory, film.AgeCategory)
			}
		})
	}
}

func TestGetGenreTitle(t *testing.T) {
	genreID := uuid.NewV4()
	genreTitle := "Drama"

	tests := []struct {
		name       string
		genreID    uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantTitle  string
		wantErr    bool
		errType    error
	}{
		{
			name:    "Success",
			genreID: genreID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"title"}).
					AddRow(genreTitle).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetGenreTitleQuery, genreID).
					Return(rows)
			},
			wantTitle: genreTitle,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			title, err := repo.GetGenreTitle(testContext(), tt.genreID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTitle, title)
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
		errType    error
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

			repo := NewFilmRepository(mockPool)
			rating, err := repo.GetFilmAvgRating(testContext(), tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRating, rating)
			}
		})
	}
}

func TestGetFilmsWithPagination(t *testing.T) {
	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()
	limit := 10
	offset := 0

	tests := []struct {
		name       string
		limit      int
		offset     int
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilms  []models.MainPageFilm
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
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
			},
			wantFilms: []models.MainPageFilm{
				{ID: filmID1, Cover: "/static/cover1.jpg", Title: "Film 1", Year: 2023, Genre: "Drama", Rating: 8.5},
				{ID: filmID2, Cover: "/static/cover2.jpg", Title: "Film 2", Year: 2022, Genre: "Comedy", Rating: 7.8},
			},
			wantErr: false,
		},
		{
			name:   "QueryError",
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsWithPaginationQuery, limit, offset).
					Return(nil, assert.AnError)
			},
			wantFilms: nil,
			wantErr:   true,
			errType:   films.ErrorInternalServerError,
		},
		{
			name:   "EmptyResult",
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mainRows := pgxpoolmock.NewRows([]string{
					"id", "cover", "title", "year", "genre_title",
				}).ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsWithPaginationQuery, limit, offset).
					Return(mainRows, nil)
			},
			wantFilms: []models.MainPageFilm{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			films, err := repo.GetFilmsWithPagination(testContext(), tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
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

func TestGetFilmFeedbacks(t *testing.T) {
	filmID := uuid.NewV4()
	userID := uuid.NewV4()
	feedbackID := uuid.NewV4()
	limit := 10
	offset := 0
	createdAt := time.Now()
	updatedAt := time.Now()

	title := "Great film!"
	text := "Amazing storyline and acting"

	tests := []struct {
		name          string
		filmID        uuid.UUID
		limit         int
		offset        int
		repoMocker    func(*pgxpoolmock.MockPgxPool)
		wantFeedbacks []models.FilmFeedback
		wantErr       bool
		errType       error
	}{
		{
			name:   "Success",
			filmID: filmID,
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				feedbackRows := pgxpoolmock.NewRows([]string{
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
					Return(feedbackRows, nil)
			},
			wantFeedbacks: []models.FilmFeedback{
				{
					ID:         feedbackID,
					UserID:     userID,
					FilmID:     filmID,
					Title:      &title,
					Text:       &text,
					Rating:     9,
					CreatedAt:  createdAt,
					UpdatedAt:  updatedAt,
					UserLogin:  "testuser",
					UserAvatar: "/static/avatar.jpg",
				},
			},
			wantErr: false,
		},
		{
			name:   "NoRows",
			filmID: filmID,
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "user_id", "film_id", "title", "text", "rating",
					"created_at", "updated_at", "user_login", "user_avatar",
				}).ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmFeedbacksQuery, filmID, limit, offset).
					Return(rows, nil)
			},
			wantFeedbacks: []models.FilmFeedback{},
			wantErr:       false,
		},
		{
			name:   "QueryError",
			filmID: filmID,
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmFeedbacksQuery, filmID, limit, offset).
					Return(nil, assert.AnError)
			},
			wantFeedbacks: nil,
			wantErr:       true,
			errType:       films.ErrorInternalServerError,
		},
		{
			name:   "NotFoundError",
			filmID: filmID,
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmFeedbacksQuery, filmID, limit, offset).
					Return(nil, pgx.ErrNoRows)
			},
			wantFeedbacks: nil,
			wantErr:       true,
			errType:       films.ErrorNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			feedbacks, err := repo.GetFilmFeedbacks(testContext(), tt.filmID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
				assert.Nil(t, feedbacks)
			} else {
				assert.NoError(t, err)
				assert.Len(t, feedbacks, len(tt.wantFeedbacks))
				if len(feedbacks) > 0 {
					assert.Equal(t, tt.wantFeedbacks[0].ID, feedbacks[0].ID)
					assert.Equal(t, tt.wantFeedbacks[0].UserID, feedbacks[0].UserID)
					assert.Equal(t, tt.wantFeedbacks[0].FilmID, feedbacks[0].FilmID)
					assert.Equal(t, tt.wantFeedbacks[0].Title, feedbacks[0].Title)
					assert.Equal(t, tt.wantFeedbacks[0].UserLogin, feedbacks[0].UserLogin)
				}
			}
		})
	}
}

func TestCheckUserFeedbackExists(t *testing.T) {
	feedbackID := uuid.NewV4()
	userID := uuid.NewV4()
	filmID := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	title := "Good film"
	text := "Nice cinematography"

	tests := []struct {
		name         string
		userID       uuid.UUID
		filmID       uuid.UUID
		repoMocker   func(*pgxpoolmock.MockPgxPool)
		wantFeedback models.FilmFeedback
		wantErr      bool
		errType      error
	}{
		{
			name:   "Success",
			userID: userID,
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
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
			},
			wantFeedback: models.FilmFeedback{
				ID:         feedbackID,
				UserID:     userID,
				FilmID:     filmID,
				Title:      &title,
				Text:       &text,
				Rating:     8,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				UserLogin:  "testuser",
				UserAvatar: "/static/avatar.jpg",
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

			repo := NewFilmRepository(mockPool)
			feedback, err := repo.CheckUserFeedbackExists(testContext(), tt.userID, tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFeedback.ID, feedback.ID)
				assert.Equal(t, tt.wantFeedback.UserID, feedback.UserID)
				assert.Equal(t, tt.wantFeedback.FilmID, feedback.FilmID)
				assert.Equal(t, tt.wantFeedback.Title, feedback.Title)
				assert.Equal(t, tt.wantFeedback.UserLogin, feedback.UserLogin)
			}
		})
	}
}

func TestUpdateFeedback(t *testing.T) {
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

	tests := []struct {
		name       string
		feedback   models.FilmFeedback
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
		errType    error
	}{
		{
			name:     "Success",
			feedback: feedback,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), UpdateFeedbackQuery, feedback.Title, feedback.Text, feedback.Rating, feedback.ID).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name:     "ExecError",
			feedback: feedback,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), UpdateFeedbackQuery, feedback.Title, feedback.Text, feedback.Rating, feedback.ID).
					Return(nil, assert.AnError)
			},
			wantErr: true,
			errType: films.ErrorInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			err := repo.UpdateFeedback(testContext(), tt.feedback)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateFeedback(t *testing.T) {
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

	tests := []struct {
		name       string
		feedback   models.FilmFeedback
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
		errType    error
	}{
		{
			name:     "Success",
			feedback: feedback,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), CreateFeedbackQuery,
						feedback.ID, feedback.UserID, feedback.FilmID,
						feedback.Title, feedback.Text, feedback.Rating).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name:     "ExecError",
			feedback: feedback,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), CreateFeedbackQuery,
						feedback.ID, feedback.UserID, feedback.FilmID,
						feedback.Title, feedback.Text, feedback.Rating).
					Return(nil, assert.AnError)
			},
			wantErr: true,
			errType: films.ErrorInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			err := repo.CreateFeedback(testContext(), tt.feedback)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetRating(t *testing.T) {
	feedbackID := uuid.NewV4()
	userID := uuid.NewV4()
	filmID := uuid.NewV4()

	feedback := models.FilmFeedback{
		ID:     feedbackID,
		UserID: userID,
		FilmID: filmID,
		Rating: 8,
	}

	tests := []struct {
		name       string
		feedback   models.FilmFeedback
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
		errType    error
	}{
		{
			name:     "Success",
			feedback: feedback,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), SetRatingQuery,
						feedback.ID, feedback.UserID, feedback.FilmID, feedback.Rating).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name:     "ExecError",
			feedback: feedback,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), SetRatingQuery,
						feedback.ID, feedback.UserID, feedback.FilmID, feedback.Rating).
					Return(nil, assert.AnError)
			},
			wantErr: true,
			errType: films.ErrorInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			err := repo.SetRating(testContext(), tt.feedback)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserByLogin(t *testing.T) {
	userID := uuid.NewV4()
	login := "testuser"
	avatar := "/static/default.jpg"
	createdAt := time.Now()
	updatedAt := time.Now()

	tests := []struct {
		name       string
		login      string
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantUser   models.User
		wantErr    bool
		errType    error
	}{
		{
			name:  "Success",
			login: login,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "version", "login", "password_hash", "avatar", "created_at", "updated_at",
				}).
					AddRow(userID, 1, login, []byte("hash"), avatar, createdAt, updatedAt).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetUserByLoginQuery, login).
					Return(rows)
			},
			wantUser: models.User{
				ID:           userID,
				Version:      1,
				Login:        login,
				PasswordHash: []byte("hash"),
				Avatar:       avatar,
				CreatedAt:    createdAt,
				UpdatedAt:    updatedAt,
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

			repo := NewFilmRepository(mockPool)
			user, err := repo.GetUserByLogin(testContext(), tt.login)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Login, user.Login)
				assert.Equal(t, tt.wantUser.Avatar, user.Avatar)
			}
		})
	}
}

func TestCheckUserLikeExists(t *testing.T) {
	likeID := uuid.NewV4()
	userID := uuid.NewV4()
	filmID := uuid.NewV4()

	tests := []struct {
		name       string
		userID     uuid.UUID
		filmID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantLike   models.FilmFeedback
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			userID: userID,
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{"id", "user_id", "film_id"}).
					AddRow(likeID, userID, filmID).
					ToPgxRows()
				rows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), CheckUserLikeExistsQuery, userID, filmID).
					Return(rows)
			},
			wantLike: models.FilmFeedback{
				ID:     likeID,
				UserID: userID,
				FilmID: filmID,
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

			repo := NewFilmRepository(mockPool)
			like, err := repo.CheckUserLikeExists(testContext(), tt.userID, tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLike.ID, like.ID)
				assert.Equal(t, tt.wantLike.UserID, like.UserID)
				assert.Equal(t, tt.wantLike.FilmID, like.FilmID)
			}
		})
	}
}

func TestSaveFilm(t *testing.T) {
	userID := uuid.NewV4()
	filmID := uuid.NewV4()

	tests := []struct {
		name       string
		userID     uuid.UUID
		filmID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			userID: userID,
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), InsertIntoSavedQuery, userID, filmID).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name:   "ExecError",
			userID: userID,
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), InsertIntoSavedQuery, userID, filmID).
					Return(nil, assert.AnError)
			},
			wantErr: true,
			errType: films.ErrorBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			err := repo.SaveFilm(testContext(), tt.userID, tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveFilm(t *testing.T) {
	userID := uuid.NewV4()
	filmID := uuid.NewV4()

	tests := []struct {
		name       string
		userID     uuid.UUID
		filmID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			userID: userID,
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), DeleteFromSavedQuery, userID, filmID).
					Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name:   "ExecError",
			userID: userID,
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Exec(gomock.Any(), DeleteFromSavedQuery, userID, filmID).
					Return(nil, assert.AnError)
			},
			wantErr: true,
			errType: films.ErrorInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			err := repo.RemoveFilm(testContext(), tt.userID, tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetFilmsForCalendar(t *testing.T) {
	limit := 10
	offset := 0

	tests := []struct {
		name       string
		limit      int
		offset     int
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilms  []models.FilmInCalendar
		wantErr    bool
		errType    error
	}{
		{
			name:   "QueryError",
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsWithDateOfReleaseQuery, limit, offset).
					Return(nil, assert.AnError)
			},
			wantFilms: nil,
			wantErr:   true,
			errType:   films.ErrorInternalServerError,
		},
		{
			name:   "EmptyResult",
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "cover", "title", "original_title", "short_description", "release_date",
				}).ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsWithDateOfReleaseQuery, limit, offset).
					Return(rows, nil)
			},
			wantFilms: []models.FilmInCalendar{},
			wantErr:   false,
		},
		{
			name:   "ScanErrorContinues",
			limit:  limit,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "cover", "title", "original_title", "short_description", "release_date",
				}).
					AddRow(
						nil, // Это вызовет ошибку сканирования
						"/static/cover.jpg",
						"Upcoming Film",
						"Original Title",
						"Short description",
						time.Now(),
					).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsWithDateOfReleaseQuery, limit, offset).
					Return(rows, nil)
			},
			wantFilms: []models.FilmInCalendar{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			films, err := repo.GetFilmsForCalendar(testContext(), tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
				assert.Nil(t, films)
			} else {
				assert.NoError(t, err)
				assert.Len(t, films, len(tt.wantFilms))
				if len(films) > 0 {
					assert.Equal(t, tt.wantFilms[0].Title, films[0].Title)
				}
			}
		})
	}
}

func TestGetUsersFavFilms(t *testing.T) {
	userID := uuid.NewV4()
	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()

	tests := []struct {
		name       string
		userID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilms  []models.FavFilm
		wantErr    bool
	}{
		{
			name:   "Success",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "title", "genre_title", "year", "duration", "image", "short_description", "rating",
				}).
					AddRow(
						filmID1,
						"Фильм 1",
						"Драма",
						2023,
						120,
						"/static/image1.jpg",
						"Короткое описание 1",
						8.5,
					).
					AddRow(
						filmID2,
						"Фильм 2",
						"Комедия",
						2022,
						110,
						"/static/image2.jpg",
						"Короткое описание 2",
						7.8,
					).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetUsersFavFilmsQuery, userID).
					Return(rows, nil)
			},
			wantFilms: []models.FavFilm{
				{
					ID:               filmID1,
					Title:            "Фильм 1",
					Genre:            "Драма",
					Year:             2023,
					Duration:         120,
					Image:            "/static/image1.jpg",
					ShortDescription: "Короткое описание 1",
					Rating:           8.5,
				},
				{
					ID:               filmID2,
					Title:            "Фильм 2",
					Genre:            "Комедия",
					Year:             2022,
					Duration:         110,
					Image:            "/static/image2.jpg",
					ShortDescription: "Короткое описание 2",
					Rating:           7.8,
				},
			},
			wantErr: false,
		},
		{
			name:   "QueryError",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetUsersFavFilmsQuery, userID).
					Return(nil, assert.AnError)
			},
			wantFilms: []models.FavFilm{},
			wantErr:   false, // Метод возвращает пустой слайс при ошибке
		},
		{
			name:   "EmptyResult",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "title", "genre_title", "year", "duration", "image", "short_description", "rating",
				}).ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetUsersFavFilmsQuery, userID).
					Return(rows, nil)
			},
			wantFilms: []models.FavFilm{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			films, err := repo.GetUsersFavFilms(testContext(), tt.userID)

			assert.NoError(t, err)
			assert.Len(t, films, len(tt.wantFilms))
			if len(films) > 0 {
				assert.Equal(t, tt.wantFilms[0].Title, films[0].Title)
				assert.Equal(t, tt.wantFilms[0].Genre, films[0].Genre)
				assert.Equal(t, tt.wantFilms[0].Rating, films[0].Rating)
			}
		})
	}
}

func TestGetFilmsWithCursorPagination(t *testing.T) {
	cursor := time.Now()
	count := 10

	tests := []struct {
		name       string
		cursor     time.Time
		count      int
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilms  []models.MainPageFilm
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			cursor: cursor,
			count:  count,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmID1 := uuid.NewV4()
				filmID2 := uuid.NewV4()

				mainRows := pgxpoolmock.NewRows([]string{
					"id", "cover", "title", "year", "genre_title", "created_at",
				}).
					AddRow(
						filmID1,
						"/static/cover1.jpg",
						"Film 1",
						2023,
						"Drama",
						cursor.Add(-24*time.Hour),
					).
					AddRow(
						filmID2,
						"/static/cover2.jpg",
						"Film 2",
						2022,
						"Comedy",
						cursor.Add(-48*time.Hour),
					).
					ToPgxRows()

				ratingRows1 := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(8.5).
					ToPgxRows()
				ratingRows1.Next()

				ratingRows2 := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(7.8).
					ToPgxRows()
				ratingRows2.Next()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsWithCursorPaginationQuery, cursor, count+1).
					Return(mainRows, nil)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID1).
					Return(ratingRows1)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID2).
					Return(ratingRows2)
			},
			wantFilms: []models.MainPageFilm{
				{
					ID:        uuid.NewV4(),
					Cover:     "/static/cover1.jpg",
					Title:     "Film 1",
					Year:      2023,
					Genre:     "Drama",
					Rating:    8.5,
					CreatedAt: cursor.Add(-24 * time.Hour),
				},
				{
					ID:        uuid.NewV4(),
					Cover:     "/static/cover2.jpg",
					Title:     "Film 2",
					Year:      2022,
					Genre:     "Comedy",
					Rating:    7.8,
					CreatedAt: cursor.Add(-48 * time.Hour),
				},
			},
			wantErr: false,
		},
		{
			name:   "QueryError",
			cursor: cursor,
			count:  count,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmsWithCursorPaginationQuery, cursor, count+1).
					Return(nil, assert.AnError)
			},
			wantFilms: nil,
			wantErr:   true,
			errType:   films.ErrorInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			films, err := repo.GetFilmsWithCursorPagination(testContext(), tt.cursor, tt.count)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
				assert.Nil(t, films)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, films)
			}
		})
	}
}

func TestGetSimilarFilms(t *testing.T) {
	filmID := uuid.NewV4()

	tests := []struct {
		name       string
		filmID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilms  []models.MainPageFilm
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				similarFilmID1 := uuid.NewV4()
				similarFilmID2 := uuid.NewV4()

				rows := pgxpoolmock.NewRows([]string{
					"id", "cover", "title", "rating", "year", "genre_title",
				}).
					AddRow(
						similarFilmID1,
						"/static/cover1.jpg",
						"Similar Film 1",
						8.5,
						2023,
						"Drama",
					).
					AddRow(
						similarFilmID2,
						"/static/cover2.jpg",
						"Similar Film 2",
						7.8,
						2022,
						"Drama",
					).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetSimilarFilmsQuery, filmID).
					Return(rows, nil)
			},
			wantFilms: []models.MainPageFilm{
				{
					ID:     uuid.NewV4(),
					Cover:  "/static/cover1.jpg",
					Title:  "Similar Film 1",
					Rating: 8.5,
					Year:   2023,
					Genre:  "Drama",
				},
				{
					ID:     uuid.NewV4(),
					Cover:  "/static/cover2.jpg",
					Title:  "Similar Film 2",
					Rating: 7.8,
					Year:   2022,
					Genre:  "Drama",
				},
			},
			wantErr: false,
		},
		{
			name:   "QueryError",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetSimilarFilmsQuery, filmID).
					Return(nil, assert.AnError)
			},
			wantFilms: nil,
			wantErr:   true,
			errType:   films.ErrorInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			films, err := repo.GetSimilarFilms(testContext(), tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
				assert.Nil(t, films)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, films)
			}
		})
	}
}

func TestGetUpdates(t *testing.T) {
	userID := uuid.NewV4()
	offset := time.Now().Add(-7 * 24 * time.Hour)

	tests := []struct {
		name        string
		userID      uuid.UUID
		offset      time.Time
		repoMocker  func(*pgxpoolmock.MockPgxPool)
		wantNews    []models.News
		wantHasMore bool
	}{
		{
			name:   "SuccessWithNews",
			userID: userID,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				newsID1 := uuid.NewV4()
				newsID2 := uuid.NewV4()
				filmID := uuid.NewV4()

				rows := pgxpoolmock.NewRows([]string{
					"id", "title", "text", "film_id", "scheduled_at",
				}).
					AddRow(
						newsID1,
						"Новость 1",
						"Текст новости 1",
						filmID,
						offset.Add(24*time.Hour),
					).
					AddRow(
						newsID2,
						"Новость 2",
						"Текст новости 2",
						filmID,
						offset.Add(48*time.Hour),
					).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetUpdatesQuery, offset, userID).
					Return(rows, nil)
			},
			wantNews: []models.News{
				{
					ID:          uuid.NewV4(),
					Title:       "Новость 1",
					Text:        "Текст новости 1",
					FilmID:      uuid.NewV4(),
					ScheduledAt: offset.Add(24 * time.Hour),
				},
				{
					ID:          uuid.NewV4(),
					Title:       "Новость 2",
					Text:        "Текст новости 2",
					FilmID:      uuid.NewV4(),
					ScheduledAt: offset.Add(48 * time.Hour),
				},
			},
			wantHasMore: true,
		},
		{
			name:   "NoNews",
			userID: userID,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "title", "text", "film_id", "scheduled_at",
				}).ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetUpdatesQuery, offset, userID).
					Return(rows, nil)
			},
			wantNews:    []models.News{},
			wantHasMore: false,
		},
		{
			name:   "QueryErrorReturnsEmpty",
			userID: userID,
			offset: offset,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetUpdatesQuery, offset, userID).
					Return(nil, assert.AnError)
			},
			wantNews:    []models.News{},
			wantHasMore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			news, hasMore := repo.GetUpdates(testContext(), tt.userID, tt.offset)

			assert.Equal(t, tt.wantHasMore, hasMore)
			assert.Len(t, news, len(tt.wantNews))
		})
	}
}

func TestGetUsersRecommendations(t *testing.T) {
	userID := uuid.NewV4()

	tests := []struct {
		name       string
		userID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilms  []models.RecFilm
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmID1 := uuid.NewV4()
				filmID2 := uuid.NewV4()
				genreID := uuid.NewV4()
				countryID := uuid.NewV4()

				rows := pgxpoolmock.NewRows([]string{
					"id", "title", "year", "genre_id", "age_category",
					"country_id", "duration", "cluster_id", "rating",
					"amount_of_reviews", "user_rating", "cover",
				}).
					AddRow(
						filmID1,
						"Рекомендованный фильм 1",
						2023,
						genreID,
						"18+",
						countryID,
						120,
						1,
						8.5,
						150,
						9.0,
						"/static/cover1.jpg",
					).
					AddRow(
						filmID2,
						"Рекомендованный фильм 2",
						2022,
						genreID,
						"16+",
						countryID,
						110,
						2,
						7.8,
						100,
						8.0,
						"/static/cover2.jpg",
					).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetUsersRecommendationsQuery, userID).
					Return(rows, nil)
			},
			wantFilms: []models.RecFilm{
				{
					ID:              uuid.NewV4(),
					Title:           "Рекомендованный фильм 1",
					Year:            2023,
					GenreID:         uuid.NewV4(),
					AgeCategory:     "18+",
					CountryID:       uuid.NewV4(),
					Duration:        120,
					ClusterID:       1,
					Rating:          8.5,
					AmountOfReviews: 150,
					UserRating:      9.0,
					Cover:           "/static/cover1.jpg",
				},
				{
					ID:              uuid.NewV4(),
					Title:           "Рекомендованный фильм 2",
					Year:            2022,
					GenreID:         uuid.NewV4(),
					AgeCategory:     "16+",
					CountryID:       uuid.NewV4(),
					Duration:        110,
					ClusterID:       2,
					Rating:          7.8,
					AmountOfReviews: 100,
					UserRating:      8.0,
					Cover:           "/static/cover2.jpg",
				},
			},
			wantErr: false,
		},
		{
			name:   "QueryError",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				mockPool.EXPECT().
					Query(gomock.Any(), GetUsersRecommendationsQuery, userID).
					Return(nil, assert.AnError)
			},
			wantFilms: nil,
			wantErr:   true,
			errType:   films.ErrorInternalServerError,
		},
		{
			name:   "ScanErrorContinues",
			userID: userID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				rows := pgxpoolmock.NewRows([]string{
					"id", "title", "year", "genre_id", "age_category",
					"country_id", "duration", "cluster_id", "rating",
					"amount_of_reviews", "user_rating", "cover",
				}).
					AddRow(
						nil, // Это вызовет ошибку сканирования
						"Рекомендованный фильм 1",
						2023,
						uuid.NewV4(),
						"18+",
						uuid.NewV4(),
						120,
						1,
						8.5,
						150,
						9.0,
						"/static/cover1.jpg",
					).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetUsersRecommendationsQuery, userID).
					Return(rows, nil)
			},
			wantFilms: []models.RecFilm{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			films, err := repo.GetUsersRecommendations(testContext(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, films)
			}
		})
	}
}

func TestGetFilmPage(t *testing.T) {
	filmID := uuid.NewV4()
	genreID := uuid.NewV4()

	originalTitle := "Original Title"
	trailerURL := "https://youtube.com/trailer"
	slogan := "Great slogan"
	image1 := "/static/image1.jpg"
	image2 := "/static/image2.jpg"
	image3 := "/static/image3.jpg"
	filmURL := "/film/" + filmID.String()

	tests := []struct {
		name       string
		filmID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilm   models.FilmPage
		wantErr    bool
		errType    error
	}{
		{
			name:   "Success",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows([]string{
					"id", "title", "original_title", "cover", "poster",
					"short_description", "description", "age_category", "budget",
					"worldwide_fees", "trailer_url", "year",
					"slogan", "duration", "image1", "image2", "image3",
					"genre_title", "genre_id", "country_title", "film_url", "number_of_ratings", "is_out",
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
						&slogan,
						120,
						&image1,
						&image2,
						&image3,
						"Drama",
						genreID,
						"USA",
						filmURL,
						150,
						true,
					).
					ToPgxRows()
				filmRows.Next()

				ratingRows := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(8.5).
					ToPgxRows()
				ratingRows.Next()

				actorRows := pgxpoolmock.NewRows([]string{
					"id", "russian_name", "original_name", "photo", "height",
					"birth_date", "death_date", "zodiac_sign", "birth_place", "marital_status",
				}).
					AddRow(
						uuid.NewV4(),
						"Актер 1",
						"Actor 1",
						"/static/actor1.jpg",
						180,
						time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
						nil,
						"Козерог",
						"Москва",
						"Холост",
					).
					ToPgxRows()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmPageQuery, filmID).
					Return(filmRows)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID).
					Return(ratingRows)

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmActorsQuery, filmID).
					Return(actorRows, nil)
			},
			wantFilm: models.FilmPage{
				ID:               filmID,
				Title:            "Test Film",
				OriginalTitle:    &originalTitle,
				Cover:            "/static/cover.jpg",
				Poster:           "/static/poster.jpg",
				ShortDescription: "Short description",
				Description:      "Full description",
				AgeCategory:      "18+",
				Budget:           1000000,
				WorldwideFees:    5000000,
				TrailerURL:       &trailerURL,
				Year:             2023,
				Slogan:           &slogan,
				Duration:         120,
				Image1:           &image1,
				Image2:           &image2,
				Image3:           &image3,
				Genre:            "Drama",
				GenreID:          genreID,
				Country:          "USA",
				FilmURL:          filmURL,
				NumberOfRatings:  150,
				IsOut:            true,
				Rating:           8.5,
				Actors: []models.Actor{
					{
						ID:            uuid.NewV4(),
						RussianName:   "Актер 1",
						Photo:         "/static/actor1.jpg",
						Height:        180,
						BirthDate:     time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
						DeathDate:     nil,
						ZodiacSign:    "Козерог",
						BirthPlace:    "Москва",
						MaritalStatus: "Холост",
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "ActorsQueryErrorReturnsError",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows([]string{
					"id", "title", "original_title", "cover", "poster",
					"short_description", "description", "age_category", "budget",
					"worldwide_fees", "trailer_url", "year",
					"slogan", "duration", "image1", "image2", "image3",
					"genre_title", "genre_id", "country_title", "film_url", "number_of_ratings", "is_out",
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
						&slogan,
						120,
						&image1,
						&image2,
						&image3,
						"Drama",
						genreID,
						"USA",
						filmURL,
						150,
						true,
					).
					ToPgxRows()
				filmRows.Next()

				ratingRows := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(8.5).
					ToPgxRows()
				ratingRows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmPageQuery, filmID).
					Return(filmRows)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID).
					Return(ratingRows)

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmActorsQuery, filmID).
					Return(nil, assert.AnError)
			},
			wantFilm: models.FilmPage{},
			wantErr:  true,
			errType:  films.ErrorInternalServerError,
		},
		{
			name:   "ActorsNotFoundError",
			filmID: filmID,
			repoMocker: func(mockPool *pgxpoolmock.MockPgxPool) {
				filmRows := pgxpoolmock.NewRows([]string{
					"id", "title", "original_title", "cover", "poster",
					"short_description", "description", "age_category", "budget",
					"worldwide_fees", "trailer_url", "year",
					"slogan", "duration", "image1", "image2", "image3",
					"genre_title", "genre_id", "country_title", "film_url", "number_of_ratings", "is_out",
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
						&slogan,
						120,
						&image1,
						&image2,
						&image3,
						"Drama",
						genreID,
						"USA",
						filmURL,
						150,
						true,
					).
					ToPgxRows()
				filmRows.Next()

				ratingRows := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(8.5).
					ToPgxRows()
				ratingRows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmPageQuery, filmID).
					Return(filmRows)

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID).
					Return(ratingRows)

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmActorsQuery, filmID).
					Return(nil, pgx.ErrNoRows)
			},
			wantFilm: models.FilmPage{
				ID:               filmID,
				Title:            "Test Film",
				OriginalTitle:    &originalTitle,
				Cover:            "/static/cover.jpg",
				Poster:           "/static/poster.jpg",
				ShortDescription: "Short description",
				Description:      "Full description",
				AgeCategory:      "18+",
				Budget:           1000000,
				WorldwideFees:    5000000,
				TrailerURL:       &trailerURL,
				Year:             2023,
				Slogan:           &slogan,
				Duration:         120,
				Image1:           &image1,
				Image2:           &image2,
				Image3:           &image3,
				Genre:            "Drama",
				GenreID:          genreID,
				Country:          "USA",
				FilmURL:          filmURL,
				NumberOfRatings:  150,
				IsOut:            true,
				Rating:           8.5,
				Actors:           []models.Actor{},
			},
			wantErr: true, // Но возвращает результат с ошибкой
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
			tt.repoMocker(mockPool)

			repo := NewFilmRepository(mockPool)
			film, err := repo.GetFilmPage(testContext(), tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.Equal(t, tt.errType, err)
				}
				if tt.name == "ActorsNotFoundError" {
					// Для этого случая метод возвращает и результат, и ошибку
					assert.NotEqual(t, models.FilmPage{}, film)
					assert.Equal(t, filmID, film.ID)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFilm.ID, film.ID)
				assert.Equal(t, tt.wantFilm.Title, film.Title)
				assert.Equal(t, tt.wantFilm.Rating, film.Rating)
			}
		})
	}
}

func TestFilmRepositoryConstructor(t *testing.T) {
	t.Run("NewFilmRepository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
		repo := NewFilmRepository(mockPool)

		assert.NotNil(t, repo)
		assert.Equal(t, mockPool, repo.db)
	})
}
