package scaler

type optionalValue struct {
	// the value
	value uint
	// is the value already known
	isKnown bool
}

func (ov *optionalValue) setValue(val uint) {
	ov.value = val
	ov.isKnown = true
}
