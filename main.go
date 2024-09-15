package main

import (
	"database/sql"
	"embed"
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"decentrala.org/events/internal/controller"
	"decentrala.org/events/internal/migration"
	"decentrala.org/events/internal/model"
	"decentrala.org/events/internal/types"
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
	dbConn := os.Getenv("DB_CONN")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")
	jwtKey := os.Getenv("JWT_SECRET")

	if jwtKey == "" {
		log.Fatalf("invalid jwt secret parameters")
	}

	connStr := ""

	if dbConn == "" {
		if dbUser == "" || dbPassword == "" || dbName == "" {
			log.Fatalf("invalid env parameters")
		}

		connStr += " user=" + dbUser
		connStr += " dbname=" + dbName
		connStr += " password=" + dbPassword

		if test == "1" {
			connStr += " sslmode=disable"
		} else {
			connStr += " host=" + dbAddress + ":" + dbPort
			connStr += " sslmode=require"
		}
	} else {
		connStr = dbConn
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

func ProcessCSV(db *sql.DB) {
	c, _ := os.Open("dogadjaji.csv")

	model := model.NewModel(db)

	s := csv.NewReader(c)

	ss, _ := s.ReadAll()

	db.Exec(`TRUNCATE TABLE event`)

	for _, row := range ss {

		var ty types.EventType

		if row[4] == " workshop" {
			ty = types.EventTypeWorkshop
		} else if row[4] == " movie" {
			ty = types.EventTypeMovie
		} else if row[4] == " hack" {
			ty = types.EvenTypeHackathon
		} else if row[4] == " lecture" {
			ty = types.EventTypeLecture
		} else {
			ty = types.EvenTypeOther
		}

		time, e := time.Parse("02-01-2006 15:04", row[0]+row[1])
		if e != nil {
			fmt.Println(e)
			continue
		}

		event := types.Event{
			StartsAt:         time,
			Name:             row[3],
			OrganizationCode: "dmz",
			EventType:        ty,
		}

		_, e = model.CreateEvent(event)
		if e != nil {
			fmt.Println(e)
		} else {
			fmt.Println("ok", event.Name)
		}
	}
}
