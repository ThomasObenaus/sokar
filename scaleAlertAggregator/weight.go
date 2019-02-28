package scaleAlertAggregator

import "time"

//weightPerSecondToWeight converts the given weight (per second) into an absolute weight
// based on the given aggregate cycle.
func weightPerSecondToWeight(weightPerSecond float32, aggregationCycle time.Duration) float32 {
	return float32(aggregationCycle.Seconds() * float64(weightPerSecond))
}

// getWeight returns the scale weight for the given alert.
// 0 is returned in case the weight for the given alert is not defined in the map
func getWeight(alertName string, weightMap ScaleAlertWeightMap) float32 {
	w, ok := weightMap[alertName]
	if !ok {
		return 0
	}
	return w
}
