package m5cli

import (
	// "flag"
	"fmt"
	S "strings"

	flag "github.com/spf13/pflag"
)

// CommandLine is the default set of command-line flags, parsed from os.Args.
// var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

func myFlagFunc(p *flag.Flag) {
	fmt.Printf("FLAG %v \n", *p)
	p.Hidden = true
}

// DisableAllFlags operates on flag.CommandLine
func DisableAllFlags() {
	flag.CommandLine.VisitAll(disableFlag)
}
func disableFlag(p *flag.Flag) {
	p.Hidden = true
}

// EnableAllFlags operates on flag.CommandLine
func EnableAllFlags() {
	flag.CommandLine.VisitAll(enableFlag)
}
func enableFlag(p *flag.Flag) {
	p.Hidden = false
}

var flagsToEnable, flagsToDisable string

// DisableFlags operates on flag.CommandLine
func DisableFlags(s string) {
	flagsToDisable = s
	flag.CommandLine.VisitAll(maybeDisableFlag)
}
func maybeDisableFlag(p *flag.Flag) {
	var thisFlag = p.Shorthand
	if S.Contains(flagsToDisable, thisFlag) {
		p.Hidden = true
	}
}

// EnableFlags operates on flag.CommandLine
func EnableFlags(s string) {
	flagsToEnable = s
	flag.CommandLine.VisitAll(maybeEnableFlag)
}
func maybeEnableFlag(p *flag.Flag) {
	var thisFlag = p.Shorthand
	if S.Contains(flagsToEnable, thisFlag) {
		p.Hidden = false
	}
}
