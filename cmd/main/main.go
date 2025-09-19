package main

import (
	"log"
	"net/http"

	"kinopoisk/internal/pkg/auth"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I am not giving any films!", http.StatusTeapot)
	})
	http.Handle("/", r)

	filmController, err := auth.NewFilmController()
	if err != nil {
		log.Fatal(err)
	}

	// регистрация/авторизация
	r.HandleFunc("/auth/signup", filmController.SignupUser).Methods("POST")
	r.HandleFunc("/auth/login", filmController.LoginUser).Methods("POST")

	// пользователи
	r.HandleFunc("/users/{id}", filmController.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", filmController.CreateUser).Methods("POST")

	// фильмы
	r.HandleFunc("/films", filmController.GetFilms).Methods("GET")
	r.HandleFunc("/films/{id}", filmController.GetFilm).Methods("GET")

	filmSrv := http.Server{
		Handler: r,
		Addr:    ":5458",
	}

	err = filmSrv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Start error: %v", err)
	}
}
