package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant/internal/database"
	"restaurant/internal/orders"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
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
	event := h.orderSvc.CreateNewOrderEvent(order, mapOrderResponse(orderReq.Items, toProductResponses(productItems)))
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
	// Initial flush for sending a signal
	_, err := fmt.Fprintf(w, "data: {\"type\":\"INITIAL\"}\n\n")
	if err != nil {
		log.Println("failed to write initial SSE payload: ", err)
		return
	}
	if err := rc.Flush(); err != nil {
		log.Println("failed to flush initial SSE headers: ", err)
		return
	}

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

func (h *Handler) ListOrdersAdmin(w http.ResponseWriter, r *http.Request) {
	orders, err := h.db.GetOrders(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, toOrdersResponse(orders))
}

func (h *Handler) ListOrderItemsAdmin(w http.ResponseWriter, r *http.Request) {
	orderItems, err := h.db.GetOrderItems(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, toOrderItemsResponse(orderItems))
}
