package sokar

const (
	// PathHealth is the url path for health end-point
	PathHealth = "/health"

	// PathMetrics path for the metrics end-point
	PathMetrics = "/metrics"

	// PathAlertmanager path for receiving alerts from the alertmanager
	PathAlertmanager = "/api/alerts"

	// PathScaleByPercentage is the scale-by end-point that allows scaling by percentage
	PathScaleByPercentage = "/api/scale-by/p" + "/:" + PathPartValue

	// PathScaleByValue is the scale-by end-point that allows scaling by value
	PathScaleByValue = "/api/scale-by/v" + "/:" + PathPartValue

	// PathPartValue represents a path part that takes a value
	PathPartValue = "value"
)
