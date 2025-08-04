package utils

import (
	"fmt"
	"strings"
)

func FormatFloat(value float64) string {
	formatted := fmt.Sprintf("%.2f", value)

	parts := strings.Split(formatted, ".")
	parts[0] = addUnderscores(parts[0])

	return strings.Join(parts, ".")
}

func FormatInt(value int64) string {
	return addUnderscores(fmt.Sprintf("%d", value))
}

func addUnderscores(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}

	var result strings.Builder
	mod := n % 3
	if mod == 0 {
		mod = 3
	}

	result.WriteString(s[:mod])
	for i := mod; i < n; i += 3 {
		result.WriteString("_")
		result.WriteString(s[i : i+3])
	}

	return result.String()
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func IntPtr(i int) *int {
	return &i
}
