package usecase

import (
	"context"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/compilations"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"

	uuid "github.com/satori/go.uuid"
)

type CompilationsUsecase struct {
	compilationRepo compilations.CompilationsRepo
}

func NewCompilationUsecase(compilationsRepo compilations.CompilationsRepo) *CompilationsUsecase {
	return &CompilationsUsecase{compilationRepo: compilationsRepo}
}

func (uc *CompilationsUsecase) GetCompilation(ctx context.Context, id uuid.UUID) (models.Compilation, error) {
	neededCompilation, err := uc.compilationRepo.GetCompilationByID(ctx, id)
	if err != nil {
		return models.Compilation{}, err
	}
	return neededCompilation, nil
}

func (uc *CompilationsUsecase) GetCompilations(ctx context.Context, pager models.Pager) ([]models.Compilation, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	allCompilations, err := uc.compilationRepo.GetCompilationsWithPagination(ctx, pager.Count, pager.Offset)
	if err != nil {
		return []models.Compilation{}, err
	}

	if len(allCompilations) == 0 {
		logger.Info("no compilations")
		return []models.Compilation{}, compilations.ErrorNotFound
	}
	return allCompilations, nil
}

func (uc *CompilationsUsecase) GetFilmsByCompilation(ctx context.Context, id uuid.UUID, pager models.Pager) ([]models.FavFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	films, err := uc.compilationRepo.GetFilmsByCompilation(ctx, id, pager.Count, pager.Offset)
	if err != nil {
		return []models.FavFilm{}, err
	}

	if len(films) == 0 {
		logger.Info("compilation has no films")
		return []models.FavFilm{}, compilations.ErrorNotFound
	}
	return films, nil
}
