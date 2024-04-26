package logger

import (
	"github.com/fatih/color"
)

var Infof = color.New(color.FgWhite).PrintfFunc()

var Successf = color.New(color.FgHiGreen).PrintfFunc()

var Messagef = color.New(color.FgCyan).PrintfFunc()

var Errorf = color.New(color.FgRed).PrintfFunc()

var Warnf = color.New(color.FgYellow).PrintfFunc()

func Checkf(format string, a ...interface{}) {
	Messagef("✔ "+format+"\n", a...)
}

// This variable gets its value in cmd/root.go from the CLI flag
var DebugEnabled bool

func Debugf(format string, a ...interface{}) {
	if DebugEnabled {
		color.New(color.FgMagenta).Printf("[DEBUG] "+format+"\n", a...)
	}
}
