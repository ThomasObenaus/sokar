package logging

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	cfg := Config{}
	lf := cfg.New()
	assert.NotNil(t, lf)
	assert.NotEmpty(t, zerolog.TimeFieldFormat)

	cfg = Config{
		UseUnixTimestampForLogging: true,
	}
	lf = cfg.New()
	assert.NotNil(t, lf)
	assert.Empty(t, zerolog.TimeFieldFormat)
}

func TestNewNamedLogger(t *testing.T) {

	loggerFactory := loggerFactoryImpl{
		UseStructuredLogging: true,
	}

	logger := loggerFactory.NewNamedLogger("MyTestLogger")
	strout := strings.Builder{}
	loggerDup := logger.Output(&strout)
	loggerDup.Info().Msg("HWLD")
	assert.Contains(t, strout.String(), "MyTestLogger")
	strout.Reset()

	loggerFactory = loggerFactoryImpl{
		UseStructuredLogging: false,
	}

	logger = loggerFactory.NewNamedLogger("MyTestLogger2")
	loggerDup = logger.Output(&strout)
	loggerDup.Info().Msg("HWLD")
	assert.Contains(t, strout.String(), "MyTestLogger2")
}

func Example() {
	cfg := Config{UseStructuredLogging: true}

	// create the factory
	loggingFactory := cfg.New()

	// create new named logger
	logger := loggingFactory.NewNamedLogger("MyLogger")
	logger.Info().Msg("Hello World")

}
