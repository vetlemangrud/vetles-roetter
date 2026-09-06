package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := os.MkdirAll("data", 0755); err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", "file:data/carrotvault.sqlite")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	carrotRepository := CarrotRepository{DB: db}
	carrotRepository.initDatabase()

	carrotHandler := CarrotHandler{Repository: carrotRepository}

	http.HandleFunc("GET /", carrotHandler.Get)
	http.HandleFunc("POST /", carrotHandler.Post)

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
