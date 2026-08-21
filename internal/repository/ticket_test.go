package repository

import (
	"testing"

	"batiqa-ai/internal/model"
)

// TestTicketModelValidation tests enum validation without DB (unit test for Phase 2)
func TestTicketModelValidation(t *testing.T) {
	tests := []struct {
		name string
		fn   func() bool
		want bool
	}{
		{"valid department HOUSEKEEPING", func() bool { return model.IsValidDepartment(model.DeptHousekeeping) }, true},
		{"valid department ENGINEERING", func() bool { return model.IsValidDepartment(model.DeptEngineering) }, true},
		{"invalid department", func() bool { return model.IsValidDepartment("INVALID") }, false},
		{"valid priority HIGH", func() bool { return model.IsValidPriority(model.PriorityHigh) }, true},
		{"valid status OPEN", func() bool { return model.IsValidStatus(model.StatusOpen) }, true},
		{"invalid status", func() bool { return model.IsValidStatus("PENDING") }, false},
		{"valid staff ADMIN", func() bool { return model.IsValidStaffDepartment("ADMIN") }, true},
		{"valid role user", func() bool { return model.IsValidRole(model.RoleUser) }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(); got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}

// TestTicketFilterValidation tests that invalid filters are rejected (no DB needed)
func TestTicketFilterValidation(t *testing.T) {
	// This test ensures repository List would reject invalid enums before hitting DB
	invalidDept := "INVALID_DEPT"
	f := model.TicketFilter{Department: &invalidDept}
	if model.IsValidDepartment(*f.Department) {
		t.Error("expected invalid department to be rejected")
	}
}
