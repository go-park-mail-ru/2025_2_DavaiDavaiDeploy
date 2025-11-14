package grpc

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	"kinopoisk/internal/pkg/genres"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcFilmsHandler struct {
	uc  films.FilmUsecase
	guc genres.GenreUsecase
	auc actors.ActorUsecase
	gen.UnimplementedFilmsServer
}

func NewGrpcFilmHandler(uc films.FilmUsecase, guc genres.GenreUsecase, auc actors.ActorUsecase) *GrpcFilmsHandler {
	return &GrpcFilmsHandler{uc: uc, guc: guc, auc: auc}
}

func (g GrpcFilmsHandler) GetPromoFilm(ctx context.Context, in *gen.EmptyRequest) (*gen.GetPromoFilmResponse, error) {
	film, err := g.uc.GetPromoFilm(ctx)
	if err != nil {
		switch {
		case errors.Is(err, films.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "film not found")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}
	film.Sanitize()
	return &gen.GetPromoFilmResponse{
		Id:               film.ID.String(),
		Image:            film.Image,
		Title:            film.Title,
		Rating:           film.Rating,
		ShortDescription: film.ShortDescription,
		Year:             int32(film.Year),
		Genre:            film.Genre,
		Duration:         int32(film.Duration),
	}, nil
}
func (g GrpcFilmsHandler) GetFilms(ctx context.Context, in *gen.GetFilmsRequest) (*gen.GetFilmsResponse, error) {
	var result []*gen.MainPageFilm
	req := models.Pager{
		Count:  int(in.Pager.Count),
		Offset: int(in.Pager.Offset),
	}
	mainPageFilms, err := g.uc.GetFilms(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, films.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "films not found")
		case errors.Is(err, films.ErrorBadRequest):
			return nil, status.Errorf(codes.Internal, "bad request")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	for i := range mainPageFilms {
		mainPageFilms[i].Sanitize()
		result = append(result, &gen.MainPageFilm{
			Id:     mainPageFilms[i].ID.String(),
			Cover:  mainPageFilms[i].Cover,
			Title:  mainPageFilms[i].Title,
			Rating: mainPageFilms[i].Rating,
			Year:   int32(mainPageFilms[i].Year),
			Genre:  mainPageFilms[i].Genre,
		})
	}

	return &gen.GetFilmsResponse{
		Films: result,
	}, nil
}
func (g GrpcFilmsHandler) GetFilm(ctx context.Context, in *gen.GetFilmRequest) (*gen.GetFilmResponse, error) {
	var actors []*gen.Actor
	id, _ := uuid.FromString(in.FilmId)
	req := models.FilmPage{
		ID: id,
	}
	film, err := g.uc.GetFilm(ctx, req.ID)
	if err != nil {
		switch {
		case errors.Is(err, films.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "films not found")
		case errors.Is(err, films.ErrorBadRequest):
			return nil, status.Errorf(codes.Internal, "bad request")
		default:
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	var userRating *int32
	if film.UserRating != nil {
		val := int32(*film.UserRating)
		userRating = &val
	}

	for i := range film.Actors {
		film.Actors[i].Sanitize()
		birthDateStr := film.Actors[i].BirthDate.String()
		var deathDateStrPtr *string
		if film.Actors[i].DeathDate != nil {
			deathDateStr := film.Actors[i].DeathDate.String()
			deathDateStrPtr = &deathDateStr
		}

		var originalName string
		if film.Actors[i].OriginalName != nil {
			originalName = *film.Actors[i].OriginalName
		}
		actors = append(actors, &gen.Actor{
			Id:            film.Actors[i].ID.String(),
			RussianName:   &film.Actors[i].RussianName,
			OriginalName:  originalName,
			Photo:         film.Actors[i].Photo,
			Height:        int32(film.Actors[i].Height),
			BirthDate:     birthDateStr,
			DeathDate:     deathDateStrPtr,
			ZodiacSign:    film.Actors[i].ZodiacSign,
			BirthPlace:    film.Actors[i].BirthPlace,
			MaritalStatus: film.Actors[i].MaritalStatus,
		})
	}

	return &gen.GetFilmResponse{
		Id:               film.ID.String(),
		Title:            film.Title,
		OriginalTitle:    film.OriginalTitle,
		Cover:            film.Cover,
		Poster:           film.Poster,
		Genre:            film.Genre,
		ShortDescription: film.ShortDescription,
		Description:      film.Description,
		AgeCategory:      film.AgeCategory,
		Budget:           int32(film.Budget),
		WorldwideFees:    int32(film.WorldwideFees),
		TrailerUrl:       film.TrailerURL,
		NumberOfRatings:  int32(film.NumberOfRatings),
		Year:             int32(film.Year),
		Rating:           film.Rating,
		Country:          film.Country,
		Slogan:           film.Slogan,
		Duration:         int32(film.Duration),
		Image1:           film.Image1,
		Image2:           film.Image2,
		Image3:           film.Image3,
		Actors:           actors,
		IsReviewed:       film.IsReviewed,
		UserRating:       userRating,
	}, nil
}
func (g GrpcFilmsHandler) GetFilmFeedbacks(ctx context.Context, in *gen.GetFilmFeedbacksRequest) (*gen.GetFilmFeedbacksResponse, error) {
	filmID, err := uuid.FromString(in.FilmId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid film ID")
	}

	pager := models.Pager{
		Count:  int(in.Pager.Count),
		Offset: int(in.Pager.Offset),
	}

	feedbacks, err := g.uc.GetFilmFeedbacks(ctx, filmID, pager)
	if err != nil {
		switch {
		case errors.Is(err, films.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "feedbacks not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get feedbacks")
		}
	}

	var result []*gen.FilmFeedback
	for i := range feedbacks {
		feedbacks[i].Sanitize()

		result = append(result, &gen.FilmFeedback{
			Id:            feedbacks[i].ID.String(),
			UserId:        feedbacks[i].UserID.String(),
			FilmId:        feedbacks[i].FilmID.String(),
			Title:         feedbacks[i].Title,
			Text:          feedbacks[i].Text,
			Rating:        int32(feedbacks[i].Rating),
			CreatedAt:     feedbacks[i].CreatedAt.String(),
			UpdatedAt:     feedbacks[i].UpdatedAt.String(),
			UserLogin:     feedbacks[i].UserLogin,
			UserAvatar:    feedbacks[i].UserAvatar,
			IsMine:        feedbacks[i].IsMine,
			NewFilmRating: feedbacks[i].NewFilmRating,
		})
	}

	return &gen.GetFilmFeedbacksResponse{
		Feedbacks: result,
	}, nil
}

func (g GrpcFilmsHandler) SendFeedback(ctx context.Context, in *gen.SendFeedbackRequest) (*gen.SendFeedbackResponse, error) {
	filmID, err := uuid.FromString(in.FilmId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid film ID")
	}

	req := models.FilmFeedbackInput{
		Title:  in.Feedback.Title,
		Text:   in.Feedback.Text,
		Rating: int(in.Feedback.Rating),
	}
	req.Sanitize()

	feedback, err := g.uc.SendFeedback(ctx, req, filmID)
	if err != nil {
		switch {
		case errors.Is(err, films.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "film not found")
		case errors.Is(err, films.ErrorBadRequest):
			return nil, status.Errorf(codes.InvalidArgument, "invalid feedback data")
		default:
			return nil, status.Errorf(codes.Internal, "failed to send feedback")
		}
	}
	feedback.Sanitize()

	return &gen.SendFeedbackResponse{
		Feedback: &gen.FilmFeedback{
			Id:            feedback.ID.String(),
			UserId:        feedback.UserID.String(),
			FilmId:        feedback.FilmID.String(),
			Title:         feedback.Title,
			Text:          feedback.Text,
			Rating:        int32(feedback.Rating),
			CreatedAt:     feedback.CreatedAt.String(),
			UpdatedAt:     feedback.UpdatedAt.String(),
			UserLogin:     feedback.UserLogin,
			UserAvatar:    feedback.UserAvatar,
			IsMine:        feedback.IsMine,
			NewFilmRating: feedback.NewFilmRating,
		},
	}, nil
}

func (g GrpcFilmsHandler) SetRating(ctx context.Context, in *gen.SetRatingRequest) (*gen.SetRatingResponse, error) {
	filmID, err := uuid.FromString(in.FilmId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid film ID")
	}

	req := models.FilmFeedbackInput{
		Rating: int(in.RatingInput.Rating),
	}

	feedback, err := g.uc.SetRating(ctx, req, filmID)
	if err != nil {
		switch {
		case errors.Is(err, films.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "film not found")
		case errors.Is(err, films.ErrorBadRequest):
			return nil, status.Errorf(codes.InvalidArgument, "invalid rating data")
		default:
			return nil, status.Errorf(codes.Internal, "failed to set rating")
		}
	}
	feedback.Sanitize()

	return &gen.SetRatingResponse{
		Feedback: &gen.FilmFeedback{
			Id:            feedback.ID.String(),
			UserId:        feedback.UserID.String(),
			FilmId:        feedback.FilmID.String(),
			Title:         feedback.Title,
			Text:          feedback.Text,
			Rating:        int32(feedback.Rating),
			CreatedAt:     feedback.CreatedAt.String(),
			UpdatedAt:     feedback.UpdatedAt.String(),
			UserLogin:     feedback.UserLogin,
			UserAvatar:    feedback.UserAvatar,
			IsMine:        feedback.IsMine,
			NewFilmRating: feedback.NewFilmRating,
		},
	}, nil
}

func (g GrpcFilmsHandler) SiteMap(ctx context.Context, in *gen.EmptyRequest) (*gen.SiteMapResponse, error) {
	urlSet, err := g.uc.SiteMap(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate sitemap")
	}

	var urlItems []*gen.URLItem
	for _, item := range urlSet.URL {
		urlItems = append(urlItems, &gen.URLItem{
			Loc:      item.Loc,
			Priority: item.Priority,
		})
	}

	return &gen.SiteMapResponse{
		Urlset: &gen.Urlset{
			Xmlns: urlSet.Xmlns,
			Url:   urlItems,
		},
	}, nil
}
func (g GrpcFilmsHandler) GetGenre(ctx context.Context, in *gen.GetGenreRequest) (*gen.GetGenreResponse, error) {
	genreID, err := uuid.FromString(in.GenreId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid genre ID")
	}

	genre, err := g.guc.GetGenre(ctx, genreID)
	if err != nil {
		switch {
		case errors.Is(err, films.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "genre not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get genre")
		}
	}
	genre.Sanitize()

	return &gen.GetGenreResponse{
		Genre: &gen.Genre{
			Id:          genre.ID.String(),
			Name:        genre.Title,
			Description: genre.Description,
			Icon:        genre.Icon,
		},
	}, nil
}

func (g GrpcFilmsHandler) GetGenres(ctx context.Context, in *gen.GetGenresRequest) (*gen.GetGenresResponse, error) {
	pager := models.Pager{
		Count:  int(in.Pager.Count),
		Offset: int(in.Pager.Offset),
	}

	genres, err := g.guc.GetGenres(ctx, pager)
	if err != nil {
		switch {
		case errors.Is(err, films.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "genres not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get genres")
		}
	}

	var result []*gen.Genre
	for i := range genres {
		genres[i].Sanitize()
		result = append(result, &gen.Genre{
			Id:          genres[i].ID.String(),
			Name:        genres[i].Title,
			Description: genres[i].Description,
			Icon:        genres[i].Icon,
		})
	}

	return &gen.GetGenresResponse{
		Genres: result,
	}, nil
}

func (g GrpcFilmsHandler) GetFilmsByGenre(ctx context.Context, in *gen.GetFilmsByGenreRequest) (*gen.GetFilmsByGenreResponse, error) {
	genreID, err := uuid.FromString(in.GenreId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid genre ID")
	}

	pager := models.Pager{
		Count:  int(in.Pager.Count),
		Offset: int(in.Pager.Offset),
	}

	films, err := g.guc.GetFilmsByGenre(ctx, genreID, pager)
	if err != nil {
		switch {
		case errors.Is(err, genres.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "films not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get films by genre")
		}
	}

	var result []*gen.MainPageFilm
	for i := range films {
		films[i].Sanitize()
		result = append(result, &gen.MainPageFilm{
			Id:     films[i].ID.String(),
			Cover:  films[i].Cover,
			Title:  films[i].Title,
			Rating: films[i].Rating,
			Year:   int32(films[i].Year),
			Genre:  films[i].Genre,
		})
	}

	return &gen.GetFilmsByGenreResponse{
		Films: result,
	}, nil
}
func (g GrpcFilmsHandler) GetActor(ctx context.Context, in *gen.GetActorRequest) (*gen.GetActorResponse, error) {
	actorID, err := uuid.FromString(in.ActorId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid actor ID")
	}

	actor, err := g.auc.GetActor(ctx, actorID)
	if err != nil {
		switch {
		case errors.Is(err, actors.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "actor not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get actor")
		}
	}
	actor.Sanitize()

	return &gen.GetActorResponse{
		Actor: &gen.ActorPage{
			Id:            actor.ID.String(),
			RussianName:   actor.RussianName,
			OriginalName:  actor.OriginalName,
			Photo:         actor.Photo,
			Height:        int32(actor.Height),
			BirthDate:     actor.BirthDate.String(),
			Age:           int32(actor.Age),
			ZodiacSign:    actor.ZodiacSign,
			BirthPlace:    actor.BirthPlace,
			MaritalStatus: actor.MaritalStatus,
			FilmsNumber:   int32(actor.FilmsNumber),
		},
	}, nil
}

func (g GrpcFilmsHandler) GetFilmsByActor(ctx context.Context, in *gen.GetFilmsByActorRequest) (*gen.GetFilmsByActorResponse, error) {
	actorID, err := uuid.FromString(in.ActorId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid actor ID")
	}

	pager := models.Pager{
		Count:  int(in.Pager.Count),
		Offset: int(in.Pager.Offset),
	}

	films, err := g.auc.GetFilmsByActor(ctx, actorID, pager)
	if err != nil {
		switch {
		case errors.Is(err, actors.ErrorNotFound):
			return nil, status.Errorf(codes.NotFound, "films not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to get films by actor")
		}
	}

	var result []*gen.MainPageFilm
	for i := range films {
		films[i].Sanitize()
		result = append(result, &gen.MainPageFilm{
			Id:     films[i].ID.String(),
			Cover:  films[i].Cover,
			Title:  films[i].Title,
			Rating: films[i].Rating,
			Year:   int32(films[i].Year),
			Genre:  films[i].Genre,
		})
	}

	return &gen.GetFilmsByActorResponse{
		Films: result,
	}, nil
}

func (g GrpcFilmsHandler) ValidateUser(ctx context.Context, in *gen.ValidateUserRequest) (*gen.ValidateUserResponse, error) {
	token := in.Token

	user, err := g.uc.ValidateAndGetUser(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	user.Sanitize()
	return &gen.ValidateUserResponse{
		ID:      user.ID.String(),
		Version: int32(user.Version),
		Login:   user.Login,
		Avatar:  user.Avatar,
	}, nil
}
