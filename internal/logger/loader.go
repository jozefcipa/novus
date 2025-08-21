package logger

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	// Catch Cmd+C and termination signals
	intSignal := make(chan os.Signal, 1)
	signal.Notify(intSignal, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-intSignal
		s.Stop()
	}()

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
