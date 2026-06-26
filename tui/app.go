package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/components"
	"github.com/thelastvideostore/tui/pages"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	baseURL  string
	screen   screen
	w, h     int
	ready    bool
	token    string
	userResp *models.UserResponse

	splash         *pages.SplashModel
	login          *pages.LoginModel
	register       *pages.RegisterModel
	browse         *pages.BrowseModel
	detail         *pages.MovieDetailModel
	rentals        *pages.MyRentalsModel
	profile        *pages.ProfileModel
	wishlist       *pages.WishlistModel
	merch          *pages.MerchModel
	inventory      *pages.InventoryModel
	tierShop       *pages.TierShopModel
	header         *components.HeaderModel
	adminMovies    *pages.AdminMoviesModel
	adminUsers     *pages.AdminUsersModel
	auditLog       *pages.AuditLogModel
	movieForm      *pages.MovieFormModel
	accessDenied   *pages.AccessDeniedModel
	snackBarMenu   *pages.SnackBarMenuModel
	snackBarOrders *pages.SnackBarOrdersModel
	snackBarManage *pages.SnackBarManageModel
	gameDetail     *pages.GameDetailModel
	gameSessions   *pages.GameSessionModel
	myPlaySessions *pages.MyPlaySessionsModel

	help help.Model

	searchBar   *components.SearchbarModel
	searching   bool
	tabs        *components.TabsModel
	genreTabs   *components.TabsModel
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
		help:      newHelpModel(),
		tabs:      components.NewTabsModel([]string{"🎬 Movies", "📺 Series", "🕹️ Games", "🍿 SnackBar"}),
		genreTabs: components.NewTabsModel([]string{"ALL", "Action", "SciFi", "Horror", "Comedy", "Drama", "Thriller", "Romance", "Animation"}),
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

func (m *Model) applyMediaTypeFilter() tea.Cmd {
	switch m.tabs.ActiveTab() {
	case "📺 Series":
		m.browse.MediaType = "series"
		m.browse.Genre = ""
		m.genreTabs = components.NewTabsModel(append([]string{"ALL"}, models.SeriesGenreList...))
	case "🕹️ Games":
		m.browse.MediaType = "game"
		m.browse.Genre = ""
		m.genreTabs = components.NewTabsModel(append([]string{"ALL"}, models.GameGenreList...))
	case "🍿 SnackBar":
		bal := 0.0
		if m.userResp != nil {
			bal = m.userResp.Balance
		}
		m.snackBarMenu = pages.NewSnackBarMenuModel(bal)
		m.snackBarMenu.SetItems(nil)
		m.screen = scrSnackBarMenu
		m.genreTabs = components.NewTabsModel([]string{})
		return m.loadSnackBarMenu()
	default:
		m.browse.MediaType = "movie"
		m.browse.Genre = ""
		m.genreTabs = components.NewTabsModel(append([]string{"ALL"}, models.GenreList...))
	}
	m.browse.Selected = -1
	return m.loadMovies(1, m.browse.Genre)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var pageCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.help.Width = msg.Width
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" || k == "ctrl+d" {
			return m, tea.Quit
		}
		if helpKeyMatches(msg) {
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
		if m.searching {
			return m, m.searchKey(msg)
		}
		if k == "q" && m.screen != scrLogin && m.screen != scrSplash && m.screen != scrBrowse && m.screen != scrTOTP && m.screen != scrMovieForm && m.screen != scrAccessDenied && m.screen != scrRegister {
			if m.screen == scrMerch || m.screen == scrInventory || m.screen == scrTierShop || m.screen == scrSnackBarMenu {
				if m.profile == nil {
					m.profile = pages.NewProfileModel(m.userResp)
				}
				m.screen = scrProfile
				return m, m.loadProfile()
			}
			if m.screen == scrSnackBarOrders || m.screen == scrSnackBarManage {
				m.screen = scrSnackBarMenu
				return m, nil
			}
			if m.screen == scrGameDetail || m.screen == scrGameSessions {
				m.screen = scrBrowse
				return m, nil
			}
			if m.screen == scrMyPlaySessions {
				m.screen = scrRentals
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
			m.browse.MediaType = "movie"
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
			case "play_sessions":
				if m.myPlaySessions == nil {
					m.myPlaySessions = pages.NewMyPlaySessionsModel()
				}
				m.screen = scrMyPlaySessions
				return m, tea.Batch(m.loadMyPlaySessions(), m.myPlaySessions.Init())
			case "profile":
				m.profile = pages.NewProfileModel(m.userResp)
				m.screen = scrProfile
				return m, m.loadProfile()
			}

		case pages.RentRequestMsg:
			return m, tea.Sequence(m.doRent(msg.MovieID), m.doRefreshMe())
		case pages.ReturnRequestMsg:
			return m, tea.Sequence(m.doReturn(msg.RentalID), m.doRefreshMe())
		case pages.ExtendRentalMsg:
			return m, tea.Sequence(m.doExtendRental(msg.RentalID), m.doRefreshMe())

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
			if m.gameDetail != nil && m.gameDetail.Game != nil {
				for i := range msg.movies {
					if msg.movies[i].ID == m.gameDetail.Game.ID {
						m.gameDetail.Game = &msg.movies[i]
						break
					}
				}
			}
		case pages.BrowseReloadMsg:
			return m, m.loadMovies(m.browse.Page, m.browse.Genre)
		case loadRentalsMsg:
			m.rentals.SetRentals(msg.rentals)
		case loadMyPlaySessionsMsg:
			if m.myPlaySessions == nil {
				m.myPlaySessions = pages.NewMyPlaySessionsModel()
			}
			m.myPlaySessions.SetSessions(msg.sessions)
		case pages.PlayTickMsg:
			if m.screen == scrMyPlaySessions {
				// The page itself handles rendering the per-second
				// countdown; we just refetch whenever a session
				// crosses its expiry so the server's "ended"
				// state replaces our locally-rendered EXPIRED row.
				if m.myPlaySessions != nil && m.myPlaySessions.HasExpired() {
					return m, m.loadMyPlaySessions()
				}
			}
			return m, nil
		case loadProfileMsg:
			m.profile.SetStats(msg.stats)
		case loadWishlistMsg:
			m.wishlist.SetItems(msg.items)
		case wishlistResultMsg:
			return m, nil
		case refreshMeMsg:
			if msg.user != nil {
				m.userResp = msg.user
				if m.snackBarMenu != nil {
					m.snackBarMenu.Balance = msg.user.Balance
				}
				if m.detail != nil {
					m.detail.Balance = msg.user.Balance
					m.detail.FreeRentals = msg.user.FreeRentals
				}
			}
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
			if m.screen == scrDetail && m.detail != nil && m.detail.Movie != nil {
				return m, tea.Batch(m.doRefreshDetail(m.detail.Movie.ID), autoRefreshCmd())
			}
			if m.screen == scrGameDetail && m.gameDetail != nil && m.gameDetail.Game != nil {
				return m, tea.Batch(m.doRefreshDetail(m.gameDetail.Game.ID), autoRefreshCmd())
			}
			if m.screen == scrMyPlaySessions {
				return m, tea.Batch(m.loadMyPlaySessions(), autoRefreshCmd())
			}
			return m, autoRefreshCmd()
		case refreshDetailMsg:
			if m.screen == scrDetail && m.detail != nil && msg.movie != nil {
				m.detail.Movie = msg.movie
			}
			if m.screen == scrGameDetail && m.gameDetail != nil && msg.movie != nil {
				m.gameDetail.Game = msg.movie
			}
		case loadInventoryMsg:
			m.inventory.SetItems(msg.items)
		case loadSnackBarMenuMsg:
			m.snackBarMenu.SetItems(msg.items)
		case loadSnackBarOrdersMsg:
			m.snackBarOrders.SetOrders(msg.orders)
		case loadSnackBarManageMsg:
			m.snackBarManage.SetItems(msg.items)
		case loadGameSessionsMsg:
			m.gameSessions.SetSessions(msg.sessions)
		case gameRefreshMsg:
			return m, nil
		case pages.SnackBarOrderMsg:
			return m, tea.Sequence(m.doSnackBarOrder(msg.ItemID), m.doRefreshMe())
		case pages.MerchRedeemMsg:
			return m, tea.Sequence(m.doRedeemMerch(msg.ItemID), m.doRefreshMe())
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

	// Audit log uses a table component (not a list) so it has its own
	// dedicated keypath. Route all messages to it.
	if m.screen == scrAuditLog && m.auditLog != nil {
		_, pc := m.auditLog.Update(msg)
		pageCmd = tea.Batch(pageCmd, pc)
	}

	// Play-sessions screen also receives non-key messages (the per-second
	// tick that drives the live countdown). Route the original message
	// through the page's Update so the tick can self-perpetuate.
	if m.screen == scrMyPlaySessions && m.myPlaySessions != nil {
		_, pc := m.myPlaySessions.Update(msg)
		pageCmd = tea.Batch(pageCmd, pc)
	}

	// Route messages to list-based pages so their built-in filtering,
	// selection and navigation work. Skip if the user is in a filter
	// input (list is handling all keys) or if the list already accepted
	// a navigation key.
	if km, ok := msg.(tea.KeyMsg); ok {
		switch m.screen {
		case scrRentals:
			if m.rentals != nil {
				_, pc := m.rentals.Update(km)
				pageCmd = tea.Batch(pageCmd, pc)
			}
		case scrMerch:
			if m.merch != nil {
				_, pc := m.merch.Update(km)
				pageCmd = tea.Batch(pageCmd, pc)
			}
		case scrTierShop:
			if m.tierShop != nil {
				_, pc := m.tierShop.Update(km)
				pageCmd = tea.Batch(pageCmd, pc)
			}
		case scrAdminMovies:
			if m.adminMovies != nil {
				_, pc := m.adminMovies.Update(km)
				pageCmd = tea.Batch(pageCmd, pc)
			}
		case scrAdminUsers:
			if m.adminUsers != nil {
				_, pc := m.adminUsers.Update(km)
				pageCmd = tea.Batch(pageCmd, pc)
			}
		case scrSnackBarMenu:
			if m.snackBarMenu != nil {
				_, pc := m.snackBarMenu.Update(km)
				pageCmd = tea.Batch(pageCmd, pc)
			}
		case scrSnackBarOrders:
			if m.snackBarOrders != nil {
				_, pc := m.snackBarOrders.Update(km)
				pageCmd = tea.Batch(pageCmd, pc)
			}
		case scrSnackBarManage:
			if m.snackBarManage != nil {
				_, pc := m.snackBarManage.Update(km)
				pageCmd = tea.Batch(pageCmd, pc)
			}
		}
	}

	return m, pageCmd
}
