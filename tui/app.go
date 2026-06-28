package tui

import (
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
		if m.detail != nil {
			_, dc := m.detail.Update(msg)
			if dc != nil {
				return m, dc
			}
		}
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
		if m.screen == scrDetail && m.detail != nil {
			_, dc := m.detail.Update(msg)
			if dc != nil {
				return m, dc
			}
		}
		if m.searching {
			return m, m.searchKey(msg)
		}
		if k == "q" && m.screen != scrLogin && m.screen != scrSplash && m.screen != scrBrowse && m.screen != scrTOTP && m.screen != scrMovieForm && m.screen != scrAccessDenied && m.screen != scrRegister {
			cmd := m.handleQKey()
			if cmd != nil {
				return m, cmd
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
		if cmd := m.handleMessage(msg); cmd != nil {
			return m, cmd
		}
	}

	cmds := append([]tea.Cmd{pageCmd}, m.dispatchToCurrentScreen(msg)...)
	return m, tea.Batch(cmds...)
}
