package pages

import (
	"fmt"
	"strconv"
	"strings"
)

func nonEmptyString(label string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errorMsg(label)
		}
		return nil
	}
}

func yearValidator(thisYear int) func(string) error {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return errorMsg("year is required")
		}
		if len(s) != 4 {
			return errorMsg("year must be 4 digits (e.g. 1999)")
		}
		for _, r := range s {
			if r < '0' || r > '9' {
				return errorMsg("year must be digits only")
			}
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return errorMsg("not a valid year")
		}
		if n < 1880 || n > thisYear+5 {
			return errorMsg(fmt.Sprintf("year must be between 1880 and %d", thisYear+5))
		}
		return nil
	}
}

func positiveIntValidator(label string, min, max int) func(string) error {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return errorMsg(label + " is required")
		}
		if !digitsOnly(s) {
			return errorMsg(label + " must be digits only")
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return errorMsg("not a valid number")
		}
		if n < min || n > max {
			return errorMsg(fmt.Sprintf("%s must be between %d and %d", label, min, max))
		}
		return nil
	}
}

func optionalIntValidator(label string, min, max int) func(string) error {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		if !digitsOnly(s) {
			return errorMsg(label + " must be digits only")
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return errorMsg("not a valid number")
		}
		if n < min || n > max {
			return errorMsg(fmt.Sprintf("%s must be between %d and %d", label, min, max))
		}
		return nil
	}
}

func priceValidator() func(string) error {
	return func(s string) error {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		for _, r := range s {
			if (r < '0' || r > '9') && r != '.' {
				return errorMsg("price must be a number (e.g. 4.99)")
			}
		}
		if strings.Count(s, ".") > 1 {
			return errorMsg("price can have at most one decimal point")
		}
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return errorMsg("not a valid amount")
		}
		if n < 0 || n > 999.99 {
			return errorMsg("price must be between 0.00 and 999.99")
		}
		return nil
	}
}

func digitsOnly(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func itoa(n int) string {
	if n == 0 {
		return ""
	}
	return strconv.Itoa(n)
}

func formatFloat(f float64) string {
	if f == 0 {
		return ""
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func joinCast(cast []string) string {
	out := ""
	for i, c := range cast {
		if i > 0 {
			out += ", "
		}
		out += c
	}
	return out
}
