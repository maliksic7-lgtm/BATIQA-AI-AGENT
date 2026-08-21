Get Tickets
GET
/api/tickets

Optional filters:

/api/tickets?department=ENGINEERING
/api/tickets?status=OPEN
/api/tickets?priority=HIGH

Response:

{
  "tickets": [
    {
      "ticket_number": "TKT-1024",
      "room_number": "305",
      "department": "ENGINEERING",
      "category": "AC_PROBLEM",
      "description": "AC tidak dingin",
      "priority": "HIGH",
      "status": "OPEN",
      "created_at": "2026-08-21T16:20:00Z"
    }
  ]
}