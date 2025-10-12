package filmHandlers

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/repo"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

// GetPromoFilm godoc
// @Summary      Get the film to the main page
// @Tags         films
// @Produce      json
// @Success      200  {object}  models.PromoFilm
// @Router       /films/promo [get]
func (c *FilmHandler) GetPromoFilm(w http.ResponseWriter, r *http.Request) {
	promoFilms := make([]models.Film, 0)

	targetFilms := []string{"Дюна: Часть вторая", "Интерстеллар", "Начало", "Темный рыцарь", "Матрица"}

	for _, film := range repo.Films {
		for _, target := range targetFilms {
			if film.Title == target {
				promoFilms = append(promoFilms, film)
				break
			}
		}
	}

	if len(promoFilms) == 0 {
		promoFilms = repo.Films
	}

	promo := promoFilms[rand.Intn(len(promoFilms))]
	var promoGenre string
	for _, genre := range repo.Genres {
		if genre.ID == promo.GenreID {
			promoGenre = genre.Title
			break
		}
	}

	response := models.PromoFilm{
		ID:               promo.ID,
		Image:            promo.Cover, // будет другая картинка
		Title:            promo.Title,
		Rating:           promo.Rating,
		ShortDescription: promo.ShortDescription,
		Year:             promo.Year,
		Genre:            promoGenre,
		Duration:         promo.Duration,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
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
	films := repo.Films
	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)
	filmTotal := len(films)
	if offset >= filmTotal {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endingIndex := min(offset+count, filmTotal)

	paginatedFilms := films[offset:endingIndex]

	if len(paginatedFilms) == 0 {
		errorResp := models.Error{
			Message: "Bad request",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader((http.StatusBadRequest))
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	var response []models.MainPageFilm
	for _, film := range paginatedFilms {
		var filmGenre string
		for _, genre := range repo.Genres {
			if genre.ID == film.GenreID {
				filmGenre = genre.Title
				break
			}
		}

		response = append(response, models.MainPageFilm{
			ID:     film.ID,
			Cover:  film.Cover,
			Title:  film.Title,
			Rating: film.Rating,
			Year:   film.Year,
			Genre:  filmGenre,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	var result models.FilmPage
	var genre string
	var country string
	var actors []models.Actor

	var foundFilm models.Film
	for _, f := range repo.Films {
		if f.ID == id {
			foundFilm = f
			break
		}
	}

	if foundFilm.ID == uuid.Nil {
		errorResp := models.Error{
			Message: "Invalid film id",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	for _, g := range repo.Genres {
		if foundFilm.GenreID == g.ID {
			genre = g.Title
			break
		}
	}

	for _, c := range repo.Countries {
		if foundFilm.CountryID == c.ID {
			country = c.Name
			break
		}
	}

	for _, aif := range repo.ActorsInFilms {
		if aif.FilmID == foundFilm.ID {
			for _, actor := range repo.Actors {
				if actor.ID == aif.ActorID {
					actors = append(actors, actor)
					break
				}
			}
		}
	}

	result = models.FilmPage{
		ID:               foundFilm.ID,
		Title:            foundFilm.Title,
		OriginalTitle:    foundFilm.OriginalTitle,
		Cover:            foundFilm.Cover,
		Poster:           foundFilm.Poster,
		Genre:            genre,
		ShortDescription: foundFilm.ShortDescription,
		Description:      foundFilm.Description,
		AgeCategory:      foundFilm.AgeCategory,
		Budget:           foundFilm.Budget,
		WorldwideFees:    foundFilm.WorldwideFees,
		TrailerURL:       foundFilm.TrailerURL,
		NumerOfRatings:   foundFilm.NumerOfRatings,
		Year:             foundFilm.Year,
		Rating:           foundFilm.Rating,
		Country:          country,
		Slogan:           foundFilm.Slogan,
		Duration:         foundFilm.Duration,
		Image1:           foundFilm.Image1,
		Image2:           foundFilm.Image2,
		Image3:           foundFilm.Image3,
		Actors:           actors,
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

	var films []models.Film
	for _, f := range repo.Films {
		if f.GenreID == neededGenre {
			films = append(films, f)
		}
	}

	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)
	if offset >= len(films) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endingIndex := min(offset+count, len(films))

	paginatedFilms := films[offset:endingIndex]

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(paginatedFilms)
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

	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)
	if offset >= len(genres) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endingIndex := min(offset+count, len(genres))

	paginatedGenres := genres[offset:endingIndex]

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(paginatedGenres)
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

	var films []models.Film
	for _, aif := range repo.ActorsInFilms {
		if aif.ActorID == neededActor {
			for _, f := range repo.Films {
				if aif.FilmID == f.ID {
					films = append(films, f)
					break
				}
			}
		}
	}

	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)
	if offset >= len(films) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	endingIndex := min(offset+count, len(films))

	paginatedFilms := films[offset:endingIndex]

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(paginatedFilms)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetActor godoc
// @Summary      Get actor by ID
// @Tags         actors
// @Produce      json
// @Param        id   path      string  true  "Actor ID"
// @Success      200  {object}  models.Actor
// @Failure      400  {object}  models.Error
// @Router       /actors/{id} [get]
func (c *FilmHandler) GetActor(w http.ResponseWriter, r *http.Request) {
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

	var actor models.Actor
	for _, a := range repo.Actors {
		if id == a.ID {
			actor = a
			break
		}
	}

	endDate := actor.DeathDate
	if actor.DeathDate.IsZero() {
		endDate = time.Now()
	}
	age := endDate.Year() - actor.BirthDate.Year()
	if endDate.YearDay() < actor.BirthDate.YearDay() {
		age--
	}

	filmsNumber := 0
	for _, aif := range repo.ActorsInFilms {
		if aif.ActorID == id {
			filmsNumber++
		}
	}

	result := models.ActorPage{
		ID:            actor.ID,
		RussianName:   actor.RussianName,
		OriginalName:  actor.OriginalName,
		Photo:         actor.Photo,
		Height:        actor.Height,
		BirthDate:     actor.BirthDate,
		Age:           age,
		ZodiacSign:    actor.ZodiacSign,
		BirthPlace:    actor.BirthPlace,
		MaritalStatus: actor.MaritalStatus,
		FilmsNumber:   filmsNumber,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
