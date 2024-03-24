package main

import (
  "fmt"
  "net/http"
)


func (app *application) healthcheckHandler(res http.ResponseWriter, req *http.Request) {
  fmt.Fprintln(res, "status: available")
  fmt.Fprintf(res, "environment: %s\n", app.config.env)
  fmt.Fprintf(res, "version: %s\n", version)
}
