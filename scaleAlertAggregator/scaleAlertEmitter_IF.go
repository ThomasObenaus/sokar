package scaleAlertAggregator

import "time"

// ScaleAlertEmitter is a component that gathers scaling alerts and provides them on demand
type ScaleAlertEmitter interface {
	Subscribe(alertChannel chan ScaleAlertPacket)
}

// ScaleAlert represents either a down or up-scale alert fired by an alerting system
type ScaleAlert struct {
	Name      string
	Firing    bool
	StartedAt time.Time
}

// ScaleAlertPacket is a container for ScaleAlerts and meta information
type ScaleAlertPacket struct {
	Receiver    string
	ScaleAlerts []ScaleAlert
}
