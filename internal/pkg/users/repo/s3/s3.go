package repo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	uuid "github.com/satori/go.uuid"
)

type S3Repository struct {
	client *s3.Client
	bucket string
}

func NewS3Repository(client *s3.Client, bucket string) *S3Repository {
	return &S3Repository{
		client: client,
		bucket: bucket,
	}
}

func (r *S3Repository) DeleteAvatar(ctx context.Context, avatarPath string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if r.client == nil || r.bucket == "" {
		return errors.New("S3 client not configured")
	}

	s3Key := filepath.Join("static", avatarPath)

	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		logger.Warn("failed to delete old avatar from S3", "error", err, "avatar_path", avatarPath)
		return fmt.Errorf("failed to delete old avatar: %w", err)
	}
	logger.Info("successfully deleted old avatar from S3", "avatar_path", avatarPath)
	return nil
}

func (r *S3Repository) UploadAvatar(ctx context.Context, userID string, buffer []byte, fileFormat string, avatarExtension string) (string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if r.client == nil || r.bucket == "" {
		return "", errors.New("S3 client not configured")
	}

	picID := uuid.NewV4().String()

	avatarKey := filepath.Join("static", "avatars", picID+avatarExtension)
	avatarDBKey := filepath.Join("avatars", picID+avatarExtension)

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(avatarKey),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(fileFormat),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		logger.Error("failed to upload avatar to S3", "error", err)
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	return avatarDBKey, nil
}

// UploadFeedbackAttachment загружает вложение для фидбэка в S3
func (r *S3Repository) UploadFeedbackAttachment(ctx context.Context, feedbackID string, buffer []byte, fileFormat string, fileExtension string) (string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if r.client == nil || r.bucket == "" {
		return "", errors.New("S3 client not configured")
	}

	attachmentKey := filepath.Join("static", "feedback", feedbackID+fileExtension)
	attachmentDBKey := filepath.Join("feedback", feedbackID+fileExtension)

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(attachmentKey),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(fileFormat),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		logger.Error("failed to upload feedback attachment to S3", "error", err)
		return "", fmt.Errorf("failed to upload feedback attachment: %w", err)
	}

	logger.Info("successfully uploaded feedback attachment to S3", "feedback_id", feedbackID)
	return attachmentDBKey, nil
}

func (r *S3Repository) DeleteFeedbackAttachment(ctx context.Context, attachmentPath string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if r.client == nil || r.bucket == "" {
		return errors.New("S3 client not configured")
	}

	s3Key := filepath.Join("static", attachmentPath)

	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		logger.Warn("failed to delete feedback attachment from S3", "error", err, "attachment_path", attachmentPath)
		return fmt.Errorf("failed to delete feedback attachment: %w", err)
	}
	logger.Info("successfully deleted feedback attachment from S3", "attachment_path", attachmentPath)
	return nil
}
