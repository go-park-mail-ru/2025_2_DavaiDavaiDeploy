package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors"
	"kinopoisk/internal/pkg/actors/mocks"
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

func TestActorUsecase_GetActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockActorRepo(ctrl)
	usecase := NewActorUsecase(mockRepo)

	actorID := uuid.NewV4()
	birthDate := time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC)
	deathDate := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	originalName := "Johnny Depp"

	tests := []struct {
		name        string
		setupMock   func()
		expected    models.ActorPage
		expectError bool
		errorType   error
	}{
		{
			name: "Success - living actor",
			setupMock: func() {
				actor := models.Actor{
					ID:            actorID,
					RussianName:   "Джонни Депп",
					OriginalName:  &originalName,
					Photo:         "photo.jpg",
					Height:        178,
					BirthDate:     birthDate,
					DeathDate:     nil,
					ZodiacSign:    "Козерог",
					BirthPlace:    "Овенсборо, Кентукки, США",
					MaritalStatus: "Разведен",
				}
				mockRepo.EXPECT().
					GetActorByID(gomock.Any(), actorID).
					Return(actor, nil)
				mockRepo.EXPECT().
					GetActorFilmsCount(gomock.Any(), actorID).
					Return(25, nil)
			},
			expected: models.ActorPage{
				ID:            actorID,
				RussianName:   "Джонни Депп",
				OriginalName:  &originalName,
				Photo:         "photo.jpg",
				Height:        178,
				BirthDate:     birthDate,
				Age:           calculateAge(birthDate, time.Now()),
				ZodiacSign:    "Козерог",
				BirthPlace:    "Овенсборо, Кентукки, США",
				MaritalStatus: "Разведен",
				FilmsNumber:   25,
			},
			expectError: false,
		},
		{
			name: "Success - deceased actor",
			setupMock: func() {
				heathName := "Heath Ledger"
				actor := models.Actor{
					ID:            actorID,
					RussianName:   "Хит Леджер",
					OriginalName:  &heathName,
					Photo:         "heath.jpg",
					Height:        185,
					BirthDate:     birthDate,
					DeathDate:     &deathDate,
					ZodiacSign:    "Телец",
					BirthPlace:    "Австралия",
					MaritalStatus: "Не женат",
				}
				mockRepo.EXPECT().
					GetActorByID(gomock.Any(), actorID).
					Return(actor, nil)
				mockRepo.EXPECT().
					GetActorFilmsCount(gomock.Any(), actorID).
					Return(15, nil)
			},
			expected: models.ActorPage{
				ID:            actorID,
				RussianName:   "Хит Леджер",
				OriginalName:  func() *string { s := "Heath Ledger"; return &s }(),
				Photo:         "heath.jpg",
				Height:        185,
				BirthDate:     birthDate,
				Age:           40,
				ZodiacSign:    "Телец",
				BirthPlace:    "Австралия",
				MaritalStatus: "Не женат",
				FilmsNumber:   15,
			},
			expectError: false,
		},
		{
			name: "Error - actor not found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetActorByID(gomock.Any(), actorID).
					Return(models.Actor{}, actors.ErrorNotFound)
			},
			expected:    models.ActorPage{},
			expectError: true,
			errorType:   actors.ErrorNotFound,
		},
		{
			name: "Error - internal server error",
			setupMock: func() {
				mockRepo.EXPECT().
					GetActorByID(gomock.Any(), actorID).
					Return(models.Actor{}, actors.ErrorInternalServerError)
			},
			expected:    models.ActorPage{},
			expectError: true,
			errorType:   actors.ErrorInternalServerError,
		},
		{
			name: "Error - no films count",
			setupMock: func() {
				actor := models.Actor{
					ID:            actorID,
					RussianName:   "Джонни Депп",
					OriginalName:  &originalName,
					Photo:         "photo.jpg",
					Height:        178,
					BirthDate:     birthDate,
					DeathDate:     nil,
					ZodiacSign:    "Козерог",
					BirthPlace:    "Овенсборо, Кентукки, США",
					MaritalStatus: "Разведен",
				}
				mockRepo.EXPECT().
					GetActorByID(gomock.Any(), actorID).
					Return(actor, nil)
				mockRepo.EXPECT().
					GetActorFilmsCount(gomock.Any(), actorID).
					Return(0, actors.ErrorInternalServerError)
			},
			expected:    models.ActorPage{},
			expectError: true,
			errorType:   actors.ErrorInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetActor(testContext(), actorID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.ErrorIs(t, err, tt.errorType)
				}
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.RussianName, result.RussianName)
				assert.Equal(t, tt.expected.Height, result.Height)
				assert.Equal(t, tt.expected.Age, result.Age)
				assert.Equal(t, tt.expected.FilmsNumber, result.FilmsNumber)
			}
		})
	}
}

func TestActorUsecase_GetFilmsByActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockActorRepo(ctrl)
	usecase := NewActorUsecase(mockRepo)

	actorID := uuid.NewV4()
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedFilms := []models.MainPageFilm{
		{
			ID:     uuid.NewV4(),
			Cover:  "film1.jpg",
			Title:  "Пираты Карибского моря",
			Rating: 8.5,
			Year:   2003,
			Genre:  "Приключения",
		},
		{
			ID:     uuid.NewV4(),
			Cover:  "film2.jpg",
			Title:  "Эдвард Руки-ножницы",
			Rating: 7.9,
			Year:   1990,
			Genre:  "Фэнтези",
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
					GetFilmsByActor(gomock.Any(), actorID, pager.Count, pager.Offset).
					Return(expectedFilms, nil)
			},
			expected:    expectedFilms,
			expectError: false,
		},
		{
			name: "Error - repository error",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, pager.Count, pager.Offset).
					Return(nil, actors.ErrorInternalServerError)
			},
			expected:    []models.MainPageFilm{},
			expectError: true,
			errorType:   actors.ErrorInternalServerError,
		},
		{
			name: "Error - actor not found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, pager.Count, pager.Offset).
					Return(nil, actors.ErrorNotFound)
			},
			expected:    []models.MainPageFilm{},
			expectError: true,
			errorType:   actors.ErrorNotFound,
		},
		{
			name: "Error - no films found",
			setupMock: func() {
				mockRepo.EXPECT().
					GetFilmsByActor(gomock.Any(), actorID, pager.Count, pager.Offset).
					Return([]models.MainPageFilm{}, nil)
			},
			expected:    []models.MainPageFilm{},
			expectError: true,
			errorType:   actors.ErrorNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := usecase.GetFilmsByActor(testContext(), actorID, pager)

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

func TestActorUsecase_GetFilmsByActor_EmptyFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockActorRepo(ctrl)
	usecase := NewActorUsecase(mockRepo)

	actorID := uuid.NewV4()
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	mockRepo.EXPECT().
		GetFilmsByActor(gomock.Any(), actorID, pager.Count, pager.Offset).
		Return([]models.MainPageFilm{}, nil)

	films, err := usecase.GetFilmsByActor(testContext(), actorID, pager)

	assert.Error(t, err)
	assert.ErrorIs(t, err, actors.ErrorNotFound)
	assert.Len(t, films, 0)
}

func TestActorUsecase_GetFilmsByActor_WithFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockActorRepo(ctrl)
	usecase := NewActorUsecase(mockRepo)

	actorID := uuid.NewV4()
	pager := models.Pager{
		Count:  10,
		Offset: 0,
	}

	expectedFilms := []models.MainPageFilm{
		{
			ID:     uuid.NewV4(),
			Cover:  "film1.jpg",
			Title:  "Тестовый фильм",
			Rating: 8.0,
			Year:   2020,
			Genre:  "Драма",
		},
	}

	mockRepo.EXPECT().
		GetFilmsByActor(gomock.Any(), actorID, pager.Count, pager.Offset).
		Return(expectedFilms, nil)

	films, err := usecase.GetFilmsByActor(testContext(), actorID, pager)

	assert.NoError(t, err)
	assert.Len(t, films, 1)
	assert.Equal(t, expectedFilms[0].Title, films[0].Title)
}

func calculateAge(birthDate, endDate time.Time) int {
	age := endDate.Year() - birthDate.Year()
	if endDate.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}
