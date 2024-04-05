package m5cli

import (
	// "errors"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	DRS "github.com/fbaube/datarepo/sqlite"
	"github.com/fbaube/m5cli/exec"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog" // Brings in global var L
	SU "github.com/fbaube/stringutils"

	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"

	DRM "github.com/fbaube/datarepo/rowmodels"
)

// The general approach (semi-OBS):
// 1. Filename via cmd line (can be Rel.FP)
// 2. Filename absolute path  (i.e. Abs.FP)
// 3. PathProps
// 4. ContentityRecord
// 5. Contentity
// 5.5 GTokens
// 6. GTree
// 7. ForesTree

// Exec does all execution, altho only after all
// prep has already been done by other funcs.
// .
func (env *XmlAppEnv) Exec() error {
	var e error
	// ======================
	//  1. PRELIMINARY STUFF
	// ======================
	/* DBG
	     // Dump out what a ContentityRow looks like in the DB
		var cntro = new(DRM.ContentityRow)
		var cptrs = DRM.ColumnPtrsFuncCNT(cntro, true)
		fmt.Fprintf(os.Stderr, "ContentityRow datarepo/TableDetails: \n")
		fmt.Fprintf(os.Stderr, "\t cntRow<%T> colPtrs <%T> \n",
			cntro, cptrs)
		for iii, ppp := range cptrs {
		    fmt.Fprintf(os.Stderr, "[%d] <%T> \n", iii, ppp)
		    }
	*/
	// Timing: // tt := MU.Into("Input file processing")
	defer func() {
		L.L.Flush()
		/* ; fmt.Printf("%s: done.", os.Args[0]) */
	}()

	// ======================
	//  2. INTRO to: PROCESS
	//       ALL INPUT ITEMS
	// ======================

	L.SetMaxLevel(LOG_LEVEL_FILE_READING)
	L.L.Okay(" ")
	L.L.Okay(SU.Rfg(SU.Ybg("===                     ===")))
	L.L.Okay(SU.Rfg(SU.Ybg("=== PROCESS INPUT PATHS ===")))
	L.L.Okay(SU.Rfg(SU.Ybg("===                     ===")))
	L.L.Okay(" ")

	// At this point, "env" has three slices
	// of variables related  to input files:
	//
	// Infiles []FU.PathProps :: is all the files
	// that were specified individually on the CLI.
	// Note that if a wildcard was used, all files
	// in the expansion appear individually here.
	//
	// Indirs []FU.PathProps :: is all the directories
	// that were specified individually on the CLI.
	//
	// IndirFSs []ContentityFS (still empty at this point)
	// :: will be filled with all the entries (both files
	// and directories) found under the directories listed
	// in Indirs, as expanded into ContentityFS's, then to
	// be flattened into slices.

	// ===========================================
	//  EVERY CLI INPUT ITEM IS COLLECTED HERE
	//  First all files named at the command line,
	//  then (recursively) all directories named
	// ===========================================
	var InputContentities []*mcfile.Contentity
	var InputContentityFSs []*mcfile.ContentityFS
	var ee []error

	// DUMP env.Indirs, Inexpandirs
	L.L.Progress("env.Infiles: [%d]: %#v \n", len(env.Infiles), env.Infiles)
	L.L.Progress("env.Indirs:: [%d]: %#v \n", len(env.Indirs), env.Indirs)
	// ALSO DUMP AS JSON
	var jout []byte
	var jerr error
	if len(env.Infiles) > 0 {
		jout, jerr = json.MarshalIndent(env.Infiles[0], "infile: ", "  ")
		if jerr != nil {
			println(jerr)
			panic(jerr)
		}
		L.L.Dbg("JSON! " + string(jout))
	}
	if len(env.Indirs) > 0 {
		jout, jerr = json.MarshalIndent(env.Indirs[0], "indirr: ", "  ")
		if jerr != nil {
			println(jerr)
			panic(jerr)
		}
		L.L.Dbg("JSON! " + string(jout))
	}
	// fmt.Printf("==> env.Inexpandirs: %#v \n", env.Inexpandirs)

	// =============================
	// 2a. FOR EVERY CLI INPUT FILE
	//     Make a new Contentity
	// =============================
	InputContentities, ee = exec.LoadFilepathsContents(env.Infiles)
	gotCtys := InputContentities != nil || len(InputContentities) > 0
	gotErrs := ee != nil || len(ee) > 0
	if gotCtys || gotErrs {
		L.L.Info("RESULTS for %d infiles: %d OK, %d not OK \n",
			len(env.Indirs), len(InputContentities), len(ee))
		for i, pC := range InputContentities {
			L.L.Info("InFile[%02d] OK! [%d] %s :: %s",
				i, len(pC.FSItem.Raw), pC.MarkupTypeOfMType(),
				pC.FSItem.FPs.ShortFP)
			/* if pCty.MarkupTypeOfMType() == SU.MU_type_UNK {
					   	s := fmt.Sprintf("INfile[%d]: [%d] %s %s",
			                        i, len(pCty.PathProps.Raw),
			                        pCty.MarkupType(), pCty.AbsFP())
						panic("UNK MarkupType in ExecuteStages; \n" + s) */
		}
		for i, eC := range ee {
			L.L.Info("InfileErr[%02d] ERR :: %T", i, eC) // .Error())
		}
	}
	// ======================================
	//  2b. FOR EVERY CLI INPUT DIRECTORY
	//      Make a new Contentity filesystem
	// ======================================
	InputContentityFSs = exec.LoadDirpathsContentFSs(env.Indirs)
	for iFS, pFS := range InputContentityFSs {
		// ==============================
		// Write out a tree rep to a file
		// ==============================
		// L.L.Warning("Skip'd wrtg out tree rep: [%d]", iFS)
		var treeFile *os.File
		var treeFilename string
		// Now write out a tree representation
		treeFilename = fmt.Sprintf("./input-tree-%d", iFS)
		treeFile, e = os.Create(treeFilename)
		if e != nil {
			// An error here does not need to be fatal
			L.L.Error("Treefile " + treeFilename + ": " + e.Error())
		} else {
			pFS.RootContentity().PrintTree(treeFile)
			L.L.Okay("Wrote input tree file: " + treeFilename)
		}
		treeFile.Close()

		// =================================
		// Write out a css-enabled html tree
		// =================================
		// L.L.Warning("SKIP'D: css-enabled tree for html, exec.go.L153")
		treeFilename = fmt.Sprintf("./css-tree-%d", iFS)
		treeFile, e = os.Create(treeFilename)
		if e != nil {
			// An error here does not need to be fatal
			L.L.Error("CssTreefile " + treeFilename + ": " + e.Error())
		} else {
			pFS.RootContentity().PrintCssTree(treeFile)
			L.L.Okay("Wrote css tree file: " + treeFilename)
		}
		treeFile.Close()
	}

	// ================================
	// 2b. FOR EVERY CLI INPUT DIR
	//     Expand it into files, which
	//     also makes new Contentities
	// ================================
	for _, pED := range InputContentityFSs {
		InputContentities = append(InputContentities, pED.AsSlice()...)
	}

	// Now we have all the inputs.
	// TODO: We could count up and tell the user
	// how many files of each valid extension.

	// =========================
	//  3. FOR EVERY CONTENTITY
	//     Prepare outputs
	// =========================
	for ii, cty := range InputContentities {
		// L.L.Info("LOOPING on IN-Cty %d", ii)
		if cty.IsDir() {
			// println("Skip dir: " + cty.AbsFP())
			continue
		}
		L.L.SetCategory(fmt.Sprintf("%02d", ii))

		// Input file VALIDATION belongs HERE!
		// See the code below, line 315.

		// Now, files for capturing debugging output:
		// func Create(name string) (*File, error) API docs:
		// Create or truncate the named file. If the file
		// already exists, it is truncated (i.e. emptied).
		// If it does not exist, it is created with mode 0666
		// (before umask). If successful, the returned File
		// (IS OPEN, apparently, AND) can be used for I/O;
		// the associated file descriptor has mode O_RDWR.
		// If there is an error, it will be of type *PathError.

		cty.GTknsWriter = io.Discard
		cty.GTreeWriter = io.Discard
		cty.GEchoWriter = io.Discard

		if env.cfg.b.TotalTextal {
			fnm := cty.AbsFP()
			if len(InputContentities) == 1 {
				fnm = "./debug"
			}
			// To name output files, append
			// "_(echo,tkns,tree)" to the entire file name.
			echoName := (fnm + "_echo")
			tknsName := (fnm + "_tkns")
			treeName := (fnm + "_tree")
			f1, e1 := os.Create(echoName)
			f2, e2 := os.Create(tknsName)
			f3, e3 := os.Create(treeName)
			if e1 == nil && e2 == nil && e3 == nil {
				L.L.Info("created %s _echo,_tkns,_tree",
					SU.Tildotted(fnm))
			} else {
				L.L.Error("Cannot open a Total Textual file")
			}
			cty.GEchoWriter = f1
			cty.GTknsWriter = f2
			cty.GTreeWriter = f3
		}
	}

	// =========================
	//  4. FOR EVERY CONTENTITY
	//     Execute all stages
	// =========================
	L.SetMaxLevel(LOG_LEVEL_EXEC_STAGES)
	L.L.Okay(" ")
	L.L.Okay(SU.Rfg(SU.Ybg("===                           ===")))
	L.L.Okay(SU.Rfg(SU.Ybg("=== EXECUTE CONTENTITY STAGES ===")))
	L.L.Okay(SU.Rfg(SU.Ybg("===                           ===")))
	L.L.Okay(" ")

	L.L.Info("Input contentities: total %d", len(InputContentities))

	for ii, cty := range InputContentities {
		if cty.IsDir() {
			continue
		}
		if cty.MarkupTypeOfMType() == SU.MU_type_UNK {
			panic("UNK MarkupType in ExecuteStages (2nd chance)")
		}

		L.L.SetCategory(fmt.Sprintf("%02d", ii))
		L.L.Info(SU.Gbg("[F%02d] %s (%d) (%s) ==="), ii,
			SU.Tildotted(cty.AbsFP()), len(cty.FSItem.Raw), cty.MType)
		cty.ExecuteStages()
	}

	L.L.Info("Accumulated errors:")
	for _, p := range InputContentities {
		if p.HasError() {
			L.L.Error("file[%s: %s", p.LogPrefix("]"), p.Error())
		}
	}
	/*
		// ================================
		//  2a. PREPARE FOR OUTPUT FILES
		//      Maybe ConfigureOutputFiles
		// ================================
		if env.GroupGenerated {
			// println("==> Creating output directories.")
			e = pMCF.ConfigureOutputFiles("_" + myAppName)
			errorbarf(e, "ConfigureOutputFiles")
		}

		pMCF.GatherGLinksInto(&AllGLinks)
		// Cross-reference and resolve the links
		println("D=> TODO: Cross-ref the GLinks")

		MU.Outa("Processing input file(s)", tt)

		/*
			fmt.Printf("==> Summary counts: %d Tags, %d Atts \n",
				mcfile.GlobalTagCount, mcfile.GlobalAttCount)
			println("--> Tags:", mcfile.GlobalTagTally.StringSortedValues())
			println("--> Atts:", mcfile.GlobalAttTally.StringSortedValues())

			println("#### GLink KEY SOURCES ####")
			for _, pGL := range AllGLinks.OutgoingKeys {
				fmt.Printf("%s@%s: %s: %s \n",
					pGL.Att, pGL.Tag, pGL.AddressMode, FU.Tilded(pGL.AbsFP.S()))
			}
			println("#### GLink KEY TARGETS ####")
			for _, pGL := range AllGLinks.IncomableKeys {
				t := pGL.Tag
				a := pGL.Att
				// if S.HasPrefix(t, "topi") || S.Contains(a, "key") || S.Contains(a, "ref") {
				fmt.Printf("%s@%s: %s: %s \n",
					a, t, pGL.AddressMode, FU.Tilded(pGL.AbsFP.S()))
				// }
			}

			println("#### GLink URI SOURCES ####")
			for _, pGL := range AllGLinks.OutgoingURIs {
				fmt.Printf("%s@%s: %s: %s \n",
					pGL.Att, pGL.Tag, pGL.AddressMode, FU.Tilded(pGL.AbsFP.S()))
			}
			println("#### GLink URI TARGETS ####")
			for _, pGL := range AllGLinks.IncomableURIs {
				t := pGL.Tag
				a := pGL.Att
				isList := (len(pGL.Tag) == 2 && S.HasSuffix(pGL.Tag, "l"))
				if !isList {
					fmt.Printf("%s@%s: %s: %s \n",
						a, t, pGL.AddressMode, FU.Tilded(pGL.AbsFP.S()))
				}
			}

		* /

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

		if /* pCA.Validate && * / env.XmlCatalogFile != nil {

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

		if env.Pritt {
			// To name output file, insert "fmtd" btwn file base name and filext.
			filext := FP.Ext(env.Infile.AbsFP())
			fnbase := S.TrimSuffix(env.Infile.AbsFP(), filext)
			newname := fnbase + ".fmtd" + filext
			println("Fmtd file:", newname)
			fd, e := os.Create(newname)
			if e != nil {
				fmt.Printf("==> Cannot open fmtd file <%s>: %s \n",
					FU.Enhomed(newname), e.Error())
			} else {
				println("==> Created fmtd file:", FU.Enhomed(newname))
				env.PrittOutput = fd
			}
		}
	*/

	// ===============================
	//   4. LOAD FILES INTO DATABASE
	//      See also above
	// ===============================

	var haveDB bool = (env.SimpleRepo != nil)
	L.L.Dbg("env.SimpleRepo: <%T> %#v", env.SimpleRepo, env.SimpleRepo)
	// var jout []byte
	// var jerr error
	jout, jerr = json.MarshalIndent(env.SimpleRepo, "SimpleRepo: ", "  ")
	if jerr != nil {
		println(jerr)
		panic(jerr)
	}
	L.L.Info("JSON! " + string(jout))
	println("JSON! " + string(jout))

	var pSR *DRS.SqliteRepo
	var ok bool
	pSR, ok = env.SimpleRepo.(*DRS.SqliteRepo)
	if !ok {
		panic("Exec: repo is not *SimpleSqliteRepo")
	}
	if pSR == nil {
		L.L.Error("Cannot proceed: SqliteRepo is not valid")
		os.Exit(1)
	}
	// if haveDB && env.cfg.b.DBdoImport {
	if env.cfg.b.DBdoImport {
	   	var batchIndex int
		if !haveDB {
			panic("Exec: doImport but not haveDB")
		}
		var err error
		L.L.Progress("Will loadicontent for import batch, ID:%d",
			batchIndex)
		// var pTx *sql.Tx
		// =====================
		//  START A TRANSACTION
		// =====================
		err = env.SimpleRepo.Begin()
		if err != nil {
			L.L.Error("cli.exec.BeginTx: %w", err)
		}
		L.L.Info("TRANSACTION IS STARTED")
		var timeNow = time.Now().UTC().Format(time.RFC3339)
		// ============================
		//  FOR EVERY Input Contentity
		// ============================
		for _, pMCF := range InputContentities {
			// Prepare a DB record for the File
			pMCF.T_Imp = timeNow
			// L.L.Info("exec.L463: Trying new INSERT Generic")
			/*
				var stmt string
				// OBS stmt, e = pSR.NewInsertStmt(&pMCF.ContentityRow)
				stmt, e = DRS.NewInsertStmtGnrcFunc(
				      pSR, &pMCF.ContentityRow)
				if e != nil {
					return mcfile.WrapAsContentityError(e,
					  "new insert contentity stmt (cli.exec)", pMCF)
				}
				var insID int
				insID, e = pSR.ExecInsertStmt(stmt)
			*/
			var insID int
			insID, e = DRS.DoInsertGeneric(pSR, &pMCF.ContentityRow)
			if e != nil {
				return mcfile.WrapAsContentityError(e,
					"insert contentity to DB (cli.exec)", pMCF)
			}
			L.L.Info("Added file to import batch, ID: %d", insID)
		}
		pIB := new(DRM.InbatchRow)
		pIB.FilCt = len(InputContentities)
		pIB.Descr = "CLI import"
		// pIB.RelFP =
		// pIB.AnsFP =
		pIB.T_Cre, pIB.T_Imp = timeNow, timeNow
		/*
			var stmt string
			stmt, e = pSR.NewInsertStmt(pIB)
			if e != nil {
				return fmt.Errorf("new insert inbatch stmt (cli.exec): %w", e)
			}
			var insID int
			insID, e = pSR.ExecInsertStmt(stmt)
		*/
		var insID int
		insID, e = DRS.DoInsertGeneric(pSR, pIB)

		if e != nil {
			return fmt.Errorf("new insert inbatch to DB (cli.exec): %w", e)
		}
		L.L.Okay("cli/exec: INSERT'd inbatch OK, ID: %d", insID)

		e = pSR.Commit()
		if e != nil {
			return mcfile.WrapAsContentityError(e,
				"commit txn to DB failed (cli.exec)", nil)
		}
		L.L.Okay("Batch imported OK: TRANSACTION SUCCEEDED")
		// env.SimpleRepo.Tx = nil
		// pp := pCA.SimpleRepo.GetFileAll()
		// fmt.Printf("    DD:Files len %d id[last] %d \n", len(pp), fileIndex)
	}
	L.L.Info("TRYING SELECT BY ID")
	stmtS, eS := pSR.NewSelectByIdStmt(&DRM.TableDetailsCNT, 1)
	if eS != nil {
		return fmt.Errorf("new select contentity by id=1 stmt (cli.exec): %w", eS)
	}
	result, e3 := DRS.ExecSelectOneStmt[*DRM.ContentityRow](pSR, stmtS)
	if e3 != nil {
		return fmt.Errorf("new select contentity by id=1 from DB (cli.exec): %w", e3)
	}
	L.L.Warning("exec.go: INSERT'd inbatch OK, ID:%d", result)
	pSR.CloseLogWriter()
	return nil
}

func prerr(e error) string {
	if e == nil {
		return "-"
	}
	return e.Error()
}
