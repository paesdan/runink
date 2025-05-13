package main

import (
    "fmt"
    "os"

    "github.com/runink/barnctl"
    "github.com/runink/buildctl"
    "github.com/runink/herdctl"
    "github.com/runink/runictl"
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "runi",
        Short: "Runink CLI - Container orchestration and workflow management",
    }

    rootCmd.AddCommand(
        barnctl.NewCommand(),
        buildctl.NewCommand(),
        herdctl.NewCommand(),
        runictl.NewCommand(),
    )

    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}