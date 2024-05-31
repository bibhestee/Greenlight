package main

import (
  "fmt"
  "net/http"
)


func (app *application) recoverPanic(next http.Handler) http.Handler {
  return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
    defer func() {
      if err := recover(); err != nil {
        res.Header().Set("Connection", "close")
        app.serverErrorResponse(res, req, fmt.Errorf("%s", err))
      }
    }()
    next.ServeHTTP(res, req)
  })
}
