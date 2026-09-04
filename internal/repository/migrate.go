package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"batiqa-ai/internal/model"
)

// Migrate creates collections/indexes and seeds initial data.
// Idempotent: safe to run multiple times. Replaces the former SQL migrations.
func Migrate(ctx context.Context, db *mongo.Database) error {
	if err := ensureIndexes(ctx, db); err != nil {
		return fmt.Errorf("ensure indexes: %w", err)
	}
	if err := seedHotelInformation(ctx, db); err != nil {
		return fmt.Errorf("seed hotel_information: %w", err)
	}
	if err := seedRecommendations(ctx, db); err != nil {
		return fmt.Errorf("seed recommendations: %w", err)
	}
	if err := seedDynamicInfo(ctx, db); err != nil {
		return fmt.Errorf("seed dynamic info: %w", err)
	}
	if err := seedStaff(ctx, db); err != nil {
		return fmt.Errorf("seed staff: %w", err)
	}
	return nil
}

// seedDynamicInfo inserts events, daily menu, and weather-context entries into
// hotel_information with categories EVENT, DAILY_MENU, and WEATHER. Idempotent
// per unique Title so re-running migrate only adds missing rows (unlike the
// count-based seedHotelInformation which skips once any rows exist).
func seedDynamicInfo(ctx context.Context, db *mongo.Database) error {
	col := db.Collection(hotelInfoCol)
	now := time.Now()
	type item struct{ Category, Title, Content string }
	// Data event & menu bersifat contoh nyata Pekanbaru/BATIQA; perbarui berkala.
	items := []item{
		{"EVENT", "Upcoming Event: Kuliner Kulineran Riau", "Pekanbaru Culinary Festival pekan ini di SKA Mall (Jl. Tuanku Tambusai): aneka kuliner Melayu, live music, dan tenant UMKM. Jumat-Minggu, sore sampai malam."},
		{"EVENT", "Upcoming Event: Laksamana Parade", "Event budaya di halaman Museum Sang Nila Utama Pekanbaru akhir pekan ini: pameran budaya Melayu, bazar, dan pertunjukan zapin."},
		{"EVENT", "Upcoming Event: Hotel Weekend Live Jazz", "FRESQA Bistro menghadirkan live acoustic & jazz tiap Jumat & Sabtu malam mulai 19.00-22.00. Reservasi meja via Front Office."},
		{"EVENT", "Upcoming Event: UMKM Pasar Rakyat", "Pasar rakyat tiap Minggu pagi di Anjungan Sungai Siak: produk lokal, kuliner khas Riau, dan area bermain anak."},
		{"DAILY_MENU", "Today's Breakfast Menu", "Prasmanan sarapan hari ini: Nasi Uduk, Bubur Ayam, Omelet, Roti, Sereal, buah segar, kopi/teh. Buka 06.00-10.00 di FRESQA Bistro."},
		{"DAILY_MENU", "Today's Lunch/Dinner Special", "Special hari ini di FRESQA Bistro: Gulai Ikan Patin, Ayam Panggang Madu, dan Sate Padang. A la carte mulai Rp45.000. Tersedia room service 24 jam."},
		{"WEATHER", "Pekanbaru Climate Note", "Pekanbaru beriklim tropis panas-lembap (30-33°C siang). Musim hujan umumnya Nov-Mar; hujan sering datang sore-malam singkat. Selalu sedia payung di kamar."},
	}
	for _, it := range items {
		var existing model.HotelInformation
		err := col.FindOne(ctx, bson.M{"category": it.Category, "title": it.Title}).Decode(&existing)
		switch {
		case err == mongo.ErrNoDocuments:
			id, err := nextID(ctx, db, hotelInfoCol)
			if err != nil {
				return err
			}
			if _, err := col.InsertOne(ctx, model.HotelInformation{
				ID: id, Category: it.Category, Title: it.Title, Content: it.Content,
				Active: true, CreatedAt: now, UpdatedAt: now,
			}); err != nil {
				return err
			}
			log.Printf("  seeded dynamic info [%s] %s", it.Category, it.Title)
		case err != nil:
			return err
		}
	}
	return nil
}

// ensureIndexes creates all required indexes (no-op if they already exist).
func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := map[string][]mongo.IndexModel{
		guestsCol: {
			mongo.IndexModel{Keys: bson.M{"session_id": 1}, Options: options.Index().SetUnique(true)},
		},
		conversationsCol: {
			mongo.IndexModel{Keys: bson.D{{Key: "session_id", Value: 1}, {Key: "created_at", Value: -1}}},
		},
		ticketsCol: {
			mongo.IndexModel{Keys: bson.M{"ticket_number": 1}, Options: options.Index().SetUnique(true)},
			mongo.IndexModel{Keys: bson.M{"department": 1}},
			mongo.IndexModel{Keys: bson.M{"status": 1}},
			mongo.IndexModel{Keys: bson.M{"priority": 1}},
			mongo.IndexModel{Keys: bson.M{"room_number": 1}},
			mongo.IndexModel{Keys: bson.M{"created_at": -1}},
		},
		hotelInfoCol: {
			mongo.IndexModel{Keys: bson.M{"category": 1}},
		},
		recommendationsCol: {
			mongo.IndexModel{Keys: bson.M{"category": 1}},
			mongo.IndexModel{Keys: bson.M{"distance_km": 1}},
		},
		staffCol: {
			mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)},
		},
		assignmentsCol: {
			mongo.IndexModel{Keys: bson.D{{Key: "ticket_id", Value: 1}, {Key: "staff_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		staffSessionsCol: {
			// TTL: Mongo auto-deletes expired sessions
			mongo.IndexModel{Keys: bson.M{"expires_at": 1}, Options: options.Index().SetExpireAfterSeconds(0)},
		},
	}
	for col, models := range indexes {
		if _, err := db.Collection(col).Indexes().CreateMany(ctx, models); err != nil {
			return fmt.Errorf("index %s: %w", col, err)
		}
		log.Printf("  indexes ensured on %s", col)
	}
	return nil
}

func seedHotelInformation(ctx context.Context, db *mongo.Database) error {
	col := db.Collection(hotelInfoCol)
	n, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if n > 0 {
		log.Printf("  hotel_information already seeded (%d docs)", n)
		return nil
	}
	now := time.Now()
	type item struct {
		Category string
		Title    string
		Content  string
	}
	// Data riil Hotel BATIQA Pekanbaru dari docs/BATIQA_REAL_WORLD_AUDIT.md
	// (sumber publik resmi; verifikasi berkala sebelum dipakai produksi).
	items := []item{
		{"BREAKFAST", "Breakfast Schedule", "Sarapan prasmanan tersedia setiap hari pukul 06.00-10.00 di FRESQA Bistro lantai 1 (sesuai paket menginap)."},
		{"BREAKFAST", "Breakfast Location", "FRESQA Bistro berada di lantai 1, beberapa langkah dari lobby."},
		{"WIFI", "Hotel WiFi", "Terhubung ke jaringan WiFi BATIQA. Password tertera pada kartu kamar Anda, atau tanyakan ke Front Office."},
		{"WIFI", "WiFi Support", "Jika WiFi bermasalah, laporkan lewat asisten AI ini atau hubungi Front Office ext 0 - tim Engineering akan menindaklanjuti."},
		{"POOL", "Swimming Pool", "Kolam renang buka pukul 07.00-19.00. Handuk kolam tersedia di pool bar."},
		{"GYM", "Gym / Fitness Center", "Fitness center berada di lantai 2, buka setiap hari pukul 06.00-22.00. Akses dengan kartu kamar."},
		{"CHECKIN", "Check-in Time", "Check-in dimulai pukul 14.00. Early check-in bergantung ketersediaan kamar - silakan minta bantuan Front Office."},
		{"CHECKOUT", "Check-out Time", "Check-out paling lambat pukul 12.00. Late check-out dapat diminta melalui Front Office (biaya dapat berlaku)."},
		{"RESTAURANT", "Hotel Restaurant", "FRESQA Bistro menyajikan menu Indonesia & Western setiap hari pukul 06.00-23.00. Room service tersedia 24 jam."},
		{"ROOM", "Room Facilities", "Setiap kamar dilengkapi AC, TV, WiFi, shower air hangat, amenities, brankas, dan meja kerja."},
		{"ROOM", "Mushalla", "Mushalla hotel berada di lantai 2, muat sekitar 20 jamaah. Mukena & sajadah tersedia di mushalla."},
		{"ROOM", "Parking & EV", "Parkir mobil gratis untuk tamu, termasuk stasiun pengisian kendaraan listrik (EV) di area parkir depan."},
		{"POLICY", "Smoking Policy", "Seluruh ruangan dalam kamar bebas asap rokok. Area merokok tersedia di luar lobby."},
		{"POLICY", "Airport Shuttle", "Shuttle bandara gratis tersedia 24 jam menuju Bandara Sultan Syarif Kasim II (sekitar 15-20 menit perjalanan). Reservasi melalui Front Office."},
		{"POLICY", "Pet Policy", "Mohon maaf, hewan peliharaan tidak diperbolehkan menginap di dalam kamar."},
	}
	docs := make([]interface{}, 0, len(items))
	for _, it := range items {
		id, err := nextID(ctx, db, hotelInfoCol)
		if err != nil {
			return err
		}
		docs = append(docs, model.HotelInformation{
			ID: id, Category: it.Category, Title: it.Title, Content: it.Content,
			Active: true, CreatedAt: now, UpdatedAt: now,
		})
	}
	if _, err := col.InsertMany(ctx, docs); err != nil {
		return err
	}
	log.Printf("  seeded %d hotel_information docs", len(docs))
	return nil
}

func seedRecommendations(ctx context.Context, db *mongo.Database) error {
	col := db.Collection(recommendationsCol)
	n, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if n > 0 {
		log.Printf("  recommendations already seeded (%d docs)", n)
		return nil
	}
	now := time.Now()
	strPtr := func(s string) *string { return &s }
	intPtr := func(i int) *int { return &i }
	f64Ptr := func(f float64) *float64 { return &f }
	type rec struct {
		name, category, desc string
		pMin, pMax           *int
		dist                 *float64
		addr                 string
		maps                 string
	}
	// Tempat nyata sekitar Jl. Jend. Sudirman / Simpang Tiga, Pekanbaru.
	// Harga bersifat estimasi publik - tandai untuk verifikasi berkala.
	// maps_link adalah URL pra-isi Google Maps (query) sehingga tamu bisa
	// langsung membuka navigasi dari HP.
	items := []rec{
		{"FRESQA Bistro - BATIQA", "restaurant", "Bistro hotel sendiri: Indonesia & Western, buka 06.00-23.00", intPtr(35000), intPtr(150000), f64Ptr(0.0), "Lantai 1, Hotel BATIQA Pekanbaru", "https://www.google.com/maps/search/?api=1&query=Hotel+BATIQA+Pekanbaru"},
		{"RM Selera Kampung", "restaurant", "Masakan Melayu-Riau legendaris, asam pedas & gulai", intPtr(25000), intPtr(80000), f64Ptr(1.8), "Jl. Diponegoro, Pekanbaru", "https://www.google.com/maps/search/?api=1&query=RM+Selera+Kampung+Pekanbaru"},
		{"Pondok Patin H. Guntur", "restaurant", "Ikan patin bakar & gulai patin khas Riau, porsi besar", intPtr(50000), intPtr(200000), f64Ptr(4.0), "Jl. Soekarno-Hatta, Pekanbaru", "https://www.google.com/maps/search/?api=1&query=Pondok+Patin+H.+Guntur+Pekanbaru"},
		{"Sapadia Coffee & Eatery", "cafe", "Kafe nyaman dengan kopi specialty & western food", intPtr(25000), intPtr(90000), f64Ptr(2.0), "Jl. Jend. Sudirman, Pekanbaru", "https://www.google.com/maps/search/?api=1&query=Sapadia+Coffee+Pekanbaru"},
		{"Kopi Koe Pekanbaru", "cafe", "Kopi lokal & suasana santai, cocok untuk bekerja", intPtr(15000), intPtr(50000), f64Ptr(3.1), "Jl. Ahmad Yani, Pekanbaru", "https://www.google.com/maps/search/?api=1&query=Kopi+Koe+Pekanbaru"},
		{"SKA Mall", "shopping", "Mal terbesar di Pekanbaru: bioskop, tenant fashion & kuliner", intPtr(0), intPtr(0), f64Ptr(2.5), "Jl. Tuanku Tambusai, Pekanbaru", "https://www.google.com/maps/search/?api=1&query=SKA+Mall+Pekanbaru"},
		{"Ciputra Seraya Mall", "shopping", "Mal modern dengan supermarket, food court & entertainment", intPtr(0), intPtr(0), f64Ptr(3.8), "Jl. Laksamana, Pekanbaru", "https://www.google.com/maps/search/?api=1&query=Ciputra+Seraya+Mall+Pekanbaru"},
		{"Anjungan Sungai Siak", "tourism", "Waterfront ikonik untuk jogging & sunset di tepi Siak", intPtr(0), intPtr(0), f64Ptr(2.2), "Jl. Raja Ali Haji, Pekanbaru", "https://www.google.com/maps/search/?api=1&query=Anjungan+Sungai+Siak+Pekanbaru"},
		{"Masjid Raya An-Nur", "tourism", "Masjid megah arsitektur Melayu, destinasi wisata religi", intPtr(0), intPtr(0), f64Ptr(1.6), "Jl. Ahmad Yani, Pekanbaru", "https://www.google.com/maps/search/?api=1&query=Masjid+Raya+An-Nur+Pekanbaru"},
		{"Taman Wisata Alam Mayang", "tourism", "Rekreasi keluarga: rekreasi air, kebun binatang mini & taman", intPtr(25000), intPtr(50000), f64Ptr(5.5), "Jl. Kaharuddin Nasution, Pekanbaru", "https://www.google.com/maps/search/?api=1&query=Taman+Wisata+Alam+Mayang+Pekanbaru"},
		{"ATM BCA Terdekat", "atm", "ATM 24 jam di lobby hotel", intPtr(0), intPtr(0), f64Ptr(0.1), "Lobby BATIQA Pekanbaru", "https://www.google.com/maps/search/?api=1&query=ATM+BCA+Hotel+BATIQA+Pekanbaru"},
		{"Airport Shuttle - SSK II", "transportation", "Shuttle gratis 24 jam ke Bandara Sultan Syarif Kasim II (~15-20 menit)", intPtr(0), intPtr(0), f64Ptr(1.9), "Reservasi via Front Office", "https://www.google.com/maps/search/?api=1&query=Sultan+Syarif+Kasim+II+Airport"},
	}
	docs := make([]interface{}, 0, len(items))
	for _, it := range items {
		id, err := nextID(ctx, db, recommendationsCol)
		if err != nil {
			return err
		}
		var pMax *int
		if it.pMax != nil && *it.pMax > 0 {
			pMax = it.pMax
		}
		docs = append(docs, model.Recommendation{
			ID: id, Name: it.name, Category: it.category, Description: strPtr(it.desc),
			PriceMin: it.pMin, PriceMax: pMax, DistanceKm: it.dist, Address: strPtr(it.addr),
			MapsLink: strPtr(it.maps), Active: true, CreatedAt: now,
		})
	}
	if _, err := col.InsertMany(ctx, docs); err != nil {
		return err
	}
	log.Printf("  seeded %d recommendation docs", len(docs))
	return nil
}

// demoStaffPassword is the bcrypt hash of "batiqa123" (cost 10).
// DEMO ONLY - rotate credentials before production per SECURITY PRINCIPLES.md
const demoStaffPassword = "$2a$10$WnP6TcAxzfFXrLBRKgoLWOgwNLEf7x9iGkuhmSO34kfBqyhsTRHCa"

func seedStaff(ctx context.Context, db *mongo.Database) error {
	col := db.Collection(staffCol)
	now := time.Now()
	staff := []model.Staff{
		{Name: "Admin BATIQA", Email: "admin@batiqa.com", PasswordHash: demoStaffPassword, Department: "ADMIN"},
		{Name: "Housekeeping Team", Email: "hk@batiqa.com", PasswordHash: demoStaffPassword, Department: model.DeptHousekeeping},
		{Name: "Engineering Team", Email: "eng@batiqa.com", PasswordHash: demoStaffPassword, Department: model.DeptEngineering},
	}
	for _, s := range staff {
		var existing model.Staff
		err := col.FindOne(ctx, bson.M{"email": s.Email}).Decode(&existing)
		switch {
		case err == mongo.ErrNoDocuments:
			// Insert with explicit integer _id (auto ObjectID would not decode into int64)
			id, err := nextID(ctx, db, staffCol)
			if err != nil {
				return err
			}
			s.ID = id
			s.CreatedAt = now
			if _, err := col.InsertOne(ctx, s); err != nil {
				return err
			}
			log.Printf("  seeded staff %s", s.Email)
		case err != nil:
			return err
		default:
			// Refresh password hash so old invalid placeholder seeds are fixed.
			if _, err := col.UpdateOne(
				ctx,
				bson.M{"_id": existing.ID},
				bson.M{"$set": bson.M{"name": s.Name, "password_hash": s.PasswordHash, "department": s.Department}},
			); err != nil {
				return err
			}
			log.Printf("  staff %s refreshed", s.Email)
		}
	}
	return nil
}
