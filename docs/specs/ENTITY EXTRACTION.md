Entity Extraction

The AI should extract relevant information when available.

Possible entities:

room_number
quantity
item
problem
budget
location
category
language
time
preference

Example:

Guest:

"Tolong antar 3 handuk ke kamar 305."

Output:

{
  "intent": "TOWEL_REQUEST",
  "room_number": "305",
  "quantity": 3,
  "item": "towel",
  "department": "HOUSEKEEPING",
  "priority": "MEDIUM",
  "requires_ticket": true
}