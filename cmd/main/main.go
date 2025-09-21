package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"kinopoisk/internal/pkg/film"

	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I am not giving any films!", http.StatusTeapot)
	})
	http.Handle("/", r)

	filmHandler := film.NewFilmHandler()

	// регистрация/авторизация
	r.HandleFunc("/auth/signup", filmHandler.SignupUser).Methods("POST")
	r.HandleFunc("/auth/login", filmHandler.LoginUser).Methods("POST")

	// пользователи
	r.HandleFunc("/users/{id}", filmHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/", filmHandler.CreateUser).Methods("POST")

	// фильмы
	r.HandleFunc("/films", filmHandler.GetFilms).Methods("GET")
	r.HandleFunc("/films/{id}", filmHandler.GetFilm).Methods("GET")

	filmSrv := http.Server{
		Handler: r,
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
	os.Exit(0)

}
