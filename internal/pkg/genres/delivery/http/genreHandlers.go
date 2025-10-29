package http

import (
	"kinopoisk/internal/pkg/genres"
	"kinopoisk/internal/pkg/helpers"
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
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
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
}

// GetGenres godoc
// @Summary Get list of all genres
// @Tags genres
// @Produce json
// @Success 200 {array} models.Genre
// @Failure 500 {object} models.Error
// @Router /genres [get]
func (g *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)
	idStr := vars["id"]

	neededGenre, err := uuid.FromString(idStr)
	if err != nil {
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
}
