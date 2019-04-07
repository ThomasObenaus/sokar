package config

import (
	"log"
	"time"
)

// ###################### Context: main ####################################################
var configFile = configEntry{
	name:         "config-file",
	bindFlag:     true,
	defaultValue: "",
	usage:        "Specifies the full path and name of the configuration file for sokar.",
}

var dryRun = configEntry{
	name:         "dry-run",
	bindFlag:     true,
	bindEnv:      true,
	defaultValue: false,
	usage:        "If true, then sokar won't execute the planned scaling action. Only scaling actions triggered via ScaleBy end-point will be executed.",
}

var port = configEntry{
	name:         "port",
	bindFlag:     true,
	bindEnv:      true,
	defaultValue: 11000,
	usage:        "Port where sokar is listening.",
}

// ###################### Context: nomad ####################################################
var nomadServerAddress = configEntry{
	name:         "nomad.server-address",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: "",
	usage:        "Specifies the address of the nomad server.",
}

// ###################### Context: job ####################################################
var jobName = configEntry{
	name:         "job.name",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: "",
	usage:        "The name of the job to be scaled.",
}

var jobMin = configEntry{
	name:         "job.min",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: 1,
	usage:        "The minimum scale of the job.",
}

var jobMax = configEntry{
	name:         "job.max",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: 10,
	usage:        "The maximum scale of the job.",
}

// ###################### Context: CapacityPlanner#########################################
var capaDownScaleCoolDown = configEntry{
	name:         "cap.down-scale-cooldown",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: time.Second * 20,
	usage:        "The time sokar waits between downscaling actions at min.",
}

var capaUpScaleCoolDown = configEntry{
	name:         "cap.up-scale-cooldown",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: time.Second * 10,
	usage:        "The time sokar waits between upscaling actions at min.",
}

// ###################### Context: Logging ################################################
var loggingStructured = configEntry{
	name:         "logging.structured",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: false,
	usage:        "Use structured logging or not.",
}

var loggingUXTS = configEntry{
	name:         "logging.unix-ts",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: false,
	usage:        "Use Unix-Timestamp representation for log entries.",
}

// ###################### Context: ScaleAlertAggregator ###################################
var saaNoAlertDamping = configEntry{
	name:         "saa.no-alert-damping",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: 1.0,
	usage:        "Damping used in case there are currently no alerts firing (neither down- nor upscaling).",
}

var saaUpThresh = configEntry{
	name:         "saa.up-thresh",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: 10.0,
	usage:        "Threshold for a upscaling event.",
}
var saaDownThresh = configEntry{
	name:         "saa.down-thresh",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: -10.0,
	usage:        "Threshold for a downscaling event.",
}

var saaEvalCylce = configEntry{
	name:         "saa.eval-cycle",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: time.Second * 1,
	usage:        "Cycle/ frequency the ScaleAlertAggregator evaluates the weights of the currently firing alerts.",
}

var saaEvalPeriodFactor = configEntry{
	name:         "saa.eval-peridod-factor",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: 10,
	usage:        "EvaluationPeriodFactor is used to calculate the evaluation period (evaluationPeriod = evaluationCycle * evaluationPeriodFactor)",
}

var saaCleanupCylce = configEntry{
	name:         "saa.cleanup-cycle",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: time.Second * 60,
	usage:        "Cycle/ frequency the ScaleAlertAggregator removes expired alerts.",
}

var bla = map[string]string{
	"n": "lsdfk",
}

var saaScaleAlerts = configEntry{
	name:         "saa.scale-alerts",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: "",
	usage:        "Cycle/ frequency the ScaleAlertAggregator removes expired alerts.",
}

var configEntries = []configEntry{
	configFile,
	port,
	dryRun,
	nomadServerAddress,
	jobName,
	jobMin,
	jobMax,
	capaDownScaleCoolDown,
	capaUpScaleCoolDown,
	loggingStructured,
	loggingUXTS,
	saaNoAlertDamping,
	saaUpThresh,
	saaDownThresh,
	saaEvalCylce,
	saaEvalPeriodFactor,
	saaCleanupCylce,
	saaScaleAlerts,
}

func (cfg *Config) fillCfgValues() {
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
	cfg.CapacityPlanner.DownScaleCooldownPeriod = cfg.viper.GetDuration(capaDownScaleCoolDown.name)
	cfg.CapacityPlanner.UpScaleCooldownPeriod = cfg.viper.GetDuration(capaUpScaleCoolDown.name)

	// Context: Logging
	cfg.Logging.Structured = cfg.viper.GetBool(loggingStructured.name)
	cfg.Logging.UxTimestamp = cfg.viper.GetBool(loggingUXTS.name)

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

	var alerts []Alert
	alertStr := cfg.viper.GetString(saaScaleAlerts.name)
	if len(alertStr) > 0 {
		alerts, _ = strToAlerts(alertStr)
	} else {

		alerts = make([]Alert, 0)

		sub := cfg.viper.Get(saaScaleAlerts.name)
		if sub != nil {
			log.Printf("SSS %v", sub)
		}
	}

	cfg.ScaleAlertAggregator.ScaleAlerts = alerts

	log.Printf("SCAAA %s", cfg.ScaleAlertAggregator.ScaleAlerts)

}

func strToAlerts(alertsStr string) ([]Alert, error) {

	log.Printf("STR %s", alertsStr)
	var alerts = make([]Alert, 0)

	return alerts, nil
}
