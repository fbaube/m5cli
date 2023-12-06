// https://pkg.go.dev/runtime/debug#BuildInfo

package main

import "runtime/debug"

func GetBuildInfo() *debug.BuildInfo {
	// func ReadBuildInfo() (info *BuildInfo, ok bool)
	pBI, ok := debug.ReadBuildInfo()
	if !ok {
		println("COULD NOT READ debug.BuildInfo")
	}
	return pBI
}

/*
type BuildInfo struct {
	GoVersion string         // Version of Go that produced this binary.
	Path      string         // The main package path
	Main      Module         // The module containing the main package
	Deps      []*Module      // Module dependencies
	Settings  []BuildSetting // Other information about the build.
}
*/
