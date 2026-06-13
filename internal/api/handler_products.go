package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"restaurant/internal/database"

	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	var err error
	var products []database.Product

	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			respondWithError(w, http.StatusUnprocessableEntity, "Invalid categoryID", err)
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
	if err := h.v.Struct(productReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	ctx := context.Background()
	product, err := h.db.CreateProduct(ctx, toProductRequest(productReq))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	productJoin, err := h.db.GetProductJoinByID(ctx, product.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	respondWithJSON(w, http.StatusCreated, toProductResponseJoin(productJoin))
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

func (h *Handler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	idsStr := r.URL.Query().Get("ids")
	if idsStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing ids parameter", nil)
		return
	}
	parts := strings.Split(idsStr, ",")

	var ids []int32
	for _, part := range parts {
		id, _ := strconv.Atoi(part)
		ids = append(ids, int32(id))
	}

	if err := h.db.DeleteProductsByID(context.Background(), ids); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("productID")
	if productIDStr == "" {
		respondWithError(w, http.StatusUnprocessableEntity, "productId in the path is necessary", fmt.Errorf(""))
		return
	}
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, err.Error(), err)
		return
	}
	var productReq ProductRequest
	if err := decodeJson(r.Body, &productReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	if err := h.v.Struct(productReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	param := database.UpdateProductByIDParams{
		Name:        productReq.Name,
		Price:       *productReq.Price,
		Discount:    pgtype.Int4{Valid: true, Int32: productReq.Discount},
		ID:          int32(productID),
		CategoryID:  *productReq.CategoryID,
		EnglishName: pgtype.Text{Valid: true, String: productReq.EnglishName},
	}

	ctx := context.Background()
	product, err := h.db.UpdateProductByID(ctx, param)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	productJoin, err := h.db.GetProductJoinByID(ctx, product.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	// TODO: Right now the category name is undefined in the front end
	// Just make the mapping category name in the front end at this point. How many helpers function do i need?
	// Kinda dilemma right now
	respondWithJSON(w, http.StatusOK, toProductResponseJoin(productJoin))
}
