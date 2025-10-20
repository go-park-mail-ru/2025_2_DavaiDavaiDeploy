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
		return models.Genre{}, errors.New("No such genre")
	}
	return neededGenre, nil
}

func (uc *GenreUsecase) GetGenres(ctx context.Context, limit, offset int) ([]models.Genre, error) {
	neededGenre, err := uc.genreRepo.GetGenresWithPagination(ctx, limit, offset)
	if err != nil {
		return []models.Genre{}, errors.New("No genres")
	}
	return neededGenre, nil
}
