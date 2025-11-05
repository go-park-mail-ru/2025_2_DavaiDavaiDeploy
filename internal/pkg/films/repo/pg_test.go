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

func TestGetFilmPage(t *testing.T) {
	filmID := uuid.NewV4()
	genre := "Drama"
	country := "USA"
	numberOfRatings := 100

	tests := []struct {
		name       string
		filmID     uuid.UUID
		repoMocker func(*pgxpoolmock.MockPgxPool)
		wantFilm   models.FilmPage
		wantErr    bool
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
					"genre", "country", "number_of_ratings",
				}).
					AddRow(
						filmID,
						"Test Film",
						nil,
						"/static/cover.jpg",
						"/static/poster.jpg",
						"Short description",
						"Full description",
						"18+",
						1000000,
						5000000,
						nil,
						2023,
						nil,
						120,
						nil,
						nil,
						nil,
						genre,
						country,
						numberOfRatings,
					).
					ToPgxRows()
				filmRows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmPageQuery, filmID).
					Return(filmRows)

				ratingRows := pgxpoolmock.NewRows([]string{"coalesce"}).
					AddRow(8.5).
					ToPgxRows()
				ratingRows.Next()

				mockPool.EXPECT().
					QueryRow(gomock.Any(), GetFilmAvgRatingQuery, filmID).
					Return(ratingRows)

				actorRows := pgxpoolmock.NewRows([]string{
					"id", "russian_name", "original_name", "photo", "height",
					"birth_date", "death_date", "zodiac_sign", "birth_place", "marital_status",
				}).
					AddRow(
						uuid.NewV4(),
						"Actor Name",
						"Actor Original Name",
						"/static/photo.jpg",
						180,
						time.Now(),
						nil,
						"Leo",
						"Moscow",
						"Single",
					).
					ToPgxRows()

				mockPool.EXPECT().
					Query(gomock.Any(), GetFilmActorsQuery, filmID).
					Return(actorRows, nil)
			},
			wantFilm: models.FilmPage{
				ID:               filmID,
				Title:            "Test Film",
				Cover:            "/static/cover.jpg",
				Poster:           "/static/poster.jpg",
				ShortDescription: "Short description",
				Description:      "Full description",
				AgeCategory:      "18+",
				Budget:           1000000,
				WorldwideFees:    5000000,
				Year:             2023,
				Duration:         120,
				Genre:            genre,
				Country:          country,
				NumberOfRatings:  numberOfRatings,
				Rating:           8.5,
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
			film, err := repo.GetFilmPage(testContext(), tt.filmID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFilm.ID, film.ID)
				assert.Equal(t, tt.wantFilm.Title, film.Title)
				assert.Equal(t, tt.wantFilm.Genre, film.Genre)
				assert.Equal(t, tt.wantFilm.Country, film.Country)
				assert.Equal(t, tt.wantFilm.Rating, film.Rating)
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
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Login, user.Login)
				assert.Equal(t, tt.wantUser.Avatar, user.Avatar)
			}
		})
	}
}
