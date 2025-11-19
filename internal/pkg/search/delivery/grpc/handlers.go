package grpc

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/search"
	"kinopoisk/internal/pkg/search/delivery/grpc/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcSearchHandler struct {
	uc search.SearchUsecase
	gen.UnimplementedSearchServer
}

func NewGrpcSearchHandler(uc search.SearchUsecase) *GrpcSearchHandler {
	return &GrpcSearchHandler{uc: uc}
}

func (g GrpcSearchHandler) SearchFilmsAndActors(ctx context.Context, in *gen.SearchFilmsAndActorsRequest) (*gen.SearchFilmsAndActorsResponse, error) {
	var filmsResult []*gen.MainPageFilm
	filmsPager := models.Pager{
		Count:  int(in.FilmsPager.Count),
		Offset: int(in.FilmsPager.Offset),
	}
	mainPageFilms, err := g.uc.GetFilmsFromSearch(ctx, in.SearchString, filmsPager)
	if err != nil {
		switch {
		case errors.Is(err, films.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "films not found")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	for i := range mainPageFilms {
		mainPageFilms[i].Sanitize()
		filmsResult = append(filmsResult, &gen.MainPageFilm{
			ID:     mainPageFilms[i].ID.String(),
			Cover:  mainPageFilms[i].Cover,
			Title:  mainPageFilms[i].Title,
			Rating: mainPageFilms[i].Rating,
			Year:   int32(mainPageFilms[i].Year),
			Genre:  mainPageFilms[i].Genre,
		})
	}

	var actorsResult []*gen.MainPageActor
	actorsPager := models.Pager{
		Count:  int(in.ActorsPager.Count),
		Offset: int(in.ActorsPager.Offset),
	}
	mainPageActors, err := g.uc.GetActorsFromSearch(ctx, in.SearchString, actorsPager)
	if err != nil {
		switch {
		case errors.Is(err, actors.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "films not found")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	for i := range mainPageActors {
		mainPageActors[i].Sanitize()
		actorsResult = append(actorsResult, &gen.MainPageActor{
			ID:          mainPageActors[i].ID.String(),
			RussianName: mainPageActors[i].RussianName,
			Photo:       mainPageActors[i].Photo,
		})
	}

	return &gen.SearchFilmsAndActorsResponse{
		Actors: actorsResult,
		Films:  filmsResult,
	}, nil
}
