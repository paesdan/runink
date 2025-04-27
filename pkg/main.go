package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"yourproject/barnctl"
	"yourproject/runinkctl"
	"yourproject/herdctl"
	"yourproject/runictl"
)

var verbose bool

func main() {
	rootCmd := &cobra.Command{
		Use:   "runi",
		Short: "Runink CLI - Manage Herds, Compile Pipelines, Execute Slices",
		Long: logo() + `
Runink is a Go-native, declarative platform to manage distributed
data pipelines, herds (namespaces), agents, secrets, and lineage governance.

Key domains:
- barnctl  → Herd, Storage, Secrets Control Plane
- runinkctl → Compile Features + Contracts into DAGs
- herdctl  → Governance, Metadata, Redactions
- runictl  → Schedule and Execute secure Slice workers
`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global persistent flags
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging output")

	// Add all subcommands
	rootCmd.AddCommand(
		barnctl.NewBarnctlCommand(),
		runinkctl.NewruninkctlCommand(),
		herdctl.NewHerdctlCommand(),
		runictl.NewRunictlCommand(),
	)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}

// logo returns an ASCII art version of the logo
func logo() string {
	return `
          .-\"\"\"-.
         / -   -  \
        |  .-. .- |
        |  \o| |o (
        \     ^    \
         '.  )--'  /
          '-'---'-'
    Runink — Herds. Pipelines. Slices. Freedom.
`
}
