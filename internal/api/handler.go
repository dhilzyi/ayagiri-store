package api

import (
	"context"
	"database/sql"
	"net/http"

	"restaurant/internal/database"
)

type Handler struct {
	db       *database.Queries
	platform string
}

func NewHandler(db *database.Queries, platform string) *Handler {
	return &Handler{db: db, platform: platform}
}

func (h *Handler) GETProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.db.GetProducts(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
	}

	respondWithJSON(w, 200, products)
}

func (h *Handler) POSTProduct(w http.ResponseWriter, r *http.Request) {
	var payloadProduct PayloadProductParams
	if err := decodeJson(r.Body, &payloadProduct); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	params := database.CreateProductParams{
		Name:  payloadProduct.Name,
		Price: int32(payloadProduct.Price),
		Discount: sql.NullInt32{
			Valid: true,
			Int32: int32(payloadProduct.Discount),
		},
	}
	product, err := h.db.CreateProduct(context.Background(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	respondWithJSON(w, http.StatusCreated, product)
}
