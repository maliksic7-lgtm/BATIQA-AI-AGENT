package ai

import (
	"context"
	"fmt"
	"strings"
)

// MockProvider is deterministic rule-based provider for testing and fallback.
// It implements the 7 steps: language -> intent -> entity -> routing -> priority -> action -> response
// without calling external API.
type MockProvider struct{}

func NewMockProvider() *MockProvider { return &MockProvider{} }

func (m *MockProvider) Name() string { return "mock" }

func (m *MockProvider) Generate(ctx context.Context, req Request) (*RawAIOutput, error) {
	// Check context
	select {
	case <-ctx.Done():
		return nil, &ProviderError{Provider: m.Name(), Err: ctx.Err()}
	default:
	}

	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		return &RawAIOutput{
			Intent:   IntentUnknown,
			Language: LangID,
			Entities: map[string]interface{}{},
			Action:   Action{Type: ActionClarify},
			Response: "Silakan tulis pesan Anda.",
		}, nil
	}

	lang := DetectLanguage(msg)
	entities := ExtractEntities(msg, req.RoomNumber)
	intent := classifyIntent(msg, lang)

	// Department & Priority
	dept := DepartmentRouting[intent]
	priority := PriorityMapping[intent]

	// Action decision
	actionType := ActionAnswer
	if RequiresTicket[intent] {
		actionType = ActionCreateTicket
	} else if intent == IntentUnknown || intent == IntentGeneralQuestion {
		// If unknown, ask clarify; but if general question may need info
		if intent == IntentUnknown {
			actionType = ActionClarify
		} else {
			actionType = ActionAnswer
		}
	}
	// For greeting/thank you, answer
	if intent == IntentGreeting || intent == IntentThankYou {
		actionType = ActionAnswer
	}

	// Build response language-aware
	response := buildResponse(intent, lang, entities, req.RoomNumber)

	// If missing room_number for ticket intents, keep action but will be validated later to ask for room
	// Do NOT invent room_number (AI SAFETY RULES)

	out := &RawAIOutput{
		Intent:   intent,
		Language: lang,
		Entities: entities,
		Action: Action{
			Type:       actionType,
			Department: dept,
			Priority:   priority,
		},
		Response: response,
	}
	return out, nil
}

// classifyIntent performs keyword-based intent classification per INTENT CATEGORIES.md
func classifyIntent(msg, lang string) string {
	lower := strings.ToLower(msg)

	// Greeting
	if lower == "hi" || lower == "hello" || lower == "halo" || strings.HasPrefix(lower, "hi ") || strings.HasPrefix(lower, "hello ") || lower == "halo kak" {
		return IntentGreeting
	}
	if strings.Contains(lower, "terima kasih") || strings.Contains(lower, "thanks") || lower == "thank you" {
		return IntentThankYou
	}

	// Breakfast
	if strings.Contains(lower, "breakfast") || strings.Contains(lower, "sarapan") {
		if strings.Contains(lower, "jam") || strings.Contains(lower, "schedule") || strings.Contains(lower, "sampai") || strings.Contains(lower, "buka") {
			return IntentBreakfastInformation
		}
		return IntentBreakfastInformation
	}
	// WiFi
	if strings.Contains(lower, "wifi") || strings.Contains(lower, "wi-fi") {
		if strings.Contains(lower, "tidak") || strings.Contains(lower, "bermasalah") || strings.Contains(lower, "rusak") || strings.Contains(lower, "problem") {
			return IntentWifiProblem
		}
		return IntentWifiInformation
	}
	// Checkin/out
	if strings.Contains(lower, "check in") || strings.Contains(lower, "check-in") || strings.Contains(lower, "checkin") {
		return IntentCheckinInformation
	}
	if strings.Contains(lower, "check out") || strings.Contains(lower, "check-out") || strings.Contains(lower, "checkout") {
		return IntentCheckoutInformation
	}
	// Pool/Gym
	if strings.Contains(lower, "pool") || strings.Contains(lower, "kolam") {
		return IntentFacilityInformation
	}
	if strings.Contains(lower, "gym") || strings.Contains(lower, "fitness") {
		return IntentFacilityInformation
	}
	// Restaurant info vs recommendation
	if strings.Contains(lower, "restaurant") || strings.Contains(lower, "restoran") {
		if strings.Contains(lower, "rekomendasi") || strings.Contains(lower, "recommend") || strings.Contains(lower, "budget") || strings.Contains(lower, "makan malam") {
			return IntentRestaurantRecommendation
		}
		return IntentRestaurantInformation
	}
	// Recommendations
	if strings.Contains(lower, "rekomendasi") || strings.Contains(lower, "recommend") {
		if strings.Contains(lower, "cafe") || strings.Contains(lower, "kafe") || strings.Contains(lower, "kopi") {
			return IntentCafeRecommendation
		}
		if strings.Contains(lower, "wisata") || strings.Contains(lower, "tourism") || strings.Contains(lower, "pantai") {
			return IntentTourismRecommendation
		}
		if strings.Contains(lower, "mall") || strings.Contains(lower, "belanja") || strings.Contains(lower, "shopping") {
			return IntentShoppingRecommendation
		}
		if strings.Contains(lower, "makan") || strings.Contains(lower, "dinner") {
			return IntentRestaurantRecommendation
		}
		return IntentGeneralQuestion
	}
	if strings.Contains(lower, "atm") {
		return IntentATMRequest
	}
	if strings.Contains(lower, "transport") || strings.Contains(lower, "taksi") || strings.Contains(lower, "ojek") || strings.Contains(lower, "airport") {
		return IntentTransportationRequest
	}

	// Housekeeping
	if strings.Contains(lower, "handuk") || strings.Contains(lower, "towel") {
		return IntentTowelRequest
	}
	if strings.Contains(lower, "bantal") || strings.Contains(lower, "pillow") || strings.Contains(lower, "selimut") || strings.Contains(lower, "amenity") || strings.Contains(lower, "amenities") {
		return IntentAmenityRequest
	}
	if strings.Contains(lower, "bersih") || strings.Contains(lower, "cleaning") || strings.Contains(lower, "bersihkan") || strings.Contains(lower, "housekeeping") {
		return IntentRoomCleaningRequest
	}
	if strings.Contains(lower, "kebersihan") {
		return IntentHousekeepingRequest
	}

	// Engineering (ID + EN)
	if strings.Contains(lower, "ac") && (strings.Contains(lower, "tidak dingin") || strings.Contains(lower, "tidak") || strings.Contains(lower, "rusak") || strings.Contains(lower, "panas") || strings.Contains(lower, "bermasalah") || strings.Contains(lower, "not working") || strings.Contains(lower, "not cold") || strings.Contains(lower, "not cooling") || strings.Contains(lower, "broken") || strings.Contains(lower, "problem")) {
		return IntentACProblem
	}
	if strings.Contains(lower, "ac tidak") {
		return IntentACProblem
	}
	if strings.Contains(lower, "tv") && (strings.Contains(lower, "rusak") || strings.Contains(lower, "mati") || strings.Contains(lower, "tidak") || strings.Contains(lower, "bermasalah") || strings.Contains(lower, "not working") || strings.Contains(lower, "broken") || strings.Contains(lower, "problem")) {
		return IntentTVProblem
	}
	if (strings.Contains(lower, "lampu") || strings.Contains(lower, "light")) && (strings.Contains(lower, "rusak") || strings.Contains(lower, "mati") || strings.Contains(lower, "bermasalah") || strings.Contains(lower, "tidak nyala") || strings.Contains(lower, "not working") || strings.Contains(lower, "broken")) {
		return IntentLightProblem
	}
	if strings.Contains(lower, "shower") && (strings.Contains(lower, "rusak") || strings.Contains(lower, "bermasalah") || strings.Contains(lower, "tidak") || strings.Contains(lower, "bocor") || strings.Contains(lower, "not working") || strings.Contains(lower, "broken") || strings.Contains(lower, "leak")) {
		return IntentShowerProblem
	}
	if strings.Contains(lower, "bocor") || strings.Contains(lower, "plumbing") || strings.Contains(lower, "leak") || (strings.Contains(lower, "air") && strings.Contains(lower, "tidak")) {
		return IntentPlumbingProblem
	}
	if strings.Contains(lower, "kamar") && (strings.Contains(lower, "rusak") || strings.Contains(lower, "bermasalah")) {
		return IntentRoomEquipmentProblem
	}
	if strings.Contains(lower, "maintenance") || strings.Contains(lower, "perbaikan") {
		return IntentGeneralMaintenance
	}

	// Hotel info general
	if strings.Contains(lower, "fasilitas") || strings.Contains(lower, "facility") || strings.Contains(lower, "hotel") && strings.Contains(lower, "apa") {
		return IntentFacilityInformation
	}
	if strings.Contains(lower, "policy") || strings.Contains(lower, "kebijakan") || strings.Contains(lower, "peraturan") {
		return IntentHotelPolicy
	}
	if strings.Contains(lower, "kamar") && (strings.Contains(lower, "fasilitas") || strings.Contains(lower, "apa")) {
		return IntentRoomInformation
	}

	// Budget recommendation
	if strings.Contains(lower, "budget") && (strings.Contains(lower, "makan") || strings.Contains(lower, "dinner")) {
		return IntentRestaurantRecommendation
	}

	// Unknown
	return IntentUnknown
}

func buildResponse(intent, lang string, entities map[string]interface{}, roomNumber string) string {
	// Language-aware responses
	isID := lang == LangID

	switch intent {
	case IntentGreeting:
		if isID {
			return "Halo! Selamat datang di BATIQA Hotels. Ada yang bisa saya bantu?"
		}
		return "Hello! Welcome to BATIQA Hotels. How can I help you?"
	case IntentThankYou:
		if isID {
			return "Sama-sama! Senang bisa membantu. Jika butuh bantuan lagi, silakan hubungi kami."
		}
		return "You're welcome! Happy to help. Let me know if you need anything else."
	case IntentBreakfastInformation:
		if isID {
			return "Breakfast tersedia mulai pukul 06:00 sampai 10:00 di restaurant lantai 1."
		}
		return "Breakfast is available from 06:00 to 10:00 at the 1st floor restaurant."
	case IntentWifiInformation:
		if isID {
			return "WiFi hotel: Connect to BATIQA HOTELS network. Password tersedia di kartu kamar atau hubungi Front Office."
		}
		return "Hotel WiFi: Connect to BATIQA HOTELS network. Password is on your room card or contact Front Office."
	case IntentCheckinInformation:
		if isID {
			return "Check-in mulai pukul 14:00. Early check-in tergantung ketersediaan."
		}
		return "Check-in starts at 14:00. Early check-in subject to availability."
	case IntentCheckoutInformation:
		if isID {
			return "Check-out pukul 12:00. Late check-out dapat diminta ke Front Office."
		}
		return "Check-out at 12:00. Late check-out can be requested at Front Office."
	case IntentFacilityInformation:
		if isID {
			return "Fasilitas hotel: Restaurant, Swimming Pool (06:00-20:00), Gym 24 jam, WiFi. Ada yang ingin ditanyakan lebih detail?"
		}
		return "Hotel facilities: Restaurant, Swimming Pool (06:00-20:00), Gym 24h, WiFi. Anything specific you'd like to know?"
	case IntentTowelRequest:
		qty := 1
		if v, ok := entities["quantity"]; ok {
			if n, ok := v.(int); ok {
				qty = n
			}
		}
		room := roomStr(entities, roomNumber)
		if isID {
			return formatID("Baik, saya akan mengirimkan permintaan %d handuk ke Housekeeping untuk kamar %s.", qty, room)
		}
		return formatEN("Sure, I will forward your request for %d towel(s) to Housekeeping for room %s.", qty, room)
	case IntentACProblem:
		room := roomStr(entities, roomNumber)
		if isID {
			if room == "Anda" {
				return "Baik, saya bisa membantu melaporkan masalah AC Anda. Boleh saya tahu nomor kamar Anda?"
			}
			return formatID("Baik, saya akan membantu melaporkan masalah AC Anda di kamar %s ke tim Engineering.", room)
		}
		if room == "your" {
			return "Sure, I can help report your AC issue. Could you please tell me your room number?"
		}
		return formatEN("Sure, I will report your AC issue in room %s to Engineering.", room)
	case IntentWifiProblem:
		if isID {
			return "Maaf atas kendala WiFi. Saya akan laporkan ke Engineering untuk segera diperiksa."
		}
		return "Sorry for the WiFi issue. I will report it to Engineering right away."
	case IntentLightProblem, IntentShowerProblem, IntentPlumbingProblem, IntentTVProblem:
		if isID {
			return "Terima kasih laporannya. Saya akan teruskan ke tim Engineering untuk segera ditangani."
		}
		return "Thank you for reporting. I will forward it to Engineering right away."
	case IntentRoomCleaningRequest:
		room := roomStr(entities, roomNumber)
		if isID {
			return formatID("Baik, permintaan pembersihan kamar %s akan saya teruskan ke Housekeeping.", room)
		}
		return formatEN("Sure, cleaning request for room %s will be forwarded to Housekeeping.", room)
	case IntentRestaurantRecommendation:
		if isID {
			return "Untuk rekomendasi restaurant, saya bisa bantu sesuai budget dan preferensi. Budget Anda berapa?"
		}
		return "For restaurant recommendations, I can help based on your budget and preference. What's your budget?"
	case IntentUnknown:
		if isID {
			return "Saya ingin memastikan saya membantu dengan tepat. Apakah Anda ingin bertanya tentang fasilitas hotel, melakukan service request, atau melaporkan masalah kamar?"
		}
		return "I want to make sure I help correctly. Are you asking about hotel facilities, making a service request, or reporting a room issue?"
	default:
		if isID {
			return "Tentu, saya bisa membantu. Ada yang bisa saya bantu lebih lanjut?"
		}
		return "Certainly, I can help. Anything else you need?"
	}
}

func roomStr(entities map[string]interface{}, fallback string) string {
	if v, ok := entities["room_number"]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	if fallback != "" {
		return fallback
	}
	// For ID vs EN, caller handles
	return ""
}

func formatID(format string, a ...interface{}) string {
	// If room is empty, return ask for room
	for _, v := range a {
		if s, ok := v.(string); ok && s == "" {
			return "Baik, saya bisa membantu. Boleh saya tahu nomor kamar Anda?"
		}
	}
	return sprintf(format, a...)
}

func formatEN(format string, a ...interface{}) string {
	for _, v := range a {
		if s, ok := v.(string); ok && s == "" {
			return "Sure, I can help. Could you please tell me your room number?"
		}
	}
	return sprintf(format, a...)
}

func sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

