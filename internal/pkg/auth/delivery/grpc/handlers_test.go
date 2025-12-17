package grpc

import (
	"context"
	"testing"

	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/auth/delivery/grpc/gen"
	auth_mocks "kinopoisk/internal/pkg/auth/mocks"
	"kinopoisk/internal/pkg/users"
	users_mocks "kinopoisk/internal/pkg/users/mocks"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcAuthHandler_SignupUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := auth_mocks.NewMockAuthUsecase(ctrl)
	mockUsersUsecase := users_mocks.NewMockUsersUsecase(ctrl)
	handler := NewGrpcAuthHandler(mockAuthUsecase, mockUsersUsecase)

	userID := uuid.NewV4()
	login := "testuser"
	password := "testpass123"

	tests := []struct {
		name           string
		input          *gen.SignupRequest
		mockSetup      func()
		expected       *gen.AuthResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.SignupRequest{
				Login:    login,
				Password: password,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().SignUpUser(gomock.Any(), models.SignUpInput{
					Login:    login,
					Password: password,
				}).Return(models.User{
					ID:      userID,
					Login:   login,
					Avatar:  "avatars/default.png",
					Version: 1,
				}, "jwt-token", nil)
			},
			expected: &gen.AuthResponse{
				User: &gen.UserResponse{
					ID:      userID.String(),
					Login:   login,
					Avatar:  "avatars/default.png",
					Version: 1,
				},
				JWTToken:  "jwt-token",
				CSRFToken: gomock.Any().String(),
			},
			expectedErr: nil,
		},
		{
			name: "Error - Bad Request",
			input: &gen.SignupRequest{
				Login:    "usr",
				Password: password,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().SignUpUser(gomock.Any(), gomock.Any()).
					Return(models.User{}, "", auth.ErrorBadRequest)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "%v", auth.ErrorBadRequest),
			expectedStatus: codes.InvalidArgument,
		},
		{
			name: "Error - Conflict",
			input: &gen.SignupRequest{
				Login:    login,
				Password: password,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().SignUpUser(gomock.Any(), gomock.Any()).
					Return(models.User{}, "", auth.ErrorConflict)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.AlreadyExists, "%v", auth.ErrorConflict),
			expectedStatus: codes.AlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.SignupUser(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.User.ID, resp.User.ID)
				assert.Equal(t, tt.expected.User.Login, resp.User.Login)
				assert.Equal(t, tt.expected.User.Avatar, resp.User.Avatar)
				assert.Equal(t, tt.expected.JWTToken, resp.JWTToken)
				assert.NotEmpty(t, resp.CSRFToken)
			}
		})
	}
}

func TestGrpcAuthHandler_SignInUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := auth_mocks.NewMockAuthUsecase(ctrl)
	mockUsersUsecase := users_mocks.NewMockUsersUsecase(ctrl)
	handler := NewGrpcAuthHandler(mockAuthUsecase, mockUsersUsecase)

	userID := uuid.NewV4()
	login := "testuser"
	password := "testpass123"
	twoFactorCode := "123456"

	tests := []struct {
		name           string
		input          *gen.SignInRequest
		mockSetup      func()
		expected       *gen.AuthResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success without 2FA",
			input: &gen.SignInRequest{
				Login:         login,
				Password:      password,
				TwoFactorCode: nil,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().SignInUser(gomock.Any(), models.SignInInput{
					Login:    login,
					Password: password,
					Code:     nil,
				}).Return(models.User{
					ID:     userID,
					Login:  login,
					Avatar: "avatars/default.png",
					Has2FA: false,
				}, "jwt-token", nil)
			},
			expected: &gen.AuthResponse{
				User: &gen.UserResponse{
					ID:     userID.String(),
					Login:  login,
					Avatar: "avatars/default.png",
					Has2Fa: false,
				},
				JWTToken:  "jwt-token",
				CSRFToken: gomock.Any().String(),
			},
			expectedErr: nil,
		},
		{
			name: "Success with 2FA",
			input: &gen.SignInRequest{
				Login:         login,
				Password:      password,
				TwoFactorCode: &twoFactorCode,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().SignInUser(gomock.Any(), models.SignInInput{
					Login:    login,
					Password: password,
					Code:     &twoFactorCode,
				}).Return(models.User{
					ID:     userID,
					Login:  login,
					Avatar: "avatars/default.png",
					Has2FA: true,
				}, "jwt-token", nil)
			},
			expected: &gen.AuthResponse{
				User: &gen.UserResponse{
					ID:     userID.String(),
					Login:  login,
					Avatar: "avatars/default.png",
					Has2Fa: true,
				},
				JWTToken:  "jwt-token",
				CSRFToken: gomock.Any().String(),
			},
			expectedErr: nil,
		},
		{
			name: "Error - Bad Request",
			input: &gen.SignInRequest{
				Login:    login,
				Password: "wrongpass",
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().SignInUser(gomock.Any(), gomock.Any()).
					Return(models.User{}, "", auth.ErrorBadRequest)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "%v", auth.ErrorBadRequest),
			expectedStatus: codes.InvalidArgument,
		},
		{
			name: "Error - Precondition Failed",
			input: &gen.SignInRequest{
				Login:    login,
				Password: password,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().SignInUser(gomock.Any(), gomock.Any()).
					Return(models.User{}, "", auth.ErrorPreconditionFailed)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.FailedPrecondition, "%v", auth.ErrorPreconditionFailed),
			expectedStatus: codes.FailedPrecondition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.SignInUser(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.User.ID, resp.User.ID)
				assert.Equal(t, tt.expected.User.Login, resp.User.Login)
				assert.Equal(t, tt.expected.User.Has2Fa, resp.User.Has2Fa)
				assert.Equal(t, tt.expected.JWTToken, resp.JWTToken)
				assert.NotEmpty(t, resp.CSRFToken)
			}
		})
	}
}

func TestGrpcAuthHandler_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := auth_mocks.NewMockAuthUsecase(ctrl)
	mockUsersUsecase := users_mocks.NewMockUsersUsecase(ctrl)
	handler := NewGrpcAuthHandler(mockAuthUsecase, mockUsersUsecase)

	userID := uuid.NewV4()
	login := "testuser"

	tests := []struct {
		name           string
		input          *gen.GetUserRequest
		mockSetup      func()
		expected       *gen.UserResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.GetUserRequest{
				ID: userID.String(),
			},
			mockSetup: func() {
				mockUsersUsecase.EXPECT().GetUser(gomock.Any(), userID).
					Return(models.User{
						ID:      userID,
						Login:   login,
						Avatar:  "avatars/default.png",
						Version: 1,
					}, nil)
			},
			expected: &gen.UserResponse{
				ID:      userID.String(),
				Login:   login,
				Avatar:  "avatars/default.png",
				Version: 1,
			},
			expectedErr: nil,
		},
		{
			name: "Error - Not Found",
			input: &gen.GetUserRequest{
				ID: userID.String(),
			},
			mockSetup: func() {
				mockUsersUsecase.EXPECT().GetUser(gomock.Any(), userID).
					Return(models.User{}, users.ErrorNotFound)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.NotFound, "%v", users.ErrorNotFound),
			expectedStatus: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.GetUser(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.ID, resp.ID)
				assert.Equal(t, tt.expected.Login, resp.Login)
				assert.Equal(t, tt.expected.Avatar, resp.Avatar)
			}
		})
	}
}

func TestGrpcAuthHandler_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := auth_mocks.NewMockAuthUsecase(ctrl)
	mockUsersUsecase := users_mocks.NewMockUsersUsecase(ctrl)
	handler := NewGrpcAuthHandler(mockAuthUsecase, mockUsersUsecase)

	userID := uuid.NewV4()
	login := "testuser"

	tests := []struct {
		name           string
		input          *gen.ChangePasswordRequest
		mockSetup      func()
		expected       *gen.AuthResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.ChangePasswordRequest{
				UserID:      userID.String(),
				OldPassword: "oldpass",
				NewPassword: "newpass",
			},
			mockSetup: func() {
				mockUsersUsecase.EXPECT().ChangePassword(gomock.Any(), userID, "oldpass", "newpass").
					Return(models.User{
						ID:      userID,
						Login:   login,
						Avatar:  "avatars/default.png",
						Version: 2,
					}, "new-jwt-token", nil)
			},
			expected: &gen.AuthResponse{
				User: &gen.UserResponse{
					ID:      userID.String(),
					Login:   login,
					Avatar:  "avatars/default.png",
					Version: 2,
				},
				JWTToken:  "new-jwt-token",
				CSRFToken: gomock.Any().String(),
			},
			expectedErr: nil,
		},
		{
			name: "Error - Bad Request",
			input: &gen.ChangePasswordRequest{
				UserID:      userID.String(),
				OldPassword: "wrongpass",
				NewPassword: "newpass",
			},
			mockSetup: func() {
				mockUsersUsecase.EXPECT().ChangePassword(gomock.Any(), userID, "wrongpass", "newpass").
					Return(models.User{}, "", users.ErrorBadRequest)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "%v", users.ErrorBadRequest),
			expectedStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.ChangePassword(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.User.ID, resp.User.ID)
				assert.Equal(t, tt.expected.JWTToken, resp.JWTToken)
				assert.NotEmpty(t, resp.CSRFToken)
			}
		})
	}
}

func TestGrpcAuthHandler_ValidateAndGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := auth_mocks.NewMockAuthUsecase(ctrl)
	mockUsersUsecase := users_mocks.NewMockUsersUsecase(ctrl)
	handler := NewGrpcAuthHandler(mockAuthUsecase, mockUsersUsecase)

	userID := uuid.NewV4()
	login := "testuser"
	token := "valid-jwt-token"

	tests := []struct {
		name           string
		input          *gen.ValidateAndGetUserRequest
		mockSetup      func()
		expected       *gen.UserResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.ValidateAndGetUserRequest{
				Token: token,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().ValidateAndGetUser(gomock.Any(), token).
					Return(models.User{
						ID:     userID,
						Login:  login,
						Avatar: "avatars/default.png",
						Has2FA: true,
					}, nil)
			},
			expected: &gen.UserResponse{
				ID:     userID.String(),
				Login:  login,
				Avatar: "avatars/default.png",
				Has2Fa: true,
			},
			expectedErr: nil,
		},
		{
			name: "Error - Unauthorized",
			input: &gen.ValidateAndGetUserRequest{
				Token: "invalid-token",
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().ValidateAndGetUser(gomock.Any(), "invalid-token").
					Return(models.User{}, users.ErrorUnauthorized)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.Unauthenticated, "%v", users.ErrorUnauthorized),
			expectedStatus: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.ValidateAndGetUser(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.ID, resp.ID)
				assert.Equal(t, tt.expected.Login, resp.Login)
				assert.Equal(t, tt.expected.Has2Fa, resp.Has2Fa)
			}
		})
	}
}

func TestGrpcAuthHandler_Enable2Fa(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := auth_mocks.NewMockAuthUsecase(ctrl)
	mockUsersUsecase := users_mocks.NewMockUsersUsecase(ctrl)
	handler := NewGrpcAuthHandler(mockAuthUsecase, mockUsersUsecase)

	userID := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.Enable2FaRequest
		mockSetup      func()
		expected       *gen.Enable2FaResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.Enable2FaRequest{
				ID:     userID.String(),
				Has2Fa: false,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().Enable2FA(gomock.Any(), userID, false).
					Return(models.EnableTwoFactorResponse{
						Has2FA: true,
						QrCode: []byte("qr-code-bytes"),
					}, nil)
			},
			expected: &gen.Enable2FaResponse{
				Has2Fa: true,
				QrCode: []byte("qr-code-bytes"),
			},
			expectedErr: nil,
		},
		{
			name: "Error - Bad Request",
			input: &gen.Enable2FaRequest{
				ID:     userID.String(),
				Has2Fa: true,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().Enable2FA(gomock.Any(), userID, true).
					Return(models.EnableTwoFactorResponse{}, users.ErrorBadRequest)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "%v", users.ErrorBadRequest),
			expectedStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.Enable2Fa(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.Has2Fa, resp.Has2Fa)
				assert.Equal(t, tt.expected.QrCode, resp.QrCode)
			}
		})
	}
}

func TestGrpcAuthHandler_Disable2Fa(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := auth_mocks.NewMockAuthUsecase(ctrl)
	mockUsersUsecase := users_mocks.NewMockUsersUsecase(ctrl)
	handler := NewGrpcAuthHandler(mockAuthUsecase, mockUsersUsecase)

	userID := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.Disable2FaRequest
		mockSetup      func()
		expected       *gen.Disable2FaResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.Disable2FaRequest{
				ID:     userID.String(),
				Has2Fa: true,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().Disable2FA(gomock.Any(), userID, true).
					Return(models.DisableTwoFactorResponse{
						Has2FA: false,
					}, nil)
			},
			expected: &gen.Disable2FaResponse{
				Has2Fa: false,
			},
			expectedErr: nil,
		},
		{
			name: "Error - Bad Request",
			input: &gen.Disable2FaRequest{
				ID:     userID.String(),
				Has2Fa: false,
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().Disable2FA(gomock.Any(), userID, false).
					Return(models.DisableTwoFactorResponse{}, users.ErrorBadRequest)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.InvalidArgument, "%v", users.ErrorBadRequest),
			expectedStatus: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.Disable2Fa(context.Background(), tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expected.Has2Fa, resp.Has2Fa)
			}
		})
	}
}

func TestGrpcAuthHandler_LogOutUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUsecase := auth_mocks.NewMockAuthUsecase(ctrl)
	mockUsersUsecase := users_mocks.NewMockUsersUsecase(ctrl)
	handler := NewGrpcAuthHandler(mockAuthUsecase, mockUsersUsecase)

	userID := uuid.NewV4()

	tests := []struct {
		name           string
		input          *gen.LogOutUserRequest
		mockSetup      func()
		expected       *gen.LogOutUserResponse
		expectedErr    error
		expectedStatus codes.Code
	}{
		{
			name: "Success",
			input: &gen.LogOutUserRequest{
				ID: userID.String(),
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().LogOutUser(gomock.Any(), userID).
					Return(nil)
			},
			expected:    &gen.LogOutUserResponse{},
			expectedErr: nil,
		},
		{
			name: "Error - Unauthorized",
			input: &gen.LogOutUserRequest{
				ID: userID.String(),
			},
			mockSetup: func() {
				mockAuthUsecase.EXPECT().LogOutUser(gomock.Any(), userID).
					Return(users.ErrorUnauthorized)
			},
			expected:       nil,
			expectedErr:    status.Errorf(codes.Unauthenticated, "%v", users.ErrorUnauthorized),
			expectedStatus: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := handler.LogOutUser(context.Background(), tt.input)

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
