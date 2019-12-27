package helper

import "fmt"

// SimpleTime just to express hours and minutes
type SimpleTime struct {
	hour   uint
	minute uint
}

// NewTime creates a new SimpleTime instance based on the given parameters
func NewTime(hour, minute uint) (SimpleTime, error) {
	if hour > 23 {
		return SimpleTime{}, fmt.Errorf("Given parameter hour is invalid (%d). The value must not be greater than 23", hour)
	}

	if minute > 59 {
		return SimpleTime{}, fmt.Errorf("Given parameter minute is invalid (%d). The value must not be greater than 59", minute)
	}

	return SimpleTime{hour, minute}, nil
}

// Minutes returns the time in minutes
func (s SimpleTime) Minutes() uint {
	return s.hour*60 + s.minute
}
