package config

type Config struct {
	ScaleAlertWeightMap ScaleAlertWeightMap `yaml:"scale_alert_weight_map"`
}

type ScaleAlertWeightMap map[string]float32
