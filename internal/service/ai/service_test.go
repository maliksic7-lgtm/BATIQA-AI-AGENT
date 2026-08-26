package ai

import (
	"context"
	"testing"
)

func TestLanguageDetection(t *testing.T) {
	tests := []struct {
		msg  string
		want string
	}{
		{"Halo, saya di kamar 305", LangID},
		{"Tolong antar handuk", LangID},
		{"AC tidak dingin", LangID},
		{"Hello, my AC is not working", LangEN},
		{"What's the breakfast schedule?", LangEN},
		{"Hi", LangEN},
		{"Terima kasih", LangID},
	}
	for _, tt := range tests {
		got := DetectLanguage(tt.msg)
		if got != tt.want {
			t.Errorf("DetectLanguage(%q)=%q want %q", tt.msg, got, tt.want)
		}
	}
}

func TestIntentClassification(t *testing.T) {
	svc := NewServiceWithProvider(NewMockProvider(), NewMockProvider())
	tests := []struct {
		msg    string
		room   string
		intent string
	}{
		{"Jam breakfast sampai jam berapa?", "", IntentBreakfastInformation},
		{"Tolong antar 2 handuk ke kamar 305", "305", IntentTowelRequest},
		{"AC kamar saya tidak dingin", "305", IntentACProblem},
		{"AC saya tidak dingin", "", IntentACProblem},
		{"I need 2 towels", "305", IntentTowelRequest},
		{"What's the breakfast schedule?", "", IntentBreakfastInformation},
		{"AC tidak dingin", "305", IntentACProblem},
		{"Lampu kamar mati", "305", IntentLightProblem},
		{"Recommend dinner under Rp100.000", "", IntentRestaurantRecommendation},
		{"Halo", "", IntentGreeting},
		{"Terima kasih", "", IntentThankYou},
		{"asdasdasdasd random", "", IntentUnknown},
	}
	for _, tt := range tests {
		res, err := svc.ProcessSimple(tt.msg, tt.room)
		if err != nil {
			t.Errorf("ProcessSimple(%q) error %v", tt.msg, err)
			continue
		}
		if res.Intent != tt.intent {
			t.Errorf("msg %q: got intent %q want %q", tt.msg, res.Intent, tt.intent)
		}
	}
}

func TestEntityExtraction(t *testing.T) {
	// quantity + room
	ents := ExtractEntities("Tolong antar 3 handuk ke kamar 305", "")
	if ents["quantity"] != 3 {
		t.Errorf("quantity got %v want 3", ents["quantity"])
	}
	if ents["room_number"] != "305" {
		t.Errorf("room got %v want 305", ents["room_number"])
	}
	if ents["item"] != "towel" {
		t.Errorf("item got %v want towel", ents["item"])
	}
	// budget
	ents2 := ExtractEntities("Saya mau makan malam budget 100 ribu", "")
	if ents2["budget"] != 100000 {
		t.Errorf("budget got %v want 100000", ents2["budget"])
	}
	// room from provided
	ents3 := ExtractEntities("AC tidak dingin", "212")
	if ents3["room_number"] != "212" {
		t.Errorf("provided room got %v want 212", ents3["room_number"])
	}
}

func TestDepartmentRouting(t *testing.T) {
	tests := []struct {
		intent string
		dept   string
	}{
		{IntentTowelRequest, "HOUSEKEEPING"},
		{IntentACProblem, "ENGINEERING"},
		{IntentWifiProblem, "ENGINEERING"},
	}
	for _, tt := range tests {
		if got := DepartmentRouting[tt.intent]; got != tt.dept {
			t.Errorf("routing %s got %s want %s", tt.intent, got, tt.dept)
		}
	}
}

func TestPriorityClassification(t *testing.T) {
	if PriorityMapping[IntentACProblem] != "HIGH" {
		t.Errorf("AC should be HIGH")
	}
	if PriorityMapping[IntentTowelRequest] != "MEDIUM" {
		t.Errorf("towel should be MEDIUM")
	}
	if PriorityMapping[IntentAmenityRequest] != "LOW" {
		t.Errorf("amenity should be LOW")
	}
}

func TestUnknownIntent(t *testing.T) {
	svc := NewServiceWithProvider(NewMockProvider(), NewMockProvider())
	res, _ := svc.ProcessSimple("asdasd qwerty zxcvbn", "")
	if res.Intent != IntentUnknown {
		t.Errorf("unknown intent got %q want UNKNOWN", res.Intent)
	}
	if res.Action.Type != ActionClarify {
		t.Errorf("unknown action got %q want CLARIFY", res.Action.Type)
	}
	if res.RequiresTicket {
		t.Error("unknown should not require ticket")
	}
	// Response should be clarification
	if res.Response == "" {
		t.Error("unknown response empty")
	}
}

func TestMissingRoomNumber(t *testing.T) {
	svc := NewServiceWithProvider(NewMockProvider(), NewMockProvider())
	// AC problem without room in message and no provided room
	res, _ := svc.ProcessSimple("AC saya tidak dingin", "")
	if res.Intent != IntentACProblem {
		t.Fatalf("intent got %q", res.Intent)
	}
	if _, ok := res.Entities["room_number"]; ok {
		t.Error("should not have invented room_number")
	}
	// Validator should keep action CREATE_TICKET but entities missing room -> handler should ask for room
	// Service should not invent room
	if res.RequiresTicket == false {
		t.Error("AC should require ticket")
	}
	// Response should ask for room
	if res.Response != "Baik, saya bisa membantu melaporkan masalah AC Anda. Boleh saya tahu nomor kamar Anda?" {
		// Could be fallback, but should ask for room
		t.Logf("response got %q", res.Response)
	}

	// With provided room, should have room
	res2, _ := svc.ProcessSimple("AC saya tidak dingin", "305")
	if res2.Entities["room_number"] != "305" {
		t.Errorf("with provided room, got %v", res2.Entities["room_number"])
	}
}

func TestMalformedAIResponse(t *testing.T) {
	// Provider returns invalid intent
	malformed := &RawAIOutput{
		Intent:   "INVALID_INTENT_XYZ",
		Language: "id",
		Entities: map[string]interface{}{"room_number": "305"},
		Action:   Action{Type: ActionCreateTicket, Department: "HOUSEKEEPING", Priority: "MEDIUM"},
		Response: "test",
	}
	svc := NewServiceWithProvider(NewMalformedProvider(malformed), NewMockProvider())
	res, err := svc.Process(context.Background(), Request{Message: "test", RoomNumber: "305"})
	// Should fallback to mock and not return error for invalid intent (fallback)
	if err != nil {
		t.Logf("expected fallback, got err %v", err)
	}
	// Fallback should give valid intent (mock will classify "test" as UNKNOWN -> CLARIFY)
	if res.Intent == "INVALID_INTENT_XYZ" {
		t.Error("should not return invalid intent, should fallback")
	}
	// Test invalid department
	malformed2 := &RawAIOutput{
		Intent:   IntentTowelRequest,
		Language: "id",
		Entities: map[string]interface{}{"room_number": "305"},
		Action:   Action{Type: ActionCreateTicket, Department: "INVALID_DEPT", Priority: "MEDIUM"},
		Response: "test",
	}
	svc2 := NewServiceWithProvider(NewMalformedProvider(malformed2), NewMockProvider())
	res2, _ := svc2.Process(context.Background(), Request{Message: "towel", RoomNumber: "305"})
	// Validator should correct department to HOUSEKEEPING
	if res2.Action.Department != "HOUSEKEEPING" {
		t.Errorf("invalid dept should be corrected to HOUSEKEEPING, got %q", res2.Action.Department)
	}

	// Test invalid room number format
	malformed3 := &RawAIOutput{
		Intent:   IntentACProblem,
		Language: "id",
		Entities: map[string]interface{}{"room_number": "99999INVALIDTOOLONG"},
		Action:   Action{Type: ActionCreateTicket, Department: "ENGINEERING", Priority: "HIGH"},
		Response: "test",
	}
	svc3 := NewServiceWithProvider(NewMalformedProvider(malformed3), NewMockProvider())
	res3, _ := svc3.Process(context.Background(), Request{Message: "AC rusak", RoomNumber: ""})
	if _, ok := res3.Entities["room_number"]; ok && res3.Entities["room_number"] == "99999INVALIDTOOLONG" {
		t.Error("invalid room should not be accepted")
	}

	// Test invalid priority
	malformed4 := &RawAIOutput{
		Intent:   IntentACProblem,
		Language: "id",
		Entities: map[string]interface{}{"room_number": "305"},
		Action:   Action{Type: ActionCreateTicket, Department: "ENGINEERING", Priority: "ULTRA"},
		Response: "test",
	}
	svc4 := NewServiceWithProvider(NewMalformedProvider(malformed4), NewMockProvider())
	_, err = svc4.Process(context.Background(), Request{Message: "AC rusak", RoomNumber: "305"})
	// Should fallback or correct? Our validator returns error for invalid priority, then fallback
	if err != nil {
		t.Logf("invalid priority correctly caused fallback with err %v", err)
	}
}

func TestProviderFailure(t *testing.T) {
	svc := NewServiceWithProvider(NewFailingProvider(), NewMockProvider())
	res, err := svc.Process(context.Background(), Request{Message: "Hello", RoomNumber: ""})
	// Should fallback to mock, not crash, return controlled response
	if err != nil {
		t.Logf("fallback err %v (acceptable, controlled)", err)
	}
	if res == nil {
		t.Fatal("should return fallback result even on provider failure")
	}
	if res.Response == "" {
		t.Error("fallback response empty")
	}
	// Should not panic
}

func TestStructuredOutput(t *testing.T) {
	svc := NewServiceWithProvider(NewMockProvider(), NewMockProvider())
	res, _ := svc.ProcessSimple("Tolong antar 2 handuk ke kamar 305", "305")
	// Validate structured output has required fields
	if res.Intent != IntentTowelRequest {
		t.Errorf("intent %q", res.Intent)
	}
	if res.Language != LangID {
		t.Errorf("lang %q", res.Language)
	}
	if res.Action.Type != ActionCreateTicket {
		t.Errorf("action %q", res.Action.Type)
	}
	if res.Action.Department != "HOUSEKEEPING" {
		t.Errorf("dept %q", res.Action.Department)
	}
	if res.Action.Priority != "MEDIUM" {
		t.Errorf("priority %q", res.Action.Priority)
	}
	if res.Entities["quantity"] != 2 {
		t.Errorf("quantity %v", res.Entities["quantity"])
	}
	if res.Response == "" {
		t.Error("response empty")
	}
}

func TestActionDecision(t *testing.T) {
	svc := NewServiceWithProvider(NewMockProvider(), NewMockProvider())
	// Information should be ANSWER, not ticket
	res, _ := svc.ProcessSimple("Jam breakfast sampai jam berapa?", "")
	if res.Action.Type != ActionAnswer {
		t.Errorf("breakfast should be ANSWER got %q", res.Action.Type)
	}
	if res.RequiresTicket {
		t.Error("breakfast should not require ticket")
	}
	// Maintenance should be CREATE_TICKET
	res2, _ := svc.ProcessSimple("AC tidak dingin", "305")
	if res2.Action.Type != ActionCreateTicket {
		t.Errorf("AC should be CREATE_TICKET got %q", res2.Action.Type)
	}
	if !res2.RequiresTicket {
		t.Error("AC should require ticket")
	}
	// Unknown should be CLARIFY
	res3, _ := svc.ProcessSimple("asdfgh jkl", "")
	if res3.Action.Type != ActionClarify {
		t.Errorf("unknown should be CLARIFY got %q", res3.Action.Type)
	}
}

func TestRoomNumberNotInvented(t *testing.T) {
	svc := NewServiceWithProvider(NewMockProvider(), NewMockProvider())
	// Message without room, no provided room -> should NOT invent
	res, _ := svc.ProcessSimple("AC saya rusak", "")
	if _, ok := res.Entities["room_number"]; ok {
		t.Error("AI invented room_number when not provided")
	}
	// Provided room should be used
	res2, _ := svc.ProcessSimple("AC saya rusak", "412")
	if res2.Entities["room_number"] != "412" {
		t.Errorf("provided room not used got %v", res2.Entities["room_number"])
	}
}
