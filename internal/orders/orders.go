package orders

import (
	"fmt"
	"restaurant/internal/database"
	"sync"

	"github.com/google/uuid"
)

type OrderService struct {
	sync.RWMutex

	// The active orders. It contains Order data and with its their channel
	orders map[uuid.UUID]*Order

	// Kitchen clients listening for ANY new order
	kitchenClients map[chan Event]interface{}
}

func NewOrderSrv() *OrderService {
	return &OrderService{
		kitchenClients: make(map[chan Event]any),
		orders:         make(map[uuid.UUID]*Order),
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

func (o *OrderService) CreateNewOrderEvent(orderData database.Order, productItems []OrderItemResponse) Event {
	return Event{
		Type: "ADD_ORDER",
		Payload: OrderKitchenResponse{
			TableID:   orderData.TableID,
			OrderID:   orderData.ID,
			Items:     productItems,
			CreatedAt: orderData.CreatedAt.Time,
		},
	}
}

func (o *OrderService) BroadcastToKitchens(e Event) error {
	o.RLock()
	if len(o.kitchenClients) <= 0 {
		o.RUnlock()
		return fmt.Errorf("no kitchen client is running")
	}
	channels := make([]chan Event, 0, len(o.kitchenClients))
	for ch := range o.kitchenClients {
		channels = append(channels, ch)
	}
	o.RUnlock()

	for _, ch := range channels {
		select {
		case ch <- e:
			// Sent successfully!
		default:
			// If the channel is blocked, we skip it so we don't freeze.
			// (The slow/dead client will be cleaned up when they disconnect).
		}
	}

	return nil
}

func (o *OrderService) SendEventToClient(e Event, ch chan Event) {
	select {
	case ch <- e:
		// Signal sent successfully because the client was listening
	default:
		// No one was listening (maybe they closed the tab).
	}
}

func (o *OrderService) CreateKitchenClient() chan Event {
	o.Lock()
	defer o.Unlock()

	ch := make(chan Event)
	o.kitchenClients[ch] = struct{}{}

	return ch
}

func (o *OrderService) CreateCustomerClient(orderID uuid.UUID) (chan Event, error) {
	o.Lock()
	defer o.Unlock()

	orderData, exists := o.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order is does not exist in map: %s", orderID)
	}

	ch := make(chan Event)
	orderData.Channel = ch

	return ch, nil
}

func (o *OrderService) DeleteKitchenClient(ch chan Event) {
	o.Lock()
	defer o.Unlock()

	delete(o.kitchenClients, ch)
	close(ch)
}

func (o *OrderService) DeleteCustomerClient(orderID uuid.UUID, ch chan Event) {
	o.Lock()
	defer o.Unlock()

	close(ch)
	delete(o.orders, orderID)
}

func (o *OrderService) GetCustomerChannel(orderID uuid.UUID) (chan Event, bool) {
	o.RLock()
	defer o.RUnlock()
	orderData, exists := o.orders[orderID]
	if !exists {
		return nil, false
	}
	if orderData.Channel == nil {
		return nil, false
	}
	return orderData.Channel, true
}
