# Phase 1 — Project & Backend Foundation — Notes

**Lokasi Implementasi:** `C:\Users\Thin 15\BATIQA-AI` (writable)
**Lokasi Asli (read-only):** `D:\AI ASSISTANT BATIQA` — butuh Fix Permissions untuk sync

## Kenapa tidak di D:\ ?
`D:\` dan `D:\AI ASSISTANT BATIQA` memiliki ACL `Users:(RX)` saja (Medium Mandatory Level, Administrators deny). Proses tidak elevated tidak bisa write.
Solusi: jalankan `FIX_PERMISSIONS.bat` sebagai Administrator, lalu copy:

```powershell
robocopy "C:\Users\Thin 15\BATIQA-AI" "D:\AI ASSISTANT BATIQA" /E /XD .git bin gopath
```

Atau jalankan VS Code / Terminal sebagai Administrator.

## Struktur Phase 1 (sesuai README.md:37 & PROJECT STRUCTURE.md)
```
C:\Users\Thin 15\BATIQA-AI\
├── cmd/api/main.go
├── internal/config/config.go
├── internal/handler/health.go, response.go
├── internal/router/router.go, middleware.go
├── internal/model/placeholder.go
├── internal/repository/placeholder.go
├── internal/service/placeholder.go (+ ai/ticket/recommendation dirs)
├── internal/router/
├── migrations/README.md
├── web/index.html, web/guest, web/staff, web/css, web/js
├── go.mod
├── .env.example
├── .gitignore
└── 57 .md docs (copied from D:\)
```

## Config
- `PORT` default 8080, `ENV` development, `LOG_LEVEL` info
- `.env` loader minimal tanpa external dep (sesuai SECURITY PRINCIPLES.md: API keys server-side only)
- `Addr()` helper, `IsProduction()`

## Router & Health
- `GET /api/health` → `{"status":"ok"}` (HEALTH CHECK.md)
- 404 JSON `{"error":{"code":"NOT_FOUND","message":"Endpoint not found"}}` (ERROR FORMAT.MD)
- 405 untuk method salah
- Middleware: logging, recovery (panic→500), CORS (* for dev)

## Main
- Graceful shutdown 10s via SIGINT/SIGTERM (context.WithTimeout)

## Testing
- `go fmt ./...` pass
- `go vet ./...` pass
- `go test ./...` pass (no test files yet)
- `go build -o bin/api.exe ./cmd/api` → 8.5 MB
- `curl /api/health` 200 OK, CORS headers, JSON valid
