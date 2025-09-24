package film

import (
	"encoding/json"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/repo"
	"net/http"
	"strconv"
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

func (c *FilmHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("god forbid")
		return
	}
	id := req.ID
	var neededUser *models.User
	for _, user := range repo.Users {
		if user.ID == id {
			neededUser = &user
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(neededUser)
}

func (c *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	films := repo.Films
	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)
	filmTotal := len(films)

	endingIndex := min(offset+count, filmTotal)

	paginatedFilms := films[offset:endingIndex]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedFilms)
}

func (c *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	var req models.Film
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("god forbid")
		return
	}
	id := req.ID

	films := repo.Films

	var result *models.Film
	for _, film := range films {
		if film.ID == id {
			result = &film
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (c *FilmHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	var req models.Film
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("god forbid")
		return
	}
	neededGenre := req.ID

	films := repo.Films

	var result []models.Film
	for _, film := range films {
		for _, genre := range film.Genres {
			if neededGenre == genre.ID {
				result = append(result, film)
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
