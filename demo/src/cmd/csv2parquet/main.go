package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/runink/runink/cmd/csv2parquet/nodes"
)

func main() {
	// Parse command-line flags
	csvPath := flag.String("csv", "", "Path to the CSV file (required)")
	parquetPath := flag.String("parquet", "", "Path to the output PARQUET file (defaults to CSV path with .parquet extension)")
	useDag := flag.Bool("use-dag", false, "Use the DAG execution engine (default: false)")
	flag.Parse()

	// Validate flags
	if *csvPath == "" {
		fmt.Println("Error: csv path is required")
		flag.Usage()
		os.Exit(1)
	}

	// If parquet path is not provided, use the same path as the CSV file but with .parquet extension
	if *parquetPath == "" {
		dir, file := filepath.Split(*csvPath)
		ext := filepath.Ext(file)
		name := file[:len(file)-len(ext)]
		*parquetPath = filepath.Join(dir, name+".parquet")
	}

	// Run the example
	fmt.Printf("Converting CSV file %s to PARQUET file %s\n", *csvPath, *parquetPath)
	
	var err error
	if *useDag {
		err = runWithDAG(*csvPath, *parquetPath)
	} else {
		err = runSimple(*csvPath, *parquetPath)
	}
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Conversion completed successfully")
}

func runSimple(csvPath, parquetPath string) error {
	// Create a CSV reader node
	csvReader, err := nodes.NewCSVReaderNode("csv_reader", map[string]interface{}{
		"path":        csvPath,
		"header":      true,
		"inferSchema": true,
	})
	if err != nil {
		return fmt.Errorf("failed to create CSV reader node: %w", err)
	}

	// Create a Parquet writer node
	parquetWriter, err := nodes.NewParquetWriterNode("parquet_writer", map[string]interface{}{
		"path":        parquetPath,
		"compression": "snappy",
		"overwrite":   true,
	})
	if err != nil {
		return fmt.Errorf("failed to create Parquet writer node: %w", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Execute the CSV reader node
	fmt.Println("Reading CSV file:", csvPath)
	csvData, err := csvReader.Execute(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to execute CSV reader node: %w", err)
	}

	// Create a channel to pass data between nodes
	ch := make(chan any, 1)
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

func runWithDAG(csvPath, parquetPath string) error {
	fmt.Println("DAG execution not implemented in this simplified version")
	return runSimple(csvPath, parquetPath)
}
