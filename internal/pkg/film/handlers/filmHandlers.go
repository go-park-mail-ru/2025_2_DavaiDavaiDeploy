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
		json.NewEncoder(w).Encode(errorResp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedFilms)
}

func (c *FilmHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader((http.StatusBadRequest))
		json.NewEncoder(w).Encode(errorResp)
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
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(neededGenre)
}

func (c *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
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
			Message: "invalid id",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (c *FilmHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	var req models.Genre
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}
	neededGenre := req.ID

	var result []models.Film
	for i, film := range repo.Films {
		for _, genre := range film.Genres {
			if neededGenre == genre.ID {
				result = append(result, repo.Films[i])
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (c *FilmHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres := repo.Genres

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}
