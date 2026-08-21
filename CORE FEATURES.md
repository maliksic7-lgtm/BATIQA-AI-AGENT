⭐ Core Features
1. 🤖 AI Hotel Concierge

AI dapat menjawab pertanyaan mengenai:

Breakfast
Restaurant
Swimming pool
Gym
WiFi
Check-in
Check-out
Hotel facilities
Room facilities
Hotel policies
2. 🛎️ Service Request

Guest dapat meminta:

Additional towel
Housekeeping
Amenities
Room service
Cleaning
Other hotel services

AI akan mengubah request menjadi structured ticket.

Example:

"Tolong antar 2 handuk ke kamar saya."

AI:

{
  "intent": "TOWEL_REQUEST",
  "department": "HOUSEKEEPING",
  "quantity": 2,
  "item": "towel",
  "priority": "MEDIUM",
  "requires_ticket": true
}

3. 🔧 Maintenance / Complaint Assistant

Guest dapat melaporkan:

AC tidak dingin.
TV bermasalah.
WiFi tidak bekerja.
Lampu rusak.
Shower bermasalah.
Fasilitas kamar lainnya.

Example:

"AC kamar saya tidak dingin."

AI:

{
  "intent": "AC_PROBLEM",
  "department": "ENGINEERING",
  "priority": "HIGH",
  "requires_ticket": true
}

4. 📍 Local Recommendation

AI dapat memberikan rekomendasi:

Restaurant
Cafe
Mall
Tourist destination
ATM
Transportation
Local food

Recommendation dapat mempertimbangkan:

Budget
Category
Distance
User preference

5. 🗣️ Multi-language

Guest dapat berkomunikasi menggunakan bahasa seperti:

Indonesian
English

AI akan memahami bahasa guest dan memberikan respons dalam bahasa yang sama.

6. 📊 Operations Dashboard

Staff memiliki dashboard untuk melihat:

Open tickets
High priority tickets
Department
Room number
Request
Timestamp
Status

Ticket lifecycle:

OPEN
 ↓
IN_PROGRESS
 ↓
RESOLVED