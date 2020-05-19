package helper

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func SendScaleAlert(alertName string, fire bool) error {
	client := http.Client{
		Timeout: time.Millisecond * 500,
	}

	fireStr := "firing"
	if !fire {
		fireStr = "expired"
	}

	bodyStr := fmt.Sprintf(`{
	"receiver": "PM",
	"status": "%s",
	"alerts": [
	  {
		"status": "%s",
		"labels": {
		  "alertname": "%s",
		  "alert-type": "scaling",
		  "scale-type": "up"
		},
		"annotations": {
		  "description": "Scales the component XYZ UP"
		},
		"startsAt": "2019-02-23T12:00:00.000+01:00",
		"endsAt": "2019-02-23T12:05:00.000+01:00",
		"generatorURL": "http://generator_url"
	  }
	],
	"groupLabels": {},
	"commonLabels": { "alertname": "%s" },
	"commonAnnotations": {},
	"externalURL": "http://externalURL",
	"version": "4",
	"groupKey": "{}:{}"
  }`, fireStr, fireStr, alertName, alertName)

	bodybytes := []byte(bodyStr)
	body := bytes.NewReader(bodybytes)
	resp, err := client.Post("http://localhost:11000/api/alerts", "application/json", body)

	if err != nil {
		return errors.Wrap(err, "Error sending request to sokar.")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Response obtained from sokar is !200 but %d", resp.StatusCode)
	}

	return nil
}
