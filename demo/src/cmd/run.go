/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute a RunInk DAG based on input files",
	Long: `The run command executes a Directed Acyclic Graph (DAG) based on the provided
input files such as contract, conf, dsl, and herd files.

Example usage:
  runink run --contract file.contract --conf file.conf --dsl file.dsl --herd file.herd
  
Advanced usage:
  runink run --contract file.contract --conf file.conf --dsl file.dsl --error-mode continue --data-pass json --monitoring verbose --execution-mode async`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("RunInk – executing DAG")
		if err := executeDAG(cmd); err != nil {
			fmt.Printf("Error executing DAG: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Define flags for input file paths
	runCmd.Flags().String("contract", "", "Path to the contract file (.contract)")
	runCmd.Flags().String("conf", "", "Path to the configuration file (.conf)")
	runCmd.Flags().String("dsl", "", "Path to the domain specific language file (.dsl)")
	runCmd.Flags().String("herd", "", "Path to the herd file (.herd)")
	
	// Basic flags
	runCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	
	// Enhanced engine flags
	runCmd.Flags().String("error-mode", "stop", "Error handling mode: 'stop' (default) or 'continue'")
	runCmd.Flags().String("data-pass", "json", "Data passing strategy: 'json' (default), 'file', or 'stdout'")
	runCmd.Flags().String("monitoring", "basic", "Monitoring level: 'none', 'basic' (default), or 'verbose'")
	runCmd.Flags().String("execution-mode", "sync", "Execution mode: 'sync' (default) or 'async'")
	runCmd.Flags().Bool("integrated-flow", true, "Enable integrated execution flow (default: true)")
	runCmd.Flags().String("isolation-level", "none", "Isolation level: 'none' (default), 'process', or 'container'")
	runCmd.Flags().String("run-id", "", "Unique identifier for this run (auto-generated if not provided)")
}
