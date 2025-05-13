// Package nodes provides node implementations for the RunInk DAG execution engine.
package nodes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

// ParquetWriterNode implements a node that writes data to a Parquet file
type ParquetWriterNode struct {
	ID          string
	Path        string
	Compression string
	Overwrite   bool
}

// NewParquetWriterNode creates a new Parquet writer node
func NewParquetWriterNode(id string, config map[string]interface{}) (*ParquetWriterNode, error) {
	// Extract path from config
	path, ok := config["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path is required for Parquet writer node")
	}

	// Extract optional parameters with defaults
	compression := "snappy"
	if compressionVal, ok := config["compression"].(string); ok {
		compression = compressionVal
	}

	overwrite := true
	if overwriteVal, ok := config["overwrite"]; ok {
		if overwriteBool, ok := overwriteVal.(bool); ok {
			overwrite = overwriteBool
		} else if overwriteStr, ok := overwriteVal.(string); ok {
			overwrite = strings.ToLower(overwriteStr) == "true"
		}
	}

	// Create the node
	return &ParquetWriterNode{
		ID:          id,
		Path:        path,
		Compression: compression,
		Overwrite:   overwrite,
	}, nil
}

// Execute writes the input data to a Parquet file
func (n *ParquetWriterNode) Execute(ctx context.Context, input <-chan any) (any, error) {
	// Get the input data
	var data *CSVData
	for item := range input {
		if csvData, ok := item.(*CSVData); ok {
			data = csvData
			break
		} else if packet, ok := item.(map[string]interface{}); ok {
			// Try to extract CSVData from a data packet
			if csvData, ok := packet["data"].(*CSVData); ok {
				data = csvData
				break
			}
		}
	}

	if data == nil {
		return nil, fmt.Errorf("no CSV data received")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(n.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists and handle overwrite
	if _, err := os.Stat(n.Path); err == nil {
		if !n.Overwrite {
			return nil, fmt.Errorf("file already exists and overwrite is disabled")
		}
		// Remove existing file
		if err := os.Remove(n.Path); err != nil {
			return nil, fmt.Errorf("failed to remove existing file: %w", err)
		}
	}

	// Create the Parquet file
	file, err := os.Create(n.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create Parquet file: %w", err)
	}
	defer file.Close()

	// Create a Parquet writer
	pw, err := createParquetWriter(file, data, n.Compression)
	if err != nil {
		return nil, fmt.Errorf("failed to create Parquet writer: %w", err)
	}

	// Write the data
	for _, row := range data.Rows {
		record := make(map[string]interface{})
		for i, header := range data.Headers {
			if i < len(row) {
				record[header] = row[i]
			}
		}
		if err := pw.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write record: %w", err)
		}
	}

	// Flush and close the writer
	if err := pw.WriteStop(); err != nil {
		return nil, fmt.Errorf("failed to finalize Parquet file: %w", err)
	}

	return fmt.Sprintf("Successfully wrote %d rows to %s", len(data.Rows), n.Path), nil
}

// createParquetWriter creates a Parquet writer with the appropriate schema
func createParquetWriter(file *os.File, data *CSVData, compression string) (*writer.JSONWriter, error) {
	// Create a schema based on the CSV data
	schema := make(map[string]string)
	for i, header := range data.Headers {
		// Get type from schema or infer from first row
		dataType, ok := data.Schema[header]
		if !ok && len(data.Rows) > 0 && i < len(data.Rows[0]) {
			// Infer type from first non-nil value
			for _, row := range data.Rows {
				if i < len(row) && row[i] != nil {
					dataType = getParquetType(reflect.TypeOf(row[i]).String())
					break
				}
			}
			if dataType == "" {
				dataType = "BYTE_ARRAY" // Default to string
			}
		} else if !ok {
			dataType = "BYTE_ARRAY" // Default to string
		} else {
			// Convert from CSV schema type to Parquet type
			dataType = convertToParquetType(dataType)
		}
		schema[header] = dataType
	}

	// Set compression
	var compressionType parquet.CompressionCodec
	switch strings.ToLower(compression) {
	case "snappy":
		compressionType = parquet.CompressionCodec_SNAPPY
	case "gzip":
		compressionType = parquet.CompressionCodec_GZIP
	case "lz4":
		compressionType = parquet.CompressionCodec_LZ4
	case "zstd":
		compressionType = parquet.CompressionCodec_ZSTD
	default:
		compressionType = parquet.CompressionCodec_SNAPPY
	}

	// Create JSON writer
	pw, err := writer.NewJSONWriter(schema, file, 4, compressionType)
	if err != nil {
		return nil, err
	}

	return pw, nil
}

// getParquetType converts a Go type to a Parquet type
func getParquetType(goType string) string {
	switch goType {
	case "int", "int8", "int16", "int32", "int64":
		return "INT64"
	case "float32", "float64":
		return "DOUBLE"
	case "bool":
		return "BOOLEAN"
	case "time.Time":
		return "INT96" // Parquet's timestamp type
	default:
		return "BYTE_ARRAY"
	}
}

// convertToParquetType converts a CSV schema type to a Parquet type
func convertToParquetType(csvType string) string {
	switch strings.ToLower(csvType) {
	case "integer", "int":
		return "INT64"
	case "float", "double":
		return "DOUBLE"
	case "boolean", "bool":
		return "BOOLEAN"
	case "date", "timestamp":
		return "INT96" // Parquet's timestamp type
	default:
		return "BYTE_ARRAY"
	}
}
