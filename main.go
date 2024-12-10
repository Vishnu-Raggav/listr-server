package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

type application struct {
	db  *sql.DB
	log *log.Logger
}

type Task struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// create logger
	logger := log.New(os.Stdout, "INFO: ", 0)

	// connect database
	db, err := sql.Open("sqlite", "tasks.db")
	if err != nil {
		logger.SetPrefix("ERROR: ")
		logger.Fatal(err)
	}

	app := application{db: db, log: logger}
	app.initDB()

	// create router
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", app.readDatabase)
	router.HandleFunc("POST /{$}", app.insertDatabase)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}
