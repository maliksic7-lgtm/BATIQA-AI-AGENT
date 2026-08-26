package ticket

import (
	"errors"
	"fmt"
	"strings"
)

// ErrTicketNotFound is returned when a ticket does not exist.
var ErrTicketNotFound = errors.New("ticket not found")

// ErrRoomRequired is returned when a ticket cannot be created because the
// room number is missing (per MISSING ROOM NUMBER.md - never invent rooms).
var ErrRoomRequired = errors.New("room_number is required")

// ErrInvalidTransition is returned for disallowed status transitions per TICKET LIFECYCLE.md.
var ErrInvalidTransition = errors.New("invalid status transition")

// ValidationError marks user-input validation failures (HTTP 400).
type ValidationError struct{ msg string }

func (e *ValidationError) Error() string { return e.msg }

// NewValidationError wraps a message as a validation error.
func NewValidationError(format string, args ...interface{}) error {
	return &ValidationError{msg: fmt.Sprintf(format, args...)}
}

// IsValidationError reports whether err (or its chain) is a ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// validateRoomNumberFormat checks room format and returns a ValidationError on failure.
func validateRoomNumberFormat(room string) error {
	if !isValidRoomNumber(room) {
		return NewValidationError("invalid room_number format: %s", room)
	}
	return nil
}

// normalizeEnum uppercases/trims an enum-ish input.
func normalizeEnum(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
