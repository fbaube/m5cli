package exec

import(
	"fmt"
	"errors"
	"github.com/fbaube/mcfile"
	FU "github.com/fbaube/fileutils"
)

func LoadFilepathsContents(ff []FU.FSItem) ([]*mcfile.Contentity, []error) {

     if ff == nil || len(ff) == 0 {
     	return nil, nil
	}
     var pCC []*mcfile.Contentity
     var pC    *mcfile.Contentity
     var ee  []error
     var e, eC error
     var path string 

     // For every input FSItem
     for i, p := range ff {
     	 path = p.FPs.AbsFP.S()
	 pC, e = mcfile.NewContentity(path)
	 if pC == nil || e != nil || pC.HasError() {
		if e == nil {
		   e = errors.New("placeholder error")
		}
		eC = fmt.Errorf("exec.loadFP: newcontentity<[%d]%s>: %w", i, path, e)
		ee = append(ee, eC)
		continue
	 }
	 if pC.MarkupType() == "UNK" {
		eC = fmt.Errorf("exec.loadFP: " +
		    "newcontentity<[%d:len%d]%s>: markupType<%s>",
		     	i, len(p.Raw), path, pC.MarkupType())
		ee = append(ee, eC)
                continue
	 }
	 pCC = append(pCC, pC)
	}
	return pCC, ee 
}