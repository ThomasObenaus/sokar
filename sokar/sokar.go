package sokar

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"
)

// Sokar component that can be used to scale scaling-objects (jobs /instances).
type Sokar struct {
	logger zerolog.Logger

	// scaleEventEmitter is the component that provides the scale alerts to sokar
	scaleEventEmitter sokarIF.ScaleEventEmitter

	// capacityplanner is the component that plans the amount of instances to be scaled
	capacityPlanner sokarIF.CapacityPlanner

	// scaler is the component that does the actual scaling by sending
	// the needed commands to the scaling target (i.e. nomad)
	scaler sokarIF.Scaler

	// schedule is the currently active scale schedule
	schedule sokarIF.ScaleSchedule

	// metrics is a collection of metrics used by the sokar
	metrics Metrics

	// channel used to signal teardown/ stop
	stopChan chan struct{}

	wg sync.WaitGroup

	// dryRunMode is a flag that defines if sokar will execute its planned
	// scale actions or not. If the flag is true, sokar won't do anything beside planning.
	dryRunMode bool

	// the interval in which sokar emits 'scheduled' ScaleEvents in order to
	// force the evaluation if the current scale satisfies the given schedule
	scheduledScaleEventCycle time.Duration
}

// Config cfg for sokar
type Config struct {
	Logger zerolog.Logger

	DryRunMode bool
}

// New creates a new instance of sokar
func (cfg *Config) New(scaleEventEmitter sokarIF.ScaleEventEmitter, capacityPlanner sokarIF.CapacityPlanner, scaler sokarIF.Scaler, schedule sokarIF.ScaleSchedule, metrics Metrics) (*Sokar, error) {
	if scaler == nil {
		return nil, fmt.Errorf("Given Scaler is nil")
	}

	if capacityPlanner == nil {
		return nil, fmt.Errorf("Given CapacityPlanner is nil")
	}

	if scaleEventEmitter == nil {
		return nil, fmt.Errorf("Given ScaleEventEmitter is nil")
	}

	if schedule == nil {
		return nil, fmt.Errorf("Given ScaleSchedule is nil")
	}

	return &Sokar{
		scaleEventEmitter:        scaleEventEmitter,
		capacityPlanner:          capacityPlanner,
		scaler:                   scaler,
		stopChan:                 make(chan struct{}, 1),
		metrics:                  metrics,
		logger:                   cfg.Logger,
		dryRunMode:               cfg.DryRunMode,
		schedule:                 schedule,
		scheduledScaleEventCycle: time.Second * 30,
	}, nil
}

// Stop tears down sokar
func (sk *Sokar) Stop() error {
	sk.logger.Info().Msg("Teardown requested")
	close(sk.stopChan)
	return nil
}

// Join blocks/ waits until sokar has been stopped
func (sk *Sokar) Join() {
	sk.wg.Wait()
}

// String returns the name of this component
func (sk *Sokar) String() string {
	return "sokar"
}

// Start sokar
func (sk *Sokar) Start() {
	scaleEventChannel := make(chan sokarIF.ScaleEvent, 10)
	sk.scaleEventEmitter.Subscribe(scaleEventChannel)

	go sk.scaleEventProcessor(scaleEventChannel)

	// start the go routine for emitting scheduled ScaleEvents
	go sk.scaleEventScheduler(scaleEventChannel)
}
