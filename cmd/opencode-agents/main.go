package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opencode/opencode-agents/internal/cli"
)

func main() {
	var (
		globalFlag  bool
		localFlag   bool
		syncAllFlag bool
		dirFlag     string
	)

	// Parse flags
	flag.BoolVar(&globalFlag, "global", false, "Update the global opencode.json (~/.config/opencode/opencode.json)")
	flag.BoolVar(&localFlag, "local", false, "Update the local opencode.json (./opencode.json)")
	flag.BoolVar(&syncAllFlag, "sync-all", false, "Bypass TUI and sync all discovered agents")
	flag.StringVar(&dirFlag, "dir", "", "Path to the agents directory (defaults to $OPENCODE_AGENTS_DIR or ~/.config/opencode-agents/agents)")

	flag.Parse()

	// Prevent conflicting flags
	if globalFlag && localFlag {
		fmt.Println("Error: Cannot use both --global and --local flags simultaneously.")
		os.Exit(1)
	}

	app := cli.NewCLI(globalFlag, localFlag, syncAllFlag, dirFlag)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
