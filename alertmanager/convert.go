package alertmanager

import (
	"strings"

	sea "github.com/thomasobenaus/sokar/scaleEventAggregator"
)

func genReceiver(name string) string {

	result := "AM"
	if len(name) > 0 {
		result += "." + name
	}

	return result
}

// amResponseToScalingAlerts extracts alerts from the response of the alertmanager
func amResponseToScalingAlerts(resp response) sea.ScaleAlertPacket {
	result := sea.ScaleAlertPacket{Receiver: genReceiver(resp.Receiver)}
	for _, alert := range resp.Alerts {
		result.ScaleAlerts = append(result.ScaleAlerts, amAlertToScalingAlert(alert))
	}

	return result
}

func amAlertToScalingAlert(alert alert) sea.ScaleAlert {

	name, ok := alert.Labels["alertname"]
	if !ok {
		name = "NO_NAME"
	}

	return sea.ScaleAlert{
		Name:      name,
		Firing:    isFiring(alert.Status),
		StartedAt: alert.StartsAt,
	}
}

func isFiring(status string) bool {
	status = strings.ToLower(status)
	status = strings.TrimSpace(status)
	return status == "firing"
}
