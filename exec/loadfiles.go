package exec

import(
	"fmt"
	"errors"
	"os"
	"github.com/fbaube/mcfile"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	L "github.com/fbaube/mlog"
)

// LoadFilepathsContents can return a nil or empty second return value.
// The items in the two arrays do not correspond. The path sets are disjoint.
func LoadFilepathsContents(inFSIs []FU.FSItem) ([]*mcfile.Contentity, []error) {

     if inFSIs == nil || len(inFSIs) == 0 {
     	L.L.Info("No filepaths to load")
     	return nil, nil
	}
     var pCC []*mcfile.Contentity
     var pC    *mcfile.Contentity
     var ee  []error
     var e, eC error
     var path string 

     // For every input FSItem
     for i, p := range inFSIs {
     	 // Use Rel.FP here, not Abs.FP, cos of
	 // use of std lib when checking path 
     	 path = p.FPs.RelFP // AbsFP
	 // println("LoadFiles: mcfile.NewContentity:", path)
	 pC, e = mcfile.NewContentity(path)
	 // We know that [NewContentity] returns exactly one nil ptr, so...
	 if pC == nil {
		eC = &fs.PathError{Op:"LoadFilepathsContents.NewContentity",
		     Err:e,Path:fmt.Sprintf("[%d]:",i)+path} 
		ee = append(ee, eC)
		L.L.Error("LoadFileOops, nil Cty, %s", path)
		continue
	 }
	 if pC.RawType() == SU.Raw_type_DIRLIKE {
	    L.L.Warning("LoadFilepathsContents: DIRLIKE: " + path)
	 }
	 if pC.RawType() == "" { // or SU.MU_type_UNK {
		eC = &fs.PathError{Op:"exec.loadFPs",
		    Err:errors.New("RawType is UNK"),Path:path}
		ee = append(ee, eC)
		L.L.Error("LoadFileOops, unk RawType, %s", path)
                continue
	 }
	 pCC = append(pCC, pC)
	 L.L.Okay("Item OK: MType<%s> RawType<%s>", pC.MType, pC.RawType())
	}
	return pCC, ee 
}