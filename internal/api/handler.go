package api

import (
	"restaurant/internal/database"
	"restaurant/internal/orders"
)

type Handler struct {
	db       *database.Queries
	platform string
	orderSvc *orders.OrderService
}

func NewHandler(db *database.Queries, platform string) *Handler {
	return &Handler{
		db:       db,
		platform: platform,
		orderSvc: orders.NewOrderSrv(),
	}
}
