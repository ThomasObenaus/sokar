package sokar

// BuildInfo contains information about the build of sokar.
type BuildInfo struct {
	Version   string
	BuildTime string
	Revision  string
	Branch    string
}

// Print prints the buildinformation using the given print function
func (bi *BuildInfo) Print(printFun func(format string, a ...interface{}) (n int, err error)) {
	printFun("-----------------------------------------------------------------\n")
	printFun("BuildInfo\n")
	printFun("-----------------------------------------------------------------\n")
	printFun("\tVersion:\t%s\n", bi.Version)
	printFun("\tBuild-Time:\t%s\n", bi.BuildTime)
	printFun("\tRevision:\t%s on %s\n", bi.Revision, bi.Branch)
	printFun("-----------------------------------------------------------------\n")
}

//// BuildInfo represents the build-info end-point of sokar
//func (sk *Sokar) BuildInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	code := http.StatusOK
//	sk.logger.Info().Str("health", http.StatusText(code)).Msg("BuildInfo Check called.")
//
//	w.WriteHeader(code)
//	io.WriteString(w, "Sokar is Healthy")
//}
