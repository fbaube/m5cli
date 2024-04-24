package m5cli

import (
	// "runtime/pprof"
	// "github.com/davecheney/profile"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog" 
	SU "github.com/fbaube/stringutils"
)

// DoSamples does demo type stuff and can be
// skipped; it does use (and demo) mlog's L.L
func DoSamples() {

	// What time it is
	L.L.Info("Run @ %s (YMDHm: %s)", SU.NowPlus(), SU.NowAsYMDHM())
	// L.L.Progress("Zulu: %s", S.Replace(time.Now().UTC().
	// 	Format(time.RFC3339), "T", "_", 1))

	/*
		f, err := os.Create("cpuprofile")
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"could not create CPU profile: %v\n", err)
			os.Exit(1)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr,
				"could not start CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	*/

	/*
		// DEMO: Let's try ERRORS
		var ERR = errors.New("My ERR text")
		var OPE os.PathError
		var PPE FU.PathPropsError
		var CTE mcfile.ContentityError

		OPE = fs.PathError{
			Op: "myPEop", Path: "my/PE/path", Err: ERR}
		L.L.Error("Demo of (os.PathError:%T).Error(): \n\t => %s",
			OPE, OPE.Error())

		p1, p1e := FU.NewPathProps("/zorkle")     // ERROR
		p2, _ := FU.NewPathProps("/Users/fbaube") // OKAY

		// NewPathPropsError(ermsg string, op string, pp *PathProps)
		PPE = FU.NewPathPropsError("My error text", "myPEop", p2)
		L.L.Error("Demo of (new  fu.PropsPathError:%T).Error(): \n\t => %s",
			PPE, PPE.Error())
		// WrapAsPathPropsError(e error, op string, pp *PathProps)
		PPE = FU.WrapAsPathPropsError(p1e, "myPEop", p1)
		L.L.Error("Demo of (rapt fu.PropsPathError:%T).Error(): \n\t => %s",
			PPE, PPE.Error())

		// !! p3 := mcfile.NewContentity("/zorkle") // ERROR
		// !! p3e := p3.Error()                           // !! ##
		// ignore the error return value, it is not being demo'd here !
		p4, _ := mcfile.NewContentity("/Users/blecch")

		// NewContentityError(ermsg string, op string, cty *Contentity)
		CTE = mcfile.NewContentityError("My error text", "myCTop", p4)
		L.L.Error("Demo of (new  mcf.ContentityError:%T).Error(): \n\t => %s",
			CTE, CTE.Error())
		// WrapAsContentityError(e error, op string, cty *Contentity)

		/* ARGHH
		CTE = mcfile.WrapAsContentityError(p3e, "myCTop", p4)
		L.L.Error("Demo of (rapt mcf.ContentityError:%T).Error(): \n\t => %s",
			CTE, CTE.Error())
	*/

	// includes color demo
	L.L.Info("ID: %s", FU.SessionSummary())

	/*
		// DEMO: Let's try a wrap'd traced error
		println("==> Wrap'd Traced Error Demo START <==")
		var E1, E2 error
		E1 = MU.TracedError(errors.New("My demo mu.TracedError"))
		E2 = fmt.Errorf("Its wrapping info: %w", E1)
		MU.ErrorTrace(os.Stdout, E2)
		// FIXME: glog.ErrorTrace(glog.SessionLogger, E2)
		println("==>", "Error Demo END", "<====")
	*/
}
