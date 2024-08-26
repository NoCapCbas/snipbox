package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

// Holds application dependencies, to inject
type application struct {
	logger *slog.Logger
}

func main() {
	var err error
	// Config command line flags
	addy := flag.String("addy", ":4000", "HTTP network address")
	flag.Parse()

	// Structured Logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", "postgres://username:password@localhost/database_name?sslmode=disable")
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		logger.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}

	// Query data from PostgreSQL
	rows, err := db.Query("SELECT * FROM savvi_snips")
	if err != nil {
		logger.Error("Failed to query database", "error", err)
		os.Exit(1)
	}
	defer rows.Close()

	// Process the query results
	var data []map[string]interface{}
	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		err := rows.Scan(pointers...)
		if err != nil {
			logger.Error("Failed to scan row", "error", err)
			continue
		}

		row := make(map[string]interface{})
		for i, colName := range columns {
			row[colName] = values[i]
		}
		data = append(data, row)
	}

	fmt.Println(data)

	// initialization of app
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.Any("addy", *addy))

	err = http.ListenAndServe(*addy, app.routes())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
