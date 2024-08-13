package main

import (
  "log/slog"
  "net/http"
  "flag"
  "os"
)

func main() {
  // Config command line flags
  addy := flag.String("addy", ":4000", "HTTP network address")
  flag.Parse()

  // Structured Logging
  logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    AddSource: true,
  }))
  mux := http.NewServeMux()

  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("GET /{$}", home) // restrict route to exact match
  mux.HandleFunc("GET /snippet/view/{id}", snippetView)
  mux.HandleFunc("GET /snippet/create", snippetCreate)
  mux.HandleFunc("POST /snippet/create", snippetCreatePost)

  logger.Info("starting server", slog.Any("addy", *addy))

  err := http.ListenAndServe(*addy, mux)
  logger.Error(err.Error())
  os.Exit(1)
}

