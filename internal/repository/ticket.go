package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"batiqa-ai/internal/model"
)

// TicketRepository handles the tickets collection.
type TicketRepository struct {
	db *mongo.Database
}

func NewTicketRepository(db *mongo.Database) *TicketRepository {
	return &TicketRepository{db: db}
}

const ticketsCol = "tickets"

func ticketCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// Create inserts a ticket and generates ticket_number as TKT-{id} (zero-padded 6 digits).
// It validates enums before insert per DATABASE.md rules.
func (r *TicketRepository) Create(t *model.Ticket) error {
	if !model.IsValidDepartment(t.Department) {
		return fmt.Errorf("invalid department: %s", t.Department)
	}
	if !model.IsValidPriority(t.Priority) {
		return fmt.Errorf("invalid priority: %s", t.Priority)
	}
	if !model.IsValidStatus(t.Status) {
		return fmt.Errorf("invalid status: %s", t.Status)
	}
	if strings.TrimSpace(t.RoomNumber) == "" {
		return fmt.Errorf("room_number is required")
	}
	if strings.TrimSpace(t.Category) == "" {
		return fmt.Errorf("category is required")
	}
	if strings.TrimSpace(t.Description) == "" {
		return fmt.Errorf("description is required")
	}

	ctx, cancel := ticketCtx()
	defer cancel()

	id, err := nextID(ctx, r.db, ticketsCol)
	if err != nil {
		return fmt.Errorf("ticket Create id: %w", err)
	}
	t.ID = id
	t.TicketNumber = fmt.Sprintf("TKT-%06d", id)
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	if t.Status == "" {
		t.Status = model.StatusOpen
	}
	if t.Priority == "" {
		t.Priority = model.PriorityMedium
	}

	if _, err := r.db.Collection(ticketsCol).InsertOne(ctx, t); err != nil {
		return fmt.Errorf("ticket Create insert: %w", err)
	}
	return nil
}

// FindByTicketNumber returns ticket by ticket_number.
func (r *TicketRepository) FindByTicketNumber(ticketNumber string) (*model.Ticket, error) {
	ctx, cancel := ticketCtx()
	defer cancel()

	var t model.Ticket
	err := r.db.Collection(ticketsCol).FindOne(ctx, bson.M{"ticket_number": ticketNumber}).Decode(&t)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("ticket FindByTicketNumber: %w", err)
	}
	return &t, nil
}

// FindByID returns ticket by id.
func (r *TicketRepository) FindByID(id int64) (*model.Ticket, error) {
	ctx, cancel := ticketCtx()
	defer cancel()

	var t model.Ticket
	err := r.db.Collection(ticketsCol).FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("ticket FindByID: %w", err)
	}
	return &t, nil
}

// List returns tickets with optional filters, ordered by created_at DESC.
func (r *TicketRepository) List(f model.TicketFilter) ([]*model.Ticket, error) {
	filter := bson.M{}
	if f.Department != nil && *f.Department != "" {
		if !model.IsValidDepartment(*f.Department) {
			return nil, fmt.Errorf("invalid department filter: %s", *f.Department)
		}
		filter["department"] = *f.Department
	}
	if f.Status != nil && *f.Status != "" {
		if !model.IsValidStatus(*f.Status) {
			return nil, fmt.Errorf("invalid status filter: %s", *f.Status)
		}
		filter["status"] = *f.Status
	}
	if f.Priority != nil && *f.Priority != "" {
		if !model.IsValidPriority(*f.Priority) {
			return nil, fmt.Errorf("invalid priority filter: %s", *f.Priority)
		}
		filter["priority"] = *f.Priority
	}
	if f.RoomNumber != nil && *f.RoomNumber != "" {
		filter["room_number"] = *f.RoomNumber
	}

	opts := options.Find().SetSort(bson.M{"created_at": -1})
	if f.Limit > 0 {
		limit := int64(f.Limit)
		if limit > 100 {
			limit = 100
		}
		opts.SetLimit(limit)
		if f.Offset > 0 {
			opts.SetSkip(int64(f.Offset))
		}
	}

	ctx, cancel := ticketCtx()
	defer cancel()

	cursor, err := r.db.Collection(ticketsCol).Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("ticket List: %w", err)
	}
	defer cursor.Close(ctx)

	var out []*model.Ticket
	if err := cursor.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("ticket List decode: %w", err)
	}
	return out, nil
}

// UpdateStatus updates ticket status with validation and sets/clears resolved_at.
func (r *TicketRepository) UpdateStatus(ticketNumber, newStatus string) (*model.Ticket, error) {
	if !model.IsValidStatus(newStatus) {
		return nil, fmt.Errorf("invalid status: %s", newStatus)
	}
	current, err := r.FindByTicketNumber(ticketNumber)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("ticket not found: %s", ticketNumber)
	}
	if !model.IsValidStatusTransition(current.Status, newStatus) {
		return nil, fmt.Errorf("invalid status transition %s -> %s", current.Status, newStatus)
	}

	set := bson.M{"status": newStatus, "updated_at": time.Now()}
	if newStatus == model.StatusResolved {
		now := time.Now()
		set["resolved_at"] = now
	} else {
		set["resolved_at"] = nil // clear when leaving RESOLVED
	}

	ctx, cancel := ticketCtx()
	defer cancel()

	var t model.Ticket
	err = r.db.Collection(ticketsCol).FindOneAndUpdate(
		ctx,
		bson.M{"ticket_number": ticketNumber},
		bson.M{"$set": set},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&t)
	if err != nil {
		return nil, fmt.Errorf("ticket UpdateStatus: %w", err)
	}
	return &t, nil
}

// UpdatePriority updates ticket priority.
func (r *TicketRepository) UpdatePriority(ticketNumber, newPriority string) (*model.Ticket, error) {
	if !model.IsValidPriority(newPriority) {
		return nil, fmt.Errorf("invalid priority: %s", newPriority)
	}
	ctx, cancel := ticketCtx()
	defer cancel()

	var t model.Ticket
	err := r.db.Collection(ticketsCol).FindOneAndUpdate(
		ctx,
		bson.M{"ticket_number": ticketNumber},
		bson.M{"$set": bson.M{"priority": newPriority, "updated_at": time.Now()}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&t)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("ticket not found: %s", ticketNumber)
		}
		return nil, fmt.Errorf("ticket UpdatePriority: %w", err)
	}
	return &t, nil
}

// TicketStats holds dashboard stats
type TicketStats struct {
	Total         int64
	Open          int64
	High          int64
	Housekeeping  int64
	Engineering   int64
	ResolvedToday int64
}

// GetStats returns dashboard statistics for staff overview.
func (r *TicketRepository) GetStats() (*TicketStats, error) {
	ctx, cancel := ticketCtx()
	defer cancel()

	col := r.db.Collection(ticketsCol)
	active := bson.M{"status": bson.M{"$nin": []string{model.StatusResolved, model.StatusCancelled}}}
	stats := &TicketStats{}

	type countJob struct {
		filter bson.M
		dst    *int64
	}
	yesterday := time.Now().Truncate(24 * time.Hour)
	jobs := []countJob{
		{bson.M{}, &stats.Total},
		{bson.M{"status": model.StatusOpen}, &stats.Open},
		{merge(active, bson.M{"priority": model.PriorityHigh}), &stats.High},
		{merge(active, bson.M{"department": model.DeptHousekeeping}), &stats.Housekeeping},
		{merge(active, bson.M{"department": model.DeptEngineering}), &stats.Engineering},
		{bson.M{"status": model.StatusResolved, "resolved_at": bson.M{"$gte": yesterday}}, &stats.ResolvedToday},
	}
	for _, j := range jobs {
		n, err := col.CountDocuments(ctx, j.filter)
		if err != nil {
			return nil, fmt.Errorf("stats count: %w", err)
		}
		*j.dst = n
	}
	return stats, nil
}

// merge combines two filter documents into one.
func merge(base, extra bson.M) bson.M {
	out := bson.M{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
