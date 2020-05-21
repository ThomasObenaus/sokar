package config

import (
	"time"

	cfg "github.com/ThomasObenaus/go-base/config"
)

// ###################### Context: main ####################################################
var configFile = cfg.NewEntry("config-file", "Specifies the full path and name of the configuration file for sokar.", cfg.Default(""), cfg.ShortName("f"))
var dryRun = cfg.NewEntry("dry-run", "If true, then sokar won't execute the planned scaling action. Only scaling\n"+
	"actions triggered via ScaleBy end-point will be executed.", cfg.Default(false))

var port = cfg.NewEntry("port", "Port where sokar is listening.", cfg.Default(11000))

// ###################### Context: scaler ####################################################

var scaMode = cfg.NewEntry("sca.mode", "Scaling target mode is either job based, instance-based or data-center\n"+
	"(worker/ instance) based scaling. In data-center (dc) mode the nomad workers\n"+
	"will be scaled. In job mode the number of allocations for this job will be adjusted.", cfg.Default("nomad-job"))

var scaWatcherInterval = cfg.NewEntry("sca.watcher-interval", "The interval the Scaler will check if the scalingObject count still matches\n"+
	"the desired state.", cfg.Default("5s"))

// ###################### Context: scaler AWS EC2 ############################################
var scaAWSEC2Profile = cfg.NewEntry("sca.aws-ec2.profile", "This parameter represents the name of the aws profile that shall be used to\n"+
	"access the resources to scale the data-center. This parameter is optional. If it\n"+
	"is empty the instance where sokar runs on has to have enough permissions to access\n"+
	"the resources (ASG) for scaling. In this case the AWSRegion parameter has to be\n"+
	"specified as well.", cfg.Default(""))

var scaAWSEC2Region = cfg.NewEntry("sca.aws-ec2.region", "This is an optional parameter and is regarded only if the parameter\n"+
	"AWSProfile is empty. The AWSRegion has to specify the region in which the\n"+
	"data-center to be scaled resides in.", cfg.Default("eu-central-1"))

var scaAWSEC2ASGTagKey = cfg.NewEntry("sca.aws-ec2.asg-tag-key", "This parameter specifies which tag on an AWS AutoScalingGroup shall be used\n"+
	"to find the ASG that should be automatically scaled.", cfg.Default("scale-object"))

// ###################### Context: scaler Nomad ###############################################
var scaNomadDataCenterAWSProfile = cfg.NewEntry("sca.nomad.dc-aws.profile", "This parameter represents the name of the aws profile that shall be used to\n"+
	"access the resources to scale the data-center. This parameter is optional. If it\n"+
	"is empty the instance where sokar runs on has to have enough permissions to access\n"+
	"the resources (ASG) for scaling. In this case the AWSRegion parameter has to be\n"+
	"specified as well.", cfg.Default(""))

var scaNomadDataCenterAWSRegion = cfg.NewEntry("sca.nomad.dc-aws.region", "This is an optional parameter and is regarded only if the parameter\n"+
	"AWSProfile is empty. The AWSRegion has to specify the region in which the data-center\n"+
	"to be scaled resides in.", cfg.Default("eu-central-1"))

var scaNomadDataCenterAWSInstanceTerminationTimeout = cfg.NewEntry("sca.nomad.dc-aws.instance-termination-timeout", "The maximum time the instance termination will be monitored before assuming\n"+
	"that this action (instance termination due to downscale) failed.", cfg.Default(time.Minute*10))

var scaNomadModeServerAddress = cfg.NewEntry("sca.nomad.server-address", "Specifies the address of the nomad server.", cfg.Default(""))

// ###################### Context: scale-object ####################################################
var scaleObjectName = cfg.NewEntry("scale-object.name", "The name of the object to be scaled.", cfg.Default(""))
var scaleObjectMin = cfg.NewEntry("scale-object.min", "The minimum count of the object to be scaled.", cfg.Default(1))
var scaleObjectMax = cfg.NewEntry("scale-object.max", "The maximum count of the object to be scaled.", cfg.Default(10))

// ###################### Context: CapacityPlanner#########################################
var capScaleSchedule = cfg.NewEntry("cap.scale-schedule", "Specifies time ranges within which it is ensured that the ScaleObject is scaled\n"+
	"to at least min and not more than max. The min/ max values specified in this\n"+
	"schedule have lower priority than the --scale-object.min/ --scale-object.max.\n"+
	"This means the sokar will ensure that the --scale-object.min/ --scale-object.max\n"+
	"are not violated no matter what is specified in the schedule.", cfg.Default(""))

var capDownScaleCoolDown = cfg.NewEntry("cap.down-scale-cooldown", "The time sokar waits between downscaling actions at min.", cfg.Default(time.Second*20))
var capUpScaleCoolDown = cfg.NewEntry("cap.up-scale-cooldown", "The time sokar waits between upscaling actions at min.", cfg.Default(time.Second*10))
var capConstantModeEnable = cfg.NewEntry("cap.constant-mode.enable", "Enable/ disable the constant mode of the CapacityPlanner. Only one of the\n"+
	"modes can be enabled at the same time.", cfg.Default(true))

var capConstantModeOffset = cfg.NewEntry("cap.constant-mode.offset", "The constant offset value that should be used to increment/ decrement the\n"+
	"count of the scale-object. Only values > 0 are valid.", cfg.Default(uint(1)))

var capLinearModeEnable = cfg.NewEntry("cap.linear-mode.enable", "Enable/ disable the linear mode of the CapacityPlanner. Only one of the modes\n"+
	"can be enabled at the same time.", cfg.Default(false))

var capLinearModeScaleFactorWeight = cfg.NewEntry("cap.linear-mode.scale-factor-weight", "This weight is used to adjust the impact of the scaleFactor during capacity\n"+
	"planning in linear mode.", cfg.Default(0.5))

// ###################### Context: Logging ################################################
var loggingStructured = cfg.NewEntry("logging.structured", "Use structured logging or not.", cfg.Default(false))
var loggingUXTS = cfg.NewEntry("logging.unix-ts", "Use Unix-Timestamp representation for log entries.", cfg.Default(false))
var loggingNoColor = cfg.NewEntry("logging.no-color", "If true colors in log out-put will be disabled.", cfg.Default(false))

// ###################### Context: ScaleAlertAggregator ###################################
var saaAlertExpirationTime = cfg.NewEntry("saa.alert-expiration-time", "Defines after which time an alert will be pruned if he did not get updated\n"+
	"again by the ScaleAlertEmitter, assuming that the alert is not relevant any more.", cfg.Default(time.Minute*10))

var saaNoAlertDamping = cfg.NewEntry("saa.no-alert-damping", "Damping used in case there are currently no alerts firing\n"+
	"(neither down- nor upscaling).", cfg.Default(1.0))

var saaUpThresh = cfg.NewEntry("saa.up-thresh", "Threshold for a upscaling event.", cfg.Default(10.0))
var saaDownThresh = cfg.NewEntry("saa.down-thresh", "Threshold for a downscaling event.", cfg.Default(-10.0))

var saaEvalCylce = cfg.NewEntry("saa.eval-cycle", "Cycle/ frequency the ScaleAlertAggregator evaluates the weights of the\n"+
	"currently firing alerts.", cfg.Default(time.Second*1))

var saaEvalPeriodFactor = cfg.NewEntry("saa.eval-period-factor", "EvaluationPeriodFactor is used to calculate the evaluation period\n"+
	"(evaluationPeriod = evaluationCycle * evaluationPeriodFactor)", cfg.Default(10))

var saaCleanupCylce = cfg.NewEntry("saa.cleanup-cycle", "Cycle/ frequency the ScaleAlertAggregator removes expired alerts.", cfg.Default(time.Second*60))
var saaScaleAlerts = cfg.NewEntry("saa.scale-alerts", "The alerts that should be used for scaling (up/down) the scale-object.", cfg.Default(""))

var configEntries = []cfg.Entry{
	configFile,
	port,
	dryRun,
	scaleObjectName,
	scaleObjectMin,
	scaleObjectMax,
	capScaleSchedule,
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
	scaMode,
	scaWatcherInterval,
	scaAWSEC2Profile,
	scaAWSEC2Region,
	scaAWSEC2ASGTagKey,
	scaNomadDataCenterAWSProfile,
	scaNomadDataCenterAWSRegion,
	scaNomadDataCenterAWSInstanceTerminationTimeout,
	scaNomadModeServerAddress,
	capConstantModeEnable,
	capConstantModeOffset,
	capLinearModeEnable,
	capLinearModeScaleFactorWeight,
}
