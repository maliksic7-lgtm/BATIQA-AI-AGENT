# BATIQA Hotel Pekanbaru — Human Verification Checklist

> **Status:** QUESTIONNAIRE — Belum ada insert ke DB. Untuk verifikasi lapangan dengan GM / Front Office / Housekeeping / Engineering / IT / F&B sebelum data masuk `hotel_information` / `ai/intent.go`.
> **Dasar:** `docs/BATIQA_REAL_WORLD_AUDIT.md` (2026-08-21) — semua `UNKNOWN` / `NEEDS HUMAN VERIFICATION` / `CONFLICT` diubah menjadi pertanyaan.
> **Cara pakai:** Print / share ke GM Alpino Indra Putra. Jawab `Answer` kolom, tandai `Verification Status` = `VERIFIED`. Kosongkan jika `NOT APPLICABLE`. Tandai `CONFIDENTIAL` untuk internal yang tidak boleh ke Guest AI.
> **Prinsip:** Jangan mengarang SOP hotel lain. Pisahkan `Guest-facing` vs `Internal operational` (lihat SOP & Emergency). Tandai `Visibility` = `PUBLIC` / `GUEST-SAFE` / `STAFF-ONLY` / `INTERNAL`.

---

## Instructions

1. **Isi `Answer` singkat** (Short text / Choice / Time / Number). Jika tidak tahu, tulis `UNKNOWN` — jangan tebak.
2. **Department** menunjukkan siapa yang paling tahu (GM, Front Office, Housekeeping, Engineering, IT, F&B, HR/Safety).
3. **Why we need this information** menjelaskan dampak ke **AI jawaban** atau **Ticket Router**.
4. **Visibility:** `PUBLIC` = boleh di website, `GUEST-SAFE` = boleh AI jawab ke guest, `STAFF-ONLY` = hanya dashboard staff, `INTERNAL` = tidak masuk DB AI (hanya SOP internal).
5. **Prioritas:** Kerjakan dulu **Critical Questions** (§ Critical) — 15 pertanyaan agar AI realistis.

---

## A. HOTEL PROFILE

**ID: HOTEL-001**
- **Department:** GM / Front Office
- **Question:** "Apakah alamat resmi yang harus AI sampaikan ke guest adalah `Jl. Jendral Sudirman No. 17, Simpang Tiga, Pekanbaru, Riau 28288`? Apakah ada alamat alternatif untuk map?"
- **Why we need this information:** AI sering ditanya lokasi & arah; harus konsisten.
- **Answer type:** Short text
- **Current source:** Official BATIQA HIGH (Jl. Jendral Sudirman no. 17) + JTB 28288 MEDIUM
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: HOTEL-002**
- **Department:** GM
- **Question:** "Konfirmasi jarak resmi ke Sultan Syarif Kasim II International Airport untuk AI: apakah `1,7 km` (Direction) atau `2,3 km / 5 menit` (Welcome) yang harus digunakan? Atau `1,7 km straight / 2,3 km driving`?"
- **Why we need this information:** Mempengaruhi jawaban AI `Direction to Hotel` — conflict di audit.
- **Answer type:** Single choice (1,7 / 2,3 / both with note)
- **Current source:** CONFLICT DETECTED (Official 1,7 vs 2,3)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: HOTEL-003**
- **Department:** GM / Marketing
- **Question:** "Tagline resmi yang boleh AI pakai: `Discover A Unique Indonesia Experience`? Apakah ada deskripsi hotel 1 kalimat yang disetujui untuk AI greeting?"
- **Why we need this information:** Hospitality personality AI.
- **Answer type:** Short text
- **Current source:** Official HIGH
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** PUBLIC

---

## B. ROOM & ACCOMMODATION

**ID: ROOM-001**
- **Department:** Front Office / GM
- **Question:** "Konfirmasi jumlah kamar aktual 2026: `132 kamar (125 Superior + 7 Suites)` atau `133 kamar (126 + 7)`? Mana yang benar untuk AI jawab `Berapa jumlah kamar?`"
- **Why we need this information:** Mempengaruhi AI `ROOM_INFORMATION` — conflict.
- **Answer type:** Single choice (132 / 133)
- **Current source:** CONFLICT (Official 125+7 vs Bisnis.com 126+7)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: ROOM-002**
- **Department:** Front Office
- **Question:** "Luas kamar Superior yang resmi: `22,5 m²` atau `22 m²`? Dan Suite `45 m²` atau `44 m²`?"
- **Why we need this information:** AI `Room Facilities`.
- **Answer type:** Short text
- **Current source:** CONFLICT (22,5 vs 22, 45 vs 44)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: ROOM-003**
- **Department:** Front Office
- **Question:** "Tipe kamar yang dijual saat ini: apakah hanya `Superior Double (1 double bed)`, `Superior Twin (2 single bed)`, `Suite (1 double)`? Atau ada varian `Superior Double include Dinner`?"
- **Why we need this information:** AI harus tidak invent tipe kamar.
- **Answer type:** Multiple choice
- **Current source:** Trip.com/Agoda MEDIUM (Superior Double/Twin, Suite)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: ROOM-004**
- **Department:** Front Office / Engineering
- **Question:** "Fasilitas kamar yang pasti ada di SEMUA kamar: apakah `AC, LCD/plasma TV internasional, WiFi, minibar, kulkas, brankas, desk, telepon, kettle, amenities, shower, slippers, hot water 24h, daily housekeeping, ada jendela` semua benar untuk semua tipe?"
- **Why we need this information:** AI `ROOM_INFORMATION` tidak boleh invent fasilitas.
- **Answer type:** Checklist (Yes/No per item)
- **Current source:** JTB 29 item HIGH (tapi EXTERNAL untuk detail)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: ROOM-005**
- **Department:** Front Office
- **Question:** "Apakah Superior di lantai 3-10 benar? Suite di lantai berapa? Apakah 7 Suites benar 1 per lantai nomor 1 (infoPKU 9 Suites)?"
- **Why we need this information:** AI tidak invent nomor lantai/kamar.
- **Answer type:** Short text
- **Current source:** Trip.com 3-10 LOW + infoPKU 9 Suites CONFLICT
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE (lantai boleh), INTERNAL (nomor kamar detail jangan invent)

---

## C. CHECK-IN / CHECK-OUT

**ID: CHECK-001**
- **Department:** Front Office
- **Question:** "Jam resmi check-in dan check-out yang boleh AI sampaikan: `Check-in 14:00` dan `Check-out 12:00`?"
- **Why we need this information:** AI `CHECKIN_INFORMATION`/`CHECKOUT_INFORMATION` paling sering ditanya.
- **Answer type:** Time
- **Current source:** Agoda MEDIUM (14:00/12:00) — belum Official.
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: CHECK-002**
- **Department:** Front Office
- **Question:** "Guest-facing: Apakah early check-in / late check-out tergantung ketersediaan tanpa biaya, atau ada biaya? Tolong 1 kalimat yang boleh AI jawab."
- **Why we need this information:** Mempengaruhi AI jawaban, tidak boleh invent biaya.
- **Answer type:** Short text
- **Current source:** Agoda `early/late depending availability` LOW
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: CHECK-003**
- **Department:** Front Office
- **Question:** "Internal: Berapa SLA normal proses check-in/out dan escalation jika kamar belum siap? (Staff-only, tidak untuk Guest AI)"
- **Why we need this information:** Untuk `Ticket` jika ada `CHECKIN_PROBLEM` routing.
- **Answer type:** Short text
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** STAFF-ONLY

---

## D. BREAKFAST & F&B

**ID: BREAK-001**
- **Department:** F&B / Front Office
- **Question:** "Breakfast buffet jam resmi: `06:00-10:00` di `FRESQA Bistro indoor & outdoor`?"
- **Why we need this information:** AI `BREAKFAST_INFORMATION` top question.
- **Answer type:** Time + location
- **Current source:** JTB HIGH (06:00-10:00) + Official `Restaurant buka 06:00-22:00`
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: BREAK-002**
- **Department:** F&B
- **Question:** "Guest-facing: Apakah breakfast berbayar atau termasuk paket? Kalimat resmi untuk AI: `Breakfast buffet 06:00-10:00, berbayar; termasuk di paket Kartini Staycation Deal, hubungi Front Office untuk promo.` apakah benar?"
- **Why we need this information:** Conflict `paid` vs `include`.
- **Answer type:** Short text
- **Current source:** CONFLICT (JTB paid vs RiauAktual include)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: BREAK-003**
- **Department:** F&B
- **Question:** "Guest-facing: Apakah 7 signature menu 11 Aug 2026 masih tersedia harian ala carte? Konfirmasi 4 makan (Arsik Baramundi, Dori Dabu-Dabu, Tahu Cabe Garam, Otak-Otak Kenangan) dan 3 minum (Kunyit Asam Squash, Lembayung Rasa, Kopiqawan) dengan harga 55-60k/45-50k masih berlaku?"
- **Why we need this information:** AI `RESTAURANT_INFORMATION` / `FRESQA_MENU`.
- **Answer type:** Yes/No + price
- **Current source:** TribunPekanbaru/RiauPos MEDIUM (GM Alpino)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: BREAK-004**
- **Department:** F&B
- **Question:** "Internal: Berapa harga breakfast walk-in terbaru yang boleh AI sebut? Jika tidak boleh sebut harga, tulis `UNKNOWN` agar AI jawab `Hubungi Front Office`."
- **Why we need this information:** Jangan invent harga.
- **Answer type:** Number / UNKNOWN
- **Current source:** UNKNOWN (JTB 110k paid, tapi LOW)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** STAFF-ONLY (jika confidential)

---

## E. HOTEL FACILITIES

**ID: FAC-001**
- **Department:** Front Office / GM
- **Question:** "Konfirmasi fasilitas yang boleh AI sebut: `FRESQA Bistro indoor & outdoor, Lounge, 4 modular meeting rooms lantai 2, Gym, WiFi high-speed, TV Kabel` — apakah ada yang sudah tidak operasional?"
- **Why we need this information:** AI `FACILITY_INFORMATION`.
- **Answer type:** Checklist
- **Current source:** Official HIGH
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: FAC-002**
- **Department:** GM
- **Question:** "Guest-facing: Apakah `Bar` (stylish/cozy) dan `Sauna` masih ada? `Rooftop lounge` city views apakah ada (Agoda sebut, tapi JTB tidak)?"
- **Why we need this information:** Mempengaruhi `FACILITY_INFORMATION` — jika tidak ada, AI harus jawab `Maaf, belum ada` bukan invent.
- **Answer type:** Yes/No
- **Current source:** Agoda LOW (bar/sauna/rooftop) — NEEDS VERIFICATION
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: FAC-003**
- **Department:** Front Office
- **Question:** "Guest-facing: Apakah `Mushalla` lantai 2 bersama gym & meeting masih ada?"
- **Why we need this information:** Banyak guest tanya mushalla.
- **Answer type:** Yes/No + location
- **Current source:** infoPKU LOW
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: FAC-004**
- **Department:** Front Office
- **Question:** "Internal: Jam operasional masing-masing fasilitas (Gym, FRESQA, Lounge) yang resmi — jika boleh share ke guest, tulis jam; jika tidak, tulis `STAFF-ONLY`."
- **Why we need this information:** AI harus jawab jam akurat, tidak karang.
- **Answer type:** Time per facility
- **Current source:** UNKNOWN (hanya breakfast 06:00-10:00 terverifikasi)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE (jika diizinkan) / STAFF-ONLY

---

## F. WIFI & TECHNOLOGY

**ID: WIFI-001**
- **Department:** IT / Front Office
- **Question:** "Apa nama SSID WiFi yang harus AI instruksikan ke guest? (contoh `BATIQA HOTELS`?)"
- **Why we need this information:** AI `WIFI_INFORMATION` top question.
- **Answer type:** Short text
- **Current source:** UNKNOWN (hanya `Free WiFi throughout` HIGH)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: WIFI-002**
- **Department:** IT / Front Office
- **Question:** "Apakah WiFi gratis tanpa password, atau perlu password di kartu kamar / Front Office? Apakah ada SSID terpisah untuk staff?"
- **Why we need this information:** AI harus beri instruksi akurat.
- **Answer type:** Single choice
- **Current source:** UNKNOWN (JTB bilang complimentary, Agoda bilang WiFi, tapi password UNKNOWN)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: WIFI-003**
- **Department:** IT
- **Question:** "Internal: Berapa kecepatan WiFi dan apakah ada troubleshooting standar (restart, hubungi IT) yang boleh AI sampaikan sebagai `WIFI_PROBLEM` routing ke Engineering vs Front Office?"
- **Why we need this information:** Mempengaruhi `WIFI_PROBLEM` priority & routing.
- **Answer type:** Short text
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL (speed) / GUEST-SAFE (troubleshooting)

---

## G. HOUSEKEEPING

**ID: HOUSE-001**
- **Department:** Housekeeping
- **Question:** "Guest-facing: Daftar permintaan yang sering ditangani Housekeeping di Pekanbaru: handuk, bantal, selimut, amenities, pembersihan harian, laundry/dry cleaning — mana yang benar untuk AI tawarkan? (centang yang ada)"
- **Why we need this information:** Mempengaruhi `TOWEL_REQUEST`, `AMENITY_REQUEST`, `ROOM_CLEANING_REQUEST` routing.
- **Answer type:** Multiple choice
- **Current source:** JTB + Trip.com HIGH (daily, laundry) — MEDIUM untuk detail item
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: HOUSE-002**
- **Department:** Housekeeping
- **Question:** "Guest-facing procedure: Jika guest minta `Tambah 2 handuk kamar 305`, bagaimana workflow normal dari request sampai selesai yang boleh AI jelaskan ke guest? (1 kalimat per step)"
- **Why we need this information:** Untuk `Ticket` `HOUSEKEEPING` `MEDIUM` dan `Guest notification` (Open→In Progress→Resolved).
- **Answer type:** Long text (step)
- **Current source:** UNKNOWN SOP
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE (step yang boleh share) / INTERNAL (detail staff)

**ID: HOUSE-003**
- **Department:** Housekeeping
- **Question:** "Internal: Berapa normal response time / SLA untuk towel (MEDIUM) vs pillow (LOW) vs cleaning (MEDIUM) yang boleh di-share ke Staff Dashboard? Jika confidential, tulis `CONFIDENTIAL`."
- **Why we need this information:** Mempengaruhi `Priority Classification` dan `Staff Dashboard` SLA.
- **Answer type:** Time per priority
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** STAFF-ONLY (jika tidak boleh ke guest)

**ID: HOUSE-004**
- **Department:** Housekeeping
- **Question:** "Internal: Kondisi eskalasi — kapan Housekeeping harus eskalasi ke Front Office atau Engineering? Contoh: handuk habis stok, kamar belum siap."
- **Why we need this information:** Mempengaruhi `Ticket routing` & `CANCELLED` workflow.
- **Answer type:** Short text
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL

**ID: HOUSE-005**
- **Department:** Housekeeping
- **Question:** "Guest-facing: Apakah laundry/dry cleaning berbayar? Kalimat resmi AI jika ditanya harga laundry?"
- **Why we need this information:** Jangan invent biaya.
- **Answer type:** Short text (price or `Hubungi Front Office`)
- **Current source:** Trip.com `biaya tambahan` LOW
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE (jika boleh sebut) / STAFF-ONLY

---

## H. ENGINEERING / MAINTENANCE

**ID: ENG-001**
- **Department:** Engineering
- **Question:** "Guest-facing: Daftar masalah yang ditangani Engineering di Pekanbaru: AC tidak dingin, TV mati, WiFi bermasalah, lampu mati, shower bocor, plumbing, peralatan kamar (brankas, kettle) — mana yang benar?"
- **Why we need this information:** Mempengaruhi `AC_PROBLEM`, `TV_PROBLEM`, `WIFI_PROBLEM`, `LIGHT_PROBLEM`, `SHOWER_PROBLEM`, `PLUMBING_PROBLEM`, `ROOM_EQUIPMENT_PROBLEM` routing `ENGINEERING`.
- **Answer type:** Multiple choice
- **Current source:** Generic (tidak ada publik SOP) — UNKNOWN untuk Pekanbaru spesifik
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: ENG-002**
- **Department:** Engineering
- **Question:** "Guest-facing procedure: Jika guest lapor `AC kamar 305 tidak dingin`, bagaimana workflow normal yang boleh AI jelaskan? (mis. `Laporan → Engineering cek 15 menit → Resolved → Guest notifikasi`)"
- **Why we need this information:** Untuk `HIGH` priority AC & guest notification.
- **Answer type:** Long text (step)
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE (step share) / INTERNAL (teknis)

**ID: ENG-003**
- **Department:** Engineering
- **Question:** "Internal: Berapa SLA untuk AC mati total (HIGH) vs bocor (HIGH) vs TV minor (MEDIUM) yang boleh tampil di Staff Dashboard? Jika confidential, tulis `CONFIDENTIAL`."
- **Why we need this information:** Mempengaruhi `Priority Classification` dan `Staff Dashboard` `High Priority` color.
- **Answer type:** Time per category
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** STAFF-ONLY

**ID: ENG-004**
- **Department:** Engineering / Front Office
- **Question:** "Jika WiFi bermasalah, apakah routing ke `ENGINEERING` atau `FRONT_OFFICE` (password)? Kapan ke masing-masing?"
- **Why we need this information:** Mempengaruhi `WIFI_PROBLEM` department routing — conflict.
- **Answer type:** Single choice
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL

---

## I. FRONT OFFICE

**ID: FRONT-001**
- **Department:** Front Office
- **Question:** "Guest-facing: Nomor Front Office yang boleh AI berikan jika info tidak ada di DB (fallback `Hubungi Front Office`): apakah `+62 761 889000` atau extension internal?"
- **Why we need this information:** AI fallback `Maaf, hubungi Front Office` harus ada nomor.
- **Answer type:** Short text (phone)
- **Current source:** Official `+62 761 889000` HIGH
- **Verification status:** NEEDS HUMAN VERIFICATION (apakah ada extension internal 0?)
- **Visibility:** GUEST-SAFE

**ID: FRONT-002**
- **Department:** Front Office
- **Question:** "Guest-facing: Apakah `early check-in / late check-out` tergantung ketersediaan tanpa biaya, atau ada biaya? Kalimat resmi untuk AI."
- **Why we need this information:** Sering ditanya, jangan invent biaya.
- **Answer type:** Short text
- **Current source:** Agoda LOW
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: FRONT-003**
- **Department:** Front Office
- **Question:** "Internal: Berapa lama proses check-in/out normal dan kapan harus buat `Ticket` vs langsung jawab AI? (mis. check-in 14:00 tapi kamar belum siap → buat `ROOM_CLEANING_REQUEST`?)"
- **Why we need this information:** Mempengaruhi `Action Decision` AI.
- **Answer type:** Short text
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL

---

## J. GUEST REQUEST

**ID: GUEST-REQ-001**
- **Department:** Front Office / Housekeeping
- **Question:** "Guest-facing: Daftar permintaan yang sering di Pekanbaru (handuk, bantal, pembersihan, amenities, room service) — apakah ada yang tidak dilayani?"
- **Why we need this information:** AI `HOUSEKEEPING_REQUEST` coverage.
- **Answer type:** Multiple choice
- **Current source:** Generic UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: GUEST-REQ-002**
- **Department:** Front Office
- **Question:** "Internal: Apakah ada batasan `room_number` yang valid (mis. 101-410, lantai 3-10)? Format `305` atau `305A`?"
- **Why we need this information:** Mempengaruhi `Entity Extraction` `room_number` validation regex.
- **Answer type:** Short text (regex)
- **Current source:** UNKNOWN (hanya 305 contoh)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL (jangan share semua nomor ke guest)

---

## K. ROOM SERVICE

**ID: ROOM-SRV-001**
- **Department:** F&B / Front Office
- **Question:** "Guest-facing: Apakah room service benar 24 jam? Apakah menu sama dengan FRESQA Bistro atau terpisah? Jam last order?"
- **Why we need this information:** AI `ROOM_SERVICE` — sering ditanya jam 2 pagi.
- **Answer type:** Short text
- **Current source:** Official HIGH (24 jam) + JTB 24h
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: ROOM-SRV-002**
- **Department:** F&B
- **Question:** "Guest-facing: Jika guest pesan room service `Saya mau pesan makan`, apakah AI harus buat `ROOM_SERVICE_REQUEST` → `FRESQA` / `Housekeeping` / `Front Office`? Department mana?"
- **Why we need this information:** Mempengaruhi `Ticket Routing` untuk Pekanbaru.
- **Answer type:** Single choice
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL

---

## L. COMPLAINT HANDLING

**ID: COMPLAINT-001**
- **Department:** Front Office / GM
- **Question:** "Guest-facing: Jika guest komplain `Kamar berisik`, `AC bocor`, bagaimana kalimat empati AI yang disetujui? (mis. `Maaf atas ketidaknyamanan... akan kami teruskan ke Engineering`)"
- **Why we need this information:** Mempengaruhi `AI Response` hospitality.
- **Answer type:** Short text (template)
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: COMPLAINT-002**
- **Department:** Front Office
- **Question:** "Internal: Kapan komplain harus jadi `HIGH` vs `MEDIUM`? Contoh: `AC mati total` HIGH, `TV minor` MEDIUM — apakah ada daftar Pekanbaru spesifik?"
- **Why we need this information:** Mempengaruhi `Priority Classification` (jangan exaggerate).
- **Answer type:** Short text
- **Current source:** Generic `PRIORITY CLASSIFICATION.md` — perlu Pekanbaru spesifik
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** STAFF-ONLY

---

## M. EMERGENCY & SAFETY

**ID: EMERG-001**
- **Department:** Safety / Front Office / GM
- **Question:** "Guest-facing: Nomor emergency yang BOLEH AI berikan ke guest (mis. Front Office `+62 761 889000`, atau `112`)? Apakah boleh sebut `extension`?"
- **Why we need this information:** Safety — jangan karang nomor.
- **Answer type:** Short text (phone or `Hubungi Front Office`)
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE (hanya yang diizinkan)

**ID: EMERG-002**
- **Department:** Safety
- **Question:** "Guest-facing: Titik kumpul (assembly point) yang BOLEH AI informasikan jika guest tanya `Jika ada kebakaran`? Jika tidak boleh, tulis `STAFF-ONLY, AI jawab Hubungi Front Office`."
- **Why we need this information:** Jangan karang evacuation.
- **Answer type:** Short text or `STAFF-ONLY`
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE (jika diizinkan) / STAFF-ONLY

**ID: EMERG-003**
- **Department:** Safety
- **Question:** "Guest-facing: Instruksi evakuasi singkat yang BOLEH AI sampaikan (mis. `Ikuti petunjuk staff, jangan gunakan lift`)? Tulis kalimat resmi atau `TIDAK BOLEH AI SAMPAIKAN, hanya staff`."
- **Why we need this information:** Jangan karang procedure.
- **Answer type:** Short text or `STAFF-ONLY`
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE / STAFF-ONLY

**ID: EMERG-004**
- **Department:** Safety / GM
- **Question:** "Internal: Siapa yang menangani emergency (Front Office, Security, GM)? Apakah AI harus buat `HIGH` ticket `SAFETY` ke `ENGINEERING`/`FRONT_OFFICE` atau langsung `Hubungi Front Office` tanpa ticket?"
- **Why we need this information:** Mempengaruhi `Ticket Routing` untuk safety (jangan salah).
- **Answer type:** Short text
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL

**ID: EMERG-005**
- **Department:** Safety
- **Question:** "Daftar informasi yang AI **TIDAK BOLEH** sampaikan ke guest (mis. detail APAR, jalur pipa, password WiFi staff, SOP internal) — tolong list."
- **Why we need this information:** Data Governance, `AI SAFETY RULES` tidak expose private.
- **Answer type:** Long text (blacklist)
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL

---

## N. PARKING & TRANSPORTATION

**ID: PARK-001**
- **Department:** Front Office / Security
- **Question:** "Guest-facing: Apakah parkir benar gratis, private, tidak perlu reservasi, multiple car parks? Kalimat resmi untuk AI."
- **Why we need this information:** AI `PARKING_INFORMATION`.
- **Answer type:** Short text
- **Current source:** Trip.com HIGH (gratis) — perlu konfirmasi
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: PARK-002**
- **Department:** Front Office
- **Question:** "Guest-facing: Apakah EV charging station gratis, berapa jumlah/tipe, perlu reservasi?"
- **Why we need this information:** AI `PARKING_INFORMATION` untuk tamu EV.
- **Answer type:** Short text
- **Current source:** Trip.com LOW
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: TRANS-001**
- **Department:** Front Office
- **Question:** "Guest-facing: Apakah airport shuttle benar gratis 24 jam? Perlu booking berapa jam sebelum? Kalimat resmi untuk AI `Gratis 24 jam, hubungi Front Office untuk booking`?"
- **Why we need this information:** AI `TRANSFER_INFORMATION` / `AIRPORT_SHUTTLE`.
- **Answer type:** Short text
- **Current source:** JTB HIGH (gratis 24h) — perlu booking detail
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: TRANS-002**
- **Department:** Front Office
- **Question:** "Internal: Berapa kapasitas shuttle, berapa lama `2,3 km / 5 menit` masih akurat 2026?"
- **Why we need this information:** Untuk `Direction` akurat.
- **Answer type:** Short text
- **Current source:** CONFLICT 1,7 vs 2,3 km
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL

---

## O. MEETING & EVENT

**ID: MEET-001**
- **Department:** Sales / Front Office
- **Question:** "Guest-facing: Apakah 4 modular meeting rooms lantai 2 masih akurat? Kapasitas per room dan apakah `lounge` termasuk?"
- **Why we need this information:** AI `MEETING_INFORMATION`.
- **Answer type:** Short text
- **Current source:** Official HIGH (4 modular)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: MEET-002**
- **Department:** Sales
- **Question:** "Guest-facing: Apakah meeting/banquet fasilitas berbayar? Kalimat resmi untuk AI jika ditanya harga meeting."
- **Why we need this information:** Jangan invent biaya.
- **Answer type:** Short text (`Hubungi Front Office` or price)
- **Current source:** Trip.com `biaya tambahan` LOW
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE (jika boleh) / STAFF-ONLY

---

## P. LOCAL RECOMMENDATION

**ID: LOCAL-001**
- **Department:** Front Office / Concierge
- **Question:** "Guest-facing: Apakah daftar rekomendasi Agoda (Warung Pempek Sikok Nak Duo, RM Puti Buana, Chilli Crab ID, The Baliview, Durian Runtuh, Pasar Wisata, Mall Pekanbaru, Ar-Rahman Mosque, Silungkang Art Centre) adalah rekomendasi resmi hotel yang boleh AI sebut? Atau ada daftar resmi lain?"
- **Why we need this information:** AI `RESTAURANT_RECOMMENDATION` tidak boleh invent resto.
- **Answer type:** Multiple choice (endorse list / ada list lain / UNKNOWN)
- **Current source:** Agoda EXTERNAL LOW — **NEEDS VERIFICATION** apakah endorsed
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: LOCAL-002**
- **Department:** Front Office
- **Question:** "Guest-facing: Untuk rekomendasi berdasarkan budget (`Saya mau makan budget 100 ribu`), apakah ada daftar harga/distance yang boleh AI pakai, atau harus jawab `Hubungi Front Office`?"
- **Why we need this information:** Mempengaruhi `RECOMMENDATION RULES` (budget, distance).
- **Answer type:** Short text
- **Current source:** UNKNOWN (Agoda list tanpa price/distance presisi)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE (jika ada) / STAFF-ONLY

---

## Q. ACCESSIBILITY

**ID: ACCESS-001**
- **Department:** Front Office / Engineering
- **Question:** "Guest-facing: Apakah hotel memiliki kamar difabel, ramp, toilet difabel, atau akses kursi roda yang boleh AI informasikan? Jika tidak ada, tulis `TIDAK ADA` agar AI jawab `Maaf, belum ada`."
- **Why we need this information:** AI `ACCESSIBILITY` — jangan invent.
- **Answer type:** Short text
- **Current source:** UNKNOWN (hanya elevator terverifikasi)
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

**ID: ACCESS-002**
- **Department:** Front Office
- **Question:** "Guest-facing: Apakah ada designated smoking area? Dimana? Apakah ada kamar non-smoking varian?"
- **Why we need this information:** AI `HOTEL_POLICY` smoking.
- **Answer type:** Short text
- **Current source:** JTB `designated smoking area` LOW
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** GUEST-SAFE

---

## R. AI / DIGITAL GUEST SERVICE

**ID: AI-001**
- **Department:** GM / Front Office / IT
- **Question:** "Guest-facing: Apakah tamu boleh scan QR di kamar untuk akses AI tanpa install app? Apakah QR sudah terpasang di semua 132 kamar?"
- **Why we need this information:** Core `Zero-Install QR Access` UX.
- **Answer type:** Yes/No
- **Current source:** README.md (konsep) — belum verifikasi fisik QR
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** PUBLIC

**ID: AI-002**
- **Department:** IT / Front Office
- **Question:** "Guest-facing: Bahasa apa yang harus AI dukung selain `id` dan `en`? Apakah perlu Melayu lokal?"
- **Why we need this information:** `LANGUAGE BEHAVIOR` id/en preserve.
- **Answer type:** Multiple choice
- **Current source:** README `Multi-language` id/en HIGH
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** PUBLIC

**ID: AI-003**
- **Department:** GM / Front Office
- **Question:** "Internal: Informasi apa yang **TIDAK BOLEH** AI sampaikan ke guest (mis. detail SOP internal, password staff, nomor kamar tamu lain, harga internal)? Tolong list blacklist."
- **Why we need this information:** `AI SAFETY RULES` tidak expose private.
- **Answer type:** Long text
- **Current source:** UNKNOWN
- **Verification status:** NEEDS HUMAN VERIFICATION
- **Visibility:** INTERNAL

**ID: AI-004**
- **Department:** Front Office
- **Question:** "Internal: Jika AI tidak tahu jawaban (UNKNOWN), apakah AI harus `Clarification` seperti `Apakah Anda ingin info hotel, request, atau lapor masalah?` atau langsung `Hubungi Front Office`?"
- **Why we need this information:** `UNKNOWN INTENT` handling.
- **Answer type:** Single choice
- **Current source:** `UNKNOWN INTENT.md` HIGH (clarification)
- **Verification status:** NEEDS HUMAN VERIFICATION (apakah sudah sesuai harapan GM?)
- **Visibility:** INTERNAL

**ID: AI-005**
- **Department:** IT
- **Question:** "Internal: Apakah AI boleh sebut `Staff telah menyelesaikan request Anda` hanya jika backend `RESOLVED` `true`, atau boleh AI karang? (Security: AI tidak boleh claim ticket selesai jika backend belum konfirmasi)"
- **Why we need this information:** `AI SAFETY RULES` — jangan claim palsu.
- **Answer type:** Yes/No
- **Current source:** `AI SAFETY RULES.md` HIGH
- **Verification status:** NEEDS HUMAN VERIFICATION (apakah aturan ini sudah dipahami staff?)
- **Visibility:** INTERNAL

---

## Verification Summary

| ID | Category | Department | Verification Status | Answer | Visibility |
|---|---|---|---|---|---|
| HOTEL-001 | Hotel Profile | GM / Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| HOTEL-002 | Hotel Profile | GM | NEEDS VERIFICATION |  | GUEST-SAFE |
| HOTEL-003 | Hotel Profile | GM / Marketing | NEEDS VERIFICATION |  | PUBLIC |
| ROOM-001 | Room & Accommodation | Front Office / GM | NEEDS VERIFICATION |  | GUEST-SAFE |
| ROOM-002 | Room & Accommodation | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| ROOM-003 | Room & Accommodation | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| ROOM-004 | Room & Accommodation | Front Office / Engineering | NEEDS VERIFICATION |  | GUEST-SAFE |
| ROOM-005 | Room & Accommodation | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE / INTERNAL |
| CHECK-001 | Check-in/out | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| CHECK-002 | Check-in/out | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| CHECK-003 | Check-in/out | Front Office | NEEDS VERIFICATION |  | STAFF-ONLY |
| BREAK-001 | Breakfast & F&B | F&B / Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| BREAK-002 | Breakfast & F&B | F&B | NEEDS VERIFICATION |  | GUEST-SAFE |
| BREAK-003 | Breakfast & F&B | F&B | NEEDS VERIFICATION |  | GUEST-SAFE |
| BREAK-004 | Breakfast & F&B | F&B | NEEDS VERIFICATION |  | STAFF-ONLY |
| FAC-001 | Hotel Facilities | Front Office / GM | NEEDS VERIFICATION |  | GUEST-SAFE |
| FAC-002 | Hotel Facilities | GM | NEEDS VERIFICATION |  | GUEST-SAFE |
| FAC-003 | Hotel Facilities | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| FAC-004 | Hotel Facilities | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE / STAFF-ONLY |
| WIFI-001 | WiFi | IT / Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| WIFI-002 | WiFi | IT / Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| WIFI-003 | WiFi | IT | NEEDS VERIFICATION |  | INTERNAL / GUEST-SAFE |
| HOUSE-001 | Housekeeping | Housekeeping | NEEDS VERIFICATION |  | GUEST-SAFE |
| HOUSE-002 | Housekeeping | Housekeeping | NEEDS VERIFICATION |  | GUEST-SAFE / INTERNAL |
| HOUSE-003 | Housekeeping | Housekeeping | NEEDS VERIFICATION |  | STAFF-ONLY |
| HOUSE-004 | Housekeeping | Housekeeping | NEEDS VERIFICATION |  | INTERNAL |
| HOUSE-005 | Housekeeping | Housekeeping | NEEDS VERIFICATION |  | GUEST-SAFE / STAFF-ONLY |
| ENG-001 | Engineering | Engineering | NEEDS VERIFICATION |  | GUEST-SAFE |
| ENG-002 | Engineering | Engineering | NEEDS VERIFICATION |  | GUEST-SAFE / INTERNAL |
| ENG-003 | Engineering | Engineering | NEEDS VERIFICATION |  | STAFF-ONLY |
| ENG-004 | Engineering / Front Office | Engineering | NEEDS VERIFICATION |  | INTERNAL |
| FRONT-001 | Front Office | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| FRONT-002 | Front Office | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| FRONT-003 | Front Office | Front Office | NEEDS VERIFICATION |  | INTERNAL |
| GUEST-REQ-001 | Guest Request | Front Office / Housekeeping | NEEDS VERIFICATION |  | GUEST-SAFE |
| GUEST-REQ-002 | Guest Request | Front Office | NEEDS VERIFICATION |  | INTERNAL |
| ROOM-SRV-001 | Room Service | F&B / Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| ROOM-SRV-002 | Room Service | F&B | NEEDS VERIFICATION |  | INTERNAL |
| COMPLAINT-001 | Complaint | Front Office / GM | NEEDS VERIFICATION |  | GUEST-SAFE |
| COMPLAINT-002 | Complaint | Front Office | NEEDS VERIFICATION |  | STAFF-ONLY |
| EMERG-001 | Emergency & Safety | Safety / Front Office / GM | NEEDS VERIFICATION |  | GUEST-SAFE |
| EMERG-002 | Emergency & Safety | Safety | NEEDS VERIFICATION |  | GUEST-SAFE / STAFF-ONLY |
| EMERG-003 | Emergency & Safety | Safety | NEEDS VERIFICATION |  | GUEST-SAFE / STAFF-ONLY |
| EMERG-004 | Emergency & Safety | Safety / GM | NEEDS VERIFICATION |  | INTERNAL |
| EMERG-005 | Emergency & Safety | Safety | NEEDS VERIFICATION |  | INTERNAL |
| PARK-001 | Parking & Transportation | Front Office / Security | NEEDS VERIFICATION |  | GUEST-SAFE |
| PARK-002 | Parking & Transportation | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| TRANS-001 | Parking & Transportation | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| TRANS-002 | Parking & Transportation | Front Office | NEEDS VERIFICATION |  | INTERNAL |
| MEET-001 | Meeting & Event | Sales / Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| MEET-002 | Meeting & Event | Sales | NEEDS VERIFICATION |  | GUEST-SAFE / STAFF-ONLY |
| LOCAL-001 | Local Recommendation | Front Office / Concierge | NEEDS VERIFICATION |  | GUEST-SAFE |
| LOCAL-002 | Local Recommendation | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE / STAFF-ONLY |
| ACCESS-001 | Accessibility | Front Office / Engineering | NEEDS VERIFICATION |  | GUEST-SAFE |
| ACCESS-002 | Accessibility | Front Office | NEEDS VERIFICATION |  | GUEST-SAFE |
| AI-001 | AI / Digital | GM / Front Office / IT | NEEDS VERIFICATION |  | PUBLIC |
| AI-002 | AI / Digital | IT / Front Office | NEEDS VERIFICATION |  | PUBLIC |
| AI-003 | AI / Digital | GM / Front Office | NEEDS VERIFICATION |  | INTERNAL |
| AI-004 | AI / Digital | Front Office | NEEDS VERIFICATION |  | INTERNAL |
| AI-005 | AI / Digital | IT | NEEDS VERIFICATION |  | INTERNAL |

> **Catatan:** Semua `NEEDS HUMAN VERIFICATION` di atas harus diubah menjadi `VERIFIED` setelah diisi `Answer`, atau `NOT APPLICABLE` jika tidak ada, atau `CONFIDENTIAL` jika tidak boleh share. `Visibility` menentukan apakah data boleh masuk `hotel_information` (GUEST-SAFE) atau hanya SOP internal.

---

## Critical Questions (15 Paling Penting — Kerjakan Dulu)

> **Agar AI realistis, kerjakan 15 ini dulu. Tanpa ini, AI akan `UNKNOWN` atau invent.**

| # | ID | Question | Why Critical |
|---|---|---|---|
| 1 | **ROOM-001** | 132 vs 133 kamar? 125 vs 126 Superior? | AI `Room Information` paling sering, conflict HIGH. |
| 2 | **CHECK-001** | Check-in 14:00 / Check-out 12:00? | Top 3 question guest. |
| 3 | **BREAK-001** | Breakfast 06:00-10:00 di FRESQA? | Top question, harus akurat. |
| 4 | **BREAK-002** | Breakfast berbayar atau termasuk? | Conflict, jangan invent harga. |
| 5 | **WIFI-001/002** | SSID & password WiFi? | Top question, AI harus instruksi tepat. |
| 6 | **HOUSE-002** | Workflow housekeeping `handuk` (request → selesai) yang boleh AI jelaskan? | Mempengaruhi `Ticket` HOUSEKEEPING & guest notification (Open→Resolved). |
| 7 | **ENG-002** | Workflow engineering `AC tidak dingin` yang boleh AI jelaskan? | Mempengaruhi `HIGH` priority & SOP. |
| 8 | **FRONT-001** | Nomor Front Office yang boleh AI berikan saat fallback? | AI fallback `Hubungi Front Office` harus ada nomor. |
| 9 | **EMERG-001** | Nomor emergency yang BOLEH AI berikan ke guest? | Safety — jangan karang. |
| 10 | **PARK-001/TRANS-001** | Parkir gratis? Shuttle gratis 24h perlu booking? | Sering ditanya, mudah invent jika tidak verifikasi. |
| 11 | **ROOM-004** | Fasilitas kamar yang pasti di semua tipe (AC, TV, WiFi, minibar, brankas)? | AI `Room Facilities` tidak boleh invent. |
| 12 | **FAC-001/002** | Fasilitas mana yang masih operasional (Bar, Sauna, Rooftop, Mushalla)? | AI `Facility Information` — jika salah, guest kecewa. |
| 13 | **HOUSE-003/ENG-003** | SLA untuk MEDIUM (handuk) vs HIGH (AC bocor) yang boleh tampil di Staff Dashboard? | Mempengaruhi `Priority Classification` operational. |
| 14 | **EMERG-005/AI-003** | Daftar blacklist informasi yang AI **TIDAK BOLEH** sampaikan? | Security, `AI SAFETY RULES`. |
| 15 | **AI-001** | Apakah QR sudah terpasang di semua 132 kamar untuk `Zero-Install`? | Core UX — jika belum, guest tidak bisa akses AI. |

---

> **Next Step:** Setelah 15 Critical diisi, audit akan di-update menjadi `VERIFIED`, lalu buat `migrations\003_batiqa_real_world.sql` **hanya untuk GUEST-SAFE/PUBLIC** (bukan INTERNAL/CONFIDENTIAL).

