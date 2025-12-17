package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"kinopoisk/internal/models"
	compilations_mocks "kinopoisk/internal/pkg/compilations/mocks"
	films_mocks "kinopoisk/internal/pkg/films/mocks"
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

	compilationsMockRepo := compilations_mocks.NewMockCompilationsRepo(ctrl)
	filmsMockRepo := films_mocks.NewMockFilmRepo(ctrl)
	usecase := NewCompilationUsecase(compilationsMockRepo, filmsMockRepo)

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
				compilationsMockRepo.EXPECT().
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
				compilationsMockRepo.EXPECT().
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

	compilationsMockRepo := compilations_mocks.NewMockCompilationsRepo(ctrl)
	filmsMockRepo := films_mocks.NewMockFilmRepo(ctrl)
	usecase := NewCompilationUsecase(compilationsMockRepo, filmsMockRepo)

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
				compilationsMockRepo.EXPECT().
					GetCompilationsWithPagination(gomock.Any(), pager.Count, pager.Offset).
					Return(expectedCompilations, nil)
			},
			expected:    expectedCompilations,
			expectError: false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				compilationsMockRepo.EXPECT().
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
				compilationsMockRepo.EXPECT().
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

	compilationsMockRepo := compilations_mocks.NewMockCompilationsRepo(ctrl)
	filmsMockRepo := films_mocks.NewMockFilmRepo(ctrl)
	usecase := NewCompilationUsecase(compilationsMockRepo, filmsMockRepo)

	compilationID := uuid.NewV4()
	userID := uuid.NewV4()
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()

	expectedFilms := []models.CompFilm{
		{
			ID:               filmID1,
			Image:            "film1.jpg",
			Title:            "Action Film 1",
			Rating:           8.5,
			Year:             2024,
			Genre:            "Action",
			ShortDescription: "Exciting action film",
			Duration:         120,
			IsLiked:          false,
		},
		{
			ID:               filmID2,
			Image:            "film2.jpg",
			Title:            "Action Film 2",
			Rating:           7.9,
			Year:             2023,
			Genre:            "Action",
			ShortDescription: "Another great action film",
			Duration:         110,
			IsLiked:          false,
		},
	}

	expectedFilmsWithLikes := []models.CompFilm{
		{
			ID:               filmID1,
			Image:            "film1.jpg",
			Title:            "Action Film 1",
			Rating:           8.5,
			Year:             2024,
			Genre:            "Action",
			ShortDescription: "Exciting action film",
			Duration:         120,
			IsLiked:          true,
		},
		{
			ID:               filmID2,
			Image:            "film2.jpg",
			Title:            "Action Film 2",
			Rating:           7.9,
			Year:             2023,
			Genre:            "Action",
			ShortDescription: "Another great action film",
			Duration:         110,
			IsLiked:          false,
		},
	}

	tests := []struct {
		name          string
		setupMock     func()
		compilationID uuid.UUID
		userID        uuid.UUID
		expected      []models.CompFilm
		expectError   bool
		errorMsg      string
	}{
		{
			name: "Success - with like checks",
			setupMock: func() {
				compilationsMockRepo.EXPECT().
					GetFilmsByCompilation(gomock.Any(), compilationID, pager.Count, pager.Offset).
					Return(expectedFilms, nil)

				filmsMockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, filmID1).
					Return(models.FilmFeedback{}, nil)

				filmsMockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, filmID2).
					Return(models.FilmFeedback{}, errors.New("not liked"))
			},
			compilationID: compilationID,
			userID:        userID,
			expected:      expectedFilmsWithLikes,
			expectError:   false,
		},
		{
			name: "Success - no likes",
			setupMock: func() {
				compilationsMockRepo.EXPECT().
					GetFilmsByCompilation(gomock.Any(), compilationID, pager.Count, pager.Offset).
					Return(expectedFilms, nil)

				filmsMockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, filmID1).
					Return(models.FilmFeedback{}, errors.New("not liked"))

				filmsMockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, filmID2).
					Return(models.FilmFeedback{}, errors.New("not liked"))
			},
			compilationID: compilationID,
			userID:        userID,
			expected:      expectedFilms,
			expectError:   false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				compilationsMockRepo.EXPECT().
					GetFilmsByCompilation(gomock.Any(), compilationID, pager.Count, pager.Offset).
					Return(nil, errors.New("database error"))
			},
			compilationID: compilationID,
			userID:        userID,
			expected:      []models.CompFilm{},
			expectError:   true,
			errorMsg:      "database error",
		},
		{
			name: "Error - no films",
			setupMock: func() {
				compilationsMockRepo.EXPECT().
					GetFilmsByCompilation(gomock.Any(), compilationID, pager.Count, pager.Offset).
					Return([]models.CompFilm{}, nil)
			},
			compilationID: compilationID,
			userID:        userID,
			expected:      []models.CompFilm{},
			expectError:   true,
			errorMsg:      "not found",
		},
		{
			name: "Success - like check fails but films still returned",
			setupMock: func() {
				compilationsMockRepo.EXPECT().
					GetFilmsByCompilation(gomock.Any(), compilationID, pager.Count, pager.Offset).
					Return(expectedFilms, nil)

				filmsMockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, filmID1).
					Return(models.FilmFeedback{}, errors.New("unexpected error"))

				filmsMockRepo.EXPECT().
					CheckUserLikeExists(gomock.Any(), userID, filmID2).
					Return(models.FilmFeedback{}, errors.New("unexpected error"))
			},
			compilationID: compilationID,
			userID:        userID,
			expected:      expectedFilms,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilmsByCompilation(testContext(), tt.compilationID, tt.userID, pager)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)

				for i, film := range result {
					assert.Equal(t, tt.expected[i].IsLiked, film.IsLiked)
				}
			}
		})
	}
}

func TestCompilationUsecase_GetFilmsByCompilation_EmptyUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	compilationsMockRepo := compilations_mocks.NewMockCompilationsRepo(ctrl)
	filmsMockRepo := films_mocks.NewMockFilmRepo(ctrl)
	usecase := NewCompilationUsecase(compilationsMockRepo, filmsMockRepo)

	compilationID := uuid.NewV4()
	emptyUserID := uuid.Nil
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedFilms := []models.CompFilm{
		{
			ID:               uuid.NewV4(),
			Image:            "film1.jpg",
			Title:            "Action Film 1",
			Rating:           8.5,
			Year:             2024,
			Genre:            "Action",
			ShortDescription: "Exciting action film",
			Duration:         120,
			IsLiked:          false,
		},
	}

	t.Run("Success with empty user ID", func(t *testing.T) {
		compilationsMockRepo.EXPECT().
			GetFilmsByCompilation(gomock.Any(), compilationID, pager.Count, pager.Offset).
			Return(expectedFilms, nil)

		filmsMockRepo.EXPECT().
			CheckUserLikeExists(gomock.Any(), emptyUserID, gomock.Any()).
			Return(models.FilmFeedback{}, errors.New("not liked")).
			AnyTimes()

		result, err := usecase.GetFilmsByCompilation(testContext(), compilationID, emptyUserID, pager)

		assert.NoError(t, err)
		assert.Equal(t, expectedFilms, result)

		for _, film := range result {
			assert.False(t, film.IsLiked)
		}
	})
}
