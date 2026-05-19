package exec

import(
	"fmt"
	"errors"
	"io/fs"
	"github.com/fbaube/cnty"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	L "github.com/fbaube/mlog"
)

// LoadFilepathsContentities turns a slice of [FSObject] into
// a slice of [Contentity]. Individual errors are returned
// via embedded struct Errer, but for convenience, a summary
// count of errors is the second return value. 
// .
func LoadFilepathsContentities(inFSOs []FU.FSObject) ([]*cnty.Contentity, int) {
     if inFSOs == nil || len(inFSOs) == 0 {
     	L.L.Info("No filepaths to load")
     	return make([]*cnty.Contentity, 0), 0
	}
     var pCC []*cnty.Contentity
     var pC    *cnty.Contentity
     var eC    error
     var errct int 
     var path  string

     // For every input FSObject
     for i, fso := range inFSOs {
     	 // If the FSO already has an error, skip it.
	 if fso.HasError() {
	    errct++
	    continue
	 }
     	 // Use Rel.FP here, not Abs.FP, cos of
	 // use of std lib when checking path 
     	 path = fso.FPs.RelFP // AbsFP
	 // println("LoadFiles: cnty.NewContentity:", path)
	 // FIXME: Contentity contains a ContentityRecord contains an
	 // FSObject, so DUH we should use the FSObject to create the 
	 // Contentity. But the Contentity also contains a Nord, so 
	 // it gets complicated. So don't worry about this too much. 
	 pC = cnty.NewContentity(path)
	 if pC.HasError() {
		eC = &fs.PathError{Op:"loadfilepathscontents.newcontentity",
		     Err:pC.GetError(),Path:fmt.Sprintf("[%d]:",i)+path}
		// if pC == nil { pC = &(cnty.Contentity{}) }
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
		if pC == nil { pC = &(cnty.Contentity{}) }
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
