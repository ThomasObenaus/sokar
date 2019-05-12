package config

import (
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
var capDownScaleCoolDown = configEntry{
	name:         "cap.down-scale-cooldown",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: time.Second * 20,
	usage:        "The time sokar waits between downscaling actions at min.",
}

var capUpScaleCoolDown = configEntry{
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

var loggingNoColor = configEntry{
	name:         "logging.no-color",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: false,
	usage:        "If true colors in log out-put will be disabled.",
}

// ###################### Context: ScaleAlertAggregator ###################################
var saaAlertExpirationTime = configEntry{
	name:         "saa.alert-expiration-time",
	bindEnv:      true,
	bindFlag:     true,
	defaultValue: time.Minute * 10,
	usage:        "Defines after which time an alert will be pruned if he did not get updated again by the ScaleAlertEmitter, assuming that the alert is not relevant any more.",
}

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
	name:         "saa.eval-period-factor",
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
	capDownScaleCoolDown,
	capUpScaleCoolDown,
	loggingStructured,
	loggingUXTS,
	loggingNoColor,
	saaNoAlertDamping,
	saaUpThresh,
	saaDownThresh,
	saaEvalCylce,
	saaEvalPeriodFactor,
	saaCleanupCylce,
	saaScaleAlerts,
	saaAlertExpirationTime,
}
