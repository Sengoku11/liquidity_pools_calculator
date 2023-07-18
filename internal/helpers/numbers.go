package helpers

const (
	nearMaxInt  = 1e18
	maxValidDec = 19
	decBase     = 10
)

// Returns a number of digits of the given Int.
// It's faster than log10 approach for small numbers.
func Digits(i int) int {
	if i >= nearMaxInt {
		// on 64 bit architectures int's max value is bigger than 1e18
		// but smaller than 1e19
		return maxValidDec
	}

	x := decBase
	count := 1

	for x <= i {
		x *= decBase
		count++
	}

	return count
}
