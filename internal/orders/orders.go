package orders

import (
	"sync"
	"time"

	"restaurant/internal/domain"
)

type OrderManager struct {
	sync.RWMutex

	// The active orders
	Orders map[string]*Order

	// Kitchen tablets listening for ANY new order
	KitchenClients map[chan Event]bool

	// Customers listening for THEIR order (Map key is OrderID)
	CustomerClients map[string]chan Event
}

type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Order struct {
	Items     []OrderItem
	CreatedAt time.Time
	Status    string
}

type OrderItem struct {
	Product domain.ProductResponse
	Amount  int32
}
