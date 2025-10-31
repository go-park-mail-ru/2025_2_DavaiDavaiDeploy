package filmHandlers

import (
	"context"
	"encoding/json"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

var (
	CookieName = "DDFilmsJWT"
)

type FilmHandler struct {
	uc films.FilmUsecase
}

func NewFilmHandler(uc films.FilmUsecase) *FilmHandler {
	return &FilmHandler{uc: uc}
}

// GetPromoFilm godoc
// @Summary Get promotional film
// @Description Get the promo film
// @Tags films
// @Produce json
// @Success 200 {object} models.PromoFilm
// @Failure 500 {object} models.Error
// @Router /films/promo [get]
func (c *FilmHandler) GetPromoFilm(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	film, err := c.uc.GetPromoFilm(r.Context())
	if err != nil {
		helpers.WriteError(w, http.StatusNotFound, err)
		return
	}

	film.Sanitize()
	helpers.WriteJSON(w, film)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}

// GetFilms godoc
// @Summary      List films
// @Tags         films
// @Produce      json
// @Param        count   query     int  false  "Number of films" default(10)
// @Param        offset  query     int  false  "Offset" default(0)
// @Success      200     {array}   models.MainPageFilm
// @Failure      400     {object}  models.Error
// @Router       /films [get]
func (c *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	pager := helpers.GetPagerFromRequest(r)

	films, err := c.uc.GetFilms(r.Context(), pager)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	for i := range films {
		films[i].Sanitize()
	}
	helpers.WriteJSON(w, films)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}

// GetFilm godoc
// @Summary      Get film by ID
// @Tags         films
// @Produce      json
// @Param        id   path      string  true  "Film ID"
// @Success      200  {object}  models.FilmPage
// @Failure      400  {object}  models.Error
// @Router       /films/{id} [get]
func (c *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	film, err := c.uc.GetFilm(r.Context(), id)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	film.Sanitize()
	helpers.WriteJSON(w, film)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}

func (c *FilmHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
		var token string
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
		}
		if token != "" {
			user, err := c.uc.ValidateAndGetUser(r.Context(), token)
			if err == nil {
				user.Sanitize()
				ctx := context.WithValue(r.Context(), auth.UserKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		log.LogHandlerInfo(logger, "Success", http.StatusOK)
		next.ServeHTTP(w, r)
	})
}

// GetFilmFeedbacks godoc
// @Summary Get film reviews
// @Tags films
// @Produce json
// @Param        id   path      string  true  "Film ID"
// @Success 200 {array} models.FilmFeedback
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /films/{id}/feedbacks [get]
func (c *FilmHandler) GetFilmFeedbacks(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	feedbacks, err := c.uc.GetFilmFeedbacks(r.Context(), id, pager)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	for i := range feedbacks {
		feedbacks[i].Sanitize()
	}

	helpers.WriteJSON(w, feedbacks)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}

// SendFeedback godoc
// @Summary Add film review
// @Tags films
// @Accept json
// @Produce json
// @Param        id   path      string  true  "Film ID"
// @Success 201 {object} models.FilmFeedback
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /films/{id}/feedbacks [post]
func (c *FilmHandler) SendFeedback(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var req models.FilmFeedbackInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid request"), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Sanitize()

	feedback, err := c.uc.SendFeedback(r.Context(), req, filmID)

	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	feedback.Sanitize()
	helpers.WriteJSON(w, feedback)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}

// SetRating godoc
// @Summary Rate a film
// @Tags films
// @Accept json
// @Produce json
// @Param        id   path      string  true  "Film ID"
// @Param input body models.FilmFeedbackInput true "Rating data (rating 1-10 is required)"
// @Success 200 {object} models.FilmFeedback
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /films/{id}/rating [post]
func (c *FilmHandler) SetRating(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var req models.FilmFeedbackInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid request"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	req.Sanitize()

	rating, err := c.uc.SetRating(r.Context(), req, filmID)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	rating.Sanitize()
	helpers.WriteJSON(w, rating)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}
