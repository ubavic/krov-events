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
	"decentrala.org/events/internal/migration"
	"decentrala.org/events/internal/model"
	"decentrala.org/events/internal/view"
)

//go:embed static/*
var staticFS embed.FS

//go:embed template/*
var templateFS embed.FS

//go:embed migrations/*
var migrationsFS embed.FS

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	test := os.Getenv("TEST_INSTANCE")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")
	jwtKey := os.Getenv("JWT_SECRET")

	if dbUser == "" || dbPassword == "" || dbName == "" || jwtKey == "" {
		log.Fatalf("invalid env parameters")
	}

	connStr := ""
	connStr += " user=" + dbUser
	connStr += " dbname=" + dbName
	connStr += " password=" + dbPassword

	if test == "1" {
		connStr += " sslmode=disable"
	} else {
		connStr += " host=" + dbAddress + ":" + dbPort
		connStr += " sslmode=require"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrator, err := migration.NewMigrator(db, &migrationsFS)
	if err != nil {
		log.Fatal(err)
	}

	err = migrator.RunMigrations()
	if err != nil {
		log.Fatal(err)
	}

	templateSub, err := fs.Sub(templateFS, "template")
	if err != nil {
		log.Fatal(err)
	}

	var staticSub fs.FS
	if test == "1" {
		staticSub = os.DirFS("static")
	} else {
		staticSub, err = fs.Sub(staticFS, "static")
		if err != nil {
			log.Fatal(err)
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
