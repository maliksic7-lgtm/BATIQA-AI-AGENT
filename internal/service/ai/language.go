package ai

import (
	"strings"
	"unicode"
)

// DetectLanguage detects guest language per LANGUAGE BEHAVIOR.md
// Simple heuristic: if message contains Indonesian keywords, return "id", else "en"
// Preserves guest language unless explicitly asked to translate.
func DetectLanguage(message string) string {
	lower := strings.ToLower(message)
	// Strong Indonesian markers (unambiguous)
	strongID := []string{
		"saya", "kamar", "tolong", "minta", "antar", "berapa", "jam", "sampai",
		"tidak dingin", "rusak", "bocor", "handuk", "bersih", "bersihkan",
		"terima kasih", "sarapan", "bisa", "mau", "butuh",
	}
	for _, m := range strongID {
		if strings.Contains(lower, m) {
			return LangID
		}
	}
	// Weak markers that alone are ambiguous (e.g., "tidak", "ac", "wifi", "hotel")
	// Only count them if combined with other ID context
	weakIDCount := 0
	weakID := []string{"tidak", "air", "lampu", "maaf", "pagi", "malam", "siang", "kapan", "dimana", "di mana"}
	for _, m := range weakID {
		if strings.Contains(lower, m) {
			weakIDCount++
		}
	}
	if weakIDCount >= 2 {
		return LangID
	}
	// Common ID structure
	idShort := []string{" yang ", " dan ", " di ", " ke ", " dari ", " untuk "}
	for _, s := range idShort {
		if strings.Contains(lower, s) {
			return LangID
		}
	}
	// Check English markers first (strong)
	enMarkers := []string{"hello", "hi", "please", "thank you", "breakfast", "what", "when", "where", "how", "need", "want", "room", "my "}
	hasEN := false
	for _, m := range enMarkers {
		if strings.Contains(lower, m) {
			hasEN = true
			break
		}
	}
	if hasEN {
		return LangEN
	}
	// Default: check if message contains only ascii without id markers -> en, else id
	// Count non-ascii
	for _, r := range lower {
		if r > unicode.MaxASCII {
			return LangID
		}
	}
	// Heuristic fallback: if short greeting in english
	if strings.HasPrefix(lower, "hi") || strings.HasPrefix(lower, "hello") {
		return LangEN
	}
	return LangID
}
