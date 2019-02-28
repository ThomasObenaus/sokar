package scaleAlertAggregator

// isScalingNeeded returns true if the current scaleCounter violates either the upScaling-
// or downScaling threshold
func (sc *ScaleAlertAggregator) isScalingNeeded() bool {
	scaleUpNeeded := sc.scaleCounter > sc.upScalingThreshold
	scaleDownNeeded := sc.scaleCounter < sc.downScalingThreshold

	return scaleDownNeeded || scaleUpNeeded
}
