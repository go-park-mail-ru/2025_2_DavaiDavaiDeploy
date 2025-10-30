package repo

import (
	"context"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"strconv"

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
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
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
		logger.Error("failed to scan actor: " + err.Error())
		return models.Actor{}, err
	}

	logger.Info("Successfully got actor")
	return actor, nil
}

func (r *ActorRepository) GetActorFilmsCount(ctx context.Context, actorID uuid.UUID) (int, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var count int
	err := r.db.QueryRow(
		ctx,
		"SELECT COUNT(1) FROM actor_in_film WHERE actor_id = $1",
		actorID,
	).Scan(&count)

	if err != nil {
		logger.Error("failed to scan films: " + err.Error())
		return 0, err
	}

	logger.Info("Successfully got number of actors films")
	return count, err
}

func (r *ActorRepository) GetFilmAvgRating(ctx context.Context, filmID uuid.UUID) (float64, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var avgRating float64
	err := r.db.QueryRow(
		ctx,
		"SELECT COALESCE(AVG(rating), 0) FROM film_feedback WHERE film_id = $1",
		filmID,
	).Scan(&avgRating)
	if err != nil {
		logger.Error("failed to scan rating" + err.Error())
		return 0, err
	}

	roundedRating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", avgRating), 64)
	logger.Info("Successfully got rating")
	return roundedRating, err
}

func (r *ActorRepository) GetFilmsByActor(ctx context.Context, actorID uuid.UUID, limit, offset int) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
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
		logger.Error("failed to get films by actor from db" + err.Error())
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
			logger.Error("failed to scan films by actor" + err.Error())
			return nil, err
		}

		rating, err := r.GetFilmAvgRating(ctx, film.ID)
		if err != nil {
			logger.Error("failed to get rating" + err.Error())
			film.Rating = 0.0
		} else {
			film.Rating = rating
		}

		films = append(films, film)
	}

	return films, nil
}
