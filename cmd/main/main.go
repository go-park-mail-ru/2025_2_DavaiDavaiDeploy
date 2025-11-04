// @title           Kinopoisk API
// @version         1.0
// @description     API для авторизации пользователей и получения фильмов/жанров/актеров.
// @host            localhost:5458
// @BasePath        /api
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	actorHandlers "kinopoisk/internal/pkg/actors/delivery/http"
	actorRepo "kinopoisk/internal/pkg/actors/repo"
	actorUsecase "kinopoisk/internal/pkg/actors/usecase"
	authHandlers "kinopoisk/internal/pkg/auth/delivery/http"
	authRepo "kinopoisk/internal/pkg/auth/repo"
	authUsecase "kinopoisk/internal/pkg/auth/usecase"
	filmHandlers "kinopoisk/internal/pkg/films/delivery/http"
	filmRepo "kinopoisk/internal/pkg/films/repo"
	filmUsecase "kinopoisk/internal/pkg/films/usecase"
	genreHandlers "kinopoisk/internal/pkg/genres/delivery/http"
	genreRepo "kinopoisk/internal/pkg/genres/repo"
	genreUsecase "kinopoisk/internal/pkg/genres/usecase"
	"kinopoisk/internal/pkg/middleware/cors"
	logger "kinopoisk/internal/pkg/middleware/logger"
	userHandlers "kinopoisk/internal/pkg/users/delivery/http"
	userRepo "kinopoisk/internal/pkg/users/repo/pg"
	storageRepo "kinopoisk/internal/pkg/users/repo/s3"
	userUsecase "kinopoisk/internal/pkg/users/usecase"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"

	_ "kinopoisk/docs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	httpSwagger "github.com/swaggo/http-swagger"
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

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
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
	ctx := context.Background()
	dbpool, err := initDB(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	s3Client, s3Bucket, err := initS3Client(ctx)
	if err != nil {
		log.Printf("Warning: Unable to connect to S3: %v\n", err)
	}

	ddLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mainRouter := mux.NewRouter()
	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	apiRouter := mainRouter.PathPrefix("/api").Subrouter()
	apiRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I am not giving any films!", http.StatusTeapot)
	})

	apiRouter.Use(cors.CorsMiddleware)
	apiRouter.Use(logger.LoggerMiddleware(ddLogger))

	// Инициализация репозиториев, usecases и handlers
	filmRepo := filmRepo.NewFilmRepository(dbpool)
	filmUsecase := filmUsecase.NewFilmUsecase(filmRepo)
	filmHandler := filmHandlers.NewFilmHandler(filmUsecase)

	authRepo := authRepo.NewAuthRepository(dbpool)
	authUsecase := authUsecase.NewAuthUsecase(authRepo)
	authHandler := authHandlers.NewAuthHandler(authUsecase)

	genreRepo := genreRepo.NewGenreRepository(dbpool)
	genreUsecase := genreUsecase.NewGenreUsecase(genreRepo)
	genreHandler := genreHandlers.NewGenreHandler(genreUsecase)

	actorRepo := actorRepo.NewActorRepository(dbpool)
	actorUsecase := actorUsecase.NewActorUsecase(actorRepo)
	actorHandler := actorHandlers.NewActorHandler(actorUsecase)

	userRepo := userRepo.NewUserRepository(dbpool)
	s3Repo := storageRepo.NewS3Repository(s3Client, s3Bucket)
	userUsecase := userUsecase.NewUserUsecase(userRepo, s3Repo)
	userHandler := userHandlers.NewUserHandler(userUsecase)

	// Auth routes
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", authHandler.SignupUser).Methods(http.MethodPost, http.MethodOptions)
	authRouter.HandleFunc("/signin", authHandler.SignInUser).Methods(http.MethodPost, http.MethodOptions)

	protectedAuthRouter := authRouter.PathPrefix("").Subrouter()
	protectedAuthRouter.Use(authHandler.Middleware)
	protectedAuthRouter.HandleFunc("/check", authHandler.CheckAuth).Methods(http.MethodGet, http.MethodOptions)
	protectedAuthRouter.HandleFunc("/logout", authHandler.LogOutUser).Methods(http.MethodPost, http.MethodOptions)

	// User routes
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/{id}", userHandler.GetUser).Methods(http.MethodGet)

	// Protected user routes
	protectedUserRouter := userRouter.PathPrefix("/change").Subrouter()
	protectedUserRouter.Use(userHandler.Middleware)
	protectedUserRouter.HandleFunc("/password", userHandler.ChangePassword).Methods(http.MethodPut, http.MethodOptions)
	protectedUserRouter.HandleFunc("/avatar", userHandler.ChangeAvatar).Methods(http.MethodPut, http.MethodOptions)

	// Film routes
	filmRouter := apiRouter.PathPrefix("/films").Subrouter()
	filmRouter.Use(filmHandler.Middleware)
	filmRouter.HandleFunc("/", filmHandler.GetFilms).Methods(http.MethodGet)
	filmRouter.HandleFunc("/promo", filmHandler.GetPromoFilm).Methods(http.MethodGet)
	filmRouter.HandleFunc("/{id}", filmHandler.GetFilm).Methods(http.MethodGet)
	filmRouter.HandleFunc("/{id}/feedbacks", filmHandler.GetFilmFeedbacks).Methods(http.MethodGet)

	// Protected film routes
	protectedFilmRouter := filmRouter.PathPrefix("").Subrouter()
	protectedFilmRouter.Use(authHandler.Middleware)
	protectedFilmRouter.HandleFunc("/{id}/feedback", filmHandler.SendFeedback).Methods(http.MethodPost, http.MethodOptions)
	protectedFilmRouter.HandleFunc("/{id}/rating", filmHandler.SetRating).Methods(http.MethodPost, http.MethodOptions)

	// Genre routes
	genreRouter := apiRouter.PathPrefix("/genres").Subrouter()
	genreRouter.HandleFunc("/", genreHandler.GetGenres).Methods(http.MethodGet)
	genreRouter.HandleFunc("/{id}", genreHandler.GetGenre).Methods(http.MethodGet)
	genreRouter.HandleFunc("/{id}/films", genreHandler.GetFilmsByGenre).Methods(http.MethodGet)

	// Actor routes
	actorRouter := apiRouter.PathPrefix("/actors").Subrouter()
	actorRouter.HandleFunc("/{id}", actorHandler.GetActor).Methods(http.MethodGet)
	actorRouter.HandleFunc("/{id}/films", actorHandler.GetFilmsByActor).Methods(http.MethodGet)

	filmSrv := http.Server{
		Handler: mainRouter,
		Addr:    ":5458",
	}

	go func() {
		log.Println("Starting server!")
		err := filmSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Server start error: %v", err)
			os.Exit(1)
		}
	}()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	<-quitChannel
	log.Printf("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = filmSrv.Shutdown(ctx)
	if err != nil {
		log.Printf("Graceful shutdown failed")
		os.Exit(1)
	}
	log.Printf("Graceful shutdown!!")
}
