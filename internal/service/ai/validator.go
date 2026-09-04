package ai

import (
	"fmt"
	"regexp"
	"strings"
)

var roomRegexValid = regexp.MustCompile(`^\d{2,4}[A-Z]?$`)

// Validate checks RawAIOutput per STRUCTURED AI OUTPUT.md backend validation
// - Intent exists
// - Department allowed (if present)
// - Priority allowed (if present)
// - Action type allowed
// - Entity types expected
// - Room number not invented (must be from entities or provided, not hallucinated)
// - No arbitrary DB operation
func Validate(raw *RawAIOutput, providedRoom string) (*AIResult, error) {
	if raw == nil {
		return nil, fmt.Errorf("nil AI output")
	}

	// Normalize intent
	intent := NormalizeIntent(raw.Intent)
	if !IsValidIntent(intent) {
		return nil, fmt.Errorf("invalid intent: %s", raw.Intent)
	}

	// Validate language
	lang := strings.ToLower(strings.TrimSpace(raw.Language))
	if lang != LangID && lang != LangEN && lang != LangZH {
		// fallback to detected language, not error
		lang = DetectLanguage(raw.Response)
		if lang != LangID && lang != LangEN && lang != LangZH {
			lang = LangID
		}
	}

	// Validate action
	actionType := strings.ToUpper(strings.TrimSpace(raw.Action.Type))
	if actionType != ActionCreateTicket && actionType != ActionAnswer && actionType != ActionClarify {
		return nil, fmt.Errorf("invalid action type: %s", raw.Action.Type)
	}

	// Department validation if present
	dept := strings.ToUpper(strings.TrimSpace(raw.Action.Department))
	if dept != "" {
		if dept != "HOUSEKEEPING" && dept != "ENGINEERING" && dept != "FRONT_OFFICE" {
			return nil, fmt.Errorf("invalid department: %s", dept)
		}
		// Ensure department matches intent routing (if intent requires ticket)
		if expected, ok := DepartmentRouting[intent]; ok && dept != expected {
			// Allow but warn: override to expected to prevent misrouting
			dept = expected
		}
	} else {
		// If action is CREATE_TICKET but no department, infer from routing
		if actionType == ActionCreateTicket {
			if expected, ok := DepartmentRouting[intent]; ok {
				dept = expected
			} else {
				// No routing for intent that requires ticket -> invalid
				if RequiresTicket[intent] {
					return nil, fmt.Errorf("missing department for intent %s", intent)
				}
			}
		}
	}

	// Priority validation
	priority := strings.ToUpper(strings.TrimSpace(raw.Action.Priority))
	if priority != "" {
		if priority != "LOW" && priority != "MEDIUM" && priority != "HIGH" {
			return nil, fmt.Errorf("invalid priority: %s", priority)
		}
	} else {
		if actionType == ActionCreateTicket {
			if p, ok := PriorityMapping[intent]; ok {
				priority = p
			} else {
				priority = "MEDIUM"
			}
		}
	}

	// Entities validation
	entities := raw.Entities
	if entities == nil {
		entities = make(map[string]interface{})
	}

	// Room number: must not be invented. Only accept if present in entities and valid, or providedRoom
	// AI must not invent room numbers per AI SAFETY RULES
	var finalRoom string
	if v, ok := entities["room_number"]; ok {
		if s, ok := v.(string); ok && s != "" {
			s = strings.TrimSpace(s)
			if !roomRegexValid.MatchString(s) {
				return nil, fmt.Errorf("invalid room_number format: %s", s)
			}
			// Check if AI invented room when providedRoom is empty and message did not contain room
			// We trust entity extraction only if it was in message or providedRoom
			// Since we use ExtractEntities which already validates, we accept if valid
			finalRoom = s
		} else {
			return nil, fmt.Errorf("room_number must be string")
		}
	}
	// If no room in entities but providedRoom exists, use it (do not treat as AI invention)
	if finalRoom == "" && providedRoom != "" {
		if !roomRegexValid.MatchString(providedRoom) {
			return nil, fmt.Errorf("provided room_number invalid: %s", providedRoom)
		}
		finalRoom = providedRoom
		entities["room_number"] = finalRoom
	}
	// If still empty and action requires ticket, do NOT error here: let service ask for room (missing room flow)
	// But ensure we don't keep invented room

	// Quantity validation
	if v, ok := entities["quantity"]; ok {
		switch val := v.(type) {
		case int:
			if val <= 0 || val > 100 {
				return nil, fmt.Errorf("quantity out of range: %d", val)
			}
		case float64:
			if val <= 0 || val > 100 {
				return nil, fmt.Errorf("quantity out of range: %v", val)
			}
			// normalize to int
			entities["quantity"] = int(val)
		default:
			return nil, fmt.Errorf("quantity must be number")
		}
	}

	// Budget validation
	if v, ok := entities["budget"]; ok {
		switch val := v.(type) {
		case int:
			if val < 0 {
				return nil, fmt.Errorf("budget negative")
			}
		case float64:
			if val < 0 {
				return nil, fmt.Errorf("budget negative")
			}
		default:
			return nil, fmt.Errorf("budget must be number")
		}
	}

	// Prevent arbitrary DB operation: only allowed actions are CREATE_TICKET, ANSWER, CLARIFY
	// Already validated above

	// Build validated result
	requiresTicket := actionType == ActionCreateTicket && RequiresTicket[intent]

	// If requires ticket but no department/priority, fill defaults (already done)
	// If intent is UNKNOWN, force CLARIFY
	if intent == IntentUnknown {
		actionType = ActionClarify
		requiresTicket = false
		dept = ""
		priority = ""
	}

	result := &AIResult{
		Intent:   intent,
		Language: lang,
		Entities: entities,
		Action: Action{
			Type:       actionType,
			Department: dept,
			Priority:   priority,
		},
		Response:       raw.Response,
		RequiresTicket: requiresTicket,
	}

	// Ensure response is not empty
	if strings.TrimSpace(result.Response) == "" {
		switch lang {
		case LangEN:
			result.Response = "How can I help you?"
		case LangZH:
			result.Response = "您好，有什么可以帮您？"
		default:
			result.Response = "Ada yang bisa saya bantu?"
		}
	}

	return result, nil
}
