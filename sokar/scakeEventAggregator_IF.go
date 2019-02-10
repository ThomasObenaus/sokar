package sokar

type ScaleEventAggregator interface {
	Substribe(eventChannel chan ScaleEvent)
}

type ScaleEvent struct {
	ScaleFactor float32
}
