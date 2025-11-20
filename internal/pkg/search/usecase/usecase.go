package usecase

import (
	"context"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/search"
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
	mainPageFilms, err := uc.searchRepo.GetFilmsFromSearch(ctx, searchString, pager.Count, pager.Offset)
	if err != nil {
		return []models.MainPageFilm{}, err
	}

	return mainPageFilms, nil
}

func (uc *SearchUsecase) GetActorsFromSearch(ctx context.Context, searchString string, pager models.Pager) ([]models.MainPageActor, error) {
	mainPageActors, err := uc.searchRepo.GetActorsFromSearch(ctx, searchString, pager.Count, pager.Offset)
	if err != nil {
		return []models.MainPageActor{}, err
	}

	return mainPageActors, nil
}
