package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "io"
  "net/http"
  "net/url"
  "strconv"
  "strings"
  "github.com/bibhestee/Greenlight/internal/validator"
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

func (app *application) readString(qr url.Values, key string, defaultValue string) string {
  s := qr.Get(key)

  if s == "" {
    return defaultValue
  }
  return s
}

func (app *application) readInt(qr url.Values, key string, defaultValue int, v *validator.Validator) int {
  s := qr.Get(key)

  if s == "" {
    return defaultValue
  }

  i, err := strconv.Atoi(s)
  if err != nil {
    v.AddError(key, "must be an integer value")
    return defaultValue
  }

  return i
}

func (app *application) readCSV(qr url.Values, key string, defaultValue []string) []string {
  csv := qr.Get(key)

  if csv == "" {
    return defaultValue
  }

  return strings.Split(csv, ",")
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
  maxBytes := 1_048_576
  req.Body = http.MaxBytesReader(res, req.Body, int64(maxBytes))

  dec := json.NewDecoder(req.Body)
  dec.DisallowUnknownFields()

  err := dec.Decode(dst)
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
    case strings.HasPrefix(err.Error(), "json: unknwon field "):
      fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
      return fmt.Errorf("body contains unknown key %s", fieldName)
    case err.Error() == "http: request body too large":
      return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
    case errors.As(err, &invalidUnmarshalError):
      panic(err)
    default:
      return err
    }
  }

  err = dec.Decode(&struct{}{})
  if err != io.EOF {
    return errors.New("body must only contain a single JSON value")
  }

  return nil
}


func (app *application) background(fn func()) {
  go func() {
    defer func() {
      if err := recover(); err != nil {
        app.logger.PrintError(fmt.Errorf("%s", err), nil)
      }
    }()

    fn()
 }()
}
