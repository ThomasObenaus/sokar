package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	//zerolog.TimeFieldFormat = ""

	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("logger", "sokar.api").Logger()
	log.Logger = zerolog.New(os.Stdout).With().Str("logger", "sokar.api").Logger()

	log.Info().Float64("duration", 29.343).Str("region", "ED01").Msg("hello world")
}
