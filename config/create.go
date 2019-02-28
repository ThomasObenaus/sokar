package config

import (
	"io"

	"gopkg.in/yaml.v2"
)

// NewConfigFromYAML reads in the configuration in yaml format
// using the provided io.Reader
func NewConfigFromYAML(reader io.Reader) (Config, error) {
	cfg := Config{}
	err := yaml.NewDecoder(reader).Decode(&cfg)
	return cfg, err
}
