package authHandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/auth/delivery/grpc/gen"
	"kinopoisk/internal/pkg/auth/mocks"
	"kinopoisk/internal/pkg/middleware/logger"

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

func TestSignupUser(t *testing.T) {
	userID := uuid.NewV4()
	userIDStr := userID.String()

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(mockClient *mocks.MockAuthClient)
		expectedStatus int
		expectBody     bool
		expectedUser   models.User
	}{
		{
			name:        "Success",
			requestBody: `{"login":"test123","password":"Pass123"}`,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					SignupUser(gomock.Any(), &gen.SignupRequest{
						Login:    "test123",
						Password: "Pass123",
					}).
					Return(&gen.AuthResponse{
						User: &gen.UserResponse{
							ID:      userIDStr,
							Version: 1,
							Login:   "test123",
							Avatar:  "",
							Has2Fa:  false,
						},
						JWTToken:  "jwt_token",
						CSRFToken: "csrf_token",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedUser: models.User{
				ID:      userID,
				Version: 1,
				Login:   "test123",
				Avatar:  "",
				Has2FA:  false,
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"login":"testuser","password":"abc123"`,
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:        "User already exists",
			requestBody: `{"login":"testuser","password":"Pass123"}`,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					SignupUser(gomock.Any(), &gen.SignupRequest{
						Login:    "testuser",
						Password: "Pass123",
					}).
					Return(nil, status.Error(codes.AlreadyExists, "user already exists"))
			},
			expectedStatus: http.StatusConflict,
			expectBody:     false,
		},
		{
			name:        "Bad request - invalid login",
			requestBody: `{"login":"usr","password":"Pass123"}`,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					SignupUser(gomock.Any(), &gen.SignupRequest{
						Login:    "usr",
						Password: "Pass123",
					}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid login"))
			},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:        "Internal server error",
			requestBody: `{"login":"testuser","password":"Pass123"}`,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					SignupUser(gomock.Any(), &gen.SignupRequest{
						Login:    "testuser",
						Password: "Pass123",
					}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockAuthClient(ctrl)
			handler := NewAuthHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			r := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(tt.requestBody)).WithContext(testContext())
			w := httptest.NewRecorder()

			handler.SignupUser(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectBody {
				var decoded models.User
				err := json.Unmarshal(w.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, decoded.ID)
				assert.Equal(t, tt.expectedUser.Login, decoded.Login)
				assert.Equal(t, tt.expectedUser.Version, decoded.Version)
				assert.Equal(t, tt.expectedUser.Avatar, decoded.Avatar)
				assert.Equal(t, tt.expectedUser.Has2FA, decoded.Has2FA)
			}
		})
	}
}

func TestSignInUser(t *testing.T) {
	userID := uuid.NewV4()
	userIDStr := userID.String()

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(mockClient *mocks.MockAuthClient)
		expectedStatus int
		expectBody     bool
		expectedUser   models.User
	}{
		{
			name:        "Success",
			requestBody: `{"login":"test123","password":"Pass123"}`,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					SignInUser(gomock.Any(), &gen.SignInRequest{
						Login:    "test123",
						Password: "Pass123",
					}).
					Return(&gen.AuthResponse{
						User: &gen.UserResponse{
							ID:      userIDStr,
							Version: 1,
							Login:   "test123",
							Avatar:  "/avatar.jpg",
							Has2Fa:  false,
						},
						JWTToken:  "jwt_token",
						CSRFToken: "csrf_token",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedUser: models.User{
				ID:      userID,
				Version: 1,
				Login:   "test123",
				Avatar:  "/avatar.jpg",
				Has2FA:  false,
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"login":"testuser","password":"abc123"`,
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:        "Wrong credentials",
			requestBody: `{"login":"testuser","password":"wrongpass"}`,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					SignInUser(gomock.Any(), &gen.SignInRequest{
						Login:    "testuser",
						Password: "wrongpass",
					}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid credentials"))
			},
			expectedStatus: http.StatusBadRequest,
			expectBody:     false,
		},
		{
			name:        "2FA required but not provided",
			requestBody: `{"login":"testuser","password":"Pass123"}`,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					SignInUser(gomock.Any(), &gen.SignInRequest{
						Login:    "testuser",
						Password: "Pass123",
					}).
					Return(nil, status.Error(codes.FailedPrecondition, "2FA required"))
			},
			expectedStatus: http.StatusPreconditionFailed,
			expectBody:     false,
		},
		{
			name:        "Internal server error",
			requestBody: `{"login":"testuser","password":"Pass123"}`,
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					SignInUser(gomock.Any(), &gen.SignInRequest{
						Login:    "testuser",
						Password: "Pass123",
					}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockAuthClient(ctrl)
			handler := NewAuthHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			r := httptest.NewRequest("POST", "/auth/signin", bytes.NewBufferString(tt.requestBody)).WithContext(testContext())
			w := httptest.NewRecorder()

			handler.SignInUser(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectBody {
				var decoded models.User
				err := json.Unmarshal(w.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, decoded.ID)
				assert.Equal(t, tt.expectedUser.Login, decoded.Login)
				assert.Equal(t, tt.expectedUser.Version, decoded.Version)
				assert.Equal(t, tt.expectedUser.Avatar, decoded.Avatar)
				assert.Equal(t, tt.expectedUser.Has2FA, decoded.Has2FA)
			}
		})
	}
}

func TestCheckAuth(t *testing.T) {
	userID := uuid.NewV4()

	tests := []struct {
		name           string
		setupContext   func(r *http.Request) *http.Request
		expectedStatus int
		expectBody     bool
		expectedUser   models.User
	}{
		{
			name: "Success",
			setupContext: func(r *http.Request) *http.Request {
				user := models.User{
					ID:      userID,
					Version: 1,
					Login:   "testuser",
					Avatar:  "/avatar.jpg",
					Has2FA:  false,
				}
				ctx := context.WithValue(r.Context(), auth.UserKey, user)
				return r.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
			expectedUser: models.User{
				ID:      userID,
				Version: 1,
				Login:   "testuser",
				Avatar:  "/avatar.jpg",
				Has2FA:  false,
			},
		},
		{
			name: "Unauthorized - no user in context",
			setupContext: func(r *http.Request) *http.Request {
				return r
			},
			expectedStatus: http.StatusUnauthorized,
			expectBody:     false,
		},
		{
			name: "Unauthorized - wrong type in context",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), auth.UserKey, "not-a-user")
				return r.WithContext(ctx)
			},
			expectedStatus: http.StatusUnauthorized,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockAuthClient(ctrl)
			handler := NewAuthHandler(mockClient)

			r := httptest.NewRequest("GET", "/auth/check", nil).WithContext(testContext())
			if tt.setupContext != nil {
				r = tt.setupContext(r)
			}
			w := httptest.NewRecorder()

			handler.CheckAuth(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectBody {
				var decoded models.User
				err := json.Unmarshal(w.Body.Bytes(), &decoded)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, decoded.ID)
				assert.Equal(t, tt.expectedUser.Login, decoded.Login)
				assert.Equal(t, tt.expectedUser.Version, decoded.Version)
				assert.Equal(t, tt.expectedUser.Avatar, decoded.Avatar)
				assert.Equal(t, tt.expectedUser.Has2FA, decoded.Has2FA)
			}
		})
	}
}

func TestEnable2FA(t *testing.T) {
	userID := uuid.NewV4()
	userIDStr := userID.String()

	tests := []struct {
		name           string
		setupContext   func(r *http.Request) *http.Request
		mockSetup      func(mockClient *mocks.MockAuthClient)
		expectedStatus int
	}{
		{
			name: "Success",
			setupContext: func(r *http.Request) *http.Request {
				user := models.User{
					ID:      userID,
					Version: 1,
					Login:   "testuser",
					Has2FA:  false,
				}
				ctx := context.WithValue(r.Context(), auth.UserKey, user)
				return r.WithContext(ctx)
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					Enable2Fa(gomock.Any(), &gen.Enable2FaRequest{
						ID:     userIDStr,
						Has2Fa: false,
					}).
					Return(&gen.Enable2FaResponse{
						QrCode: []byte("qr_code_data"),
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Unauthorized",
			setupContext: func(r *http.Request) *http.Request {
				return r
			},
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Bad request",
			setupContext: func(r *http.Request) *http.Request {
				user := models.User{
					ID:      userID,
					Version: 1,
					Login:   "testuser",
					Has2FA:  false,
				}
				ctx := context.WithValue(r.Context(), auth.UserKey, user)
				return r.WithContext(ctx)
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					Enable2Fa(gomock.Any(), &gen.Enable2FaRequest{
						ID:     userIDStr,
						Has2Fa: false,
					}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid request"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Internal server error",
			setupContext: func(r *http.Request) *http.Request {
				user := models.User{
					ID:      userID,
					Version: 1,
					Login:   "testuser",
					Has2FA:  false,
				}
				ctx := context.WithValue(r.Context(), auth.UserKey, user)
				return r.WithContext(ctx)
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					Enable2Fa(gomock.Any(), &gen.Enable2FaRequest{
						ID:     userIDStr,
						Has2Fa: false,
					}).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockAuthClient(ctrl)
			handler := NewAuthHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			r := httptest.NewRequest("POST", "/auth/enable2fa", nil).WithContext(testContext())
			if tt.setupContext != nil {
				r = tt.setupContext(r)
			}
			w := httptest.NewRecorder()

			handler.Enable2FA(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestLogOutUser(t *testing.T) {
	userID := uuid.NewV4()
	userIDStr := userID.String()

	tests := []struct {
		name           string
		setupContext   func(r *http.Request) *http.Request
		mockSetup      func(mockClient *mocks.MockAuthClient)
		expectedStatus int
	}{
		{
			name: "Success",
			setupContext: func(r *http.Request) *http.Request {
				user := models.User{
					ID:      userID,
					Version: 1,
					Login:   "testuser",
					Avatar:  "/avatar.jpg",
				}
				ctx := context.WithValue(r.Context(), auth.UserKey, user)
				return r.WithContext(ctx)
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					LogOutUser(gomock.Any(), &gen.LogOutUserRequest{
						ID:      userIDStr,
						Version: 1,
						Login:   "testuser",
						Avatar:  "/avatar.jpg",
					}).
					Return(&gen.LogOutUserResponse{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Unauthorized",
			setupContext: func(r *http.Request) *http.Request {
				return r
			},
			mockSetup:      func(mockClient *mocks.MockAuthClient) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Internal server error",
			setupContext: func(r *http.Request) *http.Request {
				user := models.User{
					ID:      userID,
					Version: 1,
					Login:   "testuser",
					Avatar:  "/avatar.jpg",
				}
				ctx := context.WithValue(r.Context(), auth.UserKey, user)
				return r.WithContext(ctx)
			},
			mockSetup: func(mockClient *mocks.MockAuthClient) {
				mockClient.EXPECT().
					LogOutUser(gomock.Any(), &gen.LogOutUserRequest{
						ID:      userIDStr,
						Version: 1,
						Login:   "testuser",
						Avatar:  "/avatar.jpg",
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
			handler := NewAuthHandler(mockClient)

			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			r := httptest.NewRequest("POST", "/auth/logout", nil).WithContext(testContext())
			if tt.setupContext != nil {
				r = tt.setupContext(r)
			}
			w := httptest.NewRecorder()

			handler.LogOutUser(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestNewAuthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockAuthClient(ctrl)

	t.Run("Success creation", func(t *testing.T) {
		handler := NewAuthHandler(mockClient)
		assert.NotNil(t, handler)
		assert.Equal(t, mockClient, handler.client)
	})

	t.Run("Creation with nil client", func(t *testing.T) {
		handler := NewAuthHandler(nil)
		assert.NotNil(t, handler)
		assert.Nil(t, handler.client)
	})
}
