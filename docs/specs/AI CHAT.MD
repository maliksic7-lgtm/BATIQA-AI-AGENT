AI Chat
POST
/api/chat

Request:

{
  "session_id": "guest-session-123",
  "room_number": "305",
  "message": "AC kamar saya tidak dingin"
}

Response:

{
  "message": "Baik, saya akan membantu melaporkan masalah AC Anda.",
  "intent": "AC_PROBLEM",
  "requires_ticket": true,
  "ticket_id": null
}