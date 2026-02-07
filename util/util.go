package util

import (
	"strings"
	"unsafe"
)

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

func SameSlice[T any](a, b []T) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	return unsafe.Pointer(&a[0]) == unsafe.Pointer(&b[0])
}
