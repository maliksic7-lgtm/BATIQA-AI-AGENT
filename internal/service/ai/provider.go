package ai

import "context"

// Provider is abstraction for AI engine per TECH STACK.md:10
// Primary: Gemini API, optional: OpenAI, also Mock/RuleBased for testing.
// Business logic must not depend on concrete provider.
type Provider interface {
	// Generate returns RawAIOutput for given prompt context. Must not panic.
	Generate(ctx context.Context, req Request) (*RawAIOutput, error)
	// Name returns provider name for logging
	Name() string
}

// ProviderError indicates provider failure (network, auth, quota)
type ProviderError struct {
	Provider string
	Err      error
}

func (e *ProviderError) Error() string {
	return "provider " + e.Provider + " failed: " + e.Err.Error()
}
