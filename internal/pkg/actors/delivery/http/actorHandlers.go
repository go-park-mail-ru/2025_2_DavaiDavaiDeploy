package http

import (
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors/repo"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type ActorHandler struct {
	actorRepo *repo.ActorRepository
}

func NewActorHandler(db *pgxpool.Pool) *ActorHandler {
	actorRepo := repo.NewActorRepository(db)
	return &ActorHandler{actorRepo: actorRepo}
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
		errorResp := models.Error{
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	actor, err := a.actorRepo.GetActorByID(r.Context(), id)

	endDate := actor.DeathDate
	if actor.DeathDate.IsZero() {
		endDate = time.Now()
	}
	age := endDate.Year() - actor.BirthDate.Year()
	if endDate.YearDay() < actor.BirthDate.YearDay() {
		age--
	}

	filmsNumber, err := a.actorRepo.GetActorFilmsCount(r.Context(), id)

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
