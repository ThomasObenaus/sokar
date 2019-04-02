package scaler

import (
	"fmt"
	"time"
)

// ScalingTicket represents a ticket/ job to express the whish to scale
// and to track the state of the scaling.
type ScalingTicket struct {
	// issuedAt reflects the point in time the scaling ticket was
	// created issued
	issuedAt time.Time

	// startedAt reflects the point in time the scaling was started
	startedAt *time.Time
	// completedAt reflects the point in time the scaling was completed
	// (failed or successful)
	completedAt *time.Time

	// If this flag is true, then no scaling is executed in the end.
	// The scaler just checks against the scaling policy if a scaling would be needed.
	dryRun bool

	desiredCount uint
	state        scaleState
}

// NewScalingTicket creates and opens/ issues a new ScalingTicket
func NewScalingTicket(desiredCount uint, dryRun bool) ScalingTicket {
	return ScalingTicket{
		issuedAt:     time.Now(),
		startedAt:    nil,
		completedAt:  nil,
		desiredCount: desiredCount,
		state:        scaleNotStarted,
		dryRun:       dryRun,
	}
}

func (st *ScalingTicket) isInProgress() bool {
	return st.state == scaleRunning
}

func (st *ScalingTicket) start() {
	now := time.Now()
	st.startedAt = &now
	st.state = scaleRunning
}

func (st *ScalingTicket) complete(state scaleState) {
	now := time.Now()
	st.completedAt = &now
	st.state = state
}

// ticketAge reprents the duration from issuing/ creating
// the ticket until now.
func (st *ScalingTicket) ticketAge() time.Duration {
	return time.Now().Sub(st.issuedAt)
}

// leadDuration reprents the duration from issuing/ creating
// the ticket until it was completed.
func (st *ScalingTicket) leadDuration() (time.Duration, error) {

	if st.completedAt == nil {
		return 0, fmt.Errorf("Ticket not completed yet")
	}

	return st.completedAt.Sub(st.issuedAt), nil
}

// processingDuration reprents the duration from starting
// the ticket until it was completed.
func (st *ScalingTicket) processingDuration() (time.Duration, error) {

	if st.startedAt == nil {
		return 0, fmt.Errorf("Ticket not started yet")
	}

	if st.completedAt == nil {
		return 0, fmt.Errorf("Ticket not completed yet")
	}

	return st.completedAt.Sub(*st.startedAt), nil
}
