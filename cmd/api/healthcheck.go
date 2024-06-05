package main

import (
  "net/http"
  "time"
)


func (app *application) healthcheckHandler(res http.ResponseWriter, req *http.Request) {
  env := envelope{
    "status": "available",
    "system_info": map[string]string{
      "environment": app.config.env,
      "version": version,
    },
  }

  time.Sleep(5*time.Second)

  err := app.writeJSON(res, http.StatusOK, env, nil)
  if err != nil {
    app.serverErrorResponse(res, req, err)
  }
}
