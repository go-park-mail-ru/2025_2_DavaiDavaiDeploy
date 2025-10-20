package repo

import (
	"context"
	"kinopoisk/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type GenreRepository struct {
	db *pgxpool.Pool
}

func NewGenreRepository(db *pgxpool.Pool) *GenreRepository {
	return &GenreRepository{db: db}
}

func (g *GenreRepository) GetGenreByID(ctx context.Context, id uuid.UUID) (models.Genre, error) {
	var genre models.Genre
	err := g.db.QueryRow(
		ctx,
		"SELECT id, title, description, icon, created_at, updated_at FROM genre WHERE id = $1",
		id,
	).Scan(
		&genre.ID, &genre.Title, &genre.Description, &genre.Icon,
		&genre.CreatedAt, &genre.UpdatedAt,
	)
	if err != nil {
		return models.Genre{}, err
	}
	return genre, nil
}

func (g *GenreRepository) GetGenresWithPagination(ctx context.Context, limit, offset int) ([]models.Genre, error) {
	query := `
        SELECT id, title, description, icon, created_at, updated_at 
        FROM genre 
        ORDER BY title
        LIMIT $1 OFFSET $2`

	rows, err := g.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(
			&genre.ID, &genre.Title, &genre.Description, &genre.Icon,
			&genre.CreatedAt, &genre.UpdatedAt,
		); err != nil {
			continue
		}
		genres = append(genres, genre)
	}
	return genres, nil
}
