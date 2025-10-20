package genres

import (
	"context"
	"kinopoisk/internal/models"

	uuid "github.com/satori/go.uuid"
)

type GenreUsecase interface {
	GetGenre(ctx context.Context, id uuid.UUID) (models.Genre, error)
	GetGenres(ctx context.Context, count int, offset int) ([]models.Genre, error)
}

type GenreRepo interface {
	GetGenreByID(ctx context.Context, id uuid.UUID) (models.Genre, error)
	GetGenresWithPagination(ctx context.Context, limit, offset int) ([]models.Genre, error)
}
