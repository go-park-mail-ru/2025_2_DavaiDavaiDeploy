package usecase

import (
	"context"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/search"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
)

type SearchUsecase struct {
	searchRepo search.SearchRepo
}

func NewSearchUsecase(repo search.SearchRepo) *SearchUsecase {
	return &SearchUsecase{
		searchRepo: repo,
	}
}

func (uc *SearchUsecase) GetFilmsFromSearch(ctx context.Context, searchString string, pager models.Pager) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	mainPageFilms, err := uc.searchRepo.GetFilmsFromSearch(ctx, searchString, pager.Count, pager.Offset)
	if err != nil {
		return []models.MainPageFilm{}, err
	}

	if len(mainPageFilms) == 0 {
		logger.Info("no films")
		return []models.MainPageFilm{}, nil
	}

	return mainPageFilms, nil
}

func (uc *SearchUsecase) GetActorsFromSearch(ctx context.Context, searchString string, pager models.Pager) ([]models.MainPageActor, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	mainPageActors, err := uc.searchRepo.GetActorsFromSearch(ctx, searchString, pager.Count, pager.Offset)
	if err != nil {
		return []models.MainPageActor{}, err
	}

	if len(mainPageActors) == 0 {
		logger.Info("no actors")
		return []models.MainPageActor{}, nil
	}

	return mainPageActors, nil
}
