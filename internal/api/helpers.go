package api

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"restaurant/internal/database"

	"github.com/goccy/go-json"
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

func toProductResponse(p database.Product) ProductResponse {
	return ProductResponse{
		ID:         p.ID,
		Name:       p.Name,
		RomajiName: p.RomajiName,
		Price:      p.Price,
		CategoryID: p.CategoryID.Int32,
		Discount:   p.Discount.Int32,
	}
}

func toProductResponses(products []database.Product) []ProductResponse {
	result := make([]ProductResponse, len(products))
	for i, p := range products {
		result[i] = toProductResponse(p)
	}
	return result
}

func toProductRequest(p ProductRequest) database.CreateProductParams {
	return database.CreateProductParams{
		Name:       p.Name,
		RomajiName: p.RomajiName,
		Price:      int32(p.Price),
		Discount:   toNullInt32(p.Discount),
		CategoryID: toNullInt32(p.CategoryID),
	}
}

func toNullInt32(v *int32) sql.NullInt32 {
	if v == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *v, Valid: true}
}
