package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/thelastvideostore/internal/ds/bitmask"
)

type splashKeys struct{}

func (splashKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start")),
		quitBinding,
	}
}
func (splashKeys) FullHelp() [][]key.Binding { return nil }

type authKeys struct {
	isRegister bool
}

func (a authKeys) ShortHelp() []key.Binding {
	submit := key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit"))
	tab := key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch field"))
	switchLabel := "ctrl+r sign up"
	if a.isRegister {
		switchLabel = "ctrl+l to login"
	}
	switcher := key.NewBinding(key.WithKeys("ctrl+r", "ctrl+l"), key.WithHelp(switchLabel, "switch screen"))
	return []key.Binding{submit, tab, switcher, quitBinding}
}
func (a authKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit form"))},
		{key.NewBinding(key.WithKeys("tab", "shift+tab"), key.WithHelp("tab", "next/prev field"))},
		{key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "go to register")),
			key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "go to login"))},
		{quitBinding},
	}
}

type browseKeys struct{}

func (browseKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
		key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "profile")),
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rentals")),
	}
}
func (browseKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "navigate"))},
		{key.NewBinding(key.WithKeys("enter", "d"), key.WithHelp("enter", "open details"))},
		{key.NewBinding(key.WithKeys("[", "]"), key.WithHelp("[ ]", "switch media")),
			key.NewBinding(key.WithKeys(",", "."), key.WithHelp(", .", "switch genre"))},
		{key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n b", "next/prev page")),
			key.NewBinding(key.WithKeys("s", "l", "a"), key.WithHelp("s l a", "staff/last/all"))},
		{key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
			key.NewBinding(key.WithKeys("f5"), key.WithHelp("f5", "refresh"))},
		{key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "profile")),
			key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rentals")),
			key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "snack bar")),
			key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "wishlist"))},
	}
}

type movieDetailKeys struct {
	rented bool
}

func (d movieDetailKeys) ShortHelp() []key.Binding {
	if d.rented {
		return []key.Binding{backBinding, quitBinding}
	}
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "rent")),
		key.NewBinding(key.WithKeys("↑↓"), key.WithHelp("↑↓", "scroll")),
		key.NewBinding(key.WithKeys("n/p"), key.WithHelp("n/p", "next/prev related")),
		key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "waitlist")),
		backBinding,
	}
}
func (d movieDetailKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "scroll body")),
			key.NewBinding(key.WithKeys("pgup", "pgdown"), key.WithHelp("pgup/dn", "page")),
			key.NewBinding(key.WithKeys("home", "end"), key.WithHelp("home/end", "top/bottom"))},
		{key.NewBinding(key.WithKeys("n", "p"), key.WithHelp("n p", "next/prev related"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "rent")),
			key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "use ticket")),
			key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "pay money"))},
		{key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "waitlist")),
			key.NewBinding(key.WithKeys("f5"), key.WithHelp("f5", "refresh"))},
		{backBinding},
	}
}

type gameDetailKeys struct{}

func (gameDetailKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rent")),
		key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "play")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "end play")),
		backBinding,
	}
}
func (gameDetailKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rent game"))},
		{key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "play in-store"))},
		{key.NewBinding(key.WithKeys("1-5"), key.WithHelp("1-5", "play duration"))},
		{key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "end play session"))},
		{backBinding},
	}
}

type rentalsKeys struct{}

func (rentalsKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "return")),
		key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "extend (30🍿)")),
		key.NewBinding(key.WithKeys("←/→"), key.WithHelp("←/→", "page")),
		key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "play sessions")),
		backBinding,
	}
}
func (rentalsKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→", "prev/next page"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "return rental"))},
		{key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "extend due date"))},
		{key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "view play sessions"))},
		{backBinding},
	}
}

type profileKeys struct{}

func (profileKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logout")),
		key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "tiers")),
		key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "snack bar")),
		key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "rewards")),
	}
}
func (profileKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logout"))},
		{key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "tier shop")),
			key.NewBinding(key.WithKeys("$"), key.WithHelp("$", "top up")),
			key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "toggle TOTP"))},
		{key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "snack bar")),
			key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "rewards shop")),
			key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "inventory"))},
		{backBinding},
	}
}

type wishlistKeys struct{}

func (wishlistKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "remove")),
		backBinding,
	}
}
func (wishlistKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "info"))},
		{key.NewBinding(key.WithKeys("d", "delete"), key.WithHelp("d", "remove"))},
		{backBinding},
	}
}

type inventoryKeys struct{}

func (inventoryKeys) ShortHelp() []key.Binding {
	return []key.Binding{backBinding}
}
func (inventoryKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{{backBinding}}
}

type merchKeys struct{}

func (merchKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "redeem")),
		pageBinding,
		backBinding,
	}
}
func (merchKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→", "page"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "redeem"))},
		{backBinding},
	}
}

type tierShopKeys struct{}

func (tierShopKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "purchase")),
		backBinding,
	}
}
func (tierShopKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select tier"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "purchase"))},
		{backBinding},
	}
}

type snackBarMenuKeys struct {
	canManage bool
}

func (s snackBarMenuKeys) ShortHelp() []key.Binding {
	bindings := []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "order")),
		key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "orders")),
	}
	if s.canManage {
		bindings = append(bindings, key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "manage")))
	}
	bindings = append(bindings, pageBinding, backBinding)
	return bindings
}
func (s snackBarMenuKeys) FullHelp() [][]key.Binding {
	cols := [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select item"))},
		{key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→", "page"))},
		{key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "place order"))},
		{key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "order history"))},
	}
	if s.canManage {
		cols = append(cols, []key.Binding{key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "manage inventory"))})
	}
	cols = append(cols, []key.Binding{backBinding})
	return cols
}

type snackBarOrdersKeys struct{}

func (snackBarOrdersKeys) ShortHelp() []key.Binding  { return []key.Binding{backBinding} }
func (snackBarOrdersKeys) FullHelp() [][]key.Binding { return [][]key.Binding{{backBinding}} }

type snackBarManageKeys struct{}

func (snackBarManageKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restock")),
		backBinding,
	}
}
func (snackBarManageKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restock +5"))},
		{backBinding},
	}
}

type myPlaySessionsKeys struct{}

func (myPlaySessionsKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
		backBinding,
	}
}
func (myPlaySessionsKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh sessions"))},
		{backBinding},
	}
}

type totpKeys struct{}

func (totpKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
		quitBinding,
	}
}
func (totpKeys) FullHelp() [][]key.Binding { return nil }

type accessDeniedKeys struct{}

func (accessDeniedKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "back")),
		quitBinding,
	}
}
func (accessDeniedKeys) FullHelp() [][]key.Binding { return nil }

type adminMoviesKeys struct {
	page int
}

func (a adminMoviesKeys) ShortHelp() []key.Binding {
	bindings := []key.Binding{
		key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "edit")),
		key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "staff pick")),
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch type")),
		key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n/b", "next/prev page")),
		backBinding,
	}
	return bindings
}
func (a adminMoviesKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab/shift+tab", "switch type"))},
		{key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n b", "next/prev server page"))},
		{key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "edit"))},
		{key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
			key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "toggle staff pick"))},
		{key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter by title"))},
		{backBinding},
	}
}

type adminUsersKeys struct{}

func (adminUsersKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "promote")),
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "demote")),
		key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "ban")),
		pageBinding,
		backBinding,
	}
}
func (adminUsersKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "select"))},
		{key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→", "page"))},
		{key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "promote")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "demote"))},
		{key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "toggle ban")),
			key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "toggle TOTP"))},
		{backBinding},
	}
}

type auditLogKeys struct{}

func (auditLogKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "verify")),
		key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n/b", "page")),
		key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "jump to break")),
		backBinding,
	}
}
func (auditLogKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("up", "down", "k", "j"), key.WithHelp("↑↓", "navigate")),
			key.NewBinding(key.WithKeys("pgup", "pgdown"), key.WithHelp("pgup/dn", "jump page"))},
		{key.NewBinding(key.WithKeys("n", "b"), key.WithHelp("n b", "next/prev page")),
			key.NewBinding(key.WithKeys("home", "end"), key.WithHelp("home/end", "first/last page"))},
		{key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "verify hash chain")),
			key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "jump to broken entry"))},
		{backBinding},
	}
}

type movieFormKeys struct{}

func (movieFormKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab/↓", "next field")),
		key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab/↑", "prev")),
		key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "submit")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}
func (movieFormKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{key.NewBinding(key.WithKeys("tab", "down"), key.WithHelp("tab/↓", "next field")),
			key.NewBinding(key.WithKeys("shift+tab", "up"), key.WithHelp("shift+tab/↑", "prev field"))},
		{key.NewBinding(key.WithKeys("ctrl+s", "ctrl+enter"), key.WithHelp("ctrl+s", "submit")),
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))},
	}
}

func (m *Model) currentScreenKeys() screenKeyMap {
	switch m.screen {
	case scrSplash:
		return helpWith(splashKeys{})
	case scrLogin:
		return helpWith(authKeys{})
	case scrRegister:
		return helpWith(authKeys{isRegister: true})
	case scrTOTP:
		return helpWith(totpKeys{})
	case scrBrowse:
		return helpWith(browseKeys{})
	case scrDetail:
		rented := m.detail != nil && m.detail.Rented
		return helpWith(movieDetailKeys{rented: rented})
	case scrGameDetail:
		return helpWith(gameDetailKeys{})
	case scrRentals:
		return helpWith(rentalsKeys{})
	case scrProfile:
		return helpWith(profileKeys{})
	case scrWishlist:
		return helpWith(wishlistKeys{})
	case scrMerch:
		return helpWith(merchKeys{})
	case scrInventory:
		return helpWith(inventoryKeys{})
	case scrTierShop:
		return helpWith(tierShopKeys{})
	case scrSnackBarMenu:
		canManage := m.userResp != nil && bitmask.CanSnackBarManage(m.userResp.Tier)
		return helpWith(snackBarMenuKeys{canManage: canManage})
	case scrSnackBarOrders:
		return helpWith(snackBarOrdersKeys{})
	case scrSnackBarManage:
		return helpWith(snackBarManageKeys{})
	case scrMyPlaySessions:
		return helpWith(myPlaySessionsKeys{})
	case scrAdminMovies:
		return helpWith(adminMoviesKeys{page: m.adminMovies.Page})
	case scrAdminUsers:
		return helpWith(adminUsersKeys{})
	case scrAuditLog:
		return helpWith(auditLogKeys{})
	}
	return helpWith(accessDeniedKeys{})
}
