package sokar

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"
)

// Sokar component that can be used to scale jobs/instances
type Sokar struct {
	logger zerolog.Logger

	scaler            sokarIF.Scaler
	capacityPlanner   sokarIF.CapacityPlanner
	scaleEventEmitter sokarIF.ScaleEventEmitter

	// channel used to signal teardown/ stop
	stopChan chan struct{}

	wg sync.WaitGroup
}

// Config cfg for sokar
type Config struct {
	Logger zerolog.Logger
}

// New creates a new instance of sokar
func (cfg *Config) New(scaleEventEmitter sokarIF.ScaleEventEmitter, capacityPlanner sokarIF.CapacityPlanner, scaler sokarIF.Scaler) (*Sokar, error) {
	if scaler == nil {
		return nil, fmt.Errorf("Given Scaler is nil")
	}

	if capacityPlanner == nil {
		return nil, fmt.Errorf("Given CapacityPlanner is nil")
	}

	if scaleEventEmitter == nil {
		return nil, fmt.Errorf("Given ScaleEventEmitter is nil")
	}

	return &Sokar{
		scaleEventEmitter: scaleEventEmitter,
		capacityPlanner:   capacityPlanner,
		scaler:            scaler,
		stopChan:          make(chan struct{}, 1),

		logger: cfg.Logger,
	}, nil
}

// Stop tears down sokar
func (sk *Sokar) Stop() {
	sk.logger.Info().Msg("Teardown requested")
	close(sk.stopChan)
}

// Join blocks/ waits until sokar has been stopped
func (sk *Sokar) Join() {
	sk.wg.Wait()
}

// GetName returns the name of this component
func (sk *Sokar) GetName() string {
	return "sokar"
}

// Run starts sokar
func (sk *Sokar) Run() {
	scaleEventChannel := make(chan sokarIF.ScaleEvent, 10)
	sk.scaleEventEmitter.Subscribe(scaleEventChannel)

	go sk.scaleEventProcessor(scaleEventChannel)
}
