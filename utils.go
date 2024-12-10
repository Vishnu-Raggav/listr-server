package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) initDB() {
	stmt := `
	CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
	)`

	_, err := app.db.Exec(stmt)
	if err != nil {
		app.log.SetPrefix("ERROR: ")
		app.log.Fatal(err)
	}
}

func (app *application) sendError(w http.ResponseWriter, statusCode int, errorMessage string) {
	jsonData := map[string]any{
		"error":   errorMessage,
		"success": false,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(jsonData)
}
