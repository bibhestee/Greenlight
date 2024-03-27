package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "io"
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


func (app *application) readJSON(res http.ResponseWriter, req *http.Request, dst interface{}) error {
  err := json.NewDecoder(req.Body).Decode(dst)
  if err != nil {
    var syntaxError *json.SyntaxError
    var unmarshalTypeError *json.UnmarshalTypeError
    var invalidUnmarshalError *json.InvalidUnmarshalError

    switch {
    case errors.As(err, &syntaxError):
      return fmt.Errorf("body contain badly-formed JSON (at character %d)", syntaxError.Offset)
    case errors.Is(err, io.ErrUnexpectedEOF):
      return errors.New("body contain badly-formed JSON")
    case errors.As(err, &unmarshalTypeError):
      if unmarshalTypeError.Field != "" {
        return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
      }
      return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
    case errors.Is(err, io.EOF):
      return errors.New("body must not be empty")
    case errors.As(err, &invalidUnmarshalError):
      panic(err)
    default:
      return err
    }
  }
  return nil
}
