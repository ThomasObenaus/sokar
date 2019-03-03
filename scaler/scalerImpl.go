package scaler

import (
	"time"
)

type jobConfig struct {
	jobName  string
	minCount uint
	maxCount uint
}

type policyCheckResult struct {
	validCount        uint
	desiredCount      uint
	minPolicyViolated bool
	maxPolicyViolated bool
}

func (s *Scaler) Run() {
	jobWatcherTicker := time.NewTicker(s.jobWatcherCycle)

	// main loop
	go func() {
		s.logger.Info().Msg("Main loop started")

	loop:
		for {
			select {
			case <-s.stopChan:
				close(s.stopChan)
				break loop

			case <-jobWatcherTicker.C:
				s.logger.Error().Msg("Check job state (not implemented yet).")
			}
		}
		s.logger.Info().Msg("Main loop left")
	}()

}

// Stop tears down scaler
func (s *Scaler) Stop() {
	s.logger.Info().Msg("Teardown requested")
	// send the stop message
	s.stopChan <- struct{}{}
}

// Join blocks/ waits until scaler has been stopped
func (s *Scaler) Join() {
	<-s.stopChan
}
