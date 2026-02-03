package util

import "strings"

// AbsInt is overflow-unsafe, but that's ok in most cases
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func PadRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
