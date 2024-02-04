package logger

import (
	"github.com/fatih/color"
)

var Infof = color.New(color.FgWhite).PrintfFunc()

var Successf = color.New(color.FgHiGreen).PrintfFunc()

var Messagef = color.New(color.FgCyan).PrintfFunc()

var Errorf = color.New(color.FgRed).PrintfFunc()

func Checkf(format string, a ...interface{}) {
	Messagef(" ✔ "+format+"\n", a...)
}

var DebugEnabled bool

func Debugf(format string, a ...interface{}) {
	if DebugEnabled {
		color.New(color.FgYellow).Printf("[DEBUG] "+format+"\n", a...)
	}
}
