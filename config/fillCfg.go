package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/thomasobenaus/sokar/helper"
)

func (cfg *Config) fillScaler() error {
	cfg.Scaler.WatcherInterval = cfg.viper.GetDuration(scaWatcherInterval.name)

	scaModeStr := cfg.viper.GetString(scaMode.name)
	scaMode, err := strToScalerMode(scaModeStr)
	if err != nil {
		return err
	}
	cfg.Scaler.Mode = scaMode

	// TODO: DEPRECATED, Remove it
	scaNomadModeStr := cfg.viper.GetString(scaNomadMode.name)
	if len(scaNomadModeStr) > 0 {
		scaMode, err = strToScalerMode(scaNomadModeStr)
		if err != nil {
			return err
		}
		cfg.Scaler.Mode = scaMode
	}

	// Context: Scaler - AWS EC2
	cfg.Scaler.AwsEc2.Profile = cfg.viper.GetString(scaAWSEC2Profile.name)
	cfg.Scaler.AwsEc2.Region = cfg.viper.GetString(scaAWSEC2Region.name)

	// Context: Scaler - Nomad
	cfg.Scaler.Nomad.ServerAddr = cfg.viper.GetString(scaNomadModeServerAddress.name)
	cfg.Scaler.Nomad.DataCenterAWS.Profile = cfg.viper.GetString(scaNomadDataCenterAWSProfile.name)
	cfg.Scaler.Nomad.DataCenterAWS.Region = cfg.viper.GetString(scaNomadDataCenterAWSRegion.name)

	return validateScaler(cfg.Scaler)
}

func validateScaler(scaler Scaler) error {

	switch mode := scaler.Mode; mode {
	case ScalerModeNomadJob:
		if len(scaler.Nomad.ServerAddr) == 0 {
			return fmt.Errorf("The parameter '%s' is missing but this is needed in Scaler.Mode '%v'", scaNomadModeServerAddress.name, mode)
		}
	case ScalerModeNomadDataCenter:
		if len(scaler.Nomad.ServerAddr) == 0 {
			return fmt.Errorf("The parameter '%s' is missing but this is needed in Scaler.Mode '%v'", scaNomadModeServerAddress.name, mode)
		}
		if len(scaler.Nomad.DataCenterAWS.Profile) == 0 {
			return fmt.Errorf("The parameter '%s' is missing but this is needed in Scaler.Mode '%v'", scaNomadDataCenterAWSProfile.name, mode)
		}
		if len(scaler.Nomad.DataCenterAWS.Region) == 0 {
			return fmt.Errorf("The parameter '%s' is missing but this is needed in Scaler.Mode '%v'", scaNomadDataCenterAWSRegion.name, mode)
		}
	case ScalerModeAwsEc2:
		if len(scaler.AwsEc2.Profile) == 0 {
			return fmt.Errorf("The parameter '%s' is missing but this is needed in Scaler.Mode '%v'", scaAWSEC2Profile.name, mode)
		}
		if len(scaler.AwsEc2.Region) == 0 {
			return fmt.Errorf("The parameter '%s' is missing but this is needed in Scaler.Mode '%v'", scaAWSEC2Region.name, mode)
		}
	default:
		return fmt.Errorf("The parameter '%s' is missing but this is needed in Scaler.Mode '%v'", scaMode.name, mode)
	}

	if scaler.WatcherInterval <= time.Millisecond*500 {
		return fmt.Errorf("'%s' can't be less then 500ms", scaWatcherInterval.name)
	}

	return nil
}

func (cfg *Config) fillCfgValues() error {
	// Context: main
	cfg.DryRunMode = cfg.viper.GetBool(dryRun.name)
	cfg.Port = cfg.viper.GetInt(port.name)

	// Context: Scaler
	cfg.fillScaler()

	// Context: scale object
	cfg.ScaleObject.Name = cfg.viper.GetString(scaleObjectName.name)
	min := cfg.viper.GetInt(scaleObjectMin.name)
	if min < 0 {
		min = 0
	}
	cfg.ScaleObject.MinCount = uint(min)

	max := cfg.viper.GetInt(scaleObjectMax.name)
	if max < 0 {
		max = 0
	}
	cfg.ScaleObject.MaxCount = uint(max)

	// Context: CapacityPlanner
	cfg.CapacityPlanner.DownScaleCooldownPeriod = cfg.viper.GetDuration(capDownScaleCoolDown.name)
	cfg.CapacityPlanner.UpScaleCooldownPeriod = cfg.viper.GetDuration(capUpScaleCoolDown.name)

	cfg.CapacityPlanner.ConstantMode.Enable = cfg.viper.GetBool(capConstantModeEnable.name)
	constModeOffset := cfg.viper.GetInt(capConstantModeOffset.name)
	if constModeOffset <= 0 {
		constModeOffset = 1
	}
	cfg.CapacityPlanner.ConstantMode.Offset = uint(constModeOffset)
	cfg.CapacityPlanner.LinearMode.Enable = cfg.viper.GetBool(capLinearModeEnable.name)
	cfg.CapacityPlanner.LinearMode.ScaleFactorWeight = cfg.viper.GetFloat64(capLinearModeScaleFactorWeight.name)

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
	cfg.ScaleAlertAggregator.AlertExpirationTime = cfg.viper.GetDuration(saaAlertExpirationTime.name)

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

func strToScalerMode(mode string) (ScalerMode, error) {

	mode = strings.TrimSpace(mode)
	if len(mode) == 0 {
		return "", fmt.Errorf("Can't parse ScalerMode since input is empty")
	}

	mode = strings.ToLower(mode)
	if mode == "job" {
		return ScalerModeNomadJob, nil
	}
	if mode == string(ScalerModeNomadJob) {
		return ScalerModeNomadJob, nil
	}
	if mode == string(ScalerModeNomadDataCenter) {
		return ScalerModeNomadDataCenter, nil
	}
	if mode == "dc" {
		return ScalerModeNomadDataCenter, nil
	}
	if mode == string(ScalerModeAwsEc2) {
		return ScalerModeAwsEc2, nil
	}

	return "", fmt.Errorf("Can't parse '%s' to ScalerMode. Given value is unknown", mode)
}
