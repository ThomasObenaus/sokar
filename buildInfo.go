package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// BuildInfo contains information about the build of sokar.
type BuildInfo struct {
	Version   string `json:"version,omitempty"`
	BuildTime string `json:"build_time,omitempty"`
	Revision  string `json:"revision,omitempty"`
	Branch    string `json:"branch,omitempty"`
}

// Print prints the build information using the given print function
func (bi *BuildInfo) Print(printFun func(format string, a ...interface{}) (n int, err error)) {
	printFun("-----------------------------------------------------------------\n")
	printFun("BuildInfo\n")
	printFun("-----------------------------------------------------------------\n")
	printFun("\tVersion:\t%s\n", bi.Version)
	printFun("\tBuild-Time:\t%s\n", bi.BuildTime)
	printFun("\tRevision:\t%s on %s\n", bi.Revision, bi.Branch)
	printFun("-----------------------------------------------------------------\n")
}

// BuildInfo represents the build-info end-point of sokar
func (bi *BuildInfo) BuildInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := http.StatusOK

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	if err := enc.Encode(bi); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
