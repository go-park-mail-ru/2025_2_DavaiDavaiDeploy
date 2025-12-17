package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"
	"kinopoisk/internal/pkg/search/mocks"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func testSearchLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func testSearchContext() context.Context {
	testLogger := testSearchLogger()
	return context.WithValue(context.Background(), logger.LoggerKey, testLogger)
}

func TestSearchUsecase_GetFilmsFromSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSearchRepo(ctrl)
	usecase := NewSearchUsecase(mockRepo)

	searchString := "пираты"
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedFilms := []models.MainPageFilm{
		{
			ID:     uuid.NewV4(),
			Cover:  "pirates.jpg",
			Title:  "Пираты Карибского моря",
			Rating: 8.5,
			Year:   2003,
			Genre:  "Приключения",
		},
		{
			ID:     uuid.NewV4(),
			Cover:  "pirates2.jpg",
			Title:  "Пираты Карибского моря 2",
			Rating: 7.8,
			Year:   2006,
			Genre:  "Приключения",
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []models.MainPageFilm
		expectError bool
		errorType   error
	}{
		{
			name: "Success - with films",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsFromSearch(gomock.Any(), searchString, pager.Count, pager.Offset).
					Return(expectedFilms, nil)
			},
			expected:    expectedFilms,
			expectError: false,
		},
		{
			name: "Success - no films found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsFromSearch(gomock.Any(), searchString, pager.Count, pager.Offset).
					Return([]models.MainPageFilm{}, nil)
			},
			expected:    []models.MainPageFilm{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilmsFromSearch(testSearchContext(), searchString, pager)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestSearchUsecase_GetActorsFromSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSearchRepo(ctrl)
	usecase := NewSearchUsecase(mockRepo)

	searchString := "джонни"
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedActors := []models.MainPageActor{
		{
			ID:          uuid.NewV4(),
			RussianName: "Джонни Депп",
			Photo:       "depp.jpg",
		},
		{
			ID:          uuid.NewV4(),
			RussianName: "Джонни Ноксвилл",
			Photo:       "knoxville.jpg",
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []models.MainPageActor
		expectError bool
		errorType   error
	}{
		{
			name: "Success - with actors",
			setupMock: func() {
				mockRepo.EXPECT().
					GetActorsFromSearch(gomock.Any(), searchString, pager.Count, pager.Offset).
					Return(expectedActors, nil)
			},
			expected:    expectedActors,
			expectError: false,
		},
		{
			name: "Success - no actors found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetActorsFromSearch(gomock.Any(), searchString, pager.Count, pager.Offset).
					Return([]models.MainPageActor{}, nil)
			},
			expected:    []models.MainPageActor{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetActorsFromSearch(testSearchContext(), searchString, pager)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
