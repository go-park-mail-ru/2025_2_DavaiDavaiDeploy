package http

import (
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ActorHandler struct {
	client gen.FilmsClient
}

func NewActorHandler(client gen.FilmsClient) *ActorHandler {
	return &ActorHandler{client: client}
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

	actor, err := a.client.GetActor(r.Context(), &gen.GetActorRequest{ActorId: id.String()})
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

	response := models.ActorPage{
		ID:            uuid.FromStringOrNil(actor.Actor.Id),
		RussianName:   actor.Actor.RussianName,
		Photo:         actor.Actor.Photo,
		Height:        int(actor.Actor.Height),
		Age:           int(actor.Actor.Age),
		ZodiacSign:    actor.Actor.ZodiacSign,
		BirthPlace:    actor.Actor.BirthPlace,
		MaritalStatus: actor.Actor.MaritalStatus,
		FilmsNumber:   int(actor.Actor.FilmsNumber),
	}

	if actor.Actor.OriginalName != nil {
		response.OriginalName = actor.Actor.OriginalName
	}

	if birthDate, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", actor.Actor.BirthDate); err == nil {
		response.BirthDate = birthDate
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
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
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of actor"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	films, err := a.client.GetFilmsByActor(r.Context(), &gen.GetFilmsByActorRequest{
		ActorId: id.String(),
		Pager: &gen.Pager{
			Count:  int32(pager.Count),
			Offset: int32(pager.Offset),
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

	response := []models.MainPageFilm{}
	for i := range films.Films {
		film := models.MainPageFilm{
			ID:     uuid.FromStringOrNil(films.Films[i].Id),
			Cover:  films.Films[i].Cover,
			Title:  films.Films[i].Title,
			Rating: films.Films[i].Rating,
			Year:   int(films.Films[i].Year),
			Genre:  films.Films[i].Genre,
		}
		response = append(response, film)
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
