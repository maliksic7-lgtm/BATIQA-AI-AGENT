package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"batiqa-ai/internal/model"
	"batiqa-ai/internal/repository"
)

// RestaurantHandler serves the room-service ordering flow:
//
//	POST /api/orders          (guest)  place an order
//	GET  /api/orders          (either) guest sees own, staff sees all
//	PATCH /api/orders/{id}/status (staff) update lifecycle status
//	GET  /api/menu            (public) room-service menu catalog
type RestaurantHandler struct {
	orders *repository.RestaurantOrderRepository
}

func NewRestaurantHandler(orders *repository.RestaurantOrderRepository) *RestaurantHandler {
	return &RestaurantHandler{orders: orders}
}

// MenuHandler serves GET /api/menu (public catalog, no DB needed).
func MenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"items": model.RoomServiceMenu})
}

// OrderRequest is the POST /api/orders body.
type OrderRequest struct {
	SessionID  string         `json:"session_id"`
	RoomNumber string         `json:"room_number,omitempty"`
	Items      []OrderItemArg `json:"items"`
	Note       string         `json:"note,omitempty"`
}

type OrderItemArg struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func (h *RestaurantHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}
	req.SessionID = strings.TrimSpace(req.SessionID)
	req.Note = strings.TrimSpace(req.Note)

	if req.SessionID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "session_id is required")
		return
	}
	// Verified token room wins, then claimed, else reject.
	room := GuestRoom(r)
	if room == "" {
		room = strings.TrimSpace(req.RoomNumber)
	}
	if room == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "room_number is required")
		return
	}
	if len(req.Items) == 0 {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "items is required")
		return
	}
	if len(req.Items) > 20 {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "too many items (max 20)")
		return
	}

	items := []model.OrderItem{}
	seen := map[string]bool{}
	for _, it := range req.Items {
		name := strings.TrimSpace(it.Name)
		mi, ok := model.MenuItemByName(name)
		if !ok {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "unknown menu item: "+name)
			return
		}
		qty := it.Quantity
		if qty < 1 || qty > 50 {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid quantity for "+name)
			return
		}
		key := strings.ToLower(mi.Name)
		if seen[key] {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "duplicate item in order")
			return
		}
		seen[key] = true
		items = append(items, model.OrderItem{Name: mi.Name, Quantity: qty, Price: mi.Price})
	}

	order := model.RestaurantOrder{
		SessionID:  req.SessionID,
		RoomNumber: room,
		Items:      items,
		Note:       req.Note,
	}
	if err := h.orders.Create(&order); err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to place order")
		return
	}
	WriteJSON(w, http.StatusCreated, order)
}

func (h *RestaurantHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	guestRoom := GuestRoom(r)
	statusF := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("status")))
	limit := 100

	var orders []model.RestaurantOrder
	var err error
	if guestRoom != "" {
		// Guest path: only their own session's orders. GuestRoom alone doesn't
		// scope by session, so filter by claimed session_id query param.
		sessionID := strings.TrimSpace(r.URL.Query().Get("session_id"))
		orders, err = h.orders.ListBySession(sessionID, 100)
	} else {
		orders, err = h.orders.ListAll(statusF, limit)
	}
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list orders")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{"orders": orders})
}

// UpdateStatus is PATCH /api/orders and /api/orders/{id}/status — staff only.
func (h *RestaurantHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	idPath := extractSegment(r.URL.Path, "/api/orders/", "/status")
	id, ierr := strconv.Atoi(idPath)
	if ierr != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid order id")
		return
	}
	type StatusBody struct {
		Status string `json:"status"`
	}
	var body StatusBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}
	order, ok, err := h.orders.UpdateStatus(int64(id), strings.ToUpper(strings.TrimSpace(body.Status)))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update order")
		return
	}
	if !ok {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid status")
		return
	}
	WriteJSON(w, http.StatusOK, order)
}
