package filmHandlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	"kinopoisk/internal/pkg/films/mocks"
	"kinopoisk/internal/pkg/middleware/logger"

	"github.com/gorilla/mux"
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

func testContextWithUser(user models.User) context.Context {
	ctx := testContext()
	return context.WithValue(ctx, auth.UserKey, user)
}

func TestGetPromoFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	promoFilmID := uuid.NewV4()
	expectedPromoFilm := models.PromoFilm{
		ID:               promoFilmID,
		Image:            "/images/promo.jpg",
		Title:            "Промо фильм",
		Rating:           8.7,
		ShortDescription: "Краткое описание промо фильма",
		Year:             2024,
		Genre:            "Драма",
		Duration:         120,
	}

	tests := []struct {
		name           string
		url            string
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name: "Success",
			url:  "/films/promo",
			mockSetup: func() {
				mockClient.EXPECT().
					GetPromoFilm(gomock.Any(), &gen.EmptyRequest{}).
					Return(&gen.GetPromoFilmResponse{
						Id:               promoFilmID.String(),
						Image:            "/images/promo.jpg",
						Title:            "Промо фильм",
						Rating:           8.7,
						ShortDescription: "Краткое описание промо фильма",
						Year:             2024,
						Genre:            "Драма",
						Duration:         120,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name: "Client not found error",
			url:  "/films/promo",
			mockSetup: func() {
				mockClient.EXPECT().
					GetPromoFilm(gomock.Any(), &gen.EmptyRequest{}).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
		{
			name: "Client internal error",
			url:  "/films/promo",
			mockSetup: func() {
				mockClient.EXPECT().
					GetPromoFilm(gomock.Any(), &gen.EmptyRequest{}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(testContext())
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/films/promo", handler.GetPromoFilm)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.PromoFilm
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedPromoFilm, decoded)
			}
		})
	}
}

func TestEasyJson(t *testing.T) {
	n := models.MainPageFilm{}
	json.Marshal(n)
}
func TestGetFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	filmID := uuid.NewV4()
	filmIDStr := filmID.String()
	userID := uuid.NewV4()
	user := models.User{ID: userID}

	originalTitle := "Original Film Title"
	slogan := "Great film slogan"
	image1 := "/images/image1.jpg"
	image2 := "/images/image2.jpg"
	image3 := "/images/image3.jpg"

	expectedFilm := models.FilmPage{
		ID:               filmID,
		Title:            "Название фильма",
		OriginalTitle:    &originalTitle,
		Cover:            "/covers/cover.jpg",
		Poster:           "/posters/poster.jpg",
		Genre:            "Драма",
		ShortDescription: "Краткое описание",
		Description:      "Полное описание фильма",
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
		Actors:           []models.Actor{},
		IsReviewed:       false,
		IsLiked:          false,
		GenreID:          uuid.NewV4(),
	}

	tests := []struct {
		name           string
		url            string
		varsID         string
		context        context.Context
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name:    "Success",
			url:     "/films/" + filmIDStr,
			varsID:  filmIDStr,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					GetFilm(gomock.Any(), &gen.GetFilmRequest{
						FilmId: filmIDStr,
						UserId: userID.String(),
					}).
					Return(&gen.GetFilmResponse{
						Id:               filmIDStr,
						Title:            "Название фильма",
						OriginalTitle:    &originalTitle,
						Cover:            "/covers/cover.jpg",
						Poster:           "/posters/poster.jpg",
						Genre:            "Драма",
						ShortDescription: "Краткое описание",
						Description:      "Полное описание фильма",
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
						Actors:           []*gen.Actor{},
						IsReviewed:       false,
						IsLiked:          false,
						GenreId:          expectedFilm.GenreID.String(),
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/films/not-a-uuid",
			varsID:         "not-a-uuid",
			context:        testContextWithUser(user),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:    "Client not found error",
			url:     "/films/" + filmIDStr,
			varsID:  filmIDStr,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					GetFilm(gomock.Any(), &gen.GetFilmRequest{
						FilmId: filmIDStr,
						UserId: userID.String(),
					}).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(tt.context)
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/films/{id}", handler.GetFilm)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.FilmPage
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedFilm.ID, decoded.ID)
				assert.Equal(t, expectedFilm.Title, decoded.Title)
			}
		})
	}
}

func TestGetFilmFeedbacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	filmID := uuid.NewV4()
	filmIDStr := filmID.String()
	userID := uuid.NewV4()
	user := models.User{ID: userID}

	title1 := "Отличный фильм!"
	text1 := "Потрясающая актерская игра и сюжет"
	title2 := "Хороший фильм"
	text2 := "Приятно посмотреть"

	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	expectedFeedbacks := []models.FilmFeedback{
		{
			ID:         uuid.NewV4(),
			UserID:     uuid.NewV4(),
			FilmID:     filmID,
			Title:      &title1,
			Text:       &text1,
			Rating:     9,
			CreatedAt:  fixedTime,
			UpdatedAt:  fixedTime,
			UserLogin:  "user1",
			UserAvatar: "/avatars/user1.jpg",
			IsMine:     false,
		},
		{
			ID:         uuid.NewV4(),
			UserID:     uuid.NewV4(),
			FilmID:     filmID,
			Title:      &title2,
			Text:       &text2,
			Rating:     8,
			CreatedAt:  fixedTime,
			UpdatedAt:  fixedTime,
			UserLogin:  "user2",
			UserAvatar: "/avatars/user2.jpg",
			IsMine:     false,
		},
	}

	tests := []struct {
		name           string
		url            string
		varsID         string
		context        context.Context
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name:    "Success",
			url:     "/films/" + filmIDStr + "/feedbacks?count=10&offset=0",
			varsID:  filmIDStr,
			context: testContextWithUser(user),
			mockSetup: func() {
				var grpcFeedbacks []*gen.FilmFeedback
				for _, feedback := range expectedFeedbacks {
					grpcFeedbacks = append(grpcFeedbacks, &gen.FilmFeedback{
						Id:         feedback.ID.String(),
						UserId:     feedback.UserID.String(),
						FilmId:     feedback.FilmID.String(),
						Title:      feedback.Title,
						Text:       feedback.Text,
						Rating:     int32(feedback.Rating),
						CreatedAt:  feedback.CreatedAt.String(),
						UpdatedAt:  feedback.UpdatedAt.String(),
						UserLogin:  feedback.UserLogin,
						UserAvatar: feedback.UserAvatar,
						IsMine:     feedback.IsMine,
					})
				}
				mockClient.EXPECT().
					GetFilmFeedbacks(gomock.Any(), &gen.GetFilmFeedbacksRequest{
						FilmId: filmIDStr,
						UserId: userID.String(),
						Pager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(&gen.GetFilmFeedbacksResponse{Feedbacks: grpcFeedbacks}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/films/not-a-uuid/feedbacks",
			varsID:         "not-a-uuid",
			context:        testContextWithUser(user),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:    "Client not found error",
			url:     "/films/" + filmIDStr + "/feedbacks?count=10&offset=0",
			varsID:  filmIDStr,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					GetFilmFeedbacks(gomock.Any(), &gen.GetFilmFeedbacksRequest{
						FilmId: filmIDStr,
						UserId: userID.String(),
						Pager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(tt.context)
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/films/{id}/feedbacks", handler.GetFilmFeedbacks)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.FilmFeedback
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, len(expectedFeedbacks), len(decoded))
			}
		})
	}
}

func TestSendFeedback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	filmID := uuid.NewV4()
	filmIDStr := filmID.String()
	userID := uuid.NewV4()
	title := "Отличный фильм!"
	text := "Очень понравилось"
	avatar := "/avatars/test.jpg"

	user := models.User{
		ID:     userID,
		Login:  "testuser",
		Avatar: avatar,
	}

	expectedFeedback := models.FilmFeedback{
		ID:         uuid.NewV4(),
		UserID:     userID,
		FilmID:     filmID,
		Title:      &title,
		Text:       &text,
		Rating:     9,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		UserLogin:  "testuser",
		UserAvatar: "/avatars/test.jpg",
		IsMine:     true,
	}

	tests := []struct {
		name           string
		url            string
		varsID         string
		body           string
		context        context.Context
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name:    "Success",
			url:     "/films/" + filmIDStr + "/feedbacks",
			varsID:  filmIDStr,
			body:    `{"title": "Отличный фильм!", "text": "Очень понравилось", "rating": 9}`,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					SendFeedback(gomock.Any(), &gen.SendFeedbackRequest{
						FilmId: filmIDStr,
						UserId: userID.String(),
						Feedback: &gen.FilmFeedbackInput{
							Title:  title,
							Text:   text,
							Rating: 9,
						},
					}).
					Return(&gen.SendFeedbackResponse{
						Feedback: &gen.FilmFeedback{
							Id:         expectedFeedback.ID.String(),
							UserId:     userID.String(),
							FilmId:     filmID.String(),
							Title:      &title,
							Text:       &text,
							Rating:     9,
							CreatedAt:  expectedFeedback.CreatedAt.String(),
							UpdatedAt:  expectedFeedback.UpdatedAt.String(),
							UserLogin:  "testuser",
							UserAvatar: "/avatars/test.jpg",
							IsMine:     true,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/films/not-a-uuid/feedbacks",
			varsID:         "not-a-uuid",
			body:           `{"rating": 9}`,
			context:        testContextWithUser(user),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:           "Invalid JSON",
			url:            "/films/" + filmIDStr + "/feedbacks",
			varsID:         filmIDStr,
			body:           `invalid json`,
			context:        testContextWithUser(user),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:    "Client not found error",
			url:     "/films/" + filmIDStr + "/feedbacks",
			varsID:  filmIDStr,
			body:    `{"rating": 9}`,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					SendFeedback(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodPost, tt.url, strings.NewReader(tt.body)).WithContext(tt.context)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/films/{id}/feedbacks", handler.SendFeedback)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.FilmFeedback
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedFeedback.ID, decoded.ID)
				assert.Equal(t, expectedFeedback.UserID, decoded.UserID)
				assert.Equal(t, expectedFeedback.FilmID, decoded.FilmID)
			}
		})
	}
}

func TestSetRating(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	filmID := uuid.NewV4()
	filmIDStr := filmID.String()
	userID := uuid.NewV4()
	avatar := "/avatars/test.jpg"

	user := models.User{
		ID:     userID,
		Login:  "testuser",
		Avatar: avatar,
	}

	expectedRating := models.FilmFeedback{
		ID:         uuid.NewV4(),
		UserID:     userID,
		FilmID:     filmID,
		Title:      nil,
		Text:       nil,
		Rating:     8,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		UserLogin:  "testuser",
		UserAvatar: "/avatars/test.jpg",
		IsMine:     true,
	}

	tests := []struct {
		name           string
		url            string
		varsID         string
		body           string
		context        context.Context
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name:    "Success",
			url:     "/films/" + filmIDStr + "/rating",
			varsID:  filmIDStr,
			body:    `{"rating": 8}`,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					SetRating(gomock.Any(), &gen.SetRatingRequest{
						FilmId: filmIDStr,
						UserId: userID.String(),
						RatingInput: &gen.FilmRatingInput{
							Rating: 8,
						},
					}).
					Return(&gen.SetRatingResponse{
						Feedback: &gen.FilmFeedback{
							Id:         expectedRating.ID.String(),
							UserId:     userID.String(),
							FilmId:     filmID.String(),
							Title:      nil,
							Text:       nil,
							Rating:     8,
							CreatedAt:  expectedRating.CreatedAt.String(),
							UpdatedAt:  expectedRating.UpdatedAt.String(),
							UserLogin:  "testuser",
							UserAvatar: "/avatars/test.jpg",
							IsMine:     true,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Invalid ID",
			url:            "/films/not-a-uuid/rating",
			varsID:         "not-a-uuid",
			body:           `{"rating": 8}`,
			context:        testContextWithUser(user),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:    "Client not found error",
			url:     "/films/" + filmIDStr + "/rating",
			varsID:  filmIDStr,
			body:    `{"rating": 8}`,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					SetRating(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodPost, tt.url, strings.NewReader(tt.body)).WithContext(tt.context)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/films/{id}/rating", handler.SetRating)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded models.FilmFeedback
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedRating.ID, decoded.ID)
				assert.Equal(t, expectedRating.UserID, decoded.UserID)
				assert.Equal(t, expectedRating.FilmID, decoded.FilmID)
				assert.Equal(t, expectedRating.Rating, decoded.Rating)
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	userID := uuid.NewV4()
	token := "valid-token"
	avatar := "/avatars/test.jpg"
	user := models.User{
		ID:     userID,
		Login:  "testuser",
		Avatar: avatar,
	}

	tests := []struct {
		name                string
		cookieValue         string
		mockSetup           func()
		expectUserInContext bool
	}{
		{
			name:        "With valid token",
			cookieValue: token,
			mockSetup: func() {
				mockClient.EXPECT().
					ValidateUser(gomock.Any(), &gen.ValidateUserRequest{Token: token}).
					Return(&gen.ValidateUserResponse{
						ID:      userID.String(),
						Version: 1,
						Login:   "testuser",
						Avatar:  "/avatars/test.jpg",
					}, nil)
			},
			expectUserInContext: true,
		},
		{
			name:        "With invalid token",
			cookieValue: "invalid-token",
			mockSetup: func() {
				mockClient.EXPECT().
					ValidateUser(gomock.Any(), &gen.ValidateUserRequest{Token: "invalid-token"}).
					Return(nil, status.Error(codes.Unauthenticated, "invalid token"))
			},
			expectUserInContext: false,
		},
		{
			name:                "Without cookie",
			cookieValue:         "",
			mockSetup:           func() {},
			expectUserInContext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectUserInContext {
					userFromContext, ok := r.Context().Value(auth.UserKey).(models.User)
					assert.True(t, ok)
					assert.Equal(t, user.ID, userFromContext.ID)
					assert.Equal(t, user.Login, userFromContext.Login)
				} else {
					userFromContext := r.Context().Value(auth.UserKey)
					assert.Nil(t, userFromContext)
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(testContext())
			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{
					Name:  CookieName,
					Value: tt.cookieValue,
				})
			}

			rec := httptest.NewRecorder()

			middleware := handler.Middleware(nextHandler)
			middleware.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestGetUsersFavFilms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	expectedFilms := []models.FavFilm{
		{
			ID:               uuid.NewV4(),
			Title:            "Избранный фильм 1",
			Image:            "/images/fav1.jpg",
			Rating:           8.5,
			Genre:            "Драма",
			Year:             2024,
			Duration:         120,
			ShortDescription: "Краткое описание 1",
		},
		{
			ID:               uuid.NewV4(),
			Title:            "Избранный фильм 2",
			Image:            "/images/fav2.jpg",
			Rating:           7.9,
			Genre:            "Комедия",
			Year:             2023,
			Duration:         110,
			ShortDescription: "Краткое описание 2",
		},
	}

	tests := []struct {
		name           string
		context        context.Context
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name:           "Unauthorized",
			context:        testContext(),
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodGet, "/users/saved", nil).WithContext(tt.context)
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/users/saved", handler.GetUsersFavFilms)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.FavFilm
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedFilms, decoded)
			}
		})
	}
}

func TestGetFilmsForCalendar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	userID := uuid.NewV4()
	user := models.User{ID: userID}

	fixedTime := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	originalTitle := "Original Film"

	expectedFilms := []models.FilmInCalendar{
		{
			ID:               uuid.NewV4(),
			Cover:            "/covers/calendar1.jpg",
			Title:            "Фильм для календаря 1",
			OriginalTitle:    &originalTitle,
			ShortDescription: "Краткое описание календаря 1",
			ReleaseDate:      fixedTime,
			IsLiked:          true,
		},
		{
			ID:               uuid.NewV4(),
			Cover:            "/covers/calendar2.jpg",
			Title:            "Фильм для календаря 2",
			OriginalTitle:    nil,
			ShortDescription: "Краткое описание календаря 2",
			ReleaseDate:      fixedTime.AddDate(0, 1, 0),
			IsLiked:          false,
		},
	}

	tests := []struct {
		name           string
		url            string
		context        context.Context
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name:    "Success",
			url:     "/films/calendar?count=10&offset=0",
			context: testContextWithUser(user),
			mockSetup: func() {
				var grpcFilms []*gen.FilmInCalendar
				for _, film := range expectedFilms {
					grpcFilms = append(grpcFilms, &gen.FilmInCalendar{
						ID:               film.ID.String(),
						Cover:            film.Cover,
						Title:            film.Title,
						OriginalTitle:    film.OriginalTitle,
						ShortDescription: film.ShortDescription,
						ReleaseDate:      film.ReleaseDate.String(),
						IsLiked:          film.IsLiked,
					})
				}
				mockClient.EXPECT().
					GetFilmsForCalendar(gomock.Any(), &gen.GetFilmsForCalendarRequest{
						UserId: userID.String(),
						Pager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(&gen.GetFilmsForCalendarResponse{Films: grpcFilms}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:    "Client not found error",
			url:     "/films/calendar?count=10&offset=0",
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					GetFilmsForCalendar(gomock.Any(), &gen.GetFilmsForCalendarRequest{
						UserId: userID.String(),
						Pager:  &gen.Pager{Count: 10, Offset: 0},
					}).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil).WithContext(tt.context)
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/films/calendar", handler.GetFilmsForCalendar)
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.FilmInCalendar
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, len(expectedFilms), len(decoded))
			}
		})
	}
}

func TestSaveFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	filmID := uuid.NewV4()
	filmIDStr := filmID.String()
	userID := uuid.NewV4()
	user := models.User{ID: userID}

	tests := []struct {
		name           string
		url            string
		varsID         string
		context        context.Context
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:    "Success",
			url:     "/films/" + filmIDStr + "/save",
			varsID:  filmIDStr,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					SaveFilm(gomock.Any(), &gen.SaveFilmRequest{
						UserId: userID.String(),
						FilmId: filmIDStr,
					}).
					Return(&gen.EmptyResponse{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized",
			url:            "/films/" + filmIDStr + "/save",
			varsID:         filmIDStr,
			context:        testContext(),
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid ID",
			url:            "/films/not-a-uuid/save",
			varsID:         "not-a-uuid",
			context:        testContextWithUser(user),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Client not found error",
			url:     "/films/" + filmIDStr + "/save",
			varsID:  filmIDStr,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					SaveFilm(gomock.Any(), &gen.SaveFilmRequest{
						UserId: userID.String(),
						FilmId: filmIDStr,
					}).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodPost, tt.url, nil).WithContext(tt.context)
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/films/{id}/save", handler.SaveFilm)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestRemoveFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilmsClient(ctrl)
	handler := NewFilmHandler(mockClient)

	filmID := uuid.NewV4()
	filmIDStr := filmID.String()
	userID := uuid.NewV4()
	user := models.User{ID: userID}

	expectedFilms := []models.FavFilm{
		{
			ID:               uuid.NewV4(),
			Title:            "Фильм1",
			Image:            "/images/film1.jpg",
			Rating:           8.2,
			Genre:            "Драма",
			Year:             2024,
			Duration:         130,
			ShortDescription: "Краткое описание фильма 1",
		},
	}

	tests := []struct {
		name           string
		url            string
		varsID         string
		context        context.Context
		mockSetup      func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name:    "Success",
			url:     "/films/" + filmIDStr + "/remove",
			varsID:  filmIDStr,
			context: testContextWithUser(user),
			mockSetup: func() {
				var grpcFilms []*gen.FavFilm
				for _, film := range expectedFilms {
					grpcFilms = append(grpcFilms, &gen.FavFilm{
						Id:               film.ID.String(),
						Title:            film.Title,
						Image:            film.Image,
						Rating:           film.Rating,
						Genre:            film.Genre,
						Year:             int32(film.Year),
						Duration:         int32(film.Duration),
						ShortDescription: film.ShortDescription,
					})
				}
				mockClient.EXPECT().
					RemoveFilm(gomock.Any(), &gen.RemoveFilmRequest{
						UserId: userID.String(),
						FilmId: filmIDStr,
					}).
					Return(&gen.GetFavFilmsResponse{Films: grpcFilms}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name:           "Unauthorized",
			url:            "/films/" + filmIDStr + "/remove",
			varsID:         filmIDStr,
			context:        testContext(),
			mockSetup:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectBody:     false,
		},
		{
			name:           "Invalid ID",
			url:            "/films/not-a-uuid/remove",
			varsID:         "not-a-uuid",
			context:        testContextWithUser(user),
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:    "Client not found error",
			url:     "/films/" + filmIDStr + "/remove",
			varsID:  filmIDStr,
			context: testContextWithUser(user),
			mockSetup: func() {
				mockClient.EXPECT().
					RemoveFilm(gomock.Any(), &gen.RemoveFilmRequest{
						UserId: userID.String(),
						FilmId: filmIDStr,
					}).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := httptest.NewRequest(http.MethodDelete, tt.url, nil).WithContext(tt.context)
			rec := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/films/{id}/remove", handler.RemoveFilm)

			if tt.varsID != "" {
				req = mux.SetURLVars(req, map[string]string{"id": tt.varsID})
			}
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectBody {
				var decoded []models.FavFilm
				err := json.Unmarshal(rec.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, expectedFilms, decoded)
			}
		})
	}
}
