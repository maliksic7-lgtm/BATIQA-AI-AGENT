Visual Personality

The application should communicate:

Hospitality
Trust
Technology
Simplicity
Efficiency

It should look like a real hotel digital product rather than a student chatbot prototype.



---


## Nah, sekarang repo lu punya fondasi yang jauh lebih enak buat vibecoding


Strukturnya jadi:


```text
BATIQA-AI/
│
├── README.md                 ← MASTER SOURCE OF TRUTH
│
├── docs/
│   ├── AI_SPEC.md            ← OTAK AI
│   ├── USER_FLOW.md          ← ALUR SISTEM
│   ├── DATABASE.md           ← DATA & MYSQL
│   ├── API.md                ← KONTRAK BACKEND
│   └── UI_SPEC.md            ← TAMPILAN
│
├── cmd/
│   └── api/
│
├── internal/
│   ├── config/
│   ├── handler/
│   ├── model/
│   ├── repository/
│   ├── router/
│   └── service/
│       ├── ai/
│       ├── ticket/
│       └── recommendation/
│
├── migrations/
│
├── web/
│   ├── guest/
│   ├── staff/
│   ├── css/
│   └── js/
│
├── .env.example
└── go.mod