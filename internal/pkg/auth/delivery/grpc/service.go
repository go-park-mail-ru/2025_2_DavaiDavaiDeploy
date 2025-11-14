package grpc

import (
	"context"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/auth/delivery/grpc/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcAuthHandler struct {
	uc auth.AuthUsecase
	gen.UnimplementedAuthServer
}

func NewGrpcAuthHandler(uc auth.AuthUsecase) *GrpcAuthHandler {
	return &GrpcAuthHandler{uc: uc}
}

func (g GrpcAuthHandler) SignupUser(context.Context, *gen.SignupRequest) (*gen.AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignupUser not implemented")
}
func (g GrpcAuthHandler) SignInUser(context.Context, *gen.SignInRequest) (*gen.AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignInUser not implemented")
}
func (g GrpcAuthHandler) CheckAuth(context.Context, *gen.CheckAuthRequest) (*gen.UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAuth not implemented")
}
func (g GrpcAuthHandler) LogOutUser(context.Context, *gen.LogOutUserRequest) (*gen.LogOutUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogOutUser not implemented")
}
func (g GrpcAuthHandler) GetUser(context.Context, *gen.GetUserRequest) (*gen.UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (g GrpcAuthHandler) ChangePassword(context.Context, *gen.ChangePasswordRequest) (*gen.AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (g GrpcAuthHandler) ChangeAvatar(context.Context, *gen.ChangeAvatarRequest) (*gen.AuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeAvatar not implemented")
}
