package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"batiqa-ai/internal/model"
)

// AssignmentRepository handles the ticket_assignments collection.
type AssignmentRepository struct {
	db *mongo.Database
}

func NewAssignmentRepository(db *mongo.Database) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

const assignmentsCol = "ticket_assignments"

// Assign creates a ticket assignment (staff to ticket).
// Returns an error if the ticket is already assigned to that staff (unique index).
func (r *AssignmentRepository) Assign(ticketID, staffID int64) (*model.TicketAssignment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := nextID(ctx, r.db, assignmentsCol)
	if err != nil {
		return nil, fmt.Errorf("assignment Assign id: %w", err)
	}
	a := &model.TicketAssignment{
		ID:         id,
		TicketID:   ticketID,
		StaffID:    staffID,
		AssignedAt: time.Now(),
	}
	if _, err := r.db.Collection(assignmentsCol).InsertOne(ctx, a); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("ticket already assigned to this staff")
		}
		return nil, fmt.Errorf("assignment Assign: %w", err)
	}
	return a, nil
}

// ListByTicket returns assignments for a ticket.
func (r *AssignmentRepository) ListByTicket(ticketID int64) ([]*model.TicketAssignment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.db.Collection(assignmentsCol).Find(
		ctx,
		bson.M{"ticket_id": ticketID},
		options.Find().SetSort(bson.M{"assigned_at": 1}),
	)
	if err != nil {
		return nil, fmt.Errorf("assignment ListByTicket: %w", err)
	}
	defer cursor.Close(ctx)

	var out []*model.TicketAssignment
	if err := cursor.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("assignment decode: %w", err)
	}
	return out, nil
}
