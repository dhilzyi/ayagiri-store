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
	// Safety check if productItems returning nothing from database
	if len(orderReq.Items) <= 0 {
		respondWithError(w, http.StatusBadRequest, "at least one item product is required", fmt.Errorf("product items is empty from request: %s", orderID))
		return
	}
	ctx := context.Background()

	// Create an order to the database
	order, err := h.db.CreateOrder(ctx, toOrderRequest(orderID, orderReq))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	// Create an order items to the database
	if _, err := h.db.BulkCreateOrderItem(ctx, toBulkOrderItemRequest(order.ID, orderReq.Items)); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	// Add new order to OrderService
	if err := h.orderSvc.AddNewOrder(order.ID, orderReq); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	// Retrieve the Products by Ids to validate. Never trust frontend as they said
	productItems, err := h.db.GetProductsByID(ctx, toProductIDs(orderReq.Items))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	// Safety check if productItems returning nothing from database
	if len(productItems) <= 0 {
		respondWithError(w, http.StatusBadRequest, "at least one item product is required", fmt.Errorf("product items is empty"))
		return
	}

	// Create new order and broadcast to all existing kitchens client
	event := h.orderSvc.CreateNewOrderEvent(order.ID, order.TableID, mapOrderResponse(orderReq.Items, toProductResponses(productItems)))
	if err := h.orderSvc.BroadcastToKitchens(event); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusCreated, "")
}

func (h *Handler) KitchenStreamListenHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr

	log.Printf("[SSE CONNECT] 🟢 Kitchen Client: %s", ip)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for client disconnection
	clientGone := r.Context().Done()

	// KitchenClient or Kitchen channel
	// This channel for receiving an orders. Other handler or customer client will send it through the channel
	kitchenClient := h.orderSvc.CreateKitchenClient()

	rc := http.NewResponseController(w)
	for {
		select {
		case <-clientGone:
			fmt.Println("Kitchen Client disconnected")
			log.Printf("[SSE DISCONNECT] 🔴 Kitchen Client left: %s", ip)

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

// order_id as parameter is required
// Customer will listen to their orderID channel. Signal will be sent from kitchen when the orders is completed
func (h *Handler) CustomerStreamListenHandler(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid", err)
		return
	}
	ip := r.RemoteAddr
	log.Printf("[SSE CONNECT] 🟢 Customer Client: %s\nOrderID: %s", ip, orderID.String())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	clientGone := r.Context().Done()
	customerChan, err := h.orderSvc.CreateCustomerClient(orderID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
	}
	rc := http.NewResponseController(w)
	for {
		select {
		case <-clientGone:
			fmt.Println("Customer client disconnected")
			log.Printf("[SSE DISCONNECT] 🔴 Customer Client left: %s\nOrderID: %s", ip, orderID.String())

			h.orderSvc.DeleteCustomerClient(orderID, customerChan)

			return
		case event, ok := <-customerChan:
			if !ok {
				return
			}

			jsonBytes, err := json.Marshal(event)
			if err != nil {
				log.Println(err)
				return
			}

			// SSE format
			_, err = fmt.Fprintf(w, "data: %s\n\n", string(jsonBytes))
			if err != nil {
				log.Println(err)
				return
			}

			err = rc.Flush()
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// order_id as parameter is required
func (h *Handler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid", err)
		return
	}

	// Update the order as complete to database
	updateOrder, err := h.db.OrderComplete(context.Background(), database.OrderCompleteParams{
		OrderComplete: true,
		ID:            orderID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	// Create a Complete event for Client and Kitchen
	event := orders.Event{
		Type: "COMPLETE_ORDER",
		Payload: orders.OrderCustomerResponse{
			OrderID: updateOrder.ID,
			TableID: updateOrder.TableID,
		},
	}
	clientCh, exist := h.orderSvc.GetCustomerChannel(updateOrder.ID)
	if exist {
		select {
		case clientCh <- event:
			// Sent to customer client successfully
		default:
			log.Printf("[WARNING] Customer %s channel blocked, skipping.", orderID)
		}
	} else {
		log.Printf("[WARNING] Customer %s channel does not exist, skipping.", orderID)

	}

	// Send the event to kitchen aswell for sync in case there's multiple kitchen clients
	// TODO: Send to all kitchen except the one who sent it
	if err := h.orderSvc.BroadcastToKitchens(event); err != nil {
		log.Println("failed to broadcast to kitchens: ", err)
	}

	respondWithJSON(w, http.StatusOK, "")
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
