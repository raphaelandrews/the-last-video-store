package pages

import (
	"strings"
	"testing"

	"github.com/thelastvideostore/internal/models"
)

func TestMyRentalsRenderAfterReturn(t *testing.T) {
	r := models.RentalResponse{
		ID:           "abc123",
		MovieTitle:   "Pulp Fiction",
		MovieFormat:  "VHS",
		RentedAt:     1782534943,
		DueDate:      1782535003,
		ReturnedAt:   1782535013,
		LateFee:      0.2,
		RewindFee:    1.0,
		Status:       "returned",
		PointsEarned: -5,
	}

	m := NewMyRentalsModel()
	m.SetRentals([]models.RentalResponse{r})

	view := m.View(120, 30)
	t.Logf("VIEW:\n%s", view)

	if !strings.Contains(view, "returned") {
		t.Error("view does not contain 'returned' status")
	}
	if !strings.Contains(view, "-5") {
		t.Error("view does not contain '-5' points")
	}
	if strings.Contains(view, "overdue") {
		t.Error("view still shows 'overdue' for returned rental")
	}
}
