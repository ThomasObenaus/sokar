package scaleEventAggregator

// ScaleAlertReceiver is a component that gathers scaling alerts and provides them on demand
type ScaleAlertReceiver interface {
	Subscribe(alertChannel chan ScaleAlert)
}

// ScaleAlert represents either a down or up-scale alert fired by an alerting system
type ScaleAlert struct {
}
