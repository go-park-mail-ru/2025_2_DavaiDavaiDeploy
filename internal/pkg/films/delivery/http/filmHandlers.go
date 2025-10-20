package filmHandlers

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/films"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type FilmHandler struct {
	uc films.FilmUsecase
}

func NewFilmHandler(uc films.FilmUsecase) *FilmHandler {
	return &FilmHandler{uc: uc}
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

// GetPromoFilm godoc
// @Summary      Get the film to the main page
// @Tags         films
// @Produce      json
// @Success      200  {object}  models.PromoFilm
// @Router       /films/promo [get]
func (c *FilmHandler) GetPromoFilm(w http.ResponseWriter, r *http.Request) {
	film, err := c.uc.GetPromoFilm(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(film)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)

	films, err := c.uc.GetFilms(r.Context(), count, offset)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	json.NewEncoder(w).Encode(films)
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
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	result, err := c.uc.GetFilm(r.Context(), id)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

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

	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)

	films, err := c.uc.GetFilmsByGenre(r.Context(), neededGenre, count, offset)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idStr := vars["id"]

	neededActor, err := uuid.FromString(idStr)
	if err != nil {
		errorResp := models.Error{
			Message: "Invalid actor id",
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)

	films, err := c.uc.GetFilmsByActor(r.Context(), neededActor, count, offset)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	err = json.NewEncoder(w).Encode(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *FilmHandler) GetFilmFeedbacks(w http.ResponseWriter, r *http.Request) {
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

	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)

	films, err := c.uc.GetFilmFeedbacks(r.Context(), id, count, offset)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	err = json.NewEncoder(w).Encode(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *FilmHandler) SendFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
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
		errorResp := models.Error{
			Message: "Failed to create feedback: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedback)

}

func (c *FilmHandler) SetRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	// feedback, err := c.filmRepo.CheckUserFeedbackExists(r.Context(), user.ID, filmID)
	// if err != nil {
	// 	fmt.Println("suslik")
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

	feedback, err := c.uc.SetRating(r.Context(), req, filmID)
	if err != nil {
		errorResp := models.Error{
			Message: "Failed to set rating: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	err = json.NewEncoder(w).Encode(feedback)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
