package m5cli

import (
	"encoding/json"
	"github.com/fbaube/m5cli/exec"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog" // Bring in global var L
	SU "github.com/fbaube/stringutils"
	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"
)

func file_reading_01(env *XmlAppEnv) error {

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
	var InfileContentities []*mcfile.Contentity   // directories
	var IndirContentityFSs []*mcfile.ContentityFS // trees
	var ee []error

	L.L.Warning(SU.Rfg(SU.Ybg("=== LOAD CLI FILE(S) ===")))
	// fmt.Fprintf(os.Stderr, "exec: env.Infiles: %#v \n", env.Infiles)
	// fmt.Fprintf(os.Stderr, "exec: env.Infiles[0]: %#v \n", *env.Infiles[0].FPs)
	InfileContentities, ee = exec.LoadFilepathsContentities(env.Infiles)
	gotCtys := InfileContentities != nil && len(InfileContentities) > 0
	gotErrs := ee != nil && len(ee) > 0
	if gotCtys || gotErrs {
		L.L.Okay("Results for %d infiles: %d OK, %d not OK \n",
			len(env.Infiles), len(InfileContentities), len(ee))
		for i, pC := range InfileContentities {
			L.L.Okay("InFile[%02d] len:%d RawTp:%s : %s",
				i, len(pC.FSItem.Raw), pC.RawType(),
				pC.FSItem.FPs.ShortFP)
			/* if pCty.RawType() == SU.Raw_type_UNK ||
			      pCty.RawType() ==  "" { {
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
	WriteContentityFStreeFiles(IndirContentityFSs)
	L.L.Info("Expanded %d file folder(s) into %d F/S(s)",
		len(env.Indirs), len(IndirContentityFSs))

	// ==============================
	//  FOR EVERY CLI INPUT DIRECTORY
	//  Expand it into files, which
	//  also makes new Contentities
	// ==============================
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
	InitContentityDebugFiles(InfileContentities, env.cfg.b.TotalTextal)

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
			if mt == "" {
				mt = "(nil MType)"
			}
			L.L.Okay("[%02d]  %s \t%s", ii, mt, cty.FSItem.FPs.ShortFP)
		}
	}
	return nil
}