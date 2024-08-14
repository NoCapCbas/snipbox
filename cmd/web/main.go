package main

import (
  "log/slog"
  "net/http"
  "flag"
  "os"
)

// Holds application dependencies, to inject
type application struct {
  logger *slog.Logger
}

func main() {
  // Config command line flags
  addy := flag.String("addy", ":4000", "HTTP network address")
  flag.Parse()

  // Structured Logging
  logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    AddSource: true,
  }))
  
  // initialization of app
  app := &application{
    logger: logger,
  }

  mux := http.NewServeMux()

  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("GET /{$}", app.home) // restrict route to exact match
  mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
  mux.HandleFunc("GET /snippet/create", app.snippetCreate)
  mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

  logger.Info("starting server", slog.Any("addy", *addy))

  err := http.ListenAndServe(*addy, mux)
  logger.Error(err.Error())
  os.Exit(1)
}

