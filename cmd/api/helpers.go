package main

import (
  "encoding/json"
  "errors"
  "net/http"
  "strconv"

  "github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}


func (app *application) readIDParam(req *http.Request) (int64, error) {
  params := httprouter.ParamsFromContext(req.Context())

  id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
  if err != nil || id < 1 {
    return 0, errors.New("invalid id parameter")
  }
  return id, nil
}


func (app *application) writeJSON(res http.ResponseWriter, status int, data envelope, headers http.Header) error {
  js, err := json.MarshalIndent(data, "", "\t")
  if err != nil {
    return err
  }

  js = append(js, '\n')

  for key, value := range(headers) {
    res.Header()[key] = value
  }

  res.Header().Set("Content-Type", "application/json")
  res.WriteHeader(status)
  res.Write(js)

  return nil

}
