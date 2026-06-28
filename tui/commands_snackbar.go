package tui

import (
	"encoding/json"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) loadSnackBarMenu() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/snackbar")
		if resp == nil {
			return loadSnackBarMenuMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Items []models.SnackBarItem `json:"items"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadSnackBarMenuMsg{items: r.Items}
	}
}

func (m *Model) doSnackBarOrder(itemID string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]interface{}{
			"item_id":  itemID,
			"quantity": 1,
		})
		resp, err := m.apiPost("/api/v1/snackbar/order", string(body))
		if err != nil {
			m.snackBarMenu.Error = err.Error()
			return nil
		}
		defer resp.Body.Close()
		var r struct {
			Error      string               `json:"error"`
			Message    string               `json:"message"`
			Order      models.SnackBarOrder `json:"order"`
			TotalSpent float64              `json:"total_spent"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if resp.StatusCode >= 400 || r.Error != "" {
			if r.Error == "" {
				r.Error = "order failed"
			}
			m.snackBarMenu.Error = r.Error
			return nil
		}
		m.snackBarMenu.Status = fmt.Sprintf("Ordered %s! $%.2f", r.Order.ItemName, r.Order.Total)
		if m.userResp != nil {
			m.userResp.Balance -= r.Order.Total
			m.snackBarMenu.Balance = m.userResp.Balance
		}
		return m.loadSnackBarMenu()()
	}
}

func (m *Model) loadSnackBarOrders() tea.Cmd {
	return func() tea.Msg {
		resp, _ := m.apiGet("/api/v1/snackbar/orders")
		if resp == nil {
			return loadSnackBarOrdersMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Orders []models.SnackBarOrder `json:"orders"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		for i, j := 0, len(r.Orders)-1; i < j; i, j = i+1, j-1 {
			r.Orders[i], r.Orders[j] = r.Orders[j], r.Orders[i]
		}
		return loadSnackBarOrdersMsg{orders: r.Orders}
	}
}

func (m *Model) loadSnackBarManage() tea.Cmd {
	return func() tea.Msg {
		if m.userResp == nil || !bitmask.CanSnackBarManage(m.userResp.Tier) {
			return pages.ErrorMsg{Message: "⛔ ACCESS DENIED — SnackBar Manager or Owner required"}
		}
		resp, _ := m.apiGet("/api/v1/snackbar")
		if resp == nil {
			return loadSnackBarManageMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Items []models.SnackBarItem `json:"items"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadSnackBarManageMsg{items: r.Items}
	}
}

func (m *Model) doSnackBarRestock(itemID string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]interface{}{
			"item_id": itemID,
			"amount":  5,
		})
		resp, err := m.apiPost("/api/v1/snackbar/restock", string(body))
		if err != nil {
			m.snackBarManage.Error = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if errMsg := handleAPIErr(resp); errMsg != nil {
			return errMsg
		}
		if resp.StatusCode >= 400 {
			var e struct {
				Error   string `json:"error"`
				Message string `json:"message"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			m.snackBarManage.Error = e.Error
			return nil
		}
		m.snackBarManage.Status = "Restocked +5 ✓"
		return m.loadSnackBarManage()()
	}
}
