package config

import "time"

// Config is a structure containing the configuration for sokar
type Config struct {
	Port                 int                  `yaml:"port"`
	Nomad                Nomad                `yaml:"nomad"`
	Logging              Logging              `yaml:"logging,omitempty"`
	Job                  Job                  `yaml:"job"`
	ScaleAlertAggregator ScaleAlertAggregator `yaml:"scale_alert_aggregator"`
}

// Nomad represents the configuration for the scaling target nomad
type Nomad struct {
	ServerAddr string `yaml:"srv_addr"`
}

// Job represents the definition for the job that should be scaled.
type Job struct {
	Name     string `yaml:"name"`
	MinCount uint   `yaml:"min"`
	MaxCount uint   `yaml:"max"`
}

// ScaleAlertAggregator is the configuration part for the ScaleAlertAggregator
type ScaleAlertAggregator struct {
	NoAlertScaleDamping    float32       `yaml:"no_alert_damping,omitempty"`
	UpScaleThreshold       float32       `yaml:"up_thresh,omitempty"`
	DownScaleThreshold     float32       `yaml:"down_thresh,omitempty"`
	ScaleAlerts            []Alert       `yaml:"scale_alerts,omitempty"`
	EvaluationCycle        time.Duration `yaml:"eval_cycle,omitempty"`
	EvaluationPeriodFactor uint          `yaml:"eval_period_factor,omitempty"`
	CleanupCycle           time.Duration `yaml:"cleanup_cycle,omitempty"`
}

// Alert represents an alert defined by its name and weight
type Alert struct {
	Name        string  `yaml:"name"`
	Weight      float32 `yaml:"weight"`
	Description string  `yaml:"description,omitempty"`
}

// Logging is used for logging configuration
type Logging struct {
	Structured  bool `yaml:"structured,omitempty"`
	UxTimestamp bool `yaml:"unix_ts,omitempty"`
}
