package m5cli

// flargdefs is flag argument definitions

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
	FollowSymLinks, DBdoImport, Debug, // DoArchive, Help, 
	// Pritt, GroupGenerated, GTokens, GTree,
	Samples, TotalTextal, Validate, DBdoZeroOut bool
}

func (b boolFlargs) String() string {
	return fmt.Sprintf("debug:%s ttlTxtl:%s import:%s "+
		"samples:%s validate:%s zeroOutDB:%s",
		SU.Yn(b.Debug), SU.Yn(b.TotalTextal),
		SU.Yn(b.DBdoImport), SU.Yn(b.Samples), 
		SU.Yn(b.Validate), SU.Yn(b.DBdoZeroOut))
		// help:%s SU.Yn(b.Help) 
}

type AllFlargs struct {
	p        pathFlargs
	b        boolFlargs
	restPort int
	webPort  int
}

func (a AllFlargs) String() string {
	return fmt.Sprintf("[PATHS]|%s| [BOOLS]|%s| WEB-port<%d> REST-port<%d>",
		a.p.String(), a.b.String(), a.webPort, a.restPort)
}
