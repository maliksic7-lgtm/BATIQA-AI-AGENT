# 🏨 BATIQA AI Guest Assistant & Operations Router

> **SSIA AI Innovation Challenge 2026 Project**  
> *Web-based AI Guest Experience & Automated Service Ticket System for BATIQA Hotels.*

---

## 📌 Project Overview

**BATIQA AI Guest Assistant** adalah solusi berbasis AI web-app yang dapat diakses tamu hotel secara instan cukup dengan **scan QR Code** di kamar tanpa perlu mengunduh aplikasi. 

Aplikasi ini bertindak sebagai *interface* pintar antara tamu dan operasional hotel untuk:
1. **AI Guest Experience:** Menjawab pertanyaan fasilitas, rekomendasi kuliner/lokal sesuai budget, dan pemanduan kamar 24/7.
2. **Automated Operations Router:** Mengubah *natural language request* dari tamu (seperti *"AC tidak dingin"* atau *"minta handuk"*) menjadi tiket operasional terstruktur secara otomatis ke departemen **Housekeeping** atau **Engineering**.

---

## 🎯 Key Features

- **📱 Zero-Install QR Access:** Akses cepat langsung melalui browser HP tamu.
- **🤖 Smart AI Concierge & Recommendation:** Pemrosesan bahasa alami untuk menjawab FAQ hotel & rekomendasi tempat sekitar berbasis budget.
- **🧹 Housekeeping & Maintenance Ticket Dispatch:** Mengonversi pesan tamu menjadi *task ticket* real-time (Room Number, Request Type, Time, Priority).
- **🗣️ Multi-Language Support:** Penerjemah otomatis antara tamu asing dan staf operasional.

---

## 🛠️ Tech Stack

- **Backend:** Golang (REST API, Webhook integration)
- **Database:** MySQL
- **Frontend:** Lightweight HTML/CSS/JS (Mobile-first web view)
- **AI Engine:** Gemini API / OpenAI API (Intent extraction & natural response)

---

## 📂 Project Structure

```text
.
├── cmd/
│   └── api/                # Entry point aplikasi Go
├── internal/
│   ├── config/             # Konfigurasi database & API Key
│   ├── handler/            # Endpoint HTTP handlers
│   ├── model/              # Struct data (Guest, Ticket, HotelData)
│   ├── repository/         # Query MySQL
│   └── service/            # Logika AI Engine & Intent Extraction
├── migrations/             # Skema & DDL MySQL
├── web/                    # Asset Frontend (HTML, CSS, JS)
├── go.mod
└── README.md