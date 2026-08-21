package ai

import "context"

// FailingProvider always returns error (for testing provider failure fallback)
type FailingProvider struct {
	name string
}

func NewFailingProvider() *FailingProvider { return &FailingProvider{name: "failing"} }
func (f *FailingProvider) Name() string { return f.name }
func (f *FailingProvider) Generate(ctx context.Context, req Request) (*RawAIOutput, error) {
	return nil, &ProviderError{Provider: f.name, Err: context.DeadlineExceeded}
}

// MalformedProvider returns invalid/malformed output for testing validation layer
type MalformedProvider struct {
	raw *RawAIOutput
}

func NewMalformedProvider(raw *RawAIOutput) *MalformedProvider { return &MalformedProvider{raw: raw} }
func (m *MalformedProvider) Name() string { return "malformed" }
func (m *MalformedProvider) Generate(ctx context.Context, req Request) (*RawAIOutput, error) {
	return m.raw, nil
}
