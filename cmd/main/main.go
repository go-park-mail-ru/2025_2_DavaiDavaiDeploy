package main

import (
	"kinopoisk/internal/pkg/auth"
	"kinopoisk/internal/pkg/film"
	"log"
	"net/http"
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
	authHandler := auth.NewAuthHandler()

	// регистрация/авторизация
	r.HandleFunc("/auth/signup", authHandler.SignupUser).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.LoginUser).Methods("POST")

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

	err := filmSrv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("Server start error: %v", err)
		os.Exit(1)
	}
}
