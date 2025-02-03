package exec

// NEED TO USE
// path.Clean (rmvs trlg slashes) 
// fs.ValidPath
// FP.IsLocal (implies ValidPath, so do VP first)


import(
	"github.com/fbaube/mcfile"
	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	L "github.com/fbaube/mlog"
)

// LoadDirpathsContentFSs turns a slice of [FSItem] into
// a slice of [ContentityFS]. Any error is returned as 
// an interface [Errer] of a ContentityFS. 
func LoadDirpathsContentFSs(ff []FU.FSItem) ([]*mcfile.ContentityFS) {
     if ff == nil || len(ff) == 0 {
     	return nil
	}
     var pFSs []*mcfile.ContentityFS
     var pFS   *mcfile.ContentityFS

     // For every input FSItem
     for iDir, pDir := range ff {
     	 var shortName = FU.EnsureTrailingPathSep(
	     SU.Tildotted(pDir.FPs.AbsFP))
	 L.L.Info("InDir[%d]: %s", iDir, shortName)
	 var e error
	 // nil is []string of OK file extensions 
	 pFS, e = mcfile.NewContentityFS(pDir.FPs.AbsFP, nil)
	 if e != nil { /*
	      	 isRillyNil := reflect.ValueOf(e).Kind() == reflect.Ptr && reflect.ValueOf(e).IsNil()
		 if isRillyNil { fmt.Printf("IT IS OK NOT ERROR") }
	      	 println(fmt.Sprintf("error %T %p", e, e))
	      	 println(fmt.Sprintf("error %+v", e))
	      	 fmt.Printf("InDir[%d]: %s: error: %s", iDir, shortName, e.Error()) */
		 L.L.Error("InDir[%d]: %s: error: %s", iDir, shortName, e.Error())
	      	 // panic("Failed: mcfile.NewContentityFS: " + pDir.FPs.AbsFP)
		 continue
	 }
	 L.L.Okay("Found %d item(s) total (%d dirs, %d files)",
	 	pFS.ItemCount(), pFS.DirCount(), pFS.FileCount())
	 if pFS.FileCount() == 0 {
	    	L.L.Warning("Found no content inputs to " +
			"process in dir: " + shortName)
		continue
	 }
	 pFSs = append(pFSs, pFS) 
     }
     return pFSs
}