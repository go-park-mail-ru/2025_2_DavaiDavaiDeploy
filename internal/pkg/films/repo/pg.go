package repo

import (
	"context"
	"errors"
	"fmt"
	"kinopoisk/internal/models"
	"strconv"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type FilmRepository struct {
	db *pgxpool.Pool
}

func NewFilmRepository(db *pgxpool.Pool) *FilmRepository {
	return &FilmRepository{db: db}
}

func (r *FilmRepository) GetFilmByID(ctx context.Context, id uuid.UUID) (models.Film, error) {
	var film models.Film
	err := r.db.QueryRow(
		ctx,
		`SELECT 
			id, title, original_title, cover, poster,
			short_description, description, age_category, budget,
			worldwide_fees, trailer_url, year, country_id,
			genre_id, slogan, duration, image1, image2,
			image3, created_at, updated_at
		FROM film WHERE id = $1`,
		id,
	).Scan(
		&film.ID, &film.Title, &film.OriginalTitle, &film.Cover, &film.Poster,
		&film.ShortDescription, &film.Description, &film.AgeCategory, &film.Budget,
		&film.WorldwideFees, &film.TrailerURL, &film.Year, &film.CountryID,
		&film.GenreID, &film.Slogan, &film.Duration, &film.Image1, &film.Image2,
		&film.Image3, &film.CreatedAt, &film.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL Error: %s, Code: %s, Detail: %s\n",
				pgErr.Message, pgErr.Code, pgErr.Detail)
		}

		fmt.Printf("Error getting film by ID %s: %v\n", id, err)
		return models.Film{}, fmt.Errorf("failed to get film: %w", err)
	}
	return film, nil
}

func (r *FilmRepository) GetGenreTitle(ctx context.Context, genreID uuid.UUID) (string, error) {
	var title string
	err := r.db.QueryRow(
		ctx,
		"SELECT title FROM genre WHERE id = $1",
		genreID,
	).Scan(&title)
	return title, err
}

func (r *FilmRepository) GetFilmAvgRating(ctx context.Context, filmID uuid.UUID) (float64, error) {
	var avgRating float64
	err := r.db.QueryRow(
		ctx,
		"SELECT COALESCE(AVG(rating), 0) FROM film_feedback WHERE film_id = $1",
		filmID,
	).Scan(&avgRating)
	roundedRating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", avgRating), 64)
	return roundedRating, err
}

func (r *FilmRepository) GetFilmsWithPagination(ctx context.Context, limit, offset int) ([]models.MainPageFilm, error) {
	query := `
        SELECT 
            f.id, f.cover, f.title, f.year, g.title as genre_title,
            COALESCE(AVG(ff.rating), 0) as avg_rating
        FROM film f
        JOIN genre g ON f.genre_id = g.id
        LEFT JOIN film_feedback ff ON f.id = ff.film_id
        GROUP BY f.id, g.title
        ORDER BY f.created_at DESC
        LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
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
			&film.Rating,
		); err != nil {
			continue
		}
		films = append(films, film)
	}
	return films, nil
}

func (r *FilmRepository) GetFilmPage(ctx context.Context, filmID uuid.UUID) (models.FilmPage, error) {
	filmsQuery := `
        SELECT 
            f.id, f.title, f.original_title, f.cover, f.poster,
            f.short_description, f.description, f.age_category, f.budget,
            f.worldwide_fees, f.trailer_url, f.year, 
            f.slogan, f.duration, f.image1, f.image2, f.image3,
            g.title as genre, c.name as country,
            COUNT(ff.id) as number_of_ratings
        FROM film f
        JOIN genre g ON f.genre_id = g.id
        JOIN country c ON f.country_id = c.id
        LEFT JOIN film_feedback ff ON f.id = ff.film_id
        WHERE f.id = $1
        GROUP BY f.id, g.title, c.name`

	var result models.FilmPage
	err := r.db.QueryRow(ctx, filmsQuery, filmID).Scan(
		&result.ID, &result.Title, &result.OriginalTitle, &result.Cover, &result.Poster,
		&result.ShortDescription, &result.Description, &result.AgeCategory, &result.Budget,
		&result.WorldwideFees, &result.TrailerURL, &result.Year,
		&result.Slogan, &result.Duration, &result.Image1, &result.Image2, &result.Image3,
		&result.Genre, &result.Country, &result.NumberOfRatings,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL Error: %s, Code: %s, Detail: %s\n",
				pgErr.Message, pgErr.Code, pgErr.Detail)
		}

		fmt.Printf("Error getting film by ID %s: %v\n", filmID, err)
		return models.FilmPage{}, fmt.Errorf("failed to get film: %w", err)
	}

	result.Rating, err = r.GetFilmAvgRating(ctx, filmID)
	if err != nil {
		result.Rating = 0
	}

	actorsQuery := `
        SELECT a.id, a.russian_name, a.original_name, a.photo, a.height,
               a.birth_date, a.death_date, a.zodiac_sign, a.birth_place, a.marital_status
        FROM actor a
        JOIN actor_in_film aif ON a.id = aif.actor_id
        WHERE aif.film_id = $1`

	rows, err := r.db.Query(ctx, actorsQuery, filmID)
	if err != nil {
		return result, nil
	}
	defer rows.Close()

	for rows.Next() {
		var actor models.Actor
		if err := rows.Scan(
			&actor.ID, &actor.RussianName, &actor.OriginalName, &actor.Photo, &actor.Height,
			&actor.BirthDate, &actor.DeathDate, &actor.ZodiacSign, &actor.BirthPlace, &actor.MaritalStatus,
		); err == nil {
			result.Actors = append(result.Actors, actor)
		}
	}

	return result, nil
}

func (r *FilmRepository) GetFilmsByGenre(ctx context.Context, genreID uuid.UUID, limit, offset int) ([]models.Film, error) {
	query := `
        SELECT 
            id, title, original_title, cover, poster,
            short_description, description, age_category, budget,
            worldwide_fees, trailer_url, year, country_id,
            genre_id, slogan, duration, image1, image2,
            image3, created_at, updated_at
        FROM film 
        WHERE genre_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, genreID, limit, offset)
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

func (r *FilmRepository) GetFilmsByActor(ctx context.Context, actorID uuid.UUID, limit, offset int) ([]models.Film, error) {
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

func (r *FilmRepository) GetFilmFeedbacks(ctx context.Context, filmID uuid.UUID, limit, offset int) ([]models.FilmFeedback, error) {
	query := `
        SELECT 
            id, user_id, film_id, title, text, rating, created_at, updated_at
        FROM film_feedback 
        WHERE film_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, filmID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []models.FilmFeedback
	for rows.Next() {
		var feedback models.FilmFeedback
		if err := rows.Scan(
			&feedback.ID, &feedback.UserID, &feedback.FilmID, &feedback.Title,
			&feedback.Text, &feedback.Rating, &feedback.CreatedAt, &feedback.UpdatedAt,
		); err != nil {
			continue
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL Error: %s, Code: %s, Detail: %s\n",
				pgErr.Message, pgErr.Code, pgErr.Detail)
		}

		fmt.Printf("Error getting feedbacks by film %v\n", err)
		return []models.FilmFeedback{}, fmt.Errorf("failed to get film: %w", err)
	}

	return feedbacks, nil
}

func (r *FilmRepository) CheckUserFeedbackExists(ctx context.Context, userID, filmID uuid.UUID) (models.FilmFeedback, error) {
	var feedback models.FilmFeedback
	err := r.db.QueryRow(
		ctx,
		"SELECT id, user_id, film_id, title, text, rating, created_at, updated_at FROM film_feedback WHERE user_id = $1 AND film_id = $2",
		userID, filmID,
	).Scan(
		&feedback.ID, &feedback.UserID, &feedback.FilmID, &feedback.Title,
		&feedback.Text, &feedback.Rating, &feedback.CreatedAt, &feedback.UpdatedAt,
	)
	if err != nil {
		return models.FilmFeedback{}, err
	}
	return feedback, nil
}

func (r *FilmRepository) UpdateFeedback(ctx context.Context, feedback models.FilmFeedback) error {
	_, err := r.db.Exec(
		ctx,
		"UPDATE film_feedback SET title = $1, text = $2, rating = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4",
		feedback.Title, feedback.Text, feedback.Rating, feedback.ID,
	)
	return err
}

func (r *FilmRepository) CreateFeedback(ctx context.Context, feedback models.FilmFeedback) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO film_feedback (id, user_id, film_id, title, text, rating) VALUES ($1, $2, $3, $4, $5, $6)",
		feedback.ID, feedback.UserID, feedback.FilmID, feedback.Title, feedback.Text, feedback.Rating,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL Error: %s, Code: %s, Detail: %s\n",
				pgErr.Message, pgErr.Code, pgErr.Detail)
		}

		fmt.Printf("Error getting film by ID %v\n", err)
		return fmt.Errorf("failed to get film: %w", err)
	}
	return err
}

func (r *FilmRepository) SetRating(ctx context.Context, feedback models.FilmFeedback) error {
	_, err := r.db.Exec(
		ctx,
		"INSERT INTO film_feedback (id, user_id, film_id, rating) VALUES ($1, $2, $3, $4)",
		feedback.ID, feedback.UserID, feedback.FilmID, feedback.Rating,
	)
	return err
}
