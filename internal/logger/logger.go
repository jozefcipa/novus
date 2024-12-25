package logger

import "fmt"

// Bash Color codes
// https://stackoverflow.com/a/69648792/4480179
const (
	RED       string = "\033[38;5;196m"
	GREEN     string = "\033[38;5;046m"
	CYAN      string = "\033[38;5;014m"
	ORANGE    string = "\033[38;5;202m"
	YELLOW    string = "\033[38;5;011m"
	MAGENTA   string = "\033[38;5;201m"
	GRAY      string = "\033[38;5;245m"
	WHITE     string = "\033[38;5;255m"
	UNDERLINE string = "\033[4m"
	RESET     string = "\033[0m"
)

type LoggerFunc func(format string, a ...interface{})

func formatInfo(format string) string {
	return GRAY + format + RESET + "\n"
}

func Infof(format string, a ...interface{}) {
	fmt.Printf(formatInfo(format), a...)
}

func Warnf(format string, a ...interface{}) {
	fmt.Printf(ORANGE+"⚠️  "+format+RESET+"\n", a...)
}

func formatCheck(format string) string {
	return GREEN + "✔ " + GRAY + format + RESET + "\n"
}

func Checkf(format string, a ...interface{}) {
	fmt.Printf(formatCheck(format), a...)
}

func Successf(format string, a ...interface{}) {
	fmt.Printf(GREEN+"✅ "+format+RESET+"\n", a...)
}

func Hintf(format string, a ...interface{}) {
	fmt.Printf(YELLOW+"💡 "+format+RESET+"\n", a...)
}

func formatError(format string) string {
	return RED + "❌ " + format + RESET + "\n"
}

func Errorf(format string, a ...interface{}) {
	fmt.Printf(formatError(format), a...)
}

// This variable gets its value in cmd/root.go from the CLI flag
var DebugEnabled bool

func Debugf(format string, a ...interface{}) {
	if DebugEnabled {
		fmt.Printf(MAGENTA+"[DEBUG] "+format+RESET+"\n", a...)
	}
}
