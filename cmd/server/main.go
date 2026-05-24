package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"restaurant/internal/api"
	"restaurant/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const (
	STORE_NAME = "Ayagiri"
)

func main() {
	godotenv.Load(".env")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dbUrl := os.Getenv("DB_URL")
	ctx := context.Background()
	db, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(ctx); err != nil {
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
	mux.HandleFunc("POST /api/products", handler.CreateProduct)
	mux.HandleFunc("POST /api/products/bulk", handler.CreateProducts)

	mux.HandleFunc("GET /api/categories", handler.ListCategories)
	mux.HandleFunc("POST /api/categories", handler.CreateCategory)
	mux.HandleFunc("POST /api/categories/bulk", handler.CreateCategories)

	// Customer client will request to this to create an order
	mux.HandleFunc("POST /api/orders", handler.CreateOrder)

	// SSE handlers. Kitchen client will listen to this for new orders from customer client side
	mux.HandleFunc("GET /api/kitchen/stream", handler.KitchenStreamListenHandler)

	fmt.Println("Server on 127.0.0.1" + port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
