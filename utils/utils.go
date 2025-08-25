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

func FormatUint64(value uint64) string {
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

func Int8Ptr(i int8) *int8 {
	return &i
}

func Uint8Ptr(i uint8) *uint8 {
	return &i
}

func GetTime(microseconds float64) string {
	if microseconds < 1_000 {
		return fmt.Sprintf("%.2f µs", microseconds)
	} else if microseconds < 1_000_000 {
		return fmt.Sprintf("%.2f ms", microseconds/1_000)
	} else if microseconds < 1_000_000_000 {
		return fmt.Sprintf("%.2f s", microseconds/1_000_000)
	} else {
		return fmt.Sprintf("%.2f ns", microseconds*1_000)
	}
}
