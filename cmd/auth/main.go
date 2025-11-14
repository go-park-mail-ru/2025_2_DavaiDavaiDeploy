// @title           Kinopoisk API
// @version         1.0
// @description     API для авторизации пользователей и получения фильмов/жанров/актеров.
// @host            localhost:5458
// @BasePath        /api
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	authHandler "kinopoisk/internal/pkg/auth/delivery/grpc"
	authRepo "kinopoisk/internal/pkg/auth/repo"
	authUsecase "kinopoisk/internal/pkg/auth/usecase"

	_ "kinopoisk/docs"

	"kinopoisk/internal/pkg/auth/delivery/grpc/gen"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func initDB(ctx context.Context) (*pgxpool.Pool, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	postgresString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	config, err := pgxpool.ParseConfig(postgresString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func initS3Client(ctx context.Context) (*s3.Client, string, error) {
	endpoint := os.Getenv("AWS_S3_ENDPOINT")
	bucket := os.Getenv("AWS_S3_BUCKET")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	region := "ru-7"

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && endpoint != "" {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	customHTTPClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Пропускаем проверку SSL
			},
		},
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithHTTPClient(customHTTPClient),
	)
	if err != nil {
		return nil, "", err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	return client, bucket, nil
}

func main() {
	_ = godotenv.Load()

	authRepo := authRepo.NewAuthRepository(nil)
	authUsecase := authUsecase.NewAuthUsecase(authRepo)
	authHandler := authHandler.NewGrpcAuthHandler(authUsecase)

	gRPCServer := grpc.NewServer() //
	gen.RegisterAuthServer(gRPCServer, authHandler)

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(": %d", 8080))
		if err != nil {
			fmt.Println(err)
		}
		if err := gRPCServer.Serve(listener); err != nil {
			fmt.Println(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	gRPCServer.GracefulStop()

}
