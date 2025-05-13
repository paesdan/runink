
package engine

import (
        "fmt"

        "github.com/runink/runink/dag"
        "github.com/runink/runink/parser"
)

// DAGConfig holds the configuration for the DAG execution
type DAGConfig struct {
        ContractFile string
        ConfFile     string
        DSLFile      string
        HerdFile     string
        Verbose      bool
}

// ExecuteDAG executes the directed acyclic graph (DAG) of tasks
func ExecuteDAG(config DAGConfig) error {
        // Print verbose information if enabled
        if config.Verbose {
                println("Executing DAG")
                println("Contract file:", config.ContractFile)
                println("Conf file:", config.ConfFile)
                println("DSL file:", config.DSLFile)
                if config.HerdFile != "" {
                        println("Herd file:", config.HerdFile)
                }
        }
        
        // 1. Parse the DSL file
        dslData, err := parser.ParseDSL(config.DSLFile)
        if err != nil {
                return fmt.Errorf("failed to parse DSL file: %v", err)
        }

        // 2. Parse contract file if provided
        var contractData parser.ContractFile
        if config.ContractFile != "" {
                contractData, err = parser.ParseContract(config.ContractFile)
                if err != nil {
                        return fmt.Errorf("failed to parse contract file: %v", err)
                }
        }

        // 3. Parse conf file if provided
        var confData parser.ConfFile
        if config.ConfFile != "" {
                confData, err = parser.ParseConf(config.ConfFile)
                if err != nil {
                        return fmt.Errorf("failed to parse conf file: %v", err)
                }
        }

        // 4. Parse herd file if provided
        var herdData parser.HerdFile
        if config.HerdFile != "" {
                herdData, err = parser.ParseHerd(config.HerdFile)
                if err != nil {
                        return fmt.Errorf("failed to parse herd file: %v", err)
                }
        }

        if config.Verbose {
                fmt.Println("Building DAG...")
        }

        // 5. Build the DAG from the parsed DSL
        dagInstance, err := dag.Build(dslData)
        if err != nil {
                return fmt.Errorf("failed to build DAG: %v", err)
        }

        if config.Verbose {
                fmt.Println("DAG structure:")
                fmt.Println(dagInstance.String())
        }

        // 6. Validate the DAG
        if err := dagInstance.Validate(); err != nil {
                return fmt.Errorf("DAG validation failed: %v", err)
        }

        // 7. Create execution config
        execConfig := ExecutionConfig{
                DAG:         dagInstance,
                Contract:    contractData,
                Conf:        confData,
                Herd:        herdData,
                Verbose:     config.Verbose,
        }

        // 8. Execute the DAG
        if config.Verbose {
                fmt.Println("Executing DAG...")
        }
        
        result, err := Execute(execConfig)
        if err != nil {
                return err
        }
        
        // We're ignoring the result here since the original function only returns error
        // In a real implementation, you might want to store or log the result
        if config.Verbose {
                fmt.Println(FormatExecutionSummary(result))
        }
        
        return nil
}
