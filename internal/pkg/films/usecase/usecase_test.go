package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
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
		setupMock   func()
		filmID      uuid.UUID
		userID      uuid.UUID
		expected    models.FilmPage
		expectError bool
	}{
		{
			name: "Success - with user feedback (has title)",
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
				mockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
			},
			filmID: filmID,
			userID: userID,
			expected: models.FilmPage{
				ID:          filmID,
				Title:       "Test Film",
				Rating:      8.5,
				Description: "Test description",
				Year:        2024,
				IsReviewed:  true,
				UserRating:  &userRating,
				IsLiked:     false,
			},
			expectError: false,
		},
		{
			name: "Success - with user like",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(expectedFilm, nil)
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, nil)
			},
			filmID: filmID,
			userID: userID,
			expected: models.FilmPage{
				ID:          filmID,
				Title:       "Test Film",
				Rating:      8.5,
				Description: "Test description",
				Year:        2024,
				IsReviewed:  false,
				UserRating:  nil,
				IsLiked:     true,
			},
			expectError: false,
		},
		{
			name: "Success - without user feedback",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(expectedFilm, nil)
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
			},
			filmID:      filmID,
			userID:      userID,
			expected:    expectedFilm,
			expectError: false,
		},
		{
			name: "Error - film not found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmPage(gomock.Any(), filmID).
					Return(models.FilmPage{}, films.ErrorNotFound)
			},
			filmID:      filmID,
			userID:      userID,
			expected:    models.FilmPage{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilm(testContext(), tt.filmID, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Title, result.Title)
				assert.Equal(t, tt.expected.IsReviewed, result.IsReviewed)
				assert.Equal(t, tt.expected.IsLiked, result.IsLiked)
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
		setupMock   func()
		filmID      uuid.UUID
		userID      uuid.UUID
		pager       models.Pager
		expected    []models.FilmFeedback
		expectError bool
	}{
		{
			name: "Success - with user feedback",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(userFeedback, nil)
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return([]models.FilmFeedback{userFeedback, otherFeedback}, nil)
			},
			filmID:      filmID,
			userID:      userID,
			pager:       pager,
			expected:    []models.FilmFeedback{userFeedback, otherFeedback},
			expectError: false,
		},
		{
			name: "Success - without user feedback",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return([]models.FilmFeedback{otherFeedback}, nil)
			},
			filmID:      filmID,
			userID:      userID,
			pager:       pager,
			expected:    []models.FilmFeedback{otherFeedback},
			expectError: false,
		},
		{
			name: "Success - with nil userID",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return([]models.FilmFeedback{otherFeedback}, nil)
			},
			filmID:      filmID,
			userID:      uuid.Nil, // nil user
			pager:       pager,
			expected:    []models.FilmFeedback{otherFeedback},
			expectError: false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return(nil, films.ErrorInternalServerError)
			},
			filmID:      filmID,
			userID:      userID,
			pager:       pager,
			expected:    []models.FilmFeedback{},
			expectError: true,
		},
		{
			name: "Error - no feedbacks",
			setupMock: func() {
				mockRepo.EXPECT().
					CheckUserFeedbackExists(gomock.Any(), userID, filmID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
				mockRepo.EXPECT().
					GetFilmFeedbacks(gomock.Any(), filmID, pager.Count, pager.Offset).
					Return([]models.FilmFeedback{}, nil)
			},
			filmID:      filmID,
			userID:      userID,
			pager:       pager,
			expected:    []models.FilmFeedback{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilmFeedbacks(testContext(), tt.filmID, tt.userID, tt.pager)

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

	validText := "This is a valid feedback text with more than 30 characters"
	validInput := models.FilmFeedbackInput{
		Title:  "Great film!",
		Text:   validText,
		Rating: 9,
	}

	tests := []struct {
		name        string
		setupMock   func()
		req         models.FilmFeedbackInput
		filmID      uuid.UUID
		userID      uuid.UUID
		expectError bool
	}{
		{
			name: "Success - create new feedback",
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
			userID:      userID,
			expectError: false,
		},
		{
			name: "Success - update existing feedback",
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
			userID:      userID,
			expectError: false,
		},
		{
			name:        "Error - invalid rating too low",
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: validText, Rating: 0},
			filmID:      filmID,
			userID:      userID,
			expectError: true,
		},
		{
			name:        "Error - invalid rating too high",
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: validText, Rating: 11},
			filmID:      filmID,
			userID:      userID,
			expectError: true,
		},
		{
			name:        "Error - title too short",
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "", Text: validText, Rating: 5},
			filmID:      filmID,
			userID:      userID,
			expectError: true,
		},
		{
			name:        "Error - title too long",
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: string(make([]byte, 101)), Text: validText, Rating: 5},
			filmID:      filmID,
			userID:      userID,
			expectError: true,
		},
		{
			name:        "Error - text too short",
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: "Short", Rating: 5},
			filmID:      filmID,
			userID:      userID,
			expectError: true,
		},
		{
			name:        "Error - text too long",
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Title: "Test", Text: string(make([]byte, 1001)), Rating: 5},
			filmID:      filmID,
			userID:      userID,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.SendFeedback(testContext(), tt.req, tt.filmID, tt.userID)

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

	validInput := models.FilmFeedbackInput{
		Rating: 8,
	}

	tests := []struct {
		name        string
		setupMock   func()
		req         models.FilmFeedbackInput
		filmID      uuid.UUID
		userID      uuid.UUID
		expectError bool
	}{
		{
			name: "Success - create new rating",
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
			userID:      userID,
			expectError: false,
		},
		{
			name: "Success - update existing rating",
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
			userID:      userID,
			expectError: false,
		},
		{
			name:        "Error - invalid rating",
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Rating: 0},
			filmID:      filmID,
			userID:      userID,
			expectError: true,
		},
		{
			name:        "Error - invalid rating too high",
			setupMock:   func() {},
			req:         models.FilmFeedbackInput{Rating: 11},
			filmID:      filmID,
			userID:      userID,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.SetRating(testContext(), tt.req, tt.filmID, tt.userID)

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

func TestFilmUsecase_GetFilmsForCalendar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFilmRepo(ctrl)
	usecase := NewFilmUsecase(mockRepo)

	userID := uuid.NewV4()
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedFilms := []models.FilmInCalendar{
		{
			ID:          uuid.NewV4(),
			Title:       "Film 1",
			ReleaseDate: time.Now(),
		},
		{
			ID:          uuid.NewV4(),
			Title:       "Film 2",
			ReleaseDate: time.Now().Add(24 * time.Hour),
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		pager       models.Pager
		userID      uuid.UUID
		expected    []models.FilmInCalendar
		expectError bool
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsForCalendar(gomock.Any(), pager.Count, pager.Offset).
					Return(expectedFilms, nil)
				mockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, gomock.Any()).
					Return(models.FilmFeedback{}, films.ErrorNotFound).AnyTimes()
			},
			pager:       pager,
			userID:      userID,
			expected:    expectedFilms,
			expectError: false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsForCalendar(gomock.Any(), pager.Count, pager.Offset).
					Return(nil, films.ErrorInternalServerError)
			},
			pager:       pager,
			userID:      userID,
			expected:    []models.FilmInCalendar{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilmsForCalendar(testContext(), tt.pager, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expected), len(result))
			}
		})
	}
}
