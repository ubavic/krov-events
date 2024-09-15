package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"decentrala.org/events/internal/controller"
	"decentrala.org/events/internal/model"
	"decentrala.org/events/internal/view"
)

//go:embed static/*
var staticFS embed.FS

//go:embed template/*
var templateFS embed.FS

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	test := os.Getenv("TEST_INSTANCE")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")
	jwtKey := os.Getenv("JWT_SECRET")

	if dbUser == "" || dbPassword == "" || dbName == "" || jwtKey == "" {
		panic("invalid env parameters")
	}

	connStr := "user=" + dbUser +
		" dbname=" + dbName +
		" password=" + dbPassword +
		" sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	templateSub, err := fs.Sub(templateFS, "template")
	if err != nil {
		panic(err)
	}

	var staticSub fs.FS
	if test == "1" {
		staticSub = os.DirFS("static")
	} else {
		staticSub, err = fs.Sub(staticFS, "static")
		if err != nil {
			panic(err)
		}
	}

	view := view.NewView(templateSub)
	model := model.NewModel(db)
	controller := controller.NewController(model, view, staticSub, jwtKey)

	mux := controller.Mux()

	if port == "" {
		port = "8080"
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("http://localhost" + server.Addr + "/")

	server.ListenAndServe()
}
