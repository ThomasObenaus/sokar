package serviceTest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type sokarMock struct {
	amRequest request
	err       error
}

func (sm *sokarMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		sm.err = fmt.Errorf("Request object is nil")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	data, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		sm.err = fmt.Errorf("Error reading body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(data, &sm.amRequest); err != nil {
		sm.err = fmt.Errorf("Error parsing body: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
