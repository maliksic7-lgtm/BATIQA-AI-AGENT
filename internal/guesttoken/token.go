// Package guesttoken creates and verifies HMAC-signed guest tokens embedded in QR codes.
// Format: b64url("ROOM=<room>|EXP=<unix>") + "." + hex(hmac-sha256(payload)).
package guesttoken

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	ErrInvalid   = errors.New("invalid guest token")
	ErrExpired   = errors.New("guest token expired")
	ErrBadFormat = errors.New("malformed guest token")
	ErrBadRoom   = errors.New("missing room in guest token")
	ErrNoSecret  = errors.New("GUEST_TOKEN_SECRET not set")
)

const DefaultDevSecret = "batiqa-dev-secret-change-me"

// Secret returns the signing secret from env (dev fallback documented in README).
func Secret() []byte {
	s := strings.TrimSpace(os.Getenv("GUEST_TOKEN_SECRET"))
	if s == "" {
		s = DefaultDevSecret
	}
	return []byte(s)
}

// New issues a signed token for a room lasting ttl.
func New(room string, ttl time.Duration, secret []byte) (string, error) {
	room = strings.TrimSpace(room)
	if room == "" {
		return "", ErrBadRoom
	}
	exp := time.Now().Add(ttl).UTC().Unix()
	payload := fmt.Sprintf("ROOM=%s|EXP=%d", room, exp)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig, nil
}

// Parse verifies signature and expiry, returning the room number.
func Parse(token string, secret []byte) (string, error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", ErrBadFormat
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", ErrInvalid
	}
	payload := string(raw)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	if !hmac.Equal([]byte(hex.EncodeToString(mac.Sum(nil))), []byte(parts[1])) {
		return "", ErrInvalid
	}
	var room string
	var exp int64
	for _, kv := range strings.Split(payload, "|") {
		switch {
		case strings.HasPrefix(kv, "ROOM="):
			room = strings.TrimPrefix(kv, "ROOM=")
		case strings.HasPrefix(kv, "EXP="):
			fmt.Sscanf(strings.TrimPrefix(kv, "EXP="), "%d", &exp)
		}
	}
	if room == "" {
		return "", ErrBadRoom
	}
	if exp > 0 && time.Now().UTC().Unix() > exp {
		return "", ErrExpired
	}
	return room, nil
}
