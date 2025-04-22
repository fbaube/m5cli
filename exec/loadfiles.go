package exec

import(
	"fmt"
	"errors"
	"io/fs"
	"github.com/fbaube/mcfile"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	L "github.com/fbaube/mlog"
)

// LoadFilepathsContentities turns a slice of [FSItem] into
// a slice of [Contentity]. Individual errors are returned
// via embedded struct Errer, but for convenience, a summary
// count of errors is the second return value. 
// .
func LoadFilepathsContentities(inFSIs []FU.FSItem) ([]*mcfile.Contentity, int) {
     if inFSIs == nil || len(inFSIs) == 0 {
     	L.L.Info("No filepaths to load")
     	return make([]*mcfile.Contentity, 0), 0
	}
     var pCC []*mcfile.Contentity
     var pC    *mcfile.Contentity
     var eC    error
     var errct int 
     var path  string

     // For every input FSItem
     for i, fsi := range inFSIs {
     	 // If the FSI already has an error, skip it.
	 if fsi.HasError() {
	    errct++
	    continue
	 }
     	 // Use Rel.FP here, not Abs.FP, cos of
	 // use of std lib when checking path 
     	 path = fsi.FPs.RelFP // AbsFP
	 // println("LoadFiles: mcfile.NewContentity:", path)
	 // FIXME: Contentity contains a ContentityRecord contains
	 // an FSItem, so DUH we should use the FSItem to create the 
	 // Contentity. But the Contentity also contains a Nord, so 
	 // it gets complicated. So don't worry about this too much. 
	 pC = mcfile.NewContentity(path)
	 if pC.HasError() {
		eC = &fs.PathError{Op:"loadfilepathscontents.newcontentity",
		     Err:pC.GetError(),Path:fmt.Sprintf("[%d]:",i)+path}
		// if pC == nil { pC = &(mcfile.Contentity{}) }
		pC.SetError(eC) 
		L.L.Error("LoadFileOops: %s: %s", pC.Error(), path)
		errct++
		continue
	 }
	 if pC.RawType() == SU.Raw_type_DIRLIKE {
	    L.L.Warning("LoadFilepathsContents: DIRLIKE: " + path)
	 }
	 if pC.RawType() == "" { // or SU.MU_type_UNK {
		eC = &fs.PathError{Op:"exec.loadFPs",
		    Err:errors.New("RawType is UNK"),Path:path}
		if pC == nil { pC = &(mcfile.Contentity{}) }
                pC.SetError(eC)
		L.L.Error("LoadFileOops, unk RawType, %s", path)
		errct++
                continue
	 }
	 pCC = append(pCC, pC)
	 L.L.Okay("Item OK: MType<%s> RawType<%s>", pC.MType, pC.RawType())
	}
	return pCC, errct
}