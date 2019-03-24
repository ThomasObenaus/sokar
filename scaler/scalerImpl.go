package scaler

import (
	"fmt"
	"time"
)

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
			s.logger.Error().Msg("Check job state (not implemented yet).")
		}
	}
}

// scaleTicketProcessor listens on the given channel for incoming
// ScalingTickets to be processed.
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

// applyScaleTicket applies the given ScalingTicket by issuing and tracking the scaling action.
func (s *Scaler) applyScaleTicket(ticket ScalingTicket) {
	ticket.start()
	result := s.scaleTo(ticket.desiredCount)
	ticket.complete(result.state)
	s.numOpenScalingTickets--

	s.metrics.scalingPolicyViolated.WithLabelValues("applied").Inc()

	// TODO: Add metric "Scaling duration"
	// TODO: Add metric "Scaling result success/failed/ignored"

	s.logger.Info().Msgf("Ticket applied. Scaling was %s (%s). New count is %d.", result.state, result.stateDescription, result.newCount)
}

// openScalingTicket opens based on the desired count a ScalingTicket
func (s *Scaler) openScalingTicket(desiredCount uint) error {

	if s.numOpenScalingTickets > s.maxOpenScalingTickets {
		s.metrics.scalingPolicyViolated.WithLabelValues("rejected").Inc()
		msg := fmt.Sprintf("Ticket rejected since currently a %d scaling tickets are open and only %d are allowed.", s.numOpenScalingTickets, s.maxOpenScalingTickets)
		s.logger.Debug().Msg(msg)
		return fmt.Errorf(msg)
	}
	s.metrics.scalingPolicyViolated.WithLabelValues("added").Inc()
	// TODO: Add metric "open scaling tickets"
	s.numOpenScalingTickets++
	s.scaleTicketChan <- NewScalingTicket(desiredCount)
	return nil
}

func (s *Scaler) scaleTo(desiredCount uint) scaleResult {
	jobName := s.job.jobName
	currentCount, err := s.scalingTarget.GetJobCount(jobName)
	if err != nil {
		return scaleResult{
			state:            scaleFailed,
			stateDescription: fmt.Sprintf("Error obtaining job count: %s.", err.Error()),
		}
	}

	return s.scale(desiredCount, currentCount)
}
