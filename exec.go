package m5cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"database/sql"
	LU "github.com/fbaube/logutils"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog" // Brings in global var L
	RS "github.com/fbaube/reposqlite"
	SU "github.com/fbaube/stringutils"
	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"
)

var LOG_LEVEL_FILE_READING = LU.LevelDbg
var LOG_LEVEL_EXEC_STAGES = LU.LevelDbg
var LOG_LEVEL_REF_LINKING = LU.LevelDbg
var LOG_LEVEL_WEB = LU.LevelDbg

// The general approach:
// 1. Filename via cmd line (can be Rel.FP)
// 2. Filename absolute path  (i.e. Abs.FP)
// 3. PathProps
// 4. ContentityRecord
// 5. Contentity
// 5.5 GTokens
// 6. GTree
// 7. ForesTree

func (env *XmlAppEnv) Exec() error {
	var e error
	// println("==> Exec: starting")
	// Timing: // tt := MU.Into("Input file processing")
	defer func() {
		L.L.Flush()
		/* ; fmt.Printf("%s: done.", os.Args[0]) */
	}()

	// ==========================
	//  1. IDENTIFY INPUT ITEMS
	// ==========================

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
	// Note that if a wildcard is used, all the files
	// in the expansion appear here.
	//
	// Indirs []FU.PathProps :: is all the directories
	// that were specified individually on the CLI.
	//
	// Inexpandirs []Expandir (still empty at this point)
	// :: is for all the entries in Indirs, to be expanded
	// into ContentityFS's and then flattened into slices.
	//
	// reference:
	//   type ExpanDir struct {
	//	   DirCt, FileCt, ItemCt  int
	//     CtyFS *mcfile.ContentityFS  }

	var pExpdDir *ExpanDir
	for iDir, pDir := range env.Indirs {
		L.L.Info("InDir[%d]: %s", iDir,
			SU.Tildotted(pDir.AbsFP.S()))
		pExpdDir = new(ExpanDir)
		pExpdDir.CtFS = mcfile.NewContentityFS(pDir.AbsFP.S(), nil)
		pExpdDir.ItemCt = pExpdDir.CtFS.Size()
		pExpdDir.FileCt = pExpdDir.CtFS.FileCount()
		pExpdDir.DirCt = pExpdDir.CtFS.DirCount()
		L.L.Okay("Found %d item(s) total (%d dirs, %d files)",
			pExpdDir.ItemCt, pExpdDir.DirCt, pExpdDir.FileCt)
		// InputFileSet.FilterInBySuffix(inputExts)
		// fmt.Printf("==> Found %d input file(s) " +
		//       "(after filtering) \n", nFiles)
		if pExpdDir.FileCt == 0 {
			L.L.Info("No content inputs to process! Aborting!")
			L.L.Close()
			// os.Exit(0)
			return errors.New("No inputs to process")
		}
		env.Inexpandirs = append(env.Inexpandirs, *pExpdDir)

		L.L.Warning("Skip'd wrtg out tree rep: [%d]", iDir)
		/*
			var f *os.File
			var filename string
			// Now write out a tree representation
			filename = fmt.Sprintf("./input-tree-%d", iDir)
			f, e = os.Create(filename)
			if e != nil {
				// An error here does not need to be fatal
				L.L.Error(filename + ": " + e.Error())
			} else {
				L.L.Progress("Writing input tree to " + filename + " ...")
				pExpdDir.CtFS.RootContentity().PrintTree(f)
				f.Close()
				L.L.Okay("Wrote input tree: " + filename)
			}
		*/
		// Now write out a css-enabled tree for html
		L.L.Warning("SKIP'D: css-enabled tree for html, exec.go.L107")
		/*
			filename = fmt.Sprintf("./css-tree", iDir)
			f, e = os.Create(filename)
			if e != nil {
				// An error here does not need to be fatal
				L.L.Error(filename + ": " + e.Error())
			} else {
				L.L.Progress("Writing css tree to " + filename + " ...")
				pExpdDir.CtFS.RootContentity().PrintTree(f)
				pExpdDir.CtFS.RootContentity().PrintCssTree(f)
				f.Close()
				L.L.Okay("Wrote css tree: " + filename)
			}
		*/
	}

	// =============================
	//  2. FOR EVERY CLI INPUT ITEM
	// =============================
	var InputCtysSlice []*mcfile.Contentity

	// =============================
	// 2a. FOR EVERY CLI INPUT FILE
	//     Make a new Contentity
	// =============================
	for i, pIF := range env.Infiles {
		var pCty *mcfile.Contentity
		L.L.Info("Input item [%02d] %s",
			i, SU.ElideHomeDir(pIF.AbsFP.S()))

		pCty, e = mcfile.NewContentity(pIF.AbsFP.S())
		if pCty == nil || e != nil || pCty.HasError() {
			if e == nil {
				e = errors.New("placeholder error")
			}
			L.L.Error("exec: newcontentity<%s>: %s",
				pIF.AbsFP, e.Error())
			continue
		}
		L.L.Info("INfile[%d]: [%d] %s %s",
			i, len(pCty.PathProps.Raw),
			pCty.MarkupType(), pCty.AbsFP())
		if pCty.MarkupType() == "UNK" {
			panic("UNK MarkupType in ExecuteStages")
		}
		InputCtysSlice = append(InputCtysSlice, pCty)
	}
	// ================================
	// 2b. FOR EVERY CLI INPUT DIR
	//     Expand it into files, which
	//     also makes new Contentities
	// ================================
	for _, pED := range env.Inexpandirs {
		InputCtysSlice = append(InputCtysSlice, pED.CtFS.AsSlice()...)
	}

	// Now we have all the inputs.
	// TODO: We could count up and tell the user
	// how many files of each valid extension.

	// =========================
	//  3. FOR EVERY CONTENTITY
	//     Prepare outputs
	// =========================
	for ii, cty := range InputCtysSlice {
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
		// os.Create(s) creates or truncates the named file.
		// If the file already exists, it is truncated.
		// If it does not exist, it is created with mode 0666
		// (before umask). If successful, the returned File
		// (IS OPEN, apparently, AND) can be used for I/O;
		// the associated file descriptor has mode O_RDWR.
		// If there is an error, it will be of type *PathError.

		cty.GTknsWriter = io.Discard
		cty.GTreeWriter = io.Discard
		cty.GEchoWriter = io.Discard

		if env.cfg.b.TotalTextal {
			fn := cty.AbsFP()
			if len(InputCtysSlice) == 1 {
				fn = "./debug"
			}
			// To name output files, append
			// "_(echo,tkns,tree)" to the entire file name.
			echoName := (fn + "_echo")
			tknsName := (fn + "_tkns")
			treeName := (fn + "_tree")
			fd1, e1 := os.Create(echoName)
			fd2, e2 := os.Create(tknsName)
			fd3, e3 := os.Create(treeName)
			if e1 == nil && e2 == nil && e3 == nil {
				L.L.Okay("created %s _echo,_tkns,_tree",
					SU.ElideHomeDir(fn))
			} else {
				L.L.Error("Cannot open a Total Textual file")
			}
			cty.GEchoWriter = fd1
			cty.GTknsWriter = fd2
			cty.GTreeWriter = fd3
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

	L.L.Info("Input contentities: total %d", len(InputCtysSlice))

	for ii, cty := range InputCtysSlice {
		if cty.IsDir() {
			continue
		}
		if cty.MarkupType() == "UNK" {
			panic("UNK MarkupType in ExecuteStages")
		}

		L.L.SetCategory(fmt.Sprintf("%02d", ii))
		L.L.Info(SU.Gbg("[F%02d] === %s (%d) ==="), ii,
			SU.ElideHomeDir(cty.AbsFP()), len(cty.PathProps.Raw))
		cty.ExecuteStages()
	}

	L.L.Info("Accumulated errors:")
	for _, p := range InputCtysSlice {
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

	var usingDB bool = (env.SimpleRepo != nil)
	var batchIndex int
	if usingDB && env.cfg.b.DBdoImport {
		var contIndex int
		var err error
		var inTx bool
		L.L.Progress("Loading content for import batch, ID:%d",
			batchIndex)
		var pTx *sql.Tx
		// pTx = env.SimpleRepo.Handle().MustBegin()
		pTx, err = env.SimpleRepo.Handle().Begin()
		if err != nil {
			// panic(err)
			fmt.Println("cli.exec.beginTx: %s", err.Error())
			L.L.Error("cli.exec.beginTx: %s", err.Error())
		} else {
			inTx = true
		}
		pSQR, ok := env.SimpleRepo.(*RS.SqliteRepo)
		if !ok {
			panic("NOT SQLITE")
		}
		for _, pMCF := range InputCtysSlice {
			// Prepare a DB record for the File
			pMCF.T_Imp = time.Now().UTC().Format(time.RFC3339)
			contIndex, e = pSQR.InsertContentityRow(&pMCF.ContentityRow) //,pTx)
			if e != nil {
				return mcfile.WrapAsContentityError(
					e, "insert contentity to DB (cli.exec)", pMCF)
			}
			L.L.Info("Added file to import batch, ID: %d", contIndex)
		}
		if inTx {
			e = pTx.Commit()
			if e != nil {
				return mcfile.WrapAsContentityError(e,
					"commit txn to DB failed (cli.exec)", nil)
			}
			L.L.Okay("Batch imported OK: TRANSACTION SUCCEEDED")
		}
		// env.SimpleRepo.Tx = nil
		// pp := pCA.SimpleRepo.GetFileAll()
		// fmt.Printf("    DD:Files len %d id[last] %d \n", len(pp), fileIndex)
	}
	return nil
}
