package sokar

// ScaleAlertAggregator represents the interface for the component that gathers and aggregates ScaleEvents.
type ScaleAlertAggregator interface {
	Subscribe(eventChannel chan ScaleEvent)
}

// ScaleEvent is an event that is created each time potentially a scale should be made (up/down)
type ScaleEvent struct {
	ScaleFactor float32
}
