package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/tui/pages"
)

func (m *Model) handleQKey() tea.Cmd {
	switch m.screen {
	case scrMerch, scrInventory, scrTierShop, scrSnackBarMenu:
		if m.profile == nil {
			m.profile = pages.NewProfileModel(m.userResp)
		}
		m.screen = scrProfile
		return m.loadProfile()
	case scrSnackBarOrders, scrSnackBarManage:
		m.screen = scrSnackBarMenu
		return nil
	case scrGameDetail:
		m.screen = scrBrowse
		return nil
	case scrMyPlaySessions:
		m.screen = scrRentals
		return nil
	case scrDetail, scrRentals, scrProfile, scrWishlist,
		scrAdminMovies, scrAdminUsers, scrAuditLog:
		m.screen = scrBrowse
		return nil
	}
	return nil
}

func (m *Model) handleMessage(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case pages.SplashDoneMsg:
		m.screen = scrLogin
		return m.login.Init()

	case pages.LoginRequestMsg:
		return m.doLogin(msg.Username, msg.Password)
	case pages.RegisterRequestMsg:
		return m.doRegister(msg.Username, msg.Password)

	case pages.LoginSuccessMsg:
		m.token = msg.AccessToken
		m.userResp = msg.User
		m.browse = pages.NewBrowseModel()
		m.browse.MediaType = "movie"
		m.screen = scrBrowse
		return tea.Batch(m.loadMovies(1, ""), autoRefreshCmd())

	case pages.NavigateMsg:
		return m.handleNavigate(msg)

	case pages.RentRequestMsg:
		return tea.Sequence(m.doRent(msg.MovieID), m.doRefreshMe())
	case pages.ReturnRequestMsg:
		return tea.Sequence(m.doReturn(msg.RentalID), m.doRefreshMe())
	case pages.ExtendRentalMsg:
		return tea.Sequence(m.doExtendRental(msg.RentalID), m.doRefreshMe())

	case loadMoviesMsg:
		return m.handleLoadMovies(msg)
	case pages.BrowseReloadMsg:
		return m.loadMovies(m.browse.Page, m.browse.Genre)

	case loadRentalsMsg:
		m.rentals.SetRentals(msg.rentals)
	case loadMyPlaySessionsMsg:
		if m.myPlaySessions == nil {
			m.myPlaySessions = pages.NewMyPlaySessionsModel()
		}
		m.myPlaySessions.SetSessions(msg.sessions)
	case pages.PlayTickMsg:
		return m.handlePlayTick()

	case loadProfileMsg:
		m.profile.SetStats(msg.stats)
	case loadWishlistMsg:
		m.wishlist.SetItems(msg.items)
	case wishlistResultMsg:
		return nil

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
		return m.doRemoveFromWishlist(msg.MovieID)
	case searchResultsMsg:
		m.searchBar.SetResults(msg.results)

	case loadAdminMoviesMsg:
		m.adminMovies.SetMovies(msg.movies, msg.total, msg.page)
	case loadCatalogOptionsMsg:
		if m.movieForm != nil {
			m.movieForm.SetOptions(msg.genres, msg.formats)
		}
	case loadAdminUsersMsg:
		m.adminUsers.SetUsers(msg.users)
	case loadAuditLogMsg:
		m.auditLog.SetEntries(msg.entries)
	case loadMerchMsg:
		m.merch.SetItems(msg.items)

	case autoRefreshMsg:
		return m.handleAutoRefresh()

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

	case gameRefreshMsg:
		return nil

	case pages.SnackBarOrderMsg:
		return tea.Sequence(m.doSnackBarOrder(msg.ItemID), m.doRefreshMe())
	case pages.MerchRedeemMsg:
		return tea.Sequence(m.doRedeemMerch(msg.ItemID), m.doRefreshMe())
	case pages.MovieFormSubmitMsg:
		if msg.Mode == pages.FormAdd {
			return m.doCreateMovie(msg)
		}
		return m.doUpdateMovie(msg)

	case pages.ErrorMsg:
		m.handleError(msg)
	}
	return nil
}

func (m *Model) handleNavigate(msg pages.NavigateMsg) tea.Cmd {
	switch msg.Page {
	case "login":
		m.token = ""
		m.userResp = nil
		m.login = pages.NewLoginModel()
		m.screen = scrLogin
		return m.login.Init()
	case "register":
		m.register = pages.NewRegisterModel()
		m.screen = scrRegister
		return m.register.Init()
	case "browse":
		m.screen = scrBrowse
	case "rentals":
		m.rentals = pages.NewMyRentalsModel()
		m.screen = scrRentals
		return m.loadRentals()
	case "play_sessions":
		if m.myPlaySessions == nil {
			m.myPlaySessions = pages.NewMyPlaySessionsModel()
		}
		m.screen = scrMyPlaySessions
		return tea.Batch(m.loadMyPlaySessions(), m.myPlaySessions.Init())
	case "profile":
		m.profile = pages.NewProfileModel(m.userResp)
		m.screen = scrProfile
		return m.loadProfile()
	}
	return nil
}

func (m *Model) handleLoadMovies(msg loadMoviesMsg) tea.Cmd {
	if msg.reqID != 0 && msg.reqID != m.browseReqID {
		return nil
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
	return nil
}

func (m *Model) handlePlayTick() tea.Cmd {
	if m.screen != scrMyPlaySessions {
		return nil
	}
	if m.myPlaySessions != nil && m.myPlaySessions.HasExpired() {
		return m.loadMyPlaySessions()
	}
	return nil
}

func (m *Model) handleAutoRefresh() tea.Cmd {
	if m.screen == scrBrowse && !m.searching {
		return tea.Batch(m.loadMovies(m.browse.Page, m.browse.Genre), autoRefreshCmd())
	}
	if m.screen == scrDetail && m.detail != nil && m.detail.Movie != nil {
		return tea.Batch(m.doRefreshDetail(m.detail.Movie.ID), autoRefreshCmd())
	}
	if m.screen == scrGameDetail && m.gameDetail != nil && m.gameDetail.Game != nil {
		return tea.Batch(m.doRefreshDetail(m.gameDetail.Game.ID), autoRefreshCmd())
	}
	if m.screen == scrMyPlaySessions {
		return tea.Batch(m.loadMyPlaySessions(), autoRefreshCmd())
	}
	return autoRefreshCmd()
}

func (m *Model) handleError(msg pages.ErrorMsg) {
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

func (m *Model) dispatchToCurrentScreen(msg tea.Msg) []tea.Cmd {
	var cmds []tea.Cmd
	switch m.screen {
	case scrSplash:
		_, pageCmd := m.splash.Update(msg)
		cmds = append(cmds, pageCmd)
	case scrLogin:
		_, pageCmd := m.login.Update(msg)
		cmds = append(cmds, pageCmd)
	case scrRegister:
		_, pageCmd := m.register.Update(msg)
		cmds = append(cmds, pageCmd)
	case scrMovieForm:
		_, pageCmd := m.movieForm.Update(msg)
		cmds = append(cmds, pageCmd)
	}
	if m.screen == scrAuditLog && m.auditLog != nil {
		_, pageCmd := m.auditLog.Update(msg)
		cmds = append(cmds, pageCmd)
	}
	if m.screen == scrMyPlaySessions && m.myPlaySessions != nil {
		_, pageCmd := m.myPlaySessions.Update(msg)
		cmds = append(cmds, pageCmd)
	}
	if km, ok := msg.(tea.KeyMsg); ok {
		switch m.screen {
		case scrRentals:
			if m.rentals != nil {
				_, pageCmd := m.rentals.Update(km)
				cmds = append(cmds, pageCmd)
			}
		case scrMerch:
			if m.merch != nil {
				_, pageCmd := m.merch.Update(km)
				cmds = append(cmds, pageCmd)
			}
		case scrTierShop:
			if m.tierShop != nil {
				_, pageCmd := m.tierShop.Update(km)
				cmds = append(cmds, pageCmd)
			}
		case scrAdminMovies:
			if m.adminMovies != nil {
				_, pageCmd := m.adminMovies.Update(km)
				cmds = append(cmds, pageCmd)
			}
		case scrAdminUsers:
			if m.adminUsers != nil {
				_, pageCmd := m.adminUsers.Update(km)
				cmds = append(cmds, pageCmd)
			}
		case scrSnackBarMenu:
			if m.snackBarMenu != nil {
				_, pageCmd := m.snackBarMenu.Update(km)
				cmds = append(cmds, pageCmd)
			}
		case scrSnackBarOrders:
			if m.snackBarOrders != nil {
				_, pageCmd := m.snackBarOrders.Update(km)
				cmds = append(cmds, pageCmd)
			}
		case scrSnackBarManage:
			if m.snackBarManage != nil {
				_, pageCmd := m.snackBarManage.Update(km)
				cmds = append(cmds, pageCmd)
			}
		}
	}
	return cmds
}
