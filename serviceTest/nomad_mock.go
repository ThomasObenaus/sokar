package serviceTest

import (
	"net/http"
)

type nomadMock struct {
	err error
}

func (sm *nomadMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}
