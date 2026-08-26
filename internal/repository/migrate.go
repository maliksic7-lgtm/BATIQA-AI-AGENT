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
	if err := seedStaff(ctx, db); err != nil {
		return fmt.Errorf("seed staff: %w", err)
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
	items := []item{
		{"BREAKFAST", "Breakfast Schedule", "Breakfast tersedia mulai pukul 06:00 sampai 10:00 di restaurant lantai 1."},
		{"BREAKFAST", "Breakfast Location", "Restaurant BATIQA di lantai 1, dekat lobby."},
		{"WIFI", "Hotel WiFi", "Connect to BATIQA HOTELS network. Password tersedia di kartu kamar atau hubungi Front Office."},
		{"WIFI", "WiFi Support", "Jika WiFi bermasalah, silakan laporkan via AI Assistant atau hubungi Front Office ext 0."},
		{"POOL", "Swimming Pool", "Kolam renang buka 06:00-20:00 di lantai 2. Tersedia handuk pool di pool bar."},
		{"GYM", "Gym / Fitness Center", "Gym buka 24 jam di lantai 2. Akses dengan kartu kamar."},
		{"CHECKIN", "Check-in Time", "Check-in mulai pukul 14:00. Early check-in tergantung ketersediaan."},
		{"CHECKOUT", "Check-out Time", "Check-out pukul 12:00. Late check-out dapat diminta ke Front Office."},
		{"RESTAURANT", "Hotel Restaurant", "Restaurant buka 06:00-22:00 menyajikan masakan Indonesia & Western."},
		{"ROOM", "Room Facilities", "Fasilitas kamar: AC, TV, WiFi, shower, amenities, brankas, minibar."},
		{"POLICY", "Smoking Policy", "Hotel bebas asap rokok di dalam kamar. Area merokok tersedia di luar lobby."},
		{"POLICY", "Pet Policy", "Hewan peliharaan tidak diperbolehkan di kamar."},
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
	}
	items := []rec{
		{"Warung Sederhana BATIQA", "restaurant", "Masakan Indonesia, budget-friendly dekat hotel", intPtr(25000), intPtr(75000), f64Ptr(0.5), "Jl. Contoh No.1"},
		{"Cafe Ceria", "cafe", "Kopi & pastry, cocok untuk meeting santai", intPtr(20000), intPtr(50000), f64Ptr(0.8), "Jl. Contoh No.2"},
		{"Mall Central", "shopping", "Pusat perbelanjaan terbesar sekitar hotel", intPtr(0), intPtr(0), f64Ptr(1.2), "Jl. Mall No.10"},
		{"Pantai Indah", "tourism", "Destinasi wisata pantai dekat hotel", intPtr(0), intPtr(30000), f64Ptr(3.5), "Jl. Pantai No.5"},
		{"ATM BCA Terdekat", "atm", "ATM 24 jam 200m dari lobby", intPtr(0), intPtr(0), f64Ptr(0.2), "Lobby BATIQA"},
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
			Active: true, CreatedAt: now,
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
