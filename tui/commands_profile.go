package tui

import (
	"encoding/json"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) loadProfile() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/rentals/history")
		if resp == nil {
			return loadProfileMsg{}
		}
		defer resp.Body.Close()
		var rentals []models.RentalResponse
		json.NewDecoder(resp.Body).Decode(&rentals)
		var late, rewind float64
		for _, r := range rentals {
			late += r.LateFee
			rewind += r.RewindFee
		}
		return loadProfileMsg{stats: &pages.RentalStats{Total: len(rentals), LateFee: late, Rewind: rewind}}
	}
}

func (m *Model) loadMerch() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/merch")
		if resp == nil {
			return loadMerchMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Items []models.MerchItem `json:"items"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadMerchMsg{items: r.Items}
	}
}

func (m *Model) doRedeemMerch(itemID string) tea.Cmd {
	return func() tea.Msg {
		body := `{"item_id":"` + itemID + `"}`
		resp, err := m.apiPost("/api/v1/merch/redeem", body)
		if err != nil {
			m.merch.Error = err.Error()
			return nil
		}
		defer resp.Body.Close()

		var result struct {
			Error       string `json:"error"`
			Message     string `json:"message"`
			PointsSpent int    `json:"points_spent"`
		}
		json.NewDecoder(resp.Body).Decode(&result)

		if resp.StatusCode >= 400 || result.Error != "" {
			m.merch.Error = result.Error
			return nil
		}

		m.merch.Error = ""
		m.merch.Status = "Redeemed! 🎉"
		if m.userResp != nil {
			m.userResp.PopcornPoints -= result.PointsSpent
			m.merch.Points = m.userResp.PopcornPoints
		}
		return m.loadMerch()()
	}
}

func (m *Model) doPurchaseTier(tierName string) tea.Cmd {
	return func() tea.Msg {
		body := `{"tier_name":"` + tierName + `"}`
		resp, err := m.apiPost("/api/v1/tiers/purchase", body)
		if err != nil {
			m.tierShop.Error = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			m.tierShop.Error = e.Error
			return nil
		}
		var r struct {
			Message string `json:"message"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		m.tierShop.Error = ""
		m.tierShop.Status = r.Message
		tier := models.TierByName(tierName)
		if tier != nil {
			m.tierShop.Current = tier.Name
			if m.userResp != nil {
				m.userResp.Balance -= tier.Price
				m.userResp.Subscription = tier.Name
				m.userResp.MaxRentals = tier.MaxConcurrent
				m.userResp.FreeRentals = tier.FreeRentals
				m.tierShop.Balance = m.userResp.Balance
			}
		}
		return nil
	}
}

func (m *Model) loadInventory() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/inventory")
		if resp == nil {
			return loadInventoryMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Items []pages.InventoryItem `json:"items"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadInventoryMsg{items: r.Items}
	}
}

func (m *Model) doTopUp() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.apiPostEmpty("/api/v1/users/me/topup")
		if err != nil {
			m.profile.StatusMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		var r struct {
			Error      string  `json:"error"`
			Message    string  `json:"message"`
			Amount     float64 `json:"amount"`
			NewBalance float64 `json:"new_balance"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if resp.StatusCode >= 400 || r.Error != "" {
			m.profile.StatusMsg = r.Error
			return nil
		}
		if m.userResp != nil {
			m.userResp.Balance = r.NewBalance
			m.userResp.LastTopUpAt = time.Now().Unix()
		}
		m.profile.StatusMsg = fmt.Sprintf("💰 %s — new balance: $%.2f", r.Message, r.NewBalance)
		return tea.Sequence(
			func() tea.Msg { return wishlistResultMsg{} },
			m.doRefreshMe(),
		)
	}
}
