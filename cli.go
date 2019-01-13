package main

import (
	"github.com/namsral/flag"
)

type cliArgs struct {
	StructuredLogging          bool
	UseUnixTimestampForLogging bool
}

func parseArgs() cliArgs {
	var structuredLogging = flag.Bool("sl", false, "Enables/ disables structured logging (using json). Defaults to false.")
	var useUnixTimestampForLogging = flag.Bool("ul", false, "Enables/ disables the usage of unix timestamp in log messages. Defaults to false.")
	flag.Parse()

	return cliArgs{
		StructuredLogging:          *structuredLogging,
		UseUnixTimestampForLogging: *useUnixTimestampForLogging,
	}
}
