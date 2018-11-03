package apis

import (
	"database/sql"
	"encoding/json"
	"github.com/eefret/gomdb"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"imdb_rest_api/config"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Movie struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	ReleasedYear int            `json:"released_year"`
	Rating       json.Number    `json:"rating"`
	Genres       pq.StringArray `json:"genres"`
}

type Year struct {
	From int `json:"from"`
	To   int `json:"to"`
}

var (
	DBConn = config.Connect().DBConn
)

type Rating struct {
	RatingPoint json.Number `json:"rating"`
	FilterBy    string      `json:"filterby"`
}

func GetMoviByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	rows, err := DBConn.Query("SELECT id, title, released_year, rating, genres FROM movies where id= $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var m Movie
	for rows.Next() {
		if err := rows.Scan(&m.ID, &m.Title, &m.ReleasedYear, &m.Rating, (*pq.StringArray)(&m.Genres)); err != nil {
			log.Fatal(err)
		}
	}
	json.NewEncoder(w).Encode(m)
	return
}

func GetMoviesByYear(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	year := params["year"]
	rows, err := DBConn.Query("SELECT id, title, released_year, rating, genres FROM movies where released_year= $1", year)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var m Movie
	var movies []Movie
	for rows.Next() {
		if err := rows.Scan(&m.ID, &m.Title, &m.ReleasedYear, &m.Rating, (*pq.StringArray)(&m.Genres)); err != nil {
			log.Fatal(err)
		}
		movies = append(movies, Movie{ID: m.ID, Title: m.Title, ReleasedYear: m.ReleasedYear, Rating: m.Rating, Genres: m.Genres})
	}
	if movies == nil {
		json.NewEncoder(w).Encode("No movies found of following Year")
	} else {
		json.NewEncoder(w).Encode(movies)
		return
	}
}

func GetMoviesByYearRange(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var y Year
	err := decoder.Decode(&y)
	if err != nil {
		panic(err)
	}
	rows, err := DBConn.Query("select id, title, released_year, rating, genres from movies where released_year between $1 and $2;", y.From, y.To)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var m Movie
	var movies []Movie
	for rows.Next() {
		if err := rows.Scan(&m.ID, &m.Title, &m.ReleasedYear, &m.Rating, (*pq.StringArray)(&m.Genres)); err != nil {
			log.Fatal(err)
		}
		movies = append(movies, Movie{ID: m.ID, Title: m.Title, ReleasedYear: m.ReleasedYear, Rating: m.Rating, Genres: m.Genres})
	}
	if movies == nil {
		json.NewEncoder(w).Encode("No movies found in following year range")
	} else {
		json.NewEncoder(w).Encode(movies)
		return
	}
}

func GetMoviesByRating(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var count Rating
	err := decoder.Decode(&count)
	if err != nil {
		panic(err)
	}
	var rows *sql.Rows
	switch count.FilterBy {
	case "":
		rows, err = DBConn.Query("select id, title, released_year, rating, genres from movies where rating= $1", count.RatingPoint)
		if err != nil {
			log.Fatal(err)
		}
	case "higher":
		rows, err = DBConn.Query("select id, title, released_year, rating, genres from movies where rating > $1", count.RatingPoint)
		if err != nil {
			log.Fatal(err)
		}
	case "lower":
		rows, err = DBConn.Query("select id, title, released_year, rating, genres from movies where rating < $1", count.RatingPoint)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer rows.Close()
	var m Movie
	var movies []Movie
	for rows.Next() {
		if err := rows.Scan(&m.ID, &m.Title, &m.ReleasedYear, &m.Rating, (*pq.StringArray)(&m.Genres)); err != nil {
			log.Fatal(err)
		}
		movies = append(movies, Movie{ID: m.ID, Title: m.Title, ReleasedYear: m.ReleasedYear, Rating: m.Rating, Genres: m.Genres})
	}
	if movies == nil {
		json.NewEncoder(w).Encode("No movies found of following ratings")
	} else {
		json.NewEncoder(w).Encode(movies)
		return
	}
}

func GetMoviesByGenres(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var genres Movie
	err := decoder.Decode(&genres)
	if err != nil {
		panic(err)
	}
	var movies []Movie
	for _, i := range genres.Genres {
		rows, err := DBConn.Query("SELECT id, title, released_year, rating, genres from (SELECT id, title, released_year, rating, genres, generate_subscripts(genres, 1) AS s FROM movies) WHERE genres[s] = $1", i)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var m Movie
		for rows.Next() {
			if err := rows.Scan(&m.ID, &m.Title, &m.ReleasedYear, &m.Rating, (*pq.StringArray)(&m.Genres)); err != nil {
				log.Fatal(err)
			}
			movies = append(movies, Movie{ID: m.ID, Title: m.Title, ReleasedYear: m.ReleasedYear, Rating: m.Rating, Genres: m.Genres})
		}
	}
	if movies == nil {
		json.NewEncoder(w).Encode("No movies found of following genre(s)")
	} else {
		json.NewEncoder(w).Encode(movies)
		return
	}
}

func GetMovieByTitle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var title Movie
	err := decoder.Decode(&title)
	if err != nil {
		panic(err)
	}
	rows := DBConn.QueryRow("SELECT id, title, released_year, rating, genres FROM movies where title= $1", title.Title)
	if err != nil {
		log.Fatal(err)
	}
	var m Movie
	switch err = rows.Scan(&m.ID, &m.Title, &m.ReleasedYear, &m.Rating, (*pq.StringArray)(&m.Genres)); err {
	case sql.ErrNoRows:
		api := gomdb.Init("9d22088a")
		query := &gomdb.QueryData{Title: title.Title, SearchType: gomdb.MovieSearch}
		res, err := api.MovieByTitle(query)
		if err != nil {
			log.Println(err)
		}
		json.NewEncoder(w).Encode(res)
		Genres := strings.Split(res.Genre, ", ")
		var GenreArray []string
		for i := range Genres {
			GenreArray = append(GenreArray, Genres[i])
		}
		if res.ImdbRating == "N/A" {
			res.ImdbRating = "0"
		}
		imdbrating, err := strconv.ParseFloat(res.ImdbRating, 64)
		_, err = DBConn.Exec("INSERT INTO movies (title, released_year, rating, genres) VALUES ($1, $2, $3, $4)", res.Title, res.Year, imdbrating, pq.Array(GenreArray))
		if err != nil {
			log.Fatal(err)
		}
		return
	case nil:
		json.NewEncoder(w).Encode(m)
		return
	default:
		panic(err)
	}
}

func UpdateRatingOfMovie(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var movie Movie
	err := decoder.Decode(&movie)
	if err != nil {
		panic(err)
	}
	rows, err := DBConn.Exec("update movies set rating= $1 where id= $2", movie.Rating, movie.ID)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := rows.RowsAffected()
	if rowsAffected == 0 {
		json.NewEncoder(w).Encode("Rating update failed for requested movie!!")
		return
	}
	if rowsAffected == 1 {
		json.NewEncoder(w).Encode("Rating successfully updated for requested movie!!")
		return
	}
}

func UpdateGenresOfMovie(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var movie Movie
	err := decoder.Decode(&movie)
	if err != nil {
		panic(err)
	}
	var rows sql.Result
	for _, gener := range movie.Genres {
		rows, err = DBConn.Exec("UPDATE movies SET genres = array_append(genres, $1) WHERE id= $2", gener, movie.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
	rowsAffected, err := rows.RowsAffected()
	if rowsAffected == 0 {
		json.NewEncoder(w).Encode("Genres update failed for requested movie!!")
		return
	}
	if rowsAffected == 1 {
		json.NewEncoder(w).Encode("Genres successfully updated for requested movie!!")
		return
	}
}
