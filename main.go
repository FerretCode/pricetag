package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	database "github.com/ferretcode/pricetag/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	_ "modernc.org/sqlite"
)

var templates *template.Template

func parseTemplates() error {
	var err error

	files := []string{
		"./views/fragments/navbar.html",
		"./views/home.html",
		"./views/error.html",
	}

	templates, err = template.ParseFiles(files...)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if _, err := os.Stat("./.env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Error("error parsing environment variable", "err", err)
			os.Exit(1)
		}
	}

	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Error("error opening connection to db", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	err = database.RunMigrations(db)
	if err != nil {
		log.Error("error running database migrations", "err", err)

		if !strings.Contains(err.Error(), "already exists") {
			os.Exit(1)
		}
	}

	err = parseTemplates()
	if err != nil {
		log.Error("error parsing templates", "err", err)
		os.Exit(1)
	}

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	registerHandlers(r)

	// TODO: change in production
	// http.ListenAndServe(":"+os.Getenv("PORT"), r)
	http.ListenAndServe("localhost:"+os.Getenv("PORT"), r)
}
