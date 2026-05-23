package api

import (
	"context"
	"net/http"
	"restaurant/internal/database"
	"restaurant/internal/orders"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.PathValue("orderID")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "uuid is invalid", err)
		return
	}
	var orderReq orders.Order
	if err := decodeJson(r.Body, &orderReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	ctx := context.Background()
	order, err := h.db.CreateOrder(ctx, toOrderRequest(orderID, orderReq))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if _, err := h.db.BulkCreateOrderItem(ctx, toBulkOrderItemRequest(order.ID, orderReq.Items)); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if err := h.orderSvc.AddNewOrder(order.ID, orderReq); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, "")
}

func toOrderRequest(orderID uuid.UUID, order orders.Order) database.CreateOrderParams {
	return database.CreateOrderParams{
		ID:      orderID,
		TableID: order.TableID,
	}
}

func toBulkOrderItemRequest(orderID uuid.UUID, orderItems []orders.OrderItem) []database.BulkCreateOrderItemParams {
	var bulks []database.BulkCreateOrderItemParams
	now := pgtype.Timestamp{Time: time.Now(), Valid: true}
	for _, ord := range orderItems {
		bulks = append(bulks, database.BulkCreateOrderItemParams{
			OrderID:   orderID,
			ProductID: ord.Product.ID,
			Quantity:  ord.Quantity,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	return bulks
}
