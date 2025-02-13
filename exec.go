package m5cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"cmp"

	DRS "github.com/fbaube/datarepo/sqlite"
	"github.com/fbaube/m5cli/exec"
	"github.com/fbaube/m5db"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog" // Brings in global var L
	SU "github.com/fbaube/stringutils"

	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"
)

// The general approach:
//  1) os.Args
//  2) AppCfg
//  3) AppEnv
//  4) Run()!
//
// How refs to files (and directories) enter the
// system (altho levels might be intermixed here):
//  1) Filename via cmd line (can be Rel.FP)
//  2) Filename absolute path  (i.e. Abs.FP)
//  3) FSItem
//  4) Loading & parsing: GTokens
//  5) PathAnalysis (i.e. content analysis) 
//  6) ContentityRow
//  7) Contentity 
//  8) GTree
//  9) ContentiTree

// Exec does all execution of all stages for
// every [mcfile.Contentity], altho only after
// all prep has already been done by other funcs.
// The key input is an [*XmlAppEnv].
// .
func (env *XmlAppEnv) Exec() error {

	L.SetMaxLevel(LOG_LEVEL_FILE_INTRO)
	var e error
	defer func() { L.L.Flush() }()
	// Timing:
	// tt := MU.Into("Input file processing")
	
	// ===================
	//  PRELIMINARY STUFF
	// ===================
	// TODO: if bool CLI flag Samples 
	if env.cfg.b.Samples { 
	     // Dump out what a ContentityRow looks like in the DB
		var cntro = new(m5db.ContentityRow)
		var cptrs = m5db.ColumnPtrsFuncCNT(cntro, true)
		L.L.Info("ContentityRow datarepo/TableDetails:")
		L.L.Info("\t instance: cntRow<%T> colPtrs <%T>", cntro, cptrs)
		for iii, ppp := range cptrs {
		    L.L.Info("\t col[%d] <%T>", iii, ppp)
		}
	}

	// At this point, "env" has three slices
	// of variables related  to input files:
	//
	// Infiles []FU.FSItem :: is all the files 
	// that were specified individually on the CLI.
	// Note that if a wildcard was used, all files
	// in the expansion appear individually here.
	//
	// Indirs []FU.FSItem :: is all the directories
	// that were specified individually on the CLI.
	//
	// IndirFSs []ContentityFS (still empty at this point)
	// :: maps to Indirs, making a ContentityFS for each, 
	// then each is flattened into a slice.

	L.SetMaxLevel(LOG_LEVEL_FILE_READING)
	// ========================================
	//  EVERY CLI INPUT ITEM IS COLLECTED HERE
	//  First all files named at the command 
	//  line, then all directories named 
	// ========================================
	// DUMP env.Indirs, Inexpandirs
	L.L.Info("AppEnv.Infiles: [%d]: %+v \n", len(env.Infiles), env.Infiles)
	L.L.Info("AppEnv.Indirs:: [%d]: %+v \n", len(env.Indirs), env.Indirs)
	if env.cfg.b.Samples {
	   	// ALSO DUMP AS JSON
		var jout []byte
		var jerr error
		if len(env.Infiles) > 0 {
		   jout, jerr = json.MarshalIndent(
		   	 env.Infiles[0], "infile: ", "  ")
		   if jerr != nil {
			println(jerr)
			panic(jerr)
		   }
		L.L.Debug("JSON! " + string(jout))
		}
		if len(env.Indirs) > 0 {
		   jout, jerr = json.MarshalIndent(
		   	 env.Indirs[0], "indirr: ", "  ")
		   if jerr != nil {
			println(jerr)
			panic(jerr)
		   }
		L.L.Debug("JSON! " + string(jout))
		}
	}

	// fmt.Printf("==> env.Inexpandirs: %#v \n", env.Inexpandirs)

	// ==========================
	//  FOR EVERY CLI INPUT FILE
	//  Make a new Contentity
	// ==========================
	var InfileContentities []*mcfile.Contentity // directories 
	var IndirContentityFSs []*mcfile.ContentityFS // trees 
	var ee []error

	L.L.Warning(SU.Rfg(SU.Ybg("=== LOAD CLI FILE(S) ===")))
	// fmt.Fprintf(os.Stderr, "exec: env.Infiles: %#v \n", env.Infiles)
	// fmt.Fprintf(os.Stderr, "exec: env.Infiles[0]: %#v \n",
	// 	*env.Infiles[0].FPs)
	InfileContentities, ee = exec.LoadFilepathsContentities(env.Infiles)
	gotCtys := InfileContentities != nil && len(InfileContentities) > 0
	gotErrs := ee != nil && len(ee) > 0
	if gotCtys || gotErrs {
		L.L.Okay("Results for %d infiles: %d OK, %d not OK \n",
			len(env.Infiles), len(InfileContentities), len(ee))
		for i, pC := range InfileContentities {
			/* L.L.Info("InFile[%02d] OK! [%d] %s :: %s",
				i, len(pC.FSItem.Raw), pC.RawType(),
				pC.FSItem.FPs.ShortFP) */
			L.L.Okay("InFile[%02d] len:%d RawTp:%s : %s",
				i, len(pC.FSItem.Raw), pC.RawType(),
				pC.FSItem.FPs.ShortFP) 
			/* if pCty.RawType() == SU.Raw_type_UNK {
				s := fmt.Sprintf("INfile[%d]: [%d] %s %s",
			             i, len(pCty.PathProps.Raw),
			             pCty.RawType(), pCty.AbsFP())
			panic("UNK RawType in ExecuteStages; \n" + s) */
		}
		for i, eC := range ee {
			L.L.Error("InfileErr[%02d] ERR :: <%T> %s", i, eC, eC)
		}
	}
	L.L.Info("Loaded %d file contentity/ies", len(InfileContentities))
	// ==================================
	//   FOR EVERY CLI INPUT DIRECTORY
	//  Make a new Contentity filesystem
	// ==================================
	L.L.Warning(SU.Rfg(SU.Ybg("=== EXPAND CLI DIR(S) ===")))
	IndirContentityFSs = exec.LoadDirpathsContentFSs(env.Indirs)
	for iFS, pFS := range IndirContentityFSs {
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
		// L.L.Warning("SKIP'D: css-enabled tree for html, exec.go.L181")
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
	L.L.Info("Expanded %d file folder(s) into %d F/S(s)",
		  len(env.Indirs), len(IndirContentityFSs))

	// =============================
	//  FOR EVERY CLI INPUT DIR
	//  Expand it into files, which
	//  also makes new Contentities
	// =============================
	L.L.Warning(SU.Rfg(SU.Ybg("=== LOAD CLI DIR(S) ===")))
	for _, pED := range IndirContentityFSs {
		InfileContentities = append(InfileContentities, pED.AsSlice()...)
	}
	L.L.Info("Expanded %d F/S(s), now have %d contentities",
		  len(IndirContentityFSs), len(InfileContentities))

	// Now we have all the inputs.
	// TODO: We could count up and tell the user
	// how many files of each valid extension.

	// ======================
	//  FOR EVERY CONTENTITY
	//  Prepare outputs
	// ======================
	for ii, cty := range InfileContentities {
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
			if len(InfileContentities) == 1 {
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
	// =======================
	//  SUMMARIZE TO THE USER 
	//  ALL CONTENTITIES THAT
	//  ARE LOADED & READY
	// =======================
	for ii, cty := range InfileContentities {
	    if cty == nil {
	       L.L.Okay("[%02d]  nil", ii)
	    } else if cty.IsDir() {
	       L.L.Okay("[%02d]  DIR \t\t%s", ii, cty.FPs.ShortFP)
	    } else {
	       mt := cty.MType
	       if mt == "" { mt = "(nil MType)" } 
	       L.L.Okay("[%02d]  %s \t%s", ii, mt, cty.FSItem.FPs.ShortFP)
	    }
	}

	// ======================
	//  FOR EVERY CONTENTITY
	//  Execute all stages
	// ======================
	L.SetMaxLevel(LOG_LEVEL_EXEC_STAGES)
	L.L.Okay(SU.Rfg(SU.Ybg("=== DO CONTENTITY STAGES ===")))

	L.L.Info("Input contentities: total %d", len(InfileContentities))

	for ii, cty := range InfileContentities {
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
		if cty.FSItem.TypedRaw != nil { rawlen = len(cty.FSItem.Raw) }
		dsp = fmt.Sprintf(" %4d  %s  %s  %s",
		      rawlen, cty.RawType(),
		      cmp.Or(cty.MType, "(nil-MType)"),
		      cmp.Or(cty.MimeType, "(nil-Mime)")) 
		L.L.Info(SU.Cyanbg(SU.Wfg(dsp)))
		
		// Now try dumping it as JSON !
		// func MarshalIndent(v any, prefix,
		//    indent string) ([]byte, error)
		var jsonOut []byte 
		jsonOut, e = json.MarshalIndent(cty, "", "  ")
		/*
		L.L.Info("================\n" +
			 "%s \n" +
			 "================", string(jsonOut))
		*/
		jsonOutFilename := cty.FPs.AbsFP + ".tmp.json"
		errr := os.WriteFile(jsonOutFilename, jsonOut, 0644)
		if errr != nil { panic("os.WriteFile json: " + jsonOutFilename) }
		L.L.Info("Wrote JSON to: " + jsonOutFilename)
		// defer jsonOutFile.Close()
		// ==================
		//  AND NOW, EXECUTE
		// ==================
		cty.ExecuteStages()
	}

	L.L.Info("Accumulated errors:")
	for _, p := range InfileContentities {
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
			pSR /* env.SimpleRepo */ , InfileContentities)
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
