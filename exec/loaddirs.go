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
	 pFS, e = mcfile.NewContentityFS(pDir.FPs.AbsFP, nil)
	 if e != nil {
	      	 panic("Failed: mcfile.NewContentityFS: " + pDir.FPs.AbsFP)
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