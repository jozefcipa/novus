package logger

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type StopFuncs struct {
	Checkf LoggerFunc
	Errorf LoggerFunc
	Infof  LoggerFunc
	Done   func()
}

var spinnerParts []string

func init() {
	spinnerParts = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
}

func Loadingf(format string, a ...interface{}) StopFuncs {
	s := spinner.New(spinnerParts, 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(GRAY+" "+format+RESET, a...)
	s.Start()

	return StopFuncs{
		Checkf: func(format string, a ...interface{}) {
			s.FinalMSG = fmt.Sprintf(formatCheck(format), a...)
			s.Stop()
		},
		Errorf: func(format string, a ...interface{}) {
			s.FinalMSG = fmt.Sprintf(formatError(format), a...)
			s.Stop()
		},
		Infof: func(format string, a ...interface{}) {
			s.FinalMSG = fmt.Sprintf(formatInfo(format), a...)
			s.Stop()
		},
		Done: func() {
			s.Stop()
		},
	}
}
