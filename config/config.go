package config

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config is a structure containing the configuration for sokar
type Config struct {
	Port                 int                  `json:"port,omitempty"`
	DummyScalingTarget   bool                 `json:"dummy_scaling_target,omitempty"`
	DryRunMode           bool                 `json:"dry_run_mode,omitempty"`
	Nomad                Nomad                `json:"nomad,omitempty"`
	Logging              Logging              `json:"logging,omitempty"`
	Job                  Job                  `json:"job,omitempty"`
	ScaleAlertAggregator ScaleAlertAggregator `json:"scale_alert_aggregator,omitempty"`
	CapacityPlanner      CapacityPlanner      `json:"capacity_planner,omitempty"`

	configEntries []configEntry

	pFlagSet *pflag.FlagSet
	viper    *viper.Viper
}

// Nomad represents the configuration for the scaling target nomad
type Nomad struct {
	ServerAddr string `json:"server_addr,omitempty"`
}

// Job represents the definition for the job that should be scaled.
type Job struct {
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
}

// NewDefaultConfig returns a default configuration without any alerts (mappings)
// or server configuration defined.
func NewDefaultConfig() Config {

	cfg := Config{
		Port:       11000,
		DryRunMode: false,
		Nomad:      Nomad{},
		Logging:    Logging{Structured: false, UxTimestamp: false},
		Job:        Job{},
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
		},
	}

	cfg.pFlagSet = pflag.NewFlagSet("sokar-config", pflag.ContinueOnError)
	cfg.viper = viper.New()
	cfg.configEntries = configEntries

	return cfg
}
