package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/thomasobenaus/sokar/helper"
)

func (cfg *Config) fillCfgValues() error {
	// Context: main
	cfg.DryRunMode = cfg.viper.GetBool(dryRun.name)
	cfg.Port = cfg.viper.GetInt(port.name)

	// Context: Nomad
	cfg.Nomad.ServerAddr = cfg.viper.GetString(nomadServerAddress.name)

	// Context: job
	cfg.Job.Name = cfg.viper.GetString(jobName.name)
	min := cfg.viper.GetInt(jobMin.name)
	if min < 0 {
		min = 0
	}
	cfg.Job.MinCount = uint(min)

	max := cfg.viper.GetInt(jobMax.name)
	if max < 0 {
		max = 0
	}
	cfg.Job.MaxCount = uint(max)

	// Context: CapacityPlanner
	cfg.CapacityPlanner.DownScaleCooldownPeriod = cfg.viper.GetDuration(capDownScaleCoolDown.name)
	cfg.CapacityPlanner.UpScaleCooldownPeriod = cfg.viper.GetDuration(capUpScaleCoolDown.name)

	// Context: Logging
	cfg.Logging.Structured = cfg.viper.GetBool(loggingStructured.name)
	cfg.Logging.UxTimestamp = cfg.viper.GetBool(loggingUXTS.name)
	cfg.Logging.NoColoredLogOutput = cfg.viper.GetBool(loggingNoColor.name)

	// Context: ScaleAlertAggregator
	cfg.ScaleAlertAggregator.NoAlertScaleDamping = float32(cfg.viper.GetFloat64(saaNoAlertDamping.name))
	cfg.ScaleAlertAggregator.UpScaleThreshold = float32(cfg.viper.GetFloat64(saaUpThresh.name))
	cfg.ScaleAlertAggregator.DownScaleThreshold = float32(cfg.viper.GetFloat64(saaDownThresh.name))
	cfg.ScaleAlertAggregator.EvaluationCycle = cfg.viper.GetDuration(saaEvalCylce.name)

	evalPeriodFactor := cfg.viper.GetInt(saaEvalPeriodFactor.name)
	if evalPeriodFactor < 0 {
		evalPeriodFactor = 1
	}
	cfg.ScaleAlertAggregator.EvaluationPeriodFactor = uint(evalPeriodFactor)
	cfg.ScaleAlertAggregator.CleanupCycle = cfg.viper.GetDuration(saaCleanupCylce.name)

	alerts, err := extractAlertsFromViper(cfg.viper)
	if err != nil {
		return err
	}
	cfg.ScaleAlertAggregator.ScaleAlerts = alerts

	return nil
}

func extractAlertsFromViper(vp *viper.Viper) ([]Alert, error) {
	var alerts = make([]Alert, 0)

	if !vp.IsSet(saaScaleAlerts.name) {
		return nil, nil
	}

	alertsAsStr := vp.GetString(saaScaleAlerts.name)

	if len(alertsAsStr) > 0 {
		return alertStrToAlerts(alertsAsStr)
	}

	alertsAsMap := helper.CastToStringMapSlice(vp.Get(saaScaleAlerts.name))
	if alertsAsMap == nil {
		return alerts, nil
	}

	alerts, err := alertMapToAlerts(alertsAsMap)
	if err != nil {
		return alerts, fmt.Errorf("Error reading alerts configuration: %s", err.Error())
	}
	return alerts, nil
}

func alertMapToAlerts(alertCfg []map[string]string) ([]Alert, error) {

	if alertCfg == nil {
		return nil, fmt.Errorf("Parameter is nil")
	}
	var alerts = make([]Alert, 0)

	for _, alert := range alertCfg {
		name := alert["name"]
		if len(name) == 0 {
			return nil, fmt.Errorf("Name for alert is missing")
		}
		weightStr := alert["weight"]
		description := alert["description"]
		weight, err := cast.ToFloat32E(weightStr)
		if err != nil {
			return nil, fmt.Errorf("Failed while reading weight for %s: %s", name, err.Error())
		}
		alerts = append(alerts, Alert{Name: name, Weight: weight, Description: description})
	}

	return alerts, nil
}

func alertStrToAlerts(alertsAsStr string) ([]Alert, error) {
	var alerts = make([]Alert, 0)

	alertStrSplit := strings.Split(alertsAsStr, ";")
	for _, element := range alertStrSplit {
		if len(element) == 0 {
			continue
		}

		alertStr := strings.Split(element, ":")
		if len(alertStr) < 2 {
			return nil, fmt.Errorf("Unable to read alert config. An alert consists of a key value pair (name:weight). This one does not '%s'", element)
		}
		name := strings.TrimSpace(alertStr[0])
		weight, err := cast.ToFloat32E(strings.TrimSpace(alertStr[1]))
		if err != nil {
			return nil, fmt.Errorf("Unable to read alert config: %s", err.Error())
		}

		description := ""

		if len(alertStr) > 2 {
			description = strings.TrimSpace(alertStr[2])
		}

		alerts = append(alerts, Alert{Name: name, Weight: weight, Description: description})
	}

	return alerts, nil
}
