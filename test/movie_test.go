package test

import (
	"net/http"
	"net/http/httptest"
	"rm_movie_backend/controllers"
	"rm_movie_backend/models"
	"strings"
	"testing"
)

var dummyIds []int

// Test for /movie/search
func TestSearchMovie(t *testing.T) {
	setupSearchMovie()
	defer teardownSearchMovie()

	expectedMovie1 := `"Title":"Dummy Movie 1","Description":"This is a dummy description for Dummy Movie 1. Insert special characters for search fsa35jifjas21W%$!#%3a","FileName":"/assets/movies/dummy1.jpg","OriginalLink":"https://www.google.com/img/asfoij2515.jpg"`
	expectedMovie2 := `"Title":"Dummy Movie 2","Description":"This is a dummy description for Dummy Movie 2. Insert special characters for search fsa3at3532tgs3$#","FileName":"/assets/movies/dummy2.jpg","OriginalLink":"https://www.google.com/img/gsa463sdyw4.jpg"`

	t.Run("keyword_title", testSearchMovie("keyword=Dummy Movie 1", []string{expectedMovie1}))
	t.Run("keyword_description", testSearchMovie("keyword=fsa3at3532tgs3$#", []string{expectedMovie2}))
	t.Run("keyword_description_2", testSearchMovie("keyword=dummy movie", []string{expectedMovie1, expectedMovie2}))
	t.Run("keyword_filename", testSearchMovie("keyword=dummy1", []string{expectedMovie1}))
	t.Run("keyword_originallink", testSearchMovie("keyword=gsa463sdyw4", []string{expectedMovie2}))
}

func testSearchMovie(query string, expecteds []string) func(*testing.T) {
	return func(t *testing.T) {
		req, err := http.NewRequest("GET", "/movie/search?"+query, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controllers.SearchMovie)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check the response body is what we expect.
		result := rr.Body.String()
		for _, expected := range expecteds {
			if !strings.Contains(result, expected) {
				t.Errorf("Handler returned unexpected body: got %v want string containing %v", result, expected)
			}
		}
	}
}

// The setup method for search movie
func setupSearchMovie() {
	// Insert two dummy records for the search
	// First dummy movie
	movie := models.Movie{
		Title:        "Dummy Movie 1",
		Description:  "This is a dummy description for Dummy Movie 1. Insert special characters for search fsa35jifjas21W%$!#%3a",
		FileName:     "/assets/movies/dummy1.jpg",
		OriginalLink: "https://www.google.com/img/asfoij2515.jpg",
	}

	db := models.GetDB()
	defer db.Close()

	db.Create(&movie)
	dummyIds = append(dummyIds, movie.ID)

	// First dummy movie
	movie = models.Movie{
		Title:        "Dummy Movie 2",
		Description:  "This is a dummy description for Dummy Movie 2. Insert special characters for search fsa3at3532tgs3$#",
		FileName:     "/assets/movies/dummy2.jpg",
		OriginalLink: "https://www.google.com/img/gsa463sdyw4.jpg",
	}

	db.Create(&movie)
	dummyIds = append(dummyIds, movie.ID)
}

// The teardown method for search movie
func teardownSearchMovie() {
	db := models.GetDB()
	defer db.Close()

	db.Where(dummyIds).Delete(&models.Movie{})
}

// Test for /tmdb/discover
func TestDiscoverMovie(t *testing.T) {
	req, err := http.NewRequest("GET", "/tmdb/discover", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.DiscoverMovie)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	result := rr.Body.String()
	expected := `"message":"20 movie(s) have been created / updated.","success":true`
	if !strings.Contains(result, expected) {
		t.Errorf("Handler returned unexpected body: got %v want string containing %v", result, expected)
	}
}
