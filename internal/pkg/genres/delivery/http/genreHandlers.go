package http

import (
	"errors"
	"kinopoisk/internal/pkg/genres"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type GenreHandler struct {
	uc genres.GenreUsecase
}

func NewGenreHandler(uc genres.GenreUsecase) *GenreHandler {
	return &GenreHandler{uc: uc}
}

// GetGenre godoc
// @Summary Get genre by ID
// @Tags genres
// @Produce json
// @Param        id   path      string  true  "Genre ID"
// @Success 200 {object} models.Genre
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /genres/{id} [get]
func (g *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of genre"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	neededGenre, err := g.uc.GetGenre(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, genres.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}
	neededGenre.Sanitize()
	helpers.WriteJSON(w, neededGenre)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetGenres godoc
// @Summary Get list of all genres
// @Tags genres
// @Produce json
// @Success 200 {array} models.Genre
// @Failure 404
// @Failure 500
// @Router /genres [get]
func (g *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	pager := helpers.GetPagerFromRequest(r)

	allGenres, err := g.uc.GetGenres(r.Context(), pager)
	if err != nil {
		switch {
		case errors.Is(err, genres.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}
	for i := range allGenres {
		allGenres[i].Sanitize()
	}
	helpers.WriteJSON(w, allGenres)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetFilmsByGenre godoc
// @Summary Get films by genre
// @Tags genres
// @Produce json
// @Param        id   path      string  true  "Genre ID"
// @Success 200 {array} models.MainPageFilm
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /genres/{id}/films [get]
func (g *GenreHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	idStr := vars["id"]

	neededGenre, err := uuid.FromString(idStr)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of genre"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	films, err := g.uc.GetFilmsByGenre(r.Context(), neededGenre, pager)
	if err != nil {
		switch {
		case errors.Is(err, genres.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}
	for i := range films {
		films[i].Sanitize()
	}
	helpers.WriteJSON(w, films)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
