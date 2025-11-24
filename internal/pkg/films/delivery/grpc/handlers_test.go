package grpc

import (
	"context"
	"testing"
	"time"

	"kinopoisk/internal/models"
	actors_mocks "kinopoisk/internal/pkg/actors/mocks"
	complitaions_mocks "kinopoisk/internal/pkg/compilations/mocks"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	films_mocks "kinopoisk/internal/pkg/films/mocks"
	genres_mocks "kinopoisk/internal/pkg/genres/mocks"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcFilmsHandler_GetPromoFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilmUsecase := films_mocks.NewMockFilmUsecase(ctrl)
	mockGenreUsecase := genres_mocks.NewMockGenreUsecase(ctrl)
	mockActorUsecase := actors_mocks.NewMockActorUsecase(ctrl)
	mockCompilationUsecase := complitaions_mocks.NewMockCompilationsUsecase(ctrl)
	handler := NewGrpcFilmHandler(mockFilmUsecase, mockGenreUsecase, mockActorUsecase, mockCompilationUsecase)

	filmID := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.EmptyRequest
		mockSetup      func()
		expected       *gen.GetPromoFilmResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name:  "Success",
			input: &gen.EmptyRequest{},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().GetPromoFilm(gomock.Any()).
					Return(models.PromoFilm{
						ID:               filmID,
						Image:            "/images/promo.jpg",
						Title:            "Промо фильм",
						Rating:           8.7,
						ShortDescription: "Краткое описание",
						Year:             2024,
						Genre:            "Драма",
						Duration:         120,
					}, nil)
			},
			expected: &gen.GetPromoFilmResponse{
				Id:               filmID.String(),
				Image:            "/images/promo.jpg",
				Title:            "Промо фильм",
				Rating:           8.7,
				ShortDescription: "Краткое описание",
				Year:             2024,
				Genre:            "Драма",
				Duration:         120,
			},
			expectedErr: nil,
		},
		{
			name:  "Error - Not Found",
			input: &gen.EmptyRequest{},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().GetPromoFilm(gomock.Any()).
					Return(models.PromoFilm{}, films.ErrorNotFound)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.NotFound, "film not found"),
			expectedStatus: codes.NotFound,
		},
		{
			name:  "Error - Internal",
			input: &gen.EmptyRequest{},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().GetPromoFilm(gomock.Any()).
					Return(models.PromoFilm{}, assert.AnError)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.Internal, "internal server error"),
			expectedStatus: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.GetPromoFilm(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected, resp)
			}
		})
	}
}

func TestGrpcFilmsHandler_GetFavFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilmUsecase := films_mocks.NewMockFilmUsecase(ctrl)
	mockGenreUsecase := genres_mocks.NewMockGenreUsecase(ctrl)
	mockActorUsecase := actors_mocks.NewMockActorUsecase(ctrl)
	mockCompilationUsecase := complitaions_mocks.NewMockCompilationsUsecase(ctrl)
	handler := NewGrpcFilmHandler(mockFilmUsecase, mockGenreUsecase, mockActorUsecase, mockCompilationUsecase)

	userID := uuid.NewV4()
	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.GetFavFilmsRequest
		mockSetup      func()
		expected       *gen.GetFavFilmsResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.GetFavFilmsRequest{
				UserId: userID.String(),
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().GetUsersFavFilms(gomock.Any(), userID).
					Return([]models.FavFilm{
						{
							ID:               filmID1,
							Title:            "Фильм 1",
							Genre:            "Драма",
							Year:             2024,
							Duration:         120,
							Image:            "/images/film1.jpg",
							ShortDescription: "Описание 1",
							Rating:           8.5,
						},
						{
							ID:               filmID2,
							Title:            "Фильм 2",
							Genre:            "Комедия",
							Year:             2023,
							Duration:         110,
							Image:            "/images/film2.jpg",
							ShortDescription: "Описание 2",
							Rating:           7.9,
						},
					}, nil)
			},
			expected: &gen.GetFavFilmsResponse{
				Films: []*gen.FavFilm{
					{
						Id:               filmID1.String(),
						Title:            "Фильм 1",
						Genre:            "Драма",
						Year:             2024,
						Duration:         120,
						Image:            "/images/film1.jpg",
						ShortDescription: "Описание 1",
						Rating:           8.5,
					},
					{
						Id:               filmID2.String(),
						Title:            "Фильм 2",
						Genre:            "Комедия",
						Year:             2023,
						Duration:         110,
						Image:            "/images/film2.jpg",
						ShortDescription: "Описание 2",
						Rating:           7.9,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error - Not Found",
			input: &gen.GetFavFilmsRequest{
				UserId: userID.String(),
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().GetUsersFavFilms(gomock.Any(), userID).
					Return(nil, films.ErrorNotFound)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "bad request"),
			expectedStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.GetFavFilms(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, len(tt.expected.Films), len(resp.Films))
			}
		})
	}
}

func TestGrpcFilmsHandler_GetFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilmUsecase := films_mocks.NewMockFilmUsecase(ctrl)
	mockGenreUsecase := genres_mocks.NewMockGenreUsecase(ctrl)
	mockActorUsecase := actors_mocks.NewMockActorUsecase(ctrl)
	mockCompilationUsecase := complitaions_mocks.NewMockCompilationsUsecase(ctrl)
	handler := NewGrpcFilmHandler(mockFilmUsecase, mockGenreUsecase, mockActorUsecase, mockCompilationUsecase)

	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.GetFilmsRequest
		mockSetup      func()
		expected       *gen.GetFilmsResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.GetFilmsRequest{
				Pager: &gen.Pager{Count: 10, Offset: 0},
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().GetFilms(gomock.Any(), models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{
						{
							ID:     filmID1,
							Cover:  "/covers/film1.jpg",
							Title:  "Фильм 1",
							Rating: 8.5,
							Year:   2024,
							Genre:  "Драма",
						},
						{
							ID:     filmID2,
							Cover:  "/covers/film2.jpg",
							Title:  "Фильм 2",
							Rating: 7.9,
							Year:   2023,
							Genre:  "Комедия",
						},
					}, nil)
			},
			expected: &gen.GetFilmsResponse{
				Films: []*gen.MainPageFilm{
					{
						Id:     filmID1.String(),
						Cover:  "/covers/film1.jpg",
						Title:  "Фильм 1",
						Rating: 8.5,
						Year:   2024,
						Genre:  "Драма",
					},
					{
						Id:     filmID2.String(),
						Cover:  "/covers/film2.jpg",
						Title:  "Фильм 2",
						Rating: 7.9,
						Year:   2023,
						Genre:  "Комедия",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error - Not Found",
			input: &gen.GetFilmsRequest{
				Pager: &gen.Pager{Count: 10, Offset: 0},
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().GetFilms(gomock.Any(), gomock.Any()).
					Return(nil, films.ErrorNotFound)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.NotFound, "films not found"),
			expectedStatus: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.GetFilms(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, len(tt.expected.Films), len(resp.Films))
			}
		})
	}
}

func TestGrpcFilmsHandler_GetFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilmUsecase := films_mocks.NewMockFilmUsecase(ctrl)
	mockGenreUsecase := genres_mocks.NewMockGenreUsecase(ctrl)
	mockActorUsecase := actors_mocks.NewMockActorUsecase(ctrl)
	mockCompilationUsecase := complitaions_mocks.NewMockCompilationsUsecase(ctrl)
	handler := NewGrpcFilmHandler(mockFilmUsecase, mockGenreUsecase, mockActorUsecase, mockCompilationUsecase)

	filmID := uuid.NewV4()
	userID := uuid.NewV4()
	genreID := uuid.NewV4()
	actorID := uuid.NewV4()

	originalTitle := "Original Title"
	slogan := "Film Slogan"
	image1 := "/images/image1.jpg"
	image2 := "/images/image2.jpg"
	image3 := "/images/image3.jpg"

	birthDate := time.Now().AddDate(-30, 0, 0)

	tests := []struct {
		name           string
		input          *gen.GetFilmRequest
		mockSetup      func()
		expected       *gen.GetFilmResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.GetFilmRequest{
				FilmId: filmID.String(),
				UserId: userID.String(),
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().GetFilm(gomock.Any(), filmID, userID).
					Return(models.FilmPage{
						ID:               filmID,
						Title:            "Название фильма",
						OriginalTitle:    &originalTitle,
						Cover:            "/covers/cover.jpg",
						Poster:           "/posters/poster.jpg",
						Genre:            "Драма",
						ShortDescription: "Краткое описание",
						Description:      "Полное описание",
						AgeCategory:      "16+",
						Budget:           10000000,
						WorldwideFees:    50000000,
						TrailerURL:       nil,
						NumberOfRatings:  1500,
						Year:             2024,
						Rating:           8.5,
						Country:          "США",
						Slogan:           &slogan,
						Duration:         120,
						Image1:           &image1,
						Image2:           &image2,
						Image3:           &image3,
						Actors: []models.Actor{
							{
								ID:            actorID,
								RussianName:   "Актер Русский",
								OriginalName:  &originalTitle,
								Photo:         "/photos/actor.jpg",
								Height:        180,
								BirthDate:     birthDate,
								DeathDate:     nil,
								ZodiacSign:    "Телец",
								BirthPlace:    "Москва",
								MaritalStatus: "Женат",
							},
						},
						IsReviewed: false,
						UserRating: nil,
						GenreID:    genreID,
						IsLiked:    true,
					}, nil)
			},
			expected: &gen.GetFilmResponse{
				Id:               filmID.String(),
				Title:            "Название фильма",
				OriginalTitle:    &originalTitle,
				Cover:            "/covers/cover.jpg",
				Poster:           "/posters/poster.jpg",
				Genre:            "Драма",
				ShortDescription: "Краткое описание",
				Description:      "Полное описание",
				AgeCategory:      "16+",
				Budget:           10000000,
				WorldwideFees:    50000000,
				TrailerUrl:       nil,
				NumberOfRatings:  1500,
				Year:             2024,
				Rating:           8.5,
				Country:          "США",
				Slogan:           &slogan,
				Duration:         120,
				Image1:           &image1,
				Image2:           &image2,
				Image3:           &image3,
				Actors: []*gen.Actor{
					{
						Id:            actorID.String(),
						RussianName:   "Актер Русский",
						OriginalName:  &originalTitle,
						Photo:         "/photos/actor.jpg",
						Height:        180,
						BirthDate:     birthDate.String(),
						DeathDate:     nil,
						ZodiacSign:    "Телец",
						BirthPlace:    "Москва",
						MaritalStatus: "Женат",
					},
				},
				IsReviewed: false,
				UserRating: nil,
				GenreId:    genreID.String(),
				IsLiked:    true,
			},
			expectedErr: nil,
		},
		{
			name: "Error - Invalid Film ID",
			input: &gen.GetFilmRequest{
				FilmId: "invalid-uuid",
				UserId: userID.String(),
			},
			mockSetup:      func() {},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "invalid film ID"),
			expectedStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.GetFilm(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.Id, resp.Id)
				assert.Equal(t, tt.expected.Title, resp.Title)
			}
		})
	}
}

func TestGrpcFilmsHandler_SaveFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilmUsecase := films_mocks.NewMockFilmUsecase(ctrl)
	mockGenreUsecase := genres_mocks.NewMockGenreUsecase(ctrl)
	mockActorUsecase := actors_mocks.NewMockActorUsecase(ctrl)
	mockCompilationUsecase := complitaions_mocks.NewMockCompilationsUsecase(ctrl)
	handler := NewGrpcFilmHandler(mockFilmUsecase, mockGenreUsecase, mockActorUsecase, mockCompilationUsecase)

	filmID := uuid.NewV4()
	userID := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.SaveFilmRequest
		mockSetup      func()
		expected       *gen.EmptyResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.SaveFilmRequest{
				FilmId: filmID.String(),
				UserId: userID.String(),
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().SaveFilm(gomock.Any(), userID, filmID).
					Return(nil)
			},
			expected:    &gen.EmptyResponse{},
			expectedErr: nil,
		},
		{
			name: "Error - Bad Request",
			input: &gen.SaveFilmRequest{
				FilmId: filmID.String(),
				UserId: userID.String(),
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().SaveFilm(gomock.Any(), userID, filmID).
					Return(films.ErrorBadRequest)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "film already saved"),
			expectedStatus: codes.InvalidArgument,
		},
		{
			name: "Error - Invalid Film ID",
			input: &gen.SaveFilmRequest{
				FilmId: "invalid-uuid",
				UserId: userID.String(),
			},
			mockSetup:      func() {},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "invalid film ID"),
			expectedStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.SaveFilm(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestGrpcFilmsHandler_RemoveFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilmUsecase := films_mocks.NewMockFilmUsecase(ctrl)
	mockGenreUsecase := genres_mocks.NewMockGenreUsecase(ctrl)
	mockActorUsecase := actors_mocks.NewMockActorUsecase(ctrl)
	mockCompilationUsecase := complitaions_mocks.NewMockCompilationsUsecase(ctrl)
	handler := NewGrpcFilmHandler(mockFilmUsecase, mockGenreUsecase, mockActorUsecase, mockCompilationUsecase)

	filmID := uuid.NewV4()
	userID := uuid.NewV4()
	remainingFilmID := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.RemoveFilmRequest
		mockSetup      func()
		expected       *gen.GetFavFilmsResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.RemoveFilmRequest{
				FilmId: filmID.String(),
				UserId: userID.String(),
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().RemoveFilm(gomock.Any(), userID, filmID).
					Return([]models.FavFilm{
						{
							ID:               remainingFilmID,
							Title:            "Оставшийся фильм",
							Genre:            "Драма",
							Year:             2024,
							Duration:         120,
							Image:            "/images/remaining.jpg",
							ShortDescription: "Описание оставшегося",
							Rating:           8.2,
						},
					}, nil)
			},
			expected: &gen.GetFavFilmsResponse{
				Films: []*gen.FavFilm{
					{
						Id:               remainingFilmID.String(),
						Title:            "Оставшийся фильм",
						Genre:            "Драма",
						Year:             2024,
						Duration:         120,
						Image:            "/images/remaining.jpg",
						ShortDescription: "Описание оставшегося",
						Rating:           8.2,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error - Bad Request",
			input: &gen.RemoveFilmRequest{
				FilmId: filmID.String(),
				UserId: userID.String(),
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().RemoveFilm(gomock.Any(), userID, filmID).
					Return(nil, films.ErrorBadRequest)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "nothing to remove"),
			expectedStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.RemoveFilm(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, len(tt.expected.Films), len(resp.Films))
			}
		})
	}
}

func TestGrpcFilmsHandler_SendFeedback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilmUsecase := films_mocks.NewMockFilmUsecase(ctrl)
	mockGenreUsecase := genres_mocks.NewMockGenreUsecase(ctrl)
	mockActorUsecase := actors_mocks.NewMockActorUsecase(ctrl)
	mockCompilationUsecase := complitaions_mocks.NewMockCompilationsUsecase(ctrl)
	handler := NewGrpcFilmHandler(mockFilmUsecase, mockGenreUsecase, mockActorUsecase, mockCompilationUsecase)

	filmID := uuid.NewV4()
	userID := uuid.NewV4()
	feedbackID := uuid.NewV4()

	title := "Отличный фильм!"
	text := "Очень понравилось"

	tests := []struct {
		name           string
		input          *gen.SendFeedbackRequest
		mockSetup      func()
		expected       *gen.SendFeedbackResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.SendFeedbackRequest{
				FilmId: filmID.String(),
				UserId: userID.String(),
				Feedback: &gen.FilmFeedbackInput{
					Title:  title,
					Text:   text,
					Rating: 9,
				},
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().SendFeedback(gomock.Any(), models.FilmFeedbackInput{
					Title:  title,
					Text:   text,
					Rating: 9,
				}, filmID, userID).
					Return(models.FilmFeedback{
						ID:            feedbackID,
						UserID:        userID,
						FilmID:        filmID,
						Title:         &title,
						Text:          &text,
						Rating:        9,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
						UserLogin:     "testuser",
						UserAvatar:    "/avatars/test.jpg",
						IsMine:        true,
						NewFilmRating: 8.7,
					}, nil)
			},
			expected: &gen.SendFeedbackResponse{
				Feedback: &gen.FilmFeedback{
					Id:            feedbackID.String(),
					UserId:        userID.String(),
					FilmId:        filmID.String(),
					Title:         &title,
					Text:          &text,
					Rating:        9,
					CreatedAt:     gomock.Any().String(),
					UpdatedAt:     gomock.Any().String(),
					UserLogin:     "testuser",
					UserAvatar:    "/avatars/test.jpg",
					IsMine:        true,
					NewFilmRating: 8.7,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error - Not Found",
			input: &gen.SendFeedbackRequest{
				FilmId: filmID.String(),
				UserId: userID.String(),
				Feedback: &gen.FilmFeedbackInput{
					Title:  title,
					Text:   text,
					Rating: 9,
				},
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().SendFeedback(gomock.Any(), gomock.Any(), filmID, userID).
					Return(models.FilmFeedback{}, films.ErrorNotFound)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.NotFound, "film not found"),
			expectedStatus: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.SendFeedback(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.Feedback.Id, resp.Feedback.Id)
				assert.Equal(t, tt.expected.Feedback.Rating, resp.Feedback.Rating)
			}
		})
	}
}

func TestGrpcFilmsHandler_ValidateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilmUsecase := films_mocks.NewMockFilmUsecase(ctrl)
	mockGenreUsecase := genres_mocks.NewMockGenreUsecase(ctrl)
	mockActorUsecase := actors_mocks.NewMockActorUsecase(ctrl)
	mockCompilationUsecase := complitaions_mocks.NewMockCompilationsUsecase(ctrl)
	handler := NewGrpcFilmHandler(mockFilmUsecase, mockGenreUsecase, mockActorUsecase, mockCompilationUsecase)

	userID := uuid.NewV4()
	token := "valid-token"

	tests := []struct {
		name           string
		input          *gen.ValidateUserRequest
		mockSetup      func()
		expected       *gen.ValidateUserResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.ValidateUserRequest{
				Token: token,
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().ValidateAndGetUser(gomock.Any(), token).
					Return(models.User{
						ID:      userID,
						Version: 1,
						Login:   "testuser",
						Avatar:  "/avatars/test.jpg",
					}, nil)
			},
			expected: &gen.ValidateUserResponse{
				ID:      userID.String(),
				Version: 1,
				Login:   "testuser",
				Avatar:  "/avatars/test.jpg",
			},
			expectedErr: nil,
		},
		{
			name: "Error - Unauthenticated",
			input: &gen.ValidateUserRequest{
				Token: "invalid-token",
			},
			mockSetup: func() {
				mockFilmUsecase.EXPECT().ValidateAndGetUser(gomock.Any(), "invalid-token").
					Return(models.User{}, films.ErrorNotFound)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.Unauthenticated, "user not found"),
			expectedStatus: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.ValidateUser(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.ID, resp.ID)
				assert.Equal(t, tt.expected.Login, resp.Login)
			}
		})
	}
}

func TestGrpcFilmsHandler_GetGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFilmUsecase := films_mocks.NewMockFilmUsecase(ctrl)
	mockGenreUsecase := genres_mocks.NewMockGenreUsecase(ctrl)
	mockActorUsecase := actors_mocks.NewMockActorUsecase(ctrl)
	mockCompilationUsecase := complitaions_mocks.NewMockCompilationsUsecase(ctrl)
	handler := NewGrpcFilmHandler(mockFilmUsecase, mockGenreUsecase, mockActorUsecase, mockCompilationUsecase)

	genreID := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.GetGenreRequest
		mockSetup      func()
		expected       *gen.GetGenreResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.GetGenreRequest{
				GenreId: genreID.String(),
			},
			mockSetup: func() {
				mockGenreUsecase.EXPECT().GetGenre(gomock.Any(), genreID).
					Return(models.Genre{
						ID:          genreID,
						Title:       "Драма",
						Description: "Драматические фильмы",
						Icon:        "/icons/drama.png",
					}, nil)
			},
			expected: &gen.GetGenreResponse{
				Genre: &gen.Genre{
					Id:          genreID.String(),
					Name:        "Драма",
					Description: "Драматические фильмы",
					Icon:        "/icons/drama.png",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error - Not Found",
			input: &gen.GetGenreRequest{
				GenreId: genreID.String(),
			},
			mockSetup: func() {
				mockGenreUsecase.EXPECT().GetGenre(gomock.Any(), genreID).
					Return(models.Genre{}, films.ErrorNotFound)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.NotFound, "genre not found"),
			expectedStatus: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.GetGenre(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.Genre.Name, resp.Genre.Name)
			}
		})
	}
}
