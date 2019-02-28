package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davecgh/go-spew/spew"
)

var data = `
alerts:
  - name: "AlertA"
    weight: 1.5
  - name: "AlertB"
    weight: -1.5
    description: "Down alert"
`

func Test_configFromYAML(t *testing.T) {

	reader := strings.NewReader(data)

	config, err := NewConfigFromYAML(reader)
	assert.NoError(t, err)
	assert.Len(t, config.Alerts, 2)
	assert.Equal(t, "AlertA", config.Alerts[0].Name)
	assert.Equal(t, float32(1.5), config.Alerts[0].Weight)
	assert.Equal(t, "", config.Alerts[0].Description)

	assert.Equal(t, "AlertB", config.Alerts[1].Name)
	assert.Equal(t, float32(-1.5), config.Alerts[1].Weight)
	assert.Equal(t, "Down alert", config.Alerts[1].Description)

	spew.Dump(config)

}
