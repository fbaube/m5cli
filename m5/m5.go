package main

// NOTE: Default log levels for the three phases of
// processing are defined at the top of file cli/exec.go

import (
	// "database/sql"
	"fmt"

	_ "github.com/fbaube/lwdx"
	CLI "github.com/fbaube/m5cli"

	// S3 "github.com/fbaube/sqlite3"
	"os"
)

// init exercises the DB and driver, to try
// to head off this obnoxious runtime error:
// Binary was compiled with 'CGO_ENABLED=0',
// go-sqlite3 requires cgo to work. This is a stub
/* func init() {
	var pDB *sql.DB
	var e error
	pDB, e = sql.Open("sqlite3", "file:./_test.db")
	if e != nil {
		panic("sql.Open failed")
	}
	e = pDB.Ping()
	if e != nil {
		panic(fmt.Sprintf("sql.DB.Ping failed: %s", e.Error()))
	}

	e = pDB.Close()
	if e != nil {
		panic("sql.DB.Close failed")
	}
} */

// CLI.CLI is in effect the application in its entirely, except
// that there is no GUI: it parses the invocation CLI arguments
// found in os.Args and then executes the application.
//
// NOTE that the variable os.Args is writeable; this means it can
// be used for testing the app or for compiling-in a configuration
// for a WASM runtime environment.
//
// If the CLI user invokes this app with no arguments, the app should
// print a usage message, per standard Unix convention. In such case
//  1. detecting the condition and generating the message must be
//     handled by CLI.CLI, where knowledge of arguments resides, and
//  2. the flarg (i.e. flag argument) definition must first be parsed
//     by package [pflag] for the help message to be correct, and
//  3. CLI.CLI must ensure that the error it returns is appropriately
//     informative and complete, because main() does not otherwise
//     embellish the returned error message.
//
// .
func main() {
	var e error

	// BI := GetBuildInfo()
	// fmt.Printf("BUILD INFO: %s %s \n"+
	//    "-------------------------------------------\n",
	//     BI.GoVersion, BI.Path)

	// NOTE: For TESTING purposes, you may modify
	// [os.Args] right here at this point in the code,
	// and then verify it with this next Printf statement.
	fmt.Printf("%s: %s  (os.Args)\n", os.Args[0], os.Args[1:])

	// S3.X() // DBSimpleTest()

	if e = CLI.CLI(); e == nil {
		// SUCCESS!
		// L.L.Flush()
		// !! os.Exit(0)
	} else {
		// L.L.Flush()
		fmt.Fprintf(os.Stderr,
			"%s encountered a fatal error: \n  %s \n",
			os.Args[0], e.Error())
		os.Exit(1)
	}
}
