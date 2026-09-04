package handler

import (
	"fmt"
	"io"
	"net/http"

	"batiqa-ai/internal/model"
)

// Photo handles POST /api/chat/photo (guest only): multipart photo of a room
// problem -> Gemini Vision -> same ChatResponse contract as text chat.
const maxPhotoBytes = 8 << 20 // 8 MB

var allowedPhotoMimes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func (h *ChatHandler) Photo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	if err := r.ParseMultipartForm(maxPhotoBytes); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid form upload")
		return
	}
	file, header, err := r.FormFile("photo")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "photo file is required")
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	if !allowedPhotoMimes[mimeType] {
		// Sniff by first bytes as some browsers send octet-stream
		buf := make([]byte, 512)
		n, _ := io.ReadFull(file, buf)
		mimeType = sniffMime(buf[:n])
		if !allowedPhotoMimes[mimeType] {
			WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Unsupported image type (use JPG/PNG/WebP)")
			return
		}
		// Rewind not possible on multipart File; read the rest for processing
		rest, _ := io.ReadAll(file)
		h.processPhoto(w, r, append(buf[:n], rest...), mimeType, r.FormValue("session_id"), r.FormValue("message"))
		return
	}
	data, err := io.ReadAll(http.MaxBytesReader(w, file, maxPhotoBytes))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Photo too large (max 8MB)")
		return
	}
	h.processPhoto(w, r, data, mimeType, r.FormValue("session_id"), r.FormValue("message"))
}

func (h *ChatHandler) processPhoto(w http.ResponseWriter, r *http.Request, image []byte, mimeType, sessionID, caption string) {
	sessionID = trimStr(sessionID)
	if sessionID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "session_id is required")
		return
	}
	// Verified room from QR token wins; guests cannot claim another room.
	room := GuestRoom(r)

	displayMsg := "[Foto] " + caption
	_ = h.convs.Create(&model.Conversation{
		SessionID: sessionID,
		Role:      model.RoleUser,
		Message:   displayMsg,
		Intent:    nil,
	})

	aiReq := h.newAIRequest(sessionID, room, caption)
	aiResult, _ := h.ai.ProcessImage(r.Context(), *aiReq, image, mimeType)
	if aiResult == nil {
		fb := aiFallbackResult()
		aiResult = &fb
	}

	// Rule-based fallback cannot see images: keep its clarify answer.
	// Persist assistant reply.
	_ = h.convs.Create(&model.Conversation{
		SessionID: sessionID,
		Role:      model.RoleAssistant,
		Message:   aiResult.Response,
		Intent:    &aiResult.Intent,
	})

	ticketID := h.maybeCreateTicket(r, aiResult, displayMsg, room)

	WriteJSON(w, http.StatusOK, ChatResponse{
		Message:        aiResult.Response,
		Intent:         aiResult.Intent,
		RequiresTicket: ticketID != nil || (aiResult.RequiresTicket && ticketID == nil),
		TicketID:       ticketID,
	})
}

func sniffMime(b []byte) string {
	switch {
	case len(b) >= 3 && b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF:
		return "image/jpeg"
	case len(b) >= 8 && b[0] == 0x89 && b[1] == 'P' && b[2] == 'N' && b[3] == 'G':
		return "image/png"
	case len(b) >= 12 && string(b[8:12]) == "WEBP":
		return "image/webp"
	default:
		return ""
	}
}

func trimStr(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

// GuestMe handles GET /api/guest/me -> {room_number} from the verified token.
func (h *ChatHandler) GuestMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	room := GuestRoom(r)
	if room == "" {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Guest token required")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"room_number": room})
}

var _ = fmt.Sprintf
