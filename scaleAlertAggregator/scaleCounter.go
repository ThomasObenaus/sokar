package scaleAlertAggregator

import "time"

type scaleCounter struct {
	startedAt time.Duration
	val       float32
}
