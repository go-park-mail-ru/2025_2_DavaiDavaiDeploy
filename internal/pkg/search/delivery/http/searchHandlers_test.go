package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"
	"kinopoisk/internal/pkg/search/delivery/grpc/gen"
	"kinopoisk/internal/pkg/search/mocks"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func testContext() context.Context {
	testLogger := testLogger()
	return context.WithValue(context.Background(), logger.LoggerKey, testLogger)
}

func TestSearchHandler_GetFilmsAndActorsFromSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockSearchClient(ctrl)
	handler := NewSearchHandler(mockClient)

	filmID1 := uuid.NewV4()
	filmID2 := uuid.NewV4()
	actorID1 := uuid.NewV4()
	actorID2 := uuid.NewV4()

	expectedResponse := models.SearchResponse{
		Films: []models.MainPageFilm{
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
		},
		Actors: []models.MainPageActor{
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
		},
	}

	tests := []struct {
		name           string
		url            string
		mockSetup      func()
		expectedStatus int
		expectBody     bool
		expectedBody   models.SearchResponse
	}{
		{
			name: "Success - Films and Actors Found",
			url:  "/search?q=matrix&films_count=10&films_offset=0&actors_count=10&actors_offset=0",
			mockSetup: func() {
				mockClient.EXPECT().
					SearchFilmsAndActors(gomock.Any(), &gen.SearchFilmsAndActorsRequest{
						SearchString: "matrix",
						FilmsPager:   &gen.Pager{Count: 10, Offset: 0},
						ActorsPager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(&gen.SearchFilmsAndActorsResponse{
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
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedBody:   expectedResponse,
		},
		{
			name: "Success - Only Films Found",
			url:  "/search?q=avatar&films_count=5&films_offset=0&actors_count=5&actors_offset=0",
			mockSetup: func() {
				mockClient.EXPECT().
					SearchFilmsAndActors(gomock.Any(), &gen.SearchFilmsAndActorsRequest{
						SearchString: "avatar",
						FilmsPager:   &gen.Pager{Count: 5, Offset: 0},
						ActorsPager:  &gen.Pager{Count: 5, Offset: 0},
					}).
					Return(&gen.SearchFilmsAndActorsResponse{
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
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedBody: models.SearchResponse{
				Films: []models.MainPageFilm{
					{
						ID:     filmID1,
						Cover:  "/covers/avatar.jpg",
						Title:  "Аватар",
						Rating: 7.8,
						Year:   2009,
						Genre:  "Фантастика",
					},
				},
				Actors: []models.MainPageActor{},
			},
		},
		{
			name: "Success - Only Actors Found",
			url:  "/search?q=leonardo&films_count=5&films_offset=0&actors_count=5&actors_offset=0",
			mockSetup: func() {
				mockClient.EXPECT().
					SearchFilmsAndActors(gomock.Any(), &gen.SearchFilmsAndActorsRequest{
						SearchString: "leonardo",
						FilmsPager:   &gen.Pager{Count: 5, Offset: 0},
						ActorsPager:  &gen.Pager{Count: 5, Offset: 0},
					}).
					Return(&gen.SearchFilmsAndActorsResponse{
						Films: []*gen.MainPageFilm{},
						Actors: []*gen.MainPageActor{
							{
								ID:          actorID1.String(),
								RussianName: "Леонардо ДиКаприо",
								Photo:       "/photos/leo.jpg",
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedBody: models.SearchResponse{
				Films: []models.MainPageFilm{},
				Actors: []models.MainPageActor{
					{
						ID:          actorID1,
						RussianName: "Леонардо ДиКаприо",
						Photo:       "/photos/leo.jpg",
					},
				},
			},
		},
		{
			name: "Success - Empty Search Results",
			url:  "/search?q=nonexistentquery&films_count=10&films_offset=0&actors_count=10&actors_offset=0",
			mockSetup: func() {
				mockClient.EXPECT().
					SearchFilmsAndActors(gomock.Any(), &gen.SearchFilmsAndActorsRequest{
						SearchString: "nonexistentquery",
						FilmsPager:   &gen.Pager{Count: 10, Offset: 0},
						ActorsPager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(&gen.SearchFilmsAndActorsResponse{
						Films:  []*gen.MainPageFilm{},
						Actors: []*gen.MainPageActor{},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedBody: models.SearchResponse{
				Films:  []models.MainPageFilm{},
				Actors: []models.MainPageActor{},
			},
		},
		{
			name: "Success - Escaped Search String",
			url:  "/search?q=matrix%20reloaded&films_count=10&films_offset=0&actors_count=10&actors_offset=0",
			mockSetup: func() {
				mockClient.EXPECT().
					SearchFilmsAndActors(gomock.Any(), &gen.SearchFilmsAndActorsRequest{
						SearchString: "matrix reloaded",
						FilmsPager:   &gen.Pager{Count: 10, Offset: 0},
						ActorsPager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(&gen.SearchFilmsAndActorsResponse{
						Films: []*gen.MainPageFilm{
							{
								ID:     filmID1.String(),
								Cover:  "/covers/matrix-reloaded.jpg",
								Title:  "Матрица: Перезагрузка",
								Rating: 7.2,
								Year:   2003,
								Genre:  "Фантастика",
							},
						},
						Actors: []*gen.MainPageActor{},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedBody: models.SearchResponse{
				Films: []models.MainPageFilm{
					{
						ID:     filmID1,
						Cover:  "/covers/matrix-reloaded.jpg",
						Title:  "Матрица: Перезагрузка",
						Rating: 7.2,
						Year:   2003,
						Genre:  "Фантастика",
					},
				},
				Actors: []models.MainPageActor{},
			},
		},
		{
			name: "Success - Default Pagination",
			url:  "/search?q=test",
			mockSetup: func() {
				mockClient.EXPECT().
					SearchFilmsAndActors(gomock.Any(), &gen.SearchFilmsAndActorsRequest{
						SearchString: "test",
						FilmsPager:   &gen.Pager{Count: 10, Offset: 0},
						ActorsPager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(&gen.SearchFilmsAndActorsResponse{
						Films:  []*gen.MainPageFilm{},
						Actors: []*gen.MainPageActor{},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedBody: models.SearchResponse{
				Films:  []models.MainPageFilm{},
				Actors: []models.MainPageActor{},
			},
		},
		{
			name: "Error - gRPC Client Error",
			url:  "/search?q=error&films_count=10&films_offset=0&actors_count=10&actors_offset=0",
			mockSetup: func() {
				mockClient.EXPECT().
					SearchFilmsAndActors(gomock.Any(), &gen.SearchFilmsAndActorsRequest{
						SearchString: "error",
						FilmsPager:   &gen.Pager{Count: 10, Offset: 0},
						ActorsPager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name: "Success - Empty Search String",
			url:  "/search?q=&films_count=10&films_offset=0&actors_count=10&actors_offset=0",
			mockSetup: func() {
				mockClient.EXPECT().
					SearchFilmsAndActors(gomock.Any(), &gen.SearchFilmsAndActorsRequest{
						SearchString: "",
						FilmsPager:   &gen.Pager{Count: 10, Offset: 0},
						ActorsPager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(&gen.SearchFilmsAndActorsResponse{
						Films:  []*gen.MainPageFilm{},
						Actors: []*gen.MainPageActor{},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedBody: models.SearchResponse{
				Films:  []models.MainPageFilm{},
				Actors: []models.MainPageActor{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			handler.GetFilmsAndActorsFromSearch(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.SearchResponse
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedBody.Films), len(decoded.Films))
				assert.Equal(t, len(tt.expectedBody.Actors), len(decoded.Actors))

				if len(tt.expectedBody.Films) > 0 {
					assert.Equal(t, tt.expectedBody.Films[0].Title, decoded.Films[0].Title)
					assert.Equal(t, tt.expectedBody.Films[0].Rating, decoded.Films[0].Rating)
				}

				if len(tt.expectedBody.Actors) > 0 {
					assert.Equal(t, tt.expectedBody.Actors[0].RussianName, decoded.Actors[0].RussianName)
				}
			}
		})
	}
}
