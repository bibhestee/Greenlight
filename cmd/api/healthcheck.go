package main

import (
  "net/http"
)


func (app *application) healthcheckHandler(res http.ResponseWriter, req *http.Request) {
  env := envelope{
    "status": "available",
    "system_info": map[string]string{
      "environment": app.config.env,
      "version": version,
    },
  }

  err := app.writeJSON(res, http.StatusOK, env, nil)
  if err != nil {
    app.logger.Println(err)
    http.Error(res, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
  }
}
