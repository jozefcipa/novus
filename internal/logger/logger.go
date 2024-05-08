package logger

import "fmt"

// Bash Color codes
// https://stackoverflow.com/a/69648792/4480179
const (
	RED     string = "\033[38;5;196m"
	GREEN   string = "\033[38;5;046m"
	CYAN    string = "\033[38;5;014m"
	ORANGE  string = "\033[38;5;202m"
	YELLOW  string = "\033[38;5;011m"
	MAGENTA string = "\033[38;5;201m"
	GRAY    string = "\033[38;5;245m"
	WHITE   string = "\033[38;5;255m"
	RESET   string = "\033[0m"
)

func Infof(format string, a ...interface{}) {
	fmt.Printf(WHITE+format+RESET+"\n", a...)
}

func Warnf(format string, a ...interface{}) {
	fmt.Printf(YELLOW+"⚠️  "+format+RESET+"\n", a...)
}

func Checkf(format string, a ...interface{}) {
	fmt.Printf(GREEN+" ✔ "+GRAY+format+RESET+"\n", a...)
}

func Successf(format string, a ...interface{}) {
	fmt.Printf(GREEN+"✅ "+format+RESET+"\n", a...)
}

func Hintf(format string, a ...interface{}) {
	fmt.Printf(ORANGE+"💡 "+format+RESET+"\n", a...)
}

func Errorf(format string, a ...interface{}) {
	fmt.Printf(RED+"❌ "+format+RESET+"\n", a...)
}

// This variable gets its value in cmd/root.go from the CLI flag
var DebugEnabled bool

func Debugf(format string, a ...interface{}) {
	if DebugEnabled {
		fmt.Printf(MAGENTA+"[DEBUG] "+format+RESET+"\n", a...)
	}
}
