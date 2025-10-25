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

func (r *ActorRepository) GetActorByID(ctx context.Context, id uuid.UUID) (models.Actor, error) {
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
		return models.Actor{}, err
	}
	return actor, nil
}

func (r *ActorRepository) GetActorFilmsCount(ctx context.Context, actorID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(
		ctx,
		"SELECT COUNT(1) FROM actor_in_film WHERE actor_id = $1",
		actorID,
	).Scan(&count)
	return count, err
}

func (r *ActorRepository) GetFilmsByActor(ctx context.Context, actorID uuid.UUID, limit, offset int) ([]models.MainPageFilm, error) {
	query := `
        SELECT 
            f.id, 
            COALESCE(f.cover, ''), 
            f.title, 
            f.year,
            g.title as genre
        FROM film f
        JOIN actor_in_film aif ON f.id = aif.film_id
        JOIN genre g ON f.genre_id = g.id
        WHERE aif.actor_id = $1
        ORDER BY f.created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, actorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []models.MainPageFilm
	for rows.Next() {
		var film models.MainPageFilm
		if err := rows.Scan(
			&film.ID,
			&film.Cover,
			&film.Title,
			&film.Year,
			&film.Genre,
		); err != nil {
			continue
		}
		films = append(films, film)
	}
	return films, nil
}
