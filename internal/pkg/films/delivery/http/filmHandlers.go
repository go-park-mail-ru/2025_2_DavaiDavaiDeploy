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
	helpers.WriteJSON(w, film)
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
	helpers.WriteJSON(w, mainPageFilms.Films)
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
	helpers.WriteJSON(w, film)
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

	var req models.SaveFilmInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid request"), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = c.client.SaveFilm(r.Context(), &gen.SaveFilmRequest{UserId: user.ID.String(), FilmId: req.FilmID.String()})
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

	var req models.RemoveFilmInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid request"), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = c.client.RemoveFilm(r.Context(), &gen.RemoveFilmRequest{UserId: user.ID.String(), FilmId: req.FilmID.String()})
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

	helpers.WriteJSON(w, feedbacks.Feedbacks)
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
	helpers.WriteJSON(w, feedback.Feedback)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// SetRating godoc
// @Summary Rate a film
// @Tags films
// @Accept json
// @Produce json
// @Param        id   path      string  true  "Film ID"
// @Param input body models.FilmRatingInput true "Rating data (rating 1-10 is required)"
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
	helpers.WriteJSON(w, rating.Feedback)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
