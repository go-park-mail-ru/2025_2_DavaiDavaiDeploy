package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/compilations/mocks"
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

func TestCompilationUsecase_GetCompilation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCompilationsRepo(ctrl)
	usecase := NewCompilationUsecase(mockRepo)

	compilationID := uuid.NewV4()
	expectedCompilation := models.Compilation{
		ID:          compilationID,
		Title:       "Best Action Films",
		Description: "Collection of the best action films",
		Icon:        "action-collection.png",
	}

	tests := []struct {
		name          string
		setupMock     func()
		compilationID uuid.UUID
		expected      models.Compilation
		expectError   bool
		errorMsg      string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetCompilationByID(gomock.Any(), compilationID).
					Return(expectedCompilation, nil)
			},
			compilationID: compilationID,
			expected:      expectedCompilation,
			expectError:   false,
		},
		{
			name: "Error - compilation not found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetCompilationByID(gomock.Any(), compilationID).
					Return(models.Compilation{}, errors.New("not found"))
			},
			compilationID: compilationID,
			expected:      models.Compilation{},
			expectError:   true,
			errorMsg:      "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetCompilation(testContext(), tt.compilationID)

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

func TestCompilationUsecase_GetCompilations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCompilationsRepo(ctrl)
	usecase := NewCompilationUsecase(mockRepo)

	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedCompilations := []models.Compilation{
		{
			ID:          uuid.NewV4(),
			Title:       "Action Collection",
			Description: "Best action films",
			Icon:        "action-collection.png",
		},
		{
			ID:          uuid.NewV4(),
			Title:       "Comedy Collection",
			Description: "Funniest comedy films",
			Icon:        "comedy-collection.png",
		},
	}

	tests := []struct {
		name        string
		setupMock   func()
		expected    []models.Compilation
		expectError bool
		errorMsg    string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetCompilationsWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return(expectedCompilations, nil)
			},
			expected:    expectedCompilations,
			expectError: false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				mockRepo.EXPECT().
					GetCompilationsWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return(nil, errors.New("database error"))
			},
			expected:    []models.Compilation{},
			expectError: true,
			errorMsg:    "database error",
		},
		{
			name: "Error - no compilations",
			setupMock: func() {
				mockRepo.EXPECT().
					GetCompilationsWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return([]models.Compilation{}, nil)
			},
			expected:    []models.Compilation{},
			expectError: true,
			errorMsg:    "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetCompilations(testContext(), pager)

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

func TestCompilationUsecase_GetFilmsByCompilation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCompilationsRepo(ctrl)
	usecase := NewCompilationUsecase(mockRepo)

	compilationID := uuid.NewV4()
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedFilms := []models.FavFilm{
		{
			ID:               uuid.NewV4(),
			Title:            "Action Film 1",
			Genre:            "Action",
			Year:             2024,
			Duration:         120,
			Image:            "film1.jpg",
			ShortDescription: "Exciting action film",
			Rating:           8.5,
		},
		{
			ID:               uuid.NewV4(),
			Title:            "Action Film 2",
			Genre:            "Action",
			Year:             2023,
			Duration:         110,
			Image:            "film2.jpg",
			ShortDescription: "Another great action film",
			Rating:           7.9,
		},
	}

	tests := []struct {
		name          string
		setupMock     func()
		compilationID uuid.UUID
		expected      []models.FavFilm
		expectError   bool
		errorMsg      string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsByCompilation(gomock.Any(), compilationID, pager.Count, pager.Offset).
					Return(expectedFilms, nil)
			},
			compilationID: compilationID,
			expected:      expectedFilms,
			expectError:   false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsByCompilation(gomock.Any(), compilationID, pager.Count, pager.Offset).
					Return(nil, errors.New("database error"))
			},
			compilationID: compilationID,
			expected:      []models.FavFilm{},
			expectError:   true,
			errorMsg:      "database error",
		},
		{
			name: "Error - no films",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsByCompilation(gomock.Any(), compilationID, pager.Count, pager.Offset).
					Return([]models.FavFilm{}, nil)
			},
			compilationID: compilationID,
			expected:      []models.FavFilm{},
			expectError:   true,
			errorMsg:      "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilmsByCompilation(testContext(), tt.compilationID, pager)

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
