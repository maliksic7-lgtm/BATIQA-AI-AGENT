# BATIQA Hotel Pekanbaru — Real World Data Audit

> **Status:** RESEARCH ONLY — Belum ada insert ke `hotel_information`, `recommendations`, atau modifikasi `migrations`. Menunggu `HUMAN REVIEW` sebelum implementasi.
> **Tanggal Audit:** 2026-08-21
> **Auditor:** AI Research Agent (Muse Spark)
> **Sumber Utama:** `https://www.batiqa.com/hotels/pekanbaru` (Official BATIQA Website) — diakses via `webfetch`/`websearch` 2026-08-21 (transport error parsial, fallback ke `websearch` highlight & Agoda/JTB/Trip.com sebagai EXTERNAL SOURCE).

---

## 1. Verified Information

| Category | Information | Source | Confidence |
|---|---|---|---|
| **1. HOTEL PROFILE** | Nama: **BATIQA Hotel Pekanbaru** — bintang 3, modern boutique, bagian dari **PT Batiqa Hotel Manajemen** (lini bisnis PT Surya Semesta Internusa Tbk). Hotel ke-5 BATIQA, properti ke-2 di Sumatera setelah Palembang. Dibuka **4 Mei 2016** (akhir medio Agustus 2016 operasional penuh per Berita Bisnis 2016-09-01). Tagline: *Discover A Unique Indonesia Experience*. | Official BATIQA `batiqa.com/hotels/pekanbaru` highlight + Bisnis.com 2016-04-06 + Berita Bisnis 2016-09-01 | HIGH |
| **2. LOCATION** | Alamat: **Jl. Jendral Sudirman No. 17, Simpang Tiga, Pekanbaru, Riau 28288** — central business district, Simpang Tiga. | Official BATIQA `hotels/pekanbaru` (Jl. Jendral Sudirman no. 17) + Agoda (28288) + JTB (Simpang Tiga) | HIGH |
| **2. LOCATION** | Kota Pekanbaru — *major economic center on eastern part of Sumatra Island* (deskripsi official). | Official BATIQA highlight | MEDIUM |
| **2. LOCATION** | Jarak Bandara: **Sultan Syarif Kasim II International Airport** — Official menyebut **1,7 km** (Direction to Hotel) dan **2,3 km / 5 menit mobil** (Welcome text). | Official BATIQA — **CONFLICT, lihat §2** | HIGH (lokasi bandara) / CONFLICT (jarak) |
| **3. CONTACT** | Telepon: **+62 761 889000** | Official BATIQA (wa.me/62761889000) + BertuahPos 2019-08-26 | HIGH |
| **3. CONTACT** | WhatsApp: **0812-68933366** (BertuahPos) / **+62 761 889000** via wa.me | BertuahPos 2019-08-26 (EXTERNAL SOURCE tapi konsisten) | MEDIUM |
| **3. CONTACT** | Email reservasi: **reservation.pekanbaru@batiqa.com** | Official BATIQA + BertuahPos | HIGH |
| **3. CONTACT** | Website: `https://www.batiqa.com/hotels/pekanbaru`, subpages `/rooms`, `/facilities`, `/dining`, `/meetings`, `/offers`, `/photos`, `/insiders-guide` | Official BATIQA highlight | HIGH |
| **3. CONTACT** | Social: FB `Batiqa_pekanbaru`, IG `Batiqa_pekanbaru` (disebut BertuahPos) | BertuahPos (EXTERNAL) | LOW |
| **4. ROOM TYPES** | Kapasitas: **132 kamar** (dominant source) terdiri **125 Superior + 7 Suites** **atau** **126 Superior + 7 Suites** — **CONFLICT** (lihat §2) | Official BATIQA article `liburan-asyik` (125+7=132) + VenueMagz 2016-08-31 (125+7, 22,5m²) + Berita Bisnis 2016-09-01 (125+7) vs Bisnis.com 2016-04-06 (133 kamar: 126+7) | HIGH (132) / CONFLICT (detail) |
| **4. ROOM TYPES** | Tipe publik: **Superior Double** (1 double bed), **Superior Twin** (2 single bed), **Suite** (1 double bed). Varian Agoda: *Superior Double/Twin, Suite, Superior Double include 1 pax Dinner* | Trip.com + Agoda (EXTERNAL SOURCE) — konsisten dengan official 126+7 | MEDIUM |
| **4. ROOM TYPES** | Suite: **7 Suites**, **45 m²** (VenueMagz, Berita Bisnis) / **45 m²** (official article) / **44 m²** (Agoda) — konflik minor 44 vs 45 | Official article + VenueMagz + Agoda (EXTERNAL) | MEDIUM (45) |
| **5. ROOM FACILITIES** | Luas: **Superior 22 m²** (Trip.com) / **22,5 m²** (VenueMagz, official article) — konflik 22 vs 22,5 | Trip.com (22) vs Official (22,5) | MEDIUM |
| **5. ROOM FACILITIES** | Fasilitas kamar (terverifikasi lintas sumber): **Air conditioning (penyejuk udara)**, **LCD/plasma TV** saluran internasional / cable channels, **Wi-Fi gratis** high-speed, **minibar + kulkas**, **safety deposit box (brankas)**, **desk + telephone**, **electric kettle + coffee/tea maker**, **bottled water gratis**, **private bathroom shower only**, **bath amenities gratis + slippers**, **hot water 24 jam**, **daily housekeeping (pembersihan setiap hari)**, **ada jendela**, **city/street view** (sebagian), **non-smoking** varian, **pemandangan kota** | Official article (`liburan-asyik`: TV Kabel, wifi, 24h room service) + JTB list (29 item) + Trip.com + Agoda | HIGH |
| **5. ROOM FACILITIES** | Lantai: **Lantai 3-10** untuk Superior (Trip.com) | Trip.com (EXTERNAL) | LOW |
| **6. HOTEL FACILITIES** | **FRESQA Bistro** indoor & outdoor (restaurant), **Lounge**, **4 modular meeting rooms** (lantai 2), **Gym / fitness center** (tempat gym), **Wi-Fi berkecepatan tinggi throughout hotel**, **TV Kabel** | Official article `liburan-asyik` + VenueMagz + Berita Bisnis | HIGH |
| **6. HOTEL FACILITIES** | **Bar** stylish/cozy (Agoda: chic bar), **Sauna** luxurious (Agoda), **Rooftop lounge** city views (Agoda EXTERNAL, butuh verifikasi) | Agoda (EXTERNAL) — **NEEDS HUMAN VERIFICATION** untuk rooftop | LOW (bar/sauna) |
| **6. HOTEL FACILITIES** | **Mushalla** di lantai 2 (bersama gym & meeting) | infoPKU 2019-07-16 (EXTERNAL) | LOW |
| **7. DINING / RESTAURANT** | **FRESQA Bistro** — 1 restaurant utama, indoor & outdoor, **specialist international cuisine**, lunch & dinner. **Coffee shop** trendy, **Bar** (Agoda) | Official + JTB (`FRESQA` lunch/dinner, multinational cuisine) + Agoda | HIGH |
| **7. DINING / RESTAURANT** | **7 Signature Menu** diluncurkan **11 Agustus 2026 serentak 7 properti BATIQA** oleh GM **Alpino Indra Putra**: Makanan (4) **Arsik Baramundi** (Sumatera Utara), **Dori Dabu-Dabu** (Manado pedas asam), **Tahu Cabe Garam** (ciri khas Pekanbaru), **Otak-Otak Kenangan**; Minuman (3) **Kunyit Asam Squash**, **Lembayung Rasa**, **Kopiqawan**. Harga makanan **Rp55.000-60.000**, minuman **Rp45.000-50.000**, tersedia harian ala carte mulai 11 Aug 2026. | TribunPekanbaru 2026-08-11 + RiauPos 2026-08-12 (EXTERNAL, tapi GM quote) | MEDIUM |
| **8. ROOM SERVICE** | **24 jam** layanan kamar | Official article `liburan-asyik` + JTB (`Room service selama 24 jam` + `24 hour front desk`) + Agoda (`24-hour room service`) | HIGH |
| **9. BREAKFAST** | **Buffet** setiap hari **06:00 – 10:00** di **FRESQA Bistro** / restaurant lantai 1. **Berbayar** (JTB: *Breakfast buffet daily 6:00-10:00 paid*, Agoda menyebut buffet). | JTB (6:00-10:00 paid) + Agoda + Official (Restaurant buka 06:00-22:00) | HIGH (jam) / MEDIUM (paid vs include) |
| **9. BREAKFAST** | Sarapan untuk 2 orang termasuk di paket **Kartini Staycation Deal** (RiauAktual 2026-04-22, Superior + smart TV + breakfast 2) | RiauAktual (EXTERNAL) | LOW |
| **10. MEETING & EVENTS** | **4 modular meeting rooms** di lantai 2, **lounge**, **meeting/banquet facilities** (biaya tambahan — Trip.com), **FRESQA Bistro** bisa untuk *Afternoon Tea Package*, *Bonding in style* | Official + VenueMagz + Trip.com (biaya tambahan) | HIGH |
| **11. AIRPORT TRANSFER** | **Gratis** airport shuttle **24 jam**, antar-jemput bandara, **multiple car parks onsite** | JTB (24h airport shuttle free) + Trip.com (`Penjemputan di bandara Gratis`) + infoPKU (`layanan antar jemput ke Bandara`) + Agoda | HIGH (gratis) |
| **12. WIFI / INTERNET** | **Free WiFi** throughout hotel, di kamar & public area, **high-speed**, complimentary | Official (Free Wifi Throughout) + JTB + Trip.com (`Wi-Fi di tempat umum`) + Agoda | HIGH |
| **13. GYM / FITNESS** | **State-of-the-art fitness center** fully equipped: cardio **treadmill, stationary bike, elliptical**, strength **free weights, weight machines**, **spacious well-lit**, **complimentary towels & water**, **staff assistance**, di lantai 2 | Agoda detail + Official (Tempat gym) + JTB | HIGH |
| **14. SPA** | **Sauna** luxurious (Agoda) dan **Spa** chic bar (Agoda) — JTB tidak list spa, hanya `Spa` di Agoda. **NEEDS HUMAN VERIFICATION** apakah spa/sauna masih operasional. | Agoda (EXTERNAL) | LOW |
| **15. CHECK-IN / CHECK-OUT** | **Check-in 14:00 (02:00 PM)**, **Check-out 12:00 PM** — generous, early/late tergantung ketersediaan. Disebut juga **24 hours check-in** untuk varian Suite di Agoda (mungkin promo, bukan SOP). | Agoda (EXTERNAL, konsisten lintas listing) | MEDIUM |
| **16. PARKING** | **Ample multiple car parks onsite**, **complimentary / gratis**, **private exclusive for guests**, **tidak perlu reservasi**, **EV charging station** (Trip.com) | Trip.com (Area Parkir Pribadi) + Agoda (car parking facilities) + JTB (`Parking (Free)`) | HIGH (gratis) |
| **17. HOUSEKEEPING** | **Daily housekeeping** (pembersihan setiap hari), **laundry / dry cleaning** (biaya tambahan per Trip.com), | JTB + Trip.com + Agoda | HIGH |
| **18. ENGINEERING / MAINTENANCE** | Tidak ada publik SOP/teknisi per departemen — **UNKNOWN** (hanya fasilitas AC, listrik, plumbing implisit) | — | — |
| **19. FRONT OFFICE** | **Resepsionis 24 jam** (reception hall), **front desk 24h**, **bell staff (baggage handling / luggage storage)**, **wake-up service (layanan bangun tidur)**, **elevator**, **safety deposit box (front desk)** | JTB (24h front desk) + Trip.com (Resepsionis 24 jam) + Agoda | HIGH |
| **20. GUEST REQUEST** | **Contoh publik:** handuk, bantal, pembersihan, amenities — tidak ada daftar SOP publik. `UNKNOWN` untuk SLA/prosedur. | — | — |
| **21. HOTEL SERVICES** | **Layanan:** laundry/dry cleaning, luggage storage, safety deposit box, smoking area (designated), business center (implisit), convenience facilities | Trip.com + Agoda + JTB | MEDIUM |
| **22. LOCAL RECOMMENDATIONS** | **Sekitar hotel (Simpang Tiga, central business district):** shopping malls, culinary centers, karaoke bars, local markets (Official), **Mall Pekanbaru**, **Pasar Wisata** (sate Pekanbaru, gulai kambing), **Ar-Rahman Mosque** walking distance, **Silungkang Art Centre** short drive, **Taman Wisata?** — **Gastronomic:** Warung Pempek Sikok Nak Duo, RM Puti Buana, Chilli Crab ID, The Baliview Luxury Villas, Durian Runtuh, Jimbaran Pool Resto & BBQ, Gubug Jowo, Sate Taichan | Official + Agoda + Trip.com | MEDIUM (lokasi) / LOW (jarak presisi) |
| **23. EMERGENCY / SAFETY INFORMATION** | **Fire safety, emergency exit, procedure** — **UNKNOWN** (tidak ada publik). Hanya **safety deposit box** dan **smoking area** yang terverifikasi. | — | — |
| **24. ACCESSIBILITY** | **Elevator**, **designated smoking area**, **non-smoking room varian** — tidak ada detail kursi roda, ramp, difabel. **UNKNOWN** untuk aksesibilitas lengkap. | JTB (Elevator) + Trip.com | LOW |
| **25. OTHER INFORMATION** | **Child policy:** Anak **4-10 tahun gratis** (Agoda) — **EXTERNAL SOURCE, needs verification** apakah masih berlaku. **Harga kamar** lebih terjangkau dari rata-rata kota ($20) — tidak relevan untuk AI. **Rating** 8.x — tidak relevan. | Agoda (EXTERNAL) | LOW |

---

## 2. Conflicting Information

### CONFLICT #1 — Jumlah Kamar

**Source A (HIGH confidence - Official BATIQA + 2 berita):**
> 132 kamar — terdiri 125 Superior + 7 Suites, 22,5 m² & 45 m²
- `https://www.batiqa.com/id/hotels/pekanbaru/read-article/liburan-asyik-di-pekanbaru-bersama-batiqa-hotel-pekanbaru` (Official BATIQA article — Indonesian)
- `https://venuemagz.com/hotel/batiqa-hotel-pekanbaru-properti-kedua-di-pulau-sumatera/` (VenueMagz 2016-08-31, quote BATIQA, 125+7)
- `https://www.berita-bisnis.com/buka-hotel-di-pekanbaru-batiqa-kini-kelola-lima-batiqa-hotel/` (Berita Bisnis 2016-09-01, quote Wadir Surya Internusa Hotels: 125+7)

**Source B (MEDIUM confidence - Official BATIQA press release awal):**
> 133 kamar — terdiri 126 Superior + 7 Suites
- `https://ekonomi.bisnis.com/read/20160406/12/535326/garap-pasar-sumatra-batiqa-hotel-hadir-di-pekanbaru` (Bisnis.com 2016-04-06, quote Michael Tjahaja, Wakil Presdir: 126+7)
- `https://www.berita-bisnis.com/awal-mei-2016-batiqa-hotel-pekanbaru-mulai-beroperasi/` (Berita Bisnis 2016-04-09, 126+7)

**Source C (LOW — EXTERNAL, Agoda):**
> 132 kamar total (tidak rincian 126/125), Suite 44 m² (bukan 45)
- `https://www.agoda.com/batiqa-hotel-pekanbaru/hotel/pekanbaru-id.html` (Agoda, Generative AI assisted)

**Recommended action:** HUMAN VERIFICATION REQUIRED — Tanya Front Office / GM Alpino Indra Putra: berapa konfigurasi aktual 2026? 125 vs 126 Superior, 44 vs 45 Suites. Untuk AI, gunakan **132 kamar** sebagai default (dominant HIGH), tandai `NEEDS HUMAN VERIFICATION` untuk detail.

### CONFLICT #2 — Luas Kamar Superior

**Source A:** 22,5 m²
- Official article + VenueMagz (22,5)

**Source B:** 22 m²
- Trip.com (22 m² | Lantai 3-10)
- Agoda (22 m²)

**Recommended action:** HUMAN VERIFICATION REQUIRED — Perbedaan 0,5 m² mungkin pembulatan. Gunakan **22,5 m²** (Official) dengan confidence MEDIUM, catat alternatif 22 m².

### CONFLICT #3 — Jarak Bandara

**Source A:** 1,7 km (Direction to Hotel — instruksi: Head north on Jalan Bandara SSK II, U-turn, 250m left)
- Official BATIQA `hotels/pekanbaru` Direction

**Source B:** 2,3 km / 5 menit mobil (Welcome text)
- Official BATIQA `hotels/pekanbaru` Welcome + VenueMagz (2,3km 5 menit)

**Recommended action:** HUMAN VERIFICATION REQUIRED — Keduanya dari Official (mungkin 1,7 km linear vs 2,3 km driving). Sajikan sebagai `1,7 km (straight) / 2,3 km (driving, ±5 menit)` dengan catatan conflict.

### CONFLICT #4 — Suite Count Detail

**Source A:** 7 Suites 45 m² (Official)
**Source B:** Agoda Suite 44 m²
**Recommended action:** Gunakan 45 m² HIGH, catat 44 m² sebagai EXTERNAL rounding.

### CONFLICT #5 — Breakfast Paid vs Include

**Source A (JTB):** Buffet 06:00-10:00 **paid** (Breakfast buffet daily 6:00-10:00 paid)
**Source B (Promo):** Breakfast untuk 2 termasuk di paket Kartini Staycation Deal (RiauAktual 2026-04-22)
**Recommended action:** Tidak conflict — Breakfast **umum berbayar** (paid), **termasuk** hanya di paket promo. AI harus jawab: `Breakfast buffet 06:00-10:00, berbayar; termasuk di paket tertentu, hubungi Front Office untuk promo.`

---

## 3. Missing Information

Kategori berikut **tidak ditemukan di publik** — **UNKNOWN** (jangan diarang):

- **SOP Housekeeping:** Prosedur permintaan handuk/bantal/pembersihan, SLA waktu respon, prioritas, stok amenities — UNKNOWN
- **SOP Engineering:** Prosedur AC/TV/WiFi/lampu/shower/plumbing, SLA HIGH (AC mati, bocor) vs MEDIUM, eskalasi, teknisi on-duty — UNKNOWN
- **Front Office SOP:** Jam operasional Front Office (implisit 24 jam tapi tidak eksplisit), nomor extension Front Office, prosedur check-in/out detail (early/late fee), deposit, identitas — UNKNOWN
- **Harga & Biaya:** Harga kamar real-time, harga breakfast, harga laundry, biaya airport transfer (gratis tapi apakah ada syarat?), biaya meeting room, biaya parkir EV — UNKNOWN (hanya signature menu 55-60k/45-50k dari Tribun, bukan tarif kamar)
- **Extension & Internal Phone:** Nomor kamar internal, extension housekeeping/engineering — UNKNOWN
- **Kebijakan Hotel (Policy):** Smoking policy detail, pet policy, child policy detail (hanya Agoda 4-10 gratis, butuh verifikasi), cancellation policy, extra bed, deposit — UNKNOWN
- **Emergency/Safety:** Jalur evakuasi, assembly point, nomor darurat internal, prosedur kebakaran, P3K, APAR — UNKNOWN
- **Accessibility:** Detail kursi roda, ramp, kamar difabel, toilet difabel, lift braille — UNKNOWN (hanya elevator terverifikasi)
- **Room Numbers:** Daftar nomor kamar per lantai (mis. 301-310 = lantai 3), nomor suite per lantai — UNKNOWN
- **Staff Structure:** Nama staff selain GM Alpino Indra Putra, struktur organisasi, shift — UNKNOWN
- **Housekeeping Inventory:** Daftar amenities lengkap (sabun, shampoo, sikat gigi, sandal, dll) — UNKNOWN (hanya general: bath amenities gratis)
- **Engineering Inventory:** Merk AC/TV, tipe Wi-Fi, kapasitas listrik — UNKNOWN
- **WiFi Detail:** SSID exact (`BATIQA HOTELS` per seed? belum verifikasi), password, kecepatan Mbps — UNKNOWN
- **Gym Detail Jam:** 24 jam (Agoda bilang gym 24 jam? Official bilang tempat gym, JTB bilang fitness, Agoda bilang spacious well-lit, tapi jam 24 jam hanya untuk gym di seed, perlu verifikasi) — UNKNOWN
- **Spa/Sauna operasional:** Apakah sauna masih ada, jam, biaya — UNKNOWN
- **Parking Detail:** Kapasitas jumlah mobil, lokasi car park multiple, EV station jumlah/tipe — UNKNOWN
- **Local Recommendations Detail:** Jarak presisi (0,5km, 0,8km) ke mall, harga rekomendasi eksternal — UNKNOWN (Agoda list Warung Pempek dll tanpa jarak)
- **Insider's Guide:** Konten `hotels/pekanbaru/insiders-guide` tidak ter-fetch (transport error) — UNKNOWN

---

## 4. Information Requiring Human Verification

Semua yang bertanda **NEEDS HUMAN VERIFICATION** di §1 + **CONFLICT** di §2, plus:

| Category | Question untuk GM / Front Office Pekanbaru |
|---|---|
| **Room Count** | 132 vs 133? 125 vs 126 Superior? |
| **Room Size** | 22 vs 22,5 m²? 44 vs 45 m² Suite? |
| **Airport Distance** | 1,7 vs 2,3 km? Mana yang publikasi resmi terbaru? |
| **Breakfast** | Apakah buffet 06:00-10:00 masih paid? Harga? Apakah termasuk paket? |
| **Signature Menu** | Apakah 7 menu 11 Aug 2026 masih tersedia harian ala carte? Harga masih 55-60k/45-50k? |
| **Bar/Sauna/Rooftop** | Apakah Bar, Sauna, Rooftop lounge masih operasional? Jam & biaya? |
| **WiFi** | SSID/password, kecepatan, apakah gratis tanpa password? |
| **Gym** | Jam operasional exact (24 jam?), perlu booking? |
| **Spa** | Apakah ada spa/sauna, atau hanya gym? |
| **Check-in/out** | 14:00/12:00 masih? Biaya early/late? 24h check-in untuk Suite promo masih? |
| **Parking** | Benar gratis? EV station gratis? Kapasitas? |
| **Housekeeping SOP** | SLA untuk towel (MEDIUM) vs pillow (LOW) vs cleaning? |
| **Engineering SOP** | SLA AC HIGH berapa menit? Teknisi standby? |
| **Emergency** | Nomor darurat, jalur evakuasi, APAR? |
| **Accessibility** | Kamar difabel, ramp, elevator akses? |
| **Child Policy** | 4-10 gratis masih berlaku? Syarat? |
| **Suite Distribution** | Benar 7 Suites, 1 per lantai? Nomor kamar suite 1 per lantai? |
| **Mushalla** | Apakah mushalla lantai 2 masih ada? |
| **Contact** | Apakah WA 0812-68933366 masih? Email reservation.pekanbaru@batiqa.com masih? |
| **Local Recommendations** | Rekomendasi resmi hotel untuk tamu (Agoda list Warung Pempek etc. apakah endorsed hotel?) |

**Action:** Jadwalkan 30 menit interview dengan **GM Alpino Indra Putra** atau **Front Office Manager** Pekanbaru, bawa daftar ini.

---

## 5. Recommended Database Entities

> **JANGAN INSERT SEKARANG.** Rekomendasi untuk `migrations\003_batiqa_real_world.sql` **setelah HUMAN REVIEW**.

### A. `hotel_information` — Perluas kategori (existing: BREAKFAST, WIFI, etc.)

| category | title (contoh) | content (verified, jangan karang) |
|---|---|---|
| `BREAKFAST` | Breakfast Schedule | Breakfast buffet 06:00-10:00 di FRESQA Bistro indoor & outdoor. Berbayar; termasuk di paket Kartini Staycation Deal. Hubungi Front Office untuk promo. (Source: JTB + RiauAktual) |
| `CHECKIN` | Check-in Time | 14:00 (02:00 PM). Early check-in tergantung ketersediaan. (Agoda) |
| `CHECKOUT` | Check-out Time | 12:00 PM. Late check-out hubungi Front Office. (Agoda) |
| `ROOM` | Superior Room 22,5 m² | 125 Superior (22,5 m²) Double & Twin, Suites 7 (45 m²). AC, LCD TV internasional, Wi-Fi high-speed, minibar, kulkas, brankas, desk, telepon, kettle, amenities, shower, slippers. Lantai 3-10. (Official 132) — flag CONFLICT |
| `WIFI` | Hotel WiFi | Free WiFi throughout hotel, di kamar & public area, high-speed, complimentary. SSID/password: NEEDS HUMAN VERIFICATION. (Official + JTB) |
| `GYM` | Fitness Center Lantai 2 | State-of-the-art, treadmill, bike, elliptical, free weights, towels & water, spacious well-lit. (Agoda) — jam NEEDS VERIFICATION |
| `RESTAURANT` | FRESQA Bistro | Indoor & outdoor, international cuisine, lunch & dinner, coffee shop, bar. (Official) |
| `RESTAURANT` | 7 Signature Menu 11 Aug 2026 | Arsik Baramundi, Dori Dabu-Dabu, Tahu Cabe Garam, Otak-Otak Kenangan; Minuman Kunyit Asam Squash, Lembayung Rasa, Kopiqawan. Rp55-60k makan, 45-50k minum. Ala carte harian. (Tribun/RiauPos, GM Alpino) |
| `PARKING` | Car Park | Multiple car parks onsite, complimentary, private exclusive, no reservation, EV station. (Trip.com) |
| `TRANSFER` | Airport Shuttle | Gratis 24 jam, 1,7 km / 2,3 km (5 menit) ke Sultan Syarif Kasim II. Hubungi Front Office. (JTB + Official — CONFLICT distance) |
| `FRONT_OFFICE` | 24h Front Desk | Reception hall, bell staff, luggage storage, wake-up service, elevator, safety deposit box. (JTB + Trip.com) |
| `HOUSEKEEPING` | Daily Housekeeping | Pembersihan setiap hari, laundry/dry cleaning (biaya tambahan). (JTB) |
| `POLICY` | Child Policy | Anak 4-10 gratis (Agoda, NEEDS VERIFICATION) |
| `LOCATION` | Simpang Tiga | Central business district, shopping malls, culinary centers, karaoke, local markets, Mall Pekanbaru, Pasar Wisata, Ar-Rahman Mosque walking distance. (Official + Agoda) |

**Rekomendasi kolom baru (jika perlu, tapi jangan ubah schema tanpa review):** `source_url`, `confidence` (HIGH/MEDIUM/LOW), `last_verified` (DATE), `needs_human_verification` (BOOLEAN) — untuk Data Governance.

### B. `recommendations` — Hanya yang terverifikasi atau EXTERNAL jelas

| name | category | description | price_min | price_max | distance_km | address | source | confidence |
|---|---|---|---|---|---|---|---|---|
| Warung Pempek Sikok Nak Duo | restaurant | Pempek tradisional, local favorite | UNKNOWN | UNKNOWN | UNKNOWN | UNKNOWN (dekat hotel) | Agoda (EXTERNAL) | LOW |
| RM Puti Buana | restaurant | Masakan Indonesia cozy | UNKNOWN | UNKNOWN | UNKNOWN | UNKNOWN | Agoda (EXTERNAL) | LOW |
| Chilli Crab ID | restaurant | Seafood crab pedas | UNKNOWN | UNKNOWN | UNKNOWN | UNKNOWN | Agoda (EXTERNAL) | LOW |
| The Baliview | restaurant | Luxury Villas & Resto | UNKNOWN | UNKNOWN | UNKNOWN | UNKNOWN | Agoda (EXTERNAL) | LOW |
| Pasar Wisata | tourism | Sate Pekanbaru, gulai kambing, oleh-oleh | UNKNOWN | UNKNOWN | UNKNOWN | UNKNOWN | Agoda (EXTERNAL) | LOW |
| Mall Pekanbaru | shopping | Pusat perbelanjaan | UNKNOWN | UNKNOWN | UNKNOWN | UNKNOWN | Official (mall) | MEDIUM |

**Catatan:** Untuk recommendations, **jangan isi price/distance jika tidak ada sumber** — biarkan `NULL` (sesuai `recommendations.price_min` nullable). Lebih baik `UNKNOWN` daripada karang.

### C. `staff` — Tidak perlu ubah, tapi pastikan `department` sesuai operasional Pekanbaru (HOUSEKEEPING, ENGINEERING, FRONT_OFFICE, ADMIN)

### D. `tickets.category` — Tambah kategori baru jika diperlukan untuk Pekanbaru real (mis. `MUSHALLA_REQUEST`, `PARKING_REQUEST`) — **NEEDS HUMAN VERIFICATION** apakah housekeeping/engineering/front_office cukup.

---

## 6. Recommended AI Knowledge Categories

> Untuk `ai/intent.go` dan `hotel_information.category` — sinkron dengan DB.

**Existing (tetap):** `BREAKFAST_INFORMATION`, `WIFI_INFORMATION`, `CHECKIN_INFORMATION`, `CHECKOUT_INFORMATION`, `FACILITY_INFORMATION`, `ROOM_INFORMATION`, `HOTEL_POLICY`, `RESTAURANT_INFORMATION` (baru, untuk FRESQA), `PARKING_INFORMATION` (baru), `TRANSFER_INFORMATION` (baru), `MEETING_INFORMATION` (baru)

**Baru untuk Pekanbaru real:**
- `FRESQA_MENU` — 7 signature menu (makanan/minuman, harga, jam)
- `SIGNATURE_MENU` — Arsik Baramundi etc.
- `MUSHALLA_INFORMATION` — lantai 2, gym & meeting
- `CAR_PARK_INFORMATION` — multiple, gratis, EV
- `AIRPORT_SHUTTLE` — gratis 24 jam, 1,7/2,3 km

**Intent mapping:**
- `RESTAURANT_RECOMMENDATION` → `FRESQA Bistro` + 7 signature
- `TRANSPORTATION_REQUEST` → `AIRPORT_SHUTTLE` + `CAR_PARK`
- `PARKING_REQUEST` → `CAR_PARK`

**Bahasa:** Tetap `LANGUAGE BEHAVIOR.md` id/en, tambah Melayu local? Tidak perlu, cukup id/en.

---

## 7. Recommended Operational Workflows

> Untuk `ticket routing` dan `priority` — sinkron dengan `TICKET ROUTING.md` & `PRIORITY CLASSIFICATION.md`.

**Housekeeping (HOUSEKEEPING):**
- `TOWEL_REQUEST` (MEDIUM, SLA NEEDS VERIFICATION) — 2 handuk kamar 305 → Housekeeping → `OPEN→IN_PROGRESS (staff ambil handuk) → RESOLVED (antar)`
- `AMENITY_REQUEST` (LOW) — bantal/selimut/amenities → HOUSEKEEPING
- `ROOM_CLEANING_REQUEST` (MEDIUM) — pembersihan harian → HOUSEKEEPING
- `LAUNDRY_REQUEST` (MEDIUM, biaya tambahan) — **Baru, perlu HUMAN VERIFICATION** apakah housekeeping atau front office
- `MUSHALLA_REQUEST` — **UNKNOWN** departemen

**Engineering (ENGINEERING):**
- `AC_PROBLEM` (HIGH) — AC tidak dingin kamar 305 → ENGINEERING → `HIGH` karena `AC completely unavailable` (PRIORITY HIGH) — SLA NEEDS VERIFICATION (mis. 15 menit)
- `TV_PROBLEM` (MEDIUM) — TV mati
- `WIFI_PROBLEM` (MEDIUM-HIGH?) — WiFi tidak bekerja → ENGINEERING (tapi bisa FRONT_OFFICE untuk password) — **NEEDS VERIFICATION** routing
- `LIGHT_PROBLEM` (MEDIUM)
- `SHOWER_PROBLEM` / `PLUMBING_PROBLEM` (HIGH jika bocor)
- `ROOM_EQUIPMENT_PROBLEM` (MEDIUM) — brankas, kettle, telepon

**Front Office (FRONT_OFFICE):**
- `CHECKIN_INFORMATION` / `CHECKOUT_INFORMATION` → FRONT_OFFICE (tidak buat ticket, hanya jawab)
- `BREAKFAST_INFORMATION` → FRONT_OFFICE (tidak ticket, tapi jika request antar ke kamar? → ROOM_SERVICE)
- `PARKING_REQUEST` / `AIRPORT_SHUTTLE` → FRONT_OFFICE (gratis, tapi perlu booking) — **Baru**
- `MEETING_REQUEST` → FRONT_OFFICE (4 modular rooms) — **Baru**

**Priority (jangan exaggerate):**
- `HIGH`: AC mati total, bocor, safety — perlu HUMAN VERIFICATION SLA
- `MEDIUM`: towel, cleaning, TV minor
- `LOW`: extra pillow, amenities

**Confirmation Policy:** Untuk Pekanbaru real, `AC rusak` → konfirmasi `Baik, saya laporkan AC kamar 305 tidak dingin ke Engineering. Benar?` (sudah ada), untuk `handuk 2` → langsung `Baik, saya kirim ke Housekeeping`.

---

## 8. Recommended Guest Scenarios

> Untuk `USER_FLOW.md` dan testing `AI CHAT FLOW.md` — semua harus tanpa training, natural language.

**S1 — Hotel Info (No Ticket):**
- Guest: `Breakfast sampai jam berapa?` → AI `BREAKFAST_INFORMATION` → DB `06:00-10:00 FRESQA` → `No ticket` — **TESTED Phase 7 PASS**

**S2 — Housekeeping (Towel):**
- Guest: `Tolong antar 2 handuk ke kamar 305.` → `TOWEL_REQUEST` `HOUSEKEEPING` `MEDIUM` → `TKT-...` → Staff HK dashboard → Guest `My Requests` `OPEN` — **TESTED PASS**

**S3 — Engineering Missing Room:**
- Guest: `AC kamar saya tidak dingin.` (tanpa room) → AI `AC_PROBLEM` → `Boleh saya tahu nomor kamar Anda?` → Guest `305` → `TKT-...` `ENGINEERING HIGH` — **TESTED PASS**

**S4 — Engineering EN:**
- Guest: `My AC is not working in room 212.` → `AC_PROBLEM` `ENGINEERING HIGH` en → **TESTED PASS**

**S5 — Unknown:**
- Guest: `asdasd qwerty` → `UNKNOWN` `Clarification` no ticket — **TESTED PASS**

**S6 — AI Failure:**
- Guest: `Hello` saat `GEMINI_API_KEY` invalid → fallback `Maaf, layanan AI gangguan` no ticket, server 200 — **TESTED PASS**

**S7 — Pekanbaru Real — FRESQA Menu (Baru, NEEDS VERIFICATION):**
- Guest: `Apa menu signature BATIQA?` → AI `RESTAURANT_INFORMATION` → DB `7 signature` (Arsik Baramundi etc., 55-60k) → no ticket
- Expected: `FRESQA Bistro 7 signature: Arsik Baramundi... Rp55-60k, tersedia harian ala carte.`

**S8 — Pekanbaru Real — Airport Shuttle (Baru):**
- Guest: `Ada antar jemput bandara gratis?` → AI `TRANSFER_INFORMATION` → DB `Gratis 24 jam, 1,7/2,3 km` → no ticket

**S9 — Pekanbaru Real — Parking (Baru):**
- Guest: `Parkir mobil gratis? Ada EV charger?` → AI `PARKING_INFORMATION` → DB `Multiple car parks gratis, private, EV station` → no ticket

**S10 — Pekanbaru Real — Gym Lantai 2:**
- Guest: `Gym buka jam berapa?` → AI `FACILITY_INFORMATION` → DB `Gym lantai 2, peralatan lengkap, towels & water, well-lit` — **Perlu HUMAN VERIFICATION jam 24 jam atau tidak**

**S11 — Pekanbaru Real — Mushalla:**
- Guest: `Ada mushalla?` → AI `MUSHALLA_INFORMATION` → DB `Lantai 2 bersama gym & meeting` (infoPKU) — **LOW confidence, NEEDS VERIFICATION**

**S12 — Pekanbaru Real — Room Service 24h:**
- Guest: `Bisa pesan room service jam 2 pagi?` → AI `ROOM_SERVICE` → DB `24 jam` → no ticket, tapi jika pesan makanan → `ROOM_SERVICE_REQUEST` → `HOUSEKEEPING` atau `FRESQA`? **NEEDS VERIFICATION** routing

---

## 9. Data Governance Rules

> **PENTING — JANGAN MENGARANG** (sesuai instruksi):

1.  **Sumber:** Hanya `Official BATIQA` (batiqa.com) sebagai **HIGH** confidence. `Agoda/JTB/Trip.com/Bisnis.com` sebagai **EXTERNAL SOURCE** **MEDIUM/LOW**, harus tandai.
2.  **Harga/Jam/SOP:** Jika tidak di Official, tulis `UNKNOWN` atau `NEEDS HUMAN VERIFICATION`, jangan karang `Rp` atau `jam`.
3.  **Konflik:** Jika `Source A != Source B` (mis. 132 vs 133 kamar), tulis `CONFLICT DETECTED` dengan kedua sumber + `HUMAN VERIFICATION REQUIRED`, jangan pilih otomatis.
4.  **Confidence:** `HIGH` = Official + konsisten lintas sumber, `MEDIUM` = Official tapi konflik minor atau EXTERNAL konsisten, `LOW` = EXTERNAL saja atau single source.
5.  **Insert:** **JANGAN** `INSERT` ke `hotel_information`/`recommendations` sebelum `HUMAN REVIEW` audit ini disetujui. Buat `migrations\003_batiqa_real_world.sql` **setelah** approval, dengan kolom `source_url`, `confidence`, `last_verified`.
6.  **Update:** Jika GM verifikasi 125 vs 126, update `migrations` + `docs` + `AI knowledge`, bukan hardcode di `mock_provider.go`.
7.  **Privacy:** Jangan publish `staff` internal, `extension`, `SOP` internal ke publik.
8.  **Versioning:** Setiap perubahan DB real-world harus `migrations` baru, bukan edit `001/002`.

---

## 10. Source References

### Official BATIQA (HIGH)

1.  **BATIQA Hotel Pekanbaru Main** — `https://www.batiqa.com/hotels/pekanbaru` — Alamat, Direction 1,7km & 2,3km, FRESQA, Gym, 132 kamar highlight — *websearch highlight 2026-08-21* (transport error parsial, fallback highlight)
2.  **BATIQA Article Liburan Asyik** — `https://www.batiqa.com/id/hotels/pekanbaru/read-article/liburan-asyik-di-pekanbaru-bersama-batiqa-hotel-pekanbaru` — 132 kamar (125+7, 22,5m² & 45m²), 4 meeting rooms, FRESQA Bistro, gym, wifi high-speed, TV Kabel, room service 24h — *websearch highlight*
3.  **BATIQA Subpages** — `/rooms`, `/facilities`, `/dining`, `/meetings` — disebut di highlight tapi tidak ter-fetch (transport error) — **NEEDS HUMAN VERIFICATION** untuk detail subpage

### Berita & Press (MEDIUM — Quote GM, tapi EXTERNAL)

4.  **TribunPekanbaru 2026-08-11** — `https://pekanbaru.tribunnews.com/adv/1111778/...-tujuh-menu-signature` — 7 signature menu, GM Alpino Indra Putra, 4 makan 3 minum, 55-60k/45-50k
5.  **RiauPos 2026-08-12** — `https://riaupos.jawapos.com/ekonomi/2608120019/...-7-menu-baru` — 7 menu sama, GM Alpino
6.  **Bisnis.com 2016-04-06** — `https://ekonomi.bisnis.com/read/20160406/12/535326/...` — 133 kamar (126+7), Michael Tjahaja Wakil Presdir, buka 4 Mei 2016
7.  **VenueMagz 2016-08-31** — `https://venuemagz.com/hotel/batiqa-hotel-pekanbaru-properti-kedua-di-pulau-sumatera/` — 132 kamar (125+7, 22,5 & 45), lounge, 4 meeting, FRESQA, gym, 24h room service, 2,3km
8.  **Berita Bisnis 2016-09-01** — `https://www.berita-bisnis.com/buka-hotel-di-pekanbaru-batiqa-kini-kelola-lima-batiqa-hotel/` — 132 kamar (125+7, 45m²), lounge, 4 meeting, FRESQA, gym
9.  **RiauAktual 2026-04-22** — Kartini Staycation Deal (Superior + smart TV + breakfast 2) — promo
10. **BertuahPos 2019-08-26** — `https://bertuahpos.com/...hut-ke-3-tahun.html` — WA 0812-68933366, reservation.pekanbaru@batiqa.com

### EXTERNAL SOURCE (LOW — Agoda/JTB/Trip.com, Generative AI assisted, perlu verifikasi)

11. **Agoda BATIQA Hotel Pekanbaru** — `https://www.agoda.com/batiqa-hotel-pekanbaru/hotel/pekanbaru-id.html` — 132 rooms, check-in 14:00, check-out 12:00, child 4-10 free, bar, sauna, parking, AC, TV, WiFi, 22m², 44m² Suite, breakfast buffet, business center — *Generative AI assisted* — **LOW**
12. **JTB BATIQA Hotel Pekanbaru** — `https://www.jtb.co.jp/ovs_htl/detail/search_detail/175089/` — 132 rooms, breakfast 6:00-10:00 paid 110k, child 75k, gym, WiFi free, parking free, smoking area, elevator, laundry, 24h front desk, safety box, etc. — **LOW**
13. **Trip.com BATIQA Hotel Pekanbaru** — `https://id.trip.com/hotels/bukit-raya-hotel-detail-5707095/batiqa-hotel-pekanbaru/` — Superior Double Twin 22m² lantai 3-10, Suite, 6, 8, parking private, EV station, 24h reception, wake-up, meeting (biaya tambahan) — **LOW**
14. **infoPKU 2019-07-16** — `https://infopku.com/menginap-di-batiqa-pekanbaru/30107/` — 9 Suites (1 per lantai, nomor 1), meeting lantai 2, gym, mushalla, cafe resto, airport shuttle — **LOW** (9 Suites konflik dengan 7)
15. **Trip.com fasilitas list** — `https://id.trip.com/...` — EV station, wake-up, etc. — **LOW**

### Info Pekanbaru (Riau) — Umum

16. **Batiqa.com About BHM** — PT Batiqa Hotel Manajemen, PT Surya Semesta Internusa Tbk, 5 hotels, 664 rooms total — *Berita Bisnis*
17. **Lokasi Kota** — Pekanbaru *major economic center on eastern Sumatra* — Official highlight

---

> **Next Step:** Tunggu **HUMAN REVIEW** dari user (GM/staff Pekanbaru) untuk `CONFLICT` & `NEEDS VERIFICATION` di §2 & §4. Setelah approved, buat `migrations\003_batiqa_real_world.sql` dengan `INSERT` verified + `source_url` + `confidence`, update `ai/intent.go` & `hotel_information` kategori, dan `docs` — **jangan langsung INSERT sekarang**.



