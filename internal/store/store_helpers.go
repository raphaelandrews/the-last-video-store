package store

func hasBytePrefix(b, prefix []byte) bool {
	if len(b) < len(prefix) {
		return false
	}
	for i := range prefix {
		if b[i] != prefix[i] {
			return false
		}
	}
	return true
}

func splitCompositeKey(key string) []string {
	var parts []string
	start := 0
	count := 0
	for i := 0; i < len(key) && count < 3; i++ {
		if key[i] == ':' {
			parts = append(parts, key[start:i])
			start = i + 1
			count++
		}
	}
	if start < len(key) {
		parts = append(parts, key[start:])
	}
	return parts
}
