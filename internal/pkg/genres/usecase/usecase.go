package usecase

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/genres"

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
	genres, err := uc.genreRepo.GetGenresWithPagination(ctx, pager.Count, pager.Offset)
	if err != nil {
		return []models.Genre{}, errors.New("no genres")
	}

	if len(genres) == 0 {
		return []models.Genre{}, errors.New("no genres")
	}
	return genres, nil
}

func (uc *GenreUsecase) GetFilmsByGenre(ctx context.Context, id uuid.UUID, pager models.Pager) ([]models.MainPageFilm, error) {
	films, err := uc.genreRepo.GetFilmsByGenre(ctx, id, pager.Count, pager.Offset)
	if err != nil {
		return []models.MainPageFilm{}, errors.New("no films")
	}

	if len(films) == 0 {
		return []models.MainPageFilm{}, errors.New("no films")
	}
	return films, nil
}
