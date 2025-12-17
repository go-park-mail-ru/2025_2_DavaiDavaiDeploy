package search

import (
	"context"
	"kinopoisk/internal/models"
)

type SearchRepo interface {
	GetFilmsFromSearch(ctx context.Context, searchString string, limit, offset int) ([]models.MainPageFilm, error)
	GetActorsFromSearch(ctx context.Context, searchString string, limit, offset int) ([]models.MainPageActor, error)
}

type SearchUsecase interface {
	GetFilmsFromSearch(ctx context.Context, searchString string, pager models.Pager) ([]models.MainPageFilm, error)
	GetActorsFromSearch(ctx context.Context, searchString string, pager models.Pager) ([]models.MainPageActor, error)
}
