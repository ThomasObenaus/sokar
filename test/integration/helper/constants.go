package helper

import (
	"flag"
	"sync"
)

var doObtainServerAddresses sync.Once = sync.Once{}
var obtainedSokarAddress string
var obtainedNomadAddress string

func ServerAddresses() (sokarAddr string, nomadAddr string) {

	doObtainServerAddresses.Do(func() {
		flag.StringVar(&obtainedSokarAddress, "sokar-addr", "http://localhost:11000", "Address of sokar (e.g. http://localhost:11000)")
		flag.StringVar(&obtainedNomadAddress, "nomad-addr", "http://localhost:4646", "Address of nomad (e.g. http://localhost:4646)")
	})

	return obtainedSokarAddress, obtainedNomadAddress
}
