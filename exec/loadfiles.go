package exec

import(
	   "os"
	   "errors"
	   "github.com/fbaube/cnty"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
 	 L "github.com/fbaube/mlog"
)

// LoadFSOsIntoContentities turns a slice of [FSObject] into
// a slice of [Contentity]. Individual errors are returned
// via embedded struct [Errer]. The input and output slices 
// are the same length, for a one-to-one mapping. 
// .
func LoadFSOsIntoContentities(inFSOs []FU.FSObject) ([]*cnty.Contentity) {
     if inFSOs == nil || len(inFSOs) == 0 {
     	L.L.Info("LoadFSOsIntoContentities: no filepaths to load")
     	return nil 
	}
     var Cntys []*cnty.Contentity
     var pCnty *cnty.Contentity
     var path  string

     // For every input FSObject
     for iFso, pFso := range inFSOs {
     	 // If the FSO already has an error, skip it.
	 if pFso.HasError() {
            Cntys = append(Cntys, nil) // cnty.NewContentity(pFso.FPs.CreatPath()))
            continue
	 }
	 // -------------------
	 //  Prepare variebles
	 // -------------------
     	 // Use Rel.FP here, not Abs.FP, cos of
	 // use of stdlib when checking path. 
	 // 2026.05 Change to Abs.FP as more reliable.
     	 path = pFso.FPs.AbsFP
	 L.L.Info("InFile[%d]: %s", iFso, path)
	 pPE := new(os.PathError{Path:pFso.FPs.CreatPath()})
	 // --------
         //  Create
         // --------
	 pCnty = cnty.NewContentity(pFso.FPs.RelFP)
	 // Error? 
	 if pCnty.HasError() {
	    	pPE.Op = "loadfiles:newconty"
		pPE.Err = pCnty.GetError()
		pCnty.SetError(pPE) 
		L.L.Error("InFile[%d](%s) error: %s", iFso, path, pPE) // .Error())
		Cntys = append(Cntys, nil) 
		continue
	 }
	 // -----------------------------------
         //  Now that the Contentity has been 
         //  created, it has its own valid FSO. 
         // -----------------------------------
	 if pCnty.RawType() == SU.Raw_type_DIRLIKE {
	    L.L.Warning("LoadFiles: DIRLIKE: " + path)
	 }
	 if pCnty.RawType() == "" { // or SU.MU_type_UNK {
	    	pPE.Op = "exec:loadfiles"
		pPE.Err = errors.New("RawType is UNK")
                pCnty.SetError(pPE)
		L.L.Error("LoadFileOops, unk RawType, %s", path)
                continue
	 }
	 Cntys = append(Cntys, pCnty)
	 L.L.Okay("File OK: MType<%s> RawType<%s>",
	 	pCnty.MType, pCnty.RawType())
	}
	return Cntys 
}
