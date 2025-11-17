package filmHandlers

import (
	"context"
	"encoding/json"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/status"
)

var (
	CookieName = "DDFilmsJWT"
)

type FilmHandler struct {
	client gen.FilmsClient
}

func NewFilmHandler(client gen.FilmsClient) *FilmHandler {
	return &FilmHandler{client: client}
}

// GetPromoFilm godoc
// @Summary Get promotional film
// @Description Get the promo film
// @Tags films
// @Produce json
// @Success 200 {object} models.PromoFilm
// @Failure 404
// @Failure 500
// @Router /films/promo [get]
func (c *FilmHandler) GetPromoFilm(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	film, err := c.client.GetPromoFilm(r.Context(), &gen.EmptyRequest{})
	if err != nil {
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			log.LogHandlerError(logger, err, http.StatusNotFound)
			helpers.WriteError(w, http.StatusNotFound)
		default:
			log.LogHandlerError(logger, err, http.StatusInternalServerError)
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := models.PromoFilm{
		ID:               uuid.FromStringOrNil(film.Id),
		Image:            film.Image,
		Title:            film.Title,
		Rating:           film.Rating,
		ShortDescription: film.ShortDescription,
		Year:             int(film.Year),
		Genre:            film.Genre,
		Duration:         int(film.Duration),
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (c *FilmHandler) GetUsersFavFilms(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	user, _ := r.Context().Value(auth.UserKey).(models.User)
	favFilms, err := c.client.GetFavFilms(r.Context(), &gen.GetFavFilmsRequest{UserId: user.ID.String()})
	if err != nil {
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			log.LogHandlerError(logger, err, http.StatusNotFound)
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			log.LogHandlerError(logger, err, http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			log.LogHandlerError(logger, err, http.StatusInternalServerError)
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := []models.FavFilm{}
	for i := range favFilms.Films {
		var film models.FavFilm
		film.ID = uuid.FromStringOrNil(favFilms.Films[i].Id)
		film.Title = favFilms.Films[i].Title
		film.Image = favFilms.Films[i].Image
		film.Rating = favFilms.Films[i].Rating
		film.Genre = favFilms.Films[i].Genre
		film.Year = int(favFilms.Films[i].Year)
		film.Duration = int(favFilms.Films[i].Duration)
		film.ShortDescription = favFilms.Films[i].ShortDescription
		response = append(response, film)
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetFilms godoc
// @Summary      List films
// @Tags         films
// @Produce      json
// @Param        count   query     int  false  "Number of films" default(10)
// @Param        offset  query     int  false  "Offset" default(0)
// @Success      200     {array}   models.MainPageFilm
// @Failure      400
// @Failure 	 404
// @Failure 	 500
// @Router       /films [get]
func (c *FilmHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	pager := helpers.GetPagerFromRequest(r)

	mainPageFilms, err := c.client.GetFilms(r.Context(), &gen.GetFilmsRequest{
		Pager: &gen.Pager{Count: int32(pager.Count), Offset: int32(pager.Offset)},
	})

	if err != nil {
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			log.LogHandlerError(logger, err, http.StatusNotFound)
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			log.LogHandlerError(logger, err, http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			log.LogHandlerError(logger, err, http.StatusInternalServerError)
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := []models.MainPageFilm{}
	for i := range mainPageFilms.Films {
		var film models.MainPageFilm
		film.ID = uuid.FromStringOrNil(mainPageFilms.Films[i].Id)
		film.Cover = mainPageFilms.Films[i].Cover
		film.Title = mainPageFilms.Films[i].Title
		film.Rating = mainPageFilms.Films[i].Rating
		film.Genre = mainPageFilms.Films[i].Genre
		film.Year = int(mainPageFilms.Films[i].Year)
		response = append(response, film)
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (c *FilmHandler) GetFilmsForCalendar(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	pager := helpers.GetPagerFromRequest(r)

	user, _ := r.Context().Value(auth.UserKey).(models.User)

	filmsForCalendar, err := c.client.GetFilmsForCalendar(r.Context(), &gen.GetFilmsForCalendarRequest{
		Pager:  &gen.Pager{Count: int32(pager.Count), Offset: int32(pager.Offset)},
		UserId: user.ID.String(),
	})

	if err != nil {
		st, _ := status.FromError(err)

		switch st.Code() {
		case codes.NotFound:
			log.LogHandlerError(logger, err, http.StatusNotFound)
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			log.LogHandlerError(logger, err, http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			log.LogHandlerError(logger, err, http.StatusInternalServerError)
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := []models.FilmInCalendar{}
	for i := range filmsForCalendar.Films {
		var film models.FilmInCalendar

		film.ID = uuid.FromStringOrNil(filmsForCalendar.Films[i].ID)
		film.Cover = filmsForCalendar.Films[i].Cover
		film.Title = filmsForCalendar.Films[i].Title
		film.IsLiked = filmsForCalendar.Films[i].IsLiked

		if filmsForCalendar.Films[i].OriginalTitle != nil {
			film.OriginalTitle = filmsForCalendar.Films[i].OriginalTitle
		}

		film.ShortDescription = filmsForCalendar.Films[i].ShortDescription

		releaseDate, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", filmsForCalendar.Films[i].ReleaseDate)
		if err != nil {
			log.LogHandlerError(logger, err, http.StatusInternalServerError)
			helpers.WriteError(w, http.StatusInternalServerError)
			return
		}

		film.ReleaseDate = releaseDate
		response = append(response, film)
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetFilm godoc
// @Summary      Get film by ID
// @Tags         films
// @Produce      json
// @Param        id   path      string  true  "Film ID"
// @Success      200  {object}  models.FilmPage
// @Failure      400
// @Failure 	 404
// @Failure 	 500
// @Router       /films/{id} [get]
func (c *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	user, _ := r.Context().Value(auth.UserKey).(models.User)

	film, err := c.client.GetFilm(r.Context(), &gen.GetFilmRequest{
		FilmId: id.String(),
		UserId: user.ID.String(),
	})

	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			log.LogHandlerError(logger, err, http.StatusNotFound)
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			log.LogHandlerError(logger, err, http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			log.LogHandlerError(logger, err, http.StatusInternalServerError)
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	actors := make([]models.Actor, 0, len(film.Actors))
	for _, actor := range film.Actors {
		mappedActor := models.Actor{
			ID:            uuid.FromStringOrNil(actor.Id),
			Photo:         actor.Photo,
			Height:        int(actor.Height),
			ZodiacSign:    actor.ZodiacSign,
			BirthPlace:    actor.BirthPlace,
			MaritalStatus: actor.MaritalStatus,
			RussianName:   actor.RussianName,
		}

		if actor.OriginalName != nil {
			mappedActor.OriginalName = actor.OriginalName
		}

		if birthDate, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", actor.BirthDate); err == nil {
			mappedActor.BirthDate = birthDate
		}
		if actor.DeathDate != nil {
			if deathDate, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", *actor.DeathDate); err == nil {
				mappedActor.DeathDate = &deathDate
			}
		}

		actors = append(actors, mappedActor)
	}

	response := models.FilmPage{
		ID:               uuid.FromStringOrNil(film.Id),
		Title:            film.Title,
		Cover:            film.Cover,
		Poster:           film.Poster,
		Genre:            film.Genre,
		ShortDescription: film.ShortDescription,
		Description:      film.Description,
		AgeCategory:      film.AgeCategory,
		Budget:           int(film.Budget),
		WorldwideFees:    int(film.WorldwideFees),
		TrailerURL:       film.TrailerUrl,
		NumberOfRatings:  int(film.NumberOfRatings),
		Year:             int(film.Year),
		Rating:           film.Rating,
		Country:          film.Country,
		Slogan:           film.Slogan,
		Duration:         int(film.Duration),
		Image1:           film.Image1,
		Image2:           film.Image2,
		Image3:           film.Image3,
		Actors:           actors,
		IsReviewed:       film.IsReviewed,
		IsLiked:          film.IsLiked,
		GenreID:          uuid.FromStringOrNil(film.GenreId),
	}

	if film.UserRating != nil {
		userRating := int(*film.UserRating)
		response.UserRating = &userRating
	}

	if film.OriginalTitle != nil {
		response.OriginalTitle = film.OriginalTitle
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (c *FilmHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
		var token string
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
			logger.Info("successfully got token")
		}

		if token != "" {
			user, err := c.client.ValidateUser(r.Context(), &gen.ValidateUserRequest{Token: token})
			if err == nil {
				neededUser := models.User{
					ID:      uuid.FromStringOrNil(user.ID),
					Version: int(user.Version),
					Login:   user.Login,
					Avatar:  user.Avatar,
				}
				ctx := context.WithValue(r.Context(), auth.UserKey, neededUser)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				log.LogHandlerError(logger, err, http.StatusUnauthorized)
			}
		}

		log.LogHandlerInfo(logger, "success", http.StatusOK)
		next.ServeHTTP(w, r)
	})
}

func (c *FilmHandler) SaveFilm(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	user, ok := r.Context().Value(auth.UserKey).(models.User)
	if !ok {
		log.LogHandlerError(logger, errors.New("user unauthorized"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	_, err = c.client.SaveFilm(r.Context(), &gen.SaveFilmRequest{UserId: user.ID.String(), FilmId: filmID.String()})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (c *FilmHandler) RemoveFilm(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	user, ok := r.Context().Value(auth.UserKey).(models.User)
	if !ok {
		log.LogHandlerError(logger, errors.New("user unauthorized"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	_, err = c.client.RemoveFilm(r.Context(), &gen.RemoveFilmRequest{UserId: user.ID.String(), FilmId: filmID.String()})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (c *FilmHandler) GetFilmFeedbacks(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	pager := helpers.GetPagerFromRequest(r)
	user, _ := r.Context().Value(auth.UserKey).(models.User)

	feedbacks, err := c.client.GetFilmFeedbacks(r.Context(), &gen.GetFilmFeedbacksRequest{
		FilmId: id.String(),
		Pager:  &gen.Pager{Count: int32(pager.Count), Offset: int32(pager.Offset)},
		UserId: user.ID.String(),
	})

	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := []models.FilmFeedback{}
	for i := range feedbacks.Feedbacks {
		feedback := models.FilmFeedback{
			ID:            uuid.FromStringOrNil(feedbacks.Feedbacks[i].Id),
			UserID:        uuid.FromStringOrNil(feedbacks.Feedbacks[i].UserId),
			FilmID:        uuid.FromStringOrNil(feedbacks.Feedbacks[i].FilmId),
			Rating:        int(feedbacks.Feedbacks[i].Rating),
			UserLogin:     feedbacks.Feedbacks[i].UserLogin,
			UserAvatar:    feedbacks.Feedbacks[i].UserAvatar,
			IsMine:        feedbacks.Feedbacks[i].IsMine,
			NewFilmRating: feedbacks.Feedbacks[i].NewFilmRating,
		}

		if createdAt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", feedbacks.Feedbacks[i].CreatedAt); err == nil {
			feedback.CreatedAt = createdAt
		}
		if updatedAt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", feedbacks.Feedbacks[i].UpdatedAt); err == nil {
			feedback.UpdatedAt = updatedAt
		}

		if feedbacks.Feedbacks[i].Title != nil {
			feedback.Title = feedbacks.Feedbacks[i].Title
		}
		if feedbacks.Feedbacks[i].Text != nil {
			feedback.Text = feedbacks.Feedbacks[i].Text
		}

		response = append(response, feedback)
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// SendFeedback godoc
// @Summary Add film review
// @Tags films
// @Accept json
// @Produce json
// @Param        id   path      string  true  "Film ID"
// @Param input body models.FilmFeedbackInput true "Feedback data"
// @Success 200 {object} models.FilmFeedback
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /films/{id}/feedbacks [post]
func (c *FilmHandler) SendFeedback(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	user, ok := r.Context().Value(auth.UserKey).(models.User)
	if !ok {
		log.LogHandlerError(logger, errors.New("user unauthorized"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	var req models.FilmFeedbackInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid request"), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Sanitize()

	feedback, err := c.client.SendFeedback(r.Context(), &gen.SendFeedbackRequest{
		FilmId: filmID.String(),
		UserId: user.ID.String(),
		Feedback: &gen.FilmFeedbackInput{
			Title:  req.Title,
			Text:   req.Text,
			Rating: int32(req.Rating),
		},
	})

	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := models.FilmFeedback{
		ID:            uuid.FromStringOrNil(feedback.Feedback.Id),
		UserID:        uuid.FromStringOrNil(feedback.Feedback.UserId),
		FilmID:        uuid.FromStringOrNil(feedback.Feedback.FilmId),
		Rating:        int(feedback.Feedback.Rating),
		UserLogin:     feedback.Feedback.UserLogin,
		UserAvatar:    feedback.Feedback.UserAvatar,
		IsMine:        feedback.Feedback.IsMine,
		NewFilmRating: feedback.Feedback.NewFilmRating,
	}

	if createdAt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", feedback.Feedback.CreatedAt); err == nil {
		response.CreatedAt = createdAt
	}
	if updatedAt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", feedback.Feedback.UpdatedAt); err == nil {
		response.UpdatedAt = updatedAt
	}

	if feedback.Feedback.Title != nil {
		response.Title = feedback.Feedback.Title
	}
	if feedback.Feedback.Text != nil {
		response.Text = feedback.Feedback.Text
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// SetRating godoc
// @Summary Rate a film
// @Tags films
// @Accept json
// @Produce json
// @Param        id   path      string  true  "Film ID"
// @Param input body gen.FilmRatingInput true "Rating data (rating 1-10 is required)"
// @Success 200 {object} models.FilmFeedback
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /films/{id}/rating [post]
func (c *FilmHandler) SetRating(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	user, ok := r.Context().Value(auth.UserKey).(models.User)
	if !ok {
		log.LogHandlerError(logger, errors.New("user unauthorized"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	filmID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of film"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	var req models.FilmFeedbackInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid request"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	rating, err := c.client.SetRating(r.Context(), &gen.SetRatingRequest{
		FilmId: filmID.String(),
		UserId: user.ID.String(),
		RatingInput: &gen.FilmRatingInput{
			Rating: int32(req.Rating),
		},
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := models.FilmFeedback{
		ID:            uuid.FromStringOrNil(rating.Feedback.Id),
		UserID:        uuid.FromStringOrNil(rating.Feedback.UserId),
		FilmID:        uuid.FromStringOrNil(rating.Feedback.FilmId),
		Rating:        int(rating.Feedback.Rating),
		UserLogin:     rating.Feedback.UserLogin,
		UserAvatar:    rating.Feedback.UserAvatar,
		IsMine:        rating.Feedback.IsMine,
		NewFilmRating: rating.Feedback.NewFilmRating,
	}

	if createdAt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", rating.Feedback.CreatedAt); err == nil {
		response.CreatedAt = createdAt
	}
	if updatedAt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", rating.Feedback.UpdatedAt); err == nil {
		response.UpdatedAt = updatedAt
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
