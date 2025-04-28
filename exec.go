package m5cli

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	S "strings"

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
	// =========================
	//  FOR EVERY CONTENTITY
	// ======================
	L.SetMaxLevel(LOG_LEVEL_EXEC_STAGES)
	L.L.Okay(SU.Rfg(SU.Ybg("=== DO CONTENTITY STAGES ===")))

	L.L.Info("Input contentities: total %d", len(env.AllCntys))

	for ii, cty := range env.AllCntys {
		// Skip directories entirely
		if cty.IsDir() {
			continue
		}
		// Complain loudly if the contentity type is unidentified
		if cty.RawType() == "" { // or SU.Raw_type_UNK {
			L.L.Error("UNK RawType in ExecuteStages (2nd chance)")
		}
		var dsp string
		L.L.SetCategory(fmt.Sprintf("%02d", ii))
		dsp = fmt.Sprintf("[F%02d] %s", ii, SU.Tildotted(cty.AbsFP()))
		L.L.Info(SU.Cyanbg(SU.Wfg(dsp)))
		var rawlen int
		if cty.FSItem.TypedRaw != nil {
			rawlen = len(cty.FSItem.Raw)
		}
		dsp = fmt.Sprintf(" len:%4d  %s  %s  %s",
			rawlen, cty.RawType(),
			cmp.Or(cty.MType, "(nil-MType)"),
			cmp.Or(cty.MimeType, "(nil-Mime)"))
		L.L.Info(SU.Cyanbg(SU.Wfg(dsp)))

		// Now try dumping it as JSON !
		// func MarshalIndent(v any, prefix,
		//    indent string) ([]byte, error)
		var jsonOut []byte
		jsonOut, _ = json.MarshalIndent(cty, "", "  ")
		/*
			L.L.Info("================\n" +
				 "%s \n" +
				 "================", string(jsonOut))
		*/
		jsonOutFilename := cty.FPs.AbsFP + ".tmp.json"
		errr := os.WriteFile(jsonOutFilename, jsonOut, 0644)
		if errr != nil {
			panic("os.WriteFile json: " + jsonOutFilename)
		}
		L.L.Info("Wrote JSON to: " + jsonOutFilename)
		// defer jsonOutFile.Close()
		// ==================
		//  AND NOW, EXECUTE
		// ==================
		cty.ExecuteStages()
	}

	// =========================
	//  DUMP ACCUMULATED ERRORS
	// =========================
	L.L.Info("Accumulated errors:")
	for _, p := range env.AllCntys {
		if p.HasError() {
			L.L.Error("file[%s: %s", p.LogPrefix("]"), p.Error())
		}
	}
	/*
		// ================================
		//  2a. PREPARE FOR OUTPUT FILES
		//      Maybe ConfigureOutputFiles
		// ================================
		(for every InfileContentity)
		if env.GroupGenerated {
			// println("==> Creating output directories.")
			e = pMCF.ConfigureOutputFiles("_" + myAppName)
			errorbarf(e, "ConfigureOutputFiles")
			}
		}
	*/
	// ============================
	// ============================
	// TOP LEVEL: INTRA-FILE AND
	// INTER-FILE REFERENCE LINKING
	// ============================
	// ============================
	for _, p := range env.AllCntys {

		// 2025.04 FIXME FIXME FIXME
		p.GatherXmlGLinks() // (&AllGLinks)

		// Cross-reference and resolve the links
		println("D=> TODO: Cross-ref the GLinks")

		// MU.Outa("Processing input file(s)", tt)

		/*
			fmt.Printf("==> Summary counts: %d Tags, %d Atts \n",
				mcfile.GlobalTagCount, mcfile.GlobalAttCount)
			println("--> Tags:", mcfile.GlobalTagTally.StringSortedValues())
			println("--> Atts:", mcfile.GlobalAttTally.StringSortedValues())
		*/
		println("#### GLink KEY SOURCES ####")
		for _, pGL := range AllGLinks.KeyRefncs {
			fmt.Printf("%s@%s: %s: %s \n", pGL.Att, pGL.Tag,
				pGL.AddressMode, pGL.AbsFP.Tildotted())
		}
		println("#### GLink KEY TARGETS ####")
		for _, pGL := range AllGLinks.KeyRefnts {
			t := pGL.Tag
			a := pGL.Att
			// if S.HasPrefix(t, "topi") ||
			// S.Contains(a, "key") || S.Contains(a, "ref") {
			fmt.Printf("%s@%s: %s: %s \n",
				a, t, pGL.AddressMode, pGL.AbsFP.Tildotted())
			// }
		}
		println("#### GLink URI SOURCES ####")
		for _, pGL := range AllGLinks.UriRefncs {
			fmt.Printf("%s@%s: %s: %s \n", pGL.Att, pGL.Tag,
				pGL.AddressMode, pGL.AbsFP.Tildotted())
		}
		println("#### GLink URI TARGETS ####")
		for _, pGL := range AllGLinks.UriRefnts {
			t := pGL.Tag
			a := pGL.Att
			isList := (len(pGL.Tag) == 2 &&
				S.HasSuffix(pGL.Tag, "l"))
			if !isList {
				fmt.Printf("%s@%s: %s: %s \n", a, t,
					pGL.AddressMode, pGL.AbsFP.Tildotted())
			}
		}
	}

	// Cross-reference and resolve the links

	// ===========================
	//   3. VALIDATE INPUT FILES
	//  This code actually belongs
	//     above: see line 145.
	// ===========================

	// We can use xmllint to validate here, but we don't
	// want to rely on schema files already in the system
	// - no using its normal catalogs, or anything at or
	// under `/etc/xml/catalog` - and we can't use the
	// envar `XML_CATALOG_FILES`. But if we have our own
	// catalog file, we can pass it our own value as (say)
	// envar `MMMC_XML_CATALOG_FILES`.

	/* if pCA.Validate && * / env.XmlCatalogFile != nil {

		println(" ")
		tt = MU.Into("")
		println("==> Validating input file(s)...")

		var dtdStatus, docStatus, errStatus string

		print("==> Text file validation statuses: \n")
		for _, pMCF = range MCFiles {
			if pMCF.IsXML() {
				if FU.MTypeSub(pMCF.MType, 0) == "img" {
					continue
				}
			}
			if !pMCF.IsXML() {
				continue
			}
			var dtdDesc string

			// DO THE VALIDATION
			dtdStatus, docStatus, errStatus = pMCF.DoValidation(env.XmlCatalogFile)

			if pMCF.XmlDoctypeFields != nil {
				dtdDesc = pMCF.XmlDoctypeFields.PIDSIDcatalogFileRecord.PublicTextDesc
			}
			fmt.Printf("%s/%s/%s %s %s :: %s :: %s  \n",
				pMCF.MType[0], pMCF.MType[1], pMCF.MType[2], dtdStatus,
				docStatus, pMCF.AbsFilePath, dtdDesc)
			if errStatus != "" {
				println(errStatus)
			}
		}
		MU.Outa("Validating input file(s)", tt)
	}
	*/

	// ===============================
	//   4. LOAD FILES INTO DATABASE
	//      See also above
	// ===============================

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

func prerr(e error) string {
	if e == nil {
		return "-"
	}
	return e.Error()
}
