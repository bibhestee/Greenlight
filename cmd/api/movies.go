package main

import (
  "fmt"
  "errors"
  "net/http"
  "github.com/bibhestee/Greenlight/internal/data"
  "github.com/bibhestee/Greenlight/internal/validator"
)


func (app *application) showMovieHandler(res http.ResponseWriter, req *http.Request) {
  id, err := app.readIDParam(req)
  if err != nil {
    app.notFoundResponse(res, req)
    return
  }

  movie, err := app.models.Movies.Get(id)
  if err != nil {
    switch {
    case errors.Is(err, data.ErrRecordNotFound):
      app.notFoundResponse(res, req)
   default:
      app.serverErrorResponse(res, req, err)
    }
    return
  }

  err = app.writeJSON(res, http.StatusOK, envelope{"movie": movie}, nil)
  if err != nil {
    app.serverErrorResponse(res, req, err)
    }
}

func (app *application) listMoviesHandler(res http.ResponseWriter, req *http.Request) {
  var input struct {
    Title string
    Genres []string
    data.Filters
  }

  v := validator.New()

  queryString := req.URL.Query()

  input.Title = app.readString(queryString, "title", "")
  input.Genres = app.readCSV(queryString, "genres", []string{})
  input.Filters.Page = app.readInt(queryString, "page", 1, v)
  input.Filters.PageSize = app.readInt(queryString, "page_size", 20, v)
  input.Filters.Sort = app.readString(queryString, "sort", "id")
  input.Filters.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

  if data.ValidateFilters(v, input.Filters); !v.Valid() {
    app.failedValidationResponse(res, req, v.Errors)
    return
  }

  movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
  if err != nil {
    app.badRequestResponse(res, req, err)
  }

  err = app.writeJSON(res, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
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

func (app *application) updateMovieHandler(res http.ResponseWriter, req *http.Request) {
  id, err := app.readIDParam(req)
  if err != nil {
    app.notFoundResponse(res, req)
  }

  movie, err := app.models.Movies.Get(id)
  if err != nil {
    switch {
    case errors.Is(err, data.ErrRecordNotFound):
      app.notFoundResponse(res, req)
   default:
      app.serverErrorResponse(res, req, err)
    }
    return
  }

  var input struct {
    Title   *string    `json:"title"`
    Year    *int32     `json:"year"`
    Runtime *data.Runtime     `json:"runtime"`
    Genres  []string  `json:"genres"`
  }

  err = app.readJSON(res, req, &input)
  if err != nil {
    app.badRequestResponse(res, req, err)
    return
  }

  if input.Title != nil {
    movie.Title = *input.Title
  }

  if input.Year != nil {
    movie.Year = *input.Year
  }

  if input.Runtime != nil {
    movie.Runtime = *input.Runtime
  }

  if input.Genres != nil {
    movie.Genres = input.Genres
  }

  // Validate the updated movie record, sending the client a 422 Unprocessable Entity
  // response if any checks fail.
  v := validator.New()
  if data.ValidateMovie(v, movie); !v.Valid() {
    app.failedValidationResponse(res, req, v.Errors)
    return
  }
  // Pass the updated movie record to our new Update() method.
  err = app.models.Movies.Update(movie)
  if err != nil {
    switch {
      case errors.Is(err, data.ErrEditConflict):
        app.editConflictResponse(res, req)
      default:
        app.serverErrorResponse(res, req, err)
      }
    return
  }
   // Write the updated movie record in a JSON response.
  err = app.writeJSON(res, http.StatusOK, envelope{"movie": movie}, nil)
  if err != nil {
    app.serverErrorResponse(res, req, err)
  }
}


func (app *application) deleteMovieHandler(res http.ResponseWriter, req *http.Request) {
  id, err := app.readIDParam(req)
  if err != nil {
    app.notFoundResponse(res, req)
    return
  }

  err = app.models.Movies.Delete(id)
  if err != nil {
    switch {
    case errors.Is(err, data.ErrRecordNotFound):
      app.notFoundResponse(res, req)
    default:
      app.serverErrorResponse(res, req, err)
    }
    return
  }

  err = app.writeJSON(res, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
  if err != nil {
    app.serverErrorResponse(res, req, err)
  }
}
