package m5cli

import (
	"errors"
	"fmt"
	"os/signal" 
	L "github.com/fbaube/mlog"
	"github.com/fbaube/rest"
	WU "github.com/fbaube/wasmutils"
	"context" 

	// flag "github.com/spf13/pflag"
	"os"
	"runtime/debug"
	"time"
)

// ctx is a global to quiet the compiler.
var ctx context.Context
var cancelFunc context.CancelFunc 

// The general approach:
// 1. Filename via command line (RelFP = relative filepath)
// 2. Filename absolute path (AbsFP)
// 3. PathProps
// 4. ContentityRecord
// 5. MCFile
// 6. GTree
// 7. ForesTree

// CLI parses the arguments passed in (tipicly [os.Args]) and then
// processes them; therefore it is very easy to pass in whatever
// arguments are suitable for testing.
//
// NOTE: [os.Args] is writable so you can also assign your own set
// of argument - this might be useful for testing or for WASM usage.
//
// An error from this func is returned unmodified and unprocessed,
// and so it is up to the caller to sort out its severity and how to
// handle it (and perhaps what to return to the shell via [os.Exit]).
//
// NOTE: Do not use logging in the code until the point
// where the command line invocation has been sorted out.
//
// Outline of operation (possibly OBS): 
//  0. flargs (command line "flag arguments") are defined
//     in a func init(), so that they are available for
//     a no-CLI-arg help message
//  1. Check for no-CLI-arg invocation that gets a help
//     message and exits
//  2. InitLogging
//  3. Samples() can do library demos and other very meta stuff 
//  4. NewXmlAppCfg creates XmlAppCfg config from CLI arguments
//  5. NewXmlAppEnv creates XmlAppEnv env'mt from XmlAppCfg
//  6. XmlAppEnv.Exec() to Get Things Done
//  7. If REST port nr given, run web UI
//
// .
func CLI(args []string) error {

     	// Top-Level Panic Recovery 
	defer func() {
	   if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr,
			"Panic recovery stack dump: %#v \n", err)
		println("STACK DUMP:")
		debug.PrintStack()
		}
	}()

	// Unix-style request for help ? 
	if len(args) < 2 && !WU.IsBrowser() {
		myUsage()
		return errors.New("No arguments. Nothing to do.")
	}
	
	// Initialize logging 
	InitLogging(args[0])
	// ===============
	//  IT'S NOW OKAY
	//  TO USE LOGGING
	// ===============
	L.SetMaxLevel(LOG_LEVEL_FILE_INTRO)

	// ==============
	//  CONVERT args
	//  TO AN AppCfg
	// ==============
	// TODO: Use this rather than edit the
	// flargs init stuff: DisableFlags("hDgpr")
	var cfg *XmlAppCfg
	var e error
	if cfg, e = newXmlAppCfg(args); e != nil {
		L.L.Flush()
		return e
	}
	// ===============
	//  NOW IT'S OKAY
	//  TO USE FLARGS
	// ===============

	if cfg.b.Samples {
		DoSamples()
		// We're not really using contexts yet, but...
		// Let's think about catching Control-C
		//
		//  func NotifyContext(
		//  	  parent context.Context, signals ...os.Signal)
		//  (ctx context.Context, stop context.CancelFunc)
		//
		// NotifyContext returns a copy of the parent context
		// that will be marked done (i.e. its Done channel 
		// will be closed) when the first of these happens: 
		//  - one of the listed signals arrives, or
		//  - the returned stop function is called, or
		//  - the parent context's Done channel is closed
		ctx = context.Background()		
		ctx, cancelFunc =
		     signal.NotifyContext(ctx, os.Interrupt)
		defer cancelFunc()
		// fmt.Printf("SignalCtx: %+v \n", ctx)
	}

	// ===================
	//  CONVERT AppCfg TO
	//  A RUNNABLE AppEnv
	// ===================
	var env *XmlAppEnv
	if env, e = cfg.newXmlAppEnv(); e != nil {
		L.L.Flush()
		L.L.Error("AppEnv cannot Exec():", e.Error())
		return e
	}
	L.L.Debug("OK to Exec()...")

	// =====
	//  RUN
	// =====
	if e = env.Exec(); e != nil {
		L.L.Flush()
		println("Exec:", e.Error())
		L.L.Error("Exec: " + e.Error()) 
		// return e
	}
	L.L.Flush()
	// Give messages a chance to get visible.
	time.Sleep(300 * time.Millisecond)

	if cfg.AllFlargs.webPort != 0 {
		RunWeb(cfg.AllFlargs.webPort)
	} else if cfg.AllFlargs.restPort != 0 {
		rest.RunRest(cfg.AllFlargs.restPort, env.SimpleRepo)
	}
	return e
}
