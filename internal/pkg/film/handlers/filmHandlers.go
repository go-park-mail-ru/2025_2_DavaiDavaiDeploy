package filmHandlers

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/repo"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type FilmHandler struct {
}

func NewFilmHandler() *FilmHandler {
	return &FilmHandler{}
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
	films := repo.Films
	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)
	filmTotal := len(films)

	endingIndex := min(offset+count, filmTotal)

	paginatedFilms := films[offset:endingIndex]

	if len(paginatedFilms) == 0 {
		errorResp := models.Error{
			Message: "bad request",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader((http.StatusBadRequest))
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(paginatedFilms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetGenre godoc
// @Summary      Get genre by ID
// @Tags         genres
// @Produce      json
// @Param        id   path      string  true  "Genre ID"
// @Success      200  {object}  models.Genre
// @Failure      400  {object}  models.Error
// @Router       /genres/{id} [get]
func (c *FilmHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader((http.StatusBadRequest))
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	var neededGenre models.Genre
	for i, genre := range repo.Genres {
		if genre.ID == id {
			neededGenre = repo.Genres[i]
		}
	}

	if neededGenre.ID == uuid.Nil {
		errorResp := models.Error{
			Message: "invalid id",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(neededGenre)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetFilm godoc
// @Summary      Get film by ID
// @Tags         films
// @Produce      json
// @Param        id   path      string  true  "Film ID"
// @Success      200  {object}  models.Film
// @Failure      400  {object}  models.Error
// @Router       /films/{id} [get]
func (c *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var result models.Film
	for i, film := range repo.Films {
		if film.ID == id {
			result = repo.Films[i]
		}
	}

	if result.ID == uuid.Nil {
		errorResp := models.Error{
			Message: "Invalid id",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
		errorResp := models.Error{
			Message: "Invalid genre id",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var result []models.Film
	for _, film := range repo.Films {
		for _, genre := range film.Genres {
			if neededGenre == genre.ID {
				result = append(result, film)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
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
func (c *FilmHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres := repo.Genres

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(genres)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
