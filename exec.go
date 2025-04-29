package m5cli

import (
	"fmt"
	"os"
	
	DRS "github.com/fbaube/datarepo/sqlite"
	"github.com/fbaube/m5cli/exec"
	L "github.com/fbaube/mlog" // Bring in global var L
	SU "github.com/fbaube/stringutils"
	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"

	"errors"
	"io/fs"
	FP "path/filepath"
)

func fpt(path string) string {
     var A, V, L bool
     var sA, sL string
     var eA, eL error 
     A = FP.IsAbs(path)
     V = fs.ValidPath(path)
     L = FP.IsLocal(path)
     sA, eA = FP.Abs(path)
     sL, eL = FP.Localize(path)
     if eA == nil { eA = errors.New("OK") }
     if eL == nil { eL = errors.New("OK") }
     nF, nE := os.Open(path)
     rF, rE :=  R.Open(path) // this line barfs on symlink to ".."!
     nF.Close()
     rF.Close()
     return fmt.Sprintf("Path: %s \n" +
     	    "Rel:%s LV:%s%s Abs<%s:%s> Lcl<%s:%s> \n" +
     	    "norm.Open.error: %s \n" +
	    "root.Open.error: %s \n", 
     	    path, SU.Yn(!A), SU.Yn(L), SU.Yn(V), sA, eA, sL, eL, nE, rE)
}

// The general approach:
//  1) os.Args
//  2) AppCfg
//  3) AppEnv
//  4) Run()!
//
// How refs to files (and directories) enter the
// system (altho levels might be intermixed here)
// (and this might be kinda outa date):
//  1) Filename via cmd line (can be Rel.FP)
//  2) Filename absolute path  (i.e. Abs.FP)
//  3) FSItem
//  4) Loading & parsing: GTokens
//  5) PathAnalysis (i.e. content analysis)
//  6) ContentityRow
//  7) Contentity
//  8) GTree
//  9) ContentiTree

var R *os.Root

// Exec does all execution of all stages for
// every [mcfile.Contentity], altho only after
// all prep has already been done by other funcs.
func (env *XmlAppEnv) Exec() error {

        var e error 
     	R, e = os.OpenRoot(".")
	if e != nil { panic("OOPS Root") }
     	println(fpt(""))
     	println(fpt("."))
     	println(fpt(".."))
     	println(fpt("../"))
     	println(fpt("../../"))
     	println(fpt("/"))
     	println(fpt("/etc"))
     	println(fpt("/etc/"))
     	println(fpt("derf"))
     	println(fpt("derf/derf2"))
	println(fpt("/Users/fbaube/src/m5app/m5/m5"))
	println(fpt("/Users/fbaube/src/m5app/m5/m5/derf/"))
	println(fpt("tstat/L-etc"))
	println(fpt("tstat/L-file-Nexist"))
	println(fpt("tstat/L-file-OK"))
	println("=> tilde")
	println(fpt("tstat/L-tilde"))
	// println("=> double dot:")
	// println(fpt("tstat/L-par-dbldot"))
	
	// =====================
	// =====================
	// TOP LEVEL: FILE INTRO
	//    file_reading_01
	// =====================
	// =====================
	L.SetMaxLevel(LOG_LEVEL_FILE_INTRO)
	defer func() { L.L.Flush() }()
	// Timing:
	// tt := MU.Into("Input file processing")
	
	e01 := file_reading_01(&(env.InputPathItems))
	if e01 != nil {
	   L.L.Error("File reading failed: %s", e01)
	   return fmt.Errorf("exec.filereading: %w", e01)
	}
	InitContentityDebugFiles(env.AllCntys, env.cfg.b.TotalTextal)
	
	// =========================
	// =========================
	// TOP LEVEL: EXECUTE STAGES
	// =========================
	e02 := exec_stages_2(env.AllCntys)
	if e02 != nil {
	   L.L.Error("Exec stages failed: %s", e02)
	   return fmt.Errorf("exec.execstages: %w", e02)
	}
	
	// ============================
	// ============================
	// TOP LEVEL: INTRA-FILE AND
	// INTER-FILE REFERENCE LINKING
	// ============================
	// ============================
	e03 := ref_linking_03(env.AllCntys)
	if e03 != nil {
	   L.L.Error("Ref linking failed: %s", e03)
	   return fmt.Errorf("exec.reflinking: %w", e03)
	}
	// =======================
	//   VALIDATE INPUT FILES
	// =======================
	e04 := validateInputFiles(env) 
	if e04 != nil {
	   L.L.Error("Input validation failed: %s", e04)
	   return fmt.Errorf("exec.valdateinputs: %w", e04)
	}
	
	// ==========================
	//  LOAD FILES INTO DATABASE
	// ==========================
	if env.cfg.b.DBdoImport {
	
		if haveDB := (env.SimpleRepo != nil); !haveDB {
			L.L.Error("Cannot proceed: SqliteRepo is not valid")
			os.Exit(1)
		}
		// ================================
		//   Verify the type of the repo
		//  (future-proofing) and then do
		//  type conversion (lame hackery)
		// ===============================
		var pSR *DRS.SqliteRepo
		var ok bool
		pSR, ok = env.SimpleRepo.(*DRS.SqliteRepo)
		if !ok {
			panic("Exec: repo is not *SimpleSqliteRepo")
		}
		// =================
		//   NOW IT'S OKAY
		//    TO PROCEED
		//  WITH THE IMPORT
		// =================
		importError := exec.ImportBatchIntoDB(
			pSR /* env.SimpleRepo */, env.AllCntys)
		if importError != nil {
			L.L.Error("exec.ImportBatchIntoDB failed: %w", importError)
		}
		pSR.CloseLogWriter()
		pSR.Flush()
	}
	return nil
}

