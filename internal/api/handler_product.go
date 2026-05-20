package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"restaurant/internal/database"
)

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	var err error
	var products []database.Product

	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid categoryID", err)
			return
		}
		products, err = h.db.GetProductByCategoryID(context.Background(), sql.NullInt32{Int32: int32(categoryID), Valid: true})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error(), err)
			return
		}

		respondWithJSON(w, http.StatusOK, toProductResponses(products))
		return
	} else {
		products, err = h.db.GetProducts(context.Background())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error(), err)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, toProductResponses(products))
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productReq ProductRequest
	if err := decodeJson(r.Body, &productReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	product, err := h.db.CreateProduct(context.Background(), toProductRequest(productReq))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	respondWithJSON(w, http.StatusCreated, product)
}

func (h *Handler) CreateMultipleProducts(w http.ResponseWriter, r *http.Request) {
	var productReq []ProductRequest
	if err := decodeJson(r.Body, &productReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	var productRes []ProductResponse
	for i := range productReq {
		product, err := h.db.CreateProduct(context.Background(), toProductRequest(productReq[i]))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error(), err)
			return
		}
		productRes = append(productRes, toProductResponse(product))
	}
	respondWithJSON(w, http.StatusCreated, productRes)
}
