package repository

import (
	"database/sql"
	"fmt"

	"batiqa-ai/internal/model"
)

// StaffRepository handles staff table.
type StaffRepository struct {
	db *sql.DB
}

func NewStaffRepository(db *sql.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

// FindByEmail returns staff by email.
func (r *StaffRepository) FindByEmail(email string) (*model.Staff, error) {
	query := `SELECT id, name, email, password_hash, department, created_at FROM staff WHERE email = ?`
	row := r.db.QueryRow(query, email)
	var s model.Staff
	err := row.Scan(&s.ID, &s.Name, &s.Email, &s.PasswordHash, &s.Department, &s.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("staff FindByEmail: %w", err)
	}
	return &s, nil
}

// FindByID returns staff by id.
func (r *StaffRepository) FindByID(id int64) (*model.Staff, error) {
	query := `SELECT id, name, email, password_hash, department, created_at FROM staff WHERE id = ?`
	row := r.db.QueryRow(query, id)
	var s model.Staff
	err := row.Scan(&s.ID, &s.Name, &s.Email, &s.PasswordHash, &s.Department, &s.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("staff FindByID: %w", err)
	}
	return &s, nil
}

// Create inserts new staff (parameterized, password_hash must be bcrypt).
func (r *StaffRepository) Create(s *model.Staff) error {
	if !model.IsValidStaffDepartment(s.Department) {
		return fmt.Errorf("invalid staff department: %s", s.Department)
	}
	query := `INSERT INTO staff (name, email, password_hash, department) VALUES (?, ?, ?, ?)`
	res, err := r.db.Exec(query, s.Name, s.Email, s.PasswordHash, s.Department)
	if err != nil {
		return fmt.Errorf("staff Create: %w", err)
	}
	id, _ := res.LastInsertId()
	s.ID = id
	return nil
}

// List returns all staff optionally filtered by department.
func (r *StaffRepository) List(department *string) ([]*model.Staff, error) {
	base := `SELECT id, name, email, password_hash, department, created_at FROM staff WHERE 1=1`
	args := []interface{}{}
	if department != nil && *department != "" {
		if !model.IsValidStaffDepartment(*department) {
			return nil, fmt.Errorf("invalid department filter: %s", *department)
		}
		base += ` AND department = ?`
		args = append(args, *department)
	}
	base += ` ORDER BY created_at ASC`
	rows, err := r.db.Query(base, args...)
	if err != nil {
		return nil, fmt.Errorf("staff List: %w", err)
	}
	defer rows.Close()

	var out []*model.Staff
	for rows.Next() {
		var s model.Staff
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.PasswordHash, &s.Department, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("staff scan: %w", err)
		}
		out = append(out, &s)
	}
	return out, rows.Err()
}
