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
// via embedded struct Errer.
//
// FIXME: Return second value = nErrors
// .
func LoadFilepathsContentities(inFSIs []FU.FSItem) []*mcfile.Contentity {
     if inFSIs == nil || len(inFSIs) == 0 {
     	L.L.Info("No filepaths to load")
     	return nil 
	}
     var pCC []*mcfile.Contentity
     var pC    *mcfile.Contentity
     var e, eC error
     var path string 

     // For every input FSItem
     for i, p := range inFSIs {
     	 // Use Rel.FP here, not Abs.FP, cos of
	 // use of std lib when checking path 
     	 path = p.FPs.RelFP // AbsFP
	 // println("LoadFiles: mcfile.NewContentity:", path)
	 // FIXME: We should use the FSItem to create the Contentity. 
	 pC, e = mcfile.NewContentity(path)
	 // We know that [NewContentity] returns exactly one nil ptr, so...
	 // if pC == nil {
	 if e != nil {
		eC = &fs.PathError{Op:"LoadFilepathsContents.NewContentity",
		     Err:e,Path:fmt.Sprintf("[%d]:",i)+path}
		if pC == nil { pC = &(mcfile.Contentity{}) }
		pC.SetError(eC.Error())
		L.L.Error("LoadFileOops: %s: %s", pC.Error(), path)
		continue
	 }
	 if pC.RawType() == SU.Raw_type_DIRLIKE {
	    L.L.Warning("LoadFilepathsContents: DIRLIKE: " + path)
	 }
	 if pC.RawType() == "" { // or SU.MU_type_UNK {
		eC = &fs.PathError{Op:"exec.loadFPs",
		    Err:errors.New("RawType is UNK"),Path:path}
		if pC == nil { pC = &(mcfile.Contentity{}) }
                pC.SetError(eC.Error())
		L.L.Error("LoadFileOops, unk RawType, %s", path)
                continue
	 }
	 pCC = append(pCC, pC)
	 L.L.Okay("Item OK: MType<%s> RawType<%s>", pC.MType, pC.RawType())
	}
	return pCC 
}