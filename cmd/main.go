package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"imdb_rest_api/apis"
	"imdb_rest_api/config"
	"log"
	"net/http"
)

func main() {
	db := config.Connect()
	if _, err := db.DBConn.Exec("CREATE TABLE IF NOT EXISTS movies (id serial not null primary key, title string not null, released_year int, rating decimal, genres string[] null)"); err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/movie/{id}", apis.GetMoviByID).Methods("GET")
	router.HandleFunc("/movies-by-year/{year}", apis.GetMoviesByYear).Methods("GET")
	router.HandleFunc("/movies-by-year-range", apis.GetMoviesByYearRange).Methods("POST")
	router.HandleFunc("/movies-by-rating", apis.GetMoviesByRating).Methods("POST")
	router.HandleFunc("/movies-by-genres", apis.GetMoviesByGenres).Methods("POST")
	router.HandleFunc("/movie-by-title", apis.GetMovieByTitle).Methods("POST")
	router.HandleFunc("/update-ratings", apis.UpdateRatingOfMovie).Methods("POST")
	router.HandleFunc("/update-genres", apis.UpdateGenresOfMovie).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
