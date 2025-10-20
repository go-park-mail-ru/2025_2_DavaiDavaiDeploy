package http

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/genres"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type GenreHandler struct {
	uc genres.GenreUsecase
}

func NewGenreHandler(uc genres.GenreUsecase) *GenreHandler {
	return &GenreHandler{uc: uc}
}

func GetParameter(r *http.Request, s string, defaultValue int) int {
	strValue := r.URL.Query().Get(s)
	if strValue == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(strValue)
	if err != nil || result <= 0 {
		return defaultValue
	}
	return result
}

// GetGenre godoc
// @Summary      Get genre by ID
// @Tags         genres
// @Produce      json
// @Param        id   path      string  true  "Genre ID"
// @Success      200  {object}  models.Genre
// @Failure      400  {object}  models.Error
// @Router       /genres/{id} [get]
func (g *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.WriteHeader((http.StatusBadRequest))
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	neededGenre, err := g.uc.GetGenre(r.Context(), id)
	if err != nil {
		errorResp := models.Error{
			Message: "invalid id",
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = json.NewEncoder(w).Encode(neededGenre)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetGenres godoc
// @Summary      List all genres
// @Tags         genres
// @Produce      json
// @Success      200  {array}  models.Genre
// @Router       /genres [get]
func (g *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)

	genres, err := g.uc.GetGenres(r.Context(), count, offset)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(genres)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
