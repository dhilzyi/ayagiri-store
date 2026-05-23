package orders

import (
	"fmt"
	"sync"

	"restaurant/internal/domain"

	"github.com/google/uuid"
)

type OrderService struct {
	sync.RWMutex

	// The active orders
	orders map[uuid.UUID]*Order

	// Kitchen tablets listening for ANY new order
	kitchenClients map[chan Event]bool

	// Customers listening for THEIR order (Map key is OrderID)
	customerClients map[uuid.UUID]chan Event
}

type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Order struct {
	TableID int32       `json:"table_id"`
	Items   []OrderItem `json:"items"`
}

type OrderItem struct {
	Product  domain.ProductResponse `json:"product"`
	Quantity int32                  `json:"quantity"`
}

func NewOrderSrv() *OrderService {
	return &OrderService{
		kitchenClients:  make(map[chan Event]bool),
		customerClients: make(map[uuid.UUID]chan Event),
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

	return nil
}

func (o *OrderService) CompleteOrder(orderID uuid.UUID) error {
	o.Lock()
	defer o.Unlock()

	_, exists := o.orders[orderID]
	if !exists {
		return fmt.Errorf("order does not exist")
	}

	return nil
}
