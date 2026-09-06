package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	if err := os.MkdirAll("data", 0755); err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", "file:data/carrotvault.sqlite")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	carrotRepository := CarrotRepository{DB: db}
	if err := carrotRepository.initDatabase(); err != nil {
		panic(err)
	}

	carrotHandler := CarrotHandler{Repository: carrotRepository}

	http.HandleFunc("GET /api/", carrotHandler.Get)
	http.HandleFunc("POST /api/", requireAuth(carrotHandler.Post))

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
