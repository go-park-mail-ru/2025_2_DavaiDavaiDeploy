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
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	actorHandlers "kinopoisk/internal/pkg/actors/delivery/http"
	authHandlers "kinopoisk/internal/pkg/auth/delivery/http"
	filmHandlers "kinopoisk/internal/pkg/films/delivery/http"
	genreHandlers "kinopoisk/internal/pkg/genres/delivery/http"
	"kinopoisk/internal/pkg/middleware/cors"
	logger "kinopoisk/internal/pkg/middleware/logger"
	searchHandlers "kinopoisk/internal/pkg/search/delivery/http"
	userHandlers "kinopoisk/internal/pkg/users/delivery/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/gorilla/mux"

	_ "kinopoisk/docs"

	authGen "kinopoisk/internal/pkg/auth/delivery/grpc/gen"
	filmGen "kinopoisk/internal/pkg/films/delivery/grpc/gen"
	searchGen "kinopoisk/internal/pkg/search/delivery/grpc/gen"

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

	authConn, err := grpc.Dial("auth:5459", grpc.WithInsecure())
	if err != nil {
		log.Printf("unable to connect to auth microservice: %v\n", err)
		return
	}
	defer authConn.Close()

	filmConn, err := grpc.Dial("films:5460", grpc.WithInsecure())
	if err != nil {
		log.Printf("unable to connect to films microservice: %v\n", err)
		return
	}
	defer filmConn.Close()

	searchConn, err := grpc.Dial("films:5462", grpc.WithInsecure())
	if err != nil {
		log.Printf("unable to connect to search microservice: %v\n", err)
		return
	}
	defer searchConn.Close()

	filmClient := filmGen.NewFilmsClient(filmConn)
	authClient := authGen.NewAuthClient(authConn)
	searchClient := searchGen.NewSearchClient(searchConn)

	authHandler := authHandlers.NewAuthHandler(authClient)
	userHandler := userHandlers.NewUserHandler(authClient)
	genreHandler := genreHandlers.NewGenreHandler(filmClient)
	actorHandler := actorHandlers.NewActorHandler(filmClient)
	filmHandler := filmHandlers.NewFilmHandler(filmClient)
	searchHandler := searchHandlers.NewSearchHandler(searchClient)

	ddLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mainRouter := mux.NewRouter()
	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	mainRouter.PathPrefix("/metrics").Handler(promhttp.Handler())

	apiRouter := mainRouter.PathPrefix("/api").Subrouter()
	apiRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I am not giving any films!", http.StatusTeapot)
	})

	apiRouter.Use(cors.CorsMiddleware)
	apiRouter.Use(logger.LoggerMiddleware(ddLogger))

	apiRouter.HandleFunc("/sitemap.xml", filmHandler.SiteMap).Methods(http.MethodGet)

	apiRouter.HandleFunc("/search", searchHandler.GetFilmsAndActorsFromSearch).Methods(http.MethodGet)
	// Auth routes
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", authHandler.SignupUser).Methods(http.MethodPost, http.MethodOptions)
	authRouter.HandleFunc("/signin", authHandler.SignInUser).Methods(http.MethodPost, http.MethodOptions)

	protectedAuthRouter := authRouter.PathPrefix("").Subrouter()
	protectedAuthRouter.Use(authHandler.Middleware)
	protectedAuthRouter.HandleFunc("/check", authHandler.CheckAuth).Methods(http.MethodGet, http.MethodOptions)
	protectedAuthRouter.HandleFunc("/logout", authHandler.LogOutUser).Methods(http.MethodPost, http.MethodOptions)
	protectedAuthRouter.HandleFunc("/enable2fa", authHandler.Enable2FA).Methods(http.MethodPost, http.MethodOptions)
	protectedAuthRouter.HandleFunc("/disable2fa", authHandler.Disable2FA).Methods(http.MethodPost, http.MethodOptions)

	// User routes
	userRouter := apiRouter.PathPrefix("/users").Subrouter()

	// Protected user routes
	protectedUserRouter := userRouter.PathPrefix("").Subrouter()
	protectedUserRouter.Use(userHandler.Middleware)
	protectedUserRouter.HandleFunc("/change/password", userHandler.ChangePassword).Methods(http.MethodPut, http.MethodOptions)
	protectedUserRouter.HandleFunc("/change/avatar", userHandler.ChangeAvatar).Methods(http.MethodPut, http.MethodOptions)
	protectedUserRouter.HandleFunc("/saved", filmHandler.GetUsersFavFilms).Methods(http.MethodGet)

	userRouter.HandleFunc("/{id}", userHandler.GetUser).Methods(http.MethodGet)

	// Film routes
	filmRouter := apiRouter.PathPrefix("/films").Subrouter()
	filmRouter.Use(filmHandler.Middleware)
	filmRouter.HandleFunc("/", filmHandler.GetFilms).Methods(http.MethodGet)
	filmRouter.HandleFunc("/promo", filmHandler.GetPromoFilm).Methods(http.MethodGet)
	filmRouter.HandleFunc("/calendar", filmHandler.GetFilmsForCalendar).Methods(http.MethodGet)
	filmRouter.HandleFunc("/{id}", filmHandler.GetFilm).Methods(http.MethodGet)
	filmRouter.HandleFunc("/{id}/feedbacks", filmHandler.GetFilmFeedbacks).Methods(http.MethodGet)

	// Protected film routes
	protectedFilmRouter := filmRouter.PathPrefix("").Subrouter()
	protectedFilmRouter.Use(authHandler.Middleware)
	protectedFilmRouter.HandleFunc("/{id}/feedback", filmHandler.SendFeedback).Methods(http.MethodPost, http.MethodOptions)
	protectedFilmRouter.HandleFunc("/{id}/rating", filmHandler.SetRating).Methods(http.MethodPost, http.MethodOptions)
	protectedFilmRouter.HandleFunc("/{id}/save", filmHandler.SaveFilm).Methods(http.MethodPost, http.MethodOptions)
	protectedFilmRouter.HandleFunc("/{id}/remove", filmHandler.RemoveFilm).Methods(http.MethodDelete, http.MethodOptions)

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
