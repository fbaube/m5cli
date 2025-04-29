package m5cli

import (
	"fmt"
	"io/fs"

	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
	"github.com/fbaube/mcfile"
)

// InputPathItems is for gathering, expanding (directories),
// verifying, and loading files and directories specified 
// on the command line, and then organising everything as 
// one large array of [mcfile.Contentity].
//  - [NamedPaths] is an input slice of paths of files and
//    directories; a path to a file that ends with "/" (or
//    os.Sep) throws a panic 
//  - [NamedFiles] is a slice of [fileutils.FSItem]
//    for files named e.g. on the CLI
//  - [NamedDirs]  is a slice of [fileutils.FSItem]
//    for dirs  named e.g. on the CLI
//  - [DirCntyFSs] is a slice of [mcfile.ContentityFS],
//    one per element of [NamedDirs]
//  - [AllCntys] is an output slice of [mcfile.Contentity] that 
//    collects all Contentities (a) named by [NamedFiles], and
//    (b) gathered by expanding [NamedDirs] and then walking
//    their [DirCntyFSs]
//  - Everything should implement interface [Errer]
//
// FIXME: Add NamedSymls
// . 
type InputPathItems struct {
        NamedPaths []string    // copied from input arg 
	NamedFiles []FU.FSItem // was: env.Infiles
	NamedDirrs []FU.FSItem // was: env.Indirs
	NamedMiscs []FU.FSItem // new 
	DirCntyFSs []mcfile.ContentityFS // was: env.InDirFSs
	AllCntys   []*mcfile.Contentity
}

// DoInpaths processes a list of paths of any type - files, directories,
// symlinks, "other". Its processing is pretty straightforward: 
//  - Use input []string to generate []FSItem
//  - Use FSItem.IsDir (and other funcs) to append each 
//    FSItem to the correct slice: files/dirrs/miscs 
//  - Check for errors along the way, and use the Errer
//    embedded in each FSItem 
//  - Resolve symlinks, appending them to files or dirrs, 
//    but keep them sandboxed by using [os.Root]
// .
func DoInpaths(inPaths []string) *InputPathItems {

     	var pIPI *InputPathItems // return value 
	var path string
	var FSI FU.FSItem
	var i, errct int
	var inPathItems []FU.FSItem // temp 

	pIPI = new(InputPathItems)
	pIPI.NamedFiles = make ([]FU.FSItem, 0)
	pIPI.NamedDirrs = make ([]FU.FSItem, 0)
	pIPI.NamedMiscs = make ([]FU.FSItem, 0)
	inPathItems = make ([]FU.FSItem, 0)

	for i, path = range inPaths {
	        L.L.Debug("doinpaths[%02d]: " + path, i)
		pFSI := FU.NewFSItem(path)
		inPathItems = append(inPathItems, *pFSI)
		// ERROR? 
		if pFSI.HasError() {
		   errct++
		   pFSI.SetError(&fs.PathError{
		   	Op: "NewFSItem", Path: path, Err: pFSI.GetError() })
		   L.L.Error(pFSI.Error() + ": " + path)
		 }
	}	
	L.L.Info("%d input path(s) had %d error(s)", len(inPaths), errct)
	
	L.L.Warning(SU.Rfg(SU.Ybg("=== CLI F/S ITEM(S) ===")))
	for i, FSI = range inPathItems {
	        if FSI.HasError() { continue }
		path = SU.Tildotted(FSI.FPs.AbsFP)
		var msg, sNote string
		msg = fmt.Sprintf("[%d]<%s>: ", i, path)

		switch FSI.FSItem_type {
		// REGULAR FILE?
		case FU.FSItem_type_FILE: // if FSI.IsFile() {
		     pIPI.NamedFiles = append(pIPI.NamedFiles, FSI)
		// DIRECTORY?
		case FU.FSItem_type_DIRR: // if FSI.IsDir() {
		     pIPI.NamedDirrs = append(pIPI.NamedDirrs, FSI)
		     sNote = ": to process recursively"
		// SYMLINK?
		case FU.FSItem_type_SYML:
		  // Should not happen! Cos we use Stat not Lstat
		     panic("path-clxn-00 FU.FSItem_type_SYML") /*
		     pIPI.NamedMiscs = append(pIPI.NamedDirrs, FSI)
		     symlS, symlE := 
		     sNote = ": processing is TBD" */
		  // Now this is where it gets tricky. We may or 
		  // may not want to follow a symlink, but we can 
		  // use funcs EvalSymlink & IsLocal, and os.Root.
		  // And, anything besides symlinks, fuggeddabouddit.
		  // FIXME: For now, we just attach any¨of
		  // these (incl.  symlinks) to NamedMiscs.
		 default:
		     pIPI.NamedMiscs = append(pIPI.NamedMiscs, FSI)
		     sNote = ": (TODO: check CLI symlink flag)"
		 }
		L.L.Info(msg + string(FSI.FSItem_type) + sNote)
	}
	L.L.Okay("Summary: Detected %d files, %d dirs, %d other",
		len(pIPI.NamedFiles), len(pIPI.NamedDirrs), len(pIPI.NamedMiscs))
	// if len(inPathItems) == 1 {
	//	env.IsSingleFile = true
	// }
	return pIPI
}
