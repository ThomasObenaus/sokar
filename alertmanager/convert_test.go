package alertmanager

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_IsFiring(t *testing.T) {

	assert.True(t, isFiring("Firing"))
	assert.True(t, isFiring("firing"))
	assert.True(t, isFiring("FIRING"))
	assert.True(t, isFiring(" firing "))
	assert.False(t, isFiring("Resolved"))
	assert.False(t, isFiring("Firin"))
}

func Test_AlertToScalingAlert(t *testing.T) {

	name := "ABC"
	labels := map[string]string{"alertname": name}
	startedAt := time.Now()

	al := alert{
		Status:   "firing",
		StartsAt: startedAt,
		Labels:   labels,
	}

	convAlert := amAlertToScalingAlert(al)

	assert.Equal(t, name, convAlert.Name)
	assert.True(t, convAlert.Firing)
	assert.Equal(t, startedAt, convAlert.StartedAt)

	al = alert{
		Status:   "firing",
		StartsAt: startedAt,
	}

	convAlert = amAlertToScalingAlert(al)

	assert.Equal(t, "NO_NAME", convAlert.Name)
	assert.True(t, convAlert.Firing)
	assert.Equal(t, startedAt, convAlert.StartedAt)
}

func Test_AlertsToScalingAlertList(t *testing.T) {

	resp := response{Receiver: "sokar"}

	name1 := "ABC"
	labels := map[string]string{"alertname": name1}
	startedAt := time.Now()

	al1 := alert{
		Status:   "firing",
		StartsAt: startedAt,
		Labels:   labels,
	}
	resp.Alerts = append(resp.Alerts, al1)

	name2 := "XYZ"
	labels = map[string]string{"alertname": name2}

	al2 := alert{
		Status:   "resolved",
		StartsAt: startedAt,
		Labels:   labels,
	}

	resp.Alerts = append(resp.Alerts, al2)

	emitter, pkg := amResponseToScalingAlerts(resp)
	alerts := pkg.ScaleAlerts

	assert.Equal(t, "AM.sokar", emitter)

	assert.Equal(t, 2, len(alerts))
	assert.Equal(t, name1, alerts[0].Name)
	assert.True(t, alerts[0].Firing)

	assert.Equal(t, name2, alerts[1].Name)
	assert.False(t, alerts[1].Firing)
}
