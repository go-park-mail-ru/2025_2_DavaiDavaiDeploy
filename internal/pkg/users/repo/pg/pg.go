package repo

import (
	"context"
	"errors"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/users"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
)

type UserRepository struct {
	db pgxtype.Querier
}

func NewUserRepository(db pgxtype.Querier) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var user models.User
	err := u.db.QueryRow(
		ctx,
		GetUserByIDQuery,
		id,
	).Scan(
		&user.ID, &user.Version, &user.Login,
		&user.PasswordHash, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("user not exists")
			return models.User{}, users.ErrorNotFound
		}
		logger.Error("failed to scan user: " + err.Error())
		return models.User{}, users.ErrorInternalServerError
	}

	logger.Info("succesfully got user by id from db")
	return user, nil
}

func (u *UserRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var user models.User
	err := u.db.QueryRow(
		ctx,
		GetUserByLoginQuery,
		login,
	).Scan(
		&user.ID, &user.Version, &user.Login,
		&user.PasswordHash, &user.Avatar, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("user not exists")
			return models.User{}, users.ErrorNotFound
		}
		logger.Error("failed to scan user: " + err.Error())
		return models.User{}, users.ErrorInternalServerError
	}

	logger.Info("succesfully got user by login from db")
	return user, nil
}

func (u *UserRepository) UpdateUserPassword(ctx context.Context, version int, userID uuid.UUID, passwordHash []byte) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	_, err := u.db.Exec(
		ctx,
		UpdateUserPasswordQuery,
		passwordHash, version, userID,
	)
	if err != nil {
		logger.Error("failed to update password: " + err.Error())
		return users.ErrorInternalServerError
	}

	logger.Info("succesfully updated password of user from db")
	return err
}

func (u *UserRepository) UpdateUserAvatar(ctx context.Context, version int, userID uuid.UUID, avatarPath string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	_, err := u.db.Exec(
		ctx,
		UpdateUserAvatarQuery,
		avatarPath, version, userID,
	)
	if err != nil {
		logger.Error("failed to update avatar: " + err.Error())
		return users.ErrorInternalServerError
	}

	logger.Info("succesfully updated avatar from db")
	return nil
}

func (u *UserRepository) CreateFeedback(ctx context.Context, feedback *models.SupportFeedback) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	err := u.db.QueryRow(
		ctx,
		CreateFeedbackQuery,
		feedback.UserID,
		feedback.Description,
		feedback.Category,
		feedback.Attachment,
	).Scan(
		&feedback.ID,
		&feedback.CreatedAt,
		&feedback.UpdatedAt,
	)
	if err != nil {
		logger.Error("failed to create feedback: " + err.Error())
		return users.ErrorInternalServerError
	}

	logger.Info("successfully created feedback")
	return nil
}

func (u *UserRepository) GetFeedbackByID(ctx context.Context, id uuid.UUID) (models.SupportFeedback, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	var feedback models.SupportFeedback
	err := u.db.QueryRow(
		ctx,
		GetFeedbackByIDQuery,
		id,
	).Scan(
		&feedback.ID,
		&feedback.UserID,
		&feedback.Description,
		&feedback.Category,
		&feedback.Status,
		&feedback.Attachment,
		&feedback.CreatedAt,
		&feedback.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("feedback not found")
			return models.SupportFeedback{}, users.ErrorNotFound
		}
		logger.Error("failed to get feedback: " + err.Error())
		return models.SupportFeedback{}, users.ErrorInternalServerError
	}

	logger.Info("successfully got feedback by id")
	return feedback, nil
}

func (u *UserRepository) GetFeedbacksByUserID(ctx context.Context, userID uuid.UUID) ([]models.SupportFeedback, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	rows, err := u.db.Query(
		ctx,
		GetFeedbacksByUserIDQuery,
		userID,
	)
	if err != nil {
		logger.Error("failed to query feedbacks: " + err.Error())
		return nil, users.ErrorInternalServerError
	}
	defer rows.Close()

	var feedbacks []models.SupportFeedback
	for rows.Next() {
		var feedback models.SupportFeedback
		err := rows.Scan(
			&feedback.ID,
			&feedback.UserID,
			&feedback.Description,
			&feedback.Category,
			&feedback.Status,
			&feedback.Attachment,
			&feedback.CreatedAt,
			&feedback.UpdatedAt,
		)
		if err != nil {
			logger.Error("failed to scan feedback: " + err.Error())
			return nil, users.ErrorInternalServerError
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err = rows.Err(); err != nil {
		logger.Error("error iterating feedback rows: " + err.Error())
		return nil, users.ErrorInternalServerError
	}

	logger.Info("successfully got feedbacks by user id", "count", len(feedbacks))
	return feedbacks, nil
}

func (u *UserRepository) UpdateFeedback(ctx context.Context, feedback *models.SupportFeedback) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	err := u.db.QueryRow(
		ctx,
		UpdateFeedbackQuery,
		feedback.Description,
		feedback.Category,
		feedback.Status,
		feedback.Attachment,
		feedback.ID,
	).Scan(
		&feedback.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println(err)
			logger.Error("feedback not found for update")
			return users.ErrorNotFound
		}
		logger.Error("failed to update feedback: " + err.Error())
		return users.ErrorInternalServerError
	}

	logger.Info("successfully updated feedback")
	return nil
}

func (u *UserRepository) GetFeedbackStats(ctx context.Context) (models.FeedbackStats, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	var stats models.FeedbackStats
	err := u.db.QueryRow(
		ctx,
		GetFeedbackStatsQuery,
	).Scan(
		&stats.Total,
		&stats.Open,
		&stats.InProgress,
		&stats.Closed,
		&stats.Bugs,
		&stats.FeatureReqs,
		&stats.Complaints,
		&stats.Questions,
	)
	if err != nil {
		logger.Error("failed to get feedback stats: " + err.Error())
		return models.FeedbackStats{}, users.ErrorInternalServerError
	}

	logger.Info("successfully got feedback stats")
	return stats, nil
}

func (u *UserRepository) GetUserFeedbackStats(ctx context.Context, userID uuid.UUID) (models.FeedbackStats, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	var stats models.FeedbackStats
	err := u.db.QueryRow(
		ctx,
		GetUserFeedbackStatsQuery,
		userID,
	).Scan(
		&stats.Total,
		&stats.Open,
		&stats.InProgress,
		&stats.Closed,
		&stats.Bugs,
		&stats.FeatureReqs,
		&stats.Complaints,
		&stats.Questions,
	)
	if err != nil {
		logger.Error("failed to get user feedback stats: " + err.Error())
		return models.FeedbackStats{}, users.ErrorInternalServerError
	}

	logger.Info("successfully got user feedback stats")
	return stats, nil
}
