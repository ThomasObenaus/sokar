package logging

import (
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	lf := New(false, false, false)
	assert.NotNil(t, lf)
	assert.NotEmpty(t, zerolog.TimeFieldFormat)

	lf = New(false, true, false)
	assert.NotNil(t, lf)
	assert.Empty(t, zerolog.TimeFieldFormat)
}

func TestNewNamedLogger(t *testing.T) {

	loggerFactory := New(true, false, false)

	logger := loggerFactory.NewNamedLogger("MyTestLogger")
	strout := strings.Builder{}
	loggerDup := logger.Output(&strout)
	loggerDup.Info().Msg("HWLD")
	assert.Contains(t, strout.String(), "MyTestLogger")
	strout.Reset()

	loggerFactory = New(false, false, false)

	logger = loggerFactory.NewNamedLogger("MyTestLogger2")
	loggerDup = logger.Output(&strout)
	loggerDup.Info().Msg("HWLD")
	assert.Contains(t, strout.String(), "MyTestLogger2")
}

func ExampleNew() {
	// create the factory
	loggingFactory := New(true, false, false)

	// create new named logger
	logger := loggingFactory.NewNamedLogger("MyLogger")
	logger.Info().Msg("Hello World")

}
