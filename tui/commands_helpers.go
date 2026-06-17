package tui

import "strings"

var tierOrder = []string{"Bronze", "Silver", "Gold", "Employee", "Supervisor", "Manager", "Owner"}

func nextTier(current string) string {
	for i, t := range tierOrder {
		if t == current && i+1 < len(tierOrder) {
			return strings.ToLower(tierOrder[i+1])
		}
	}
	return strings.ToLower(current)
}

func prevTier(current string) string {
	for i, t := range tierOrder {
		if t == current && i > 0 {
			return strings.ToLower(tierOrder[i-1])
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
