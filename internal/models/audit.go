package models

const (
	ActionLogin              = "login"
	ActionLogout             = "logout"
	ActionRegister           = "register"
	ActionRent               = "rent"
	ActionReturn             = "return"
	ActionPromote            = "promote"
	ActionDemote             = "demote"
	ActionBan                = "ban"
	ActionUnban              = "unban"
	ActionAddMovie           = "add_movie"
	ActionEditMovie          = "edit_movie"
	ActionDeleteMovie        = "delete_movie"
	ActionTOTPEnabled        = "totp_enabled"
	ActionTOTPDisabled       = "totp_disabled"
	ActionAddToWishlist      = "add_to_wishlist"
	ActionRemoveFromWishlist = "remove_from_wishlist"
	ActionAddStaffPick       = "add_staff_pick"
	ActionRemoveStaffPick    = "remove_staff_pick"
)

type AuditEntry struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Action    string `json:"action"`
	ActorID   string `json:"actor_id"`
	TargetID  string `json:"target_id"`
	Data      string `json:"data,omitempty"`
	Hash      []byte `json:"hash"`
	PrevHash  []byte `json:"prev_hash"`
}
