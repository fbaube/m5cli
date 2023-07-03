package mcm_cli

import (
	"errors"
	L "github.com/fbaube/mlog"
	WU "github.com/fbaube/wasmutils"
	// flag "github.com/spf13/pflag"
	"os"
	"runtime/debug"
	"time"
)

// The general approach:
// 1. Filename via command line (RelFP = relative filepath)
// 2. Filename absolute path (AbsFP)
// 3. PathProps
// 4. ContentityRecord
// 5. MCFile
// 6. GTree
// 7. ForesTree

// CLI parses the contents of [os.Args] and then processes them.
//
// NOTE: You can assign your own set of mock arguments to os.Args -
// it is writeable! This can be useful for testing, or for WASM usage.
//
// An error from this func is returned unmodified and unprocessed,
// and so it is up to the caller to sort out its severity and how
// to handle it.
//
// NOTE: Do not use logging in the code until the point
// where the command line invocation has been sorted out.
//
// Outline of operation:
//  0. flargs (command line flag arguments) are defined
//     in a func init(), so that they are available for
//     a no-CLI-arg help message
//  1. Check for no-CLI-arg invocation that gets a help
//     message and exits
//  2. InitLogging
//  3. PreParse() does preliminaries like library demos
//  4. NewXmlAppCfg creates XmlAppCfg config from CLI arguments
//  5. NewXmlAppEnv creates XmlAppEnv env'mt from XmlAppCfg
//  6. XmlAppEnv.Exec() to Get Things Done
//  7. If REST port nr given, run web UI
//
// .
func CLI() error {

	defer func() {
		if err := recover(); err != nil {
			println("RUNTIME PANIC! cli.CLI() failed:", err)
			debug.PrintStack()
		}
	}()
	var cfg *XmlAppCfg
	var env *XmlAppEnv
	var e error

	// ===(1)===
	if len(os.Args) < 2 && !WU.IsBrowser() {
		myUsage()
		return errors.New("No arguments. Nothing to do.")
	}
	// ===(2)===
	InitLogging(os.Args[0])
	// ===(3)===
	PreParse()

	L.L.Dbg("=============================================")

	// ===(4)===
	// There is no need to pass os.Args cos they are the default
	// but anyways they can be overwritten by any code anywhere.
	// HOWEVER note that [pflags] has been initialised to write
	// to this package's variable [allFlargs], obtainable from
	// func [getAllFlargs].
	// TODO: Re-enable this: // DisableFlags("hDgpr")
	cfg, e = newXmlAppCfg(nil)

	if e != nil {
		L.L.Flush()
		return e
	}
	// ===(5)===
	env, e = cfg.newXmlAppEnv()
	if e != nil {
		L.L.Flush()
		L.L.Progress("Cannot Exec():", e.Error())
		return e
	}
	L.L.Progress("OK to Exec()...")
	// (6)
	e = env.Exec()
	if e != nil {
		L.L.Flush()
		println("GACK! Exec() returned an error:", e.Error())
		// return e
	}
	L.L.Flush()
	time.Sleep(500 * time.Millisecond)

	if cfg.AllFlargs.webPort != 0 {
		RunWeb(cfg.AllFlargs.webPort)
	} else if cfg.AllFlargs.restPort != 0 {
		RunRest(cfg.AllFlargs.restPort)
	}
	return e
}
