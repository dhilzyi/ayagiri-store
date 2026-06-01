package api

import (
	"restaurant/internal/database"
	"restaurant/internal/orders"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	db       *database.Queries
	platform string
	orderSvc *orders.OrderService
	v        *validator.Validate
}

func NewHandler(db *database.Queries, platform string) *Handler {
	return &Handler{
		db:       db,
		platform: platform,
		orderSvc: orders.NewOrderSrv(),
		v:        validator.New(),
	}
}
