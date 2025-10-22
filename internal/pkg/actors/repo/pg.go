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

func (r *ActorRepository) GetFilmsByActor(ctx context.Context, actorID uuid.UUID, limit, offset int) ([]models.Film, error) {
	query := `
        SELECT 
            f.id, f.title, f.original_title, f.cover, f.poster,
            f.short_description, f.description, f.age_category, f.budget,
            f.worldwide_fees, f.trailer_url, f.year, f.country_id,
            f.genre_id, f.slogan, f.duration, f.image1, f.image2,
            f.image3, f.created_at, f.updated_at
        FROM film f
        JOIN actor_in_film aif ON f.id = aif.film_id
        WHERE aif.actor_id = $1
        ORDER BY f.created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, actorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []models.Film
	for rows.Next() {
		var film models.Film
		if err := rows.Scan(
			&film.ID, &film.Title, &film.OriginalTitle, &film.Cover, &film.Poster,
			&film.ShortDescription, &film.Description, &film.AgeCategory, &film.Budget,
			&film.WorldwideFees, &film.TrailerURL, &film.Year, &film.CountryID,
			&film.GenreID, &film.Slogan, &film.Duration, &film.Image1, &film.Image2,
			&film.Image3, &film.CreatedAt, &film.UpdatedAt,
		); err != nil {
			continue
		}
		films = append(films, film)
	}
	return films, nil
}
