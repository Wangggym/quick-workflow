package main

import (
	"fmt"
	"os"

	"github.com/Wangggym/quick-workflow/cmd/qk/commands"
)

var (
	// Version is set during build
	Version = "dev"
	// BuildTime is set during build
	BuildTime = "unknown"
)

func main() {
	commands.Version = Version
	commands.BuildTime = BuildTime

	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

