package api

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"time"

	"restaurant/internal/database"
	"restaurant/internal/domain"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func decodeJson(raw io.Reader, placeholder any) error {
	decoder := json.NewDecoder(raw)
	err := decoder.Decode(placeholder)
	if err != nil {
		return err
	}

	return nil
}

func toProductResponse(p database.Product) domain.ProductResponse {
	return domain.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		EnglishName: p.EnglishName,
		Price:       p.Price,
		CategoryID:  p.CategoryID,
		Discount:    p.Discount,
	}
}

func toProductResponses(products []database.Product) []domain.ProductResponse {
	result := make([]domain.ProductResponse, len(products))
	for i, p := range products {
		result[i] = toProductResponse(p)
	}
	return result
}

func toProductRequest(p ProductRequest) database.CreateProductParams {
	return database.CreateProductParams{
		Name:        p.Name,
		EnglishName: p.EnglishName,
		Price:       p.Price,
		Discount:    p.Discount,
		CategoryID:  p.CategoryID,
	}
}

func toNullInt32(v *int32) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *v, Valid: true}
}

func toCategoryRequest(c CategoryRequest) database.CreateCategoryParams {
	return database.CreateCategoryParams{
		Name:        c.Name,
		EnglishName: c.EnglishName,
	}
}

func toCategoryResponse(c database.Category) CategoryResponse {
	return CategoryResponse{
		ID:          int(c.ID),
		Name:        c.Name,
		EnglishName: c.EnglishName,
	}
}

func toBulkCategories(cList []CategoryRequest) []database.BulkCreateCategoriesParams {
	var bulks []database.BulkCreateCategoriesParams
	now := time.Now()
	for _, category := range cList {
		timestamp := pgtype.Timestamp{
			Time:  now,
			Valid: true,
		}
		bulks = append(bulks, database.BulkCreateCategoriesParams{
			Name:        category.Name,
			EnglishName: category.EnglishName,
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
		})
	}

	return bulks
}

func toBulkProducts(pList []ProductRequest) []database.BulkCreateProductsParams {
	var bulks []database.BulkCreateProductsParams
	now := time.Now()
	for _, product := range pList {
		timestamp := pgtype.Timestamp{
			Time:  now,
			Valid: true,
		}
		bulks = append(bulks, database.BulkCreateProductsParams{
			Name:        product.Name,
			EnglishName: product.EnglishName,
			CreatedAt:   timestamp,
			UpdatedAt:   timestamp,
			Price:       product.Price,
			Discount:    product.Discount,
			CategoryID:  product.CategoryID,
		})
	}

	return bulks
}
