package m5cli

import (
	"errors"
	"fmt"
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

// InputPathItems is for gathering, expanding (directories),
// verifying, and loading files and directories specified 
// on the command line, and then organising everything as 
// one large array of [mcfile.Contentity].
//  - [NamedPaths] is an input slice of paths of files and
//    directories; a path to a file that ends with "/" (or
//    os.Sep) throws a panic 
//  - [NamedFiles] is a slice of [fileutils.FSItem]
//    for files named e.g. on the CLI
//  - [NamedDirs]  is a slice of [fileutils.FSItem]
//    for dirs  named e.g. on the CLI
//  - [DirCntyFSs] is a slice of [mcfile.ContentityFS],
//    one per element of [NamedDirs]
//  - [AllCntys] is an output slice of [mcfile.Contentity] that 
//    collects all Contentities (a) named by [NamedFiles], and
//    (b) gathered by expanding [NamedDirs] and then walking
//    their [DirCntyFSs]
// . 
type InputPathItems struct {
        NamedPaths []string    // Input 
	NamedFiles []FU.FSItem // or FSItemInfo ?
	NamedDirs  []FU.FSItem // or FSItemInfo ?
	DirCntyFSs []mcfile.ContentityFS
	AllCntys   []mcfile.Contentity
}

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
	// ====================================
	//   PROCESS INPUT PATHS, to get info
	//   about paths, existence, and types
	// ====================================
	L.L.Warning(SU.Rfg(SU.Ybg("=== CLI PATH(S) ===")))
	L.L.Debug("AppCfg.sInpaths: %+v", cfg.p.sInpaths)
	var InputFSItems []*FU.FSItem
	var EE []error
	for _, path := range cfg.p.sInpaths {
	        L.L.Info("AppEnv: do input path: " + path)
		npp, err := FU.NewFSItem(path)
		if err != nil {
		   fmt.Printf("AppEnv: GOT ERROR")
		 } else {
		   InputFSItems = append(InputFSItems, npp)
		   // FIXME: BAD HACK - about doubly-nil interfaces 
		   if InputFSItems != nil {
		      EE = append(EE, nil)
		   } else {
		     EE = append(EE, err)
		   }
		 }
	}
	L.L.Info("%d input path(s) yielded %d F/S item(s)",
		len(cfg.p.sInpaths), len(InputFSItems))
	
	L.L.Warning(SU.Rfg(SU.Ybg("=== CLI F/S ITEM(S) ===")))
	for i, pp := range InputFSItems {
	        if pp == nil || pp.HasError() { panic("NIL or GOT ERROR") }
		inp := SU.Tildotted(pp.FPs.AbsFP)
		msg := fmt.Sprintf("[%d:%s] is ", i, inp)
		if EE[i] != nil {
			L.L.Error("TRIGR'D! EE[i] :: %T %#v", EE[i], EE[i])
			L.L.Error(msg + "ERROR: " + EE[i].Error())
			continue
		}
		// TRUE?! It's a problem here because we need to check 
		// the file type without yet having a FileInfo 
		// to work with. So we do it the hard way.
		var sType string
		// sType = pp.Code4L()
		if pp.IsDir()  { sType = "DIRR" } else
		if pp.IsFile() { sType = "FILE" } else
		if pp.IsSymlink() { sType = "SYML" } else
		{ sType = "OTHR" }
		// println(">>>", msg, sType)
		var sNote string
		switch sType {
		case "DIRR":
			env.Indirs = append(env.Indirs, *pp)
			sNote = "to process recursively"
			// L.L.Info("Directory, to be processed recursively")
		case "FILE":
			env.Infiles = append(env.Infiles, *pp)
			// L.L.Info("File")
		case "SYML":
			sNote = "(TODO: check CLI symlink flag)"
			// L.L.Info("File")
		case "OTHR":
			sNote = "Unknown type: not file, not dir, not symlink"
			L.L.Error(msg + sNote)
			return env, errors.New(sNote + "Bad input")
		default: // and case "Non-existent"
			sNote = "Does not exist or is extremely weird"
			L.L.Error(msg + sNote)
			return env, errors.New(sNote + "Bad input")
		}
		L.L.Info(msg+"%s: %s", sType, sNote)
	}
	if len(env.Infiles) == 1 && len(env.Indirs) == 0 {
		env.IsSingleFile = true
	}
	L.L.Info("CLI arguments (unexpanded): %d file(s), %d folder(s)",
		len(env.Infiles), len(env.Indirs))
	return env, nil
}
