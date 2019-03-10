package scaler

import (
	"fmt"
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

func (s *Scaler) jobWatcher(cycle time.Duration) {
	s.wg.Add(1)
	defer s.wg.Done()

	jobWatcherTicker := time.NewTicker(cycle)

	for {
		select {
		case <-s.stopChan:
			s.logger.Info().Msg("JobWatcher Closed.")
			return
		case <-jobWatcherTicker.C:
			s.logger.Error().Msgf("Check job state (not implemented yet). Desired %d.", s.desiredCount)
		}
	}
}

// scaleTicketProcessor
func (s *Scaler) scaleTicketProcessor(ticketChan <-chan ScalingTicket) {
	s.wg.Add(1)
	defer s.wg.Done()
	s.logger.Info().Msg("ScaleTicketProcessor started.")

	for ticket := range ticketChan {
		// TODO: Stop jobwatcher here
		s.applyScaleTicket(ticket)
		// TODO: Start jobwatcher here
	}

	s.logger.Info().Msg("ScaleTicketProcessor closed.")
}

// Run starts/ runs the scaler
func (s *Scaler) Run() {
	// handler that processes incoming scaling tickets
	go s.scaleTicketProcessor(s.scaleTicketChan)
	// handler that checks periodically if the desired count is still valid
	go s.jobWatcher(s.jobWatcherCycle)
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

func (s *Scaler) applyScaleTicket(ticket ScalingTicket) {
	ticket.start()
	result := s.scaleTo(ticket.desiredCount)
	ticket.complete(result.State)
	s.numOpenScalingTickets--

	s.logger.Info().Msgf("Ticket applied. Scaling was %s (%s). New count is %d.", result.State, result.StateDescription, result.NewCount)
}

func (s *Scaler) openScalingTicket(desiredCount uint) error {

	if s.numOpenScalingTickets > s.maxOpenScalingTickets {
		msg := fmt.Sprintf("Ticket rejected since currently a %d scaling tickets are open and only %d are allowed.", s.numOpenScalingTickets, s.maxOpenScalingTickets)
		s.logger.Debug().Msg(msg)
		return fmt.Errorf(msg)
	}
	s.numOpenScalingTickets++
	s.scaleTicketChan <- NewScalingTicket(desiredCount)
	return nil
}
