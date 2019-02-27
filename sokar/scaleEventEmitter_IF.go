package sokar

// ScaleEventEmitter represents the interface for the component that fires ScaleEvents if needed.
type ScaleEventEmitter interface {
	Subscribe(eventChannel chan ScaleEvent)
}

// ScaleEvent is an event that is created each time potentially a scale should be made (up/down)
type ScaleEvent struct {
	ScaleFactor float32
}
