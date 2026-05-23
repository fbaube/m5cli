package m5cli

import (
	"fmt"
	"github.com/fbaube/cnty"
	// L "github.com/fbaube/mlog" // Bring in global var L
	S "strings"
	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"
)

func step3_link_references(InfileContentities []*cnty.Contentity) error {

	// ============================
	// ============================
	// TOP LEVEL: INTRA-FILE AND
	// INTER-FILE REFERENCE LINKING
	// ============================
	// ============================
	for _, p := range InfileContentities {

		// 2025.04 FIXME FIXME FIXME
		p.GatherXmlGLinks() // (&AllGLinks)

		// Cross-reference and resolve the links
		println("D=> TODO: Cross-ref the GLinks")

		// MU.Outa("Processing input file(s)", tt)

		/*
			fmt.Printf("==> Summary counts: %d Tags, %d Atts \n",
				cnty.GlobalTagCount, mcfile.GlobalAttCount)
			println("--> Tags:", cnty.GlobalTagTally.StringSortedValues())
			println("--> Atts:", cnty.GlobalAttTally.StringSortedValues())
		*/
		println("#### GLink KEY SOURCES ####")
		for _, pGL := range AllGLinks.KeyRefncs {
			fmt.Printf("%s@%s: %s: %s \n", pGL.Att, pGL.Tag,
				pGL.AddressMode, pGL.AbsFP.Tildotted())
		}
		println("#### GLink KEY TARGETS ####")
		for _, pGL := range AllGLinks.KeyRefnts {
			t := pGL.Tag
			a := pGL.Att
			// if S.HasPrefix(t, "topi") ||
			// S.Contains(a, "key") || S.Contains(a, "ref") {
			fmt.Printf("%s@%s: %s: %s \n",
				a, t, pGL.AddressMode, pGL.AbsFP.Tildotted())
			// }
		}
		println("#### GLink URI SOURCES ####")
		for _, pGL := range AllGLinks.UriRefncs {
			fmt.Printf("%s@%s: %s: %s \n", pGL.Att, pGL.Tag,
				pGL.AddressMode, pGL.AbsFP.Tildotted())
		}
		println("#### GLink URI TARGETS ####")
		for _, pGL := range AllGLinks.UriRefnts {
			t := pGL.Tag
			a := pGL.Att
			b := S.HasSuffix(pGL.Tag, "l") 
			isList := b && (len(pGL.Tag) == 2)
			if !isList {
				fmt.Printf("%s@%s: %s: %s \n", a, t,
					pGL.AddressMode, pGL.AbsFP.Tildotted())
			}
		}
	}

	// Cross-reference and resolve the links

	return nil
}

