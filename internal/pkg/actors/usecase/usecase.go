package usecase

import (
	"context"
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type ActorUsecase struct {
	actorRepo actors.ActorRepo
}

func NewActorUsecase(repo actors.ActorRepo) *ActorUsecase {
	return &ActorUsecase{
		actorRepo: repo,
	}
}

func (uc *ActorUsecase) GetActor(ctx context.Context, id uuid.UUID) (models.ActorPage, error) {
	actor, err := uc.actorRepo.GetActorByID(ctx, id)
	if err != nil {
		return models.ActorPage{}, errors.New("actor not exists")
	}

	var endDate time.Time

	if actor.DeathDate == nil || actor.DeathDate.IsZero() {
		endDate = time.Now()
	} else {
		endDate = *actor.DeathDate
	}

	age := endDate.Year() - actor.BirthDate.Year()
	if endDate.YearDay() < actor.BirthDate.YearDay() {
		age--
	}

	filmsNumber, err := uc.actorRepo.GetActorFilmsCount(ctx, id)
	if err != nil {
		return models.ActorPage{}, errors.New("no films")
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
	return result, nil
}
