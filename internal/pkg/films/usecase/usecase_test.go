package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/films/mocks"
	"kinopoisk/internal/pkg/middleware/logger"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func testContext() context.Context {
	testLogger := testLogger()
	return context.WithValue(context.Background(), logger.LoggerKey, testLogger)
}

func testContextWithUser(user models.User) context.Context {
	ctx := testContext()
	return context.WithValue(ctx, auth.UserKey, user)
}

func TestFilmUsecase_GetPromoFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFilmRepo(ctrl)
	usecase := NewFilmUsecase(mockRepo)

	promoFilmIDs := []string{
		"8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c",
		"2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	}

	tests := []struct {
		name        string
		setupMock   func()
		expectError bool
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetPromoFilmByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id uuid.UUID) (models.PromoFilm, error) {
						return models.PromoFilm{
							ID:               id,
							Image:            "promo.jpg",
							Title:            "Test Promo Film",
							ShortDescription: "Test description",
							Year:             2024,
							Genre:            "Action",
							Duration:         120,
						}, nil
					}).AnyTimes()
				mockRepo.EXPECT().
					GetFilmAvgRating(gomock.Any(), gomock.Any()).
					Return(8.5, nil).AnyTimes()
			},
			expectError: false,
		},
		{
			name: "Success - rating error returns zero rating",
			setupMock: func() {
				mockRepo.EXPECT().
					GetPromoFilmByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id uuid.UUID) (models.PromoFilm, error) {
						return models.PromoFilm{
							ID:               id,
							Image:            "promo.jpg",
							Title:            "Test Promo Film",
							ShortDescription: "Test description",
							Year:             2024,
							Genre:            "Action",
							Duration:         120,
						}, nil
					}).AnyTimes()
				mockRepo.EXPECT().
					GetFilmAvgRating(gomock.Any(), gomock.Any()).
					Return(0.0, films.ErrorNotFound).AnyTimes()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetPromoFilm(testContext())

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.Contains(t, promoFilmIDs, result.ID.String())
				assert.Equal(t, "Test Promo Film", result.Title)
			}
		})
	}
}

func TestFilmUsecase_GetFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFilmRepo(ctrl)
	usecase := NewFilmUsecase(mockRepo)

	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedFilms := []models.MainPageFilm{
		{
			ID:     uuid.NewV4(),
			Cover:  "film1.jpg",
			Title:  "Film 1",
			Rating: 8.5,
			Year:   2024,
			Genre:  "Action",
		},
		{
			ID:     uuid.NewV4(),
			Cover:  "film2.jpg",
			Title:  "Film 2",
			Rating: 7.9,
			Year:   2023,
			Genre:  "Drama",
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []models.MainPageFilm
		expectError bool
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return(expectedFilms, nil)
			},
			expected:    expectedFilms,
			expectError: false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return(nil, films.ErrorInternalServerError)
			},
			expected:    []models.MainPageFilm{},
			expectError: true,
		},
		{
			name: "Error - no films",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return([]models.MainPageFilm{}, nil)
			},
			expected:    []models.MainPageFilm{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilms(testContext(), pager)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestFilmUsecase_GetFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFilmRepo(ctrl)
	usecase := NewFilmUsecase(mockRepo)

	filmID := uuid.NewV4()
	userID := uuid.NewV4()
	user := models.User{ID: userID}
	expectedFilm := models.FilmPage{
		ID:          filmID,
		Title:       "Test Film",
		Rating:      8.5,
		Description: "Test description",
		Year:        2024,
	}

	title := "User review"
	userRating := 9

	tests := []struct {
		name        string
		ctx         context.Context
		setupMock   func()
		filmID      uuid.UUID
		expected    models.FilmPage
		expectError bool
	}{
		{
			name: "Success - with user feedback (has title)",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(expectedFilm, nil)
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{
						Title:  &title,
						Rating: userRating,
					}, nil)
			},
			filmID: filmID,
			expected: models.FilmPage{
				ID:          filmID,
				Title:       "Test Film",
				Rating:      8.5,
				Description: "Test description",
				Year:        2024,
				IsReviewed:  true,
				UserRating:  &userRating,
			},
			expectError: false,
		},
		{
			name: "Success - without user feedback",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(expectedFilm, nil)
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
			},
			filmID:      filmID,
			expected:    expectedFilm,
			expectError: false,
		},
		{
			name: "Error - film not found",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(models.FilmPage{}, films.ErrorNotFound)
			},
			filmID:      filmID,
			expected:    models.FilmPage{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilm(tt.ctx, tt.filmID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Title, result.Title)
				assert.Equal(t, tt.expected.IsReviewed, result.IsReviewed)
				if tt.expected.UserRating != nil {
					assert.Equal(t, *tt.expected.UserRating, *result.UserRating)
				}
			}
		})
	}
}

func TestFilmUsecase_GetFilmFeedbacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFilmRepo(ctrl)
	usecase := NewFilmUsecase(mockRepo)

	filmID := uuid.NewV4()
	userID := uuid.NewV4()
	user := models.User{ID: userID}
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	title1 := "Great film!"
	text1 := "Amazing acting and story with more than 30 characters"
	title2 := "Good film"
	text2 := "Enjoyed watching it with more than 30 characters"

	userFeedback := models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    userID,
		FilmID:    filmID,
		Title:     &title1,
		Text:      &text1,
		Rating:    9,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserLogin: "user1",
	}

	otherFeedback := models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    uuid.NewV4(),
		FilmID:    filmID,
		Title:     &title2,
		Text:      &text2,
		Rating:    8,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserLogin: "user2",
	}

	tests := []struct {
		name        string
		ctx         context.Context
		setupMock   func()
		expected    []models.FilmFeedback
		expectError bool
	}{
		{
			name: "Success - with user feedback",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(userFeedback, nil)
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return([]models.FilmFeedback{userFeedback, otherFeedback}, nil)
			},
			expected:    []models.FilmFeedback{userFeedback, otherFeedback},
			expectError: false,
		},
		{
			name: "Success - without user feedback",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return([]models.FilmFeedback{otherFeedback}, nil)
			},
			expected:    []models.FilmFeedback{otherFeedback},
			expectError: false,
		},
		{
			name: "Error - repository error",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return(nil, films.ErrorInternalServerError)
			},
			expected:    []models.FilmFeedback{},
			expectError: true,
		},
		{
			name: "Error - no feedbacks",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return([]models.FilmFeedback{}, nil)
			},
			expected:    []models.FilmFeedback{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilmFeedbacks(tt.ctx, filmID, pager)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(result))
			}
		})
	}
}

func TestFilmUsecase_SendFeedback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFilmRepo(ctrl)
	usecase := NewFilmUsecase(mockRepo)

	userID := uuid.NewV4()
	filmID := uuid.NewV4()
	user := models.User{ID: userID}

	validText := "This is a valid feedback text with more than 30 characters"
	validInput := models.FilmFeedbackInput{
		Title:  "Great film!",
		Text:   validText,
		Rating: 9,
	}

	tests := []struct {
		name        string
		ctx         context.Context
		setupMock   func()
		req         models.FilmFeedbackInput
		filmID      uuid.UUID
		expectError bool
	}{
		{
			name: "Success - create new feedback",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					CreateFeedback(gomock.Any(), gomock.Any()).
					Return(nil)
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(models.FilmPage{Rating: 8.7}, nil)
			},
			req:         validInput,
			filmID:      filmID,
			expectError: false,
		},
		{
			name: "Success - update existing feedback",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				oldTitle := "Old title"
				oldText := "Old text that was previously written by user"
				existingFeedback := models.FilmFeedback{
					ID:     uuid.NewV4(),
					UserID: userID,
					FilmID: filmID,
					Title:  &oldTitle,
					Text:   &oldText,
					Rating: 7,
				}
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(existingFeedback, nil)
				mockRepo.EXPECT().
					UpdateFeedback(gomock.Any(), gomock.Any()).
					Return(nil)
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(models.FilmPage{Rating: 8.2}, nil)
			},
			req:         validInput,
			filmID:      filmID,
			expectError: false,
		},
		{
			name:        "Error - no user in context",
			ctx:         testContext(),
			setupMock:   func() {},
			req:         validInput,
			filmID:      filmID,
			expectError: true,
		},
		{
			name:        "Error - invalid rating too low",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: validText, Rating: 0},
			filmID:      filmID,
			expectError: true,
		},
		{
			name:        "Error - invalid rating too high",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: validText, Rating: 11},
			filmID:      filmID,
			expectError: true,
		},
		{
			name:        "Error - title too short",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "", Text: validText, Rating: 5},
			filmID:      filmID,
			expectError: true,
		},
		{
			name:        "Error - title too long",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: string(make([]byte, 101)), Text: validText, Rating: 5},
			filmID:      filmID,
			expectError: true,
		},
		{
			name:        "Error - text too short",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: "Short", Rating: 5},
			filmID:      filmID,
			expectError: true,
		},
		{
			name:        "Error - text too long",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: string(make([]byte, 1001)), Rating: 5},
			filmID:      filmID,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.SendFeedback(tt.ctx, tt.req, tt.filmID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.Equal(t, tt.req.Rating, result.Rating)
			}
		})
	}
}

func TestFilmUsecase_SetRating(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFilmRepo(ctrl)
	usecase := NewFilmUsecase(mockRepo)

	userID := uuid.NewV4()
	filmID := uuid.NewV4()
	user := models.User{ID: userID}

	validInput := models.FilmFeedbackInput{
		Rating: 8,
	}

	tests := []struct {
		name        string
		ctx         context.Context
		setupMock   func()
		req         models.FilmFeedbackInput
		filmID      uuid.UUID
		expectError bool
	}{
		{
			name: "Success - create new rating",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					CreateFeedback(gomock.Any(), gomock.Any()).
					Return(nil)
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(models.FilmPage{Rating: 8.1}, nil)
			},
			req:         validInput,
			filmID:      filmID,
			expectError: false,
		},
		{
			name: "Success - update existing rating",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				existingFeedback := models.FilmFeedback{
					ID:     uuid.NewV4(),
					UserID: userID,
					FilmID: filmID,
					Rating: 7,
				}
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(existingFeedback, nil)
				mockRepo.EXPECT().
					UpdateFeedback(gomock.Any(), gomock.Any()).
					Return(nil)
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(models.FilmPage{Rating: 7.8}, nil)
			},
			req:         validInput,
			filmID:      filmID,
			expectError: false,
		},
		{
			name:        "Error - no user in context",
			ctx:         testContext(),
			setupMock:   func() {},
			req:         validInput,
			filmID:      filmID,
			expectError: true,
		},
		{
			name:        "Error - invalid rating",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Rating: 0},
			filmID:      filmID,
			expectError: true,
		},
		{
			name:        "Error - invalid rating too high",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Rating: 11},
			filmID:      filmID,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.SetRating(tt.ctx, tt.req, tt.filmID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.Equal(t, tt.req.Rating, result.Rating)
			}
		})
	}
}
