package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/immanuel-254/myauth/cmd"
	"github.com/immanuel-254/myauth/internal"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func main() {
	// connect to db
	db, err := sql.Open("sqlite3", os.Getenv("DB"))
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer func() {
		if closeError := db.Close(); closeError != nil {
			fmt.Println("Error closing database", closeError)
			if err == nil {
				err = closeError
			}
		}
	}()

	internal.DB = db

	// migrate to database
	goose.SetDialect("sqlite3")

	// Apply all "up" migrations
	err = goose.Up(internal.DB, "internal/migrations")
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")

	// get alpine js
	var client = &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Create or open a file for writing
	file, err := os.Create("static/styles.css")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(string(body))
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 1 {
		panic("There has to be exactly one argument")
	} else {
		if os.Args[1] == "createadmin" {
			cmd.CreateAdminUser()
		} else if os.Args[1] == "runserver" {
			cmd.Api()
		} else {
			panic("Invalid Argument")
		}
	}

}
