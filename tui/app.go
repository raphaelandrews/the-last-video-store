package tui

import (
	"encoding/json"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/components"
	"github.com/thelastvideostore/tui/pages"
	"github.com/thelastvideostore/tui/styles"
)

type screen int

const (
	scrSplash screen = iota
	scrLogin
	scrRegister
	scrBrowse
	scrDetail
	scrRentals
	scrProfile
)

type loadMoviesMsg struct {
	movies []models.MovieResponse
	total  int
}
type loadRentalsMsg struct {
	rentals []models.RentalResponse
}
type loadProfileMsg struct {
	stats *pages.RentalStats
}

type Model struct {
	baseURL  string
	screen   screen
	w, h     int
	ready    bool
	token    string
	userResp *models.UserResponse

	splash   *pages.SplashModel
	login    *pages.LoginModel
	register *pages.RegisterModel
	browse   *pages.BrowseModel
	detail   *pages.MovieDetailModel
	rentals  *pages.MyRentalsModel
	profile  *pages.ProfileModel
	header   *components.HeaderModel
}

func NewModel(baseURL string) *Model {
	return &Model{
		baseURL: baseURL,
		screen:  scrSplash,
		splash:  pages.NewSplashModel(),
		login:   pages.NewLoginModel(),
		header:  components.NewHeaderModel(),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen, m.splash.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var pageCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" || k == "ctrl+d" {
			return m, tea.Quit
		}
		if k == "q" && m.screen != scrLogin && m.screen != scrSplash && m.screen != scrBrowse {
			if m.screen == scrDetail || m.screen == scrRentals || m.screen == scrProfile || m.screen == scrRegister {
				m.screen = scrBrowse
				return m, nil
			}
		}
		if k == "ctrl+r" && (m.screen == scrLogin || m.screen == scrRegister) {
			return m, func() tea.Msg { return pages.NavigateMsg{Page: "register"} }
		}
		if k == "ctrl+l" && m.screen == scrRegister {
			return m, func() tea.Msg { return pages.NavigateMsg{Page: "login"} }
		}
		pageCmd = m.pageKey(msg)

	default:
		switch msg := msg.(type) {
		case pages.SplashDoneMsg:
			m.screen = scrLogin
			return m, m.login.Init()

		case pages.LoginRequestMsg:
			return m, m.doLogin(msg.Username, msg.Password)

		case pages.RegisterRequestMsg:
			return m, m.doRegister(msg.Username, msg.Password)

		case pages.LoginSuccessMsg:
			m.token = msg.AccessToken
			m.userResp = msg.User
			m.browse = pages.NewBrowseModel()
			m.screen = scrBrowse
			return m, m.loadMovies()

		case pages.NavigateMsg:
			switch msg.Page {
			case "login":
				m.token = ""
				m.userResp = nil
				m.login = pages.NewLoginModel()
				m.screen = scrLogin
				return m, m.login.Init()
			case "register":
				m.register = pages.NewRegisterModel()
				m.screen = scrRegister
				return m, m.register.Init()
			case "browse":
				m.screen = scrBrowse
			case "rentals":
				m.rentals = pages.NewMyRentalsModel()
				m.screen = scrRentals
				return m, m.loadRentals()
			case "profile":
				m.profile = pages.NewProfileModel(m.userResp)
				m.screen = scrProfile
				return m, m.loadProfile()
			}

		case pages.RentRequestMsg:
			return m, m.doRent(msg.MovieID)
		case pages.ReturnRequestMsg:
			return m, m.doReturn(msg.RentalID)

		case loadMoviesMsg:
			m.browse.SetMovies(msg.movies, msg.total)
			if m.detail != nil && m.detail.Movie != nil {
				for i := range msg.movies {
					if msg.movies[i].ID == m.detail.Movie.ID {
						m.detail.Movie = &msg.movies[i]
						break
					}
				}
			}
		case pages.BrowseReloadMsg:
			return m, m.loadMovies()
		case loadRentalsMsg:
			m.rentals.SetRentals(msg.rentals)
		case loadProfileMsg:
			m.profile.SetStats(msg.stats)

		case pages.ErrorMsg:
			if m.login != nil {
				m.login.SetError(msg.Message)
			}
			if m.register != nil {
				m.register.SetError(msg.Message)
			}
			if m.detail != nil {
				m.detail.SetError(msg.Message)
			}
		}
	}

	switch m.screen {
	case scrSplash:
		_, pageCmd = m.splash.Update(msg)
	case scrLogin:
		_, pageCmd = m.login.Update(msg)
	case scrRegister:
		_, pageCmd = m.register.Update(msg)
	}

	return m, pageCmd
}

func (m *Model) pageKey(msg tea.KeyMsg) tea.Cmd {
	k := msg.String()

	switch m.screen {
	case scrBrowse:
		switch k {
		case "down", "j":
			m.browse.MoveDown()
		case "up", "k":
			m.browse.MoveUp()
		case "enter":
			mv := m.browse.SelectedMovie()
			if mv != nil {
				m.detail = pages.NewMovieDetailModel(mv)
				m.screen = scrDetail
			}
		case "d":
			mv := m.browse.SelectedMovie()
			if mv != nil {
				m.detail = pages.NewMovieDetailModel(mv)
				m.screen = scrDetail
			}
		case "w":
			mv := m.browse.SelectedMovie()
			if mv != nil {
				return m.doAddToWishlist(mv.ID, false)
			}
		case "r":
			return func() tea.Msg { return pages.NavigateMsg{Page: "rentals"} }
		case "p":
			return func() tea.Msg { return pages.NavigateMsg{Page: "profile"} }
		case "f5":
			return m.loadMovies()
		}
	case scrDetail:
		switch k {
		case "enter":
			if m.detail != nil && !m.detail.Rented && m.detail.Movie.Available {
				return func() tea.Msg { return pages.RentRequestMsg{MovieID: m.detail.Movie.ID} }
			} else if m.detail != nil && !m.detail.Rented && !m.detail.Movie.Available {
				m.detail.ErrMsg = "No copies available — press [W] to join the waitlist"
			}
		case "w":
			if m.detail != nil {
				return m.doAddToWishlist(m.detail.Movie.ID, true)
			}
		case "f5":
			return m.loadMovies()
		}
	case scrRentals:
		switch k {
		case "down", "j":
			m.rentals.MoveDown()
		case "up", "k":
			m.rentals.MoveUp()
		case "enter":
			r := m.rentals.SelectedRental()
			if r != nil && r.Status != "returned" {
				return func() tea.Msg { return pages.ReturnRequestMsg{RentalID: r.ID} }
			}
		}
	case scrProfile:
		if k == "l" {
			return func() tea.Msg { return pages.NavigateMsg{Page: "login"} }
		}
	}
	return nil
}

func (m *Model) View() string {
	if !m.ready {
		return "loading..."
	}
	fh := lipgloss.Height(m.footerView())
	hh := lipgloss.Height(m.headerView())
	ch := m.h - hh - fh
	if ch < 5 {
		ch = 5
	}

	var body string
	switch m.screen {
	case scrSplash:
		return m.splash.View(m.w, m.h)
	case scrLogin:
		body = m.login.View(m.w, ch)
	case scrRegister:
		body = m.register.View(m.w, ch)
	case scrBrowse:
		body = m.browse.View(m.w, ch)
	case scrDetail:
		body = m.detail.View(m.w, ch)
	case scrRentals:
		body = m.rentals.View(m.w, ch)
	case scrProfile:
		body = m.profile.View(m.w, ch)
	}

	return lipgloss.JoinVertical(lipgloss.Top, m.headerView(), body, m.footerView())
}

func (m *Model) headerView() string {
	un, tier := "", ""
	pts := 0
	loggedIn := m.userResp != nil
	if loggedIn {
		un = m.userResp.Username
		tier = m.userResp.TierName
		pts = m.userResp.PopcornPoints
	}
	return m.header.View(m.w, loggedIn, un, tier, pts)
}

func (m *Model) footerView() string {
	var hints string
	switch m.screen {
	case scrSplash:
		hints = "[ENTER] start  [Ctrl+C] quit"
	case scrLogin:
		hints = "[TAB] switch  [ENTER] login  [Ctrl+R] sign up  [Ctrl+C] quit"
	case scrRegister:
		hints = "[TAB] switch  [ENTER] create account  [Ctrl+L] back to login"
	case scrBrowse:
		hints = "[↑↓] navigate  [ENTER] details  [R] rentals  [P] profile  [W] waitlist  [F5] refresh  [Ctrl+C] quit"
	case scrDetail:
		if m.detail != nil && !m.detail.Rented {
			hints = "[ENTER] rent  [W] waitlist  [F5] refresh  [Q] back  [Ctrl+C] quit"
		} else {
			hints = "[W] waitlist  [F5] refresh  [Q] back  [Ctrl+C] quit"
		}
	case scrRentals:
		hints = "[↑↓] select  [ENTER] return  [Q] back"
	case scrProfile:
		hints = "[L] logout  [Q] back"
	default:
		hints = "[Q] back  [Ctrl+C] quit"
	}
	return lipgloss.NewStyle().Background(styles.BgBlue).Foreground(styles.TextMedium).Width(m.w).Padding(0, 1).Render(hints)
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
		}
		json.NewDecoder(resp.Body).Decode(&r)
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

func (m *Model) loadMovies() tea.Cmd {
	return func() tea.Msg {
		req, _ := http.NewRequest("GET", m.baseURL+"/api/v1/movies?page_size=50", nil)
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
		return loadMoviesMsg{movies: r.Movies, total: r.Total}
	}
}

func (m *Model) doRent(movieID string) tea.Cmd {
	return func() tea.Msg {
		body := `{"movie_id":"` + movieID + `"}`
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
		m.browse.Status = "Rented! " + rental.MovieTitle
		return m.loadMovies()()
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
			m.detail.ErrMsg = "Added to waitlist ✓"
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
		http.DefaultClient.Do(req)
		m.rentals.Status = "Returned!"
		return tea.Batch(m.loadRentals(), m.loadMovies())()
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
