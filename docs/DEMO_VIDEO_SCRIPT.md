# Skrip Video Demo — BATIQA AI Guest Assistant (150 Detik)

Kompetisi: **SSIA AI Innovation** · Video demo produk untuk juri.
Total durasi: **2 menit 30 detik** · Bahasa narasi: Indonesia.

| # | Timestamp | Durasi | Aksi Layar | Voice-Over (ID) | Overlay Teks |
|---|---|---|---|---|---|
| 1 | 00:00–00:15 | 15s | Judul + b-roll resepsionis hotel / telepon FO berdering terus. Cut cepat 3 suasana: tamu telepon, staf mencatat manual, tamu asing bingung bahasa. | "Setiap hari, front office hotel dibanjiri panggilan untuk hal sepele: handuk ekstra, jam sarapan, password WiFi. Ditambah hambatan bahasa, dan request lisan yang tidak pernah tercatat." | `Masalah: panggilan sepele · hambatan bahasa · request tanpa jejak` |
| 2 | 00:15–00:35 | 20s | Kamera HP mengarah ke QR di dinding kamar (mock). Jari scan → landing chat BATIQA terbuka di browser HP. Tidak ada install, tidak ada login manual. | "Inilah BATIQA AI Guest Assistant. Tamu cukup scan satu QR code di kamar — langsung masuk lewat browser, tanpa download aplikasi apa pun." | `1 QR scan → tanpa aplikasi` |
| 3 | 00:35–00:55 | 20s | Ketik di chat: *"What time is breakfast?"* → jawaban EN muncul instan dari data hotel. Lanjut ketik versi ID *"sarapan jam berapa?"* → jawaban ID. | "Tamu bisa bertanya dalam bahasa apa pun yang mereka kuasai. AI mendeteksi bahasa otomatis dan menjawab akurat — karena semua info hotel diambil dari database terverifikasi, bukan karangan AI." | `Multi-bahasa ID / EN` · `Jawaban dari DB terverifikasi` |
| 4 | 00:55–01:20 | 25s | Ketik *"AC kamar saya mati"* + lampirkan foto unit AC. AI merespons dengan empati, klasifikasi prioritas HIGH → konfirmasi tiket **TKT-0000XX – ENGINEERING** tampil. | "Sekarang kasus nyata: AC mati. Tamu cukup kirim pesan plus foto kerusakannya. AI membuat tiket lengkap — nomor tiket, departemen Engineering, prioritas HIGH — semuanya tervalidasi backend sebelum tersimpan." | `Photo-to-Ticket Vision` · `TKT-0000XX · HIGH · Engineering` |
| 5 | 01:20–01:50 | 30s | Cut ke laptop: Operations Dashboard staff. Tiket baru **muncul live tanpa refresh** (SSE). Staff klik assign ke teknisi → ubah status IN_PROGRESS → kerjakan → RESOLVED. Sisipan singkat: status juga terlihat oleh tamu di chat. | "Di sisi staff, dashboard operasional menerima tiket secara real-time. Satu klik untuk assign ke teknisi, update progres, hingga selesai. Setiap langkah terekam sebagai jejak audit — dan tamu bisa melihat statusnya langsung dari chat." | `Real-time via SSE` · `Assign → In Progress → Resolved` |
| 6 | 01:50–02:10 | 20s | Buka tab statistik dashboard: jumlah tiket per departemen, distribusi prioritas, tren harian. Scroll ringkas. | "Manajemen mendapat analytics harian otomatis: volume permintaan, kategori terbanyak, distribusi prioritas, dan waktu penanganan. Tidak ada lagi rekap manual." | `Analytics & SLA otomatis` |
| 7 | 02:10–02:30 | 20s | Closing: logo wordmark BATIQA, tagline, alur mini QR→Chat→Ticket→Staff, stack line. End card kontak tim. | "BATIQA AI Guest Assistant — smart concierge for a smarter stay. Dibangun dengan Go, MongoDB, dan Gemini. Request tamu tercepat yang pernah ada: dari kamar, langsung ditangani." | `BATIQA AI Guest Assistant` · `Go · MongoDB · Gemini` · `[kontak tim]` |

## Catatan Timing & Ritme

- Total = 150 detik persis; jaga tiap segmen ±2 detik dari rencana saat editing.
- Aturan praktis: maksimal ~35 kata voice-over per 15 detik agar tidak terburu-buru.
- Transisi antar shot: hard cut sederhana atau crossfade 8 frame; hindari efek berlebihan.
- Musik: instrumental tenang volume −20 dB di bawah voice-over; fade out di end card.

## Tips Rekaman

### Persiapan Data Demo
- Reset database ke kondisi bersih agar tiket rapi: `docker-compose down -v`, lalu jalankan ulang `run.bat` (migrasi + seed idempotent).
- Siapkan foto kerusakan AC yang realistis (jangan stok foto internet).
- Uji dulu seluruh alur secara live sebelum merekam — jangan improviseasi saat rekam.

### Layar & Browser
- Resolusi **1920×1080**, frame rate 30 fps (cukup) atau 60 fps bila ingin gerakan super halus.
- Zoom browser **110%** (Ctrl + `+`) agar teks chat/dashboard terbaca di layar kecil/juri.
- Sembunyikan bookmarks bar (Ctrl+Shift+B), tutup tab lain, gunakan profil browser bersih tanpa notifikasi.
- Aktifkan **Focus Assist / Do Not Disturb** Windows agar popup tidak merusak take.
- Gunakan akun demo dari README (`admin@batiqa.com`) — pastikan tidak menampilkan data sensitif.

### Rekam HP (Shot #2)
- Rekam layar HP terpisah dengan resolusi tertinggi, lalu masukkan sebagai picture-in-picture / split screen saat editing; ATAU mirror layar HP ke PC via scrcpy/Vysor lalu rekam satu software saja.
- Pastikan QR mock tercetak besar dan tajam; jarak scan ±30 cm, cahaya cukup.

### OBS Studio
- Scene 1: *Display Capture* (alur desktop); Scene 2: *Window Capture* browser (fallback crop lebih rapi).
- Settings → Output: MP4/MKV (MKV aman dari crash, remux ke MP4 setelahnya), encoder x264, bitrate 8000–12000 Kbps, keyframe 2s.
- Aktifkan filter *Color Correction* ringan + highlight kursor agar aksi klik terlihat jelas.
- Hotkey Start/Stop recording supaya bisa ambil beberapa take dan pilih yang terbaik.

### Audio & Ekspor
- Mic dinamis/condenser dekat mulut, ruangan tanpa gema; rekam suara di ruang kecil berisi barang lunak.
- Target loudness sekitar −14 LUFS, noise gate halus, potong "emm" dan jeda panjang.
- Ekspor final: H.264 MP4, 1080p, audio AAC 320 kbps; tonton ulang full sekali sebelum submit.
