package model

import "time"

// OrderStatus enum per RESTAURANT order flow.
const (
	OrderNew       = "NEW"
	OrderPreparing = "PREPARING"
	OrderCompleted = "COMPLETED"
	OrderCancelled = "CANCELLED"
)

// OrderItem is a single line item within a restaurant order.
type OrderItem struct {
	Name     string `json:"name" bson:"name"`
	Quantity int    `json:"quantity" bson:"quantity"`
	Price    int    `json:"price" bson:"price"`
}

// RestaurantOrder represents an order placed by a guest via the chat/AI or the
// room app. ID doubles as the MongoDB _id; order_number is ORD-{id} zero-padded.
type RestaurantOrder struct {
	ID          int64       `json:"id" bson:"_id"`
	OrderNumber string      `json:"order_number" bson:"order_number"`
	RoomNumber  string      `json:"room_number" bson:"room_number"`
	SessionID   string      `json:"session_id" bson:"session_id"`
	Status      string      `json:"status" bson:"status"`
	Items       []OrderItem `json:"items" bson:"items"`
	TotalPrice  int         `json:"total_price" bson:"total_price"`
	Note        string      `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt   time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" bson:"updated_at"`
}

// IsValidOrderStatus validates an order status enum value.
func IsValidOrderStatus(s string) bool {
	switch s {
	case OrderNew, OrderPreparing, OrderCompleted, OrderCancelled:
		return true
	default:
		return false
	}
}
