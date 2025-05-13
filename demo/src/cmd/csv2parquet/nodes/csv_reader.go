// Package nodes provides node implementations for the CSV to PARQUET conversion.
package nodes

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// CSVData represents the data read from a CSV file
type CSVData struct {
	Headers []string
	Rows    [][]interface{}
	Schema  map[string]string // Maps column names to their types
}

// CSVReaderNode implements a node that reads data from a CSV file
type CSVReaderNode struct {
	ID          string
	Path        string
	HasHeader   bool
	Delimiter   rune
	InferSchema bool
}

// NewCSVReaderNode creates a new CSV reader node
func NewCSVReaderNode(id string, config map[string]interface{}) (*CSVReaderNode, error) {
	// Extract path from config
	path, ok := config["path"].(string)
	if !ok {
		return nil, fmt.Errorf("path is required for CSV reader node")
	}

	// Extract optional parameters with defaults
	hasHeader := true
	if headerVal, ok := config["header"]; ok {
		if headerBool, ok := headerVal.(bool); ok {
			hasHeader = headerBool
		} else if headerStr, ok := headerVal.(string); ok {
			hasHeader = strings.ToLower(headerStr) == "true"
		}
	}

	delimiter := ','
	if delimiterVal, ok := config["delimiter"].(string); ok && len(delimiterVal) > 0 {
		delimiter = rune(delimiterVal[0])
	}

	inferSchema := false
	if inferSchemaVal, ok := config["inferSchema"]; ok {
		if inferSchemaBool, ok := inferSchemaVal.(bool); ok {
			inferSchema = inferSchemaBool
		} else if inferSchemaStr, ok := inferSchemaVal.(string); ok {
			inferSchema = strings.ToLower(inferSchemaStr) == "true"
		}
	}

	// Create the node
	return &CSVReaderNode{
		ID:          id,
		Path:        path,
		HasHeader:   hasHeader,
		Delimiter:   delimiter,
		InferSchema: inferSchema,
	}, nil
}

// Execute reads data from the CSV file and returns it as a CSVData struct
func (n *CSVReaderNode) Execute(ctx context.Context, _ <-chan any) (any, error) {
	// Open the CSV file
	file, err := os.Open(n.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)
	reader.Comma = n.Delimiter

	// Read the header row if present
	var headers []string
	if n.HasHeader {
		headers, err = reader.Read()
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV header: %w", err)
		}
	}

	// Initialize schema
	schema := make(map[string]string)

	// Read all rows
	var rows [][]interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		// If we don't have headers yet, use the first row as headers
		if headers == nil {
			headers = make([]string, len(record))
			for i := range record {
				headers[i] = fmt.Sprintf("col%d", i+1)
			}
		}

		// Convert strings to appropriate types based on schema or inference
		row := make([]interface{}, len(record))
		for i, value := range record {
			if i >= len(headers) {
				continue
			}

			columnName := headers[i]
			columnType, hasType := schema[columnName]

			if !hasType && n.InferSchema {
				// Infer type if not in schema
				columnType = inferType(value)
				schema[columnName] = columnType
			}

			// Convert value based on type
			row[i] = convertValue(value, columnType)
		}

		rows = append(rows, row)
	}

	// Create and return the CSV data
	csvData := &CSVData{
		Headers: headers,
		Rows:    rows,
		Schema:  schema,
	}

	return csvData, nil
}

// inferType attempts to infer the data type of a string value
func inferType(value string) string {
	// Check if empty
	if value == "" {
		return "string"
	}

	// Check if boolean
	if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
		return "boolean"
	}

	// Check if integer
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return "integer"
	}

	// Check if float
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "float"
	}

	// Check if date (simple check for YYYY-MM-DD format)
	if _, err := time.Parse("2006-01-02", value); err == nil {
		return "date"
	}

	// Default to string
	return "string"
}

// convertValue converts a string value to the appropriate type
func convertValue(value string, dataType string) interface{} {
	if value == "" {
		return nil
	}

	switch strings.ToLower(dataType) {
	case "integer", "int":
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	case "float", "double":
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	case "boolean", "bool":
		if strings.ToLower(value) == "true" {
			return true
		} else if strings.ToLower(value) == "false" {
			return false
		}
	case "date":
		if t, err := time.Parse("2006-01-02", value); err == nil {
			return t
		}
	}

	// Default to string
	return value
}
