package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	NAME = "Ayagiri"
)

type Handlers struct {
	// db *database.
	platform string
}

func main() {
	godotenv.Load(".env")

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db.Ping())
	fmt.Println(db.Stats())

	fmt.Println(dbUrl)
}
