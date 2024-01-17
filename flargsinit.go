package m5cli

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

// myUsage displays (1) the app name (or "(wasm)"), plus (2) a usage
// summary (see the func body), plus (3) the flags' usage message.
// TODO: Should not return info for flags that are Hidden (i.e. disabled).
func myUsage() {
	fmt.Println(os.Args[0], "[-a] [-h] [-m] [-t] [-v] [-z] [-D] [-d dbdir] [-r port] Infile")
	fmt.Println("   Process mixed content XML, XHTML/XDITA, HTML/HDITA, and Markdown/MDITA input.")
	fmt.Println("   Inpath(s)s are file or directory paths; no wildcards (?,*).")
	fmt.Println("           A directory is processed recursively.")
	fmt.Println("   The first Inpath may be \"-\" for Stdin: input that is typed (or pasted)")
	fmt.Println("           interactively is written to file ./Stdin.xml for processing")
	flag.Usage()
}

// For each DOCTYPE. also an XSL, thus XSL catalog file.
// Maybe also same for CSS.

// usageCatalog provides (here) a quick option reference,
// and (below) a means to make the actual code clearer.
// e.g. commits := map[string]int { "rsc": 3711, "r": 2138, }
// .
var usageCatalog = map[string]string{
	// Simple BOOLs
	"h": "Show extended help message and exit",
	"t": "Total textal dbg: also generate filnam.ext_[text,tkns,tree]",
	// "g": "Group all generated files in same-named folder \n" +
	//	"(e.g. ./Filenam.xml maps to ./Filenam.xml_gxml/Filenam.*)",
	// "p": "Pretty-print to file with \"fmtd-\" prepended to file ext'n",
	// "t": "gTree written to file with \"gtr-\" prepended to file ext'n",
	// "k": "gTokens written to file with \"gtk-\" prepended to file ext'n",
	"a": "Archive input file(s) to sqlar file (SQLite archive)",
	"m": "Import input file(s) to database",
	"v": "Validate input file(s) (using xmllint) (with flag \"-c\" or \"-s\")",
	"z": "Zero out the database",
	"D": "Turn on debugging",
	"L": "Follow symbolic links in directory recursion",
	// All others
	"c": "XML `catalog_path` (do not use with \"-s\" flag)",
	"d": "Database mmmc.db `dir_path`",
	"l": "`log_levels` (per processing stage, TBS)",
	"o": "`output_dir_path` (possibly ignored, depending on command)",
	"r": "Run REST server on `rest_port`",
	"s": "DTD schema file(s) `dir_path` (.dtd, .mod)",
	"w": "Run WEB server on `web_port`",
}

var allFlargs *AllFlargs

// getAllFlargs returns a global singleton.
func getAllFlargs() AllFlargs {
	return *allFlargs
}

// init loads definitions into the global [flag.CommandLine],
// and enables them all (which may or may not be appropriate).
func init() {
	allFlargs = new(AllFlargs)
	var af *AllFlargs
	af = allFlargs

	// PATH (string) FLAGS
	// f.StringVarP(&inArg, "infile", "i", "", usageCatalog["i"])
	flag.StringVarP(&af.p.sDbdir, "db-dir", "d", "", usageCatalog["d"])
	flag.StringVarP(&af.p.sOutdir, "outdir", "o", "", usageCatalog["o"])
	flag.StringVarP(&af.p.sXmlcatlgfile, "catalog-path", "c", "", usageCatalog["c"])
	flag.StringVarP(&af.p.sXmlschemasdir, "schemas-dir", "s", "", usageCatalog["s"])

	// BOOL FLAGS
	flag.BoolVarP(&af.b.Help, "help", "h", false, usageCatalog["h"])
	flag.BoolVarP(&af.b.TotalTextal, "textal", "t", false, usageCatalog["t"])
	flag.BoolVarP(&af.b.Debug, "debug", "D", false, usageCatalog["D"])
	flag.BoolVarP(&af.b.Validate, "validate", "v", false, usageCatalog["v"])
	flag.BoolVarP(&af.b.DoArchive, "archive", "a", false, usageCatalog["a"])
	flag.BoolVarP(&af.b.DBdoImport, "import", "m", false, usageCatalog["m"])
	flag.BoolVarP(&af.b.DBdoZeroOut, "zero-out", "z", false, usageCatalog["z"])
	flag.BoolVarP(&af.b.FollowSymLinks, "symlinks", "L", true, usageCatalog["L"])

	// INT FLAGS
	flag.IntVarP(&af.restPort, "rest-port", "r", 0, usageCatalog["r"])
	flag.IntVarP(&af.webPort, "web-port", "w", 0, usageCatalog["w"])

	EnableAllFlags()
	println("Initialised flargs OK")
}
