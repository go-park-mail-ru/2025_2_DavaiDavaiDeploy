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

func (r *S3Repository) UploadAvatar(ctx context.Context, userID string, buffer []byte, fileFormat string, avatarExtension string) (string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if r.client == nil || r.bucket == "" {
		return "", errors.New("S3 client not configured")
	}

	avatarKey := filepath.Join("static", "avatars", userID+avatarExtension)
	avatarDBKey := filepath.Join("avatars", userID+avatarExtension)

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
