# BATIQA AI Guest Assistant

AI-powered guest assistant untuk **Hotel BATIQA Pekanbaru** — dibuat untuk SSIA AI Innovation Competition.

Tamu scan QR di kamar → web app (tanpa install aplikasi) → chat bahasa natural → AI deteksi intent & ekstraksi entitas → tiket otomatis ter-route ke Housekeeping/Engineering → staff kelola via Operations Dashboard.

## Fitur

- **AI Concierge Chat** — info hotel, service request, laporan kerusakan, rekomendasi lokal (ID/EN)
- **Anti-halusinasi** — jawaban info hotel & rekomendasi WAJIB dari database terverifikasi, bukan karangan AI
- **Ticket System** — lifecycle OPEN → IN_PROGRESS → RESOLVED (+ CANCELLED), prioritas LOW/MEDIUM/HIGH
- **Operations Dashboard** — statistik, filter departemen/status/prioritas, ubah status & prioritas, assign tiket
- **Riwayat Chat** — tersimpan per sesi tamu

## Tech Stack

| Layer | Teknologi |
|---|---|
| Backend | Go 1.24 (`net/http` murni, tanpa framework) |
| Database | MySQL 8 |
| AI | Gemini API via REST + provider abstraction (fallback rule-based) |
| Frontend | HTML/CSS/JS vanilla — tema resmi BATIQA (#033C5A / #C9AC85) |

Arsitektur berlapis: `Handler → Service → Repository → MySQL`. AI tidak pernah menyentuh DB langsung — semua output tervalidasi backend.

## Struktur Proyek

```
├── cmd/
│   ├── api/main.go        # entry point server
│   └── migrate/main.go    # CLI migrasi
├── internal/
│   ├── config/            # env-based configuration
│   ├── handler/           # HTTP handlers (chat, ticket, auth, dst.)
│   ├── model/             # struct DB + enum validasi
│   ├── repository/        # akses MySQL (parameterized queries)
│   ├── router/            # registrasi route, middleware, static files
│   └── service/
│       ├── ai/            # pipeline AI: language→intent→entity→routing→priority→action
│       └── ticket/        # business logic tiket (validasi + sentinel errors)
├── migrations/            # SQL skema + seed (idempotent)
├── web/                   # frontend guest (mobile-first) + staff dashboard
├── docker-compose.yml     # MySQL saja; API jalan lokal
└── go.mod
```

## Quick Start

### Prasyarat
- Go 1.24+
- Docker (untuk MySQL) ATAU MySQL lokal — tanpa keduanya server tetap jalan dalam mode degraded

### Cara Paling Gampang (Windows)

```bat
run.bat
```

Satu perintah untuk semuanya: nyalakan MySQL via Docker (jika ada), build, migrasi + seed, jalankan server. Tinggal buka:

- Tamu: http://localhost:8080/
- Staff: http://localhost:8080/staff/login.html (`admin@batiqa.com` / `batiqa123`)

### Manual (Linux/macOS/CI)

```bash
cp .env.example .env          # sesuaikan DB_DSN, isi GEMINI_API_KEY jika ada
docker-compose up -d          # MySQL
go run ./cmd/migrate          # migrasi + seed
go run ./cmd/api              # server di :8080
```

### Kredensial Demo Staff

| Email | Password | Departemen |
|---|---|---|
| admin@batiqa.com | batiqa123 | ADMIN |
| hk@batiqa.com | batiqa123 | HOUSEKEEPING |
| eng@batiqa.com | batiqa123 | ENGINEERING |

> DEMO ONLY — ganti kredensial sebelum production.

## API

| Endpoint | Method | Auth | Fungsi |
|---|---|---|---|
| `/api/chat` | POST | - | Chat AI tamu `{session_id, room_number?, message}` |
| `/api/conversations?session_id=` | GET | - | Riwayat chat per sesi |
| `/api/tickets` | POST/GET | - | Buat tiket / daftar tiket (filter dept/status/priority/room) |
| `/api/tickets/{ticket_number}` | GET | - | Detail tiket |
| `/api/tickets/{id}/status` | PATCH | Staff | Ubah status (transisi tervalidasi) |
| `/api/tickets/{id}/priority` | PATCH | Staff | Ubah prioritas |
| `/api/tickets/{id}/assign` | POST | Staff | Assign tiket ke staff |
| `/api/tickets/{id}/assignments` | GET | - | Daftar assignment |
| `/api/tickets/stats` | GET | Staff | Statistik dashboard |
| `/api/hotel-info?category=` | GET | - | Info hotel terverifikasi |
| `/api/recommendations?category=&max_price=` | GET | - | Rekomendasi lokal terverifikasi |
| `/api/staff/login` | POST | - | Login staff (bcrypt) |
| `/api/staff/me` · `/logout` | GET/POST | Staff | Profil & logout |
| `/api/health` | GET | - | Health check |

Format error seragam: `{"error": {"code": "INVALID_REQUEST", "message": "..."}}`.

## Testing

```bash
go vet ./...
go test ./...      # test DB otomatis skip jika MySQL tidak tersedia
gofmt -l .         # harus kosong
```

## Aturan Bisnis Kunci

1. AI dilarang mengarang info hotel, harga, nomor kamar, atau ticket ID
2. Room number wajib dari tamu — jika hilang, AI bertanya dulu
3. Intent UNKNOWN → minta klarifikasi, jangan menebak
4. Tiket hanya dibuat jika request sukses divalidasi backend
5. Prioritas proporsional: LOW (bantal ekstra), MEDIUM (handuk), HIGH (AC mati/bocor)
6. RESOLVED & CANCELLED bersifat terminal (tidak bisa diubah lagi)
7. Endpoint tulis staff wajib autentikasi Bearer token

## Dokumentasi

Spesifikasi lengkap ada di file `.MD` di root repo (GOALS, MVP SCOPE, AI ARCHITECTURE, TICKET LIFECYCLE, dst.) dan `docs/`.

## License

MIT — Built for SSIA AI Innovation 2024
