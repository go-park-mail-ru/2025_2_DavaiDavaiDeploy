package actors

import (
	"context"
	"kinopoisk/internal/models"

	uuid "github.com/satori/go.uuid"
)

type ActorUsecase interface {
	GetActor(ctx context.Context, id uuid.UUID) (models.ActorPage, error)
	GetFilmsByActor(ctx context.Context, id uuid.UUID, pager models.Pager) ([]models.MainPageFilm, error)
}

type ActorRepo interface {
	GetActorByID(ctx context.Context, id uuid.UUID) (models.Actor, error)
	GetActorFilmsCount(ctx context.Context, actorID uuid.UUID) (int, error)
	GetFilmsByActor(ctx context.Context, actorID uuid.UUID, limit, offset int) ([]models.MainPageFilm, error)
}
