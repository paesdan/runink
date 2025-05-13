package cmd

import (
	"context"
	"fmt"
	"time"
	
	"github.com/google/uuid"
	"github.com/runink/runink/internal/engine"
	"github.com/spf13/cobra"
)

// getRunFlags retrieves the flags from the run command
func getRunFlags(cmd *cobra.Command) (engine.RunConfig, error) {
	config := engine.RunConfig{}
	
	// Get basic file paths
	contractFile, err := cmd.Flags().GetString("contract")
	if err != nil {
		return config, err
	}
	
	confFile, err := cmd.Flags().GetString("conf")
	if err != nil {
		return config, err
	}
	
	dslFile, err := cmd.Flags().GetString("dsl")
	if err != nil {
		return config, err
	}
	
	herdFile, err := cmd.Flags().GetString("herd")
	if err != nil {
		return config, err
	}
	
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return config, err
	}
	
	// Get enhanced engine flags
	errorMode, err := cmd.Flags().GetString("error-mode")
	if err != nil {
		return config, err
	}
	
	dataPass, err := cmd.Flags().GetString("data-pass")
	if err != nil {
		return config, err
	}
	
	monitoring, err := cmd.Flags().GetString("monitoring")
	if err != nil {
		return config, err
	}
	
	executionMode, err := cmd.Flags().GetString("execution-mode")
	if err != nil {
		return config, err
	}
	
	integratedFlow, err := cmd.Flags().GetBool("integrated-flow")
	if err != nil {
		return config, err
	}
	
	isolationLevel, err := cmd.Flags().GetString("isolation-level")
	if err != nil {
		return config, err
	}
	
	runID, err := cmd.Flags().GetString("run-id")
	if err != nil {
		return config, err
	}
	
	// Generate run ID if not provided
	if runID == "" {
		runID = fmt.Sprintf("run-%s", uuid.New().String()[:8])
	}
	
	// Populate the config
	config = engine.RunConfig{
		Files: engine.FileConfig{
			ContractFile: contractFile,
			ConfFile:     confFile,
			DSLFile:      dslFile,
			HerdFile:     herdFile,
		},
		Execution: engine.ExecutionConfig{
			ErrorMode:      parseErrorMode(errorMode),
			DataPass:       parseDataPass(dataPass),
			Monitoring:     parseMonitoring(monitoring),
			ExecutionMode:  parseExecutionMode(executionMode),
			IntegratedFlow: integratedFlow,
			IsolationLevel: parseIsolationLevel(isolationLevel),
			Verbose:        verbose,
			RunID:          runID,
		},
	}
	
	return config, nil
}

// Helper functions to parse string flags into engine types
func parseErrorMode(mode string) engine.ErrorMode {
	switch mode {
	case "continue":
		return engine.ContinueOnError
	default:
		return engine.StopOnError
	}
}

func parseDataPass(strategy string) engine.DataPassStrategy {
	switch strategy {
	case "file":
		return engine.FileDataPass
	case "stdout":
		return engine.StdoutDataPass
	default:
		return engine.JSONDataPass
	}
}

func parseMonitoring(level string) engine.MonitoringLevel {
	switch level {
	case "none":
		return engine.NoMonitoring
	case "verbose":
		return engine.VerboseMonitoring
	default:
		return engine.BasicMonitoring
	}
}

func parseExecutionMode(mode string) engine.ExecutionMode {
	switch mode {
	case "async":
		return engine.AsyncMode
	default:
		return engine.SyncMode
	}
}

func parseIsolationLevel(level string) engine.IsolationLevel {
	switch level {
	case "process":
		return engine.ProcessIsolation
	case "container":
		return engine.ContainerIsolation
	default:
		return engine.NoIsolation
	}
}

// executeDAG is a wrapper function that calls the enhanced engine's Execute function
func executeDAG(cmd *cobra.Command) error {
	// Get configuration from flags
	config, err := getRunFlags(cmd)
	if err != nil {
		return err
	}
	
	// Validate required files
	if config.Files.ContractFile == "" || config.Files.ConfFile == "" || config.Files.DSLFile == "" {
		return fmt.Errorf("contract, conf, and dsl files are required")
	}
	
	if config.Execution.Verbose {
		fmt.Printf("RunInk Execution - Run ID: %s\n", config.Execution.RunID)
		fmt.Printf("Processing files:\n")
		fmt.Printf("  Contract: %s\n", config.Files.ContractFile)
		fmt.Printf("  Conf: %s\n", config.Files.ConfFile)
		fmt.Printf("  DSL: %s\n", config.Files.DSLFile)
		if config.Files.HerdFile != "" {
			fmt.Printf("  Herd: %s\n", config.Files.HerdFile)
		}
		fmt.Printf("Execution configuration:\n")
		fmt.Printf("  Error Mode: %v\n", config.Execution.ErrorMode)
		fmt.Printf("  Data Pass Strategy: %v\n", config.Execution.DataPass)
		fmt.Printf("  Monitoring Level: %v\n", config.Execution.Monitoring)
		fmt.Printf("  Execution Mode: %v\n", config.Execution.ExecutionMode)
		fmt.Printf("  Integrated Flow: %v\n", config.Execution.IntegratedFlow)
		fmt.Printf("  Isolation Level: %v\n", config.Execution.IsolationLevel)
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	
	// Parse files and build DAG
	dag, err := engine.BuildDAG(config.Files, config.Execution.IntegratedFlow, config.Execution.RunID)
	if err != nil {
		return fmt.Errorf("failed to build DAG: %w", err)
	}
	
	// Create monitor based on monitoring level
	var monitor engine.Monitor
	switch config.Execution.Monitoring {
	case engine.NoMonitoring:
		monitor = engine.NewNullMonitor()
	case engine.VerboseMonitoring:
		monitor = engine.NewVerboseMonitor(config.Execution.RunID)
	default:
		monitor = engine.NewDefaultMonitor()
	}
	
	// Create error handler
	errorHandler := engine.NewDefaultErrorHandler(config.Execution.ErrorMode == engine.StopOnError)
	
	// Execute the DAG with the enhanced engine
	startTime := time.Now()
	err = engine.Execute(
		ctx,
		dag,
		config.Execution.ExecutionMode,
		monitor,
		errorHandler,
		engine.WithDataPassStrategy(config.Execution.DataPass),
		engine.WithIsolationLevel(config.Execution.IsolationLevel),
	)
	
	if config.Execution.Verbose {
		fmt.Printf("Execution completed in %v\n", time.Since(startTime))
	}
	
	return err
}
