package logger

import "github.com/fatih/color"

var Infof = color.New(color.FgWhite).PrintfFunc()

var Successf = color.New(color.FgHiGreen).PrintfFunc()

var Messagef = color.New(color.FgCyan).PrintfFunc()

var Errorf = color.New(color.FgRed).PrintfFunc()
