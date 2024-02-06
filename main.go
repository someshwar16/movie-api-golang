package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var movies []Movie

func main() {
	fmt.Println("Starting the server...")
	router := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Title: "Inception", Director: &Director{FirstName: "Christopher", LastName: "Nolan"}},
		Movie{ID: "2", Title: "The Dark Knight", Director: &Director{FirstName: "Christopher", LastName: "Nolan"}})

	router.HandleFunc("/", welcomeUser).Methods("GET")
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST") // Create a new movie
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func welcomeUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Welcome to the movie API"}`))
}

// Get All Movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Get Single Movie By Id
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for _, movie := range movies {
		if movie.ID == id {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

}

// Crete a Movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newMovie Movie

	err := json.NewDecoder(r.Body).Decode(&newMovie)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request payload")
		return
	}
	newMovie.ID = strconv.Itoa(rand.Intn(100))
	movies = append(movies, newMovie)
	json.NewEncoder(w).Encode("Movie created successfully" + " " + newMovie.Title)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	for index, movie := range movies {
		if movie.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			var updatedMovie Movie
			err := json.NewDecoder(r.Body).Decode(&updatedMovie)
			checkError(err)
			updatedMovie.ID = id
			movies = append(movies, updatedMovie)
			json.NewEncoder(w).Encode("Movie updated successfully" + " " + id)
			return
		}
	}
	json.NewEncoder(w).Encode("Movie not found")
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Invalid movie ID"}`))
		return
	}
	for index, movie := range movies {
		if movie.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			json.NewEncoder(w).Encode("Movie deleted successfully" + " " + id)
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
