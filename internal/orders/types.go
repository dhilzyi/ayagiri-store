package orders

import (
	"github.com/google/uuid"
	"restaurant/internal/domain"
)

type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Order struct {
	TableID int32              `json:"table_id"`
	Items   []OrderItemRequest `json:"items"`
	Channel chan Event
}

type OrderKitchenResponse struct {
	OrderID uuid.UUID           `json:"order_id"`
	TableID int32               `json:"table_id"`
	Items   []OrderItemResponse `json:"items"`
}

type OrderCustomerResponse struct {
	OrderID uuid.UUID `json:"order_id"`
	TableID int32     `json:"table_id"`
}

type OrderItemRequest struct {
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type OrderItemResponse struct {
	Quantity int32                  `json:"quantity"`
	Products domain.ProductResponse `json:"products"`
}
