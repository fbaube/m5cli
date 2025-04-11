package m5cli

import (
	L "github.com/fbaube/mlog" // Brings in global var L log.Logger
)

// InitLogging starts with Targets[] empty, then
//  1) adds a console logger (i.e. Stdout),
//  2) adds a file logger named "[appName].log",
//  3) opens them both,
//  4) outputs a few demo log records
func InitLogging(appName string) {
	// (1) CONSOLE LOGGER
	L.L.Targets = append(L.L.Targets, L.NewConsoleTarget())
	// (2) FILE LOGGER
	fileTarget := L.NewFileTarget()
	if appName == "" {
		appName = "mmmc"
	}
	fileTarget.FileName = appName + ".log"
	L.L.Targets = append(L.L.Targets, fileTarget)
	// (3) OPEN (and list) LOGGING TARGETS
	L.L.Open()
	L.L.Info("Log targets: %d", len(L.L.Targets))
	for ii, pp := range L.L.Targets {
		L.L.Info("\t Log target [%d]: %T", ii, pp)
	}
	/*
        // RE-ENABLE AS NEEDED
	// DEMO MESSAGES
	L.L.Dbg("example msg %v dbg", 10)
	L.L.Progress("example msg %v prog", 20)
	L.L.Info("example msg %v info", 30)
	L.L.Okay("example msg %v okay", 40)
	L.L.Warning("example msg %v warng", 40)
	L.L.Error("example msg %v error", 50)
	L.L.Panic("example msg %v panic", 60)
	L.L.Dbg("- end - %v - end -", "***")
	*/
}
