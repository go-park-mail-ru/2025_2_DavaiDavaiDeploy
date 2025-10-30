package usecase

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/genres"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"

	uuid "github.com/satori/go.uuid"
)

type GenreUsecase struct {
	genreRepo genres.GenreRepo
}

func NewGenreUsecase(genreRepo genres.GenreRepo) *GenreUsecase {
	return &GenreUsecase{genreRepo: genreRepo}
}

func (uc *GenreUsecase) GetGenre(ctx context.Context, id uuid.UUID) (models.Genre, error) {
	neededGenre, err := uc.genreRepo.GetGenreByID(ctx, id)
	if err != nil {
		return models.Genre{}, errors.New("no such genre")
	}
	return neededGenre, nil
}

func (uc *GenreUsecase) GetGenres(ctx context.Context, pager models.Pager) ([]models.Genre, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	genres, err := uc.genreRepo.GetGenresWithPagination(ctx, pager.Count, pager.Offset)
	if err != nil {
		return []models.Genre{}, errors.New("no genres")
	}

	if len(genres) == 0 {
		logger.Info("no genres")
		return []models.Genre{}, errors.New("no genres")
	}
	return genres, nil
}

func (uc *GenreUsecase) GetFilmsByGenre(ctx context.Context, id uuid.UUID, pager models.Pager) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	films, err := uc.genreRepo.GetFilmsByGenre(ctx, id, pager.Count, pager.Offset)
	if err != nil {
		return []models.MainPageFilm{}, errors.New("no films")
	}

	if len(films) == 0 {
		logger.Info("genre has no films")
		return []models.MainPageFilm{}, errors.New("no films")
	}
	return films, nil
}
