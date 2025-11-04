package http

import (
	"errors"
	"kinopoisk/internal/pkg/actors"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type ActorHandler struct {
	uc actors.ActorUsecase
}

func NewActorHandler(uc actors.ActorUsecase) *ActorHandler {
	return &ActorHandler{uc: uc}
}

// GetActor godoc
// @Summary      Get actor by ID
// @Tags         actors
// @Produce      json
// @Param        id   path      string  true  "Actor ID"
// @Success      200  {object}  models.ActorPage
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /actors/{id} [get]
func (a *ActorHandler) GetActor(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of actor"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	actor, err := a.uc.GetActor(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, actors.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		case errors.Is(err, actors.ErrorInternalServerError):
			helpers.WriteError(w, http.StatusInternalServerError)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}
	actor.Sanitize()
	helpers.WriteJSON(w, actor)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}

// GetFilmsByActor godoc
// @Summary      Get films by actor ID
// @Tags         actors
// @Produce      json
// @Param        id   path      string  true  "Actor ID"
// @Success      200  {array}   models.MainPageFilm
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /actors/{id}/films [get]
func (a *ActorHandler) GetFilmsByActor(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	idStr := vars["id"]

	neededActor, err := uuid.FromString(idStr)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of actor"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	films, err := a.uc.GetFilmsByActor(r.Context(), neededActor, pager)
	if err != nil {
		switch {
		case errors.Is(err, actors.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		case errors.Is(err, actors.ErrorInternalServerError):
			helpers.WriteError(w, http.StatusInternalServerError)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}
	for i := range films {
		films[i].Sanitize()
	}

	helpers.WriteJSON(w, films)
	log.LogHandlerInfo(logger, "Success", http.StatusOK)
}
