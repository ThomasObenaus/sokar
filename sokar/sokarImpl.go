package sokar

import (
	"fmt"

	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"
)

func (sk *Sokar) scaleEventProcessor(scaleEventChannel <-chan sokarIF.ScaleEvent) {
	sk.wg.Add(1)
	defer sk.wg.Done()

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

func scaleValueToScaleDir(scaleValue float32) (scaleDown bool) {
	if scaleValue < 0 {
		return true
	}
	return false
}

func (sk *Sokar) handleScaleEvent(scaleEvent sokarIF.ScaleEvent) {
	sk.logger.Info().Msgf("Scale Event received: %v", scaleEvent)

	scaleFactor := scaleEvent.ScaleFactor

	sk.metrics.scaleEventsTotal.Inc()
	sk.metrics.scaleFactor.Set(float64(scaleFactor))

	// this method is used for automatic mode only --> force has to be false
	err := sk.triggerScale(false, scaleFactor, sk.capacityPlanner.Plan)
	if err != nil {
		sk.logger.Error().Err(err).Msg("Failed to scale.")
	}
}

func (sk *Sokar) triggerScale(force bool, scaleValue float32, planFun func(scaleValue float32, currentScale uint) uint) error {

	scaleDown := scaleValueToScaleDir(scaleValue)

	preScaleJobCount, err := sk.scaler.GetCount()
	if err != nil {
		sk.metrics.failedScalingTotal.Inc()
		return fmt.Errorf("Failed to obtain current count. %s", err.Error())
	}
	sk.metrics.preScaleJobCount.Set(float64(preScaleJobCount))

	// Don't scale if sokar is in cool down mode
	if sk.capacityPlanner.IsCoolingDown(sk.scaler.GetTimeOfLastScaleAction(), scaleDown) {
		sk.metrics.skippedScalingDuringCooldownTotal.Inc()
		sk.logger.Info().Msg("Skip scale event. Sokar is cooling down.")
		return nil
	}

	// plan
	plannedJobCount := planFun(scaleValue, preScaleJobCount)
	sk.metrics.plannedJobCount.Set(float64(plannedJobCount))

	err = sk.scaler.ScaleTo(plannedJobCount, force)

	// HACK: For now we ignore all rejected scaling tickets
	if err != nil {
		sk.metrics.failedScalingTotal.Inc()
		return err
	}

	sk.logger.Info().Uint("preScaleCnt", preScaleJobCount).Uint("plannedCnt", plannedJobCount).Msg("Scaling triggered. Scaler will apply the planned count.")
	return nil
}
