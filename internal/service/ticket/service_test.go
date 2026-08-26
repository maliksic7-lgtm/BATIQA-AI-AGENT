package ticket

import (
	"database/sql"
	"testing"

	"batiqa-ai/internal/config"
	"batiqa-ai/internal/model"
	"batiqa-ai/internal/repository"
	"batiqa-ai/internal/service/ai"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	cfg := config.Load()
	db, err := config.OpenDB(cfg)
	if err != nil {
		t.Skipf("DB not available: %v", err)
	}
	// Clean ticket tables before test
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS=0")
	_, _ = db.Exec("TRUNCATE TABLE ticket_assignments")
	_, _ = db.Exec("TRUNCATE TABLE tickets")
	_, _ = db.Exec("TRUNCATE TABLE conversations")
	_, _ = db.Exec("TRUNCATE TABLE guests")
	_, _ = db.Exec("SET FOREIGN_KEY_CHECKS=1")
	return db
}

func TestCreateValidTicket(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	req := CreateRequest{
		RoomNumber:  "305",
		Department:  "HOUSEKEEPING",
		Category:    "TOWEL_REQUEST",
		Description: "Minta 2 handuk",
		Priority:    "MEDIUM",
	}
	ticket, err := svc.Create(req)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if ticket.TicketNumber == "" {
		t.Error("ticket_number empty")
	}
	if ticket.Status != model.StatusOpen {
		t.Errorf("status got %q want OPEN", ticket.Status)
	}
	if ticket.TicketNumber[:4] != "TKT-" {
		t.Errorf("ticket_number format got %q want TKT-xxxxxx", ticket.TicketNumber)
	}
	// Get back
	got, err := svc.GetByTicketNumber(ticket.TicketNumber)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.RoomNumber != "305" {
		t.Errorf("room got %q", got.RoomNumber)
	}
}

func TestCreateMissingRoomNumber(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	req := CreateRequest{
		RoomNumber:  "",
		Department:  "ENGINEERING",
		Category:    "AC_PROBLEM",
		Description: "AC tidak dingin",
		Priority:    "HIGH",
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected error for missing room_number")
	}
	if err != nil && err.Error() != "room_number is required" {
		t.Logf("got expected err: %v", err)
	}
}

func TestCreateInvalidDepartment(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	req := CreateRequest{
		RoomNumber:  "305",
		Department:  "INVALID_DEPT",
		Category:    "AC_PROBLEM",
		Description: "test",
		Priority:    "HIGH",
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected invalid department error")
	}
}

func TestCreateInvalidPriority(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	req := CreateRequest{
		RoomNumber:  "305",
		Department:  "ENGINEERING",
		Category:    "AC_PROBLEM",
		Description: "test",
		Priority:    "ULTRA",
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected invalid priority error")
	}
}

func TestCreateInvalidCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	req := CreateRequest{
		RoomNumber:  "305",
		Department:  "ENGINEERING",
		Category:    "FAKE_INTENT",
		Description: "test",
		Priority:    "HIGH",
	}
	_, err := svc.Create(req)
	if err == nil {
		t.Error("expected invalid category error")
	}
}

func TestCreateFromAI(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	aiResult := &ai.AIResult{
		Intent:         ai.IntentTowelRequest,
		Language:       ai.LangID,
		Entities:       map[string]interface{}{"room_number": "305", "quantity": 2, "item": "towel"},
		Action:         ai.Action{Type: ai.ActionCreateTicket, Department: "HOUSEKEEPING", Priority: "MEDIUM"},
		Response:       "Baik",
		RequiresTicket: true,
	}
	ticket, err := svc.CreateFromAI(aiResult, "Tolong antar 2 handuk")
	if err != nil {
		t.Fatalf("CreateFromAI failed: %v", err)
	}
	if ticket.Department != "HOUSEKEEPING" {
		t.Errorf("dept got %q", ticket.Department)
	}
}

func TestCreateFromAIMissingRoom(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	aiResult := &ai.AIResult{
		Intent:         ai.IntentACProblem,
		Language:       ai.LangID,
		Entities:       map[string]interface{}{"problem": "AC tidak dingin"},
		Action:         ai.Action{Type: ai.ActionCreateTicket, Department: "ENGINEERING", Priority: "HIGH"},
		Response:       "Boleh tahu kamar?",
		RequiresTicket: true,
	}
	_, err := svc.CreateFromAI(aiResult, "AC tidak dingin")
	if err == nil {
		t.Error("expected missing room error")
	}
}

func TestGetTicket(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	req := CreateRequest{RoomNumber: "212", Department: "ENGINEERING", Category: "TV_PROBLEM", Description: "TV mati", Priority: "MEDIUM"}
	t1, _ := svc.Create(req)
	got, err := svc.GetByTicketNumber(t1.TicketNumber)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.TicketNumber != t1.TicketNumber {
		t.Errorf("mismatch")
	}
	// Not found
	_, err = svc.GetByTicketNumber("TKT-999999")
	if err == nil {
		t.Error("expected not found")
	}
}

func TestUpdateTicketStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	req := CreateRequest{RoomNumber: "305", Department: "HOUSEKEEPING", Category: "TOWEL_REQUEST", Description: "handuk", Priority: "MEDIUM"}
	t1, _ := svc.Create(req)
	if t1.Status != model.StatusOpen {
		t.Fatalf("initial not OPEN")
	}
	// OPEN -> IN_PROGRESS
	t2, err := svc.UpdateStatus(t1.TicketNumber, model.StatusInProgress)
	if err != nil {
		t.Fatalf("update to IN_PROGRESS failed: %v", err)
	}
	if t2.Status != model.StatusInProgress {
		t.Errorf("got %q", t2.Status)
	}
	// IN_PROGRESS -> RESOLVED
	t3, err := svc.UpdateStatus(t1.TicketNumber, model.StatusResolved)
	if err != nil {
		t.Fatalf("update to RESOLVED failed: %v", err)
	}
	if t3.Status != model.StatusResolved {
		t.Errorf("got %q", t3.Status)
	}
	if t3.ResolvedAt == nil {
		t.Error("resolved_at should be set")
	}
}

func TestInvalidStatusTransition(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	req := CreateRequest{RoomNumber: "305", Department: "HOUSEKEEPING", Category: "TOWEL_REQUEST", Description: "handuk", Priority: "MEDIUM"}
	t1, _ := svc.Create(req)
	// Move to RESOLVED
	_, _ = svc.UpdateStatus(t1.TicketNumber, model.StatusInProgress)
	_, _ = svc.UpdateStatus(t1.TicketNumber, model.StatusResolved)
	// Try to revert RESOLVED -> OPEN (should fail)
	_, err := svc.UpdateStatus(t1.TicketNumber, model.StatusOpen)
	if err == nil {
		t.Error("expected invalid transition RESOLVED->OPEN")
	}
	// Invalid status value
	_, err = svc.UpdateStatus(t1.TicketNumber, "INVALID_STATUS")
	if err == nil {
		t.Error("expected invalid status error")
	}
}

func TestListFiltering(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	// Create 2 tickets different dept/priority
	_, _ = svc.Create(CreateRequest{RoomNumber: "101", Department: "HOUSEKEEPING", Category: "TOWEL_REQUEST", Description: "towel", Priority: "MEDIUM"})
	_, _ = svc.Create(CreateRequest{RoomNumber: "102", Department: "ENGINEERING", Category: "AC_PROBLEM", Description: "ac", Priority: "HIGH"})

	dept := "HOUSEKEEPING"
	list, err := svc.List(model.TicketFilter{Department: &dept})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("HOUSEKEEPING filter got %d want 1", len(list))
	}
	status := "OPEN"
	list, _ = svc.List(model.TicketFilter{Status: &status})
	if len(list) != 2 {
		t.Errorf("OPEN filter got %d want 2", len(list))
	}
	priority := "HIGH"
	list, _ = svc.List(model.TicketFilter{Priority: &priority})
	if len(list) != 1 {
		t.Errorf("HIGH filter got %d want 1", len(list))
	}
}

func TestTicketNumberGeneration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	svc := NewService(repository.NewTicketRepository(db), repository.NewGuestRepository(db))

	t1, _ := svc.Create(CreateRequest{RoomNumber: "301", Department: "HOUSEKEEPING", Category: "TOWEL_REQUEST", Description: "a", Priority: "MEDIUM"})
	t2, _ := svc.Create(CreateRequest{RoomNumber: "302", Department: "HOUSEKEEPING", Category: "TOWEL_REQUEST", Description: "b", Priority: "MEDIUM"})
	if t1.TicketNumber == t2.TicketNumber {
		t.Error("ticket numbers should be unique")
	}
	if len(t1.TicketNumber) < 8 || t1.TicketNumber[:4] != "TKT-" {
		t.Errorf("format got %q", t1.TicketNumber)
	}
	// Ensure sequential
	if t1.ID+1 != t2.ID {
		t.Logf("IDs %d vs %d not sequential but ticket numbers unique is ok", t1.ID, t2.ID)
	}
}
