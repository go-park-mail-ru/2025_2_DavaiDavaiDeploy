package http

import (
	"errors"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GenreHandler struct {
	client gen.FilmsClient
}

func NewGenreHandler(client gen.FilmsClient) *GenreHandler {
	return &GenreHandler{client: client}
}

// GetGenre godoc
// @Summary Get genre by ID
// @Tags genres
// @Produce json
// @Param        id   path      string  true  "Genre ID"
// @Success 200 {object} models.Genre
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /genres/{id} [get]
func (g *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of genre"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	genre, err := g.client.GetGenre(r.Context(), &gen.GetGenreRequest{GenreId: id.String()})
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
	helpers.WriteJSON(w, genre.Genre)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetGenres godoc
// @Summary Get list of all genres
// @Tags genres
// @Produce json
// @Success 200 {array} models.Genre
// @Failure 404
// @Failure 500
// @Router /genres [get]
func (g *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	pager := helpers.GetPagerFromRequest(r)

	genres, err := g.client.GetGenres(r.Context(), &gen.GetGenresRequest{
		Pager: &gen.Pager{Count: int32(pager.Count), Offset: int32(pager.Offset)},
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}
	helpers.WriteJSON(w, genres.Genres)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetFilmsByGenre godoc
// @Summary Get films by genre
// @Tags genres
// @Produce json
// @Param        id   path      string  true  "Genre ID"
// @Success 200 {array} models.MainPageFilm
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /genres/{id}/films [get]
func (g *GenreHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of genre"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	films, err := g.client.GetFilmsByGenre(r.Context(), &gen.GetFilmsByGenreRequest{
		GenreId: id.String(),
		Pager:   &gen.Pager{Count: int32(pager.Count), Offset: int32(pager.Offset)},
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
	helpers.WriteJSON(w, films.Films)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
