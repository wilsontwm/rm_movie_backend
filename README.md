## About

This is a movie API which does two main things:

1. To extract photo information, retrieve photos from [The Movie DB](https://www.themoviedb.org/documentation/api) and save photos in local filesystem
2. To accept input from users to query photo name, description, filename or original link from the photo information stored in MySQL database

## Project setup

1. Clone the current repository
2. Change the environment variables in root directory and /test directory
3. Create the database in MySQL server
4. Run the following command to initiate the Go web server

```
go run main.go
```

## API Examples

GET /tmdb/discover :

This will send a GET request to https://api.themoviedb.org/3/discover/movie to get the list of movies. Upon retrieval of the movies, the information of the photo of each movie will be retrieved and stored in local file system and database

GET /movie/search :

This will search through the movie photos stored in the database. It also allows searching by keyword in photo name,description, filename or original link

Examples:
```
curl http://localhost:3000/movie/search?keyword=abc&limit=5
curl http://localhost:3000/movie/search?page=2&limit=5
```

## Testing

The unit test involves the testing of the below:
1. /tmdb/discover - To check if the movie photos information have been retrieved and stored in database correctly
2. /movie/search - To check if the searching of the movie photos in database is performed correctly

The unit test can be executed via:

```
cd test
go test
```
