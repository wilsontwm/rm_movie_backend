package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"rm_movie_backend/utils"
	"strings"
	"time"
)

type Movie struct {
	ID           int    `gorm:"primary_key;"`
	Title        string `sql:"type:longtext"`
	Description  string `sql:"type:longtext"`
	FileName     string
	OriginalLink string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SearchMovieFilter struct {
	ID      int
	Keyword string
	Limit   int
	Page    int
}

// Create / update the movie
func CreateUpdateMovie(input Movie) (*Movie, error) {
	movie := input

	// Download the image to local filesystem
	var folder = "assets/movies"
	var fileName = strings.TrimPrefix(movie.OriginalLink, "/")
	var fullPath = folder + "/" + fileName
	movie.OriginalLink = "https://image.tmdb.org/t/p/w500" + movie.OriginalLink

	err := utils.DownloadFile(movie.OriginalLink, folder, fileName)

	if err != nil {
		return nil, fmt.Errorf("Error downloading image.")
	}

	movie.FileName = fullPath

	db := GetDB()
	defer db.Close()

	// Create the movie
	db.Where(Movie{ID: movie.ID}).Assign(movie).FirstOrCreate(&movie)

	if movie.ID <= 0 {
		return nil, fmt.Errorf("Movie is not created.")
	}

	return &movie, nil
}

// Create / update multiple movies
func CreateUpdateMultipleMovies(movies []Movie) ([]Movie, error) {
	// Create the jobs and worker pools
	var numOfJobs = len(movies)
	var numOfWorkers = 4
	var outputs []Movie
	if numOfJobs > 0 {
		movieJobs := make(chan Movie, numOfJobs) // Accept the movie input
		results := make(chan *Movie, numOfJobs)  // Return the movie that is successfully stored

		for w := 1; w <= numOfWorkers; w++ {
			go movieCreateUpdateWorker(movieJobs, results)
		}

		for _, m := range movies {
			movieJobs <- m
		}

		close(movieJobs)

		for a := 1; a <= numOfJobs; a++ {
			output := <-results
			if output != nil {
				outputs = append(outputs, *output)
			}
		}
	}

	return outputs, nil
}

// Search the movie by keywords / ID
func SearchMovie(filter SearchMovieFilter) ([]Movie, error) {
	// Get movies based on filter
	var movies []Movie

	db := GetDB()
	defer db.Close()

	var offset = filter.Limit * (filter.Page - 1)
	// Get the list of movies
	db.Table("movies").
		Scopes(filterByID(filter.ID), filterByKeyword(filter.Keyword)).
		Limit(filter.Limit).
		Offset(offset).
		Find(&movies)

	return movies, nil
}

// Filter by ID
func filterByID(id int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id > 0 {
			return db.Where("id = ?", id)
		}

		return db
	}
}

// Filter by keyword
func filterByKeyword(keyword string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(keyword) > 0 {
			keyword = "%" + keyword + "%"
			return db.Where("title LIKE ? or description LIKE ? or original_link LIKE ? or file_name LIKE ?", keyword, keyword, keyword, keyword)
		}

		return db
	}
}

// Worker to trigger the creation / update of the movie
func movieCreateUpdateWorker(movies <-chan Movie, results chan<- *Movie) {
	for movie := range movies {
		m, _ := CreateUpdateMovie(movie)
		results <- m
	}
}
