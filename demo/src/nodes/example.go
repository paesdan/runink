// Package nodes provides node implementations for the RunInk DAG execution engine.
package nodes

import (
	"context"
	"fmt"
	"time"

	"github.com/runink/runink/internal/engine"
)

// RunCSVToParquetExample demonstrates how to use the CSV reader and Parquet writer nodes
func RunCSVToParquetExample(csvPath, parquetPath string) error {
	// Create a channel manager for data flow
	channelManager := engine.NewChannelManager()

	// Create a CSV reader node
	csvReader, err := NewCSVReaderNode("csv_reader", map[string]interface{}{
		"path":        csvPath,
		"header":      true,
		"inferSchema": true,
	})
	if err != nil {
		return fmt.Errorf("failed to create CSV reader node: %w", err)
	}

	// Create a Parquet writer node
	parquetWriter, err := NewParquetWriterNode("parquet_writer", map[string]interface{}{
		"path":        parquetPath,
		"compression": "snappy",
		"overwrite":   true,
	})
	if err != nil {
		return fmt.Errorf("failed to create Parquet writer node: %w", err)
	}

	// Create a data edge between the nodes
	edge := &engine.DataEdge{
		From:        "csv_reader",
		To:          "parquet_writer",
		BufferSize:  1,
		CodecName:   "json",
		RetryPolicy: engine.DefaultRetryPolicy,
	}

	// Register the edge with the channel manager
	channelManager.RegisterEdge(edge)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Execute the CSV reader node
	fmt.Println("Reading CSV file:", csvPath)
	csvData, err := csvReader.Execute(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to execute CSV reader node: %w", err)
	}

	// Send the CSV data to the Parquet writer node
	ch, err := channelManager.GetChannel("csv_reader", "parquet_writer")
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}

	// Send the CSV data to the channel
	ch <- csvData
	close(ch)

	// Execute the Parquet writer node
	fmt.Println("Writing Parquet file:", parquetPath)
	result, err := parquetWriter.Execute(ctx, ch)
	if err != nil {
		return fmt.Errorf("failed to execute Parquet writer node: %w", err)
	}

	fmt.Println("Result:", result)
	return nil
}

// RunCSVToParquetWithDAG demonstrates how to use the CSV reader and Parquet writer nodes with a DAG
func RunCSVToParquetWithDAG(csvPath, parquetPath string) error {
	// Create nodes for the DAG
	nodes := []*engine.Node{
		{
			ID: "csv_reader",
			Execute: func(ctx context.Context, _ <-chan any) (any, error) {
				csvReader, err := NewCSVReaderNode("csv_reader", map[string]interface{}{
					"path":        csvPath,
					"header":      true,
					"inferSchema": true,
				})
				if err != nil {
					return nil, err
				}
				return csvReader.Execute(ctx, nil)
			},
			Dependencies: []string{},
		},
		{
			ID: "parquet_writer",
			Execute: func(ctx context.Context, input <-chan any) (any, error) {
				parquetWriter, err := NewParquetWriterNode("parquet_writer", map[string]interface{}{
					"path":        parquetPath,
					"compression": "snappy",
					"overwrite":   true,
				})
				if err != nil {
					return nil, err
				}
				return parquetWriter.Execute(ctx, input)
			},
			Dependencies: []string{"csv_reader"},
		},
	}

	// Build the DAG
	dag := engine.BuildDAG(nodes, false, "csv_to_parquet")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create a monitor
	monitor := engine.NewDefaultMonitor()

	// Create an error handler
	errorHandler := engine.NewDefaultErrorHandler(true)

	// Execute the DAG
	fmt.Println("Executing CSV to Parquet conversion DAG")
	err := engine.Execute(ctx, dag, engine.SyncMode, monitor, errorHandler)
	if err != nil {
		return fmt.Errorf("failed to execute DAG: %w", err)
	}

	fmt.Println("CSV to Parquet conversion completed successfully")
	return nil
}
