package serviceTest

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BuildAlertRequest(t *testing.T) {

	// no alert
	alerts := make([]alert, 0)
	request := buildAlertRequest(alerts)
	assert.Empty(t, request.Alerts)

	assert.Equal(t, "sokar", request.Receiver)
	assert.Equal(t, "firing", request.Status)
	assert.Equal(t, "localhost", request.ExternalURL)
	assert.Empty(t, request.GroupLabels)
	assert.Empty(t, request.CommonAnnotations)
	assert.Empty(t, request.CommonLabels)

	// one alert
	alerts = append(alerts, alert{
		Labels: map[string]string{"alertname": "Alert A"},
	})
	request = buildAlertRequest(alerts)
	assert.NotEmpty(t, request.Alerts)

	assert.Equal(t, "sokar", request.Receiver)
	assert.Equal(t, "firing", request.Status)
	assert.Equal(t, "localhost", request.ExternalURL)
	assert.Empty(t, request.GroupLabels)
	assert.Empty(t, request.CommonAnnotations)
	assert.Empty(t, request.CommonLabels)
}

func Test_New(t *testing.T) {

	am := newAlertManager("http://localhost", time.Second*2)

	require.NotNil(t, am.client)
	assert.Equal(t, time.Second*2, am.client.Timeout)
	assert.Equal(t, "http://localhost", am.sokarAddress)
}

func Test_RequestToStr(t *testing.T) {
	// one alert
	alerts := make([]alert, 0)
	alerts = append(alerts, alert{
		Labels: map[string]string{"alertname": "Alert A"},
	})
	request := buildAlertRequest(alerts)

	s, err := requestToStr(request)
	require.NoError(t, err)
	assert.NotEmpty(t, s)
	assert.Equal(t, `{"receiver":"sokar","status":"firing","alerts":[{"labels":{"alertname":"Alert A"},"startsAt":"0001-01-01T00:00:00Z","endsAt":"0001-01-01T00:00:00Z"}],"externalURL":"localhost"}`, s)
}

func Test_SendAlertmanagerRequest(t *testing.T) {

	// one alert
	alerts := make([]alert, 0)
	alerts = append(alerts, alert{
		Labels: map[string]string{"alertname": "Alert A"},
	})
	request := buildAlertRequest(alerts)
	require.NotEmpty(t, request.Alerts)

	sm := sokarMock{}
	sokarSrv := httptest.NewServer(&sm)
	sokarURL := sokarSrv.URL
	defer sokarSrv.Close()

	am := newAlertManager(sokarURL, time.Second*2)
	s, err := requestToStr(request)
	require.NoError(t, err)
	_, err = am.sendAlertmanagerRequest(s)
	require.NoError(t, err)

	require.NoError(t, sm.err)
	assert.NotEmpty(t, sm.amRequest.Alerts)
	assert.Equal(t, "Alert A", sm.amRequest.Alerts[0].Labels["alertname"])
}
