package grpc

import (
	"context"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcFilmsHandler struct {
	uc films.FilmUsecase
	gen.UnimplementedFilmsServer
}

func NewGrpcFilmHandler(uc films.FilmUsecase) *GrpcFilmsHandler {
	return &GrpcFilmsHandler{uc: uc}
}

func (g GrpcFilmsHandler) GetPromoFilm(context.Context, *gen.EmptyRequest) (*gen.GetPromoFilmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPromoFilm not implemented")
}
func (g GrpcFilmsHandler) GetFilms(context.Context, *gen.GetFilmsRequest) (*gen.GetFilmsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilms not implemented")
}
func (g GrpcFilmsHandler) GetFilm(context.Context, *gen.GetFilmRequest) (*gen.GetFilmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilm not implemented")
}
func (g GrpcFilmsHandler) GetFilmFeedbacks(context.Context, *gen.GetFilmFeedbacksRequest) (*gen.GetFilmFeedbacksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilmFeedbacks not implemented")
}
func (g GrpcFilmsHandler) SendFeedback(context.Context, *gen.SendFeedbackRequest) (*gen.SendFeedbackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFeedback not implemented")
}
func (g GrpcFilmsHandler) SetRating(context.Context, *gen.SetRatingRequest) (*gen.SetRatingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetRating not implemented")
}
func (g GrpcFilmsHandler) SiteMap(context.Context, *gen.EmptyRequest) (*gen.SiteMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SiteMap not implemented")
}
func (g GrpcFilmsHandler) GetGenre(context.Context, *gen.GetGenreRequest) (*gen.GetGenreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGenre not implemented")
}
func (g GrpcFilmsHandler) GetGenres(context.Context, *gen.GetGenresRequest) (*gen.GetGenresResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGenres not implemented")
}
func (g GrpcFilmsHandler) GetFilmsByGenre(context.Context, *gen.GetFilmsByGenreRequest) (*gen.GetFilmsByGenreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilmsByGenre not implemented")
}
func (g GrpcFilmsHandler) GetActor(context.Context, *gen.GetActorRequest) (*gen.GetActorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActor not implemented")
}
func (g GrpcFilmsHandler) GetFilmsByActor(context.Context, *gen.GetFilmsByActorRequest) (*gen.GetFilmsByActorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilmsByActor not implemented")
}
