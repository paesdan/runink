// Package nodes provides node implementations for the RunInk DAG execution engine.
package nodes

import (
	"context"
	"fmt"

	"github.com/runink/runink/dag"
	"github.com/runink/runink/internal/engine"
)

// IntegrateWithEngine integrates the CSV and Parquet nodes with the RunInk engine
func IntegrateWithEngine() {
	// Register our nodes with the engine
	fmt.Println("Integrating CSV and Parquet nodes with RunInk engine")
	
	// Register with the default registry
	RegisterWithEngine(DefaultRegistry)
	
	// Register with the engine's node registry if available
	// This would be done by the engine during initialization
	// engine.RegisterNodeType("csv_reader", CreateCSVReaderDAGNode)
	// engine.RegisterNodeType("parquet_writer", CreateParquetWriterDAGNode)
}

// CreateCSVReaderDAGNode creates a DAG node for reading CSV files
func CreateCSVReaderDAGNode(id string, config map[string]interface{}) (*dag.Node, error) {
	// Create the CSV reader node implementation
	csvReader, err := NewCSVReaderNode(id, config)
	if err != nil {
		return nil, err
	}

	// Create a DAG node that wraps our implementation
	return &dag.Node{
		ID:     id,
		Name:   "CSV Reader",
		Type:   "csv_reader",
		Config: config,
		Function: func() error {
			// This is a simplified implementation
			// In a real implementation, we would need to handle the context and channels
			_, err := csvReader.Execute(context.Background(), nil)
			return err
		},
	}, nil
}

// CreateParquetWriterDAGNode creates a DAG node for writing Parquet files
func CreateParquetWriterDAGNode(id string, config map[string]interface{}) (*dag.Node, error) {
	// Create the Parquet writer node implementation
	parquetWriter, err := NewParquetWriterNode(id, config)
	if err != nil {
		return nil, err
	}

	// Create a DAG node that wraps our implementation
	return &dag.Node{
		ID:     id,
		Name:   "Parquet Writer",
		Type:   "parquet_writer",
		Config: config,
		Function: func() error {
			// This is a simplified implementation
			// In a real implementation, we would need to handle the context and channels
			_, err := parquetWriter.Execute(context.Background(), nil)
			return err
		},
	}, nil
}

// RegisterWithEngine registers our nodes with the engine
func RegisterWithEngine(registry *NodeRegistry) {
	// Register CSV reader node
	registry.Register("csv_reader", func(id string, config map[string]interface{}) (interface{}, error) {
		return NewCSVReaderNode(id, config)
	})

	// Register Parquet writer node
	registry.Register("parquet_writer", func(id string, config map[string]interface{}) (interface{}, error) {
		return NewParquetWriterNode(id, config)
	})
}

// CreateEngineNode creates an engine node from a node implementation
func CreateEngineNode(nodeImpl interface{}, id string, dependencies []string) *engine.Node {
	// Create an engine node that wraps our implementation
	return &engine.Node{
		ID:           id,
		Dependencies: dependencies,
		Execute: func(ctx context.Context, input <-chan any) (any, error) {
			// Call the Execute method of the node implementation
			switch node := nodeImpl.(type) {
			case *CSVReaderNode:
				return node.Execute(ctx, input)
			case *ParquetWriterNode:
				return node.Execute(ctx, input)
			default:
				return nil, fmt.Errorf("unsupported node type: %T", nodeImpl)
			}
		},
	}
}
