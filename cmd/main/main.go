package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	authHandlers "kinopoisk/internal/pkg/auth/handlers"
	filmHandlers "kinopoisk/internal/pkg/film/handlers"
	"kinopoisk/internal/pkg/middleware/cors"
	"kinopoisk/internal/pkg/repo"
	"os"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	godotenv.Load()

	mainRouter := mux.NewRouter()
	fs := http.FileServer(http.Dir("/opt/static/"))
	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r := mainRouter.PathPrefix("/api").Subrouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I am not giving any films!", http.StatusTeapot)
	})

	r.Use(cors.CorsMiddleware)
	repo.InitRepo()

	filmHandler := filmHandlers.NewFilmHandler()
	authHandler := authHandlers.NewAuthHandler()

	// регистрация/авторизация
	r.HandleFunc("/auth/signup", authHandler.SignupUser).Methods(http.MethodPost)
	r.HandleFunc("/auth/signin", authHandler.SignInUser).Methods(http.MethodPost)
	r.Handle("/auth/check", authHandler.Middleware(http.HandlerFunc(authHandler.CheckAuth))).Methods(http.MethodGet)

	// пользователи
	r.HandleFunc("/users/{id}", authHandler.GetUser).Methods(http.MethodGet)

	// фильмы
	r.HandleFunc("/films", filmHandler.GetFilms).Methods(http.MethodGet)
	r.HandleFunc("/films/{id}", filmHandler.GetFilm).Methods(http.MethodGet)
	r.HandleFunc("/films/{genre-id}", filmHandler.GetFilmsByGenre).Methods(http.MethodGet)

	// жанры
	r.HandleFunc("/genres", filmHandler.GetGenres).Methods(http.MethodGet)
	r.HandleFunc("/genres/{id}", filmHandler.GetGenre).Methods(http.MethodGet)

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

	err := filmSrv.Shutdown(ctx)
	if err != nil {
		log.Printf("Graceful shutdown failed")
		os.Exit(1)
	}
	log.Printf("Graceful shutdown!")
}
