package m5cli

import (
	"errors"
	"fmt"
	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog"
	RS "github.com/fbaube/reposqlite"
	RM "github.com/fbaube/rowmodels"
	_ "github.com/fbaube/sqlite3"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
	"os"
	// "runtime/pprof"
	// "github.com/davecheney/profile"
	// "github.com/fbaube/tags"
	_ "database/sql"
)

var DB_FILENAME = "mmmc.db"

// FIXME:
// hide a flag by specifying its name
// flags.MarkHidden("secretFlag")

/* FIXME:
// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
}
*/

var multipleXmlCatalogFiles []*XU.XmlCatalogFile

// ProcessDatabaseArgs should be able to process
// either a new DB OR an existing DB.
// .
func (env *XmlAppEnv) ProcessDatabaseArgs() error {

	// type-checking
	// var _ repo.SimpleRepo = (*RS.SqliteRepo)(nil)

	var mustAccessTheDB bool
	var e error
	mustAccessTheDB = env.cfg.b.DBdoImport ||
		env.cfg.b.DBdoZeroOut || env.cfg.p.sDbdir != ""
	if !mustAccessTheDB {
		return nil
	}
	// Start by checking on the status of the filename.
	// This all assumes that the DB is SQLite, a single file.
	// Note that a path is used to derive a FILE path.
	var dbFilepath string
	// println("misc.go: BEFOR:", env.cfg.p.sDbdir)
	// NOTE that if env.cfg.p.sDbdir is "", ResolvePath won't fix it!
	if env.cfg.p.sDbdir == "" {
		env.cfg.p.sDbdir = "."
	}
	// println("misc.go: BEFOR:", env.cfg.p.sDbdir)
	dbFilepath = FU.ResolvePath(
		env.cfg.p.sDbdir + FU.PathSep + DB_FILENAME)
	L.L.Info("DB resolved path: " + dbFilepath)
	errPfx := fmt.Errorf("processDBargs(%s):", dbFilepath)
	// func IsFileAtPath(aPath string) (bool, *os.FileInfo, error) {

	var fileinfo os.FileInfo
	filexist, fileinfo, filerror := FU.IsFileAtPath(dbFilepath)
	if filerror != nil {
		panic("L71")
		return fmt.Errorf("%s file error: %w", errPfx, filerror)
	}
	s := SU.ElideHomeDir(dbFilepath)
	if filexist {
		L.L.Info("DB exists: " + s)
		if fileinfo.Size() == 0 {
			L.L.Info("DB is empty: " + s)
			e = os.Remove(dbFilepath)
			if e != nil {
				panic(e)
			}
			filexist = false
		} else {
			env.SimpleRepo, e = RS.OpenRepoAtPath(dbFilepath)
			// If the DB exists and we want to open
			// it as-is, i.e. without zeroing it out,
			// then this is where we return success:
			if e == nil && !env.cfg.b.DBdoZeroOut {
				L.L.Info("DB opened: " + s)
				return nil
			}
		}
	}
	if !filexist {
		L.L.Info("Creating DB: " + s)
		if env.cfg.b.DBdoZeroOut {
			L.L.Info("Zeroing out the DB is redundant")
		}
		env.SimpleRepo, e = RS.NewRepoAtPath(dbFilepath)
	}
	if e != nil {
		return fmt.Errorf("%s DB failure: %w", errPfx, e)
	}
	repoAbsPath := env.SimpleRepo.Path()
	L.L.Info("DB OK: " + SU.ElideHomeDir(repoAbsPath))

	pSQR, ok := env.SimpleRepo.(*RS.SqliteRepo)
	if !ok {
		panic("L100")
		return errors.New("processDBargs: is not sqlite")
	}
	e = pSQR.SetAppTables("", RM.MmmcTableDescriptors)
	/* type RepoAppTables interface {
		// SetAppTables specifies schemata
		SetAppTables(string, []U.TableConfig) error
		// EmptyAllTables deletes (app-level) data
		EmptyAppTables() error
		// CreateTables creates/empties the app's tables
		CreateAppTables() error
	} */
	if !filexist {
		// env.SimpleRepo.ForceExistDBandTables()
		e = pSQR.CreateAppTables()

	} else if env.cfg.b.DBdoZeroOut {
		L.L.Progress("Zeroing out DB")
		_, e := env.SimpleRepo.CopyToBackup()
		if e != nil {
			panic(e)
		}
		pSQR.EmptyAppTables()
	}
	return nil
}

// The general approach:
// 1. rel.filepath
// 2. abs.filepath
// 3. PathProps
// 4. ContentityRecord
// 5. MCFile
// 6. GTree
// 7. ForesTree

// inputExts more than covers the file types associated with the LwDITA spec.
// Of course, when we check for them we do so case-insensitively.
var inputExts = []string{
	".dita", ".map", ".ditamap", ".xml",
	".md", ".markdown", ".mdown", ".mkdn",
	".html", ".htm", ".xhtml", ".png", ".gif", ".jpg"}

// AllGLinks gathers all GLinks in the current run's input set.
var AllGLinks mcfile.GLinks
