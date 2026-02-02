package util

// AbsInt is overflow-unsafe, but that's ok in most cases
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
