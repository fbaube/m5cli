package m5cli

import (
	"github.com/fbaube/mcfile"
	XU "github.com/fbaube/xmlutils"
)

// FIXME:
// hide a flag by specifying its name
// flags.MarkHidden("secretFlag")

/* FIXME:
// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (f *FlagSet) SetOutput(output io.Writer) {
	f.output = output
}
*/

var multipleXmlCatalogFiles []*XU.XmlCatalogFile

// InputExts is a whitelist that more than covers 
// the file types associated with the LwDITA spec.
// Of course, check for them case-insensitively.
var InputExts = []string{
	".dita", ".map", ".ditamap", ".xml",
	".md", ".markdown", ".mdown", ".mkdn",
	".html", ".htm", ".xhtml",
	".png", ".gif", ".jpg", ".svg" }

// AllGLinks gathers all [mcfile.GLinks] in the current
// run's input set. This should actually be re-entrant,
// like [mcfile.ContentityEngine].
var AllGLinks mcfile.GLinks

