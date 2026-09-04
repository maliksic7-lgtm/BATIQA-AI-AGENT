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
