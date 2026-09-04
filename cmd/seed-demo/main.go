// Command seed-demo populates realistic operational demo data so the staff
// dashboard, analytics, and live feed look alive during presentations.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"batiqa-ai/internal/config"
)

type demoTicket struct {
	Room       string
	Department string
	Category   string
	Desc       string
	Priority   string
	Status     string
	DaysAgo    int
	Hour       int
	ResolveHrs float64 // 0 = not resolved
}

func main() {
	force := flag.Bool("force", false, "seed even if tickets already exist")
	flag.Parse()

	cfg := config.Load()
	db, closeDB, err := config.ConnectMongo(cfg)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer closeDB()
	ctx := context.Background()
	tickets := db.Collection("tickets")

	n, _ := tickets.CountDocuments(ctx, bson.M{})
	if n > 10 && !*force {
		fmt.Printf("Tickets already exist (%d) - skip seeding (use -force to override)\n", n)
		return
	}
	if *force {
		tickets.DeleteMany(ctx, bson.M{})
		db.Collection("ticket_assignments").DeleteMany(ctx, bson.M{})
	}

	demo := []demoTicket{
		{"305", "ENGINEERING", "AC_PROBLEM", "AC kamar tidak dingin sama sekali sejak siang", "HIGH", "RESOLVED", 6, 9, 3.5},
		{"308", "HOUSEKEEPING", "TOWEL_REQUEST", "Tolong antar 2 handuk tambahan", "MEDIUM", "RESOLVED", 6, 11, 1.2},
		{"410", "ENGINEERING", "SHOWER_PROBLEM", "Air shower tidak panas", "HIGH", "RESOLVED", 5, 8, 4.0},
		{"302", "HOUSEKEEPING", "ROOM_CLEANING_REQUEST", "Minta cleaning saat tamu keluar jam 13.00", "MEDIUM", "RESOLVED", 5, 12, 2.5},
		{"305", "ENGINEERING", "TV_PROBLEM", "Remote TV tidak merespon", "MEDIUM", "RESOLVED", 4, 15, 2.0},
		{"301", "HOUSEKEEPING", "AMENITY_REQUEST", "Tambah amenities: sabun & shampoo", "LOW", "RESOLVED", 4, 19, 1.0},
		{"308", "ENGINEERING", "LIGHT_PROBLEM", "Lampu kamar berkedip", "MEDIUM", "RESOLVED", 3, 10, 3.0},
		{"410", "HOUSEKEEPING", "TOWEL_REQUEST", "Handuk basah tidak diganti semalam", "MEDIUM", "RESOLVED", 3, 14, 1.5},
		{"302", "ENGINEERING", "PLUMBING_PROBLEM", "Bak toilet mampet parah", "HIGH", "RESOLVED", 2, 7, 5.5},
		{"301", "ENGINEERING", "WIFI_PROBLEM", "WiFi sering putus dari laptop", "MEDIUM", "RESOLVED", 2, 16, 2.8},
		{"305", "HOUSEKEEPING", "ROOM_CLEANING_REQUEST", "Cleaning rutin harian", "MEDIUM", "IN_PROGRESS", 1, 9, 0},
		{"308", "ENGINEERING", "AC_PROBLEM", "AC kurang dingin, thermostat tidak merespon", "HIGH", "IN_PROGRESS", 1, 11, 0},
		{"302", "HOUSEKEEPING", "AMENITY_REQUEST", "Minta air mineral tambahan 2 botol", "LOW", "IN_PROGRESS", 1, 17, 0},
		{"410", "ENGINEERING", "GENERAL_MAINTENANCE", "Pintu lemari berderit keras", "LOW", "OPEN", 0, 8, 0},
		{"305", "HOUSEKEEPING", "TOWEL_REQUEST", "Tolong ganti handuk 3 buah, ada tamu tambahan", "MEDIUM", "OPEN", 0, 12, 0},
		{"301", "ENGINEERING", "AC_PROBLEM", "AC mati total, ruangan panas sekali", "HIGH", "OPEN", 0, 13, 0},
		{"308", "HOUSEKEEPING", "ROOM_CLEANING_REQUEST", "Turunkan service, tamu keluar rapat", "MEDIUM", "OPEN", 0, 14, 0},
	}

	inserted := 0
	for i, t := range demo {
		seq := int64(i + 1)
		number := fmt.Sprintf("TKT-%06d", seq)
		day := time.Now().AddDate(0, 0, -t.DaysAgo)
		minute := (i*13 + 7) % 60
		created := time.Date(day.Year(), day.Month(), day.Day(), t.Hour, minute, 0, 0, day.Location())
		doc := bson.M{
			"_id":           seq,
			"ticket_number": number,
			"room_number":   t.Room,
			"department":    t.Department,
			"category":      t.Category,
			"description":   t.Desc,
			"priority":      t.Priority,
			"status":        t.Status,
			"created_at":    created,
			"updated_at":    created,
		}
		if t.Status == "RESOLVED" {
			resolvedAt := created.Add(time.Duration(t.ResolveHrs * float64(time.Hour)))
			doc["resolved_at"] = resolvedAt
			doc["updated_at"] = resolvedAt
		}
		if _, err := tickets.InsertOne(ctx, doc); err != nil {
			log.Fatalf("insert %s: %v", number, err)
		}
		inserted++
	}

	// Sync counter past the highest seeded id so live-created tickets continue after TKT-000017
	counters := db.Collection("counters")
	counters.UpdateOne(ctx,
		bson.M{"_id": "tickets", "seq": bson.M{"$lt": int64(len(demo))}},
		bson.M{"$set": bson.M{"seq": int64(len(demo))}}, nil)

	fmt.Printf("Seeded %d operational demo tickets across %d rooms\n", inserted, 4)
	fmt.Println("Dashboard, analytics & SSE are now alive - open /staff/")
}
