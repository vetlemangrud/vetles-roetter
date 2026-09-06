package main

import (
	"database/sql"
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

//go:embed templates/carrots.html
var carrotsHTML string

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

	tmpl := template.Must(template.New("carrotsHome").Parse(carrotsHTML))
	carrotHandler := CarrotHandler{Repository: carrotRepository, HomeTemplate: tmpl}

	http.HandleFunc("GET /api/", carrotHandler.APIGet)
	http.HandleFunc("POST /api/", requireAuth(carrotHandler.APIPost))

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
