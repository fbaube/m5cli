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
// on the command line, and then organising everything as 
// one large array of [cnty.Contentity].
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
type InputPathItems struct {
        NamedPaths []string    // copied from input arg 
	NamedFiles []FU.FSObject // was: env.Infiles
	NamedDirrs []FU.FSObject // was: env.Indirs
	NamedMiscs []FU.FSObject // new 
	DirCntyFSs []cnty.ContentityFS // was: env.InDirFSs
	AllCntys   []*cnty.Contentity
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
func DoInpaths(inPaths []string) *InputPathItems {

     	var pIPI *InputPathItems // return value 
	var path string
	var FSI FU.FSObject
	var i, errct int
	var inPathItems []FU.FSObject // temp 

	pIPI = new(InputPathItems)
	pIPI.NamedFiles = make ([]FU.FSObject, 0)
	pIPI.NamedDirrs = make ([]FU.FSObject, 0)
	pIPI.NamedMiscs = make ([]FU.FSObject, 0)
	inPathItems = make ([]FU.FSObject, 0)

	for i, path = range inPaths {
	        L.L.Debug("doinpaths[%02d]: " + path, i)
		pFSI := FU.NewFSObject(path)
		inPathItems = append(inPathItems, *pFSI)
		// ERROR? 
		if pFSI.HasError() {
		   errct++
		   pFSI.SetError(&fs.PathError{
		   	Op: "NewFSObject", Path: path, Err: pFSI.GetError() })
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

		switch FSI.FSO_type {
		// REGULAR FILE?
		case FU.FSO_type_FILE: // if FSI.IsFile() {
		     pIPI.NamedFiles = append(pIPI.NamedFiles, FSI)
		// DIRECTORY?
		case FU.FSO_type_DIRR: // if FSI.IsDir() {
		     pIPI.NamedDirrs = append(pIPI.NamedDirrs, FSI)
		     sNote = ": to process recursively"
		// SYMLINK?
		case FU.FSO_type_SYML:
		  // Should not happen! Cos we use Stat not Lstat
		     panic("path-clxn-00 FU.FSO_type_SYML") /*
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
		L.L.Info(msg + string(FSI.FSO_type) + sNote)
	}
	L.L.Okay("Summary: Detected %d files, %d dirs, %d other",
		len(pIPI.NamedFiles), len(pIPI.NamedDirrs), len(pIPI.NamedMiscs))
	// if len(inPathItems) == 1 {
	//	env.IsSingleFile = true
	// }
	return pIPI
}
