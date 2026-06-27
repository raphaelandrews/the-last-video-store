package tui

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
)

func (m *Model) doReturn(rentalID string) tea.Cmd {
	return func() tea.Msg {
		body := `{"rental_id":"` + rentalID + `"}`
		resp, err := m.apiPost("/api/v1/rentals/return", body)
		if err != nil {
			m.rentals.Status = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		var rental models.RentalResponse
		json.NewDecoder(resp.Body).Decode(&rental)
		m.rentals.Status = formatReturnStatus(rental)
		return tea.Batch(m.loadRentals(), m.loadMovies(m.browse.Page, m.browse.Genre))()
	}
}

func formatReturnStatus(r models.RentalResponse) string {
	parts := []string{"✓ Returned"}
	totalFee := r.LateFee + r.RewindFee
	if totalFee > 0 {
		parts = append(parts, fmt.Sprintf("-$%.2f fees", totalFee))
	}
	if r.PointsEarned > 0 {
		parts = append(parts, fmt.Sprintf("+%d🍿", r.PointsEarned))
	} else if r.PointsEarned < 0 {
		parts = append(parts, fmt.Sprintf("%d🍿", r.PointsEarned))
	}
	return strings.Join(parts, " · ")
}

func (m *Model) loadRentals() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/rentals/history")
		if resp == nil {
			return loadRentalsMsg{}
		}
		defer resp.Body.Close()
		var rentals []models.RentalResponse
		json.NewDecoder(resp.Body).Decode(&rentals)
		// Sort most-recent-first. The backend does not guarantee a
		// stable order, and reversing the slice would only be correct
		// if the backend already returned oldest-first. Sorting
		// explicitly on RentedAt makes the order robust to backend
		// changes.
		sort.SliceStable(rentals, func(i, j int) bool {
			return rentals[i].RentedAt > rentals[j].RentedAt
		})
		return loadRentalsMsg{rentals: rentals}
	}
}

func (m *Model) loadMyPlaySessions() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/games/my-sessions")
		if resp == nil {
			return loadMyPlaySessionsMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Sessions []models.GameSession `json:"sessions"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadMyPlaySessionsMsg{sessions: r.Sessions}
	}
}

func (m *Model) doExtendRental(rentalID string) tea.Cmd {
	return func() tea.Msg {
		body := `{"rental_id":"` + rentalID + `"}`
		resp, err := m.apiPost("/api/v1/rentals/extend", body)
		if err != nil {
			m.rentals.Status = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		var r struct {
			Message string `json:"message"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		m.rentals.Status = r.Message
		if m.userResp != nil {
			m.userResp.PopcornPoints -= 30
		}
		return m.loadRentals()()
	}
}
