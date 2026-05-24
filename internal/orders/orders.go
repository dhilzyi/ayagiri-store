package orders

import (
	"fmt"
	"restaurant/internal/domain"
	"sync"

	"github.com/google/uuid"
)

type OrderService struct {
	sync.RWMutex

	// The active orders
	orders map[uuid.UUID]*Order

	// Kitchen clients listening for ANY new order
	kitchenClients map[chan Event]interface{}

	// Customers listening for THEIR order (Map key is OrderID)
	customerClients map[uuid.UUID]chan interface{}
}

type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Order struct {
	TableID int32              `json:"table_id"`
	Items   []OrderItemRequest `json:"items"`
}

type OrderKitchenResponse struct {
	OrderID uuid.UUID           `json:"order_id"`
	TableID int32               `json:"table_id"`
	Items   []OrderItemResponse `json:"items"`
}

type OrderItemRequest struct {
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type OrderItemResponse struct {
	Quantity int32                  `json:"quantity"`
	Products domain.ProductResponse `json:"products"`
}

func NewOrderSrv() *OrderService {
	return &OrderService{
		kitchenClients:  make(map[chan Event]any),
		customerClients: make(map[uuid.UUID]chan any),
		orders:          make(map[uuid.UUID]*Order),
	}
}

func (o *OrderService) AddNewOrder(orderID uuid.UUID, order Order) error {
	o.Lock()
	defer o.Unlock()

	_, exists := o.orders[orderID]
	if exists {
		return fmt.Errorf("order is already exist")
	}

	o.orders[orderID] = &order
	o.customerClients[orderID] = make(chan any)

	return nil
}

func (o *OrderService) SendOrderToKitchen(orderID uuid.UUID, tableID int32, productItems []OrderItemResponse) error {
	if len(o.kitchenClients) <= 0 {
		return fmt.Errorf("no kitchen client is running")
	}

	event := Event{
		Type: "ORDER",
		Payload: OrderKitchenResponse{
			TableID: tableID,
			OrderID: orderID,
			Items:   productItems,
		},
	}

	for ch := range o.kitchenClients {
		ch <- event
	}

	return nil
}

func (o *OrderService) CompleteOrder(orderID uuid.UUID) error {
	o.Lock()
	defer o.Unlock()

	clientCh, exists := o.customerClients[orderID]
	if !exists {
		return fmt.Errorf("order does not exist")
	}
	clientCh <- struct{}{}

	return nil
}

func (o *OrderService) CreateKitchenClient() chan Event {
	o.Lock()
	defer o.Unlock()

	ch := make(chan Event)
	o.kitchenClients[ch] = struct{}{}

	return ch
}

func (o *OrderService) DeleteKitchenClient(ch chan Event) {
	o.Lock()
	defer o.Unlock()

	delete(o.kitchenClients, ch)
	close(ch)
}
