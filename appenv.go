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
	Infiles       []FU.PathProps
	Indirs        []FU.PathProps
	IndirFSs      []mcfile.ContentityFS // was: []ExpanDir
	Outdir, Dbdir FU.PathProps // NOT ptr! Barfs at startup.
	Xmlcatfile    FU.PathProps // NOT ptr! Barfs at startup.
	Xmlschemasdir FU.PathProps // NOT ptr! Barfs at startup.
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
func NewPathPropsFromRelFP (PP,err <= string,err) [dummy err]
func NewMCEfromPathProps(MCE,err <= PP,err)
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
	pOF, e := FU.NewPathProps(cfg.p.sOutdir) // CA.Out.RelFilePath)
	env.Outdir = *pOF
	if e != nil {
		return nil, errors.New("Could not process output file argument")
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
	var PP []*FU.PathProps
	var EE []error
	// var errct int
	// string = typeof  input arg+array,
	// *FU.PP = typeof output arg+array,
	// sInpaths = the input array
	// NOTE that MapAwerr is written such that the
	// func requires a ptr type, and as a side effect,
	// R/W access is assured.
	/* DROP THIS
	PP, EE, errct = LOE.MapAwerr(cfg.p.sInpaths,
		// NOTE that type inference works OK
		// here and we did not need to put this
		// btwn MapAwerr and cfg.sInpaths:
		// [string, *FU.PathProps]
		func(s string) (*FU.PathProps, error) {
			return FU.NewPathProps(s)
		})
	if errct > 0 {
		L.L.Error("CLI file/dir argument processing got %d error(s)", errct)
	}
	*/
	for _, path := range cfg.p.sInpaths {
		npp, err := FU.NewPathProps(path)
		PP = append(PP, npp)
		EE = append(EE, err)
	}
	for i, pp := range PP {
		inp := SU.Tildotted(pp.AbsFP.S())
		msg := fmt.Sprintf("[%d:%s] ", i, inp)
		if EE[i] != nil {
			L.L.Error(msg + "ERROR: " + EE[i].Error())
			continue
		}
		var sType = pp.IsWhat()
		// println(">>>", msg, sType)
		var sNote string
		switch sType {
		case "DIR":
			env.Indirs = append(env.Indirs, *pp)
			sNote = " (to process recursively)"
			// L.L.Info("Directory, to be processed recursively")
		case "FILE":
			env.Infiles = append(env.Infiles, *pp)
			// L.L.Info("File")
		case "SYMLINK":
			sNote = " (TODO: check CLI symlink flag)"
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
