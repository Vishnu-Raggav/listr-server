package main

import (
	"encoding/json"
	"net/http"
)

// Read Tasks
func (app *application) readTasks(w http.ResponseWriter, r *http.Request) {
	// Read Database
	stmt := `SELECT * FROM tasks`

	rows, err := app.db.Query(stmt)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, "Database Error")
	}
	defer rows.Close()

	taskSlice := []Task{}
	for rows.Next() {
		var newTask Task
		if err = rows.Scan(&newTask.Id, &newTask.Name); err != nil {
			app.sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		taskSlice = append(taskSlice, newTask)
	}

	// Send Data as JSON
	jsonData := map[string]any{
		"tasks":   taskSlice,
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Create Task
func (app *application) insertTasks(w http.ResponseWriter, r *http.Request) {
	var newTask struct {
		Name string `json:"name"`
	}

	// Read JSON Data
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		app.sendError(w, http.StatusBadRequest, "Invalid JSON Format")
		return
	}

	// Validate Data
	if newTask.Name == "" {
		app.sendError(w, http.StatusBadRequest, "'name' Parameter is Expected")
		return
	}

	// Insert Into Database
	res, err := app.db.Exec(`INSERT INTO tasks (name) VALUES (?)`, newTask.Name)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, "Database Error")
		return
	}

	// Get Last Inserted Task
	id, err := res.LastInsertId()
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, "Database Error")
		return
	}

	var lastId int
	var lastName string

	row := app.db.QueryRow(`SELECT * FROM tasks WHERE id=?`, id)
	if err = row.Scan(&lastId, &lastName); err != nil {
		app.sendError(w, http.StatusInternalServerError, "Database Error")
		return
	}

	// Send JSON Data
	jsonData := map[string]any{
		"task": map[string]any{
			"id":   lastId,
			"name": lastName,
		},
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) deleteTasks(w http.ResponseWriter, r *http.Request) {
	// Parse Path Param
	taskId := r.PathValue("taskId")

	// Query Task
	queryStmt, err := app.db.Prepare("SELECT * FROM tasks WHERE id = ?")
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, "Database Error")
		return
	}
	defer queryStmt.Close()

	var id int
	var name string

	row := queryStmt.QueryRow(taskId)
	if err = row.Scan(&id, &name); err != nil {
		app.sendError(w, http.StatusInternalServerError, "Invalid Task ID")
		return
	}

	// Delete Task
	deleteStmt, err := app.db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, "Database Error")
		return
	}
	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(taskId)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, "Database Error")
		return
	}

	// Send JSON Data
	jsonData := map[string]any{
		"task": map[string]any{
			"id":   id,
			"name": name,
		},
		"success": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonData)
}
