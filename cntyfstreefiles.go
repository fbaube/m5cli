package m5cli

import(
	"os"
	"fmt"
	"github.com/fbaube/cnty"
	L "github.com/fbaube/mlog"
)

func WriteContentityFStreeFiles(IndirContentityFSs []cnty.ContentityFS) {

     var e error 
     for iFS, pFS := range IndirContentityFSs {
		// ==============================
		// Write out a tree rep to a file
		// ==============================
		var treeFile *os.File
		var treeFilename string
		// Now write out a tree representation
		treeFilename = fmt.Sprintf("./input-tree-%d", iFS)
		treeFile, e = os.Create(treeFilename)
		if e != nil {
			// An error here does not need to be fatal
			L.L.Error("Treefile " + treeFilename + ": " + e.Error())
		} else {
			pFS.RootContentity().PrintTree(treeFile)
			L.L.Okay("Wrote input tree file: " + treeFilename)
		}
		treeFile.Close()

		// =================================
		// Write out a css-enabled html tree
		// =================================
		// L.L.Warning("SKIP'D: css-enabled tree for html, exec.go.L181")
		treeFilename = fmt.Sprintf("./css-tree-%d", iFS)
		treeFile, e = os.Create(treeFilename)
		if e != nil {
			// An error here does not need to be fatal
			L.L.Error("CssTreefile " + treeFilename + ": " + e.Error())
		} else {
		        L.L.Error("cntyfsfreefiles.go FIXME")
		//	pFS.RootContentity().PrintCssTree(treeFile)
		//	L.L.Okay("Wrote css tree file: " + treeFilename)
		}
		treeFile.Close()
	}
}
