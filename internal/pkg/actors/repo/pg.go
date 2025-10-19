package repo

import (
	"context"
	"kinopoisk/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type ActorRepository struct {
	db *pgxpool.Pool
}

func NewActorRepository(db *pgxpool.Pool) *ActorRepository {
	return &ActorRepository{db: db}
}

func (r *ActorRepository) GetActorByID(ctx context.Context, id uuid.UUID) (*models.Actor, error) {
	var actor models.Actor
	err := r.db.QueryRow(
		ctx,
		"SELECT id, russian_name, original_name, photo, height, birth_date, death_date, zodiac_sign, birth_place, marital_status FROM actor WHERE id = $1",
		id,
	).Scan(
		&actor.ID, &actor.RussianName, &actor.OriginalName, &actor.Photo, &actor.Height,
		&actor.BirthDate, &actor.DeathDate, &actor.ZodiacSign, &actor.BirthPlace, &actor.MaritalStatus,
	)
	if err != nil {
		return nil, err
	}
	return &actor, nil
}

func (r *ActorRepository) GetActorFilmsCount(ctx context.Context, actorID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM actor_in_film WHERE actor_id = $1",
		actorID,
	).Scan(&count)
	return count, err
}
