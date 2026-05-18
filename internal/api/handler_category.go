package api

import (
	"context"
	"net/http"

	"restaurant/internal/database"
)

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var payloadCategory PayloadCategoryParams
	if err := decodeJson(r.Body, &payloadCategory); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	params := database.CreateCategoryParams{
		Name:       payloadCategory.Name,
		RomajiName: payloadCategory.RomajiName,
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
