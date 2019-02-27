package scaleAlertAggregator

import "time"

// ScaleAlertEmitter is a component emits ScaleAlerts.
type ScaleAlertEmitter interface {
	Subscribe(alertChannel chan ScaleAlertPacket)
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
	// Name of the component that has emitted the ScaleAlert's
	// of this packet.
	Emitter string

	ScaleAlerts []ScaleAlert
}
