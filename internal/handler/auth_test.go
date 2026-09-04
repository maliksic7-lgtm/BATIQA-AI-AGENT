package handler

import (
	"testing"
	"time"
)

func TestRateLimitedBlocksAfterMaxFails(t *testing.T) {
	const key = "test-block@example.com"
	defer clearFails(key)

	for i := 0; i < rateMaxFails-1; i++ {
		recordFail(key)
	}
	if limited, _ := rateLimited(key); limited {
		t.Fatalf("expected not limited before reaching max fails")
	}
	recordFail(key)
	limited, retry := rateLimited(key)
	if !limited {
		t.Fatalf("expected limited at max fails")
	}
	if retry <= 0 || retry > rateWindow {
		t.Fatalf("unexpected retry duration: %v", retry)
	}
}

func TestClearFailsResetsLimit(t *testing.T) {
	const key = "test-clear@example.com"
	defer clearFails(key)

	for i := 0; i < rateMaxFails; i++ {
		recordFail(key)
	}
	if limited, _ := rateLimited(key); !limited {
		t.Fatalf("expected limited before clear")
	}
	clearFails(key)
	if limited, _ := rateLimited(key); limited {
		t.Fatalf("expected not limited after clear")
	}
}

func TestRateLimitedIsKeyIsolated(t *testing.T) {
	const a = "test-iso-a@example.com"
	const b = "test-iso-b@example.com"
	defer clearFails(a)
	defer clearFails(b)

	for i := 0; i < rateMaxFails; i++ {
		recordFail(a)
	}
	if limited, _ := rateLimited(a); !limited {
		t.Fatalf("expected key a limited")
	}
	if limited, _ := rateLimited(b); limited {
		t.Fatalf("expected key b unaffected by key a")
	}
}

func TestExpiredFailsDontCount(t *testing.T) {
	const key = "test-expired@example.com"
	defer clearFails(key)

	rateMu.Lock()
	loginFails[key] = append(loginFails[key], time.Now().Add(-rateWindow-time.Minute))
	rateMu.Unlock()

	// One old entry should not block the key.
	if limited, _ := rateLimited(key); limited {
		t.Fatalf("expected stale fail to be pruned and not block")
	}
}
