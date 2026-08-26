package ai

import (
	"regexp"
	"strconv"
	"strings"
)

// Entity extraction patterns per ENTITY EXTRACTION.md

var (
	roomRegex     = regexp.MustCompile(`(?i)\b(kamar\s*)?(\d{2,4}[A-Z]?)\b`)
	quantityRegex = regexp.MustCompile(`(?i)\b(\d+)\s*(handuk|towel|pillow|bantal|selimut|blanket|botol|pcs|buah)\b`)
	budgetRegex   = regexp.MustCompile(`(?i)\b(\d{1,3}(?:[.,]\d{3})*)\s*(ribu|rb|k|rupiah|rp)?\b`)
)

// ExtractEntities extracts entities from message and existing roomNumber context.
// Room number is NOT invented: only extracted if present in message, else uses provided roomNumber if valid.
func ExtractEntities(message, providedRoom string) map[string]interface{} {
	entities := make(map[string]interface{})
	lower := strings.ToLower(message)

	// room_number: try to find in message first
	if m := roomRegex.FindStringSubmatch(message); m != nil {
		// m[2] is number part
		candidate := strings.TrimSpace(m[2])
		// Validate room number format: 2-4 digits optionally letter, not part of price
		if isValidRoomNumber(candidate) && !isPriceContext(lower, candidate) {
			entities["room_number"] = candidate
		}
	} else if providedRoom != "" && isValidRoomNumber(providedRoom) {
		entities["room_number"] = providedRoom
	}

	// quantity + item
	if m := quantityRegex.FindStringSubmatch(message); m != nil {
		if n, err := strconv.Atoi(m[1]); err == nil {
			entities["quantity"] = n
		}
		item := strings.ToLower(m[2])
		// normalize item
		switch item {
		case "handuk", "towel":
			entities["item"] = "towel"
		case "bantal", "pillow":
			entities["item"] = "pillow"
		default:
			entities["item"] = item
		}
	} else {
		// check for towel without quantity -> default 1
		if strings.Contains(lower, "handuk") || strings.Contains(lower, "towel") {
			entities["item"] = "towel"
			if _, ok := entities["quantity"]; !ok {
				entities["quantity"] = 1
			}
		}
	}

	// problem: capture after keywords
	problemKeywords := []string{"ac tidak dingin", "ac rusak", "ac bermasalah", "tidak dingin", "tv bermasalah", "tv tidak menyala", "wifi tidak", "lampu rusak", "lampu mati", "shower bermasalah", "air bocor", "bocor"}
	for _, kw := range problemKeywords {
		if strings.Contains(lower, kw) {
			entities["problem"] = kw
			break
		}
	}
	if _, ok := entities["problem"]; !ok {
		// fallback: if engineering intent, use message snippet
		if strings.Contains(lower, "rusak") || strings.Contains(lower, "bermasalah") || strings.Contains(lower, "tidak") {
			entities["problem"] = strings.TrimSpace(message)
		}
	}

	// budget: look for numbers near budget keywords
	if strings.Contains(lower, "budget") || strings.Contains(lower, "rp") || strings.Contains(lower, "ribu") || strings.Contains(lower, "rupiah") {
		if m := budgetRegex.FindStringSubmatch(message); m != nil {
			numStr := strings.ReplaceAll(m[1], ".", "")
			numStr = strings.ReplaceAll(numStr, ",", "")
			if n, err := strconv.Atoi(numStr); err == nil {
				// handle ribu
				if m[2] != "" {
					l := strings.ToLower(m[2])
					if l == "ribu" || l == "rb" || l == "k" {
						n *= 1000
					}
				}
				entities["budget"] = n
			}
		}
	}

	// category for recommendations
	catKeywords := map[string]string{
		"restaurant": "restaurant", "restoran": "restaurant", "makan": "restaurant",
		"cafe": "cafe", "kafe": "cafe", "kopi": "cafe",
		"mall": "shopping", "belanja": "shopping", "shopping": "shopping",
		"wisata": "tourism", "tourism": "tourism", "pantai": "tourism", "candi": "tourism",
		"atm":       "atm",
		"transport": "transportation", "taksi": "transportation", "ojek": "transportation",
	}
	for k, v := range catKeywords {
		if strings.Contains(lower, k) {
			entities["category"] = v
			break
		}
	}

	// preference (free text after "prefer" etc) - simple
	if strings.Contains(lower, "prefer") {
		entities["preference"] = message
	}

	return entities
}

func isValidRoomNumber(s string) bool {
	// Must be 2-4 digits optionally letter, and be plausible hotel room
	if len(s) < 2 || len(s) > 5 {
		return false
	}
	matched, _ := regexp.MatchString(`^\d{2,4}[A-Z]?$`, s)
	return matched
}

func isPriceContext(lower, candidate string) bool {
	// If candidate number appears near "rp", "budget", "ribu", likely price not room
	idx := strings.Index(lower, strings.ToLower(candidate))
	if idx == -1 {
		return false
	}
	contextWindow := 20
	start := idx - contextWindow
	if start < 0 {
		start = 0
	}
	end := idx + len(candidate) + contextWindow
	if end > len(lower) {
		end = len(lower)
	}
	ctx := lower[start:end]
	return strings.Contains(ctx, "rp") || strings.Contains(ctx, "budget") || strings.Contains(ctx, "ribu") || strings.Contains(ctx, "rb") || strings.Contains(ctx, "rupiah")
}
