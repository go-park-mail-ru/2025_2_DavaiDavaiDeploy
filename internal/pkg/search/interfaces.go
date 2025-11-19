package search

import (
	"context"
	"kinopoisk/internal/models"
)

type SearchRepo interface {
	GetFilmsFromSearch(ctx context.Context, limit, offset int) ([]models.MainPageFilm, error)
	GetActorsFromSearch(ctx context.Context, limit, offset int) ([]models.MainPageActor, error)
}
