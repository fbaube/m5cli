package m5cli

import (
	// "encoding/json"
	"github.com/fbaube/m5cli/exec"
	L "github.com/fbaube/mlog" // Bring in global var L
	SU "github.com/fbaube/stringutils"
	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"
)

// *InputPathItems
// func file_reading_01(env *XmlAppEnv) error {
func file_reading_01(pIPI *InputPathItems) error {

	// At this point, "env" has three slices
	// of variables related  to input files:
	//
	// Infiles []FU.FSItem :: is all the files that were
	// specified individually on the CLI. Note that if
	// a wildcard was used, unquoted, then all files in
	// the expansion appear individually here.
	//
	// Indirs []FU.FSItem :: is all the directories
	// that were specified individually on the CLI.
	//
	// IndirFSs []ContentityFS (still empty at this point)
	// :: this maps to Indirs, making a ContentityFS for each
	// Indir, and then later on, each is flattened into a slice.

	// =======================
	// =======================
	// TOP LEVEL: FILE READING
	// =======================
	// =======================
	L.SetMaxLevel(LOG_LEVEL_FILE_READING)
	// ========================================
	//  EVERY CLI INPUT ITEM IS COLLECTED HERE
	//  First all files named on the command
	//  line, then all directories named there
	// ========================================
	// DUMP pIPI.NamedDirrs, Inexpandirs
	L.L.Info("AppEnv.NamedFiles: [%d]: %+v \n",
		len(pIPI.NamedFiles), pIPI.NamedFiles)
	L.L.Info("AppEnv.NamedDirrs:: [%d]: %+v \n",
		len(pIPI.NamedDirrs), pIPI.NamedDirrs)
	/*
	if env.cfg.b.Samples {
		// ALSO DUMP AS JSON
		var jout []byte
		var jerr error
		if len(pIPI.NamedFiles) > 0 {
			jout, jerr = json.MarshalIndent(
				pIPI.NamedFiles[0], "infile: ", "  ")
			if jerr != nil {
				println(jerr)
				panic(jerr)
			}
			L.L.Debug("JSON! " + string(jout))
		}
		if len(pIPI.NamedDirrs) > 0 {
			jout, jerr = json.MarshalIndent(
				pIPI.NamedDirrs[0], "indirr: ", "  ")
			if jerr != nil {
				println(jerr)
				panic(jerr)
			}
			L.L.Debug("JSON! " + string(jout))
		}
	}
	*/
	// fmt.Printf("==> pIPI.Inexpandirs: %#v \n", pIPI.Inexpandirs)

	// ==========================
	//  FOR EVERY CLI INPUT FILE
	//  Make a new Contentity
	// ==========================
	// var InfileContentities []*mcfile.Contentity   // directories
	// var IndirContentityFSs []*mcfile.ContentityFS // trees

	L.L.Warning(SU.Rfg(SU.Ybg("=== LOAD CLI FILE(S) ===")))
	// fmt.Fprintf(os.Stderr, "exec: pIPI.NamedFiles: %#v \n", pIPI.NamedFiles)
	// fmt.Fprintf(os.Stderr, "exec: pIPI.NamedFiles[0]: %#v \n", *pIPI.NamedFiles[0].FPs)
	var errct int 
	pIPI.AllCntys, errct = exec.LoadFilepathsContentities(pIPI.NamedFiles)
	gotCtys := pIPI.AllCntys != nil && len(pIPI.AllCntys) > 0
	if gotCtys {
		L.L.Okay("Results for %d infiles: %d OK, %d not OK \n",
			len(pIPI.NamedFiles), len(pIPI.AllCntys)-errct, errct)
		for i, pC := range pIPI.AllCntys {
		        if !pC.HasError() {
			   L.L.Okay("InFile[%02d] len:%d RawTp:%s : %s",
				i, len(pC.FSItem.Raw), pC.RawType(),
				pC.FSItem.FPs.ShortFP)
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
	}
	L.L.Info("Loaded %d file contentity/ies", len(pIPI.AllCntys))
	// ==================================
	//   FOR EVERY CLI INPUT DIRECTORY
	//  Make a new Contentity filesystem
	// ==================================
	L.L.Warning(SU.Rfg(SU.Ybg("=== EXPAND CLI DIR(S) ===")))
	pIPI.DirCntyFSs = exec.LoadDirpathsContentFSs(pIPI.NamedDirrs)
	WriteContentityFStreeFiles(pIPI.DirCntyFSs)
	L.L.Info("Expanded %d file folder(s) into %d F/S(s)",
		len(pIPI.NamedDirrs), len(pIPI.DirCntyFSs))

	// ==============================
	//  FOR EVERY CLI INPUT DIRECTORY
	//  Expand it into files, which
	//  also makes new Contentities
	// ==============================
	L.L.Warning(SU.Rfg(SU.Ybg("=== LOAD CLI DIR(S) ===")))
	for _, pED := range pIPI.DirCntyFSs {
		pIPI.AllCntys = append(pIPI.AllCntys, pED.AsSlice()...)
	}
	L.L.Info("Expanded %d F/S(s), now have %d contentities",
		len(pIPI.DirCntyFSs), len(pIPI.AllCntys))

	// Now we have all the inputs.
	// TODO: We could count up and tell the user
	// how many files of each valid extension.

	// =======================
	//  SUMMARIZE TO THE USER
	//  ALL CONTENTITIES THAT
	//  ARE LOADED & READY
	// =======================
	for ii, cty := range pIPI.AllCntys {
		if cty == nil {
			L.L.Okay("[%02d]  nil", ii)
		} else if cty.IsDir() {
			L.L.Okay("[%02d]  DIR \t\t%s", ii, cty.FPs.ShortFP)
		} else {
			mt := cty.MType
			if mt == "" {
				mt = "(nil MType)"
			}
			L.L.Okay("[%02d]  %s \t%s", ii, mt, cty.FSItem.FPs.ShortFP)
		}
	}
	return nil
}