package main

import (
  "fmt"
  "net/http"
)


func (app *application) logError(req *http.Request, err error) {
  app.logger.PrintError(err, map[string]string{
    "request_method": req.Method,
    "request_url": req.URL.String(),
  })
}

func (app *application) rateLimitExceededResponse(res http.ResponseWriter, req *http.Request) {
  message := "rate limit exceeded"
  app.errorResponse(res, req, http.StatusTooManyRequests, message)
}

func (app *application) errorResponse(res http.ResponseWriter, req *http.Request, status int, message interface{}) {
  env := envelope{"error": message}

  err := app.writeJSON(res, status, env, nil)
  if err != nil {
    app.logError(req, err)
    res.WriteHeader(500)
  }
}


func (app *application) serverErrorResponse(res http.ResponseWriter, req *http.Request, err error) {
  app.logError(req, err)
  message := "the server encountered a problem and could not process your request"

  app.errorResponse(res, req, http.StatusInternalServerError, message)
}


func (app *application) notFoundResponse(res http.ResponseWriter, req *http.Request) {
  message := "the requested resource could not be found"
  app.errorResponse(res, req, http.StatusNotFound, message)
}


func (app *application) methodNotAllowedResponse(res http.ResponseWriter, req *http.Request) {
  message := fmt.Sprintf("the %s method is not supported for this resource", req.Method)
  app.errorResponse(res, req, http.StatusMethodNotAllowed, message)
}


func (app *application) badRequestResponse(res http.ResponseWriter, req *http.Request, err error) {
  app.errorResponse(res, req, http.StatusBadRequest, err.Error())
}


func (app *application) failedValidationResponse(res http.ResponseWriter, req *http.Request, errors map[string]string) {
  app.errorResponse(res, req, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(res http.ResponseWriter, req *http.Request) {
  message := "unable to update the record due to an edit conflict, please try again"
  app.errorResponse(res, req, http.StatusConflict, message)
}
