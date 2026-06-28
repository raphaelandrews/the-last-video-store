package pages

func formatAction(a string) string {
	switch a {
	case "login":
		return "🔑 LOGIN"
	case "logout":
		return "🚪 LOGOUT"
	case "rent":
		return "📼 RENT"
	case "return":
		return "📀 RETURN"
	case "register":
		return "👤 REGISTER"
	case "promote":
		return "⬆ PROMOTE"
	case "demote":
		return "⬇ DEMOTE"
	case "ban":
		return "🚫 BAN"
	case "unban":
		return "✅ UNBAN"
	case "totp_enable":
		return "🔒 TOTP+"
	case "totp_disable":
		return "🔓 TOTP-"
	case "create_movie":
		return "➕ MOVIE+"
	case "update_movie":
		return "✎ MOVIE~"
	case "delete_movie":
		return "🗑 MOVIE-"
	case "staff_pick":
		return "★ STAFF"
	case "topup":
		return "💰 TOPUP"
	case "extend_rental":
		return "⏰ EXTEND"
	case "play_start":
		return "▶ PLAY"
	case "play_end":
		return "⏹ PLAY-END"
	case "purchase_tier":
		return "🏷️ TIER"
	case "redeem_merch":
		return "🎁 REDEEM"
	case "order_snackbar":
		return "🍿 SNACK"
	case "restock":
		return "📦 RESTOCK"
	default:
		return a
	}
}

func shortID(s string) string {
	if len(s) > 12 {
		return s[:12]
	}
	return s
}

func shortHash(s string) string {
	if len(s) > 12 {
		return s[:6] + "…" + s[len(s)-4:]
	}
	return s
}
