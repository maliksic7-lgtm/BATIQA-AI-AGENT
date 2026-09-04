Create Ticket
POST
/api/tickets

Request:

{
  "room_number": "305",
  "department": "ENGINEERING",
  "category": "AC_PROBLEM",
  "description": "AC tidak dingin",
  "priority": "HIGH"
}

Response:

{
  "ticket_number": "TKT-1024",
  "room_number": "305",
  "department": "ENGINEERING",
  "category": "AC_PROBLEM",
  "priority": "HIGH",
  "status": "OPEN"
}