package main

import (
	"fmt"
	"os"

	"github.com/namsral/flag"
)

const pNomadServerAddress = "nomad-server-address"
const pCfgFile = "config-file"
const pDryMode = "dry-run"

type cliArgs struct {
	NomadServerAddr string
	CfgFile         string
	DryRunMode      bool

	flagSet *flag.FlagSet
}

func (ca *cliArgs) printDefaults() {
	fmt.Println()
	fmt.Printf("Usage of %s\n", os.Args[0])
	ca.flagSet.PrintDefaults()
}

func (ca *cliArgs) validateArgs() bool {
	success := true

	if len(ca.CfgFile) == 0 {
		fmt.Printf("Parameter '-%s' is missing\n", pCfgFile)
		success = false
	}
	return success
}

func parseArgs(args []string) (cliArgs, error) {
	if args == nil || len(args) == 0 {
		return cliArgs{}, fmt.Errorf("Args are missing")
	}

	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)

	var nomadServerAddr = flagSet.String(pNomadServerAddress, "", "Specifies the address of the nomad server.")
	var cfgFile = flagSet.String(pCfgFile, "", "Specifies the full path and name of the configuration file for sokar.")
	var dryRun = flagSet.Bool(pDryMode, false, "If true, then sokar won't execute the planned scaling action. Only scaling actions triggered via ScaleBy end-point will be executed.")

	err := flagSet.Parse(args[1:])
	return cliArgs{
		NomadServerAddr: *nomadServerAddr,
		CfgFile:         *cfgFile,
		DryRunMode:      *dryRun,
		flagSet:         flagSet,
	}, err
}
