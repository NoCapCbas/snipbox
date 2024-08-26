package main

import (
	"flag"
	"log/slog"
	"net/http"
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

	logger.Info("starting server", slog.Any("addy", *addy))

	err := http.ListenAndServe(*addy, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
