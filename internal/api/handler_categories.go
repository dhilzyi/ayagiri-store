package api

import (
	"context"
	"fmt"
	"net/http"

	"restaurant/internal/database"
)

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var payloadCategory CategoryRequest
	if err := decodeJson(r.Body, &payloadCategory); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	params := database.CreateCategoryParams{
		Name:        payloadCategory.Name,
		EnglishName: payloadCategory.EnglishName,
	}
	category, err := h.db.CreateCategory(context.Background(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	respondWithJSON(w, http.StatusCreated, category)
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categorys, err := h.db.GetCategories(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	respondWithJSON(w, 200, categorys)
}

func (h *Handler) CreateCategories(w http.ResponseWriter, r *http.Request) {
	var categoryReq []CategoryRequest
	if err := decodeJson(r.Body, &categoryReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	if len(categoryReq) <= 0 {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("request need to one or multiple items"), fmt.Errorf("request need to one or multiple items"))
		return
	}
	result, err := h.db.BulkCreateCategories(context.Background(), toBulkCategories(categoryReq))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, result)
}
