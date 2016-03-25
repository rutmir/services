package log

import (
	"testing"
)

func TestInfo(t *testing.T) {
	Info(2, 1)
}

func TestEmergency(t *testing.T) {
	Emergency("test %s %v", "Emergency", 2)
}

func TestAlert(t *testing.T) {
	Alert("test %s %v", "Alert", 3)
}

func TestCritical(t *testing.T) {
	Critical("test %s", "Critical", 4, 7, "trau")
}

func TestErr(t *testing.T) {
	Err("test %s %v", "Err", 5)
}

func TestWarn(t *testing.T) {
	Warn("test %s %v", "Warn", 6)
}

func TestNotice(t *testing.T) {
	Notice("test %s %v", "Notice", 7)
}

func TestDebug(t *testing.T) {
	Debug("test %s %v", "Debug", 8)
}
