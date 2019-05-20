package serviceTest

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	fmt.Println("Setup")
}

func shutdown() {
	fmt.Println("Shutdown")
}

func TestMain(m *testing.M) {

	sokarAddr := flag.String("sokar-address", "", "Address of sokar")
	flag.Parse()

	fmt.Printf("Sokar Address: %s\n", *sokarAddr)

	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestScale_JobDead(t *testing.T) {

	assert.Equal(t, "expected", "expected")
}

func TestScale_JobCheck(t *testing.T) {

	assert.Equal(t, "expected", "expected")
}
