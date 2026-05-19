package exec

// TODO: SHOULD USE
// path.Clean (rmvs trlg slashes) 
// fs.ValidPath
// FP.IsLocal (implies ValidPath, so do VP first)


import(
	   "github.com/fbaube/cnty"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	L  "github.com/fbaube/mlog"
)

// LoadDirpathsContentFSs turns a slice of [FSObject] into
// a slice of [ContentityFS]. Any error is returned as 
// an interface [Errer] of a ContentityFS. 
func LoadDirpathsContentFSs(inFSOs []FU.FSObject) ([]cnty.ContentityFS) {
     if inFSOs == nil || len(inFSOs) == 0 {
     	return nil
	}
     var FSs []cnty.ContentityFS
     var pFS  *cnty.ContentityFS

     // For every input FSObject
     for iDir, pDir := range inFSOs {
     	 var shortName = FU.EnsureTrailingPathSep(
	     SU.Tildotted(pDir.FPs.AbsFP))
	 L.L.Info("InDir[%d]: %s", iDir, shortName)
	 var e error
	 // nil is []string of OK file extensions 
	 pFS, e = cnty.NewContentityFS(pDir.FPs.AbsFP, nil)
	 if e != nil { /*
	      	 isRillyNil := reflect.ValueOf(e).Kind() ==
		 	       reflect.Ptr && reflect.ValueOf(e).IsNil()
		 if isRillyNil { fmt.Printf("IT IS OK NOT ERROR") }
	      	 println(fmt.Sprintf("error %T %p", e, e))
	      	 println(fmt.Sprintf("error %+v", e))
	      	 fmt.Printf("InDir[%d]: %s: error: %s", iDir, shortName, e.Error()) */
		 L.L.Error("InDir[%d]: %s: error: %s", iDir, shortName, e.Error())
	      	 // panic("Failed: cnty.NewContentityFS: " + pDir.FPs.AbsFP)
		 continue
	 }
	 L.L.Okay("Found %d item(s) total (%d dirs, %d files)",
	 	pFS.ItemCount(), pFS.DirCount(), pFS.FileCount())
	 if pFS.FileCount() == 0 {
	    	L.L.Warning("Found no content inputs to " +
			"process in dir: " + shortName)
		continue
	 }
	 FSs = append(FSs, *pFS) 
     }
     return FSs
}
