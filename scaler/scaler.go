package scaler

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var oneDayAgo = time.Now().Add(time.Hour * -24)

// Scaler is a component responsible for scaling a scalingObject
type Scaler struct {
	logger zerolog.Logger

	// scalingTarget is the component that represents
	// the system that shall be used for scaling (i.e nomad)
	scalingTarget ScalingTarget

	// ScalingObject represents the ScalingObject and relevant meta data
	scalingObject ScalingObject

	// dryRunMode active/ not active. In dry run mode no automatic scaling will
	// executed. For more information see ../doc/DryRunMode.md
	dryRunMode bool

	// LastScaleAction represents that point in time
	// when the scaler was triggered to execute a scaling
	// action the last time
	lastScaleAction time.Time

	// watcherInterval the interval the Scaler will check if
	// the scalingObject count still matches the desired state.
	watcherInterval time.Duration

	// numOpenScalingTickets represents the number
	// of Scaling Tickets that where issued but not yet
	// applied.
	numOpenScalingTickets uint

	// maxOpenScalingTickets is a number that states how
	// many scaling tickets are accepted to be in state open
	// at the same time at max.
	maxOpenScalingTickets uint
	scaleTicketChan       chan ScalingTicket

	// desiredScale reflects the last successfully applied scale
	desiredScale optionalValue

	// channel used to signal teardown/ stop
	stopChan chan struct{}

	metrics Metrics

	wg sync.WaitGroup

	// scalingObjectWatcherPaused if true the scaling object won't be tracked to check if there is an adjustment needed
	scalingObjectWatcherPaused bool
}

// Option represents an option for the Scaler
type Option func(c *Scaler)

// WithLogger adds a configured Logger to the Scaler
func WithLogger(logger zerolog.Logger) Option {
	return func(s *Scaler) {
		s.logger = logger
	}
}

// MaxOpenScalingTickets specifies how many scaling tickets can be open at the same time
func MaxOpenScalingTickets(num uint) Option {
	return func(s *Scaler) {
		s.maxOpenScalingTickets = num
	}
}

// WatcherInterval specifies the interval the scaleObjectWatcher will use to check if the scale still
// matches the current expectation
func WatcherInterval(interval time.Duration) Option {
	return func(s *Scaler) {
		s.watcherInterval = interval
	}
}

// DryRunMode can be used to activate/ deactivate the dry run mode.
// In dry run mode no automatic scaling will executed.
// For more information see ../doc/DryRunMode.md
func DryRunMode(enable bool) Option {
	return func(s *Scaler) {
		s.dryRunMode = enable
	}
}

// New creates a new instance of a scaler using the given
// ScalingTarget to send scaling events to.
func New(scalingTarget ScalingTarget, scalingObject ScalingObject, metrics Metrics, options ...Option) (*Scaler, error) {

	maxOpenScalingTickets := uint(0)
	scaler := Scaler{
		scalingTarget:         scalingTarget,
		watcherInterval:       time.Second * 5,
		scalingObject:         scalingObject,
		stopChan:              make(chan struct{}, 1),
		numOpenScalingTickets: 0,
		maxOpenScalingTickets: maxOpenScalingTickets,
		metrics:               metrics,
		desiredScale:          optionalValue{isKnown: false},
		dryRunMode:            false,
		lastScaleAction:       oneDayAgo,
	}

	// apply the options
	for _, opt := range options {
		opt(&scaler)
	}

	scaler.scaleTicketChan = make(chan ScalingTicket, scaler.maxOpenScalingTickets+1)

	if scaler.scalingTarget == nil {
		return nil, fmt.Errorf("Given ScalingTarget is nil")
	}
	if scaler.scaleTicketChan == nil {
		return nil, fmt.Errorf("Scaling Ticket Channel is nil")
	}

	if scaler.watcherInterval <= time.Second*0 {
		return nil, fmt.Errorf("WatcherInterval is %s which is a too small value and thus not supported", scaler.watcherInterval.String())
	}

	return &scaler, nil
}

// GetCount returns the number of currently deployed instances
func (s *Scaler) GetCount() (uint, error) {
	return s.scalingTarget.GetScalingObjectCount(s.scalingObject.Name)
}

// ScaleTo will scale the scalingObject to the desired count.
func (s *Scaler) ScaleTo(desiredCount uint, force bool) error {
	s.logger.Info().Msgf("Scale to %d requested (force=%t).", desiredCount, force)
	return s.openScalingTicket(desiredCount, force)
}

// GetName returns the name of this component
func (s *Scaler) GetName() string {
	return "scaler"
}

// Run starts/ runs the scaler
func (s *Scaler) Run() {
	// handler that processes incoming scaling tickets
	s.wg.Add(1)
	go s.scaleTicketProcessor(s.scaleTicketChan)

	if s.dryRunMode {
		s.logger.Info().Msg("Don't start the ScalingObjectWatcher in dry-run mode.")
	} else {
		// handler that checks periodically if the desired count is still valid
		go s.scalingObjectWatcher(s.watcherInterval)
		s.logger.Info().Msg("ScalingObjectWatcher started.")
	}
}

// Stop tears down scaler
func (s *Scaler) Stop() error {
	s.logger.Info().Msg("Teardown requested")

	close(s.scaleTicketChan)
	close(s.stopChan)

	return nil
}

// Join blocks/ waits until scaler has been stopped
func (s *Scaler) Join() {
	s.wg.Wait()
}

// GetTimeOfLastScaleAction returns that point in time where the most recent
// scaling STARTED.
func (s *Scaler) GetTimeOfLastScaleAction() time.Time {
	return s.lastScaleAction
}
