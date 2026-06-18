package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) loadSnackBarMenu() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/snackbar", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
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
		body := `{"item_id":"` + itemID + `","quantity":1}`
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/snackbar/order", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
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
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/snackbar/orders", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadSnackBarOrdersMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Orders []models.SnackBarOrder `json:"orders"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadSnackBarOrdersMsg{orders: r.Orders}
	}
}

func (m *Model) loadSnackBarManage() tea.Cmd {
	return func() tea.Msg {
		if m.userResp == nil || !bitmask.CanSnackBarManage(m.userResp.Tier) {
			return pages.ErrorMsg{Message: "⛔ ACCESS DENIED — SnackBar Manager or Owner required"}
		}
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/snackbar", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
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
		body := `{"item_id":"` + itemID + `","amount":5}`
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/snackbar/restock", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.snackBarManage.Error = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error   string `json:"error"`
				Message string `json:"message"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			if resp.StatusCode == http.StatusForbidden || strings.Contains(e.Error, "ACCESS DENIED") || strings.Contains(e.Error, "⛔") {
				return pages.ErrorMsg{Message: e.Error}
			}
			m.snackBarManage.Error = e.Error
			return nil
		}
		var r struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		m.snackBarManage.Status = "Restocked +5 ✓"
		return m.loadSnackBarManage()()
	}
}
