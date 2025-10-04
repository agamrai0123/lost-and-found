package pkg

import "strings"

func validStatus(s PostStatus) bool {
	switch s {
	case StatusLost, StatusFound, StatusClaimed:
		return true
	default:
		return false
	}
}

func containsIgnoreCase(haystack, needle string) bool {
	if needle == "" {
		return true
	}
	return caseInsensitiveContains(haystack, needle)
}

// simple case-insensitive substring check without importing strings multiple times
func caseInsensitiveContains(a, b string) bool {
	// convert to lower using ASCII-friendly approach
	la := []rune(a)
	lb := []rune(b)
	if len(lb) == 0 {
		return true
	}
	// Lowercase runes (simple) — for most use-cases this is fine.
	for i := range la {
		if la[i] >= 'A' && la[i] <= 'Z' {
			la[i] = la[i] - 'A' + 'a'
		}
	}
	for i := range lb {
		if lb[i] >= 'A' && lb[i] <= 'Z' {
			lb[i] = lb[i] - 'A' + 'a'
		}
	}

	// naive substring search
	stra := string(la)
	strb := string(lb)
	return (len(strb) == 0) || (len(stra) >= len(strb) && (indexOf(stra, strb) >= 0))
}

func indexOf(s, substr string) int {
	// fallback to built-in strings.Index (use std lib)
	return indexUsingStrings(s, substr)
}

func indexUsingStrings(s, substr string) int {
	return func() int {
		// using strings package here directly
		// this helper exists to keep top-level imports organized
		return indexStrings(s, substr)
	}()
}

func indexStrings(s, substr string) int {
	// tiny wrapper to avoid adding strings import at top-level repeated elsewhere
	// but we need strings here, so import it
	return stringsIndex(s, substr)
}

func stringsIndex(s, substr string) int { // this function uses strings.Index from stdlib
	return strings.Index(s, substr)
}
