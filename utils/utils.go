package utils

import (
	"fmt"
)

func FormatFloat(value float64) string {
	if value == 0 {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", value)
}

func FormatInt(value int64) string {
	if value == 0 {
		return "0"
	}
	return fmt.Sprintf("%d", value)
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func IntPtr(i int) *int {
	return &i
}
