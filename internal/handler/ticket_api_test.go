package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"batiqa-ai/internal/config"
	"batiqa-ai/internal/repository"
	ticketservice "batiqa-ai/internal/service/ticket"
)

func setupTicketHandler(t *testing.T) (*TicketHandler, func()) {
	t.Helper()
	cfg := config.Load()
	cfg.MongoDB = testDBName(t)
	db, closeDB, err := config.ConnectMongo(cfg)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	repo := repository.NewTicketRepository(db)
	guestRepo := repository.NewGuestRepository(db)
	svc := ticketservice.NewService(repo, guestRepo)
	h := NewTicketHandler(svc)
	cleanup := func() {
		_ = db.Drop(context.Background())
		closeDB()
	}
	return h, cleanup
}

func TestTicketAPI_CreateValid(t *testing.T) {
	h, cleanup := setupTicketHandler(t)
	defer cleanup()

	body := `{"room_number":"305","department":"HOUSEKEEPING","category":"TOWEL_REQUEST","description":"Minta 2 handuk","priority":"MEDIUM"}`
	req := httptest.NewRequest("POST", "/api/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Create valid got %d want 201, body %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["ticket_number"] == nil {
		t.Error("ticket_number missing")
	}
	if resp["status"] != "OPEN" {
		t.Errorf("status got %v", resp["status"])
	}
}

func TestTicketAPI_MissingRoom(t *testing.T) {
	h, cleanup := setupTicketHandler(t)
	defer cleanup()

	body := `{"department":"ENGINEERING","category":"AC_PROBLEM","description":"AC tidak dingin","priority":"HIGH"}`
	req := httptest.NewRequest("POST", "/api/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("missing room got %d want 400, body %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if errObj, ok := resp["error"]; !ok || errObj == nil {
		t.Error("expected error envelope")
	}
}

func TestTicketAPI_InvalidDepartment(t *testing.T) {
	h, cleanup := setupTicketHandler(t)
	defer cleanup()

	body := `{"room_number":"305","department":"INVALID","category":"AC_PROBLEM","description":"test","priority":"HIGH"}`
	req := httptest.NewRequest("POST", "/api/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid dept got %d want 400", w.Code)
	}
}

func TestTicketAPI_InvalidPriority(t *testing.T) {
	h, cleanup := setupTicketHandler(t)
	defer cleanup()

	body := `{"room_number":"305","department":"ENGINEERING","category":"AC_PROBLEM","description":"test","priority":"ULTRA"}`
	req := httptest.NewRequest("POST", "/api/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid priority got %d want 400", w.Code)
	}
}

func TestTicketAPI_GetAndUpdate(t *testing.T) {
	h, cleanup := setupTicketHandler(t)
	defer cleanup()

	// Create
	body := `{"room_number":"212","department":"ENGINEERING","category":"TV_PROBLEM","description":"TV mati","priority":"MEDIUM"}`
	req := httptest.NewRequest("POST", "/api/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)
	if w.Code != 201 {
		t.Fatalf("create failed %d", w.Code)
	}
	var created map[string]interface{}
	json.NewDecoder(w.Body).Decode(&created)
	ticketNumber := created["ticket_number"].(string)

	// Get detail
	req2 := httptest.NewRequest("GET", "/api/tickets/"+ticketNumber, nil)
	w2 := httptest.NewRecorder()
	h.GetDetail(w2, req2)
	if w2.Code != 200 {
		t.Fatalf("get detail got %d want 200, body %s", w2.Code, w2.Body.String())
	}
	var detail map[string]interface{}
	json.NewDecoder(w2.Body).Decode(&detail)
	if detail["ticket_number"] != ticketNumber {
		t.Errorf("detail mismatch")
	}

	// List with filter
	req3 := httptest.NewRequest("GET", "/api/tickets?department=ENGINEERING", nil)
	w3 := httptest.NewRecorder()
	h.List(w3, req3)
	if w3.Code != 200 {
		t.Fatalf("list got %d", w3.Code)
	}
	var list map[string]interface{}
	json.NewDecoder(w3.Body).Decode(&list)
	if tickets, ok := list["tickets"].([]interface{}); !ok || len(tickets) == 0 {
		t.Error("list should have tickets")
	}

	// Update status OPEN -> IN_PROGRESS
	patchBody := `{"status":"IN_PROGRESS"}`
	req4 := httptest.NewRequest("PATCH", "/api/tickets/"+ticketNumber+"/status", bytes.NewBufferString(patchBody))
	req4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	h.UpdateStatus(w4, req4)
	if w4.Code != 200 {
		t.Fatalf("update status got %d want 200, body %s", w4.Code, w4.Body.String())
	}
	var updated map[string]interface{}
	json.NewDecoder(w4.Body).Decode(&updated)
	if updated["status"] != "IN_PROGRESS" {
		t.Errorf("status got %v", updated["status"])
	}

	// Invalid status transition: IN_PROGRESS -> try invalid status
	patchBody2 := `{"status":"INVALID_STATUS"}`
	req5 := httptest.NewRequest("PATCH", "/api/tickets/"+ticketNumber+"/status", bytes.NewBufferString(patchBody2))
	req5.Header.Set("Content-Type", "application/json")
	w5 := httptest.NewRecorder()
	h.UpdateStatus(w5, req5)
	if w5.Code != 400 {
		t.Fatalf("invalid status should 400 got %d", w5.Code)
	}

	// RESOLVED then try to revert to OPEN (should fail)
	patchBody3 := `{"status":"RESOLVED"}`
	req6 := httptest.NewRequest("PATCH", "/api/tickets/"+ticketNumber+"/status", bytes.NewBufferString(patchBody3))
	req6.Header.Set("Content-Type", "application/json")
	w6 := httptest.NewRecorder()
	h.UpdateStatus(w6, req6)
	if w6.Code != 200 {
		t.Fatalf("resolve got %d", w6.Code)
	}
	patchBody4 := `{"status":"OPEN"}`
	req7 := httptest.NewRequest("PATCH", "/api/tickets/"+ticketNumber+"/status", bytes.NewBufferString(patchBody4))
	req7.Header.Set("Content-Type", "application/json")
	w7 := httptest.NewRecorder()
	h.UpdateStatus(w7, req7)
	if w7.Code != 400 {
		t.Fatalf("invalid transition RESOLVED->OPEN should 400 got %d, body %s", w7.Code, w7.Body.String())
	}
}
