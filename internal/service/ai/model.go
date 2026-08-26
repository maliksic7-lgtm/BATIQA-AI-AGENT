package ai

// Structured output per docs/STRUCTURED AI OUTPUT.md
type AIResult struct {
	Intent   string                 `json:"intent"`
	Language string                 `json:"language"`
	Entities map[string]interface{} `json:"entities"`
	Action   Action                 `json:"action"`
	Response string                 `json:"response"`
	// Helper fields for handler (not part of raw AI JSON but derived)
	RequiresTicket bool `json:"-"`
}

type Action struct {
	Type       string `json:"type"`
	Department string `json:"department,omitempty"`
	Priority   string `json:"priority,omitempty"`
}

// RawAIOutput is what provider returns before validation
type RawAIOutput struct {
	Intent   string                 `json:"intent"`
	Language string                 `json:"language"`
	Entities map[string]interface{} `json:"entities"`
	Action   Action                 `json:"action"`
	Response string                 `json:"response"`
}

// Turn is one message in conversation history for context.
type Turn struct {
	Role    string // "user" or "assistant"
	Content string
}

// Request to AI service
type Request struct {
	SessionID  string
	RoomNumber string   // may be empty (missing)
	Message    string   // latest guest message
	History    []Turn   // prior turns (oldest first), excludes Message
	Facts      []string // verified hotel/recommendation facts the AI may cite verbatim
}

// Known action types
const (
	ActionCreateTicket = "CREATE_TICKET"
	ActionAnswer       = "ANSWER"
	ActionClarify      = "CLARIFY"
)

// Known languages
const (
	LangID = "id"
	LangEN = "en"
)
