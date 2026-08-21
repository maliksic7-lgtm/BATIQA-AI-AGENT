package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// GeminiProvider calls Gemini API (https://generativelanguage.googleapis.com) if GEMINI_API_KEY is set.
// For Phase 3, it uses HTTP directly without external SDK to keep dependencies minimal.
// If key is missing or call fails, caller should fallback to MockProvider per error handling spec.
type GeminiProvider struct {
	apiKey string
	model  string
	client *http.Client
}

func NewGeminiProvider() *GeminiProvider {
	key := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	model := strings.TrimSpace(os.Getenv("GEMINI_MODEL"))
	if model == "" {
		model = "gemini-1.5-flash"
	}
	return &GeminiProvider{
		apiKey: key,
		model:  model,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (g *GeminiProvider) Name() string { return "gemini" }

// Generate implements Provider. It builds prompt, calls Gemini, parses JSON, and returns RawAIOutput.
// On any failure, returns ProviderError so caller can fallback.
func (g *GeminiProvider) Generate(ctx context.Context, req Request) (*RawAIOutput, error) {
	if g.apiKey == "" {
		return nil, &ProviderError{Provider: g.Name(), Err: fmt.Errorf("GEMINI_API_KEY not set")}
	}

	prompt := buildGeminiPrompt(req)

	// Gemini API request body per https://ai.google.dev/api/rest/v1/models/generateContent
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]interface{}{{"text": prompt}}},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.2,
			"maxOutputTokens": 512,
			"responseMimeType": "application/json",
		},
	}
	b, _ := json.Marshal(body)

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, &ProviderError{Provider: g.Name(), Err: err}
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(httpReq)
	if err != nil {
		return nil, &ProviderError{Provider: g.Name(), Err: err}
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, &ProviderError{Provider: g.Name(), Err: fmt.Errorf("gemini status %d: %s", resp.StatusCode, string(respBytes))}
	}

	// Parse Gemini response
	var gemResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(respBytes, &gemResp); err != nil {
		return nil, &ProviderError{Provider: g.Name(), Err: fmt.Errorf("unmarshal gemini: %w", err)}
	}
	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return nil, &ProviderError{Provider: g.Name(), Err: fmt.Errorf("empty gemini candidates")}
	}
	text := gemResp.Candidates[0].Content.Parts[0].Text
	text = strings.TrimSpace(text)
	// Try to extract JSON from markdown code block if present
	if strings.Contains(text, "```") {
		// extract between ```json and ```
		start := strings.Index(text, "{")
		end := strings.LastIndex(text, "}")
		if start != -1 && end != -1 && end > start {
			text = text[start : end+1]
		}
	}

	var raw RawAIOutput
	if err := json.Unmarshal([]byte(text), &raw); err != nil {
		return nil, &ProviderError{Provider: g.Name(), Err: fmt.Errorf("parse gemini JSON: %w text=%s", err, text)}
	}
	return &raw, nil
}

// buildGeminiPrompt creates structured prompt for Gemini to return JSON per STRUCTURED AI OUTPUT.md
func buildGeminiPrompt(req Request) string {
	// Keep prompt minimal and in Indonesian/English bilingual instruction for language preservation
	return fmt.Sprintf(`You are BATIQA AI Guest Assistant. Analyze guest message and return ONLY valid JSON.

Guest message: "%s"
Room number (if provided): "%s"
Session: %s

You must do:
1. Detect language: "id" for Indonesian, "en" for English. Respond in same language.
2. Classify intent from: %s
3. Extract entities: room_number, quantity, item, problem, budget, category, preference. Do NOT invent room_number if not in message or provided.
4. Route department: HOUSEKEEPING for TOWEL_REQUEST/HOUSEKEEPING_REQUEST/AMENITY_REQUEST/ROOM_CLEANING_REQUEST, ENGINEERING for AC_PROBLEM/TV_PROBLEM/WIFI_PROBLEM/LIGHT_PROBLEM/SHOWER_PROBLEM/PLUMBING_PROBLEM/ROOM_EQUIPMENT_PROBLEM/GENERAL_MAINTENANCE
5. Priority: HIGH for AC completely unavailable/water leakage/safety, MEDIUM for towel/cleaning/minor TV, LOW for extra pillow/amenities/general info. Do NOT exaggerate.
6. Action: CREATE_TICKET if intent requires ticket, ANSWER if information, CLARIFY if UNKNOWN.
7. Natural response in guest language.

Return JSON exactly:
{"intent":"...","language":"id|en","entities":{"room_number":"...","quantity":1,"item":"towel","problem":"...","budget":100000},"action":{"type":"CREATE_TICKET|ANSWER|CLARIFY","department":"HOUSEKEEPING|ENGINEERING","priority":"LOW|MEDIUM|HIGH"},"response":"..."}

Valid intents: %s
Do not fabricate hotel info, ticket IDs, or room numbers. If unsure, intent=UNKNOWN.
`,
		req.Message,
		req.RoomNumber,
		req.SessionID,
		strings.Join(getIntentList(), ", "),
		strings.Join(getIntentList(), ", "),
	)
}

func getIntentList() []string {
	list := make([]string, 0, len(ValidIntents))
	for k := range ValidIntents {
		list = append(list, k)
	}
	return list
}
