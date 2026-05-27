package api

import (
	"context"
	"fmt"
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
		products, err = h.db.GetProductByCategoryID(context.Background(), int32(categoryID))
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

func (h *Handler) CreateProducts(w http.ResponseWriter, r *http.Request) {
	var productReq []ProductRequest
	if err := decodeJson(r.Body, &productReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	if len(productReq) <= 0 {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("at least one product is required"), fmt.Errorf("at least one product is required"))
		return
	}
	result, err := h.db.BulkCreateProducts(context.Background(), toBulkProducts(productReq))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}

func (h *Handler) ListProductsAdmin(w http.ResponseWriter, r *http.Request) {
	products, err := h.db.GetProductsJoin(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, toProductResponsesAdmin(products))
}
