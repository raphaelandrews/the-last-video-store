package tui

import (
	"strings"
	"testing"

	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func TestMyRentalsE2EAfterReturn(t *testing.T) {
	m := NewModel("http://localhost:18080")
	m.userResp = &models.UserResponse{ID: "seed-manager", Username: "manager"}
	m.token = "dummy"
	m.screen = scrRentals
	m.ready = true
	m.w = 120
	m.h = 40
	m.rentals = pages.NewMyRentalsModel()

	oldRental := models.RentalResponse{
		ID:           "abc123",
		MovieTitle:   "Pulp Fiction",
		MovieFormat:  "VHS",
		RentedAt:     1000,
		DueDate:      2000,
		Status:       "overdue",
		PointsEarned: 0,
	}
	m.rentals.SetRentals([]models.RentalResponse{oldRental})

	view1 := m.View()
	t.Logf("BEFORE RETURN (full):\n%s", view1)
	if !strings.Contains(view1, "overdue") {
		t.Error("expected overdue in view before return")
	}

	newRental := models.RentalResponse{
		ID:           "abc123",
		MovieTitle:   "Pulp Fiction",
		MovieFormat:  "VHS",
		RentedAt:     1000,
		DueDate:      2000,
		Status:       "returned",
		PointsEarned: -5,
	}
	_, _ = m.Update(loadRentalsMsg{rentals: []models.RentalResponse{newRental}})

	view2 := m.View()
	t.Logf("AFTER RETURN (full):\n%s", view2)

	if strings.Contains(view2, "overdue") {
		t.Error("after return, view still shows 'overdue'")
	}
	if !strings.Contains(view2, "returned") {
		t.Error("after return, view does not show 'returned'")
	}
	if !strings.Contains(view2, "-5") {
		t.Error("after return, view does not show '-5' points")
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
