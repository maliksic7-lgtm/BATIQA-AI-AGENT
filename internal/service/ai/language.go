package ai

import (
	"strings"
	"unicode"
)

// isCJK reports whether a rune is within the common Chinese/Japanese/Korean
// blocks (Han unified + compatibility). Used to detect Chinese-script input
// without false-positive hits from box-drawing or emoji in the ASCII/ID checks.
func isCJK(r rune) bool {
	// CJK Unified Ideographs, Ext-A, compatibility ideographs
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	if r >= 0x3400 && r <= 0x4DBF {
		return true
	}
	if r >= 0xF900 && r <= 0xFAFF {
		return true
	}
	// Fullwidth forms (ｆｕｌｌｗｉｄｔｈ) plus halfwidth katakana are rare in
	// hotel chat; ignore to keep detection conservative.
	return false
}

// DetectLanguage detects guest language per LANGUAGE BEHAVIOR.md
// Simple heuristic: if message contains Indonesian keywords, return "id", else "en"
// Preserves guest language unless explicitly asked to translate.
func DetectLanguage(message string) string {
	lower := strings.ToLower(message)

	// Chinese script is unambiguous: if the message contains any Han ideograph,
	// treat it as Chinese (zh). Checked first so CJK never falls through to the
	// ASCII heuristic that would otherwise label it "en".
	for _, r := range lower {
		if isCJK(r) {
			return LangZH
		}
	}

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
