package config

import (
	"time"

	"github.com/rs/zerolog"

	cfg "github.com/ThomasObenaus/go-base/config"
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
	Port                 int                  `json:"port"`
	Scaler               Scaler               `json:"scaler"`
	DryRunMode           bool                 `json:"dry_run_mode"`
	Logging              Logging              `json:"logging"`
	ScaleObject          ScaleObject          `json:"scale_object"`
	ScaleAlertAggregator ScaleAlertAggregator `json:"scale_alert_aggregator"`
	CapacityPlanner      CapacityPlanner      `json:"capacity_planner"`
}

// Scaler represents the config for the scaler/ ScalingTarget
type Scaler struct {
	Mode            ScalerMode    `json:"mode"`
	Nomad           SCANomad      `json:"nomad"`
	AwsEc2          SCAAwsEc2     `json:"aws_ec2"`
	WatcherInterval time.Duration `json:"watcher_interval"`
}

// SCAAwsEc2 represents the parameters for a AWS EC2 based scaler.
type SCAAwsEc2 struct {
	Profile   string `json:"profile"`
	Region    string `json:"region"`
	ASGTagKey string `json:"asg_tag_key"`
}

// SCANomad represents the parameters for a nomad based scaler (job or data-center).
type SCANomad struct {
	ServerAddr    string                `json:"server_addr"`
	DataCenterAWS SCANomadDataCenterAWS `json:"datacenter_aws"`
}

// SCANomadDataCenterAWS represents the parameters needed for the nomad based scaler for mode data-center running on AWS.
type SCANomadDataCenterAWS struct {
	Profile                    string        `json:"profile"`
	Region                     string        `json:"region"`
	ASGTagKey                  string        `json:"asg_tag_key"`
	InstanceTerminationTimeout time.Duration `json:"instance_termination_timeout"`
}

// ScaleObject represents the definition for the object that should be scaled.
type ScaleObject struct {
	Name     string `json:"name"`
	MinCount uint   `json:"min_count"`
	MaxCount uint   `json:"max_count"`
}

// ScaleAlertAggregator is the configuration part for the ScaleAlertAggregator
type ScaleAlertAggregator struct {
	NoAlertScaleDamping    float32       `json:"no_alert_scale_damping"`
	UpScaleThreshold       float32       `json:"up_scale_threshold"`
	DownScaleThreshold     float32       `json:"down_scale_threshold"`
	ScaleAlerts            []Alert       `json:"scale_alerts"`
	EvaluationCycle        time.Duration `json:"evaluation_cycle"`
	EvaluationPeriodFactor uint          `json:"evaluation_period_factor"`
	CleanupCycle           time.Duration `json:"cleanup_cycle"`
	AlertExpirationTime    time.Duration `json:"alert_expiration_time"`
}

// Alert represents an alert defined by its name and weight
type Alert struct {
	Name        string  `json:"name"`
	Weight      float32 `json:"weight"`
	Description string  `json:"description"`
}

// Logging is used for logging configuration
type Logging struct {
	Structured         bool          `json:"structured"`
	UxTimestamp        bool          `json:"ux_timestamp"`
	NoColoredLogOutput bool          `json:"no_colored_log_output"`
	Level              zerolog.Level `json:"level"`
}

// CapacityPlanner is used for the configuration of the CapacityPlanner
type CapacityPlanner struct {
	DownScaleCooldownPeriod time.Duration        `json:"down_scale_cooldown_period"`
	UpScaleCooldownPeriod   time.Duration        `json:"up_scale_cooldown_period"`
	ConstantMode            CAPConstMode         `json:"constant_mode"`
	LinearMode              CAPLinearMode        `json:"linear_mode"`
	ScaleSchedule           []ScaleScheduleEntry `json:"scaling_schedule"`
}

// CAPLinearMode configuration parameters needed for linear mode of the CapacityPlanner
type CAPLinearMode struct {
	Enable            bool    `json:"enable"`
	ScaleFactorWeight float64 `json:"scale_factor_weight"`
}

// CAPConstMode configuration parameters needed for constant mode of the CapacityPlanner
type CAPConstMode struct {
	Enable bool `json:"enable"`
	Offset uint `json:"offset"`
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
			ScaleSchedule:           make([]ScaleScheduleEntry, 0),
		},
	}

	return cfg
}

// New creates a new Config instance based on the given cli args
func New(args []string, serviceAbbreviation string) (Config, error) {
	provider := cfg.NewProvider(configEntries, serviceAbbreviation, serviceAbbreviation)
	err := provider.ReadConfig(args)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	if err := config.fillCfgValues(provider); err != nil {
		return Config{}, err
	}

	return config, nil
}
