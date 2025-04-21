package m5cli

import (
	"errors"
	"io"

	D "github.com/fbaube/dsmnd"
	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog"
	DRP "github.com/fbaube/datarepo"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
	"github.com/fbaube/m5db"
)

// XmlAppEnv should be usable in other apps & scenarios too. 
type XmlAppEnv struct {
	cfg *XmlAppCfg
	DRP.SimpleRepo
	InputPathItems 
	Infiles       []FU.FSItem // bye 
	Indirs        []FU.FSItem // bye 
	IndirFSs      []mcfile.ContentityFS // bye 
	Outdir, Dbdir FU.FSItem // NOT ptr! Barfs at startup.
	Xmlcatfile    FU.FSItem // NOT ptr! Barfs at startup.
	Xmlschemasdir FU.FSItem // NOT ptr! Barfs at startup.
	// IsSingleFile is a convenience flag, and a
	// result of processing CLI arg for input file(s)
	IsSingleFile bool
	// Result of processing CLI args -c & -s
	*XU.XmlCatalogFile
	PrittOutput io.Writer
}

// ContentityProcessor (see ContentityStage !!) is one 
// way to go. It preserves order of processing of MCFile's 
// (unlike iterating thru a map of them), and the func 
// signature is in the Go style, and the style IS CHAINABLE.
// 
// Note that interface [fileutils.Errer] makes this type
// UNNECESSARY.
//
// Note that to declare a func that is a ContentityProcessor,
// the func's signature is the RH side of this, NOT the LH 
// side: don't try to declare a "ContentityProcessor" named 
// as such. Treat it like an interface, which is instantiated 
// by doing, not by saying. 
// . 
// type ContentityProcessor func(
//	p *mcfile.Contentity, e error) (*mcfile.Contentity, error)

// newXmlAppEnv turns an XmlAppCfg into an XmlAppEnv.
func (cfg *XmlAppCfg) newXmlAppEnv() (*XmlAppEnv, error) {
	var env *XmlAppEnv
	var e error
	env = new(XmlAppEnv)
	env.cfg = cfg

	// =======================================
	//   PROCESS DATABASE DIRECTORY ARGUMENT
	// =======================================
	L.L.Warning(SU.Rfg(SU.Ybg("=== CLI DATABASE ===")))
	dbargs := *new(DRP.Init9nArgs)
	dbargs.DB_type = D.DB_SQLite
	dbargs.BaseFilename = "m5" // DRP.DEFAULT_FILENAME // if omitted, still default! 
	// A relative filepath is OK
	dbargs.Dir = env.cfg.p.sDbdir
	dbargs.DoImport = env.cfg.b.DBdoImport
	dbargs.DoZeroOut = env.cfg.b.DBdoZeroOut
	dbargs.DoBackup = true 
	dbargs.TableDetailz = m5db.M5_TableDetails
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
	return env, nil
}
