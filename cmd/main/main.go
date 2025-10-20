package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"kinopoisk/dsn"
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
	userHandlers "kinopoisk/internal/pkg/users/delivery/http"
	userRepo "kinopoisk/internal/pkg/users/repo"
	userUsecase "kinopoisk/internal/pkg/users/usecase"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
)

func initDB(ctx context.Context) (*pgxpool.Pool, error) {
	// убрать
	postgresString := dsn.FromEnv()

	config, err := pgxpool.ParseConfig(postgresString)
	if err != nil {
		fmt.Println("bug")
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		fmt.Println("bug")
		return nil, err
	}

	return pool, nil
}

// @title           Kinopoisk API
// @version         1.0
// @description     API для авторизации пользователей и получения фильмов/жанров.
// @host            localhost:5458
// @BasePath        /api
func main() {
	_ = godotenv.Load()
	ctx := context.Background()
	dbpool, err := initDB(ctx)

	mainRouter := mux.NewRouter()
	fs := http.FileServer(http.Dir("/opt/static/"))
	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r := mainRouter.PathPrefix("/api").Subrouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I am not giving any films!", http.StatusTeapot)
	})

	r.Use(cors.CorsMiddleware)

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
	userUsecase := userUsecase.NewUserUsecase(userRepo)
	userHandler := userHandlers.NewUserHandler(userUsecase)

	// регистрация/авторизация
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", authHandler.SignupUser).Methods(http.MethodPost, http.MethodOptions)
	authRouter.HandleFunc("/signin", authHandler.SignInUser).Methods(http.MethodPost, http.MethodOptions)
	authRouter.Handle("/check", authHandler.Middleware(http.HandlerFunc(authHandler.CheckAuth))).Methods(http.MethodGet, http.MethodOptions)
	authRouter.Handle("/change/password", userHandler.Middleware(http.HandlerFunc(userHandler.ChangePassword))).Methods(http.MethodPut, http.MethodOptions)
	authRouter.Handle("/change/avatar", userHandler.Middleware(http.HandlerFunc(userHandler.ChangeAvatar))).Methods(http.MethodPut, http.MethodOptions)
	authRouter.HandleFunc("/change/avatar", userHandler.ChangeAvatar).Methods(http.MethodPut, http.MethodOptions)
	authRouter.Handle("/logout", authHandler.Middleware(http.HandlerFunc(authHandler.LogOutUser))).Methods(http.MethodPost, http.MethodOptions)

	// пользователи
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods(http.MethodGet)

	// фильмы
	r.HandleFunc("/films", filmHandler.GetFilms).Methods(http.MethodGet)
	r.HandleFunc("/films/promo", filmHandler.GetPromoFilm).Methods(http.MethodGet)
	r.HandleFunc("/films/{id}", filmHandler.GetFilm).Methods(http.MethodGet)
	r.HandleFunc("/films/genre/{id}", filmHandler.GetFilmsByGenre).Methods(http.MethodGet)
	r.HandleFunc("/films/actor/{id}", filmHandler.GetFilmsByActor).Methods(http.MethodGet)
	r.HandleFunc("/film/feedbacks/{id}", filmHandler.GetFilmFeedbacks).Methods(http.MethodGet)
	r.Handle("/films/send-feedback/{id}", authHandler.Middleware(http.HandlerFunc(filmHandler.SendFeedback))).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/films/set-rating/{id}", authHandler.Middleware(http.HandlerFunc(filmHandler.SetRating))).Methods(http.MethodPost, http.MethodOptions)

	// жанры
	r.HandleFunc("/genres", genreHandler.GetGenres).Methods(http.MethodGet)
	r.HandleFunc("/genres/{id}", genreHandler.GetGenre).Methods(http.MethodGet)

	// актеры
	r.HandleFunc("/actors/{id}", actorHandler.GetActor).Methods(http.MethodGet)

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
