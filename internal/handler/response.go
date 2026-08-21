package handler

import (
	"encoding/json"
	"net/http"
)

// ErrorEnvelope follows docs/ERROR_FORMAT.MD
// { "error": { "code": "INVALID_REQUEST", "message": "..." } }
type ErrorEnvelope struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteJSON writes JSON response with status code.
func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteError writes standardized error envelope.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, ErrorEnvelope{
		Error: ErrorDetail{Code: code, Message: message},
	})
}

// WriteOK is helper for 200 JSON.
func WriteOK(w http.ResponseWriter, payload interface{}) {
	WriteJSON(w, http.StatusOK, payload)
}
