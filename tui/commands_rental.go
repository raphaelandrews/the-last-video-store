package tui

import (
	"encoding/json"

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
		m.rentals.Status = "Returned!"
		if m.userResp != nil {
			m.userResp.RentalCount--
			if rental.LateFee == 0 && rental.RewindFee == 0 {
				m.userResp.PopcornPoints += 10
				m.rentals.Status += " (+10🍿)"
			} else {
				m.userResp.PopcornPoints -= 5
			}
		}
		return tea.Batch(m.loadRentals(), m.loadMovies(m.browse.Page, m.browse.Genre))()
	}
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

		var sessions []models.GameSession
		if m.userResp != nil {
			sessResp, _ := m.apiGet("/api/v1/games/my-sessions")
			if sessResp != nil {
				defer sessResp.Body.Close()
				var r struct {
					Sessions []models.GameSession `json:"sessions"`
				}
				json.NewDecoder(sessResp.Body).Decode(&r)
				for i := range r.Sessions {
					if r.Sessions[i].Status == "active" {
						sessions = append(sessions, r.Sessions[i])
					}
				}
			}
		}
		return loadRentalsMsg{rentals: rentals, sessions: sessions}
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
