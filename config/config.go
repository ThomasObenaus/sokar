package config

type Config struct {
	ScaleAlertWeightMap ScaleAlertWeightMap
}

type ScaleAlertWeightMap map[string]float32
