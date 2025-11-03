package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/genres/mocks"
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

func TestGenreUsecase_GetGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGenreRepo(ctrl)
	usecase := NewGenreUsecase(mockRepo)

	genreID := uuid.NewV4()
	expectedGenre := models.Genre{
		ID:          genreID,
		Title:       "Action",
		Description: "Action films",
		Icon:        "action.png",
	}

	tests := []struct {
		name        string
		setupMock   func()
		genreID     uuid.UUID
		expected    models.Genre
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetGenreByID(gomock.Any(), genreID).
					Return(expectedGenre, nil)
			},
			genreID:     genreID,
			expected:    expectedGenre,
			expectError: false,
		},
		{
			name: "Error - genre not found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetGenreByID(gomock.Any(), genreID).
					Return(models.Genre{}, errors.New("not found"))
			},
			genreID:     genreID,
			expected:    models.Genre{},
			expectError: true,
			errorMsg:    "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetGenre(testContext(), tt.genreID)

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

func TestGenreUsecase_GetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGenreRepo(ctrl)
	usecase := NewGenreUsecase(mockRepo)

	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedGenres := []models.Genre{
		{
			ID:          uuid.NewV4(),
			Title:       "Action",
			Description: "Action films",
			Icon:        "action.png",
		},
		{
			ID:          uuid.NewV4(),
			Title:       "Drama",
			Description: "Drama films",
			Icon:        "drama.png",
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []models.Genre
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetGenresWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return(expectedGenres, nil)
			},
			expected:    expectedGenres,
			expectError: false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				mockRepo.EXPECT().
					GetGenresWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return(nil, errors.New("database error"))
			},
			expected:    []models.Genre{},
			expectError: true,
			errorMsg:    "database error",
		},
		{
			name: "Error - no genres",
			setupMock: func() {
				mockRepo.EXPECT().
					GetGenresWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return([]models.Genre{}, nil)
			},
			expected:    []models.Genre{},
			expectError: true,
			errorMsg:    "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetGenres(testContext(), pager)

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

func TestGenreUsecase_GetFilmsByGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGenreRepo(ctrl)
	usecase := NewGenreUsecase(mockRepo)

	genreID := uuid.NewV4()
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedFilms := []models.MainPageFilm{
		{
			ID:     uuid.NewV4(),
			Cover:  "film1.jpg",
			Title:  "Action Film 1",
			Rating: 8.5,
			Year:   2024,
			Genre:  "Action",
		},
		{
			ID:     uuid.NewV4(),
			Cover:  "film2.jpg",
			Title:  "Action Film 2",
			Rating: 7.9,
			Year:   2023,
			Genre:  "Action",
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		genreID     uuid.UUID
		expected    []models.MainPageFilm
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsByGenre(gomock.Any(), genreID, pager.Count, pager.Offset).
					Return(expectedFilms, nil)
			},
			genreID:     genreID,
			expected:    expectedFilms,
			expectError: false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsByGenre(gomock.Any(), genreID, pager.Count, pager.Offset).
					Return(nil, errors.New("database error"))
			},
			genreID:     genreID,
			expected:    []models.MainPageFilm{},
			expectError: true,
			errorMsg:    "database error",
		},
		{
			name: "Error - no films",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsByGenre(gomock.Any(), genreID, pager.Count, pager.Offset).
					Return([]models.MainPageFilm{}, nil)
			},
			genreID:     genreID,
			expected:    []models.MainPageFilm{},
			expectError: true,
			errorMsg:    "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilmsByGenre(testContext(), tt.genreID, pager)

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
