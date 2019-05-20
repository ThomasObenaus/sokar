package serviceTest

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sokarAddr string

func setup() {
	fmt.Println("Setup")
}

func shutdown() {
	fmt.Println("Shutdown")
}

func TestMain(m *testing.M) {

	sokarAddrPtr := flag.String("sokar-address", "", "Address of sokar")
	flag.Parse()

	if sokarAddrPtr != nil {
		sokarAddr = *sokarAddrPtr
	}

	fmt.Printf("Sokar Address: %s\n", sokarAddr)

	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func Test_AlertmanagerRequest(t *testing.T) {
	// Invalid request from Alertmanager
	am := newAlertManager(sokarAddr, time.Second*2)
	code, err := am.sendAlertmanagerRequest("INVALID")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, code)
}

func Test_Alert(t *testing.T) {
	//gock.New("http://localhost:12000").

	nm := &nomadMock{}

	http.Handle("/", nm)

	go log.Fatal(http.ListenAndServe(":12000", nil))

	am := newAlertManager(sokarAddr, time.Second*2)

	alerts := make([]alert, 0)
	alerts = append(alerts, alert{
		Labels: map[string]string{"alertname": "Alert A"},
	})

	request, err := requestToStr(buildAlertRequest(alerts))
	require.NoError(t, err)

	code, err := am.sendAlertmanagerRequest(request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, code)
}
