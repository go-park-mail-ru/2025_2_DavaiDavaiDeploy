package films

import (
	"context"
	"kinopoisk/internal/models"

	uuid "github.com/satori/go.uuid"
)

type FilmUsecase interface {
	GetPromoFilm(ctx context.Context) (models.PromoFilm, error)
	GetFilms(ctx context.Context, pager models.Pager) ([]models.MainPageFilm, error)
	GetFilm(ctx context.Context, id uuid.UUID) (models.FilmPage, error)
	GetFilmFeedbacks(ctx context.Context, id uuid.UUID, pager models.Pager) ([]models.FilmFeedback, error)
	SendFeedback(ctx context.Context, req models.FilmFeedbackInput, filmID uuid.UUID) (models.FilmFeedback, error)
	SetRating(ctx context.Context, req models.FilmFeedbackInput, filmID uuid.UUID) (models.FilmFeedback, error)
	ValidateAndGetUser(ctx context.Context, token string) (models.User, error)
	SiteMap(ctx context.Context) (models.Urlset, error)
}

type FilmRepo interface {
	GetFilmByID(ctx context.Context, id uuid.UUID) (models.Film, error)
	GetGenreTitle(ctx context.Context, genreID uuid.UUID) (string, error)
	GetFilmAvgRating(ctx context.Context, filmID uuid.UUID) (float64, error)
	GetFilmsWithPagination(ctx context.Context, limit, offset int) ([]models.MainPageFilm, error)
	GetFilmPage(ctx context.Context, filmID uuid.UUID) (models.FilmPage, error)
	GetFilmFeedbacks(ctx context.Context, filmID uuid.UUID, limit, offset int) ([]models.FilmFeedback, error)
	CheckUserFeedbackExists(ctx context.Context, userID, filmID uuid.UUID) (models.FilmFeedback, error)
	UpdateFeedback(ctx context.Context, feedback models.FilmFeedback) error
	CreateFeedback(ctx context.Context, feedback models.FilmFeedback) error
	SetRating(ctx context.Context, feedback models.FilmFeedback) error
	GetPromoFilmByID(ctx context.Context, id uuid.UUID) (models.PromoFilm, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}
