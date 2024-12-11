package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

// Application Struct
type application struct {
	db  *sql.DB
	log *log.Logger
}

// Task Struct
type Task struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Create Logger
	logger := log.New(os.Stdout, "INFO: ", 0)

	// Load .Env
	err := godotenv.Load()
	if err != nil {
		logger.SetPrefix("ERROR: ")
		logger.Fatal(err)
	}

	// Connect Database
	db, err := sql.Open("sqlite", os.Getenv("DB_URL"))
	if err != nil {
		logger.SetPrefix("ERROR: ")
		logger.Fatal(err)
	}

	app := application{db: db, log: logger}
	app.initDB()

	// HTTP Router
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", app.readTasks)
	router.HandleFunc("POST /{$}", app.insertTasks)
	router.HandleFunc("DELETE /{taskId}", app.deleteTasks)

	server := http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: loggingMiddleware(corsMiddleware(router)),
	}
	server.ListenAndServe()
}
