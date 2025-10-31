package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/films/mocks"
	"kinopoisk/internal/pkg/middleware/logger"

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

func testContextWithUser(user models.User) context.Context {
	ctx := testContext()
	return context.WithValue(ctx, auth.UserKey, user)
}

func TestFilmUsecase_GetPromoFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFilmRepo(ctrl)
	usecase := NewFilmUsecase(mockRepo)

	promoFilmID := uuid.FromStringOrNil("8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c")
	expectedPromoFilm := models.PromoFilm{
		ID:               promoFilmID,
		Image:            "promo.jpg",
		Title:            "Test Promo Film",
		Rating:           8.5,
		ShortDescription: "Test description",
		Year:             2024,
		Genre:            "Action",
		Duration:         120,
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    models.PromoFilm
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetPromoFilmByID(gomock.Any(), promoFilmID).
					Return(models.PromoFilm{
						ID:               promoFilmID,
						Image:            "promo.jpg",
						Title:            "Test Promo Film",
						ShortDescription: "Test description",
						Year:             2024,
						Genre:            "Action",
						Duration:         120,
					}, nil)
				mockRepo.EXPECT().
					GetFilmAvgRating(gomock.Any(), promoFilmID).
					Return(8.5, nil)
			},
			expected:    expectedPromoFilm,
			expectError: false,
		},
		{
			name: "Error - promo film not found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetPromoFilmByID(gomock.Any(), promoFilmID).
					Return(models.PromoFilm{}, errors.New("not found"))
			},
			expected:    models.PromoFilm{},
			expectError: true,
			errorMsg:    "no films",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetPromoFilm(testContext())

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
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
		errorMsg    string
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
					Return(nil, errors.New("database error"))
			},
			expected:    []models.MainPageFilm{},
			expectError: true,
			errorMsg:    "no films",
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
			errorMsg:    "no films",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilms(testContext(), pager)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
				assert.Equal(t, tt.expected, result)
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
	expectedFilm := models.FilmPage{
		ID:          filmID,
		Title:       "Test Film",
		Rating:      8.5,
		Description: "Test description",
		Year:        2024,
	}

	tests := []struct {
		name        string
		setupMock   func()
		filmID      uuid.UUID
		expected    models.FilmPage
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(expectedFilm, nil)
			},
			filmID:      filmID,
			expected:    expectedFilm,
			expectError: false,
		},
		{
			name: "Error - film not found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(models.FilmPage{}, errors.New("not found"))
			},
			filmID:      filmID,
			expected:    models.FilmPage{},
			expectError: true,
			errorMsg:    "no films",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilm(testContext(), tt.filmID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
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
	text1 := "Amazing acting and story"
	title2 := "Good film"
	text2 := "Enjoyed watching it"

	expectedFeedbacks := []models.FilmFeedback{
		{
			ID:        uuid.NewV4(),
			UserID:    userID,
			FilmID:    filmID,
			Title:     &title1,
			Text:      &text1,
			Rating:    9,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserLogin: "user1",
			IsMine:    true,
		},
		{
			ID:        uuid.NewV4(),
			UserID:    uuid.NewV4(),
			FilmID:    filmID,
			Title:     &title2,
			Text:      &text2,
			Rating:    8,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserLogin: "user2",
			IsMine:    false,
		},
	}

	tests := []struct {
		name        string
		ctx         context.Context
		setupMock   func()
		expected    []models.FilmFeedback
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success - with user in context",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return(expectedFeedbacks, nil)
			},
			expected:    expectedFeedbacks,
			expectError: false,
		},
		{
			name: "Error - repository error",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return(nil, errors.New("database error"))
			},
			expected:    []models.FilmFeedback{},
			expectError: true,
			errorMsg:    "no feedbacks",
		},
		{
			name: "Error - no feedbacks",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return([]models.FilmFeedback{}, nil)
			},
			expected:    []models.FilmFeedback{},
			expectError: true,
			errorMsg:    "no feedbacks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilmFeedbacks(tt.ctx, filmID, pager)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
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

	validInput := models.FilmFeedbackInput{
		Title:  "Great film!",
		Text:   "Amazing acting and story",
		Rating: 9,
	}

	tests := []struct {
		name        string
		ctx         context.Context
		setupMock   func()
		req         models.FilmFeedbackInput
		filmID      uuid.UUID
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success - create new feedback",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, errors.New("not found"))
				mockRepo.EXPECT().
					CreateFeedback(gomock.Any(), gomock.Any()).
					Return(nil)
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
				oldText := "Old text"
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
			errorMsg:    "user not authenticated",
		},
		{
			name:        "Error - invalid rating too low",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: "Test", Rating: 0},
			filmID:      filmID,
			expectError: true,
			errorMsg:    "rating must be between 1 and 10",
		},
		{
			name:        "Error - invalid rating too high",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: "Test", Rating: 11},
			filmID:      filmID,
			expectError: true,
			errorMsg:    "rating must be between 1 and 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.SendFeedback(tt.ctx, tt.req, tt.filmID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
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
		errorMsg    string
	}{
		{
			name: "Success",
			ctx:  testContextWithUser(user),
			setupMock: func() {
				mockRepo.EXPECT().
					CreateFeedback(gomock.Any(), gomock.Any()).
					Return(nil)
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
			errorMsg:    "user is not authorized",
		},
		{
			name:        "Error - invalid rating",
			ctx:         testContextWithUser(user),
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Rating: 0},
			filmID:      filmID,
			expectError: true,
			errorMsg:    "rating must be between 1 and 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.SetRating(tt.ctx, tt.req, tt.filmID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.Equal(t, tt.req.Rating, result.Rating)
			}
		})
	}
}
