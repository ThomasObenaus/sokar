package helper

const (
	maxUint = ^uint(0)
	minUint = 0
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
)

// IncUint increments the given value (unsigned) by the given amount (signed)
// and returns the result as unsigned int (safe from under-/overflow)
func IncUint(value uint, by int) uint {

	// upper thresh of uint reached
	// avoid overflow
	if value >= maxUint && by >= 0 {
		return maxUint
	}

	// positive by --> just add
	if by >= 0 {
		return value + uint(by)
	}

	// negative by --> convert to uint and substract
	byUint := uint(by * -1)

	// avoid underflow
	if byUint >= value {
		return 0
	}

	return value - byUint
}

// SubUint substracts two uint's and returns the signed difference
// Avoids over-/underflow.
func SubUint(a uint, by uint) int {
	r := float64(a) - float64(by)

	if r >= float64(maxInt) {
		return maxInt
	}

	if r <= float64(minInt) {
		return minInt
	}

	return int(r)
}
