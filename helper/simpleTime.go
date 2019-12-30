package helper

import "fmt"

// SimpleTime just to express hours and minutes
type SimpleTime struct {
	Hour   uint `json:"hour,omitempty"`
	Minute uint `json:"minute,omitempty"`
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

// NewTimeFromMinute creates a new SimpleTime instance based on the given parameters
func NewTimeFromMinute(minutes uint) (SimpleTime, error) {

	hour := minutes / 60
	minute := minutes - hour*60

	return NewTime(hour, minute)
}

// Minutes returns the time in minutes
func (s SimpleTime) Minutes() uint {
	return s.Hour*60 + s.Minute
}

func (s SimpleTime) String() string {
	return fmt.Sprintf("%02d:%02d", s.Hour, s.Minute)
}
