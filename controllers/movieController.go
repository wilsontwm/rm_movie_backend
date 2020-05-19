package controllers

import (
	"encoding/json"
	"rm_movie_backend/models"
	//"fmt"
	"net/http"
	"os"
	"rm_movie_backend/utils"
	"strconv"
)

// GET: Get list of movies from TMDB via discovery
var DiscoverMovie = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	// Get the movies from TMDB API
	url := os.Getenv("TMDB_API_URL") + "discover/movie?api_key=" + os.Getenv("TMDB_KEY")
	res, err := http.Get(url)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	defer res.Body.Close()

	type Movie struct {
		ID           int    `json:"id"`
		Title        string `json:"title"`
		Description  string `json:"overview"`
		OriginalLink string `json:"poster_path"`
	}

	type Result struct {
		Page   int     `json:"page"`
		Movies []Movie `json:"results"`
	}
	result := Result{}
	json.NewDecoder(res.Body).Decode(&result)

	if len(result.Movies) == 0 {
		utils.Success(w, http.StatusOK, resp, result.Movies, "There is no movies to be discovered.")
		return
	}

	// Store the movies into database and download the picture
	var movies []models.Movie
	for _, movie := range result.Movies {
		m := &models.Movie{
			ID:           movie.ID,
			Title:        movie.Title,
			Description:  movie.Description,
			OriginalLink: movie.OriginalLink,
		}
		movies = append(movies, *m)
	}

	data, err := models.CreateUpdateMultipleMovies(movies)

	if err != nil {
		utils.Fail(w, http.StatusBadRequest, resp, err.Error())
		return
	} else if len(data) == 0 {
		utils.Fail(w, http.StatusBadRequest, resp, "No movie(s) have been created nor updated.")
		return
	}

	utils.Success(w, http.StatusOK, resp, data, strconv.Itoa(len(data))+" movie(s) have been created / updated.")
}

// GET: Search movies in database by keyword
var SearchMovie = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	id := 0
	if idParam := utils.GetParam(r, "id"); idParam != "" {
		id, _ = strconv.Atoi(idParam)
	}

	keyword := ""
	if keywordParam := utils.GetParam(r, "keyword"); keywordParam != "" {
		keyword = keywordParam
	}

	limit := 10
	if limitParam := utils.GetParam(r, "limit"); limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
	}

	page := 1
	if pageParam := utils.GetParam(r, "page"); pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	// Build the filter for the movie search result
	filter := models.SearchMovieFilter{
		ID:      id,
		Keyword: keyword,
		Limit:   limit,
		Page:    page,
	}

	// Get the movies list
	movies, err := models.SearchMovie(filter)
	if err != nil {
		utils.Fail(w, http.StatusBadRequest, resp, err.Error())
		return
	}

	response := map[string]interface{}{
		"result": movies,
		"page":   page,
		"count":  len(movies),
	}

	utils.Success(w, http.StatusOK, resp, response, "")
}
