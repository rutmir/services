package log

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"runtime"
)

type Level byte

type LogSettings struct {
	PathSeparator string
}

const (
	LEVEL_EMERGENCY Level = iota // 0
	LEVEL_ALERT                  // 1
	LEVEL_CRITICAL               // 2
	LEVEL_ERROR                  // 3
	LEVEL_WARNING                // 4
	LEVEL_NOTICE                 // 5
	LEVEL_INFO                   // 6
	LEVEL_DEBUG                  // 7
	LEVEL_FATAL                  // 8
)

var settings *LogSettings

/*API*/

// Emergency
func Emergency(v ...interface{}) {
	str := extractString(v ...)
	file, line := extractPath()

	doLocalLog(str, LEVEL_EMERGENCY, file, line)
}

// Alert
func Alert(v ...interface{}) {
	str := extractString(v ...)
	file, line := extractPath()

	doLocalLog(str, LEVEL_ALERT, file, line)
}

// Critical
func Critical(v ...interface{}) {
	str := extractString(v ...)
	file, line := extractPath()

	doLocalLog(str, LEVEL_CRITICAL, file, line)
}

// Err
func Err(v ...interface{}) {
	str := extractString(v ...)
	file, line := extractPath()

	doLocalLog(str, LEVEL_ERROR, file, line)
}

// Warn
func Warn(v ...interface{}) {
	str := extractString(v ...)
	file, line := extractPath()

	doLocalLog(str, LEVEL_WARNING, file, line)
}

// Notice
func Notice(v ...interface{}) {
	str := extractString(v ...)
	file, line := extractPath()

	doLocalLog(str, LEVEL_NOTICE, file, line)
}

// Info
func Info(v ...interface{}) {
	str := extractString(v ...)
	file, line := extractPath()

	doLocalLog(str, LEVEL_INFO, file, line)
}

// Debug
func Debug(v ...interface{}) {
	str := extractString(v ...)
	file, line := extractPath()

	doLocalLog(str, LEVEL_DEBUG, file, line)
}

// Fatal
func Fatal(v ...interface{}) {
	str := extractString(v ...)
	file, line := extractPath()

	doLocalLog(str, LEVEL_FATAL, file, line)
	os.Exit(1)
}

func extractPath() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = ""
		line = 0
	}else {
		file = strings.Split(file, settings.PathSeparator)[1]
	}
	return file, line
}

func doLocalLog(str string, level Level, file string, line int) {
	var levelStr string
	switch level{
	case LEVEL_NOTICE:
		levelStr = "Notice"
		break
	case LEVEL_WARNING:
		levelStr = "Warning"
		break
	case LEVEL_DEBUG:
		levelStr = "Debug"
		break
	case LEVEL_ALERT:
		levelStr = "Alert"
		break
	case LEVEL_EMERGENCY:
		levelStr = "Emergency"
		break
	case LEVEL_INFO:
		levelStr = "Info"
		break
	case LEVEL_ERROR:
		levelStr = "Error"
		break
	case LEVEL_CRITICAL:
		levelStr = "Critical"
		break
	case LEVEL_FATAL:
		levelStr = "Fatal"
		break
	default:
		levelStr = "-"
	}
	fmt.Println("[" + levelStr + "][" + file + ":" + strconv.Itoa(line) + "] - " + str + "\n")

}

func extractString(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}
	if first, ok := v[0].(string); ok {
		if len(v) > 1 {
			return fmt.Sprintf(first, v [1:]...)
		}
		return first
	}else {
		return fmt.Sprint(v ...)
	}
}

func init() {
	pathSeparator := os.Getenv("LOG_PATH_SEPARATOR")
	if pathSeparator == "" {
		pathSeparator = "/src"
	}
	settings = &LogSettings{
		PathSeparator :pathSeparator}
}
