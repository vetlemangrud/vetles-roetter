package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func createCarrotTable(db *sql.DB) {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS carrots (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"eaten_at" DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
	_, err := db.Exec(createTableSQL);
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "file:carrotvault.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createCarrotTable(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "LetsGo")
	})

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
