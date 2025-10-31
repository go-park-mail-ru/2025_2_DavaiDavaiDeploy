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
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /genres/{id} [get]
func (g *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of genre"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	neededGenre, err := g.uc.GetGenre(r.Context(), id)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	neededGenre.Sanitize()
	helpers.WriteJSON(w, neededGenre)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}

// GetGenres godoc
// @Summary Get list of all genres
// @Tags genres
// @Produce json
// @Success 200 {array} models.Genre
// @Failure 500 {object} models.Error
// @Router /genres [get]
func (g *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	pager := helpers.GetPagerFromRequest(r)

	genres, err := g.uc.GetGenres(r.Context(), pager)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}
	for i := range genres {
		genres[i].Sanitize()
	}
	helpers.WriteJSON(w, genres)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}

// GetFilmsByGenre godoc
// @Summary Get films by genre
// @Tags genres
// @Produce json
// @Param        id   path      string  true  "Genre ID"
// @Success 200 {array} models.MainPageFilm
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /genres/{id}/films [get]
func (g *GenreHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	idStr := vars["id"]

	neededGenre, err := uuid.FromString(idStr)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of genre"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	films, err := g.uc.GetFilmsByGenre(r.Context(), neededGenre, pager)
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
