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
     	 path = p.FPs.AbsFP
	 pC, e = mcfile.NewContentity(path)
	 // We know that [NewContentity] returns exactly one nil ptr, so...
	 if pC == nil {
		eC = &os.PathError{Op:"LoadFilepathsContents.NewContentity",
		     Err:e,Path:fmt.Sprintf("[%d]:",i)+path} 
		ee = append(ee, eC)
		L.L.Error("LoadFileOops, nil Cty, %s", path)
		continue
	 }
	 if pC.MarkupType() == SU.MU_type_DIRLIKE {
	    L.L.Warning("LoadFilepathsContents: DIRLIKE: " + path)
	 }
	 if pC.MarkupType() == SU.MU_type_UNK {
		eC = &os.PathError{Op:"exec.loadFPs",
		    Err:errors.New("MarkupType is UNK"),Path:path}
		ee = append(ee, eC)
		L.L.Error("LoadFileOops, unk MarkupType, %s", path)
                continue
	 }
	 pCC = append(pCC, pC)
	 L.L.Okay("Item OK: MType<%s> MarkupType<%s>",
	 	pC.MType, pC.MarkupType())
	}
	return pCC, ee 
}