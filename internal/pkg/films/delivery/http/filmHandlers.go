package filmHandlers

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/helpers"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type FilmHandler struct {
	uc films.FilmUsecase
}

func NewFilmHandler(uc films.FilmUsecase) *FilmHandler {
	return &FilmHandler{uc: uc}
}

// GetPromoFilm godoc
// @Summary      Get the film to the main page
// @Tags         films
// @Produce      json
// @Success      200  {object}  models.PromoFilm
// @Router       /films/promo [get]
func (c *FilmHandler) GetPromoFilm(w http.ResponseWriter, r *http.Request) {
	film, err := c.uc.GetPromoFilm(r.Context())
	if err != nil {
		helpers.WriteError(w, 500, err)
		return
	}

	helpers.WriteJSON(w, film)
}

// GetFilms godoc
// @Summary      List films
// @Tags         films
// @Produce      json
// @Param        count   query     int  false  "Number of films" default(10)
// @Param        offset  query     int  false  "Offset" default(0)
// @Success      200     {array}   models.Film
// @Failure      400     {object}  models.Error
// @Router       /films [get]
func (c *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	pager := helpers.GetPagerFromRequest(r)

	films, err := c.uc.GetFilms(r.Context(), pager)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	helpers.WriteJSON(w, films)
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
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	film, err := c.uc.GetFilm(r.Context(), id)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	helpers.WriteJSON(w, film)
}

// GetFilmsByGenre godoc
// @Summary      Get films by genre ID
// @Tags         films
// @Produce      json
// @Param        id   path      string  true  "Genre ID"
// @Success      200  {array}   models.Film
// @Failure      400  {object}  models.Error
// @Router       /films/genre/{id} [get]
func (c *FilmHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	neededGenre, err := uuid.FromString(idStr)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	films, err := c.uc.GetFilmsByGenre(r.Context(), neededGenre, pager)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	helpers.WriteJSON(w, films)
}

// GetFilmsByActor godoc
// @Summary      Get films by actor ID
// @Tags         films
// @Produce      json
// @Param        id   path      string  true  "Actor ID"
// @Success      200  {array}   models.Film
// @Failure      400  {object}  models.Error
// @Router       /films/actor/{id} [get]
func (c *FilmHandler) GetFilmsByActor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	neededActor, err := uuid.FromString(idStr)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	films, err := c.uc.GetFilmsByActor(r.Context(), neededActor, pager)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	helpers.WriteJSON(w, films)
}

func (c *FilmHandler) GetFilmFeedbacks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	films, err := c.uc.GetFilmFeedbacks(r.Context(), id, pager)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	helpers.WriteJSON(w, films)
}

func (c *FilmHandler) SendFeedback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	var req models.FilmFeedbackInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	feedback, err := c.uc.SendFeedback(r.Context(), req, filmID)

	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	helpers.WriteJSON(w, feedback)
}

func (c *FilmHandler) SetRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	// feedback, err := c.filmRepo.CheckUserFeedbackExists(r.Context(), user.ID, filmID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(feedback)
	// 	return // у нас нельзя менять рейтинг, но можно поменять отзыв
	// }

	var req models.FilmFeedbackInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rating, err := c.uc.SetRating(r.Context(), req, filmID)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	helpers.WriteJSON(w, rating)
}
