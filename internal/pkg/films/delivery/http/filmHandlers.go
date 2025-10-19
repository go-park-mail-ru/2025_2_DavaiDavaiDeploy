package filmHandlers

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/films/repo"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type FilmHandler struct {
	filmRepo *repo.FilmRepository
}

func NewFilmHandler(db *pgxpool.Pool) *FilmHandler {
	filmRepo := repo.NewFilmRepository(db)
	return &FilmHandler{filmRepo: filmRepo}
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
	film, err := c.filmRepo.GetFilmByID(r.Context(), uuid.FromStringOrNil("1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	genre, err := c.filmRepo.GetGenreTitle(r.Context(), film.GenreID)
	if err != nil {
		genre = "Unknown"
	}

	avgRating, err := c.filmRepo.GetFilmAvgRating(r.Context(), film.ID)
	if err != nil {
		avgRating = 0.0
	}

	response := models.PromoFilm{
		ID:               film.ID,
		Image:            film.Poster,
		Title:            film.Title,
		Rating:           avgRating,
		ShortDescription: film.ShortDescription,
		Year:             film.Year,
		Genre:            genre,
		Duration:         film.Duration,
	}

	err = json.NewEncoder(w).Encode(response)
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

	films, err := c.filmRepo.GetFilmsWithPagination(r.Context(), count, offset)
	if err != nil {
		return
	}

	if len(films) == 0 {
		errorResp := models.Error{Message: "No films found"}
		w.WriteHeader(http.StatusNotFound)
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

	result, err := c.filmRepo.GetFilmPage(r.Context(), id)
	if err != nil {
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

	films, err := c.filmRepo.GetFilmsByGenre(r.Context(), neededGenre, count, offset)
	if err != nil {
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
	vars := mux.Vars(r)
	idStr := vars["id"]

	neededActor, err := uuid.FromString(idStr)
	if err != nil {
		errorResp := models.Error{
			Message: "Invalid actor id",
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

	films, err := c.filmRepo.GetFilmsByActor(r.Context(), neededActor, count, offset)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *FilmHandler) GetFilmsFeedback(w http.ResponseWriter, r *http.Request) {
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

	films, err := c.filmRepo.GetFilmFeedbacks(r.Context(), id, count, offset)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(films)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *FilmHandler) SendFeedback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		errorResp := models.Error{
			Message: "User not authenticated",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

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

	if req.Rating < 1 || req.Rating > 10 {
		errorResp := models.Error{
			Message: "Rating must be between 1 and 10",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	if len(req.Title) < 1 || len(req.Title) > 100 {
		errorResp := models.Error{
			Message: "Title length must be between 1 and 100",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	if len(req.Text) < 1 || len(req.Text) > 1000 {
		errorResp := models.Error{
			Message: "Text length must be between 1 and 1000",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	existingFeedback, err := c.filmRepo.CheckUserFeedbackExists(r.Context(), user.ID, filmID)
	if err == nil {
		// отзыв существует - обновляем
		existingFeedback.Title = req.Title
		existingFeedback.Text = req.Text
		existingFeedback.Rating = req.Rating

		err := c.filmRepo.UpdateFeedback(r.Context(), existingFeedback)
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(existingFeedback)
		return
	}

	// создаем новый отзыв
	feedback := &models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    user.ID,
		FilmID:    filmID,
		Title:     req.Title,
		Text:      req.Text,
		Rating:    req.Rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := c.filmRepo.CreateFeedback(r.Context(), feedback); err != nil {
		errorResp := models.Error{Message: "Failed to create feedback: " + err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedback)

}

func (c *FilmHandler) SetRating(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		errorResp := models.Error{
			Message: "User not authenticated",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

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

	if req.Rating < 1 || req.Rating > 10 {
		errorResp := models.Error{
			Message: "Rating must be between 1 and 10",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	newFeedback := models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    user.ID,
		FilmID:    filmID,
		Title:     "",
		Text:      "",
		Rating:    req.Rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = c.filmRepo.CreateFeedback(r.Context(), &newFeedback)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(newFeedback)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
