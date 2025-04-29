package m5cli

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	// "github.com/fbaube/m5cli/exec"
	"github.com/fbaube/mcfile"
	L "github.com/fbaube/mlog" // Bring in global var L
	SU "github.com/fbaube/stringutils"
	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"
)

func exec_stages_2(InfileContentities []*mcfile.Contentity) error {

	// =========================
	// =========================
	// TOP LEVEL: EXECUTE STAGES
	// =========================
	// =========================
	//  FOR EVERY CONTENTITY
	// ======================
	L.SetMaxLevel(LOG_LEVEL_EXEC_STAGES)
	L.L.Okay(SU.Rfg(SU.Ybg("=== DO CONTENTITY STAGES ===")))

	L.L.Info("Input contentities: total %d", len(InfileContentities))

	for ii, cty := range InfileContentities {
		// Skip directories entirely
		if cty.IsDir() {
			continue
		}
		// If we still have symlinks here, note it
		if cty.IsDirlike() {
		   	L.L.Warning("execstages: got symlink: " + cty.AbsFP())
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
		dsp = fmt.Sprintf(" %4d  %s  %s  %s",
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
		// ==============================================
		//     AND NOW, EXECUTE
		//   If we want to try a fancy execution model, 
		// such as lots of gogunc's, it will happen here. 
		// ==============================================
		cty.ExecuteStages()
	}

	// =========================
	//  DUMP ACCUMULATED ERRORS
	// =========================
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
		(for every InfileContentity)
		if env.GroupGenerated {
			// println("==> Creating output directories.")
			e = pMCF.ConfigureOutputFiles("_" + myAppName)
			errorbarf(e, "ConfigureOutputFiles")
			}
		}
	*/
	return nil
}