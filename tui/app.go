package tui

import (
	"strings"
	"time"

	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/components"
	"github.com/thelastvideostore/tui/pages"

	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	scrSplash screen = iota
	scrLogin
	scrRegister
	scrTOTP
	scrBrowse
	scrDetail
	scrRentals
	scrProfile
	scrWishlist
	scrMerch
	scrInventory
	scrTierShop
	scrAdminMovies
	scrAdminUsers
	scrAuditLog
	scrMovieForm
	scrAccessDenied
)

type loadMoviesMsg struct {
	movies []models.MovieResponse
	total  int
	page   int
	reqID  int
}
type loadRentalsMsg struct {
	rentals []models.RentalResponse
}
type loadProfileMsg struct {
	stats *pages.RentalStats
}
type loadWishlistMsg struct {
	items []pages.WishlistItem
}
type searchResultsMsg struct {
	results []models.MovieResponse
}
type autoRefreshMsg struct{}
type loadAdminMoviesMsg struct {
	movies []models.MovieResponse
	total  int
	page   int
}
type loadAdminUsersMsg struct {
	users []models.UserResponse
}
type loadAuditLogMsg struct {
	entries []map[string]interface{}
}
type loadMerchMsg struct {
	items []models.MerchItem
}
type loadInventoryMsg struct {
	items []pages.InventoryItem
}

type Model struct {
	baseURL  string
	screen   screen
	w, h     int
	ready    bool
	token    string
	userResp *models.UserResponse

	splash       *pages.SplashModel
	login        *pages.LoginModel
	register     *pages.RegisterModel
	browse       *pages.BrowseModel
	detail       *pages.MovieDetailModel
	rentals      *pages.MyRentalsModel
	profile      *pages.ProfileModel
	wishlist     *pages.WishlistModel
	merch        *pages.MerchModel
	inventory    *pages.InventoryModel
	tierShop     *pages.TierShopModel
	header       *components.HeaderModel
	adminMovies  *pages.AdminMoviesModel
	adminUsers   *pages.AdminUsersModel
	auditLog     *pages.AuditLogModel
	movieForm    *pages.MovieFormModel
	accessDenied *pages.AccessDeniedModel

	searchBar   *components.SearchbarModel
	searching   bool
	tabs        *components.TabsModel
	tempToken   string
	totpCode    string
	browseReqID int
}

func NewModel(baseURL string) *Model {
	return &Model{
		baseURL:   baseURL,
		screen:    scrSplash,
		splash:    pages.NewSplashModel(),
		login:     pages.NewLoginModel(),
		header:    components.NewHeaderModel(),
		searchBar: components.NewSearchbarModel(),
		tabs:      components.NewTabsModel([]string{"ALL", "Action", "SciFi", "Horror", "Comedy", "Drama", "Thriller", "Romance", "Animation"}),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen, m.splash.Init())
}

func autoRefreshCmd() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return autoRefreshMsg{}
	})
}

func (m *Model) setDetailContext() {
	if m.detail == nil || m.userResp == nil {
		return
	}
	tier := models.TierByName(m.userResp.Subscription)
	m.detail.SetUserContext(m.userResp.FreeRentals, tier.FreeRentals, m.userResp.Balance)
	if m.detail.Movie.SequelTo != "" {
		for _, mv := range m.browse.Movies {
			if mv.ID == m.detail.Movie.SequelTo {
				m.detail.SequelTitle = mv.Title
				break
			}
		}
	}
	var franchise []models.MovieResponse
	currentID := m.detail.Movie.ID
	seen := map[string]bool{currentID: true}
	// Find prequels
	id := m.detail.Movie.SequelTo
	for id != "" && !seen[id] {
		found := false
		for _, mv := range m.browse.Movies {
			if mv.ID == id {
				seen[id] = true
				franchise = append([]models.MovieResponse{mv}, franchise...)
				id = mv.SequelTo
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	// Add current movie
	franchise = append(franchise, *m.detail.Movie)
	// Find sequels
	queue := []string{currentID}
	for len(queue) > 0 {
		prequelID := queue[0]
		queue = queue[1:]
		for _, mv := range m.browse.Movies {
			if mv.SequelTo == prequelID && !seen[mv.ID] {
				seen[mv.ID] = true
				franchise = append(franchise, mv)
				queue = append(queue, mv.ID)
			}
		}
	}
	// Only show franchise if there are 2+ movies in chain
	if len(franchise) > 1 {
		m.detail.Franchise = franchise
	} else {
		m.detail.Franchise = nil
	}

	var sameGenre []models.MovieResponse
	for _, mv := range m.browse.Movies {
		if !seen[mv.ID] && mv.Genre == m.detail.Movie.Genre {
			sameGenre = append(sameGenre, mv)
			if len(sameGenre) >= 5 {
				break
			}
		}
	}
	m.detail.Recommendations = sameGenre
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
		if m.searching {
			return m, m.searchKey(msg)
		}
		if k == "q" && m.screen != scrLogin && m.screen != scrSplash && m.screen != scrBrowse && m.screen != scrTOTP && m.screen != scrMovieForm && m.screen != scrAccessDenied && m.screen != scrRegister {
			if m.screen == scrMerch || m.screen == scrInventory || m.screen == scrTierShop {
				m.screen = scrProfile
				return m, nil
			}
			if m.screen == scrDetail || m.screen == scrRentals || m.screen == scrProfile || m.screen == scrWishlist || m.screen == scrAdminMovies || m.screen == scrAdminUsers || m.screen == scrAuditLog {
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
			return m, tea.Batch(m.loadMovies(1, ""), autoRefreshCmd())

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
		case pages.ExtendRentalMsg:
			return m, m.doExtendRental(msg.RentalID)

		case loadMoviesMsg:
			if msg.reqID != 0 && msg.reqID != m.browseReqID {
				return m, nil
			}
			m.browse.SetMovies(msg.movies, msg.total, msg.page)
			if m.detail != nil && m.detail.Movie != nil {
				for i := range msg.movies {
					if msg.movies[i].ID == m.detail.Movie.ID {
						m.detail.Movie = &msg.movies[i]
						break
					}
				}
			}
		case pages.BrowseReloadMsg:
			return m, m.loadMovies(m.browse.Page, m.browse.Genre)
		case loadRentalsMsg:
			m.rentals.SetRentals(msg.rentals)
		case loadProfileMsg:
			m.profile.SetStats(msg.stats)
		case loadWishlistMsg:
			m.wishlist.SetItems(msg.items)
		case pages.WishlistRemoveMsg:
			return m, m.doRemoveFromWishlist(msg.MovieID)
		case searchResultsMsg:
			m.searchBar.SetResults(msg.results)
		case loadAdminMoviesMsg:
			m.adminMovies.SetMovies(msg.movies, msg.total, msg.page)
		case loadAdminUsersMsg:
			m.adminUsers.SetUsers(msg.users)
		case loadAuditLogMsg:
			m.auditLog.SetEntries(msg.entries)
		case loadMerchMsg:
			m.merch.SetItems(msg.items)
		case autoRefreshMsg:
			if m.screen == scrBrowse && !m.searching {
				return m, tea.Batch(m.loadMovies(m.browse.Page, m.browse.Genre), autoRefreshCmd())
			}
			return m, autoRefreshCmd()
		case loadInventoryMsg:
			m.inventory.SetItems(msg.items)
		case pages.MerchRedeemMsg:
			return m, m.doRedeemMerch(msg.ItemID)
		case pages.MovieFormSubmitMsg:
			if msg.Mode == pages.FormAdd {
				return m, m.doCreateMovie(msg)
			}
			return m, m.doUpdateMovie(msg)

		case pages.ErrorMsg:
			if strings.Contains(msg.Message, "ACCESS DENIED") || strings.Contains(msg.Message, "⛔") {
				m.accessDenied = pages.NewAccessDeniedModel(msg.Message)
				m.screen = scrAccessDenied
			}
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
	case scrMovieForm:
		pageCmd, _ = m.movieForm.Update(msg)
	}

	return m, pageCmd
}
