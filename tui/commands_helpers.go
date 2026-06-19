package tui

import (
	"strings"

	"github.com/thelastvideostore/internal/ds/bitmask"
)

var promotionChains = [][]bitmask.Permission{
	bitmask.TierPromotionOrder,
	bitmask.StaffPromotionOrder,
	bitmask.SnackBarPromotionOrder,
	bitmask.GamePromotionOrder,
}

func nextTier(current string) string {
	for _, order := range promotionChains {
		for i, t := range order {
			if bitmask.TierLabels[t] == current && i+1 < len(order) {
				return strings.ToLower(bitmask.TierLabels[order[i+1]])
			}
		}
	}
	return strings.ToLower(current)
}

func prevTier(current string) string {
	for _, order := range promotionChains {
		for i, t := range order {
			if bitmask.TierLabels[t] == current && i > 0 {
				return strings.ToLower(bitmask.TierLabels[order[i-1]])
			}
		}
	}
	return strings.ToLower(current)
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
