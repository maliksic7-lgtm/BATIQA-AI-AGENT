Update Ticket Status
PATCH
/api/tickets/:id/status

Request:

{
  "status": "IN_PROGRESS"
}

Allowed values:

OPEN
IN_PROGRESS
RESOLVED
CANCELLED

Response:

{
  "ticket_number": "TKT-1024",
  "status": "IN_PROGRESS"
}