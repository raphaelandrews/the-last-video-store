package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) doSearch(query string) tea.Cmd {
	if query == "" {
		return nil
	}
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/movies/search?q="+query, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return searchResultsMsg{}
		}
		defer resp.Body.Close()
		var results []models.MovieResponse
		json.NewDecoder(resp.Body).Decode(&results)
		return searchResultsMsg{results: results}
	}
}

func (m *Model) doSubmitTOTP(tempToken, code string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{"temp_token": tempToken, "code": code})
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/auth/login/totp", strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var r struct {
			AccessToken  string              `json:"access_token"`
			RefreshToken string              `json:"refresh_token"`
			Error        string              `json:"error"`
			User         models.UserResponse `json:"user"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if r.Error != "" || resp.StatusCode >= 400 {
			if r.Error == "" {
				r.Error = "invalid code"
			}
			m.totpCode = ""
			return pages.ErrorMsg{Message: r.Error}
		}
		return pages.LoginSuccessMsg{
			AccessToken:  r.AccessToken,
			RefreshToken: r.RefreshToken,
			User:         &r.User,
		}
	}
}

func (m *Model) doLogin(u, p string) tea.Cmd {
	return func() tea.Msg {
		body, _ := json.Marshal(map[string]string{"username": u, "password": p})
		resp, err := http.Post(m.baseURL+"/api/v1/auth/login", "application/json", strings.NewReader(string(body)))
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var r struct {
			AccessToken  string              `json:"access_token"`
			RefreshToken string              `json:"refresh_token"`
			Error        string              `json:"error"`
			User         models.UserResponse `json:"user"`
			TOTPRequired bool                `json:"totp_required"`
			TempToken    string              `json:"temp_token"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if r.TOTPRequired && r.TempToken != "" {
			m.tempToken = r.TempToken
			m.totpCode = ""
			m.screen = scrTOTP
			return nil
		}
		if r.Error != "" || resp.StatusCode >= 400 {
			if r.Error == "" {
				r.Error = "invalid credentials"
			}
			return pages.ErrorMsg{Message: r.Error}
		}
		return pages.LoginSuccessMsg{AccessToken: r.AccessToken, RefreshToken: r.RefreshToken, User: &r.User}
	}
}

func (m *Model) doRegister(u, p string) tea.Cmd {
	return func() tea.Msg {
		if len(u) < 3 {
			return pages.ErrorMsg{Message: "Username must be at least 3 characters"}
		}
		if len(p) < 6 {
			return pages.ErrorMsg{Message: "Password must be at least 6 characters"}
		}
		body, _ := json.Marshal(map[string]string{"username": u, "password": p})
		resp, err := http.Post(m.baseURL+"/api/v1/auth/register", "application/json", strings.NewReader(string(body)))
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var r struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		if r.Error != "" || resp.StatusCode >= 400 {
			if r.Error == "" {
				r.Error = "registration failed"
			}
			return pages.ErrorMsg{Message: r.Error}
		}
		m.register = nil
		return pages.NavigateMsg{Page: "login"}
	}
}

func (m *Model) loadMovies(page int) tea.Cmd {
	ps := m.browse.PageSize
	if ps <= 0 {
		ps = 20
	}
	return func() tea.Msg {
		url := fmt.Sprintf("%s/api/v1/movies?page_size=%d&page=%d", m.baseURL, ps, page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadMoviesMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Movies []models.MovieResponse `json:"movies"`
			Total  int                    `json:"total"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadMoviesMsg{movies: r.Movies, total: r.Total, page: page}
	}
}

func (m *Model) loadStaffPicks() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/movies/staff-picks", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadMoviesMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Movies []models.MovieResponse `json:"movies"`
			Total  int                    `json:"total"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadMoviesMsg{movies: r.Movies, total: r.Total, page: 1}
	}
}

func (m *Model) loadLastChance() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/movies/last-chance", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadMoviesMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Movies []models.MovieResponse `json:"movies"`
			Total  int                    `json:"total"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadMoviesMsg{movies: r.Movies, total: r.Total, page: 1}
	}
}

func (m *Model) loadAdminMovies(page int) tea.Cmd {
	ps := m.adminMovies.PageSize
	if ps <= 0 {
		ps = 30
	}
	return func() tea.Msg {
		url := fmt.Sprintf("%s/api/v1/movies?page_size=%d&page=%d", m.baseURL, ps, page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadAdminMoviesMsg{}
		}
		defer resp.Body.Close()
		var r struct {
			Movies []models.MovieResponse `json:"movies"`
			Total  int                    `json:"total"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		return loadAdminMoviesMsg{movies: r.Movies, total: r.Total, page: page}
	}
}

func (m *Model) loadAdminUsers() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadAdminUsersMsg{}
		}
		defer resp.Body.Close()
		var users []models.UserResponse
		json.NewDecoder(resp.Body).Decode(&users)
		return loadAdminUsersMsg{users: users}
	}
}

func (m *Model) doVerifyAuditChain() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/audit/verify", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.auditLog.VerifyMsg = "⚠️ Verification failed: " + err.Error()
			return nil
		}
		defer resp.Body.Close()
		var r struct {
			ChainIntact bool   `json:"chain_intact"`
			Message     string `json:"message"`
		}
		json.NewDecoder(resp.Body).Decode(&r)
		m.auditLog.VerifyMsg = r.Message
		return nil
	}
}

func (m *Model) loadAuditLog() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/audit", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadAuditLogMsg{}
		}
		defer resp.Body.Close()
		var entries []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&entries)
		return loadAuditLogMsg{entries: entries}
	}
}

func (m *Model) doRent(movieID string) tea.Cmd {
	return func() tea.Msg {
		useTicket := m.detail != nil && m.detail.UseTicket
		body := fmt.Sprintf(`{"movie_id":"%s","use_ticket":%v}`, movieID, useTicket)
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/rentals/rent", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct{ Error string }
			json.NewDecoder(resp.Body).Decode(&e)
			return pages.ErrorMsg{Message: e.Error}
		}
		var rental models.RentalResponse
		json.NewDecoder(resp.Body).Decode(&rental)
		if m.detail != nil {
			m.detail.SetRental(&rental)
		}
		if m.userResp != nil {
			if m.detail != nil && m.detail.UseTicket {
				m.userResp.FreeRentals--
				m.browse.Status = "Rented! " + rental.MovieTitle + " (🎟️ free rental)"
			} else {
				cost := models.MovieCost(0, rental.MovieFormat)
				m.userResp.Balance -= cost
				m.userResp.RentalCount++
				m.browse.Status = fmt.Sprintf("Rented! %s (💵 $%.2f)", rental.MovieTitle, cost)
			}
			m.detail.UseTicket = false
		}
		return m.loadMovies(m.browse.Page)()
	}
}

func (m *Model) doAddToWishlist(movieID string, fromDetail bool) tea.Cmd {
	return func() tea.Msg {
		body := `{"movie_id":"` + movieID + `"}`
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/wishlist", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if fromDetail {
				m.detail.ErrMsg = err.Error()
			} else {
				m.browse.Status = err.Error()
			}
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode == 409 {
			if fromDetail {
				m.detail.ErrMsg = "Already in waitlist"
			} else {
				m.browse.Status = "Already in waitlist"
			}
			return nil
		}
		if resp.StatusCode >= 400 {
			var e struct{ Error string }
			json.NewDecoder(resp.Body).Decode(&e)
			if fromDetail {
				m.detail.ErrMsg = e.Error
			} else {
				m.browse.Status = e.Error
			}
			return nil
		}
		if fromDetail {
			m.detail.StatusMsg = "Added to waitlist ✓"
		} else {
			m.browse.Status = "Added to waitlist ✓"
		}
		return nil
	}
}

func (m *Model) doReturn(rentalID string) tea.Cmd {
	return func() tea.Msg {
		body := `{"rental_id":"` + rentalID + `"}`
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/rentals/return", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
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
		}
		return tea.Batch(m.loadRentals(), m.loadMovies(m.browse.Page))()
	}
}

func (m *Model) loadWishlist() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/wishlist", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		var items []pages.WishlistItem
		json.NewDecoder(resp.Body).Decode(&items)
		return loadWishlistMsg{items: items}
	}
}

func (m *Model) doRemoveFromWishlist(movieID string) tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("DELETE", m.baseURL+"/api/v1/wishlist/"+movieID, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return pages.ErrorMsg{Message: err.Error()}
		}
		defer resp.Body.Close()
		m.wishlist.RemoveSelected()
		m.wishlist.Status = "Removed from wishlist"
		if resp.StatusCode >= 400 {
			var e struct{ Error string }
			json.NewDecoder(resp.Body).Decode(&e)
			m.wishlist.Status = e.Error
		}
		return nil
	}
}

func (m *Model) loadRentals() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/rentals/history", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp == nil {
			return loadRentalsMsg{}
		}
		defer resp.Body.Close()
		var rentals []models.RentalResponse
		json.NewDecoder(resp.Body).Decode(&rentals)
		return loadRentalsMsg{rentals: rentals}
	}
}

func (m *Model) loadProfile() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/rentals/history", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
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
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/merch", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
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
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/merch/redeem", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.merch.Error = err.Error()
			return nil
		}
		defer resp.Body.Close()

		var result struct {
			Error       string `json:"error"` // server sends "error" for errors
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

func (m *Model) doExtendRental(rentalID string) tea.Cmd {
	return func() tea.Msg {
		body := `{"rental_id":"` + rentalID + `"}`
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/rentals/extend", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.rentals.Status = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			m.rentals.Status = e.Error
			return nil
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

func (m *Model) doPurchaseTier(tierName string) tea.Cmd {
	return func() tea.Msg {
		body := `{"tier_name":"` + tierName + `"}`
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/tiers/purchase", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
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
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/inventory", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
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

func (m *Model) doCreateMovie(msg pages.MovieFormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		cast := parseCast(msg.Cast)
		body, _ := json.Marshal(map[string]interface{}{
			"title":        msg.Title,
			"year":         msg.Year,
			"genre":        msg.Genre,
			"format":       msg.Format,
			"director":     msg.Director,
			"cast":         cast,
			"synopsis":     msg.Synopsis,
			"copies_total": msg.Copies,
			"rental_price": msg.Price,
		})
		req, _ := http.NewRequest("POST", m.baseURL+"/api/v1/movies", strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.movieForm.ErrMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			m.movieForm.ErrMsg = e.Error
			return nil
		}
		m.moveToAdminMovies()
		return m.loadAdminMovies(m.adminMovies.Page)()
	}
}

func (m *Model) doUpdateMovie(msg pages.MovieFormSubmitMsg) tea.Cmd {
	return func() tea.Msg {
		cast := parseCast(msg.Cast)
		body, _ := json.Marshal(map[string]interface{}{
			"title":        msg.Title,
			"year":         msg.Year,
			"genre":        msg.Genre,
			"format":       msg.Format,
			"director":     msg.Director,
			"cast":         cast,
			"synopsis":     msg.Synopsis,
			"copies_total": msg.Copies,
			"rental_price": msg.Price,
		})
		req, _ := http.NewRequest("PUT", m.baseURL+"/api/v1/movies/"+msg.MovieID, strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.movieForm.ErrMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			m.movieForm.ErrMsg = e.Error
			return nil
		}
		m.moveToAdminMovies()
		return m.loadAdminMovies(m.adminMovies.Page)()
	}
}

func (m *Model) doUpdateUser(userID, action string) tea.Cmd {
	return func() tea.Msg {
		var body string
		u := m.adminUsers.SelectedUser()
		if action == "promote" {
			next := nextTier(u.TierName)
			body = `{"tier":"` + next + `"}`
		} else if action == "demote" {
			prev := prevTier(u.TierName)
			body = `{"tier":"` + prev + `"}`
		} else if action == "ban" {
			if u.Banned {
				body = `{"banned":false}`
			} else {
				body = `{"banned":true}`
			}
		}
		req, _ := http.NewRequest("PUT", m.baseURL+"/api/v1/users/"+userID, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			m.adminUsers.ErrMsg = err.Error()
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			var e struct {
				Error string `json:"error"`
			}
			json.NewDecoder(resp.Body).Decode(&e)
			m.adminUsers.ErrMsg = e.Error
			return nil
		}
		m.adminUsers.ErrMsg = ""
		return m.loadAdminUsers()()
	}
}

var tierOrder = []string{"Bronze", "Silver", "Gold", "Employee", "Supervisor", "Manager", "Owner"}

func nextTier(current string) string {
	for i, t := range tierOrder {
		if t == current && i+1 < len(tierOrder) {
			return strings.ToLower(tierOrder[i+1])
		}
	}
	return strings.ToLower(current)
}

func prevTier(current string) string {
	for i, t := range tierOrder {
		if t == current && i > 0 {
			return strings.ToLower(tierOrder[i-1])
		}
	}
	return strings.ToLower(current)
}

func (m *Model) doToggleStaffPick(movieID string, current bool) tea.Cmd {
	method := "POST"
	if current {
		method = "DELETE"
	}
	return func() tea.Msg {
		req, _ := http.NewRequest(method, m.baseURL+"/api/v1/movies/"+movieID+"/staff-pick", nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
		return m.loadAdminMovies(m.adminMovies.Page)()
	}
}

func (m *Model) doDeleteMovie(movieID string) tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("DELETE", m.baseURL+"/api/v1/movies/"+movieID, nil)
		req.Header.Set("Authorization", "Bearer "+m.token)
		resp, _ := http.DefaultClient.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
		return m.loadAdminMovies(m.adminMovies.Page)()
	}
}

func (m *Model) moveToAdminMovies() {
	m.screen = scrAdminMovies
	m.movieForm = nil
}

func parseCast(s string) []string {
	if s == "" {
		return nil
	}
	var cast []string
	for _, c := range splitCSV(s) {
		c = strings.TrimSpace(c)
		if c != "" {
			cast = append(cast, c)
		}
	}
	return cast
}

func splitCSV(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
