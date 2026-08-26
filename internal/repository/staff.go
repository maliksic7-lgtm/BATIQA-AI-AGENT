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

// StaffRepository handles the staff collection.
type StaffRepository struct {
	db *mongo.Database
}

func NewStaffRepository(db *mongo.Database) *StaffRepository {
	return &StaffRepository{db: db}
}

const staffCol = "staff"

func staffCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// FindByEmail returns staff by email.
func (r *StaffRepository) FindByEmail(email string) (*model.Staff, error) {
	ctx, cancel := staffCtx()
	defer cancel()

	var s model.Staff
	err := r.db.Collection(staffCol).FindOne(ctx, bson.M{"email": email}).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("staff FindByEmail: %w", err)
	}
	return &s, nil
}

// FindByID returns staff by id.
func (r *StaffRepository) FindByID(id int64) (*model.Staff, error) {
	ctx, cancel := staffCtx()
	defer cancel()

	var s model.Staff
	err := r.db.Collection(staffCol).FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("staff FindByID: %w", err)
	}
	return &s, nil
}

// Create inserts new staff (password_hash must be bcrypt).
func (r *StaffRepository) Create(s *model.Staff) error {
	if !model.IsValidStaffDepartment(s.Department) {
		return fmt.Errorf("invalid staff department: %s", s.Department)
	}
	ctx, cancel := staffCtx()
	defer cancel()

	id, err := nextID(ctx, r.db, staffCol)
	if err != nil {
		return fmt.Errorf("staff Create id: %w", err)
	}
	s.ID = id
	s.CreatedAt = time.Now()
	_, err = r.db.Collection(staffCol).InsertOne(ctx, s)
	if err != nil {
		return fmt.Errorf("staff Create: %w", err)
	}
	return nil
}

// List returns all staff optionally filtered by department.
func (r *StaffRepository) List(department *string) ([]*model.Staff, error) {
	filter := bson.M{}
	if department != nil && *department != "" {
		if !model.IsValidStaffDepartment(*department) {
			return nil, fmt.Errorf("invalid department filter: %s", *department)
		}
		filter["department"] = *department
	}
	ctx, cancel := staffCtx()
	defer cancel()

	cursor, err := r.db.Collection(staffCol).Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		return nil, fmt.Errorf("staff List: %w", err)
	}
	defer cursor.Close(ctx)

	var out []*model.Staff
	if err := cursor.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("staff decode: %w", err)
	}
	return out, nil
}
