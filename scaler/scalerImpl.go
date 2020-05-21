package scaler

import (
	"fmt"
	"strings"
	"time"

	m "github.com/thomasobenaus/sokar/metrics"
)

func (s *Scaler) scalingObjectWatcher(cycle time.Duration) {
	s.wg.Add(1)
	defer s.wg.Done()

	scalingObjectWatcherTicker := time.NewTicker(cycle)

	for {
		select {
		case <-s.stopChan:
			s.logger.Info().Msg("ScaleObjectWatcher Closed.")
			return
		case <-scalingObjectWatcherTicker.C:
			// Skip/ ignore the events for checking the current scale in case the
			// watcher is paused. This is usually the case if already a scaling is ongoing
			if !s.scalingObjectWatcherPaused {
				s.logger.Debug().Bool("watcher", true).Msgf("Check scalingObject state")
				if err := s.ensureScalingObjectCount(); err != nil {
					s.logger.Error().Bool("watcher", true).Msgf("Check scalingObject state failed: %s", err)
				}
			}
		}
	}
}

// scaleTicketProcessor listens on the given channel for incoming
// ScalingTickets to be processed.
func (s *Scaler) scaleTicketProcessor(ticketChan <-chan ScalingTicket) {

	defer s.wg.Done()
	s.logger.Info().Msg("ScaleTicketProcessor started.")

	for ticket := range ticketChan {
		s.scalingObjectWatcherPaused = true
		s.applyScaleTicket(ticket)
		s.scalingObjectWatcherPaused = false
	}

	s.logger.Info().Msg("ScaleTicketProcessor closed.")
}

func updateDesiredScale(sResult scaleResult, desiredScale *optionalValue) error {
	if desiredScale == nil {
		return fmt.Errorf("desiredScale parameter is nil")
	}

	if sResult.state != scaleDone {
		return nil
	}

	desiredScale.setValue(sResult.newCount)
	return nil
}

func updateScaleResultMetric(result scaleResult, scaleResultCounter m.CounterVec) {

	switch result.state {
	case scaleFailed:
		scaleResultCounter.WithLabelValues("failed").Inc()
	case scaleDone:
		scaleResultCounter.WithLabelValues("done").Inc()
	case scaleIgnored:
		scaleResultCounter.WithLabelValues("ignored").Inc()
	default:
		scaleResultCounter.WithLabelValues("other").Inc()
	}
}

// openScalingTicket opens based on the desired count a ScalingTicket
func (s *Scaler) openScalingTicket(desiredCount uint, force bool) error {

	if s.numOpenScalingTickets > s.maxOpenScalingTickets {
		s.metrics.scalingTicketCount.WithLabelValues("rejected").Inc()
		msg := fmt.Sprintf("Ticket rejected since currently a %d scaling tickets are open and only %d are allowed.", s.numOpenScalingTickets, s.maxOpenScalingTickets)
		s.logger.Debug().Msg(msg)
		return fmt.Errorf(msg)
	}

	s.metrics.scalingTicketCount.WithLabelValues("added").Inc()
	// TODO: Add metric "open scaling tickets"
	s.numOpenScalingTickets++
	s.scaleTicketChan <- NewScalingTicket(desiredCount, force)
	return nil
}

// applyScaleTicket applies the given ScalingTicket by issuing and tracking the scaling action.
func (s *Scaler) applyScaleTicket(ticket ScalingTicket) {
	ticket.start()
	result := s.scaleTo(ticket.desiredCount, ticket.force)
	if err := updateDesiredScale(result, &s.desiredScale); err != nil {
		s.logger.Error().Err(err).Msg("Failed updating desired scale.")
	}

	ticket.complete(result.state)
	s.numOpenScalingTickets--

	s.metrics.scalingTicketCount.WithLabelValues("applied").Inc()

	dur, _ := ticket.processingDuration()
	s.metrics.scalingDurationSeconds.Observe(float64(dur.Seconds()))
	updateScaleResultMetric(result, s.metrics.scaleResultCounter)

	logEvent := s.logger.Info()
	if result.state != scaleDone && result.state != scaleIgnored {
		logEvent = s.logger.Error()
	}
	logEvent.Msgf("Ticket applied. Scaling was %s (%s). New count is %d. Scaling in %f .", strings.ToUpper(string(result.state)), result.stateDescription, result.newCount, dur.Seconds())
}

func (s *Scaler) scaleTo(desiredCount uint, force bool) scaleResult {
	scalingObjectName := s.scalingObject.Name
	currentCount, err := s.scalingTarget.GetScalingObjectCount(scalingObjectName)
	if err != nil {
		return scaleResult{
			state:            scaleFailed,
			stateDescription: fmt.Sprintf("Error obtaining scalingObject count: %s.", err.Error()),
		}
	}

	return s.scale(desiredCount, currentCount, force)
}
