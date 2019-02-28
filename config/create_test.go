package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = `
scale_alert_weight_map:
  AlertA: 1.0
  AlertB: 100
  AlertC: -1.0
  AlertD: 0.009
`

func Test_configFromYAML(t *testing.T) {

	reader := strings.NewReader(data)

	config, err := NewConfigFromYAML(reader)
	assert.NoError(t, err)
	assert.Len(t, config.ScaleAlertWeightMap, 4)

}
