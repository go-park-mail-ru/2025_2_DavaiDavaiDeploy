package http

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kinopoisk/internal/pkg/auth/delivery/grpc/gen"
	"kinopoisk/internal/pkg/auth/mocks"
	"kinopoisk/internal/pkg/middleware/logger"
	"kinopoisk/internal/pkg/users"

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

func TestGetUser(t *testing.T) {
	userID := uuid.NewV4()
	userIDStr := userID.String()

	tests := []struct {
		name           string
		userID         string
		mockSetup      func(mockClient *mocks.MockAuthClient)
		expectedStatus int
	}{
		{
			name:   "Success",
			userID: userIDStr,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					GetUser(gomock.Any(), &gen.GetUserRequest{ID: userIDStr}).
					Return(&gen.UserResponse{
						ID:      userIDStr,
						Version: 1,
						Login:   "testuser",
						Avatar:  "/avatar.jpg",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid UUID",
			userID:         "invalid-uuid",
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "User not found",
			userID: userIDStr,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					GetUser(gomock.Any(), &gen.GetUserRequest{ID: userIDStr}).
					Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Internal server error",
			userID: userIDStr,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					GetUser(gomock.Any(), &gen.GetUserRequest{ID: userIDStr}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockAuthClient(ctrl)
			handler := NewUserHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			router := mux.NewRouter()
			router.HandleFunc("/users/{id}", handler.GetUser)

			r := httptest.NewRequest("GET", "/users/"+tt.userID, nil).WithContext(testContext())
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestChangePassword(t *testing.T) {
	userID := uuid.NewV4()
	userIDStr := userID.String()

	tests := []struct {
		name           string
		requestBody    string
		setupContext   func(r *http.Request) *http.Request
		mockSetup      func(mockClient *mocks.MockAuthClient)
		expectedStatus int
	}{
		{
			name:        "Success",
			requestBody: `{"old_password":"oldPass123","new_password":"newPass123"}`,
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), users.UserKey, userID)
				return r.WithContext(ctx)
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					ChangePassword(gomock.Any(), &gen.ChangePasswordRequest{
						UserID:      userIDStr,
						OldPassword: "oldPass123",
						NewPassword: "newPass123",
					}).
					Return(&gen.AuthResponse{
						User: &gen.UserResponse{
							ID:      userIDStr,
							Version: 1,
							Login:   "testuser",
							Avatar:  "/avatar.jpg",
						},
						JWTToken:  "new_jwt_token",
						CSRFToken: "new_csrf_token",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Invalid JSON",
			requestBody: `{"old_password":"oldPass123","new_password":"newPass123"`,
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), users.UserKey, userID)
				return r.WithContext(ctx)
			},
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "No user in context",
			requestBody:    `{"old_password":"oldPass123","new_password":"newPass123"}`,
			setupContext:   func(r *http.Request) *http.Request { return r },
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:        "Wrong old password",
			requestBody: `{"old_password":"wrongOldPass","new_password":"newPass123"}`,
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), users.UserKey, userID)
				return r.WithContext(ctx)
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					ChangePassword(gomock.Any(), &gen.ChangePasswordRequest{
						UserID:      userIDStr,
						OldPassword: "wrongOldPass",
						NewPassword: "newPass123",
					}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid old password"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Internal server error",
			requestBody: `{"old_password":"oldPass123","new_password":"newPass123"}`,
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), users.UserKey, userID)
				return r.WithContext(ctx)
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					ChangePassword(gomock.Any(), &gen.ChangePasswordRequest{
						UserID:      userIDStr,
						OldPassword: "oldPass123",
						NewPassword: "newPass123",
					}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockAuthClient(ctrl)
			handler := NewUserHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			r := httptest.NewRequest("PUT", "/users/password", bytes.NewBufferString(tt.requestBody)).WithContext(testContext())
			if tt.setupContext != nil {
				r = tt.setupContext(r)
			}
			w := httptest.NewRecorder()

			handler.ChangePassword(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
