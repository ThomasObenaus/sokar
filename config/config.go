package config

// Config is a structure containing the configuration for sokar
type Config struct {
	Alerts []Alert `yaml:"alerts"`

	//ScaleAlertWeightMap ScaleAlertWeightMap `yaml:"scale_alert_weight_map"`
}

type Alert struct {
	Name        string  `yaml:"name"`
	Weight      float32 `yaml:"weight"`
	Description string  `yaml:"description,omitempty"`
}

// ScaleAlertWeightMap is a map that provides a mapping of an alert name to
// its weight.
type ScaleAlertWeightMap map[string]float32
