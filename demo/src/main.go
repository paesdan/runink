/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/runink/runink/cmd"
	"github.com/runink/runink/nodes"
)

func main() {
	// Register nodes with the engine
	nodes.IntegrateWithEngine()
	
	// Execute the command
	cmd.Execute()
}
