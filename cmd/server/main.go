package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"restaurant/internal/api"
	"restaurant/internal/database"
	"restaurant/internal/middleware"

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

	if err := bootstrapAdmin(dbQueries); err != nil {
		log.Fatal(err)
	}

	// Static web files handler
	mux.Handle("/", http.FileServer(http.Dir("./web/customer")))
	mux.Handle("/kitchen/", http.StripPrefix("/kitchen", http.FileServer(http.Dir("./web/kitchen"))))

	mux.Handle("GET /api/products", middleware.Chain(
		http.HandlerFunc(handler.ListProducts),
		middleware.Logging,
	))
	mux.HandleFunc("POST /api/login", handler.Login)

	// mux.HandleFunc("GET /api/products", handler.ListProducts)
	mux.HandleFunc("POST /api/products", handler.CreateProduct)
	mux.HandleFunc("DELETE /api/products", handler.DeleteProducts)
	mux.HandleFunc("PUT /api/products/{productID}", handler.UpdateProduct)
	mux.HandleFunc("POST /api/products/bulk", handler.CreateProducts)

	// Admin handler. For database admin interface purposes
	// Some of the query use JOIN method to others table
	mux.HandleFunc("GET /api/admin/products", handler.ListProductsAdmin)
	mux.HandleFunc("GET /api/admin/orders", handler.ListOrdersAdmin)
	mux.HandleFunc("GET /api/admin/order_items", handler.ListOrderItemsAdmin)
	mux.HandleFunc("GET /api/admin/categories", handler.ListCategories)

	mux.HandleFunc("POST /api/categories", handler.CreateCategory)
	mux.HandleFunc("DELETE /api/categories", handler.DeleteCategories)
	mux.HandleFunc("POST /api/categories/bulk", handler.CreateCategories)

	// Customer client will send a request to this to create an order
	mux.HandleFunc("POST /api/orders", handler.CreateOrder)
	// Kitchen client will send a signal to this endpoints to complete an orders
	mux.HandleFunc("POST /api/complete_orders", handler.CompleteOrder)

	// SSE handlers.
	// Kitchen client will have connection SSE open to this for new orders from customer client side
	mux.HandleFunc("GET /api/kitchen/stream", handler.KitchenStreamListenHandler)
	// Customer client will have connection open for their orders signal
	mux.HandleFunc("GET /api/orders/stream", handler.CustomerStreamListenHandler)

	fmt.Println("Server on 127.0.0.1" + port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
