package main

import (
	"encoding/json"
	"net/http"
)

// Read Tasks
func (app *application) readDatabase(w http.ResponseWriter, r *http.Request) {
	// Read Database
	stmt := `SELECT * FROM tasks`

	rows, err := app.db.Query(stmt)
	if err != nil {
		jsonData := map[string]any{
			"error": "Database Error",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonData)
	}
	defer rows.Close()

	taskSlice := []Task{}

	for rows.Next() {
		var newTask Task
		if err = rows.Scan(&newTask.Id, &newTask.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		taskSlice = append(taskSlice, newTask)
	}

	// Send Data as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(taskSlice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create Task
func (app *application) insertDatabase(w http.ResponseWriter, r *http.Request) {
	var newTask struct {
		Name string `json:"name"`
	}

	// Read JSON Data
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		jsonData := map[string]any{
			"error":   "Invalid JSON Format",
			"created": false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(jsonData)
		return
	}

	// Validate Data
	if newTask.Name == "" {
		jsonData := map[string]any{
			"error":   `"name" Parameter Is Expected`,
			"created": false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(jsonData)
		return
	}

	// Insert Into Database
	res, err := app.db.Exec(`INSERT INTO tasks (name) VALUES (?)`, newTask.Name)
	if err != nil {
		jsonData := map[string]any{
			"error":   "Database Error",
			"created": false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonData)
		return
	}

	// Get Last Inserted Task
	id, err := res.LastInsertId()
	if err != nil {
		jsonData := map[string]any{
			"error":   "Database Error",
			"created": false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonData)
		return
	}

	var lastId int
	var lastName string

	row := app.db.QueryRow(`SELECT * FROM tasks WHERE id=?`, id)
	if err = row.Scan(&lastId, &lastName); err != nil {
		jsonData := map[string]any{
			"error":   "Database Error",
			"created": false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonData)
		return
	}

	// Send JSON Data
	jsonData := map[string]any{
		"task": map[string]any{
			"id":   lastId,
			"name": lastName,
		},
		"created": true,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
