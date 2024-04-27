package m5cli

import (
	"errors"
	// "fmt"
	"os"
	// "flag"
	flag "github.com/spf13/pflag"
)

// parseFlargs parses the flargs ("flag arguments") passed in 
// (using [pflag.Parse]) and returns an error if the caller
// should abort execution (i.e. if no inputfile(s) specified).
//
// NOTE: Include the app's invocation name in args[0]. 
//
// The func is not exported cos it requires set-up performed
// by NewXmlAppConfig(), and so parseFlargs should be called
// by that func only.
//
// NOTE: [pflags] has been initialised to write to this
// package's variable [allFlargs], gettable from func [getAllFlargs]
//
// Both spf13/pflag and the stdib package "flags" parse os.Args .
// But here they are passed in explicitly, which can be helpful
// for testing purposes or maybe to configure a WASM execution
// environment. 
// .
func parseFlargs(args []string) (*AllFlargs, error) {
	// Disable this check, because we want to be able
	// to (hackily) use os.Args in the browser too.
	// if isBrowser() {
	//	return nil
	// }
	if len(args) < 2 {
		println("parseFlags is calling myUsage...")
		myUsage()
		return nil, errors.New("No arguments. Nothing to do.")
	}
	// Process CLI invocation flags. This is where the sausage is made.
	//
	// func Parse(): API docs: Parse parses the command-line flags
	// from os.Args[1:]. Must be called after all flags are defined
	// and before flags are accessed by the program.
	if args != nil && len(args) > 1 {
	   os.Args = args
	}
	flag.Parse()
	// fmt.Printf("parseFlags: flag.Args(): %+v \n", flag.Args())

	// Now examine the arguments not associated with
	// flargs, which should be input file(s) and dir(s)
	var paths []string
	// func Args() []string: API docs:
	// Args returns the non-flag command-line arguments.
	paths = flag.Args()
	// If no non-flag args then we have no
	// input file/directory specifier(s)
	if (allFlargs.restPort == 0) &&
		(allFlargs.webPort == 0) &&
		(nil == paths || (0 == len(paths))) {
		return nil, errors.New("Argument parsing failed. " +
			"Did not specify input file(s) and/or server?")
	}
	allFlargs.p.sInpaths = paths

	// FIXME - pos'l arg OR "-i" OR stdin OR "-"

	return allFlargs, nil
}
