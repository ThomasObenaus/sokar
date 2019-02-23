package alertmanager

import "time"

// response from the alertmanager
type response struct {
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
