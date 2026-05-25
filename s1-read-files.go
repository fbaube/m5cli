package m5cli

import (
	// "encoding/json"
	"github.com/fbaube/m5cli/exec"
	L "github.com/fbaube/mlog" // Bring in global var L
	SU "github.com/fbaube/stringutils"
	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"
)

func step1_read_files(pIPO *InputPathObjex) error {

	// At this point, "env" has three slices
	// of variables related  to input files:
	//
	// Infiles []FU.FSO :: is all the files that were
	// specified individually on the CLI. Note that if
	// a wildcard was used, unquoted, then all files in
	// the expansion appear individually here.
	//
	// Indirs []FU.FSO :: is all the directories
	// that were specified individually on the CLI.
	//
	// IndirFSs []ContentityFS (still empty at this point)
	// :: this maps to Indirs, making a ContentityFS for each
	// Indir, and then later on, each is flattened into a slice.

	// =================================
	//  1) READ FILES INTO Contentities 
	// =================================
	L.SetMaxLevel(LOG_LEVEL_FILE_READING)
	// ==========================================
	//  EVERY CLI INPUT ITEM IS COLLECTED HERE.
	//  First all files named on the command line.
	//  Then all directories named there.
	// ==========================================
	// DUMP pIPO.NamedDirrs, Inexpandirs
	L.L.Info("AppEnv: %d named files, %d named dirs", 
		len(pIPO.NamedFiles), pIPO.NamedDirrs)
	/*
	if env.cfg.b.Samples {
		// ALSO DUMP AS JSON
		var jout []byte
		var jerr error
		if len(pIPO.NamedFiles) > 0 {
			jout, jerr = json.MarshalIndent(
				pIPO.NamedFiles[0], "infile: ", "  ")
			if jerr != nil {
				println(jerr)
				panic(jerr)
			}
			L.L.Debug("JSON! " + string(jout))
		}
		if len(pIPO.NamedDirrs) > 0 {
			jout, jerr = json.MarshalIndent(
				pIPO.NamedDirrs[0], "indirr: ", "  ")
			if jerr != nil {
				println(jerr)
				panic(jerr)
			}
			L.L.Debug("JSON! " + string(jout))
		}
	}
	*/
	// fmt.Printf("==> pIPO.Inexpandirs: %#v \n", pIPO.Inexpandirs)

	// =============================
	//  2) FOR EVERY CLI INPUT FILE
	//        Make a new Contentity
	// =============================
	L.L.Warning(SU.Rfg(SU.Ybg("=== (stg1) LOAD CLI FILE(S) ===")))
	var errct int
	// ----------------------------------
	//  This func does the heavy lifting 
	// ----------------------------------
	pIPO.AllCntys = exec.LoadFSOsIntoContentities(pIPO.NamedFiles)
	gotCtys := pIPO.AllCntys != nil && len(pIPO.AllCntys) > 0

	if gotCtys {
		L.L.Okay("Results for %d infiles: %d OK, %d not OK \n",
			len(pIPO.NamedFiles), len(pIPO.AllCntys)-errct, errct)
		for i, pC := range pIPO.AllCntys {
		        if !pC.HasError() {
			   L.L.Okay("InFile[%02d] len:%d RawTp:%s : %s",
				i, len(pC.FSO.Raw), pC.RawType(),
				pC.FSO.FPs.ShortFP)
			/* if pCty.RawType() == SU.Raw_type_UNK ||
			      pCty.RawType() ==  "" { {
				s := fmt.Sprintf("INfile[%d]: [%d] %s %s",
			             i, len(pCty.PathProps.Raw),
			             pCty.RawType(), pCty.AbsFP())
				panic("UNK RawType in ExecuteStages; \n" + s) */
			} else {
			  L.L.Error("InFile[%02d] ERROR: %s",
                                 i, pC.GetError())
			}
		}
		L.L.Info("Loaded %d file contentity(s)", len(pIPO.AllCntys))
	} else {
	        L.L.Info("Loaded NO file contentitys")
	}
	// ===================================
	//   3) FOR EVERY CLI INPUT DIRECTORY
	//     Make a new Contentity FS
	// ===================================
	L.L.Warning(SU.Rfg(SU.Ybg("=== (stg1) CLI DIR(S) form FSs ===")))
	pIPO.DirCntyFSs = exec.LoadFSOsIntoContentityFSs(pIPO.NamedDirrs)
	WriteContentityFStreeFiles(pIPO.DirCntyFSs)
	L.L.Info("Expanded %d file folder(s) into %d F/S(s)",
		len(pIPO.NamedDirrs), len(pIPO.DirCntyFSs))

	// ==================================
	//  4) FOR EVERY CLI INPUT DIRECTORY
	//    (and Contentity FS) Expand it
	//    into files, which also makes
	//    new Contentities
	// ==================================
	L.L.Warning(SU.Rfg(SU.Ybg("=== (stg1) LOAD CLI DIR(S) ===")))
	for _, pED := range pIPO.DirCntyFSs {
		pIPO.AllCntys = append(pIPO.AllCntys, pED.AsSlice()...)
	}
	L.L.Info("Expanded %d F/S(s), now have %d contentities",
		len(pIPO.DirCntyFSs), len(pIPO.AllCntys))

	// Now we have all the inputs.
	// TODO: We could count up and tell the user
	// how many files of each valid extension.

	// =======================
	//  SUMMARIZE TO THE USER
	//    ALL CONTENTITIES 
	// =======================
	for ii, cty := range pIPO.AllCntys {
	    	if cty.HasError() {
		   	L.L.Error("[%02d] %s", ii, cty.Error())
		} else if cty == nil {
			panic("NIL CNTY ?!") // L.L.Okay("[%02d]  nil", ii)
		} else if cty.FSO.FPs.IsDir {
			L.L.Okay("[%02d]  DIR \t\t%s", ii, cty.FSO.FPs.ShortFP)
		} else if cty.ContentAnalysis == nil {
			L.L.Okay("[%02d]  nilContentAnalysis \t%s",
				ii,  cty.FSO.FPs.ShortFP)
		} else { 
			mt := cty.MType
			if mt == "" {
				mt = "(nil MType)"
			}
			L.L.Okay("[%02d]  %s \t%s", ii, mt, cty.FSO.FPs.ShortFP)
		}
	}
	return nil
}