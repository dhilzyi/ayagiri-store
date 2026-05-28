package api

import (
	"time"

	"github.com/google/uuid"
)

type ProductRequest struct {
	Name        string
	EnglishName string `json:"english_name"`
	Price       int32
	Discount    int32
	CategoryID  int32 `json:"category_id"`
}

type CategoryRequest struct {
	Name        string
	EnglishName string `json:"english_name"`
}

type CategoryResponse struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	EnglishName string    `json:"english_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrderResponse struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	TableID       int32     `json:"table_id"`
	OrderComplete bool      `json:"order_complete"`
}

type OrderItemsResponse struct {
	ID            int32     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	OrderID       uuid.UUID `json:"order_id"`
	ProductID     int32     `json:"product_id"`
	Quantity      int32     `json:"quantity"`
	ProductName   string    `json:"product_name"`
	OrderComplete bool      `json:"order_complete"`
}
