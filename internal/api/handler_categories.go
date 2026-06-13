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

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var payloadCategory CategoryRequest
	if err := decodeJson(r.Body, &payloadCategory); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	if err := h.v.Struct(payloadCategory); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	params := database.CreateCategoryParams{
		Name:        payloadCategory.Name,
		EnglishName: pgtype.Text{String: payloadCategory.EnglishName, Valid: true},
	}
	category, err := h.db.CreateCategory(context.Background(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	respondWithJSON(w, http.StatusCreated, toCategoryResponse(category))
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.db.GetCategories(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	respondWithJSON(w, http.StatusOK, toCategoriesResponse(categories))
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

func (h *Handler) DeleteCategories(w http.ResponseWriter, r *http.Request) {
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

	if err := h.db.DeleteCategoriesByID(context.Background(), ids); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
