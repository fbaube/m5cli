package mcm_cli

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"errors"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	WU "github.com/fbaube/wasmutils"
	flag "github.com/spf13/pflag"
)

// XmlAppCfg can probably be used with various 3rd-party utilities.
type XmlAppCfg struct {
	AppName     string
	CmdTail     []string
	XmllintPath string
	AllFlargs
}

// newXmlAppCfg parses and processes CLI arguments passed to it.
// Be sure to include args[0], which is the CLI invocation name of
// the program. If no args are passed (i.e. args is nil or 0-length),
// it defaults to parsing os.Args instead. Else (i.e. len(args)>0)
// it copies args to os.Args before parsing begins.
//
// NOTE that args is probably nil !!
func newXmlAppCfg(args []string) (*XmlAppCfg, error) {
	var cfg *XmlAppCfg
	var e error

	fmt.Printf("DBG ARGS :: %+v \n", os.Args)

	if args != nil && len(args) > 0 {
		println("Arguments provided internally are overriding os.Args")
		os.Args = args
	}
	cfg = new(XmlAppCfg)
	// This if-test should be unnecessary, cos we should have
	// already caught a no-args invocation, issued a usage
	// message, and exited. But, well, it is future-proofing.
	if os.Args != nil && len(os.Args) > 0 {
		cfg.AppName = os.Args[0]
	}
	var pFlargs *AllFlargs
	if pFlargs, e = parseFlargs(); e != nil {
		return nil, e
	}
	cfg.AllFlargs = *pFlargs

	// At this point, package [pflag] has parsed the command
	// line and loaded its singletons, including flag.CommandLine.

	// Comment this out, cos if we only use "-r" there are no files.
	/* if len(flag.Args()) == 0 {
		panic("OOPS")
	} */
	cfg.CmdTail = flag.Args()
	cfg.p.sInpaths = flag.Args()
	L.L.Dbg("CLI bool flargs: " + cfg.b.String())
	L.L.Dbg("CLI path flargs: %s", cfg.p.String())
	L.L.Dbg("Cmd tail: %+v  (flag.Args)", cfg.CmdTail)
	L.L.Dbg("In_paths: %+v  (flag.Args)", cfg.p.sInpaths)

	// Handle case where XML comes from standard input i.e. os.Stdin
	if flag.Args() != nil && len(flag.Args()) > 0 &&
		flag.Args()[0] == "-" {

		if WU.IsBrowser() {
			return cfg, errors.New("FIXME/wasm: " +
				"Trying to read from Stdin")
		}
		var stat fs.FileInfo
		stat, e = os.Stdin.Stat()
		if e != nil {
			return cfg, errors.New("Cannot Stat() Stdin")
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			return cfg, errors.New("Cannot read Stdin " +
				"when not interactive (file or pipe)")
		}
		println("==> Reading from Stdin; " +
			"press ^D right after a newline to end")
		// bb, e := ReadAll(os.Stdin)
		var stdIn string
		stdIn, e = FU.GetStringFromStdin()
		if e != nil {
			return cfg, errors.New("Cannot read from Stdin")
		}
		e = os.WriteFile("Stdin.xml", []byte(stdIn), 0666)
		if e != nil {
			return cfg, errors.New("Cannot write to ./Stdin.xml")
		}
		// cfg.Infile = *FU.NewPathProps("Stdin.xml") // .RelFilePath = "Stdin.xml"
		cfg.p.sInpaths[0] = "./Stdin.xml"
	}
	if !cfg.b.Validate {
		return cfg, nil
	}
	if WU.IsBrowser() {
		return cfg, errors.New("Validation not possible: " +
			"tools not available in browser")
	}
	// Locate xmllint for doing XML validations
	cfg.XmllintPath, e = exec.LookPath("xmllint")
	if e != nil {
		return cfg, errors.New("Validation not possible: " +
			"xmllint cannot be found")
	}
	L.L.Info("xmllint found at: " + cfg.XmllintPath)
	return cfg, nil
}
