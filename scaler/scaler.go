package scaler

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Scaler is a component responsible for scaling a scalingObject
type Scaler struct {
	logger zerolog.Logger

	// scalingTarget is the component that represents
	// the system that shall be used for scaling (i.e nomad)
	scalingTarget ScalingTarget

	// scalingObjectCfg is the configuration for the scalingObject
	scalingObjectCfg scalingObjectConfig

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
}

// Config is the configuration for the Scaler
type Config struct {
	Name                  string
	MinCount              uint
	MaxCount              uint
	Logger                zerolog.Logger
	MaxOpenScalingTickets uint
	WatcherInterval       time.Duration
}

// scalingObjectConfig config of the scalingObject to be scaled
type scalingObjectConfig struct {
	name     string
	minCount uint
	maxCount uint
}

// New creates a new instance of a scaler using the given
// ScalingTarget to send scaling events to.
func (cfg Config) New(scalingTarget ScalingTarget, metrics Metrics) (*Scaler, error) {
	if scalingTarget == nil {
		return nil, fmt.Errorf("Given ScalingTarget is nil")
	}

	if cfg.WatcherInterval <= time.Second*0 {
		return nil, fmt.Errorf("WatcherInterval is %s which is a too small value and thus not supported", cfg.WatcherInterval.String())
	}

	return &Scaler{
		logger:          cfg.Logger,
		scalingTarget:   scalingTarget,
		watcherInterval: cfg.WatcherInterval,
		scalingObjectCfg: scalingObjectConfig{
			name:     cfg.Name,
			minCount: cfg.MinCount,
			maxCount: cfg.MaxCount,
		},
		stopChan:              make(chan struct{}, 1),
		numOpenScalingTickets: 0,
		maxOpenScalingTickets: cfg.MaxOpenScalingTickets,
		scaleTicketChan:       make(chan ScalingTicket, cfg.MaxOpenScalingTickets+1),
		metrics:               metrics,
		desiredScale:          optionalValue{isKnown: false},
	}, nil
}

// GetCount returns the number of currently deployed instances
func (s *Scaler) GetCount() (uint, error) {
	return s.scalingTarget.GetScalingObjectCount(s.scalingObjectCfg.name)
}

// ScaleTo will scale the scalingObject to the desired count.
func (s *Scaler) ScaleTo(desiredCount uint, dryRun bool) error {
	s.logger.Info().Msgf("Scale to %d requested (dryRun=%t).", desiredCount, dryRun)
	return s.openScalingTicket(desiredCount, dryRun)
}

// GetName returns the name of this component
func (s *Scaler) GetName() string {
	return "scaler"
}

// Run starts/ runs the scaler
func (s *Scaler) Run() {
	// handler that processes incoming scaling tickets
	go s.scaleTicketProcessor(s.scaleTicketChan)
	// handler that checks periodically if the desired count is still valid
	go s.scalingObjectWatcher(s.watcherInterval)
}

// Stop tears down scaler
func (s *Scaler) Stop() {
	s.logger.Info().Msg("Teardown requested")

	close(s.scaleTicketChan)
	close(s.stopChan)
}

// Join blocks/ waits until scaler has been stopped
func (s *Scaler) Join() {
	s.wg.Wait()
}
