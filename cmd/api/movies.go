package main

import (
  "fmt"
  "net/http"
  "time"
  "github.com/bibhestee/Greenlight/internal/data"
  "github.com/bibhestee/Greenlight/internal/validator"
)


func (app *application) showMovieHandler(res http.ResponseWriter, req *http.Request) {
  id, err := app.readIDParam(req)
  if err != nil {
    app.notFoundResponse(res, req)
    return
  }

  movie := data.Movie{
    ID: id,
    CreatedAt: time.Now(),
    Title: "Casablanca",
    Runtime: 102,
    Genres: []string{"drama", "romance", "war"},
    Version: 1,
  }

  err = app.writeJSON(res, http.StatusOK, envelope{"movie": movie}, nil)
  if err != nil {
    app.serverErrorResponse(res, req, err)
    }
}


func (app *application) createMovieHandler(res http.ResponseWriter, req *http.Request) {
  var input struct {
    Title   string    `json:"title"`
    Year    int32     `json:"year"`
    Runtime data.Runtime     `json:"runtime"`
    Genres  []string  `json:"genres"`
  }

  err := app.readJSON(res, req, &input)
  if err != nil {
    app.badRequestResponse(res, req, err)
    return
  }

  movie := &data.Movie{
    Title: input.Title,
    Year: input.Year,
    Runtime: input.Runtime,
    Genres: input.Genres,
  }

  v := validator.New()

  if data.ValidateMovie(v, movie); !v.Valid() {
    app.failedValidationResponse(res, req, v.Errors)
    return
  }

  err = app.models.Movies.Insert(movie)
  if err != nil {
    app.serverErrorResponse(res, req, err)
    return
  }

  headers := make(http.Header)
  headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

  err = app.writeJSON(res, http.StatusCreated, envelope{"movie": movie}, headers)
  if err != nil {
    app.serverErrorResponse(res, req, err)
  }
}

