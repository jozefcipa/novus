package main

import (
	"time"

	"github.com/jozefcipa/novus/cmd"
)

// These values will be set by GoReleaser during the build process.
var (
	version = "DEV"
	date    = ""
)

func main() {
	var buildDate time.Time
	if date == "" {
		buildDate = time.Now()
	} else {
		buildDate, _ = time.Parse(time.RFC3339, date)
	}

	cmd.Execute(version, buildDate.Format("2006-01-02"))
}
