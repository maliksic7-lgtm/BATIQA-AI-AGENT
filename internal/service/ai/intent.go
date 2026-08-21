package ai

import "strings"

// Intent constants per docs/INTENT CATEGORIES.md
const (
	// Hotel Information (no ticket)
	IntentHotelInformation      = "HOTEL_INFORMATION"
	IntentFacilityInformation   = "FACILITY_INFORMATION"
	IntentBreakfastInformation  = "BREAKFAST_INFORMATION"
	IntentRestaurantInformation = "RESTAURANT_INFORMATION"
	IntentWifiInformation       = "WIFI_INFORMATION"
	IntentCheckinInformation    = "CHECKIN_INFORMATION"
	IntentCheckoutInformation   = "CHECKOUT_INFORMATION"
	IntentRoomInformation       = "ROOM_INFORMATION"
	IntentHotelPolicy           = "HOTEL_POLICY"

	// Housekeeping -> HOUSEKEEPING
	IntentHousekeepingRequest = "HOUSEKEEPING_REQUEST"
	IntentTowelRequest        = "TOWEL_REQUEST"
	IntentAmenityRequest      = "AMENITY_REQUEST"
	IntentRoomCleaningRequest = "ROOM_CLEANING_REQUEST"

	// Engineering -> ENGINEERING
	IntentACProblem             = "AC_PROBLEM"
	IntentTVProblem             = "TV_PROBLEM"
	IntentWifiProblem           = "WIFI_PROBLEM"
	IntentLightProblem          = "LIGHT_PROBLEM"
	IntentShowerProblem         = "SHOWER_PROBLEM"
	IntentPlumbingProblem       = "PLUMBING_PROBLEM"
	IntentRoomEquipmentProblem  = "ROOM_EQUIPMENT_PROBLEM"
	IntentGeneralMaintenance    = "GENERAL_MAINTENANCE"

	// Recommendations (no ticket)
	IntentRestaurantRecommendation = "RESTAURANT_RECOMMENDATION"
	IntentCafeRecommendation       = "CAFE_RECOMMENDATION"
	IntentTourismRecommendation    = "TOURISM_RECOMMENDATION"
	IntentShoppingRecommendation   = "SHOPPING_RECOMMENDATION"
	IntentATMRequest               = "ATM_REQUEST"
	IntentTransportationRequest    = "TRANSPORTATION_REQUEST"

	// General
	IntentGeneralQuestion = "GENERAL_QUESTION"
	IntentGreeting        = "GREETING"
	IntentThankYou        = "THANK_YOU"
	IntentUnknown         = "UNKNOWN"
)

// ValidIntents is set for validation layer per STRUCTURED AI OUTPUT.md
var ValidIntents = map[string]bool{
	IntentHotelInformation: true, IntentFacilityInformation: true, IntentBreakfastInformation: true,
	IntentRestaurantInformation: true, IntentWifiInformation: true, IntentCheckinInformation: true,
	IntentCheckoutInformation: true, IntentRoomInformation: true, IntentHotelPolicy: true,
	IntentHousekeepingRequest: true, IntentTowelRequest: true, IntentAmenityRequest: true, IntentRoomCleaningRequest: true,
	IntentACProblem: true, IntentTVProblem: true, IntentWifiProblem: true, IntentLightProblem: true,
	IntentShowerProblem: true, IntentPlumbingProblem: true, IntentRoomEquipmentProblem: true, IntentGeneralMaintenance: true,
	IntentRestaurantRecommendation: true, IntentCafeRecommendation: true, IntentTourismRecommendation: true,
	IntentShoppingRecommendation: true, IntentATMRequest: true, IntentTransportationRequest: true,
	IntentGeneralQuestion: true, IntentGreeting: true, IntentThankYou: true, IntentUnknown: true,
}

// DepartmentRouting per TICKET ROUTING.md
var DepartmentRouting = map[string]string{
	IntentHousekeepingRequest: "HOUSEKEEPING",
	IntentTowelRequest:        "HOUSEKEEPING",
	IntentAmenityRequest:      "HOUSEKEEPING",
	IntentRoomCleaningRequest: "HOUSEKEEPING",
	IntentACProblem:           "ENGINEERING",
	IntentTVProblem:           "ENGINEERING",
	IntentWifiProblem:         "ENGINEERING",
	IntentLightProblem:        "ENGINEERING",
	IntentShowerProblem:       "ENGINEERING",
	IntentPlumbingProblem:     "ENGINEERING",
	IntentRoomEquipmentProblem: "ENGINEERING",
	IntentGeneralMaintenance:  "ENGINEERING",
}

// PriorityMapping per PRIORITY CLASSIFICATION.md
var PriorityMapping = map[string]string{
	IntentTowelRequest:        "MEDIUM",
	IntentHousekeepingRequest: "MEDIUM",
	IntentAmenityRequest:      "LOW",
	IntentRoomCleaningRequest: "MEDIUM",
	IntentACProblem:           "HIGH",
	IntentTVProblem:           "MEDIUM",
	IntentWifiProblem:         "MEDIUM",
	IntentLightProblem:        "MEDIUM",
	IntentShowerProblem:       "HIGH",
	IntentPlumbingProblem:     "HIGH",
	IntentRoomEquipmentProblem: "MEDIUM",
	IntentGeneralMaintenance:  "MEDIUM",
}

// RequiresTicket intents
var RequiresTicket = map[string]bool{
	IntentHousekeepingRequest: true, IntentTowelRequest: true, IntentAmenityRequest: true, IntentRoomCleaningRequest: true,
	IntentACProblem: true, IntentTVProblem: true, IntentWifiProblem: true, IntentLightProblem: true,
	IntentShowerProblem: true, IntentPlumbingProblem: true, IntentRoomEquipmentProblem: true, IntentGeneralMaintenance: true,
}

// IsValidIntent checks if intent is known
func IsValidIntent(s string) bool {
	return ValidIntents[s]
}

// NormalizeIntent uppercases and trims
func NormalizeIntent(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
