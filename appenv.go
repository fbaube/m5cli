package m5cli

import (
	"errors"
	"fmt"
	"io"

	D "github.com/fbaube/dsmnd"
	// DI "github.com/fbaube/dbinit"
	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog"
	DR "github.com/fbaube/datarepo"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
	DRM "github.com/fbaube/datarepo/rowmodels"
)

type XmlAppEnv struct {
	cfg *XmlAppCfg
	DR.SimpleRepo
	Infiles       []FU.FSItem
	Indirs        []FU.FSItem
	IndirFSs      []mcfile.ContentityFS // was: []ExpanDir
	Outdir, Dbdir FU.FSItem // NOT ptr! Barfs at startup.
	Xmlcatfile    FU.FSItem // NOT ptr! Barfs at startup.
	Xmlschemasdir FU.FSItem // NOT ptr! Barfs at startup.
	// Result of processing CLI arg for input file(s)
	IsSingleFile bool
	// Result of processing CLI args -c & -s
	*XU.XmlCatalogFile
	PrittOutput io.Writer
}

// ContentityProcessor (see ContentityStage !!)
// is most probably the best way to go.
// It preserves order of processing of MCFile's (unlike
// iterating thru a map of them), and the func signature
// is most def in the Go style, and the style IS CHAINABLE.
// Note that when a ContentityProcessor is declared, the
// func signature is the RH side of this, NOT the LH side
// ("ContentityProcessor").
type ContentityProcessor func(
	p *mcfile.Contentity, e error) (*mcfile.Contentity, error)

/*
We want errors to propogate end-to-end, thru all
call chains (and maybe all type conversions too).
So, it would look something like
func NewFSItemFromRelFP (PP,err <= string,err) [dummy err]
func NewMCEfromFSItem(MCE,err <= PP,err)
func MCEprocessor(MCE,err <= MCE,err) [CHAINABLE!]
Then we plug these into samber.lo/Map(..)
(or variations on it).
*/

func (cfg *XmlAppCfg) newXmlAppEnv() (*XmlAppEnv, error) {
	var env *XmlAppEnv
	var e error
	env = new(XmlAppEnv)
	env.cfg = cfg

	// =======================================
	//   PROCESS DATABASE DIRECTORY ARGUMENT
	// =======================================
	// A relative filepath is OK
	// e = env.ProcessDatabaseArgs()
	dbargs := *new(DR.Init9nArgs)
	dbargs.DB_type = D.DB_SQLite
	dbargs.BaseFilename = "m5" // DR.DEFAULT_FILENAME // if omitted, still default! 
	dbargs.Dir = env.cfg.p.sDbdir
	dbargs.DoImport = env.cfg.b.DBdoImport
	dbargs.DoZeroOut = env.cfg.b.DBdoZeroOut
	dbargs.DoBackup = true 
	dbargs.TableDetailz = DRM.M5_TableDetails
	env.SimpleRepo, e = dbargs.ProcessInit9nArgs()

	if e != nil {
		return nil, errors.New(
			"Bad DB directory argument(s): " + e.Error())
	}
	
	// ===================================
	//   PROCESS XML CATALOG ARGUMENT(S)
	// ===================================
	L.L.Info("XML catalog processing is disabled")
	// checkbarf(e, "Could not process XML catalog argument(s)")

	// =================================
	//   PROCESS OUTPUT-FILES ARGUMENT
	// =================================
	// A relative filepath is OK
	if cfg.p.sOutdir != "" {
	   pOF, e := FU.NewFSItem(cfg.p.sOutdir) // CA.Out.RelFilePath)
	   env.Outdir = *pOF
	   if e != nil {
		return nil, errors.New("Could not process output file argument")
		}
	}
	L.L.Okay(" ")
	L.L.Okay(SU.Rfg(SU.Ybg("===                     ===")))
	L.L.Okay(SU.Rfg(SU.Ybg("=== COLLECT INPUT PATHS ===")))
	L.L.Okay(SU.Rfg(SU.Ybg("===                     ===")))
	L.L.Okay(" ")
	// ====================================
	//   PROCESS INPUT PATHS, to get info
	//   about paths, existence, and types
	// ====================================
	// fmt.Printf("cfg.sInpaths: %+v \n", cfg.sInpaths)
	var FF []*FU.FSItem
	var EE []error
	// var errct int
	// string = typeof  input arg+array,
	// *FU.FF = typeof output arg+array,
	// sInpaths = the input array
	for _, path := range cfg.p.sInpaths {
	        L.L.Info("AppEnv: do path: " + path)
		npp, err := FU.NewFSItem(path)
		FF = append(FF, npp)
		// FIXME: BAD HACK - about doubly-nil interfaces 
		if FF != nil {
		   EE = append(EE, nil)
		} else {
		   EE = append(EE, err)
		}
	}
	for i, pp := range FF {
		inp := SU.Tildotted(pp.FPs.AbsFP.S())
		msg := fmt.Sprintf("[%d:%s] ", i, inp)
		if EE[i] != nil {
			L.L.Error("TRIGRD! EE[i] :: %T %#v", EE[i], EE[i])
			L.L.Error(msg + "ERROR: " + EE[i].Error())
			continue
		}
		var sType = pp.IsWhat()
		// println(">>>", msg, sType)
		var sNote string
		switch sType {
		case "DIR":
			env.Indirs = append(env.Indirs, *pp)
			sNote = "(to process recursively)"
			// L.L.Info("Directory, to be processed recursively")
		case "FILE":
			env.Infiles = append(env.Infiles, *pp)
			// L.L.Info("File")
		case "SYMLINK":
			sNote = "(TODO: check CLI symlink flag)"
			// L.L.Info("File")
		case "UnknownType":
			sNote = "Unknown type: not file, not dir, not symlink"
			L.L.Error(msg + sNote)
			return env, errors.New(sNote + "Bad input")
		default: // and case "Non-existent"
			sNote = "Does not exist or is extremely weird"
			L.L.Error(msg + sNote)
			return env, errors.New(sNote + "Bad input")
		}
		// L.L.Okay(msg+"%s%s \n      \\\\ %+v", sType, sNote, *pp)
		L.L.Okay(msg+"%s: %s", sType, sNote)
	}
	if len(env.Infiles) == 1 && len(env.Indirs) == 0 {
		env.IsSingleFile = true
	}
	L.L.Info("CLI args (unexpanded): %d files, %d directories",
		len(env.Infiles), len(env.Indirs))
	return env, nil
}
