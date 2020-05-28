package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"

	cfglib "github.com/ThomasObenaus/go-base/config"
	"github.com/spf13/cast"
	"github.com/thomasobenaus/sokar/helper"
)

func (cfg *Config) fillScaler(provider cfglib.Provider) error {
	cfg.Scaler.WatcherInterval = provider.GetDuration(scaWatcherInterval.Name())

	scaModeStr := provider.GetString(scaMode.Name())
	scaMode, err := strToScalerMode(scaModeStr)
	if err != nil {
		return err
	}
	cfg.Scaler.Mode = scaMode

	// Context: Scaler - AWS EC2
	cfg.Scaler.AwsEc2.Profile = provider.GetString(scaAWSEC2Profile.Name())
	cfg.Scaler.AwsEc2.Region = provider.GetString(scaAWSEC2Region.Name())
	cfg.Scaler.AwsEc2.ASGTagKey = provider.GetString(scaAWSEC2ASGTagKey.Name())
	// Context: Scaler - Nomad
	cfg.Scaler.Nomad.ServerAddr = provider.GetString(scaNomadModeServerAddress.Name())
	cfg.Scaler.Nomad.DataCenterAWS.Profile = provider.GetString(scaNomadDataCenterAWSProfile.Name())
	cfg.Scaler.Nomad.DataCenterAWS.Region = provider.GetString(scaNomadDataCenterAWSRegion.Name())
	cfg.Scaler.Nomad.DataCenterAWS.InstanceTerminationTimeout = provider.GetDuration(scaNomadDataCenterAWSInstanceTerminationTimeout.Name())

	return validateScaler(cfg.Scaler)
}

func validateScaler(scaler Scaler) error {
	const parameterMissingErrorPattern = "The parameter '%s' is missing but this is needed in Scaler.Mode '%v'"

	switch mode := scaler.Mode; mode {
	case ScalerModeNomadJob:
		if len(scaler.Nomad.ServerAddr) == 0 {
			return fmt.Errorf(parameterMissingErrorPattern, scaNomadModeServerAddress.Name(), mode)
		}
	case ScalerModeNomadDataCenter:
		hasRegion := len(scaler.Nomad.DataCenterAWS.Region) > 0
		hasProfile := len(scaler.Nomad.DataCenterAWS.Profile) > 0
		if len(scaler.Nomad.ServerAddr) == 0 {
			return fmt.Errorf(parameterMissingErrorPattern, scaNomadModeServerAddress.Name(), mode)
		}
		if !hasProfile && !hasRegion {
			return fmt.Errorf("The parameter '%s' and '%s' are missing but one of both is needed in Scaler.Mode '%v'", scaNomadDataCenterAWSProfile.Name(), scaNomadDataCenterAWSRegion.Name(), mode)
		}
	case ScalerModeAwsEc2:
		hasRegion := len(scaler.AwsEc2.Region) > 0
		hasProfile := len(scaler.AwsEc2.Profile) > 0

		if !hasProfile && !hasRegion {
			return fmt.Errorf("The parameter '%s' and '%s' are missing but one of both is needed in Scaler.Mode '%v'", scaAWSEC2Profile.Name(), scaAWSEC2Region.Name(), mode)
		}
		if len(scaler.AwsEc2.ASGTagKey) == 0 {
			return fmt.Errorf(parameterMissingErrorPattern, scaAWSEC2ASGTagKey.Name(), mode)
		}
	default:
		return fmt.Errorf(parameterMissingErrorPattern, scaMode.Name(), mode)
	}

	if scaler.WatcherInterval <= time.Millisecond*500 {
		return fmt.Errorf("'%s' can't be less then 500ms", scaWatcherInterval.Name())
	}

	return nil
}

func (cfg *Config) fillCapacityPlanner(provider cfglib.Provider) error {

	// Context: CapacityPlanner
	cfg.CapacityPlanner.DownScaleCooldownPeriod = provider.GetDuration(capDownScaleCoolDown.Name())
	cfg.CapacityPlanner.UpScaleCooldownPeriod = provider.GetDuration(capUpScaleCoolDown.Name())

	cfg.CapacityPlanner.ConstantMode.Enable = provider.GetBool(capConstantModeEnable.Name())
	constModeOffset := provider.GetInt(capConstantModeOffset.Name())
	if constModeOffset <= 0 {
		constModeOffset = 1
	}
	cfg.CapacityPlanner.ConstantMode.Offset = uint(constModeOffset)
	cfg.CapacityPlanner.LinearMode.Enable = provider.GetBool(capLinearModeEnable.Name())
	cfg.CapacityPlanner.LinearMode.ScaleFactorWeight = provider.GetFloat64(capLinearModeScaleFactorWeight.Name())

	if cfg.CapacityPlanner.LinearMode.Enable {
		cfg.CapacityPlanner.ConstantMode.Enable = false
	}

	entries, err := extractScaleScheduleFromViper(provider)
	if err != nil {
		return err
	}
	cfg.CapacityPlanner.ScaleSchedule = entries

	return validateCapacityPlanner(cfg.CapacityPlanner)
}

func parseScalingScheduleEntries(raw string) ([]ScaleScheduleEntry, error) {
	parts := strings.Split(raw, "|")
	result := make([]ScaleScheduleEntry, 0)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			continue
		}
		entry, err := NewScaleScheduleEntry(part)
		if err != nil {
			return make([]ScaleScheduleEntry, 0), err
		}
		result = append(result, entry)
	}

	return result, nil
}

func validateCapacityPlanner(capacityPlanner CapacityPlanner) error {

	if capacityPlanner.ConstantMode.Enable && capacityPlanner.LinearMode.Enable {
		return fmt.Errorf("constant and linear mode are set at the same time, this is not allowed")
	}

	if !capacityPlanner.ConstantMode.Enable && !capacityPlanner.LinearMode.Enable {
		return fmt.Errorf("neither constant nor linear mode are set, this is not allowed")
	}

	if capacityPlanner.LinearMode.Enable && capacityPlanner.LinearMode.ScaleFactorWeight <= 0 {
		return fmt.Errorf("invalid scale factor (%f) for linear mode, it has to be greater than 0", capacityPlanner.LinearMode.ScaleFactorWeight)
	}

	if capacityPlanner.ConstantMode.Enable && capacityPlanner.ConstantMode.Offset == 0 {
		return fmt.Errorf("invalid offset (%d) for constant mode, it has to be greater than 0", capacityPlanner.ConstantMode.Offset)
	}

	return nil
}

func (cfg *Config) fillCfgValues(provider cfglib.Provider) error {
	// Context: main
	cfg.DryRunMode = provider.GetBool(dryRun.Name())
	cfg.Port = provider.GetInt(port.Name())

	// Context: Scaler
	err := cfg.fillScaler(provider)
	if err != nil {
		return err
	}

	// Context: scale object
	cfg.ScaleObject.Name = provider.GetString(scaleObjectName.Name())
	min := provider.GetInt(scaleObjectMin.Name())
	if min < 0 {
		min = 0
	}
	cfg.ScaleObject.MinCount = uint(min)

	max := provider.GetInt(scaleObjectMax.Name())
	if max < 0 {
		max = 0
	}
	cfg.ScaleObject.MaxCount = uint(max)

	// Context: CapacityPlanner
	err = cfg.fillCapacityPlanner(provider)
	if err != nil {
		return err
	}
	// Context: Logging
	err = cfg.fillLoggingContext(provider)
	if err != nil {
		return err
	}

	// Context: ScaleAlertAggregator
	cfg.ScaleAlertAggregator.NoAlertScaleDamping = float32(provider.GetFloat64(saaNoAlertDamping.Name()))
	cfg.ScaleAlertAggregator.UpScaleThreshold = float32(provider.GetFloat64(saaUpThresh.Name()))
	cfg.ScaleAlertAggregator.DownScaleThreshold = float32(provider.GetFloat64(saaDownThresh.Name()))
	cfg.ScaleAlertAggregator.EvaluationCycle = provider.GetDuration(saaEvalCylce.Name())

	evalPeriodFactor := provider.GetInt(saaEvalPeriodFactor.Name())
	if evalPeriodFactor < 0 {
		evalPeriodFactor = 1
	}
	cfg.ScaleAlertAggregator.EvaluationPeriodFactor = uint(evalPeriodFactor)
	cfg.ScaleAlertAggregator.CleanupCycle = provider.GetDuration(saaCleanupCylce.Name())

	alerts, err := extractAlertsFromViper(provider)
	if err != nil {
		return err
	}
	cfg.ScaleAlertAggregator.ScaleAlerts = alerts
	cfg.ScaleAlertAggregator.AlertExpirationTime = provider.GetDuration(saaAlertExpirationTime.Name())

	return nil
}

func (cfg *Config) fillLoggingContext(provider cfglib.Provider) error {
	cfg.Logging.Structured = provider.GetBool(loggingStructured.Name())
	cfg.Logging.UxTimestamp = provider.GetBool(loggingUXTS.Name())
	cfg.Logging.NoColoredLogOutput = provider.GetBool(loggingNoColor.Name())

	level, err := strToLogLevel(provider.GetString(loggingLevel.Name()))
	cfg.Logging.Level = level
	return err
}

func strToLogLevel(v string) (zerolog.Level, error) {

	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	switch v {
	case "debug":
		return zerolog.DebugLevel, nil
	case "info":
		return zerolog.InfoLevel, nil
	case "warn":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	case "fatal":
		return zerolog.FatalLevel, nil
	case "off":
		return zerolog.Disabled, nil
	}

	return zerolog.NoLevel, fmt.Errorf("Invalid loglevel '%s'. Only debug, info, warn, error, fatal and off is supported", v)
}

func extractAlertsFromViper(provider cfglib.Provider) ([]Alert, error) {
	var alerts = make([]Alert, 0)

	if !provider.IsSet(saaScaleAlerts.Name()) {
		return nil, nil
	}

	alertsAsStr := provider.GetString(saaScaleAlerts.Name())

	if len(alertsAsStr) > 0 {
		return alertStrToAlerts(alertsAsStr)
	}

	alertsAsMap := helper.CastToStringMapSlice(provider.Get(saaScaleAlerts.Name()))
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

func extractScaleScheduleFromViper(provider cfglib.Provider) ([]ScaleScheduleEntry, error) {
	var scaleSchedule = make([]ScaleScheduleEntry, 0)

	if !provider.IsSet(capScaleSchedule.Name()) {
		return nil, nil
	}

	scaleScheduleAsStr := provider.GetString(capScaleSchedule.Name())
	if len(scaleScheduleAsStr) > 0 {
		return parseScalingScheduleEntries(scaleScheduleAsStr)
	}

	scaleScheduleAsMap := helper.CastToStringMapSlice(provider.Get(capScaleSchedule.Name()))
	if scaleScheduleAsMap == nil {
		return scaleSchedule, nil
	}

	scaleSchedule, err := scaleScheduleMapToScaleSchedule(scaleScheduleAsMap)
	if err != nil {
		return scaleSchedule, fmt.Errorf("Error reading scale schedule configuration: %s", err.Error())
	}
	return scaleSchedule, nil
}

func scaleScheduleMapToScaleSchedule(scaleScheduleCfg []map[string]string) ([]ScaleScheduleEntry, error) {

	if scaleScheduleCfg == nil {
		return nil, fmt.Errorf("Parameter is nil")
	}
	var scaleSchedule = make([]ScaleScheduleEntry, 0)

	for _, scheduleEntry := range scaleScheduleCfg {
		days := strings.TrimSpace(scheduleEntry["days"])
		if len(days) == 0 {
			return nil, fmt.Errorf("Days is missing for scale schedule entry")
		}

		startTime := strings.TrimSpace(scheduleEntry["start-time"])
		if len(startTime) == 0 {
			return nil, fmt.Errorf("StartTime is missing for scale schedule entry")
		}

		endTime := strings.TrimSpace(scheduleEntry["end-time"])
		if len(endTime) == 0 {
			return nil, fmt.Errorf("EndTime is missing for scale schedule entry")
		}

		min := strings.TrimSpace(scheduleEntry["min"])
		if len(min) == 0 {
			return nil, fmt.Errorf("Min is missing for scale schedule entry")
		}

		max := strings.TrimSpace(scheduleEntry["max"])
		if len(max) == 0 {
			return nil, fmt.Errorf("Max is missing for scale schedule entry")
		}

		spec := fmt.Sprintf("%s %s %s %s-%s", days, startTime, endTime, min, max)
		sse, err := NewScaleScheduleEntry(spec)
		if err != nil {
			return nil, fmt.Errorf("ScaleScheduleSpec malformed (%s): %s", spec, err.Error())
		}
		scaleSchedule = append(scaleSchedule, sse)
	}

	return scaleSchedule, nil
}
