package main

import (
	"database/sql"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

//go:embed static
var staticFS embed.FS

//go:embed templates/home.html
var homeHTML string

//go:embed templates/vetle.html
var vetleHTML string

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

	homeTemplate := template.Must(template.New("home").Parse(homeHTML))
	vetleTemplate := template.Must(template.New("vetle").Parse(vetleHTML))
	carrotHandler := CarrotHandler{Repository: carrotRepository, HomeTemplate: homeTemplate, VetleTemplate: vetleTemplate}

	http.HandleFunc("GET /{$}", carrotHandler.HomeGet)
	http.HandleFunc("GET /vetle", carrotHandler.VetleGet)
	http.HandleFunc("POST /vetle", carrotHandler.VetlePost)
	http.HandleFunc("GET /api/", carrotHandler.APIGet)
	http.HandleFunc("POST /api/", requireAuth(carrotHandler.APIPost))
	http.Handle("GET /static/", http.FileServerFS(staticFS))

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
