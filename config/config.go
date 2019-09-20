package config

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	Port                 int                  `json:"port,omitempty"`
	Scaler               Scaler               `json:"scaler,omitempty"`
	DryRunMode           bool                 `json:"dry_run_mode,omitempty"`
	Logging              Logging              `json:"logging,omitempty"`
	ScaleObject          ScaleObject          `json:"scale_object,omitempty"`
	ScaleAlertAggregator ScaleAlertAggregator `json:"scale_alert_aggregator,omitempty"`
	CapacityPlanner      CapacityPlanner      `json:"capacity_planner,omitempty"`

	configEntries []configEntry

	pFlagSet *pflag.FlagSet
	viper    *viper.Viper
}

// Scaler represents the config for the scaler/ ScalingTarget
type Scaler struct {
	Mode            ScalerMode    `json:"mode,omitempty"`
	Nomad           SCANomad      `json:"nomad,omitempty"`
	AwsEc2          SCAAwsEc2     `json:"aws_ec2,omitempty"`
	WatcherInterval time.Duration `json:"watcher_interval,omitempty"`
}

// SCAAwsEc2 represents the parameters for a AWS EC2 based scaler.
type SCAAwsEc2 struct {
	Profile   string `json:"profile,omitempty"`
	Region    string `json:"region,omitempty"`
	ASGTagKey string `json:"asg_tag_key,omitempty"`
}

// SCANomad represents the parameters for a nomad based scaler (job or data-center).
type SCANomad struct {
	ServerAddr    string                `json:"server_addr,omitempty"`
	DataCenterAWS SCANomadDataCenterAWS `json:"datacenter_aws,omitempty"`
}

// SCANomadDataCenterAWS represents the parameters needed for the nomad based scaler for mode data-center running on AWS.
type SCANomadDataCenterAWS struct {
	Profile   string `json:"profile,omitempty"`
	Region    string `json:"region,omitempty"`
	ASGTagKey string `json:"asg_tag_key,omitempty"`
}

// ScaleObject represents the definition for the object that should be scaled.
type ScaleObject struct {
	Name     string `json:"name,omitempty"`
	MinCount uint   `json:"min_count,omitempty"`
	MaxCount uint   `json:"max_count,omitempty"`
}

// ScaleAlertAggregator is the configuration part for the ScaleAlertAggregator
type ScaleAlertAggregator struct {
	NoAlertScaleDamping    float32       `json:"no_alert_scale_damping,omitempty"`
	UpScaleThreshold       float32       `json:"up_scale_threshold,omitempty"`
	DownScaleThreshold     float32       `json:"down_scale_threshold,omitempty"`
	ScaleAlerts            []Alert       `json:"scale_alerts,omitempty"`
	EvaluationCycle        time.Duration `json:"evaluation_cycle,omitempty"`
	EvaluationPeriodFactor uint          `json:"evaluation_period_factor,omitempty"`
	CleanupCycle           time.Duration `json:"cleanup_cycle,omitempty"`
	AlertExpirationTime    time.Duration `json:"alert_expiration_time,omitempty"`
}

// Alert represents an alert defined by its name and weight
type Alert struct {
	Name        string  `json:"name,omitempty"`
	Weight      float32 `json:"weight,omitempty"`
	Description string  `json:"description,omitempty"`
}

// Logging is used for logging configuration
type Logging struct {
	Structured         bool `json:"structured,omitempty"`
	UxTimestamp        bool `json:"ux_timestamp,omitempty"`
	NoColoredLogOutput bool `json:"no_colored_log_output,omitempty"`
}

// CapacityPlanner is used for the configuration of the CapacityPlanner
type CapacityPlanner struct {
	DownScaleCooldownPeriod time.Duration `json:"down_scale_cooldown_period,omitempty"`
	UpScaleCooldownPeriod   time.Duration `json:"up_scale_cooldown_period,omitempty"`
	ConstantMode            CAPConstMode  `json:"constant_mode,omitempty"`
	LinearMode              CAPLinearMode `json:"linear_mode,omitempty"`
}

// CAPLinearMode configuration parameters needed for linear mode of the CapacityPlanner
type CAPLinearMode struct {
	Enable            bool    `json:"enable,omitempty"`
	ScaleFactorWeight float64 `json:"scale_factor_weight,omitempty"`
}

// CAPConstMode configuration parameters needed for constant mode of the CapacityPlanner
type CAPConstMode struct {
	Enable bool `json:"enable,omitempty"`
	Offset uint `json:"offset,omitempty"`
}

// NewDefaultConfig returns a default configuration without any alerts (mappings)
// or server configuration defined.
func NewDefaultConfig() Config {

	cfg := Config{
		Port:        11000,
		DryRunMode:  false,
		Logging:     Logging{Structured: false, UxTimestamp: false},
		ScaleObject: ScaleObject{},
		Scaler: Scaler{
			Mode:            ScalerModeNomadJob,
			Nomad:           SCANomad{},
			WatcherInterval: time.Second * 5,
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

	cfg.pFlagSet = pflag.NewFlagSet("sokar-config", pflag.ContinueOnError)
	cfg.viper = viper.New()
	cfg.configEntries = configEntries

	return cfg
}
