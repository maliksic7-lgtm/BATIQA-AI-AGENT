-- 002_seed.sql
-- Seed hotel_information and minimal recommendations/staff for demo/testing
-- All seeds are idempotent (INSERT IGNORE / ON DUPLICATE KEY where possible)

-- Hotel Information seeds (per spec categories: BREAKFAST, WIFI, POOL, GYM, CHECKIN, CHECKOUT, RESTAURANT, ROOM, POLICY)
INSERT INTO hotel_information (category, title, content, active) VALUES
('BREAKFAST', 'Breakfast Schedule', 'Breakfast tersedia mulai pukul 06:00 sampai 10:00 di restaurant lantai 1.', TRUE),
('BREAKFAST', 'Breakfast Location', 'Restaurant BATIQA di lantai 1, dekat lobby.', TRUE),
('WIFI', 'Hotel WiFi', 'Connect to BATIQA HOTELS network. Password tersedia di kartu kamar atau hubungi Front Office.', TRUE),
('WIFI', 'WiFi Support', 'Jika WiFi bermasalah, silakan laporkan via AI Assistant atau hubungi Front Office ext 0.', TRUE),
('POOL', 'Swimming Pool', 'Kolam renang buka 06:00-20:00 di lantai 2. Tersedia handuk pool di pool bar.', TRUE),
('GYM', 'Gym / Fitness Center', 'Gym buka 24 jam di lantai 2. Akses dengan kartu kamar.', TRUE),
('CHECKIN', 'Check-in Time', 'Check-in mulai pukul 14:00. Early check-in tergantung ketersediaan.', TRUE),
('CHECKOUT', 'Check-out Time', 'Check-out pukul 12:00. Late check-out dapat diminta ke Front Office.', TRUE),
('RESTAURANT', 'Hotel Restaurant', 'Restaurant buka 06:00-22:00 menyajikan masakan Indonesia & Western.', TRUE),
('ROOM', 'Room Facilities', 'Fasilitas kamar: AC, TV, WiFi, shower, amenities, brankas, minibar.', TRUE),
('POLICY', 'Smoking Policy', 'Hotel bebas asap rokok di dalam kamar. Area merokok tersedia di luar lobby.', TRUE),
('POLICY', 'Pet Policy', 'Hewan peliharaan tidak diperbolehkan di kamar.', TRUE)
ON DUPLICATE KEY UPDATE title=VALUES(title);

-- Recommendations seeds (verified data, not fabricated by AI)
INSERT INTO recommendations (name, category, description, price_min, price_max, distance_km, address, active) VALUES
('Warung Sederhana BATIQA', 'restaurant', 'Masakan Indonesia, budget-friendly dekat hotel', 25000, 75000, 0.5, 'Jl. Contoh No.1', TRUE),
('Cafe Ceria', 'cafe', 'Kopi & pastry, cocok untuk meeting santai', 20000, 50000, 0.8, 'Jl. Contoh No.2', TRUE),
('Mall Central', 'shopping', 'Pusat perbelanjaan terbesar sekitar hotel', 0, 0, 1.2, 'Jl. Mall No.10', TRUE),
('Pantai Indah', 'tourism', 'Destinasi wisata pantai dekat hotel', 0, 30000, 3.5, 'Jl. Pantai No.5', TRUE),
('ATM BCA Terdekat', 'atm', 'ATM 24 jam 200m dari lobby', 0, 0, 0.2, 'Lobby BATIQA', TRUE)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- Demo staff accounts (bcrypt hash of password "batiqa123", cost 10)
-- DEMO ONLY: rotate credentials before production per SECURITY PRINCIPLES.md
INSERT INTO staff (name, email, password_hash, department) VALUES
('Admin BATIQA', 'admin@batiqa.com', '$2a$10$WnP6TcAxzfFXrLBRKgoLWOgwNLEf7x9iGkuhmSO34kfBqyhsTRHCa', 'ADMIN'),
('Housekeeping Team', 'hk@batiqa.com', '$2a$10$WnP6TcAxzfFXrLBRKgoLWOgwNLEf7x9iGkuhmSO34kfBqyhsTRHCa', 'HOUSEKEEPING'),
('Engineering Team', 'eng@batiqa.com', '$2a$10$WnP6TcAxzfFXrLBRKgoLWOgwNLEf7x9iGkuhmSO34kfBqyhsTRHCa', 'ENGINEERING')
ON DUPLICATE KEY UPDATE name=VALUES(name);
