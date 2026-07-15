package ai

import "strings"

func extractJSON(s string) string {
	s = strings.TrimSpace(s)

	start := strings.IndexAny(s, "[{")
	if start == -1 {
		return s
	}

	endObj := strings.LastIndex(s, "}")
	endArr := strings.LastIndex(s, "]")

	end := max(endObj, endArr)
	if end == -1 {
		return s
	}

	return s[start : end+1]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
