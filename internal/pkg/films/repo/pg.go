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

func (r *FilmRepository) GetPromoFilmByID(ctx context.Context, id uuid.UUID) (models.PromoFilm, error) {
	var film models.PromoFilm

	err := r.db.QueryRow(
		ctx,
		GetPromoFilmByIDQuery,
		id,
	).Scan(
		&film.ID,
		&film.Image,
		&film.Title,
		&film.ShortDescription,
		&film.Year,
		&film.Genre,
		&film.Duration,
		&film.CreatedAt,
		&film.UpdatedAt,
	)

	if err != nil {
		return models.PromoFilm{}, fmt.Errorf("failed to get promo film: %w", err)
	}

	return film, nil
}

func (r *FilmRepository) GetFilmByID(ctx context.Context, id uuid.UUID) (models.Film, error) {
	var film models.Film
	err := r.db.QueryRow(
		ctx,
		GetFilmByIDQuery,
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
		GetGenreTitleQuery,
		genreID,
	).Scan(&title)
	return title, err
}

func (r *FilmRepository) GetFilmAvgRating(ctx context.Context, filmID uuid.UUID) (float64, error) {
	var avgRating float64
	err := r.db.QueryRow(
		ctx,
		GetFilmAvgRatingQuery,
		filmID,
	).Scan(&avgRating)
	roundedRating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", avgRating), 64)
	return roundedRating, err
}

func (r *FilmRepository) GetFilmsWithPagination(ctx context.Context, limit, offset int) ([]models.MainPageFilm, error) {
	rows, err := r.db.Query(ctx, GetFilmsWithPaginationQuery, limit, offset)
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
		rating, err := r.GetFilmAvgRating(ctx, film.ID)
		if err != nil {
			film.Rating = 0.0
		} else {
			film.Rating = rating
		}
		films = append(films, film)
	}
	return films, nil
}

func (r *FilmRepository) GetFilmPage(ctx context.Context, filmID uuid.UUID) (models.FilmPage, error) {
	var result models.FilmPage
	err := r.db.QueryRow(ctx, GetFilmPageQuery, filmID).Scan(
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

	rows, err := r.db.Query(ctx, GetFilmActorsQuery, filmID)
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

func (r *FilmRepository) GetFilmFeedbacks(ctx context.Context, filmID uuid.UUID, limit, offset int) ([]models.FilmFeedback, error) {
	rows, err := r.db.Query(ctx, GetFilmFeedbacksQuery, filmID, limit, offset)
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
			&feedback.UserLogin, &feedback.UserAvatar,
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
		CheckUserFeedbackExistsQuery,
		userID, filmID,
	).Scan(
		&feedback.ID, &feedback.UserID, &feedback.FilmID, &feedback.Title,
		&feedback.Text, &feedback.Rating, &feedback.CreatedAt, &feedback.UpdatedAt,
		&feedback.UserLogin, &feedback.UserAvatar,
	)
	if err != nil {
		return models.FilmFeedback{}, err
	}
	return feedback, nil
}

func (r *FilmRepository) UpdateFeedback(ctx context.Context, feedback models.FilmFeedback) error {
	_, err := r.db.Exec(
		ctx,
		UpdateFeedbackQuery,
		feedback.Title, feedback.Text, feedback.Rating, feedback.ID,
	)
	return err
}

func (r *FilmRepository) CreateFeedback(ctx context.Context, feedback models.FilmFeedback) error {
	_, err := r.db.Exec(
		ctx,
		CreateFeedbackQuery,
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
		SetRatingQuery,
		feedback.ID, feedback.UserID, feedback.FilmID, feedback.Rating,
	)
	return err
}

func (r *FilmRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		ctx,
		GetUserByLoginQuery,
		login,
	).Scan(
		&user.ID, &user.Version, &user.Login,
		&user.PasswordHash, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL Error: %s, Code: %s, Detail: %s\n",
				pgErr.Message, pgErr.Code, pgErr.Detail)
		}

		fmt.Printf("Error getting user by ID %v: %v\n", user, err)
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
