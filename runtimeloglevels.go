package m5cli

import (
	LU "github.com/fbaube/logutils"
)

/* REF
// RFC5424 log message levels.
Debug, Info, G Y R, Panic 
const (
        LevelPanic    Level = iota + 2 // i.e. 0 + 2 
        LevelError          // 3 R
        LevelWarning        // 4 Y
        LevelOkay           // 5 G
        LevelInfo           // 6
        LevelDebug          // 7 
*/

var LOG_LEVEL_FILE_INTRO   = LU.LevelInfo  // 6
var LOG_LEVEL_FILE_READING = LU.LevelOkay  // 5
var LOG_LEVEL_EXEC_STAGES  = LU.LevelDebug
var LOG_LEVEL_REF_LINKING  = LU.LevelInfo  // 6 
var LOG_LEVEL_WEB          = LU.LevelDebug // 7

