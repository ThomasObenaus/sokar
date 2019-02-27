package scaleAlertAggregator

import "time"

// ScaleAlertEmitter is a component emits ScaleAlerts.
type ScaleAlertEmitter interface {
	// Register is used to register the given handler func.
	// The ScaleAlertHandleFunc is called each time the ScaleAlertEmitter wants to promote
	// received alerts.
	Register(handleFunc ScaleAlertHandleFunc)
}

// ScaleAlert represents either a down or up-scale alert fired by an alerting system
type ScaleAlert struct {
	// Name of the alert.
	Name string
	// Firing is true if the alert is active, false otherwise.
	Firing bool
	// StartedAt represents the point in time the alert was created.
	StartedAt time.Time
}

// ScaleAlertPacket is a container for ScaleAlerts and meta information
type ScaleAlertPacket struct {
	ScaleAlerts []ScaleAlert
}

// ScaleAlertHandleFunc is a handler for received ScaleAlerts
type ScaleAlertHandleFunc func(emitter string, scaleAlerts ScaleAlertPacket)
