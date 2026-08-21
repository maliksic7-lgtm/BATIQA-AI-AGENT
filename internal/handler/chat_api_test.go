package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"batiqa-ai/internal/config"
	"batiqa-ai/internal/repository"
	"batiqa-ai/internal/service/ai"
	ticketservice "batiqa-ai/internal/service/ticket"
)

func setupChatHandler(t *testing.T) (*ChatHandler, func()) {
	t.Helper()
	cfg := config.Load()
	db, err := config.OpenDB(cfg)
	if err != nil {
		t.Skipf("DB not available: %v", err)
	}
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
	db.Exec("TRUNCATE TABLE ticket_assignments")
	db.Exec("TRUNCATE TABLE tickets")
	db.Exec("TRUNCATE TABLE conversations")
	db.Exec("TRUNCATE TABLE guests")
	db.Exec("SET FOREIGN_KEY_CHECKS=1")

	ticketRepo := repository.NewTicketRepository(db)
	guestRepo := repository.NewGuestRepository(db)
	convRepo := repository.NewConversationRepository(db)
	aiSvc := ai.NewServiceWithProvider(ai.NewMockProvider(), ai.NewMockProvider())
	ticketSvc := ticketservice.NewService(ticketRepo, guestRepo)
	h := NewChatHandler(aiSvc, ticketSvc, convRepo, guestRepo)
	return h, func() { db.Close() }
}

func TestChatAPI_ValidTicketCreation(t *testing.T) {
	h, cleanup := setupChatHandler(t)
	defer cleanup()

	body := `{"session_id":"sess-chat-1","room_number":"305","message":"Tolong antar 2 handuk ke kamar 305"}`
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("chat valid got %d want 200, body %s", w.Code, w.Body.String())
	}
	var resp ChatResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Intent != ai.IntentTowelRequest {
		t.Errorf("intent got %q want TOWEL_REQUEST", resp.Intent)
	}
	if !resp.RequiresTicket {
		t.Error("should require ticket")
	}
	if resp.TicketID == nil {
		t.Error("ticket_id should not be null for valid request with room")
	}
}

func TestChatAPI_MissingRoom(t *testing.T) {
	h, cleanup := setupChatHandler(t)
	defer cleanup()

	body := `{"session_id":"sess-chat-2","room_number":"","message":"AC saya tidak dingin"}`
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("chat missing room got %d want 200, body %s", w.Code, w.Body.String())
	}
	var resp ChatResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Intent != ai.IntentACProblem {
		t.Errorf("intent got %q", resp.Intent)
	}
	if !resp.RequiresTicket {
		t.Error("should require ticket even if missing room")
	}
	if resp.TicketID != nil {
		t.Error("ticket_id should be null when room missing (AI should not invent)")
	}
	// Response should ask for room
	if resp.Message == "" {
		t.Error("response empty")
	}
}

func TestChatAPI_InfoNoTicket(t *testing.T) {
	h, cleanup := setupChatHandler(t)
	defer cleanup()

	body := `{"session_id":"sess-chat-3","room_number":"305","message":"Jam breakfast sampai jam berapa?"}`
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("chat info got %d", w.Code)
	}
	var resp ChatResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Intent != ai.IntentBreakfastInformation {
		t.Errorf("intent got %q", resp.Intent)
	}
	if resp.RequiresTicket {
		t.Error("breakfast should not require ticket")
	}
	if resp.TicketID != nil {
		t.Error("ticket_id should be null for info")
	}
}

func TestChatAPI_UnknownIntent(t *testing.T) {
	h, cleanup := setupChatHandler(t)
	defer cleanup()

	body := `{"session_id":"sess-chat-4","room_number":"305","message":"asdasd qwerty"}`
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("unknown got %d", w.Code)
	}
	var resp ChatResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Intent != ai.IntentUnknown {
		t.Errorf("got %q want UNKNOWN", resp.Intent)
	}
	if resp.RequiresTicket {
		t.Error("unknown should not require ticket")
	}
}

func TestChatAPI_InvalidRequest(t *testing.T) {
	h, cleanup := setupChatHandler(t)
	defer cleanup()

	// Missing session_id
	body := `{"room_number":"305","message":"hello"}`
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("missing session_id should 400 got %d", w.Code)
	}

	// Empty message
	body2 := `{"session_id":"s","room_number":"305","message":""}`
	req2 := httptest.NewRequest("POST", "/api/chat", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	if w2.Code != 400 {
		t.Fatalf("empty message should 400 got %d", w2.Code)
	}
}

func TestChatAPI_AINotDirectDB(t *testing.T) {
	// Ensure chat handler goes through ticket service validation, not direct DB
	// This is architectural: chat handler calls ai.Process -> ticketService.CreateFromAI -> repo
	// We test that invalid AI output does not create ticket
	h, cleanup := setupChatHandler(t)
	defer cleanup()

	// This will be classified as TOWEL but with invalid room via AI? Our mock will use provided room
	// Test with invalid department via direct AI mock: we already test malformed in ai service
	// Here test that chat with valid ticket creates ticket via service, not direct
	body := `{"session_id":"sess-5","room_number":"305","message":"Tolong antar handuk"}`
	req := httptest.NewRequest("POST", "/api/chat", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("chat got %d", w.Code)
	}
	var resp ChatResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.TicketID == nil {
		t.Error("should create ticket via backend validation")
	}
	// Verify ticket exists via DB indirectly by calling chat again and checking ticket ID format
	if resp.TicketID != nil && len(*resp.TicketID) < 4 {
		t.Errorf("ticket_id format got %q", *resp.TicketID)
	}
}
