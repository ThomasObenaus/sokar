package main

import (
	"fmt"
	"os"

	"github.com/namsral/flag"
)

const pNomadServerAddress = "nomad-server-address"
const pCfgFile = "config-file"

type cliArgs struct {
	NomadServerAddr string
	CfgFile         string

	flagSet *flag.FlagSet
}

func (ca *cliArgs) validateArgs() bool {
	success := true

	if len(ca.CfgFile) == 0 {
		fmt.Printf("Parameter '-%s' is missing\n", pCfgFile)
		success = false
	}

	if !success {
		fmt.Println()
		fmt.Printf("Usage of %s\n", os.Args[0])
		ca.flagSet.PrintDefaults()
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

	err := flagSet.Parse(args[1:])
	return cliArgs{
		NomadServerAddr: *nomadServerAddr,
		CfgFile:         *cfgFile,
		flagSet:         flagSet,
	}, err
}
