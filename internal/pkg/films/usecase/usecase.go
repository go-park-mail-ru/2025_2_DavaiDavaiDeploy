package usecase

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/films"
	"time"

	uuid "github.com/satori/go.uuid"
)

type FilmUsecase struct {
	filmRepo films.FilmRepo
}

func NewFilmUsecase(repo films.FilmRepo) *FilmUsecase {
	return &FilmUsecase{
		filmRepo: repo,
	}
}

func (uc *FilmUsecase) GetPromoFilm(ctx context.Context) (models.PromoFilm, error) {
	film, err := uc.filmRepo.GetFilmByID(ctx, uuid.FromStringOrNil("8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c"))
	if err != nil {
		return models.PromoFilm{}, errors.New("no films")
	}

	genre, err := uc.filmRepo.GetGenreTitle(ctx, film.GenreID)
	if err != nil {
		genre = "Unknown"
	}

	avgRating, err := uc.filmRepo.GetFilmAvgRating(ctx, film.ID)
	if err != nil {
		avgRating = 0.0
	}

	promoFilm := models.PromoFilm{
		ID:               film.ID,
		Image:            *film.Poster,
		Title:            film.Title,
		Rating:           avgRating,
		ShortDescription: film.ShortDescription,
		Year:             film.Year,
		Genre:            genre,
		Duration:         film.Duration,
	}
	return promoFilm, nil
}

func (uc *FilmUsecase) GetFilms(ctx context.Context, pager models.Pager) ([]models.MainPageFilm, error) {

	films, err := uc.filmRepo.GetFilmsWithPagination(ctx, pager.Count, pager.Offset)
	if err != nil {
		return []models.MainPageFilm{}, errors.New("no films")
	}

	if len(films) == 0 {
		return []models.MainPageFilm{}, errors.New("no films")
	}

	return films, nil
}

func (uc *FilmUsecase) GetFilm(ctx context.Context, id uuid.UUID) (models.FilmPage, error) {
	film, err := uc.filmRepo.GetFilmPage(ctx, id)
	if err != nil {
		return models.FilmPage{}, errors.New("no films")
	}

	return film, nil
}

func (uc *FilmUsecase) GetFilmFeedbacks(ctx context.Context, id uuid.UUID, pager models.Pager) ([]models.FilmFeedback, error) {
	feedbacks, err := uc.filmRepo.GetFilmFeedbacks(ctx, id, pager.Count, pager.Offset)
	if err != nil {
		return []models.FilmFeedback{}, errors.New("no feedbacks")
	}

	if len(feedbacks) == 0 {
		return []models.FilmFeedback{}, errors.New("no feedbacks")
	}
	return feedbacks, nil
}

func (uc *FilmUsecase) SendFeedback(ctx context.Context, req models.FilmFeedbackInput, filmID uuid.UUID) (models.FilmFeedback, error) {
	user, ok := ctx.Value(auth.UserKey).(models.User)
	if !ok {
		return models.FilmFeedback{}, errors.New("user not authenticated")
	}

	if req.Rating < 1 || req.Rating > 10 {
		return models.FilmFeedback{}, errors.New("rating must be between 1 and 10")
	}

	if len(req.Title) < 1 || len(req.Title) > 100 {
		return models.FilmFeedback{}, errors.New("title length must be between 1 and 100")
	}

	if len(req.Text) < 1 || len(req.Text) > 1000 {
		return models.FilmFeedback{}, errors.New("text length must be between 1 and 1000")
	}

	existingFeedback, err := uc.filmRepo.CheckUserFeedbackExists(ctx, user.ID, filmID)
	if err == nil {
		// отзыв существует - обновляем
		existingFeedback.Title = &req.Title
		existingFeedback.Text = &req.Text
		existingFeedback.Rating = req.Rating

		err := uc.filmRepo.UpdateFeedback(ctx, existingFeedback)
		if err != nil {
			return models.FilmFeedback{}, err
		}

		return existingFeedback, nil
	}

	// создаем новый отзыв
	feedback := models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    user.ID,
		FilmID:    filmID,
		Title:     &req.Title,
		Text:      &req.Text,
		Rating:    req.Rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := uc.filmRepo.CreateFeedback(ctx, feedback); err != nil {
		return models.FilmFeedback{}, nil
	}
	return feedback, nil
}

func (uc *FilmUsecase) SetRating(ctx context.Context, req models.FilmFeedbackInput, filmID uuid.UUID) (models.FilmFeedback, error) {
	user, ok := ctx.Value(auth.UserKey).(models.User)
	if !ok {
		return models.FilmFeedback{}, errors.New("user not authenticated")
	}

	// feedback, err := c.filmRepo.CheckUserFeedbackExists(r.Context(), user.ID, filmID)
	// if err != nil {
	// 	fmt.Println("suslik")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(feedback)
	// 	return // у нас нельзя менять рейтинг, но можно поменять отзыв
	// }

	if req.Rating < 1 || req.Rating > 10 {
		return models.FilmFeedback{}, errors.New("rating must be between 1 and 10")
	}

	newFeedback := models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    user.ID,
		FilmID:    filmID,
		Rating:    req.Rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := uc.filmRepo.CreateFeedback(ctx, newFeedback)
	if err != nil {
		return models.FilmFeedback{}, errors.New("no feedback")
	}

	return newFeedback, nil
}
