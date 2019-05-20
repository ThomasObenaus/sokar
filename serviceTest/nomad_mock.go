package serviceTest

import (
	"fmt"
	"net/http"
)

type nomadMock struct {
	err error
}

func (sm *nomadMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Request %v\n", *r)

	w.WriteHeader(http.StatusOK)
}
