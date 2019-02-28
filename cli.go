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
		flag.PrintDefaults()
	}

	return success
}

func parseArgs() cliArgs {
	var nomadServerAddr = flag.String(pNomadServerAddress, "", "Specifies the address of the nomad server.")
	var cfgFile = flag.String(pCfgFile, "", "Specifies the full path and name of the configuration file for sokar.")
	flag.Parse()

	return cliArgs{
		NomadServerAddr: *nomadServerAddr,
		CfgFile:         *cfgFile,
	}
}
