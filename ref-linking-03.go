package m5cli

import (
	"fmt"
	"github.com/fbaube/mcfile"
	// L "github.com/fbaube/mlog" // Bring in global var L
	S "strings"
	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"
)

func ref_linking_03(env *XmlAppEnv, InfileContentities []*mcfile.Contentity) error {

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
				mcfile.GlobalTagCount, mcfile.GlobalAttCount)
			println("--> Tags:", mcfile.GlobalTagTally.StringSortedValues())
			println("--> Atts:", mcfile.GlobalAttTally.StringSortedValues())
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

	// ===========================
	//   3. VALIDATE INPUT FILES
	//  This code actually belongs
	//     above: see line 145.
	// ===========================

	// We can use xmllint to validate here, but we don't
	// want to rely on schema files already in the system
	// - no using its normal catalogs, or anything at or
	// under `/etc/xml/catalog` - and we can't use the
	// envar `XML_CATALOG_FILES`. But if we have our own
	// catalog file, we can pass it our own value as (say)
	// envar `MMMC_XML_CATALOG_FILES`.

	/* if pCA.Validate && * / env.XmlCatalogFile != nil {

		println(" ")
		tt = MU.Into("")
		println("==> Validating input file(s)...")

		var dtdStatus, docStatus, errStatus string

		print("==> Text file validation statuses: \n")
		for _, pMCF = range MCFiles {
			if pMCF.IsXML() {
				if FU.MTypeSub(pMCF.MType, 0) == "img" {
					continue
				}
			}
			if !pMCF.IsXML() {
				continue
			}
			var dtdDesc string

			// DO THE VALIDATION
			dtdStatus, docStatus, errStatus = pMCF.DoValidation(env.XmlCatalogFile)

			if pMCF.XmlDoctypeFields != nil {
				dtdDesc = pMCF.XmlDoctypeFields.PIDSIDcatalogFileRecord.PublicTextDesc
			}
			fmt.Printf("%s/%s/%s %s %s :: %s :: %s  \n",
				pMCF.MType[0], pMCF.MType[1], pMCF.MType[2], dtdStatus,
				docStatus, pMCF.AbsFilePath, dtdDesc)
			if errStatus != "" {
				println(errStatus)
			}
		}
		MU.Outa("Validating input file(s)", tt)
	}
	*/
	return nil
}

