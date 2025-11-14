package grpc

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/auth/delivery/grpc/gen"
	"kinopoisk/internal/pkg/users"
	"os"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcAuthHandler struct {
	JWTSecret string
	auc       auth.AuthUsecase
	uuc       users.UsersUsecase
	gen.UnimplementedAuthServer
}

func NewGrpcAuthHandler(auc auth.AuthUsecase, uuc users.UsersUsecase) *GrpcAuthHandler {
	return &GrpcAuthHandler{auc: auc, uuc: uuc, JWTSecret: os.Getenv("JWT_SECRET")}
}

func (g GrpcAuthHandler) SignupUser(ctx context.Context, in *gen.SignupRequest) (*gen.AuthResponse, error) {
	req := models.SignUpInput{
		Login:    in.Login,
		Password: in.Password,
	}
	req.Sanitize()

	user, token, err := g.auc.SignUpUser(ctx, req)
	if err != nil {
		switch err {
		case auth.ErrorBadRequest:
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		case auth.ErrorConflict:
			return nil, status.Errorf(codes.AlreadyExists, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	user.Sanitize()

	csrfToken := uuid.NewV4().String()

	userResponse := &gen.UserResponse{
		ID:      user.ID.String(),
		Version: int32(user.Version),
		Login:   user.Login,
		Avatar:  user.Avatar,
	}

	return &gen.AuthResponse{
		User:      userResponse,
		JWTToken:  token,
		CSRFToken: csrfToken,
	}, err
}

func (g GrpcAuthHandler) SignInUser(ctx context.Context, in *gen.SignInRequest) (*gen.AuthResponse, error) {
	req := models.SignInInput{
		Login:    in.Login,
		Password: in.Password,
	}
	req.Sanitize()

	user, token, err := g.auc.SignInUser(ctx, req)
	if err != nil {
		switch err {
		case auth.ErrorBadRequest:
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	user.Sanitize()

	userResponse := &gen.UserResponse{
		ID:      user.ID.String(),
		Version: int32(user.Version),
		Login:   user.Login,
		Avatar:  user.Avatar,
	}

	csrfToken := uuid.NewV4().String()

	return &gen.AuthResponse{
		User:      userResponse,
		JWTToken:  token,
		CSRFToken: csrfToken,
	}, err
}

func (g GrpcAuthHandler) CheckAuth(ctx context.Context, in *gen.CheckAuthRequest) (*gen.UserResponse, error) {
	user, err := g.auc.CheckAuth(ctx)
	if err != nil {
		switch err {
		case auth.ErrorUnauthorized:
			return nil, status.Errorf(codes.Unauthenticated, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	user.Sanitize()

	return &gen.UserResponse{
		ID:      user.ID.String(),
		Version: int32(user.Version),
		Login:   user.Login,
		Avatar:  user.Avatar,
	}, err
}

func (g GrpcAuthHandler) GetUser(ctx context.Context, in *gen.GetUserRequest) (*gen.UserResponse, error) {
	neededUser, err := g.uuc.GetUser(ctx, uuid.FromStringOrNil(in.ID))
	if err != nil {
		switch err {
		case users.ErrorNotFound:
			return nil, status.Errorf(codes.NotFound, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	neededUser.Sanitize()

	return &gen.UserResponse{
		ID:      neededUser.ID.String(),
		Version: int32(neededUser.Version),
		Login:   neededUser.Login,
		Avatar:  neededUser.Avatar,
	}, err
}

func (g GrpcAuthHandler) ChangePassword(ctx context.Context, in *gen.ChangePasswordRequest) (*gen.AuthResponse, error) {
	userID, ok := ctx.Value(users.UserKey).(uuid.UUID)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "%v", errors.New("no user"))
	}

	var req models.ChangePasswordInput
	req.Sanitize()

	user, token, err := g.uuc.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case users.ErrorBadRequest:
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	user.Sanitize()

	csrfToken := uuid.NewV4().String()

	userResponse := &gen.UserResponse{
		ID:      user.ID.String(),
		Version: int32(user.Version),
		Login:   user.Login,
		Avatar:  user.Avatar,
	}

	return &gen.AuthResponse{
		User:      userResponse,
		JWTToken:  token,
		CSRFToken: csrfToken,
	}, err
}

func (g GrpcAuthHandler) ChangeAvatar(ctx context.Context, in *gen.ChangeAvatarRequest) (*gen.AuthResponse, error) {
	userID, ok := ctx.Value(users.UserKey).(uuid.UUID)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "%v", errors.New("no user"))
	}

	user, token, err := g.uuc.ChangeUserAvatar(ctx, userID, in.Avatar, in.FileFormat)
	if err != nil {
		switch err {
		case users.ErrorBadRequest:
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	user.Sanitize()

	csrfToken := uuid.NewV4().String()
	userResponse := &gen.UserResponse{
		ID:      user.ID.String(),
		Version: int32(user.Version),
		Login:   user.Login,
		Avatar:  user.Avatar,
	}

	return &gen.AuthResponse{
		User:      userResponse,
		JWTToken:  token,
		CSRFToken: csrfToken,
	}, err
}

func (g GrpcAuthHandler) LogOutUser(ctx context.Context, in *gen.LogOutUserRequest) (*gen.LogOutUserResponse, error) {
	err := g.auc.LogOutUser(ctx)
	if err != nil {
		switch err {
		case users.ErrorUnauthorized:
			return nil, status.Errorf(codes.Unauthenticated, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	response := &gen.LogOutUserResponse{}
	return response, nil
}

func (g GrpcAuthHandler) ValidateAndGetUser(ctx context.Context, in *gen.ValidateAndGetUserRequest) (*gen.UserResponse, error) {
	user, err := g.auc.ValidateAndGetUser(ctx, in.Token)
	if err != nil {
		switch err {
		case users.ErrorUnauthorized:
			return nil, status.Errorf(codes.Unauthenticated, "%v", err)
		default:
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	user.Sanitize()

	return &gen.UserResponse{
		ID:      user.ID.String(),
		Version: int32(user.Version),
		Login:   user.Login,
		Avatar:  user.Avatar,
	}, err
}
