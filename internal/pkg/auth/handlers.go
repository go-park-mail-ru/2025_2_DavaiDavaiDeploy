package auth

import (
	"net/http"
)

type FilmController struct {
}

func NewFilmController() (*FilmController, error) {
	return &FilmController{}, nil
}

func (c *FilmController) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *FilmController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *FilmController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *FilmController) GetUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *FilmController) GetFilms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *FilmController) GetFilm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
