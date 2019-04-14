package config

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config is a structure containing the configuration for sokar
type Config struct {
	Port                 int
	DryRunMode           bool
	Nomad                Nomad
	Logging              Logging
	Job                  Job
	ScaleAlertAggregator ScaleAlertAggregator
	CapacityPlanner      CapacityPlanner

	configEntries []configEntry

	pFlagSet *pflag.FlagSet
	viper    *viper.Viper
}

// Nomad represents the configuration for the scaling target nomad
type Nomad struct {
	ServerAddr string
}

// Job represents the definition for the job that should be scaled.
type Job struct {
	Name     string
	MinCount uint
	MaxCount uint
}

// ScaleAlertAggregator is the configuration part for the ScaleAlertAggregator
type ScaleAlertAggregator struct {
	NoAlertScaleDamping    float32
	UpScaleThreshold       float32
	DownScaleThreshold     float32
	ScaleAlerts            []Alert
	EvaluationCycle        time.Duration
	EvaluationPeriodFactor uint
	CleanupCycle           time.Duration
}

// Alert represents an alert defined by its name and weight
type Alert struct {
	Name        string
	Weight      float32
	Description string
}

// Logging is used for logging configuration
type Logging struct {
	Structured  bool
	UxTimestamp bool
}

// CapacityPlanner is used for the configuration of the CapacityPlanner
type CapacityPlanner struct {
	DownScaleCooldownPeriod time.Duration
	UpScaleCooldownPeriod   time.Duration
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
