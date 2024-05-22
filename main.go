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
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	versionString := fmt.Sprintf("Novus %s (%s) on %s/%s\n", version, date, runtime.GOOS, runtime.GOARCH)
	cmd.Execute(versionString)
}
