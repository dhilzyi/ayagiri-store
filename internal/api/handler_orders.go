package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant/internal/database"
	"restaurant/internal/domain"
	"restaurant/internal/orders"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid", err)
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
	productItems, err := h.db.GetProductsByID(ctx, toProductIDs(orderReq.Items))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	if len(productItems) <= 0 {
		respondWithError(w, http.StatusBadRequest, "at least one item product is required", fmt.Errorf("product items is return nothing"))
		return
	}

	if err := h.orderSvc.SendOrderToKitchen(order.ID, order.TableID, mapOrderResponse(orderReq.Items, toProductResponses(productItems))); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, "")
}

func (h *Handler) KitchenStreamListenHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr

	log.Printf("[SSE CONNECT] 🟢 Client: %s", ip)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// You may need this locally for CORS requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for client disconnection
	clientGone := r.Context().Done()

	// KitchenClient or Kitchen channel
	// This channel for receiving an orders. Other handler will send it through the channel
	kitchenClient := h.orderSvc.CreateKitchenClient()

	rc := http.NewResponseController(w)
	for {
		select {
		case <-clientGone:
			fmt.Println("Client disconnected")
			log.Printf("[SSE DISCONNECT] 🔴 Client left: %s", ip)

			h.orderSvc.DeleteKitchenClient(kitchenClient)

			return
		case event, ok := <-kitchenClient:
			if !ok {
				return
			}

			jsonBytes, err := json.Marshal(event)
			if err != nil {
				log.Println(err)
				return
			}

			// IMPORTANT: Wrap the JSON inside the SSE protocol format!
			// Notice the "data: %s\n\n"
			_, err = fmt.Fprintf(w, "data: %s\n\n", string(jsonBytes))
			if err != nil {
				log.Println(err)
				return
			}

			// Flush the buffer
			err = rc.Flush()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (h *Handler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("order_id") // Read their query param
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid", err)
		return
	}

	if _, err := h.db.OrderComplete(context.Background(), database.OrderCompleteParams{
		OrderComplete: true,
		ID:            orderID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if err := h.orderSvc.CompleteOrder(orderID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
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
