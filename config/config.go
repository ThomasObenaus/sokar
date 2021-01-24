package config

import (
	"fmt"
	"reflect"
	"time"

	cfglib "github.com/ThomasObenaus/go-base/config"
	cfglibIf "github.com/ThomasObenaus/go-base/config/interfaces"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// ScalerMode represents the mode the Scaler can be set to
type ScalerMode string

const (
	// ScalerModeNomadJob that the number of allocations of a job will be scaled
	ScalerModeNomadJob ScalerMode = "nomad-job"
	// ScalerModeNomadDataCenter that the number of instances/ workers of a data-center will be scaled
	ScalerModeNomadDataCenter ScalerMode = "nomad-dc"
	// ScalerModeAwsEc2 that the number of instances/ workers of a AWS EC2 ASG will be scaled
	ScalerModeAwsEc2 ScalerMode = "aws-ec2"
)

// Config is a structure containing the configuration for sokar
type Config struct {
	Port                 int                  `json:"port" cfg:"{'name':'port','desc':'Port where sokar is listening.','default':11000}"`
	Scaler               Scaler               `json:"scaler" cfg:"{'name':'sca'}"`
	DryRunMode           bool                 `json:"dry_run_mode" cfg:"{'name':'dry-run','desc':'If true, then sokar wont execute the planned scaling action. Only scaling actions triggered via ScaleBy end-point will be executed.','default':false}"`
	Logging              Logging              `json:"logging" cfg:"{'name':'logging'}"`
	ScaleObject          ScaleObject          `json:"scale_object" cfg:"{'name':'scale-object'}"`
	ScaleAlertAggregator ScaleAlertAggregator `json:"scale_alert_aggregator" cfg:"{'name':'saa'}"`
	CapacityPlanner      CapacityPlanner      `json:"capacity_planner" cfg:"{'name':'cap'}"`
}

// Scaler represents the config for the scaler/ ScalingTarget
type Scaler struct {
	Mode            ScalerMode    `json:"mode" cfg:"{'name':'mode','desc':'Scaling target mode is either job based, instance-based or data-center (worker/ instance) based scaling. In data-center (dc) mode the nomad workers will be scaled. In job mode the number of allocations for this job will be adjusted.','default':'nomad-job'}"`
	Nomad           SCANomad      `json:"nomad" cfg:"{'name':'nomad'}"`
	AwsEc2          SCAAwsEc2     `json:"aws_ec2" cfg:"{'name':'aws-ec2'}"`
	WatcherInterval time.Duration `json:"watcher_interval" cfg:"{'name':'watcher-interval','desc':'The interval the Scaler will check if the scalingObject count still matches the desired state.','default':'5s'}"`
}

// SCAAwsEc2 represents the parameters for a AWS EC2 based scaler.
type SCAAwsEc2 struct {
	Profile   string `json:"profile" cfg:"{'name':'profile','desc':'This parameter represents the name of the aws profile that shall be used to access the resources to scale the data-center. This parameter is optional. If it is empty the instance where sokar runs on has to have enough permissions to access the resources (ASG) for scaling. In this case the AWSRegion parameter has to be specified as well.','default':''}"`
	Region    string `json:"region" cfg:"{'name':'region','desc':'This is an optional parameter and is regarded only if the parameter AWSProfile is empty. The AWSRegion has to specify the region in which the data-center to be scaled resides in.','default':'eu-central-1'}"`
	ASGTagKey string `json:"asg_tag_key" cfg:"{'name':'asg-tag-key','desc':'This parameter specifies which tag on an AWS AutoScalingGroup shall be used to find the ASG that should be automatically scaled.','default':'scale-object'}"`
}

// SCANomad represents the parameters for a nomad based scaler (job or data-center).
type SCANomad struct {
	ServerAddr    string                `json:"server_addr" cfg:"{'name':'server-address','desc':'Specifies the address of the nomad server.','default':'http://localhost:4646'}"`
	DataCenterAWS SCANomadDataCenterAWS `json:"datacenter_aws" cfg:"{'name':'dc-aws'}"`
}

// SCANomadDataCenterAWS represents the parameters needed for the nomad based scaler for mode data-center running on AWS.
type SCANomadDataCenterAWS struct {
	Profile                    string        `json:"profile" cfg:"{'name':'profile','desc':'This parameter represents the name of the aws profile that shall be used to access the resources to scale the data-center. This parameter is optional. If it is empty the instance where sokar runs on has to have enough permissions to access the resources (ASG) for scaling. In this case the AWSRegion parameter has to be specified as well.','default':''}"`
	Region                     string        `json:"region" cfg:"{'name':'region','desc':'This is an optional parameter and is regarded only if the parameter AWSProfile is empty. The AWSRegion has to specify the region in which the data-center to be scaled resides in.','default':'eu-central-1'}"`
	ASGTagKey                  string        `json:"asg_tag_key" cfg:"{'name':'asg-tag-key','desc':'This parameter specifies which tag on an AWS AutoScalingGroup shall be used to find the ASG that should be automatically scaled.','default':'scale-object'}"`
	InstanceTerminationTimeout time.Duration `json:"instance_termination_timeout" cfg:"{'name':'instance-termination-timeout','desc':'The intervall the Scaler will check if the scalingObject count still matches the desired state.','default':'5s'}"`
}

// ScaleObject represents the definition for the object that should be scaled.
type ScaleObject struct {
	Name     string `json:"name" cfg:"{'name':'name','desc':'The name of the object to be scaled.','default':''}"`
	MinCount uint   `json:"min_count" cfg:"{'name':'min','desc':'The minimum count of the object to be scaled.','default':1}"`
	MaxCount uint   `json:"max_count" cfg:"{'name':'max','desc':'The maximum count of the object to be scaled.','default':10}"`
}

// ScaleAlertAggregator is the configuration part for the ScaleAlertAggregator
type ScaleAlertAggregator struct {
	NoAlertScaleDamping    float32       `json:"no_alert_scale_damping" cfg:"{'name':'no-alert-damping','desc':'Damping used in case there are currently no alerts firing (neither down- nor upscaling).','default':1.0}"`
	UpScaleThreshold       float32       `json:"up_scale_threshold" cfg:"{'name':'up-thresh','desc':'Threshold for a upscaling event.','default':10.0}"`
	DownScaleThreshold     float32       `json:"down_scale_threshold" cfg:"{'name':'down-thresh','desc':'Threshold for a downscaling event.','default':-10.0}"`
	ScaleAlerts            []Alert       `json:"scale_alerts" cfg:"{'name':'scale-alerts','desc':'The alerts that should be used for scaling (up/down) the scale-object.','default':[]}"`
	EvaluationCycle        time.Duration `json:"evaluation_cycle" cfg:"{'name':'eval-cycle','desc':'Cycle/ frequency the ScaleAlertAggregator evaluates the weights of the currently firing alerts.','default':'1s'}"`
	EvaluationPeriodFactor uint          `json:"evaluation_period_factor" cfg:"{'name':'eval-period-factor','desc':'EvaluationPeriodFactor is used to calculate the evaluation period (evaluationPeriod = evaluationCycle * evaluationPeriodFactor).','default':10}"`
	CleanupCycle           time.Duration `json:"cleanup_cycle" cfg:"{'name':'cleanup-cycle','desc':'Cycle/ frequency the ScaleAlertAggregator removes expired alerts.','default':'60s'}"`
	AlertExpirationTime    time.Duration `json:"alert_expiration_time" cfg:"{'name':'alert-expiration-time','desc':'Defines after which time an alert will be pruned if he did not get updated again by the ScaleAlertEmitter, assuming that the alert is not relevant any more.','default':'10m'}"`
}

// Alert represents an alert defined by its name and weight
type Alert struct {
	Name        string  `json:"name" cfg:"{'name':'name','desc':'Name of the alert on prometheus alertmanager. This name is used to identify the alert that should be used for scaling.','default':''}"`
	Weight      float32 `json:"weight" cfg:"{'name':'weight','desc':'The weight of the alert, hence the impact it should have on the scaling.','default':1.0}"`
	Description string  `json:"description" cfg:"{'name':'description','default':''}"`
}

// Logging is used for logging configuration
type Logging struct {
	Structured         bool          `json:"structured" cfg:"{'name':'structured','desc':'Use structured logging or not.','default':false}"`
	UxTimestamp        bool          `json:"ux_timestamp" cfg:"{'name':'unix-ts','desc':'Use Unix-Timestamp representation for log entries.','default':false}"`
	NoColoredLogOutput bool          `json:"no_colored_log_output" cfg:"{'name':'no-color','desc':'If true colors in log out-put will be disabled.','default':false}"`
	Level              zerolog.Level `json:"level" cfg:"{'name':'level','desc':'The level that should be used for logs. Valid entries are debug, info, warn, error, fatal and off.','default':'info','mapfun':'strToLoglevel'}"`
}

// CapacityPlanner is used for the configuration of the CapacityPlanner
type CapacityPlanner struct {
	DownScaleCooldownPeriod time.Duration        `json:"down_scale_cooldown_period" cfg:"{'name':'down-scale-cooldown','desc':'The time sokar waits between downscaling actions at min.','default':'20s'}"`
	UpScaleCooldownPeriod   time.Duration        `json:"up_scale_cooldown_period" cfg:"{'name':'up-scale-cooldown','desc':'The time sokar waits between upscaling actions at min.','default':'10s'}"`
	ConstantMode            CAPConstMode         `json:"constant_mode" cfg:"{'name':'constant-mode'}"`
	LinearMode              CAPLinearMode        `json:"linear_mode" cfg:"{'name':'linear-mode'}"`
	ScaleSchedule           []ScaleScheduleEntry `json:"scaling_schedule" cfg:"{'name':'scale-schedule','desc':'Specifies time ranges within which it is ensured that the ScaleObject is scaled to at least min and not more than max. The min/ max values specified in this schedule have lower priority than the --scale-object.min/ --scale-object.max. This means the sokar will ensure that the --scale-object.min/ --scale-object.max are not violated no matter what is specified in the schedule.','mapfun':'strToScaleSchedule','default':''}"`
}

// CAPLinearMode configuration parameters needed for linear mode of the CapacityPlanner
type CAPLinearMode struct {
	Enable            bool    `json:"enable" cfg:"{'name':'enable','desc':'Enable/ disable the linear mode of the CapacityPlanner. Only one of the modes can be enabled at the same time.','default':false}"`
	ScaleFactorWeight float64 `json:"scale_factor_weight" cfg:"{'name':'scale-factor-weight','desc':'This weight is used to adjust the impact of the scaleFactor during capacity planning in linear mode.','default':0.5}"`
}

// CAPConstMode configuration parameters needed for constant mode of the CapacityPlanner
type CAPConstMode struct {
	Enable bool `json:"enable" cfg:"{'name':'enable','desc':'Enable/ disable the constant mode of the CapacityPlanner. Only one of the modes can be enabled at the same time.','default':true}"`
	Offset uint `json:"offset" cfg:"{'name':'offset','desc':'The constant offset value that should be used to increment/ decrement the count of the scale-object. Only values > 0 are valid.','default':1}"`
}

// NewDefaultConfig returns a default configuration without any alerts (mappings)
// or server configuration defined.
func NewDefaultConfig() Config {

	cfg := Config{
		Port:        11000,
		DryRunMode:  false,
		Logging:     Logging{Structured: false, UxTimestamp: false, Level: zerolog.InfoLevel},
		ScaleObject: ScaleObject{},
		Scaler: Scaler{
			Mode:            ScalerModeNomadJob,
			Nomad:           SCANomad{},
			WatcherInterval: time.Second * 5,
			AwsEc2:          SCAAwsEc2{ASGTagKey: "scale-object"},
		},
		ScaleAlertAggregator: ScaleAlertAggregator{
			EvaluationCycle:        time.Second * 1,
			EvaluationPeriodFactor: 10,
			CleanupCycle:           time.Second * 60,
			NoAlertScaleDamping:    1,
			UpScaleThreshold:       10,
			DownScaleThreshold:     -10,
			ScaleAlerts:            make([]Alert, 0),
			AlertExpirationTime:    time.Minute * 10,
		},
		CapacityPlanner: CapacityPlanner{
			DownScaleCooldownPeriod: time.Second * 80,
			UpScaleCooldownPeriod:   time.Second * 60,
			ConstantMode:            CAPConstMode{Enable: true, Offset: 1},
			LinearMode:              CAPLinearMode{Enable: false},
		},
	}

	return cfg
}

// New creates a new Config instance based on the given cli args
func New(args []string, serviceAbbreviation string) (Config, error) {

	config := Config{}

	provider, err := cfglib.NewConfigProvider(
		&config,
		serviceAbbreviation,
		serviceAbbreviation,
		cfglib.CustomConfigEntries(configEntries),
		cfglib.Logger(cfglibIf.InfoLogger),
	)
	if err != nil {
		return Config{}, err
	}

	if err := provider.RegisterMappingFunc("strToLoglevel", strToLoglevel); err != nil {
		return Config{}, err
	}
	if err := provider.RegisterMappingFunc("strToScaleSchedule", strToScaleSchedule); err != nil {
		return Config{}, err
	}

	err = provider.ReadConfig(args)
	if err != nil {
		fmt.Print(provider.Usage())
		fmt.Println()
		return Config{}, err
	}

	if err := config.fillCfgValues(provider); err != nil {
		return Config{}, err
	}

	return config, nil
}

func strToLoglevel(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error) {
	asString, ok := rawUntypedValue.(string)
	if !ok {
		return nil, fmt.Errorf("Expected type string but was %T", rawUntypedValue)
	}

	if asString == "off" {
		return zerolog.Disabled, nil
	}
	return zerolog.ParseLevel(asString)
}

func strToScaleSchedule(rawUntypedValue interface{}, targetType reflect.Type) (interface{}, error) {
	asString, ok := rawUntypedValue.(string)
	if !ok {
		return nil, fmt.Errorf("Expected type string but was %T", rawUntypedValue)
	}

	ssEntries, err := parseScalingScheduleEntries(asString)
	if err != nil {
		return nil, errors.Wrapf(err, "Parsing scale Schedule entries from '%s'", asString)
	}

	return ssEntries, nil
}
