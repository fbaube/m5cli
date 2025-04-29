package m5cli

import (
	// "fmt"
	// "github.com/fbaube/mcfile"
	// L "github.com/fbaube/mlog" // Bring in global var L
	// S "strings"
	// mime "github.com/fbaube/fileutils/contentmime"
	// "github.com/fbaube/tags"
)

func validateInputFiles(env *XmlAppEnv) error {

	// ======================
	//  VALIDATE INPUT FILES
	// ======================
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

