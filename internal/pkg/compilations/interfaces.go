package compilations

import (
	"context"
	"kinopoisk/internal/models"

	uuid "github.com/satori/go.uuid"
)

type CompilationsUsecase interface {
	GetCompilation(ctx context.Context, id uuid.UUID) (models.Compilation, error)
	GetCompilations(ctx context.Context, pager models.Pager) ([]models.Compilation, error)
	GetFilmsByCompilation(ctx context.Context, id uuid.UUID, pager models.Pager) ([]models.FavFilm, error)
}

type CompilationsRepo interface {
	GetCompilationByID(ctx context.Context, id uuid.UUID) (models.Compilation, error)
	GetCompilationsWithPagination(ctx context.Context, limit, offset int) ([]models.Compilation, error)
	GetFilmsByCompilation(ctx context.Context, compilationID uuid.UUID, limit, offset int) ([]models.FavFilm, error)
}
