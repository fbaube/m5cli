package m5cli

import (
	"fmt"
	"io/fs"

	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
	"github.com/fbaube/cnty"
)

// InputPathItems is for gathering, expanding (directories),
// verifying, and loading files and directories specified 
// on the command line, and then [still true±] organising
// or flattening everything into one large array of
// [cnty.Contentity].
//  - [NamedPaths] is an input slice of paths of files and
//    directories; a path to a file that ends with "/" (or
//    os.Sep) throws a panic 
//  - [NamedFiles] is a slice of [fileutils.FSObject]
//    for files named e.g. on the CLI
//  - [NamedDirs]  is a slice of [fileutils.FSObject]
//    for dirs  named e.g. on the CLI
//  - [DirCntyFSs] is a slice of [cnty.ContentityFS],
//    one per element of [NamedDirs]
//  - [AllCntys] is an output slice of [cnty.Contentity] that 
//    collects all Contentities (a) named by [NamedFiles], and
//    (b) gathered by expanding [NamedDirs] and then walking
//    their [DirCntyFSs]
//  - Everything should implement interface [Errer]
//
// FIXME: Add NamedSymls
// . 
type InputPathObjex struct {
        NamedPaths []string    // copied from input arg 
	NamedFiles []*FU.FSObject // was: env.Infiles
	NamedDirrs []*FU.FSObject // was: env.Indirs
	NamedMiscs []*FU.FSObject // new 
	DirCntyFSs []*cnty.ContentityFS // was: env.InDirFSs
	AllCntys   []*cnty.Contentity
}

func (p *InputPathObjex) String() string {
     	return fmt.Sprintf("paths:%d files:%d " +
	       "dirs:%d miscs:%d dirFSs:%d allCntys:%d",
	       len(p.NamedPaths), len(p.NamedFiles), len(p.NamedDirrs),
	       len(p.NamedMiscs), len(p.DirCntyFSs), len(p.AllCntys))
}

// DoInpaths processes a list of paths of any type - files, directories,
// symlinks, "other". Its processing is pretty straightforward: 
//  - Use input []string to generate []FSObject
//  - Use FSObject.IsDir (and other funcs) to append each 
//    FSObject to the correct slice: files/dirrs/miscs 
//  - Check for errors along the way, and use the Errer
//    embedded in each FSObject 
//  - Resolve symlinks, appending them to files or dirrs, 
//    but keep them sandboxed by using [os.Root]
// .
func DoInpaths(inPaths []string) *InputPathObjex {

     	var pIPO *InputPathObjex // return value 
	var path string
	var pFSO *FU.FSObject
	var i, errct int
	var inPathObjex []*FU.FSObject // temp

	pIPO = new(InputPathObjex)
	pIPO.NamedFiles = make ([]*FU.FSObject, 0)
	pIPO.NamedDirrs = make ([]*FU.FSObject, 0)
	pIPO.NamedMiscs = make ([]*FU.FSObject, 0)
	inPathObjex = make ([]*FU.FSObject, 0)

	for i, path = range inPaths {
	        L.L.Debug("doinpaths[%02d]: " + path, i)
		pFSO := FU.NewFSObject(path)
		inPathObjex = append(inPathObjex, pFSO)
		// ERROR? 
		if pFSO.HasError() {
		   errct++
		   pFSO.SetError(&fs.PathError{
		   	Op: "NewFSObject", Path: path, Err: pFSO.GetError() })
		   L.L.Error(pFSO.Error() + ": " + path)
		 }
	}	
	L.L.Info("%d input path(s) had %d error(s)", len(inPaths), errct)
	
	L.L.Warning(SU.Rfg(SU.Ybg("=== CLI INPATH(S) & F/S ITEM(S) ===")))
	for i, pFSO = range inPathObjex {
	        if pFSO.HasError() { continue }
		path = SU.Tildotted(pFSO.FPs.AbsFP)
		var msg, sNote string
		msg = fmt.Sprintf("[%d]<%s>", i, path)

		// REGULAR FILE?
		if pFSO.IsFile() {
		     pIPO.NamedFiles = append(pIPO.NamedFiles, pFSO)
		} else
		// DIRECTORY?
		if pFSO.IsDir() {
		     pIPO.NamedDirrs = append(pIPO.NamedDirrs, pFSO)
		     sNote = "to process recursively"
		} else
		// SYMLINK?
		if pFSO.IsSymlink() {
		  // Should not happen! Cos we use Stat not Lstat
		     panic("path-clxn-00 FU.FSO_type_SYML") /*
		     pIPO.NamedMiscs = append(pIPO.NamedDirrs, pFSO)
		     symlS, symlE := 
		     sNote = "processing is TBD" */
		  // Now this is where it gets tricky. We may or 
		  // may not want to follow a symlink, but we can 
		  // use funcs EvalSymlink & IsLocal, and os.Root.
		  // And, anything besides symlinks, fuggeddabouddit.
		  // FIXME: For now, we just attach any¨of
		  // these (incl.  symlinks) to NamedMiscs.
		 } else { 
		     pIPO.NamedMiscs = append(pIPO.NamedMiscs, pFSO)
		     sNote = "[TODO: check CLI symlink flag]"
		 }
		L.L.Info(msg + ": " + string(pFSO.FSO_type) + sNote)
	}
	L.L.Info("DoInpaths OUT: " + pIPO.String())
	L.L.Okay("Summary: Detected %d files, %d dirs, %d other",
		len(pIPO.NamedFiles), len(pIPO.NamedDirrs), len(pIPO.NamedMiscs))
	// if len(inPathObjex) == 1 {
	//	env.IsSingleFile = true
	// }
	return pIPO
}
