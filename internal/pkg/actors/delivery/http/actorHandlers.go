package http

import (
	"kinopoisk/internal/pkg/actors"
	"kinopoisk/internal/pkg/helpers"
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
// @Success      200  {object}  models.Actor
// @Failure      400  {object}  models.Error
// @Router       /actors/{id} [get]
func (a *ActorHandler) GetActor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	actor, err := a.uc.GetActor(r.Context(), id)
	if err != nil {
		helpers.WriteError(w, 400, err)
		return
	}

	helpers.WriteJSON(w, actor)
}
