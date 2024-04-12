package m5cli

import (
	LU "github.com/fbaube/logutils"
)

/* REF
// RFC5424 log message levels.
Dbg, Prog, Info, G Y R, Panic 
const (
        LevelPanic    Level = iota + 2 // i.e. 0 + 2 
        LevelError          // 3 R
        LevelWarning        // 4 Y
        LevelOkay           // 5 G
        LevelInfo           // 6
        LevelProgress       // 7
        LevelDbg            // misspelled cos 8 != RFC5424 "7"
*/

var LOG_LEVEL_FILE_INTRO = LU.LevelOkay // 5
var LOG_LEVEL_FILE_READING = LU.LevelOkay // 5
var LOG_LEVEL_EXEC_STAGES = LU.LevelDbg
var LOG_LEVEL_REF_LINKING = LU.LevelInfo // 6 
var LOG_LEVEL_WEB = LU.LevelDbg

