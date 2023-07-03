// flargdefs is flag argument definitions

package mcm_cli

import (
	"fmt"
	SU "github.com/fbaube/stringutils"
)

type pathFlargs struct {
	sInpaths []string
	sOutdir, sDbdir, sXmlcatlgfile,
	sXmlschemasdir string
}

func (p pathFlargs) String() string {
	return fmt.Sprintf("db<%s> cat<%s> sch<%s> out<%s> in's<%+v>",
		p.sDbdir, p.sXmlcatlgfile, p.sXmlschemasdir,
		p.sOutdir, p.sInpaths)
}

type boolFlargs struct {
	FollowSymLinks, DoArchive, DBdoImport, Help, Debug,
	// Pritt, GroupGenerated, GTokens, GTree,
	TotalTextal, Validate, DBdoZeroOut bool
}

func (b boolFlargs) String() string {
	return fmt.Sprintf("debug:%s ttlTxtl:%s help:%s import:%s "+
		"validate:%s zeroOutDB:%s",
		SU.Yn(b.Debug), SU.Yn(b.TotalTextal),
		SU.Yn(b.Help), SU.Yn(b.DBdoImport),
		SU.Yn(b.Validate), SU.Yn(b.DBdoZeroOut))
}

type AllFlargs struct {
	p        pathFlargs
	b        boolFlargs
	restPort int
	webPort  int
}

func (a AllFlargs) String() string {
	return fmt.Sprintf("[PATHS]<%s> [BOOLS]<%s> WEB-port<%d> REST-port<%d>",
		a.p.String(), a.b.String(), a.webPort, a.restPort)
}
