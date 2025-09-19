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

	r.HandleFunc("/login", filmController.LoginUser).Methods("POST")
	r.HandleFunc("/register", filmController.RegisterUser).Methods("POST")
	r.HandleFunc("/create", filmController.CreateUser).Methods("POST")
	r.HandleFunc("/user/:id", filmController.GetUser).Methods("GET")
	r.HandleFunc("/films", filmController.GetFilms).Methods("GET")
	r.HandleFunc("/film/:id", filmController.GetFilm).Methods("GET")

	filmSrv := http.Server{
		Handler: r,
		Addr:    ":5458",
	}

	err = filmSrv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Start error: %v", err)
	}
}
