package api

import (
	"io"
	"log"
	"net/http"
	"time"

	"restaurant/internal/database"
	"restaurant/internal/domain"
	"restaurant/internal/orders"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
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

func toProductResponsesAdmin(p []database.GetProductsJoinRow) []domain.ProductResponseAdmin {
	results := make([]domain.ProductResponseAdmin, 0, len(p))
	for i := range p {
		inst := p[i]
		results = append(results, domain.ProductResponseAdmin{
			ProductResponse: toProductResponse(database.Product{
				ID:          inst.ID,
				Name:        inst.Name,
				EnglishName: inst.EnglishName,
				Price:       inst.Price,
				CategoryID:  inst.CategoryID,
				Discount:    inst.Discount,
			}),
			CreatedAt:    inst.CreatedAt.Time,
			UpdatedAt:    inst.UpdatedAt.Time,
			CategoryName: inst.CategoryName,
		})
	}

	return results
}

func toOrderRequest(orderID uuid.UUID, order orders.Order) database.CreateOrderParams {
	return database.CreateOrderParams{
		ID:      orderID,
		TableID: order.TableID,
	}
}

func toBulkOrderItemRequest(orderID uuid.UUID, orderItems []orders.OrderItemRequest) []database.BulkCreateOrderItemParams {
	var bulks []database.BulkCreateOrderItemParams
	now := pgtype.Timestamp{Time: time.Now(), Valid: true}
	for _, ord := range orderItems {
		bulks = append(bulks, database.BulkCreateOrderItemParams{
			OrderID:   orderID,
			ProductID: ord.ProductID,
			Quantity:  ord.Quantity,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	return bulks
}

func toProductIDs(items []orders.OrderItemRequest) []int32 {
	var productIDs []int32
	for _, i := range items {
		productIDs = append(productIDs, i.ProductID)
	}

	return productIDs
}

func mapOrderResponse(reqItems []orders.OrderItemRequest, dbProducts []domain.ProductResponse) []orders.OrderItemResponse {
	qtyMap := make(map[int32]int32)
	for _, reqItem := range reqItems {
		qtyMap[reqItem.ProductID] = reqItem.Quantity
	}

	responses := make([]orders.OrderItemResponse, len(dbProducts))
	for i, prod := range dbProducts {
		quantity := qtyMap[prod.ID]

		responses[i] = orders.OrderItemResponse{
			Quantity: quantity,
			Products: prod,
		}
	}

	return responses
}

func toOrdersResponse(orders []database.Order) []OrderResponse {
	responses := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		responses = append(responses, OrderResponse{
			ID:            o.ID,
			CreatedAt:     o.CreatedAt.Time,
			UpdatedAt:     o.UpdatedAt.Time,
			TableID:       o.TableID,
			OrderComplete: o.OrderComplete,
		})
	}

	return responses
}

func toOrderItemsResponse(orderItems []database.GetOrderItemsRow) []OrderItemsResponse {
	responses := make([]OrderItemsResponse, 0, len(orderItems))
	for _, o := range orderItems {
		responses = append(responses, OrderItemsResponse{
			ID:            o.ID,
			CreatedAt:     o.CreatedAt.Time,
			UpdatedAt:     o.UpdatedAt.Time,
			OrderID:       o.OrderID,
			ProductID:     o.ProductID,
			Quantity:      o.Quantity,
			ProductName:   o.ProductName,
			OrderComplete: o.OrderComplete,
		})
	}

	return responses
}

func toCategoriesResponse(categories []database.Category) []CategoryResponse {
	responses := make([]CategoryResponse, 0, len(categories))
	for _, c := range categories {
		responses = append(responses, CategoryResponse{
			ID:          c.ID,
			CreatedAt:   c.CreatedAt.Time,
			UpdatedAt:   c.UpdatedAt.Time,
			EnglishName: c.EnglishName,
			Name:        c.Name,
		})
	}
	return responses
}
