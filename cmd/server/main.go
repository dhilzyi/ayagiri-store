package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"restaurant/internal/api"
	"restaurant/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	STORE_NAME = "Ayagiri"
)

func main() {
	godotenv.Load(".env")

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	handler := api.NewHandler(dbQueries, os.Getenv("PLATFORM"))
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    port,
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir("./web/customer")))
	mux.Handle("/kitchen/", http.StripPrefix("/kitchen", http.FileServer(http.Dir("./web/kitchen"))))

	mux.HandleFunc("GET /api/product", handler.ListProducts)
	mux.HandleFunc("POST /api/product", handler.CreateProduct)
	mux.HandleFunc("POST /api/products", handler.CreateMultipleProducts)

	mux.HandleFunc("GET /api/category", handler.ListCategories)
	mux.HandleFunc("POST /api/category", handler.CreateCategory)

	fmt.Println("Server on 127.0.0.1" + port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
