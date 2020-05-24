package main

import bInfo "github.com/ThomasObenaus/go-base/buildinfo"

// the following variables are set via go build-flags
// -X main.version=... -X main.buildTime=...
// -X main.revision=... -X main.branch=...
var version string
var buildTime string
var revision string
var branch string

var buildinfo = bInfo.BuildInfo{
	Version:   version,
	BuildTime: buildTime,
	Revision:  revision,
	Branch:    branch,
}
