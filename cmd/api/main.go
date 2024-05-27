package main

import (
  "context"
  "database/sql"
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"
  "time"
  "github.com/bibhestee/Greenlight/internal/data"
  "github.com/bibhestee/Greenlight/internal/jsonlog"
  _ "github.com/lib/pq"
)

// App version
const version = "1.0.0"

// Application configuration
type config struct {
  port int
  env string
  db  struct {
    dsn string
  }
}

// Application dependency injector
type application struct {
  config config
  logger *jsonlog.Logger
  models data.Models
}

func main() {
  var cfg config

  flag.IntVar(&cfg.port, "port", 4000, "API server port")
  flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
  flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

  flag.Parse()

  logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

  db, err := openDB(cfg)
  if err != nil {
    logger.PrintFatal(err, nil)
  }

  defer db.Close()

  logger.PrintInfo("database connection pool established", nil)

  app := &application{
    config: cfg,
    logger: logger,
    models: data.NewModels(db),
  }

  srv := &http.Server{
    Addr: fmt.Sprintf(":%d", cfg.port),
    Handler: app.routes(),
    ErrorLog: log.New(logger, "", 0),
    IdleTimeout: time.Minute,
    ReadTimeout: 10 * time.Second,
    WriteTimeout: 30 * time.Second,
  }

  logger.PrintInfo("starting server", map[string]string{
    "addr": srv.Addr,
    "env": cfg.env,
    })
  err = srv.ListenAndServe()
  logger.PrintFatal(err, nil)
}


func openDB(cfg config) (*sql.DB, error) {
  db, err := sql.Open("postgres", cfg.db.dsn)
  if err != nil {
    return nil, err
  }

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  err = db.PingContext(ctx)
  if err != nil {
    return nil, err
  }

  return db, nil
}
