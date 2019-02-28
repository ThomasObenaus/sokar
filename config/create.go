package config

import (
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// NewConfigFromYAML reads in the configuration in yaml format
// using the provided io.Reader
func NewConfigFromYAML(reader io.Reader) (Config, error) {
	cfg := NewDefaultConfig()
	err := yaml.NewDecoder(reader).Decode(&cfg)
	return cfg, err
}

// NewDefaultConfig returns a default configuration without any alerts (mappings)
// or server configuration defined.
func NewDefaultConfig() Config {
	return Config{
		Logging{Structured: false, UxTimestamp: false},
		Job{},
		ScaleAlertAggregator{
			EvaluationCycle:        time.Second * 1,
			EvaluationPeriodFactor: 10,
			CleanupCycle:           time.Second * 60,
			NoAlertScaleDamping:    1,
			UpScaleThreshold:       10,
			DownScaleThreshold:     -10,
		},
	}
}

// NewConfigFromYAMLFile reads the configuration from a file
func NewConfigFromYAMLFile(file *os.File) (Config, error) {
	return NewConfigFromYAML(file)
}
