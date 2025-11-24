package repo

import (
	"context"
	"errors"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/compilations"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"strconv"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
)

type CompilationRepository struct {
	db pgxtype.Querier
}

func NewCompilationRepository(db pgxtype.Querier) *CompilationRepository {
	return &CompilationRepository{db: db}
}

func (g *CompilationRepository) GetCompilationByID(ctx context.Context, id uuid.UUID) (models.Compilation, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var compilation models.Compilation
	err := g.db.QueryRow(
		ctx,
		GetCompilationByIDQuery,
		id,
	).Scan(
		&compilation.ID, &compilation.Title, &compilation.Description, &compilation.Icon,
		&compilation.CreatedAt, &compilation.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("compilation is not found: " + err.Error())
			return models.Compilation{}, compilations.ErrorNotFound
		}
		logger.Error("failed to scan compilation: " + err.Error())
		return models.Compilation{}, compilations.ErrorInternalServerError
	}

	logger.Info("succesfully got compilation by id from db")
	return compilation, nil
}

func (g *CompilationRepository) GetCompilationsWithPagination(ctx context.Context, limit, offset int) ([]models.Compilation, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if limit <= 0 || offset < 0 {
		return nil, compilations.ErrorBadRequest
	}

	rows, err := g.db.Query(ctx, GetCompilationsWithPaginationQuery, limit, offset)
	if err != nil {
		logger.Error("failed to get rows: " + err.Error())
		return nil, compilations.ErrorInternalServerError
	}
	defer rows.Close()

	var compilations []models.Compilation
	for rows.Next() {
		var compilation models.Compilation
		if err := rows.Scan(
			&compilation.ID, &compilation.Title, &compilation.Description, &compilation.Icon,
			&compilation.CreatedAt, &compilation.UpdatedAt,
		); err != nil {
			logger.Error("failed to scan compilation: " + err.Error())
			continue
		}
		compilations = append(compilations, compilation)
	}

	logger.Info("succesfully got compilations from db")
	return compilations, nil
}

func (g *CompilationRepository) GetFilmsByCompilation(ctx context.Context, compilationID uuid.UUID, limit, offset int) ([]models.FavFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	rows, err := g.db.Query(ctx, GetFilmsByCompilationQuery, compilationID, limit, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("compilation is not found: " + err.Error())
			return nil, compilations.ErrorNotFound
		}
		logger.Error("failed to scan compilation: " + err.Error())
		return nil, compilations.ErrorInternalServerError
	}
	defer rows.Close()

	var films []models.FavFilm
	for rows.Next() {
		var film models.FavFilm
		if err := rows.Scan(
			&film.ID,
			&film.Title,
			&film.Genre,
			&film.Year,
			&film.Duration,
			&film.Image,
			&film.ShortDescription,
			&film.Rating,
		); err != nil {
			logger.Error("failed to scan film: " + err.Error())
			continue
		}
		rating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", film.Rating), 64)
		film.Rating = rating
		films = append(films, film)
	}

	logger.Info("succesfully got films by compilation from db")
	return films, nil
}
