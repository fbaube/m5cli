package m5cli

import(
        "os"
        "fmt"
	"io"
        "github.com/fbaube/mcfile"
        L "github.com/fbaube/mlog"
        SU "github.com/fbaube/stringutils"
)

func InitContentityDebugFiles(InfileContentities []*mcfile.Contentity, doTotalTextal bool) {

     for ii, cty := range InfileContentities {
		// L.L.Info("LOOPING on IN-Cty %d", ii)
		if cty.IsDir() {
			// println("Skip dir: " + cty.AbsFP())
			continue
		}
		L.L.SetCategory(fmt.Sprintf("%02d", ii))

		// Input file VALIDATION belongs HERE!
		// See the code below, line 315.

		// Now, files for capturing debugging output:
		// func Create(name string) (*File, error) API docs:
		// Create or truncate the named file. If the file
		// already exists, it is truncated (i.e. emptied).
		// If it does not exist, it is created with mode 0666
		// (before umask). If successful, the returned File
		// (IS OPEN, apparently, AND) can be used for I/O;
		// the associated file descriptor has mode O_RDWR.
		// If there is an error, it will be of type *PathError.

		cty.GTknsWriter = io.Discard
		cty.GTreeWriter = io.Discard
		cty.GEchoWriter = io.Discard

		if doTotalTextal { // env.cfg.b.TotalTextal {
			fnm := cty.AbsFP()
			if len(InfileContentities) == 1 {
				fnm = "./debug"
			}
			// To name output files, append
			// "_(echo,tkns,tree)" to the entire file name.
			echoName := (fnm + "_echo")
			tknsName := (fnm + "_tkns")
			treeName := (fnm + "_tree")
			f1, e1 := os.Create(echoName)
			f2, e2 := os.Create(tknsName)
			f3, e3 := os.Create(treeName)
			if e1 == nil && e2 == nil && e3 == nil {
				L.L.Info("created %s _echo,_tkns,_tree",
					SU.Tildotted(fnm))
			} else {
				L.L.Error("Cannot open a Total Textual file")
			}
			cty.GEchoWriter = f1
			cty.GTknsWriter = f2
			cty.GTreeWriter = f3
		}
	}
}