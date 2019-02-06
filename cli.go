package main

import (
	"fmt"
	"os"

	"github.com/namsral/flag"
)

const pNomadServerAddress = "nomad-server-address"
const pJobName = "job-name"
const pMinCount = "job-min-count"
const pMaxCount = "job-max-count"
const pScaleBy = "scale-by"

type cliArgs struct {
	StructuredLogging          bool
	UseUnixTimestampForLogging bool
	NomadServerAddr            string
	JobName                    string
	JobMinCount                uint
	JobMaxCount                uint
	ScaleBy                    int
}

func (ca *cliArgs) validateArgs() bool {
	success := true

	if len(ca.NomadServerAddr) == 0 {
		fmt.Printf("Parameter '-%s' is missing\n", pNomadServerAddress)
		success = false
	}

	if len(ca.JobName) == 0 {
		fmt.Printf("Parameter '-%s' is missing\n", pJobName)
		success = false
	}

	if !success {
		fmt.Println()
		fmt.Printf("Usage of %s\n", os.Args[0])
		flag.PrintDefaults()
	}

	return success
}

func parseArgs() cliArgs {
	var structuredLogging = flag.Bool("logging-structured", false, "Enables/ disables structured logging (using json). Defaults to false.")
	var useUnixTimestampForLogging = flag.Bool("logging-ux-ts", false, "Enables/ disables the usage of unix timestamp in log messages. Defaults to false.")
	var nomadServerAddr = flag.String(pNomadServerAddress, "", "Specifies the address of the nomad server.")
	var jobName = flag.String(pJobName, "", "The name of the job to be scaled.")
	var minCount = flag.Uint(pMinCount, 1, "Specifies the minimum number of instances this job shall have (default 1).")
	var maxCount = flag.Uint(pMaxCount, 2, "Specifies the minimum number of instances this job shall have (default 2).")
	var scaleBy = flag.Int(pScaleBy, 0, "Specifies the amount the job shall be scaled. A positive number means scale up and a negative means scale down by the specified amount.")
	flag.Parse()

	return cliArgs{
		StructuredLogging:          *structuredLogging,
		UseUnixTimestampForLogging: *useUnixTimestampForLogging,
		NomadServerAddr:            *nomadServerAddr,
		JobName:                    *jobName,
		JobMinCount:                *minCount,
		JobMaxCount:                *maxCount,
		ScaleBy:                    *scaleBy,
	}
}
