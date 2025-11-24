package grpc

import (
	"context"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/search/delivery/grpc/gen"
	search_mocks "kinopoisk/internal/pkg/search/mocks"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcSearchHandler_SearchFilmsAndActors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearchUsecase := search_mocks.NewMockSearchUsecase(ctrl)
	handler := NewGrpcSearchHandler(mockSearchUsecase)

	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()
	actorID1 := uuid.NewV4()
	actorID2 := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.SearchFilmsAndActorsRequest
		mockSetup      func()
		expected       *gen.SearchFilmsAndActorsResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success - Films and Actors Found",
			input: &gen.SearchFilmsAndActorsRequest{
				SearchString: "matrix",
				FilmsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
				ActorsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockSearchUsecase.EXPECT().GetFilmsFromSearch(gomock.Any(), "matrix", models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{
						{
							ID:     filmID1,
							Cover:  "/covers/matrix.jpg",
							Title:  "Матрица",
							Rating: 8.7,
							Year:   1999,
							Genre:  "Фантастика",
						},
						{
							ID:     filmID2,
							Cover:  "/covers/matrix-reloaded.jpg",
							Title:  "Матрица: Перезагрузка",
							Rating: 7.2,
							Year:   2003,
							Genre:  "Фантастика",
						},
					}, nil)

				mockSearchUsecase.EXPECT().GetActorsFromSearch(gomock.Any(), "matrix", models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageActor{
						{
							ID:          actorID1,
							RussianName: "Киану Ривз",
							Photo:       "/photos/keanu.jpg",
						},
						{
							ID:          actorID2,
							RussianName: "Лоренс Фишберн",
							Photo:       "/photos/laurence.jpg",
						},
					}, nil)
			},
			expected: &gen.SearchFilmsAndActorsResponse{
				Films: []*gen.MainPageFilm{
					{
						ID:     filmID1.String(),
						Cover:  "/covers/matrix.jpg",
						Title:  "Матрица",
						Rating: 8.7,
						Year:   1999,
						Genre:  "Фантастика",
					},
					{
						ID:     filmID2.String(),
						Cover:  "/covers/matrix-reloaded.jpg",
						Title:  "Матрица: Перезагрузка",
						Rating: 7.2,
						Year:   2003,
						Genre:  "Фантастика",
					},
				},
				Actors: []*gen.MainPageActor{
					{
						ID:          actorID1.String(),
						RussianName: "Киану Ривз",
						Photo:       "/photos/keanu.jpg",
					},
					{
						ID:          actorID2.String(),
						RussianName: "Лоренс Фишберн",
						Photo:       "/photos/laurence.jpg",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Success - Only Films Found",
			input: &gen.SearchFilmsAndActorsRequest{
				SearchString: "avatar",
				FilmsPager: &gen.Pager{
					Count:  5,
					Offset: 0,
				},
				ActorsPager: &gen.Pager{
					Count:  5,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockSearchUsecase.EXPECT().GetFilmsFromSearch(gomock.Any(), "avatar", models.Pager{Count: 5, Offset: 0}).
					Return([]models.MainPageFilm{
						{
							ID:     filmID1,
							Cover:  "/covers/avatar.jpg",
							Title:  "Аватар",
							Rating: 7.8,
							Year:   2009,
							Genre:  "Фантастика",
						},
					}, nil)

				mockSearchUsecase.EXPECT().GetActorsFromSearch(gomock.Any(), "avatar", models.Pager{Count: 5, Offset: 0}).
					Return([]models.MainPageActor{}, nil)
			},
			expected: &gen.SearchFilmsAndActorsResponse{
				Films: []*gen.MainPageFilm{
					{
						ID:     filmID1.String(),
						Cover:  "/covers/avatar.jpg",
						Title:  "Аватар",
						Rating: 7.8,
						Year:   2009,
						Genre:  "Фантастика",
					},
				},
				Actors: []*gen.MainPageActor{},
			},
			expectedErr: nil,
		},
		{
			name: "Success - Only Actors Found",
			input: &gen.SearchFilmsAndActorsRequest{
				SearchString: "leonardo",
				FilmsPager: &gen.Pager{
					Count:  5,
					Offset: 0,
				},
				ActorsPager: &gen.Pager{
					Count:  5,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockSearchUsecase.EXPECT().GetFilmsFromSearch(gomock.Any(), "leonardo", models.Pager{Count: 5, Offset: 0}).
					Return([]models.MainPageFilm{}, nil)

				mockSearchUsecase.EXPECT().GetActorsFromSearch(gomock.Any(), "leonardo", models.Pager{Count: 5, Offset: 0}).
					Return([]models.MainPageActor{
						{
							ID:          actorID1,
							RussianName: "Леонардо ДиКаприо",
							Photo:       "/photos/leo.jpg",
						},
					}, nil)
			},
			expected: &gen.SearchFilmsAndActorsResponse{
				Films: []*gen.MainPageFilm{},
				Actors: []*gen.MainPageActor{
					{
						ID:          actorID1.String(),
						RussianName: "Леонардо ДиКаприо",
						Photo:       "/photos/leo.jpg",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Success - Empty Search Results",
			input: &gen.SearchFilmsAndActorsRequest{
				SearchString: "nonexistentquery",
				FilmsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
				ActorsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockSearchUsecase.EXPECT().GetFilmsFromSearch(gomock.Any(), "nonexistentquery", models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{}, nil)
				mockSearchUsecase.EXPECT().GetActorsFromSearch(gomock.Any(), "nonexistentquery", models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageActor{}, nil)
			},
			expected: &gen.SearchFilmsAndActorsResponse{
				Films:  []*gen.MainPageFilm{},
				Actors: []*gen.MainPageActor{},
			},
			expectedErr: nil,
		},
		{
			name: "Error - Films Search Fails",
			input: &gen.SearchFilmsAndActorsRequest{
				SearchString: "error",
				FilmsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
				ActorsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockSearchUsecase.EXPECT().GetFilmsFromSearch(gomock.Any(), "error", models.Pager{Count: 10, Offset: 0}).
					Return(nil, assert.AnError)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.Internal, "internal server error"),
			expectedStatus: codes.Internal,
		},
		{
			name: "Error - Actors Search Fails",
			input: &gen.SearchFilmsAndActorsRequest{
				SearchString: "error",
				FilmsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
				ActorsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockSearchUsecase.EXPECT().GetFilmsFromSearch(gomock.Any(), "error", models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{}, nil)

				mockSearchUsecase.EXPECT().GetActorsFromSearch(gomock.Any(), "error", models.Pager{Count: 10, Offset: 0}).
					Return(nil, assert.AnError)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.Internal, "internal server error"),
			expectedStatus: codes.Internal,
		},
		{
			name: "Success - With Pagination",
			input: &gen.SearchFilmsAndActorsRequest{
				SearchString: "test",
				FilmsPager: &gen.Pager{
					Count:  5,
					Offset: 10,
				},
				ActorsPager: &gen.Pager{
					Count:  3,
					Offset: 5,
				},
			},
			mockSetup: func() {
				mockSearchUsecase.EXPECT().GetFilmsFromSearch(gomock.Any(), "test", models.Pager{Count: 5, Offset: 10}).
					Return([]models.MainPageFilm{
						{
							ID:     filmID1,
							Cover:  "/covers/test1.jpg",
							Title:  "Тестовый фильм 1",
							Rating: 6.5,
							Year:   2020,
							Genre:  "Драма",
						},
					}, nil)

				mockSearchUsecase.EXPECT().GetActorsFromSearch(gomock.Any(), "test", models.Pager{Count: 3, Offset: 5}).
					Return([]models.MainPageActor{
						{
							ID:          actorID1,
							RussianName: "Тестовый актер",
							Photo:       "/photos/test.jpg",
						},
					}, nil)
			},
			expected: &gen.SearchFilmsAndActorsResponse{
				Films: []*gen.MainPageFilm{
					{
						ID:     filmID1.String(),
						Cover:  "/covers/test1.jpg",
						Title:  "Тестовый фильм 1",
						Rating: 6.5,
						Year:   2020,
						Genre:  "Драма",
					},
				},
				Actors: []*gen.MainPageActor{
					{
						ID:          actorID1.String(),
						RussianName: "Тестовый актер",
						Photo:       "/photos/test.jpg",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Success - Empty Search String",
			input: &gen.SearchFilmsAndActorsRequest{
				SearchString: "",
				FilmsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
				ActorsPager: &gen.Pager{
					Count:  10,
					Offset: 0,
				},
			},
			mockSetup: func() {
				mockSearchUsecase.EXPECT().GetFilmsFromSearch(gomock.Any(), "", models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageFilm{}, nil)

				mockSearchUsecase.EXPECT().GetActorsFromSearch(gomock.Any(), "", models.Pager{Count: 10, Offset: 0}).
					Return([]models.MainPageActor{}, nil)
			},
			expected: &gen.SearchFilmsAndActorsResponse{
				Films:  []*gen.MainPageFilm{},
				Actors: []*gen.MainPageActor{},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.SearchFilmsAndActors(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, len(tt.expected.Films), len(resp.Films))
				assert.Equal(t, len(tt.expected.Actors), len(resp.Actors))

				if len(tt.expected.Films) > 0 {
					assert.Equal(t, tt.expected.Films[0].Title, resp.Films[0].Title)
					assert.Equal(t, tt.expected.Films[0].Rating, resp.Films[0].Rating)
				}

				if len(tt.expected.Actors) > 0 {
					assert.Equal(t, tt.expected.Actors[0].RussianName, resp.Actors[0].RussianName)
				}
			}
		})
	}
}
