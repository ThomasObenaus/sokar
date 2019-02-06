package nomadConnector

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnector(t *testing.T) {

	cfg := Config{NomadServerAddress: "http://1.2.3.4"}
	connector, err := cfg.New()

	assert.NotNil(t, connector)
	assert.NoError(t, err)

	cfg = Config{}
	connector, err = cfg.New()

	assert.Nil(t, connector)
	assert.Error(t, err)
}

func ExampleConfig_New() {
	cfg := NewDefaultConfig("http://1.2.3.4")
	conn, err := cfg.New()

	if err != nil {
		log.Fatalf("Unable to create connector: %s.", err.Error())
	}

	// just to avoid the not used error
	_ = conn
}
