package sokar

import (
	"time"

	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"
)

func (sk *Sokar) scaleEventScheduler(scaleEventChannel chan sokarIF.ScaleEvent) {
	sk.wg.Add(1)
	defer sk.wg.Done()

	evaluationTicker := time.NewTicker(time.Second * 2)

	sk.logger.Info().Msg("ScaleEventScheduler started.")

	for {
		select {
		case <-sk.stopChan:
			sk.logger.Info().Msg("ScaleEventScheduler stopped.")
			return
		case <-evaluationTicker.C:

			if sk.shouldFireAlert(time.Now()) {
				sk.logger.Info().Msg("############# Fire scheduled scaling alert")
				event := sokarIF.ScaleEvent{ScaleFactor: 0}
				scaleEventChannel <- event
			}
		}
	}
}

func (sk *Sokar) shouldFireAlert(now time.Time) bool {
	// TODO: Fill
	//day := now.Weekday()
	//hour := now.Hour()
	//minute := now.Minute()
	//at, err := helper.NewTime(uint(hour), uint(minute))
	//
	//if err != nil {
	//	sk.logger.Warn().Msgf("Could not evaluate scaling schedule: %s", err.Error())
	//	return false
	//}
	//
	//return sk.schedule.IsActiveAt(day, at)
	return true
}
