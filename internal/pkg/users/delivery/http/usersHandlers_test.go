package http

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/middleware/logger"
	"kinopoisk/internal/pkg/users"
	"kinopoisk/internal/pkg/users/mocks"

	"github.com/gorilla/mux"
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

func TestGetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		ucErr          error
		expectedStatus int
	}{
		{
			name:           "Success",
			userID:         uuid.NewV4().String(),
			ucErr:          nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid UUID",
			userID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Internal server error",
			userID:         uuid.NewV4().String(),
			ucErr:          users.ErrorInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Generic error",
			userID:         uuid.NewV4().String(),
			ucErr:          errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mocks.NewMockUsersUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "Invalid UUID" {
				userUUID, _ := uuid.FromString(tt.userID)
				mockUsecase.EXPECT().GetUser(gomock.Any(), userUUID).
					Return(models.User{
						ID:    userUUID,
						Login: "testuser",
					}, tt.ucErr)
			}

			router := mux.NewRouter()
			handler := NewUserHandler(mockUsecase)
			router.HandleFunc("/users/{id}", handler.GetUser)

			r := httptest.NewRequest("GET", "/users/"+tt.userID, nil).WithContext(testContext())
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestChangePassword(t *testing.T) {
	type args struct {
		oldPassword string
		newPassword string
	}

	tests := []struct {
		name           string
		requestBody    string
		args           args
		userID         uuid.UUID
		ucErr          error
		expectedStatus int
	}{
		{
			name:        "Success",
			requestBody: `{"old_password":"oldPass123","new_password":"newPass123"}`,
			args: args{
				oldPassword: "oldPass123",
				newPassword: "newPass123",
			},
			userID:         uuid.NewV4(),
			ucErr:          nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"old_password":"oldPass123","new_password":"newPass123"`,
			userID:         uuid.NewV4(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Wrong old password",
			requestBody: `{"old_password":"wrongOldPass","new_password":"newPass123"}`,
			args: args{
				oldPassword: "wrongOldPass",
				newPassword: "newPass123",
			},
			userID:         uuid.NewV4(),
			ucErr:          users.ErrorBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "No user in context",
			requestBody:    `{"old_password":"oldPass123","new_password":"newPass123"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:        "Internal server error",
			requestBody: `{"old_password":"oldPass123","new_password":"newPass123"}`,
			args: args{
				oldPassword: "oldPass123",
				newPassword: "newPass123",
			},
			userID:         uuid.NewV4(),
			ucErr:          users.ErrorInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mocks.NewMockUsersUsecase(ctrl)
			defer ctrl.Finish()

			ctx := testContext()
			if tt.name != "No user in context" {
				ctx = context.WithValue(ctx, users.UserKey, tt.userID)
			}

			if tt.name != "Invalid JSON" && tt.name != "No user in context" {
				mockUsecase.EXPECT().ChangePassword(gomock.Any(), tt.userID, tt.args.oldPassword, tt.args.newPassword).
					Return(models.User{
						ID:    tt.userID,
						Login: "testuser",
					}, "new_jwt_token", tt.ucErr)
			}

			r := httptest.NewRequest("PUT", "/users/password", bytes.NewBufferString(tt.requestBody)).WithContext(ctx)
			w := httptest.NewRecorder()

			handler := NewUserHandler(mockUsecase)
			handler.ChangePassword(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		setupMocks     func(mockUsecase *mocks.MockUsersUsecase)
		expectedStatus int
	}{
		{
			name: "Success",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "csrf-token"})
				req.AddCookie(&http.Cookie{Name: CookieName, Value: "jwt-token"})
				req.Header.Set("X-CSRF-Token", "csrf-token")
				return req.WithContext(testContext())
			},
			setupMocks: func(mockUsecase *mocks.MockUsersUsecase) {
				mockUsecase.EXPECT().ValidateAndGetUser(gomock.Any(), "jwt-token").
					Return(models.User{ID: uuid.NewV4()}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "No CSRF cookie",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.AddCookie(&http.Cookie{Name: CookieName, Value: "jwt-token"})
				req.Header.Set("X-CSRF-Token", "csrf-token")
				return req.WithContext(testContext())
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "No CSRF token in header or form",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "csrf-token"})
				req.AddCookie(&http.Cookie{Name: CookieName, Value: "jwt-token"})
				return req.WithContext(testContext())
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "CSRF token mismatch",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "csrf-token"})
				req.AddCookie(&http.Cookie{Name: CookieName, Value: "jwt-token"})
				req.Header.Set("X-CSRF-Token", "different-csrf-token")
				return req.WithContext(testContext())
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid JWT token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "csrf-token"})
				req.AddCookie(&http.Cookie{Name: CookieName, Value: "invalid-jwt-token"})
				req.Header.Set("X-CSRF-Token", "csrf-token")
				return req.WithContext(testContext())
			},
			setupMocks: func(mockUsecase *mocks.MockUsersUsecase) {
				mockUsecase.EXPECT().ValidateAndGetUser(gomock.Any(), "invalid-jwt-token").
					Return(models.User{}, errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "CSRF token from form",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("csrftoken=csrf-token"))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "csrf-token"})
				req.AddCookie(&http.Cookie{Name: CookieName, Value: "jwt-token"})
				return req.WithContext(testContext())
			},
			setupMocks: func(mockUsecase *mocks.MockUsersUsecase) {
				mockUsecase.EXPECT().ValidateAndGetUser(gomock.Any(), "jwt-token").
					Return(models.User{ID: uuid.NewV4()}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mocks.NewMockUsersUsecase(ctrl)
			defer ctrl.Finish()

			if tt.setupMocks != nil {
				tt.setupMocks(mockUsecase)
			}

			handler := NewUserHandler(mockUsecase)
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			r := tt.setupRequest()
			w := httptest.NewRecorder()

			middleware := handler.Middleware(nextHandler)
			middleware.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
