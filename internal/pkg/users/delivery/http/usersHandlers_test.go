package http

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
			name:   "Unauthorized",
			userID: userIDStr,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					GetUser(gomock.Any(), &gen.GetUserRequest{ID: userIDStr}).
					Return(nil, status.Error(codes.Unauthenticated, "unauthorized"))
			},
			expectedStatus: http.StatusUnauthorized,
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
			name:        "User not found",
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
					Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			expectedStatus: http.StatusNotFound,
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

func TestChangeAvatar(t *testing.T) {
	userID := uuid.NewV4()

	tests := []struct {
		name           string
		setupContext   func(r *http.Request) *http.Request
		setupRequest   func() *http.Request
		mockSetup      func(mockClient *mocks.MockAuthClient)
		expectedStatus int
	}{
		{
			name:           "No user in context",
			setupContext:   func(r *http.Request) *http.Request { return r },
			setupRequest:   func() *http.Request { return httptest.NewRequest("PUT", "/users/avatar", nil) },
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "No avatar file",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), users.UserKey, userID)
				return r.WithContext(ctx)
			},
			setupRequest: func() *http.Request {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.WriteField("other_field", "value")
				writer.Close()

				req := httptest.NewRequest("PUT", "/users/avatar", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "File too large",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), users.UserKey, userID)
				return r.WithContext(ctx)
			},
			setupRequest: func() *http.Request {
				// Create a request with body larger than 10MB
				largeContent := strings.Repeat("a", 11*1024*1024) // 11MB
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				io.WriteString(part, largeContent)
				writer.Close()

				req := httptest.NewRequest("PUT", "/users/avatar", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req
			},
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name: "Invalid request body",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), users.UserKey, userID)
				return r.WithContext(ctx)
			},
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("PUT", "/users/avatar", bytes.NewBufferString("invalid body"))
				req.Header.Set("Content-Type", "multipart/form-data")
				return req
			},
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusBadRequest,
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

			r := tt.setupRequest().WithContext(testContext())
			if tt.setupContext != nil {
				r = tt.setupContext(r)
			}
			w := httptest.NewRecorder()

			handler.ChangeAvatar(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestJWTMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		mockSetup      func(mockClient *mocks.MockAuthClient)
		expectedStatus int
	}{
		{
			name: "Success with cookie",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				cookie := &http.Cookie{Name: CookieName, Value: "valid_jwt_token"}
				req.AddCookie(cookie)
				return req
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					ValidateAndGetUser(gomock.Any(), &gen.ValidateAndGetUserRequest{Token: "valid_jwt_token"}).
					Return(&gen.UserResponse{
						ID:      uuid.NewV4().String(),
						Version: 1,
						Login:   "testuser",
						Avatar:  "/avatar.jpg",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Success without cookie",
			setupRequest: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					ValidateAndGetUser(gomock.Any(), &gen.ValidateAndGetUserRequest{Token: ""}).
					Return(&gen.UserResponse{
						ID:      uuid.NewV4().String(),
						Version: 1,
						Login:   "testuser",
						Avatar:  "/avatar.jpg",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Unauthenticated",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				cookie := &http.Cookie{Name: CookieName, Value: "invalid_token"}
				req.AddCookie(cookie)
				return req
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					ValidateAndGetUser(gomock.Any(), &gen.ValidateAndGetUserRequest{Token: "invalid_token"}).
					Return(nil, status.Error(codes.Unauthenticated, "invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Internal server error",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				cookie := &http.Cookie{Name: CookieName, Value: "token"}
				req.AddCookie(cookie)
				return req
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					ValidateAndGetUser(gomock.Any(), &gen.ValidateAndGetUserRequest{Token: "token"}).
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

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			r := tt.setupRequest().WithContext(testContext())
			w := httptest.NewRecorder()

			middleware := handler.JWTMiddleware(nextHandler)
			middleware.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		mockSetup      func(mockClient *mocks.MockAuthClient)
		expectedStatus int
	}{
		{
			name: "Success with header CSRF token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("POST", "/test", nil)
				req.Header.Set("X-CSRF-Token", "valid_csrf_token")
				req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "valid_csrf_token"})
				req.AddCookie(&http.Cookie{Name: CookieName, Value: "valid_jwt_token"})
				return req
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					ValidateAndGetUser(gomock.Any(), &gen.ValidateAndGetUserRequest{Token: "valid_jwt_token"}).
					Return(&gen.UserResponse{
						ID:      uuid.NewV4().String(),
						Version: 1,
						Login:   "testuser",
						Avatar:  "/avatar.jpg",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Success with form CSRF token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("POST", "/test", nil)
				req.Form = map[string][]string{"csrftoken": {"valid_csrf_token"}}
				req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "valid_csrf_token"})
				req.AddCookie(&http.Cookie{Name: CookieName, Value: "valid_jwt_token"})
				return req
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					ValidateAndGetUser(gomock.Any(), &gen.ValidateAndGetUserRequest{Token: "valid_jwt_token"}).
					Return(&gen.UserResponse{
						ID:      uuid.NewV4().String(),
						Version: 1,
						Login:   "testuser",
						Avatar:  "/avatar.jpg",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "No CSRF cookie",
			setupRequest: func() *http.Request {
				return httptest.NewRequest("POST", "/test", nil)
			},
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "No CSRF token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("POST", "/test", nil)
				req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "valid_csrf_token"})
				return req
			},
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusUnauthorized,
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

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			r := tt.setupRequest().WithContext(testContext())
			w := httptest.NewRecorder()

			middleware := handler.Middleware(nextHandler)
			middleware.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
