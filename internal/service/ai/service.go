package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

// Service is the main AI service abstraction per AI_SPEC.md
// It orchestrates: Language Detection -> Intent -> Entity -> Routing -> Priority -> Action -> Response
// and handles validation + fallback. Provider is swappable.
type Service struct {
	primary  Provider
	fallback Provider
}

// NewService creates AI service. Provider is chosen via AI_PROVIDER env: "gemini", "mock", or auto.
// If GEMINI_API_KEY is set and provider is gemini, it will try Gemini first, fallback to mock on failure.
func NewService() *Service {
	providerName := strings.ToLower(strings.TrimSpace(os.Getenv("AI_PROVIDER")))
	if providerName == "" {
		// Auto: use gemini if key exists, else mock
		if strings.TrimSpace(os.Getenv("GEMINI_API_KEY")) != "" {
			providerName = "gemini"
		} else {
			providerName = "mock"
		}
	}

	var primary Provider
	var fallback Provider = NewMockProvider()

	switch providerName {
	case "gemini":
		primary = NewGeminiProvider()
	case "mock", "rulebased", "rule":
		primary = NewMockProvider()
		fallback = NewMockProvider()
	default:
		primary = NewMockProvider()
	}

	return &Service{primary: primary, fallback: fallback}
}

// NewServiceWithProvider is for testing (inject mock/failing provider)
func NewServiceWithProvider(primary, fallback Provider) *Service {
	if fallback == nil {
		fallback = NewMockProvider()
	}
	return &Service{primary: primary, fallback: fallback}
}

// Process is the main entry: transforms natural language to structured AIResult with validation.
// It never panics; on provider failure returns controlled fallback (UNKNOWN + clarification).
func (s *Service) Process(ctx context.Context, req Request) (result *AIResult, err error) {
	// Recover from panic per ERROR FLOW.md: never crash, return safe fallback
	defer func() {
		if r := recover(); r != nil {
			result = fallbackResult(req, langFromReq(req))
			err = fmt.Errorf("recovered from panic in AI provider %s: %v", s.primary.Name(), r)
		}
	}()

	if strings.TrimSpace(req.Message) == "" {
		return &AIResult{
			Intent:   IntentUnknown,
			Language: LangID,
			Entities: map[string]interface{}{},
			Action:   Action{Type: ActionClarify},
			Response: "Silakan tulis pesan Anda.",
		}, nil
	}

	// Create context with timeout 8s if not set
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 8*time.Second)
		defer cancel()
	}

	// Try primary provider
	raw, err := s.primary.Generate(ctx, req)
	if err != nil {
		// Fallback to mock provider (controlled error, don't crash)
		// Log fallback internally (caller can log)
		fallbackRaw, fallbackErr := s.fallback.Generate(ctx, req)
		if fallbackErr != nil {
			// Even fallback failed -> return safe UNKNOWN
			return fallbackResult(req, langFromReq(req)), fmt.Errorf("primary %s failed: %v; fallback also failed: %v", s.primary.Name(), err, fallbackErr)
		}
		// Validate fallback output
		validated, vErr := Validate(fallbackRaw, req.RoomNumber)
		if vErr != nil {
			return fallbackResult(req, fallbackRaw.Language), fmt.Errorf("validation after fallback failed: %v", vErr)
		}
		return validated, nil
	}

	// Validate primary output
	validated, err := Validate(raw, req.RoomNumber)
	if err != nil {
		// Invalid AI output (malformed, invented room, etc.) -> fallback to mock logic
		fallbackRaw, fallbackErr := s.fallback.Generate(ctx, req)
		if fallbackErr != nil {
			return fallbackResult(req, LangID), fmt.Errorf("validation failed: %v; fallback failed: %v", err, fallbackErr)
		}
		validated2, v2Err := Validate(fallbackRaw, req.RoomNumber)
		if v2Err != nil {
			return fallbackResult(req, fallbackRaw.Language), fmt.Errorf("validation failed: %v; fallback validation failed: %v", err, v2Err)
		}
		return validated2, nil
	}

	return validated, nil
}

func langFromReq(req Request) string {
	lang := DetectLanguage(req.Message)
	if lang == "" {
		return LangID
	}
	return lang
}

func fallbackResult(req Request, lang string) *AIResult {
	if lang != LangEN {
		lang = LangID
	}
	resp := "Maaf, layanan AI sedang mengalami gangguan. Silakan coba kembali beberapa saat lagi."
	if lang == LangEN {
		resp = "Sorry, AI service is temporarily unavailable. Please try again shortly."
	}
	// Per ERROR FLOW.md: fallback response, do not create ticket
	return &AIResult{
		Intent:   IntentUnknown,
		Language: lang,
		Entities: map[string]interface{}{},
		Action:   Action{Type: ActionClarify},
		Response: resp,
	}
}

// ProcessSimple is helper for tests without context
func (s *Service) ProcessSimple(message, roomNumber string) (*AIResult, error) {
	return s.Process(context.Background(), Request{
		SessionID:  "test-session",
		RoomNumber: roomNumber,
		Message:    message,
	})
}
