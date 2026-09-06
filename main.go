package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := os.Mkdir("data", 0755)
	db, err := sql.Open("sqlite3", "file:data/carrotvault.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	carrotRepository := CarrotRepository{DB: db}
	carrotRepository.initDatabase()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "LetsGo")
	})

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
