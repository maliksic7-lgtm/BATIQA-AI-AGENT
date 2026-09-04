Structured AI Output

The AI should return structured data internally.

Example:

{
  "intent": "AC_PROBLEM",
  "language": "id",
  "entities": {
    "room_number": "305",
    "problem": "AC tidak dingin"
  },
  "action": {
    "type": "CREATE_TICKET",
    "department": "ENGINEERING",
    "priority": "HIGH"
  },
  "response": "Baik, saya akan membantu melaporkan masalah AC Anda."
}

Backend must validate:

Intent exists.
Department is allowed.
Priority is allowed.
Room number is valid.
Action type is allowed.