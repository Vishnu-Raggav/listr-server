package main

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
