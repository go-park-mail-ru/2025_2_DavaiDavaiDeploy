package usecase

import (
	"context"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type FilmUsecase struct {
	filmRepo films.FilmRepo
	secret   string
}

func NewFilmUsecase(repo films.FilmRepo) *FilmUsecase {
	return &FilmUsecase{
		filmRepo: repo,
		secret:   os.Getenv("JWT_SECRET"),
	}
}

func (uc *FilmUsecase) GetPromoFilm(ctx context.Context) (models.PromoFilm, error) {
	films := []string{"8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c", "2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c", "6ba7b810-9dad-11d1-80b4-00c04fd430c8"}

	randomIndex := rand.Intn(len(films))
	randomFilmID := films[randomIndex]

	film, err := uc.filmRepo.GetPromoFilmByID(ctx, uuid.FromStringOrNil(randomFilmID))
	if err != nil {
		return models.PromoFilm{}, err
	}

	avgRating, err := uc.filmRepo.GetFilmAvgRating(ctx, film.ID)
	if err != nil {
		avgRating = 0.0
	}

	promoFilm := models.PromoFilm{
		ID:               film.ID,
		Image:            film.Image,
		Title:            film.Title,
		Rating:           avgRating,
		ShortDescription: film.ShortDescription,
		Year:             film.Year,
		Genre:            film.Genre,
		Duration:         film.Duration,
	}
	return promoFilm, nil
}

func (uc *FilmUsecase) GetFilms(ctx context.Context, pager models.Pager) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	mainPageFilms, err := uc.filmRepo.GetFilmsWithPagination(ctx, pager.Count, pager.Offset)
	if err != nil {
		return []models.MainPageFilm{}, err
	}

	if len(mainPageFilms) == 0 {
		logger.Error("no films")
		return []models.MainPageFilm{}, films.ErrorNotFound
	}

	return mainPageFilms, nil
}

func (uc *FilmUsecase) GetUsersFavFilms(ctx context.Context, id uuid.UUID) ([]models.FavFilm, error) {
	favFilms, _ := uc.filmRepo.GetUsersFavFilms(ctx, id)
	return favFilms, nil
}

func (uc *FilmUsecase) GetFilmsForCalendar(ctx context.Context, pager models.Pager, userID uuid.UUID) ([]models.FilmInCalendar, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	filmsInCalendar, err := uc.filmRepo.GetFilmsForCalendar(ctx, pager.Count, pager.Offset)
	if err != nil {
		return []models.FilmInCalendar{}, err
	}

	if len(filmsInCalendar) == 0 {
		logger.Error("no films")
		return []models.FilmInCalendar{}, films.ErrorNotFound
	}

	for i := range filmsInCalendar {
		_, err = uc.filmRepo.CheckUserLikeExists(ctx, userID, filmsInCalendar[i].ID)
		filmsInCalendar[i].IsLiked = false
		if err == nil {
			filmsInCalendar[i].IsLiked = true
		}
	}

	return filmsInCalendar, nil
}

func (uc *FilmUsecase) GetFilm(ctx context.Context, id uuid.UUID, userID uuid.UUID) (models.FilmPage, error) {
	film, err := uc.filmRepo.GetFilmPage(ctx, id)
	if err != nil {
		return models.FilmPage{}, err
	}

	feedback, err := uc.filmRepo.CheckUserFeedbackExists(ctx, userID, id)
	film.IsReviewed = false
	emptyFeedback := ""
	if err == nil && feedback.Title != &emptyFeedback {
		film.IsReviewed = true
		film.UserRating = &feedback.Rating
	} else if err == nil {
		film.UserRating = &feedback.Rating
	}

	_, err = uc.filmRepo.CheckUserLikeExists(ctx, userID, film.ID)
	film.IsLiked = false
	if err == nil {
		film.IsLiked = true
	}

	return film, nil
}

func (uc *FilmUsecase) GetFilmFeedbacks(ctx context.Context, id uuid.UUID, userID uuid.UUID, pager models.Pager) ([]models.FilmFeedback, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	result := make([]models.FilmFeedback, 0, pager.Count+1)
	emptyFeedback := ""
	usersFeedbackLogin := ""
	if userID != uuid.Nil {
		feedback, err := uc.filmRepo.CheckUserFeedbackExists(ctx, userID, id)
		if err == nil && feedback.Text != &emptyFeedback && feedback.Text != nil {
			feedback.IsMine = true
			usersFeedbackLogin = feedback.UserLogin
			result = append(result, feedback)
		}
	}

	feedbacks, err := uc.filmRepo.GetFilmFeedbacks(ctx, id, pager.Count, pager.Offset)
	if err != nil {
		return []models.FilmFeedback{}, err
	}

	if len(feedbacks) == 0 {
		logger.Error("no feedbacks")
		return []models.FilmFeedback{}, films.ErrorNotFound
	}

	for i := range feedbacks {
		feedbacks[i].IsMine = false
		if feedbacks[i].UserLogin != usersFeedbackLogin {
			result = append(result, feedbacks[i])
		}
	}

	return result, nil
}

func (uc *FilmUsecase) SendFeedback(ctx context.Context, req models.FilmFeedbackInput, filmID uuid.UUID, userID uuid.UUID) (models.FilmFeedback, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if req.Rating < 1 || req.Rating > 10 {
		logger.Error("invalid rating")
		return models.FilmFeedback{}, films.ErrorBadRequest
	}

	if len(req.Title) < 1 || len(req.Title) > 100 {
		logger.Error("invalid length of title")
		return models.FilmFeedback{}, films.ErrorBadRequest
	}

	if len(req.Text) < 30 || len(req.Text) > 1000 {
		logger.Error("invalid length of text")
		return models.FilmFeedback{}, films.ErrorBadRequest
	}

	existingFeedback, err := uc.filmRepo.CheckUserFeedbackExists(ctx, userID, filmID)
	if err == nil {
		// отзыв существует - обновляем
		existingFeedback.Title = &req.Title
		existingFeedback.Text = &req.Text
		existingFeedback.Rating = req.Rating

		err := uc.filmRepo.UpdateFeedback(ctx, existingFeedback)
		if err != nil {
			return models.FilmFeedback{}, err
		}

		updatedFilm, _ := uc.filmRepo.GetFilmPage(ctx, filmID)
		existingFeedback.NewFilmRating = updatedFilm.Rating

		return existingFeedback, nil
	}

	// создаем новый отзыв
	feedback := models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    userID,
		FilmID:    filmID,
		Title:     &req.Title,
		Text:      &req.Text,
		Rating:    req.Rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := uc.filmRepo.CreateFeedback(ctx, feedback); err != nil {
		return models.FilmFeedback{}, err
	}

	updatedFilm, _ := uc.filmRepo.GetFilmPage(ctx, filmID)
	feedback.NewFilmRating = updatedFilm.Rating
	return feedback, nil
}

func (uc *FilmUsecase) SaveFilm(ctx context.Context, userID uuid.UUID, filmID uuid.UUID) error {
	return uc.filmRepo.SaveFilm(ctx, userID, filmID)
}

func (uc *FilmUsecase) RemoveFilm(ctx context.Context, userID uuid.UUID, filmID uuid.UUID) ([]models.FavFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var favFilms []models.FavFilm

	err := uc.filmRepo.RemoveFilm(ctx, userID, filmID)

	if err == nil {
		favFilms, err = uc.filmRepo.GetUsersFavFilms(ctx, userID)
		if err != nil {
			logger.Error("bad request")
			return []models.FavFilm{}, films.ErrorBadRequest
		}
	}
	return favFilms, nil
}

func (uc *FilmUsecase) SetRating(ctx context.Context, req models.FilmFeedbackInput, filmID uuid.UUID, userID uuid.UUID) (models.FilmFeedback, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if req.Rating < 1 || req.Rating > 10 {
		logger.Error("invalid rating")
		return models.FilmFeedback{}, films.ErrorBadRequest
	}

	existingFeedback, err := uc.filmRepo.CheckUserFeedbackExists(ctx, userID, filmID)
	if err == nil {
		// запись существует - обновляем рейтинг
		existingFeedback.Rating = req.Rating

		err := uc.filmRepo.UpdateFeedback(ctx, existingFeedback)
		if err != nil {
			return models.FilmFeedback{}, err
		}

		updatedFilm, _ := uc.filmRepo.GetFilmPage(ctx, filmID)
		existingFeedback.NewFilmRating = updatedFilm.Rating

		return existingFeedback, nil
	}

	newFeedback := models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    userID,
		FilmID:    filmID,
		Rating:    req.Rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = uc.filmRepo.CreateFeedback(ctx, newFeedback)
	if err != nil {
		return models.FilmFeedback{}, err
	}

	updatedFilm, _ := uc.filmRepo.GetFilmPage(ctx, filmID)
	newFeedback.NewFilmRating = updatedFilm.Rating

	return newFeedback, nil
}

func (uc *FilmUsecase) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.secret), nil
	})
}

func (uc *FilmUsecase) ValidateAndGetUser(ctx context.Context, token string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if token == "" {
		logger.Error("user is not authorized")
		return models.User{}, films.ErrorUnauthorized
	}

	parsedToken, err := uc.ParseToken(token)
	if err != nil || !parsedToken.Valid {
		logger.Error("invalid token")
		return models.User{}, films.ErrorUnauthorized
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("invalid claims")
		return models.User{}, films.ErrorUnauthorized
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		logger.Error("invalid exp claim")
		return models.User{}, films.ErrorUnauthorized
	}

	login, ok := claims["login"].(string)
	if !ok || login == "" {
		logger.Error("invalid login claim")
		return models.User{}, films.ErrorUnauthorized
	}

	user, err := uc.filmRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, films.ErrorUnauthorized
	}

	return user, nil
}

func (uc *FilmUsecase) SiteMap(ctx context.Context) (models.Urlset, error) {
	var urlSet models.Urlset

	urlSet.Xmlns = "https://www.sitemaps.org/schemas/sitemap/0.9/"
	urlSet.URL = append(urlSet.URL, models.URLItem{Loc: "https://ddfilms.online/"})
	mainPageFilms, err := uc.filmRepo.GetFilmsWithPagination(ctx, 10, 0)
	if err != nil {
		return models.Urlset{}, err
	}
	for _, v := range mainPageFilms {
		var item models.URLItem
		item.Loc, _ = url.JoinPath("https://ddfilms.online/films/", v.ID.String())
		item.Priority = 1.0
		urlSet.URL = append(urlSet.URL, item)
	}
	return urlSet, nil
}
