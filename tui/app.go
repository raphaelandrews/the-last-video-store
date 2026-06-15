package tui

import (
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
	scrAdminMovies
	scrAdminUsers
	scrAuditLog
)

type loadMoviesMsg struct {
	movies []models.MovieResponse
	total  int
	page   int
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
type loadAdminMoviesMsg struct {
	movies []models.MovieResponse
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

	splash      *pages.SplashModel
	login       *pages.LoginModel
	register    *pages.RegisterModel
	browse      *pages.BrowseModel
	detail      *pages.MovieDetailModel
	rentals     *pages.MyRentalsModel
	profile     *pages.ProfileModel
	wishlist    *pages.WishlistModel
	merch       *pages.MerchModel
	inventory   *pages.InventoryModel
	header      *components.HeaderModel
	adminMovies *pages.AdminMoviesModel
	adminUsers  *pages.AdminUsersModel
	auditLog    *pages.AuditLogModel

	searchBar *components.SearchbarModel
	searching bool
	tempToken string
	totpCode  string
}

func NewModel(baseURL string) *Model {
	return &Model{
		baseURL:   baseURL,
		screen:    scrSplash,
		splash:    pages.NewSplashModel(),
		login:     pages.NewLoginModel(),
		header:    components.NewHeaderModel(),
		searchBar: components.NewSearchbarModel(),
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
		if m.searching {
			return m, m.searchKey(msg)
		}
		if k == "q" && m.screen != scrLogin && m.screen != scrSplash && m.screen != scrBrowse && m.screen != scrTOTP {
			if m.screen == scrDetail || m.screen == scrRentals || m.screen == scrProfile || m.screen == scrRegister || m.screen == scrWishlist || m.screen == scrMerch || m.screen == scrInventory || m.screen == scrAdminMovies || m.screen == scrAdminUsers || m.screen == scrAuditLog {
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
			return m, m.loadMovies(1)

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
			return m, m.loadMovies(m.browse.Page)
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
			m.adminMovies.SetMovies(msg.movies)
		case loadAdminUsersMsg:
			m.adminUsers.SetUsers(msg.users)
		case loadAuditLogMsg:
			m.auditLog.SetEntries(msg.entries)
		case loadMerchMsg:
			m.merch.SetItems(msg.items)
		case loadInventoryMsg:
			m.inventory.SetItems(msg.items)
		case pages.MerchRedeemMsg:
			return m, m.doRedeemMerch(msg.ItemID)

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
