package main

import "fmt"

// BuildInfo contains information about the build of sokar.
type BuildInfo struct {
	Version   string
	BuildTime string
	Revision  string
	Branch    string
}

// ToStdOut print the buildinformation to stdout
func (bi *BuildInfo) ToStdOut() {
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("BuildInfo")
	fmt.Println("-----------------------------------------------------------------")
	fmt.Printf("\tVersion:\t%s\n", bi.Version)
	fmt.Printf("\tBuild-Time:\t%s\n", bi.BuildTime)
	fmt.Printf("\tRevision:\t%s on %s\n", bi.Revision, bi.Branch)
	fmt.Println("-----------------------------------------------------------------")
}
