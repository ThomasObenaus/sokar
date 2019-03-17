package sokar

import (
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"
)

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

func (sk *Sokar) scaleEventProcessor(scaleEventChannel <-chan sokarIF.ScaleEvent) {
	sk.logger.Info().Msg("ScaleEventProcessor started.")

	for {
		select {
		case <-sk.stopChan:
			sk.logger.Info().Msg("ScaleEventProcessor stopped.")
			return
		case se := <-scaleEventChannel:
			sk.handleScaleEvent(se)
		}
	}
}

// Stop tears down sokar
func (sk *Sokar) Stop() {
	sk.logger.Info().Msg("Teardown requested")
	close(sk.stopChan)
}

// Join blocks/ waits until sokar has been stopped
func (sk *Sokar) Join() {
	<-sk.stopChan
}

func (sk *Sokar) handleScaleEvent(scaleEvent sokarIF.ScaleEvent) {

	sk.logger.Info().Msgf("SCALE-EVENT TRIGGERED: %v", scaleEvent)

	currentCount, err := sk.scaler.GetCount()
	if err != nil {
		sk.logger.Error().Err(err).Msg("Scaling ignored. Failed to obtain current count.")
		return
	}

	// plan
	plannedCount := sk.capacityPlanner.Plan(scaleEvent.ScaleFactor, currentCount)
	err = sk.scaler.ScaleTo(plannedCount)

	// HACK: For now we ignore all rejected scaling tickets
	if err != nil {
		sk.logger.Error().Err(err).Msg("Failed to scale.")
	}
}
