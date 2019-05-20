package serviceTest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/thomasobenaus/sokar/sokar"
)

type alertmanager struct {
	sokarAddress string
	client       *http.Client
}

// request send from an alertmanger to sokar
type request struct {
	Receiver          string  `json:"receiver,omitempty"`
	Status            string  `json:"status,omitempty"`
	Alerts            []alert `json:"alerts,omitempty"`
	GroupLabels       kvp     `json:"groupLabels,omitempty"`
	CommonLabels      kvp     `json:"commonLabels,omitempty"`
	CommonAnnotations kvp     `json:"commonAnnotations,omitempty"`
	ExternalURL       string  `json:"externalURL,omitempty"`
	Version           string  `json:"version,omitempty"`
	GroupKey          string  `json:"groupKey,omitempty"`
}

// alert data from alertmanager
type alert struct {
	Status       string    `json:"status,omitempty"`
	Labels       kvp       `json:"labels,omitempty"`
	Annotations  kvp       `json:"annotations,omitempty"`
	StartsAt     time.Time `json:"startsAt,omitempty"`
	EndsAt       time.Time `json:"endsAt,omitempty"`
	GeneratorURL string    `json:"generatorURL,omitempty"`
}

// kvp a key value pair
type kvp map[string]string

func newAlertManager(sokarAddress string, timeout time.Duration) alertmanager {
	return alertmanager{
		sokarAddress: sokarAddress,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func buildAlertRequest(alerts []alert) request {
	request := request{
		Receiver:    "sokar",
		Status:      "firing",
		ExternalURL: "localhost",
		Alerts:      alerts,
	}

	return request
}

func requestToStr(request request) (string, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	return string(data[:len(data)]), nil
}

func (am *alertmanager) sendAlertmanagerRequest(request string) (int, error) {

	req, err := http.NewRequest("POST", am.sokarAddress+sokar.PathAlertmanager, strings.NewReader(request))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := am.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
