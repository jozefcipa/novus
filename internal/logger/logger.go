package logger

import (
	"github.com/fatih/color"
)

var Messagef = color.New(color.FgCyan).PrintfFunc()

var Warnf = color.New(color.FgYellow).PrintfFunc()

func Checkf(format string, a ...interface{}) {
	Messagef("✔ "+format+"\n", a...)
}

func Hintf(format string, a ...interface{}) {
	Warnf("💡 "+format+"\n", a...)
}

func Successf(format string, a ...interface{}) {
	color.New(color.FgHiGreen).PrintfFunc()("✅ "+format+"\n", a...)
}

func Errorf(format string, a ...interface{}) {
	color.New(color.FgRed).PrintfFunc()("❌ "+format+"\n", a...)
}

// This variable gets its value in cmd/root.go from the CLI flag
var DebugEnabled bool

func Debugf(format string, a ...interface{}) {
	if DebugEnabled {
		color.New(color.FgMagenta).Printf("[DEBUG] "+format+"\n", a...)
	}
}
