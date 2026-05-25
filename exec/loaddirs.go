package exec

// TODO: SHOULD USE
// path.Clean (rmvs trlg slashes) 
// fs.ValidPath
// FP.IsLocal (implies ValidPath, so do VP first)


import(
	   "fmt"
	   "github.com/fbaube/cnty"
	FU "github.com/fbaube/fileutils"
	 L "github.com/fbaube/mlog"

)

// LoadFSOsIntoContentityFSs turns a slice of [FSObject] into
// a slice of [ContentityFS]. Individual errors are returned
// via embedded struct [Errer]. The input and output slices
// are the same length, for a one-to-one mapping.
// .
func LoadFSOsIntoContentityFSs(inFSOs []*FU.FSObject) ([]*cnty.ContentityFS) {
     if inFSOs == nil || len(inFSOs) == 0 {
     	L.L.Info("LoadFSOsIntoContentityFSs: no filepaths to load")
     	return nil
	}
     var CFSs []*cnty.ContentityFS
     var pCFS   *cnty.ContentityFS
     var path  string

     //  For every input FSObject
     for iFso, pFso := range inFSOs {
     	 // If the FSO already has an error, copy it into an
	 // empty ContentityFS and skip further processing.
         if pFso.HasError() {
	 // tmp, _ := cnty.NewContentityFS(pFso.FPs.CreatPath(), nil)
	    CFSs = append(CFSs, nil) // tmp) 
            continue
         }
	 // -------------------
	 //  Prepare variebles
	 // -------------------
	 // AbsFP might be more reliable, but use 
	 // RelFPbecause we will be using [os.Root]. 
	 path = pFso.FPs.RelFP
	 L.L.Info("InDir[%d]: %s", iFso, path) // FIXME pFso.FPs.CreatPath()) 
	 var e error
	 // --------
	 //  Create
	 // --------
	 // nil is []string of OK file extensions 
	 pCFS, e = cnty.NewContentityFS(path, nil)
	 // Error?
	 if e != nil { 
	      	 fmt.Printf("InDir[%d](%s) error: %s", iFso, e.Error(), path)
		 L.L.Error ("InDir[%d](%s) error: %s", iFso, e.Error(), path)
	 	 CFSs = append(CFSs, nil)
		 continue
	 }
	 // --------------------------------------------------
	 //   Now that the ContentityFS has been created, it
	 //  has its own valid FSO and a valid RootContentity.
	 // --------------------------------------------------
	 L.L.Okay("Got %d item(s) total (%d dirs, %d files)",
	 	pCFS.ItemCount(), pCFS.DirCount(), pCFS.FileCount())
	 if pCFS.FileCount() == 0 {
	    	L.L.Warning("Found no content inputs in/under dir: " + path)
		continue
	 }
	 CFSs = append(CFSs, pCFS) 
     }
     return CFSs
}
