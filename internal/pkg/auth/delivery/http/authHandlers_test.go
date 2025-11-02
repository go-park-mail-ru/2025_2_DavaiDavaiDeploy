package authHandlers

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
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/auth/mocks"
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

func TestSignupUser(t *testing.T) {
	type args struct {
		login    string
		password string
	}

	tests := []struct {
		name           string
		requestBody    string
		args           args
		ucErr          error
		expectedStatus int
	}{
		{
			name:        "Success",
			requestBody: `{"login":"test123","password":"Pass123"}`,
			args: args{
				login:    "test123",
				password: "Pass123",
			},
			ucErr:          nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"login":"testuser","password":"abc123"`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "User already exists",
			requestBody: `{"login":"testuser","password":"Pass123"}`,
			args: args{
				login:    "testuser",
				password: "Pass123",
			},
			ucErr:          auth.ErrorConflict,
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "Bad request - invalid login",
			requestBody: `{"login":"usr","password":"Pass123"}`,
			args: args{
				login:    "usr",
				password: "Pass123",
			},
			ucErr:          auth.ErrorBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Internal server error",
			requestBody: `{"login":"testuser","password":"Pass123"}`,
			args: args{
				login:    "testuser",
				password: "Pass123",
			},
			ucErr:          errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mocks.NewMockAuthUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "Invalid JSON" {
				mockUsecase.EXPECT().SignUpUser(gomock.Any(), models.SignUpInput{
					Login:    tt.args.login,
					Password: tt.args.password,
				}).Return(models.User{
					ID:    uuid.NewV4(),
					Login: tt.args.login,
				}, "jwt_token", tt.ucErr)
			}

			r := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(tt.requestBody)).WithContext(testContext())
			w := httptest.NewRecorder()

			handler := NewAuthHandler(mockUsecase)
			handler.SignupUser(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestSignInUser(t *testing.T) {
	type args struct {
		login    string
		password string
	}

	tests := []struct {
		name           string
		requestBody    string
		args           args
		ucErr          error
		expectedStatus int
	}{
		{
			name:        "Success",
			requestBody: `{"login":"test123","password":"Pass123"}`,
			args: args{
				login:    "test123",
				password: "Pass123",
			},
			ucErr:          nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"login":"testuser","password":"abc123"`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Wrong credentials",
			requestBody: `{"login":"testuser","password":"wrongpass"}`,
			args: args{
				login:    "testuser",
				password: "wrongpass",
			},
			ucErr:          auth.ErrorBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Internal server error",
			requestBody: `{"login":"testuser","password":"Pass123"}`,
			args: args{
				login:    "testuser",
				password: "Pass123",
			},
			ucErr:          errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mocks.NewMockAuthUsecase(ctrl)
			defer ctrl.Finish()

			if tt.name != "Invalid JSON" {
				mockUsecase.EXPECT().SignInUser(gomock.Any(), models.SignInInput{
					Login:    tt.args.login,
					Password: tt.args.password,
				}).Return(models.User{
					ID:    uuid.NewV4(),
					Login: tt.args.login,
				}, "jwt_token", tt.ucErr)
			}

			r := httptest.NewRequest("POST", "/auth/signin", bytes.NewBufferString(tt.requestBody)).WithContext(testContext())
			w := httptest.NewRecorder()

			handler := NewAuthHandler(mockUsecase)
			handler.SignInUser(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCheckAuth(t *testing.T) {
	tests := []struct {
		name           string
		ucErr          error
		expectedStatus int
	}{
		{
			name:           "Success",
			ucErr:          nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized",
			ucErr:          auth.ErrorUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Internal server error",
			ucErr:          errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mocks.NewMockAuthUsecase(ctrl)
			defer ctrl.Finish()

			mockUsecase.EXPECT().CheckAuth(gomock.Any()).
				Return(models.User{
					ID:    uuid.NewV4(),
					Login: "testuser",
				}, tt.ucErr)

			r := httptest.NewRequest("GET", "/auth/check", nil).WithContext(testContext())
			w := httptest.NewRecorder()

			handler := NewAuthHandler(mockUsecase)
			handler.CheckAuth(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestLogOutUser(t *testing.T) {
	tests := []struct {
		name           string
		ucErr          error
		expectedStatus int
	}{
		{
			name:           "Success",
			ucErr:          nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized",
			ucErr:          auth.ErrorUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Internal server error",
			ucErr:          errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsecase := mocks.NewMockAuthUsecase(ctrl)
			defer ctrl.Finish()

			mockUsecase.EXPECT().LogOutUser(gomock.Any()).
				Return(tt.ucErr)

			r := httptest.NewRequest("POST", "/auth/logout", nil).WithContext(testContext())
			w := httptest.NewRecorder()

			handler := NewAuthHandler(mockUsecase)
			handler.LogOutUser(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
