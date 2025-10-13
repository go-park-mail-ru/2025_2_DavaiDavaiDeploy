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

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Kinopoisk API
// @version         1.0
// @description     API для авторизации пользователей и получения фильмов/жанров.
// @host            localhost:5458
// @BasePath        /api
func main() {
	_ = godotenv.Load()

	mainRouter := mux.NewRouter()
	fs := http.FileServer(http.Dir("/opt/static/"))
	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r := mainRouter.PathPrefix("/api").Subrouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I am not giving any films!", http.StatusTeapot)
	})

	r.Use(cors.CorsMiddleware)
	repo.InitRepo()

	filmHandler := filmHandlers.NewFilmHandler()
	authHandler := authHandlers.NewAuthHandler()

	// регистрация/авторизация
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/signup", authHandler.SignupUser).Methods(http.MethodPost, http.MethodOptions)
	authRouter.HandleFunc("/signin", authHandler.SignInUser).Methods(http.MethodPost, http.MethodOptions)
	authRouter.Handle("/check", authHandler.Middleware(http.HandlerFunc(authHandler.CheckAuth))).Methods(http.MethodGet, http.MethodOptions)
	authRouter.HandleFunc("/change/password", authHandler.ChangePassword).Methods(http.MethodPut, http.MethodOptions)
	authRouter.HandleFunc("/change/avatar", authHandler.ChangeAvatar).Methods(http.MethodPut, http.MethodOptions)
	authRouter.Handle("/logout", authHandler.Middleware(http.HandlerFunc(authHandler.LogOutUser))).Methods(http.MethodGet, http.MethodOptions)

	// пользователи
	r.HandleFunc("/users/{id}", authHandler.GetUser).Methods(http.MethodGet)

	// фильмы
	r.HandleFunc("/films", filmHandler.GetFilms).Methods(http.MethodGet)
	r.HandleFunc("/films/promo", filmHandler.GetPromoFilm).Methods(http.MethodGet)
	r.HandleFunc("/films/{id}", filmHandler.GetFilm).Methods(http.MethodGet)
	r.HandleFunc("/films/genre/{id}", filmHandler.GetFilmsByGenre).Methods(http.MethodGet)
	r.HandleFunc("/films/actor/{id}", filmHandler.GetFilmsByActor).Methods(http.MethodGet)
	r.HandleFunc("/films/feedback/{id}", filmHandler.GetFilmsFeedback).Methods(http.MethodGet)
	r.Handle("/films/send-feedback/{id}", authHandler.Middleware(http.HandlerFunc(filmHandler.SendFeedback))).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/films/set-rating/{id}", authHandler.Middleware(http.HandlerFunc(filmHandler.SetRating))).Methods(http.MethodPost, http.MethodOptions)

	// жанры
	r.HandleFunc("/genres", filmHandler.GetGenres).Methods(http.MethodGet)
	r.HandleFunc("/genres/{id}", filmHandler.GetGenre).Methods(http.MethodGet)

	// актеры
	r.HandleFunc("/actors/{id}", filmHandler.GetActor).Methods(http.MethodGet)

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
	log.Printf("Graceful shutdown!!")
}
