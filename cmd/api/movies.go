package main

import (
  "fmt"
  "net/http"
  "strconv"
  "github.com/julienschmidt/httprouter"
)


func (app *application) showMovieHandler(res http.ResponseWriter, req *http.Request) {
  params := httprouter.ParamsFromContext(req.Context())

  id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
  if err != nil || id < 1 {
    http.NotFound(res, req)
    return
  }

  fmt.Fprintf(res, "show the details of movie %d\n", id)
}


func (app *application) createMovieHandler(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintln(res, "create a new movie")
}
