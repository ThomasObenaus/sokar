package helper

import (
	"os"
	"sync"
)

var doObtainServerAddresses sync.Once = sync.Once{}
var obtainedSokarAddress string
var obtainedNomadAddress string

func ServerAddresses() (sokarAddr string, nomadAddr string) {

	doObtainServerAddresses.Do(func() {
		obtainedSokarAddress = os.Getenv("SOKAR_ADDR")
		if len(obtainedSokarAddress) == 0 {
			obtainedSokarAddress = "http://localhost:11000"
		}
		obtainedNomadAddress = os.Getenv("NOMAD_ADDR")
		if len(obtainedNomadAddress) == 0 {
			obtainedNomadAddress = "http://localhost:4646"
		}
	})

	return obtainedSokarAddress, obtainedNomadAddress
}
