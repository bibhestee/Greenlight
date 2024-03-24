package main

import (
  "fmt"
  "net/http"
  "strconv"
  "github.com/julienschmidt/httprouter"
)


func (app *application) showMovieHandler(res http.ResponseWriter, req *http.Request) {
  id, err := app.readIDParam(req)
  if err != nil {
    http.NotFound(res, req)
  }

  fmt.Fprintf(res, "show the details of movie %d\n", id)
}


func (app *application) createMovieHandler(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintln(res, "create a new movie")
}