package repository

import "testing"

// TestInfographicIntentSetsDisjoint ensures a ticket category is never counted
// as both a "complaint" and a "borrowed" item, so dashboard slices don't double
// count the same ticket.
func TestInfographicIntentSetsDisjoint(t *testing.T) {
	for k := range complaintIntents {
		if borrowedIntents[k] {
			t.Errorf("intent %q is in both complaintIntents and borrowedIntents", k)
		}
	}
}

// TestInfographicIntentSetsNonEmpty guards against an accidental empty set that
// would silently produce an empty dashboard slice.
func TestInfographicIntentSetsNonEmpty(t *testing.T) {
	if len(complaintIntents) == 0 {
		t.Errorf("complaintIntents is empty")
	}
	if len(borrowedIntents) == 0 {
		t.Errorf("borrowedIntents is empty")
	}
}

// TestCountOrderedKeywordsCountsMentions checks the pure F&B matching rules:
// a guest message mentioning a menu item is counted once for that item, keys
// are case-insensitive, and unrelated messages add no count.
func TestCountOrderedKeywordsCountsMentions(t *testing.T) {
	msgs := []string{
		"can I order gulai patin and some iced tea please",
		"UMA sate padang untuk meja 3",
		"beri aku 2 omlet dan kopi",
	}
	counts := countOrderedKeywords(msgs)

	if counts["Gulai Ikan Patin"] != 1 {
		t.Errorf("expected Gulai Ikan Patin x1, got %v", counts["Gulai Ikan Patin"])
	}
	if counts["Sate Padang"] != 1 {
		t.Errorf("expected Sate Padang x1, got %v", counts["Sate Padang"])
	}
	if counts["Omelet"] != 1 {
		t.Errorf("expected Omelet x1, got %v", counts["Omelet"])
	}
	if counts["Kopi"] != 1 {
		t.Errorf("expected Kopi x1, got %v", counts["Kopi"])
	}
	// "teh" keyword appears in "iced tea"; "patin" already counted once for that
	// message but "Teh" should also register because the keywords include "tea".
	if counts["Teh"] != 1 {
		t.Errorf("expected Teh x1 (iced tea), got %v", counts["Teh"])
	}
}
