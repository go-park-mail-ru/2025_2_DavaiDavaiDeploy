package repo

import (
	"context"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"strconv"

	"github.com/jackc/pgtype/pgxtype"
	uuid "github.com/satori/go.uuid"
)

type ActorRepository struct {
	db pgxtype.Querier
}

func NewActorRepository(db pgxtype.Querier) *ActorRepository {
	return &ActorRepository{db: db}
}

func (r *ActorRepository) GetActorByID(ctx context.Context, id uuid.UUID) (models.Actor, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var actor models.Actor
	err := r.db.QueryRow(
		ctx,
		GetActorByID,
		id,
	).Scan(
		&actor.ID, &actor.RussianName, &actor.OriginalName, &actor.Photo, &actor.Height,
		&actor.BirthDate, &actor.DeathDate, &actor.ZodiacSign, &actor.BirthPlace, &actor.MaritalStatus,
	)
	if err != nil {
		logger.Error("failed to scan actor: " + err.Error())
		return models.Actor{}, err
	}
	return actor, nil
}

func (r *ActorRepository) GetActorFilmsCount(ctx context.Context, actorID uuid.UUID) (int, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var count int
	err := r.db.QueryRow(
		ctx,
		GetActorFilmsCount,
		actorID,
	).Scan(&count)
	if err != nil {
		logger.Error("failed to scan films of actor: " + err.Error())
		return 0, err
	}
	return count, err
}

func (r *ActorRepository) GetFilmAvgRating(ctx context.Context, filmID uuid.UUID) (float64, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var avgRating float64
	err := r.db.QueryRow(
		ctx,
		GetFilmAvgRating,
		filmID,
	).Scan(&avgRating)
	if err != nil {
		logger.Error("failed to scan rating: " + err.Error())
		return 0.0, err
	}
	roundedRating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", avgRating), 64)
	return roundedRating, err
}

func (r *ActorRepository) GetFilmsByActor(ctx context.Context, actorID uuid.UUID, limit, offset int) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	rows, err := r.db.Query(ctx, GetFilmsByActor, actorID, limit, offset)
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
			logger.Error("failed to scan films: " + err.Error())
			return nil, err
		}

		rating, err := r.GetFilmAvgRating(ctx, film.ID)
		if err != nil {
			logger.Error("failed to scan rating of film: " + err.Error())
			film.Rating = 0.0
		} else {
			film.Rating = rating
		}

		films = append(films, film)
	}

	return films, nil
}
