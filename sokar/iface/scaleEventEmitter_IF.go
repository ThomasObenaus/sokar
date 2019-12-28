package sokar

import "fmt"

// ScaleEventEmitter represents the interface for the component that fires ScaleEvents if needed.
type ScaleEventEmitter interface {
	Subscribe(eventChannel chan ScaleEvent)
}

// scaleEventType represents the type of the ScaleEvent
type scaleEventType string

const (
	// scaleEventRegular denotes a regular ScaleEvent. This means a ScaleEvent which was triggered because of incoming scale alerts.
	scaleEventRegular scaleEventType = "regular"
	// scaleEventScheduled denotes a scheduled ScaleEvent. This means a ScaleEvent which was triggered according to an active schedule.
	scaleEventScheduled scaleEventType = "scheduled"
)

// ScaleEvent is an event that is created each time potentially a scale should be made (up/down)
type ScaleEvent struct {
	scaleFactor float32
	sType       scaleEventType
}

// NewScheduledScaleEvent creates a new ScaleEvent of type scheduled. This means a ScaleEvent which was triggered according to an active schedule.
func NewScheduledScaleEvent() ScaleEvent {
	return ScaleEvent{sType: scaleEventScheduled}
}

// NewScaleEvent creates a new ScaleEvent of type regular. This means a ScaleEvent which was triggered because of incoming scale alerts.
func NewScaleEvent(scaleFactor float32) ScaleEvent {
	return ScaleEvent{sType: scaleEventRegular, scaleFactor: scaleFactor}
}

// ScaleFactor returns the scalefactor
func (se ScaleEvent) ScaleFactor() float32 {
	return se.scaleFactor
}

func (se ScaleEvent) String() string {
	sfStr := fmt.Sprintf("%2f", se.scaleFactor)
	if se.sType == scaleEventScheduled {
		sfStr = "n/a"
	}
	return fmt.Sprintf("{type=%s,scaleFactor=%s}", se.sType, sfStr)
}
