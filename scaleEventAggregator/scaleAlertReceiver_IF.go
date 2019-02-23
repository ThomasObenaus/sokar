package scaleEventAggregator

import "time"

// ScaleAlertReceiver is a component that gathers scaling alerts and provides them on demand
type ScaleAlertReceiver interface {
	Subscribe(alertChannel chan ScaleAlertList)
}

// ScaleAlert represents either a down or up-scale alert fired by an alerting system
type ScaleAlert struct {
	Name      string
	Firing    bool
	StartedAt time.Time
}

// ScaleAlertList is a slice of ScaleAlert's
type ScaleAlertList []ScaleAlert
