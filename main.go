package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/jozefcipa/novus/cmd"
)

// These values are replaced when building the binary by GoReleaser
var (
	version = "DEVELOPMENT"
	date    = ""
)

func main() {
	var buildDate time.Time
	if date == "" {
		buildDate = time.Now()
	} else {
		buildDate, _ = time.Parse(time.RFC3339, date)
	}

	versionString := fmt.Sprintf("Novus v%s (built on %s) on %s/%s\n", version, buildDate.Format("2006-01-02"), runtime.GOOS, runtime.GOARCH)
	cmd.Execute(versionString)
}
