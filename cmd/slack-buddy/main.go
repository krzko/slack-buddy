package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/getlantern/systray"
	sh "github.com/krzko/slack-buddy/pkg/systray"
)

var (
	Version     = "v0.0.0"
	ShortCommit = "0000000"
	CommitDate  = "1970-01-01T00:00:00Z"
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "Show version number")

	flag.Parse()

	if showVersion {
		version := fmt.Sprintf("%s-%s (%s)", Version, ShortCommit, CommitDate)
		fmt.Println("Slack Buddy version:", version)
		os.Exit(0)
	}

	systray.Run(sh.OnReady, sh.OnExit)
}
