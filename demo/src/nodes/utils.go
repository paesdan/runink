
// Package nodes provides node implementations for the RunInk DAG execution engine.
package nodes

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseValue parses a string value into the appropriate type based on the given data type
func ParseValue(value string, dataType string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}

	switch strings.ToLower(dataType) {
	case "integer", "int":
		return strconv.ParseInt(value, 10, 64)
	case "float", "double":
		return strconv.ParseFloat(value, 64)
	case "boolean", "bool":
		return strconv.ParseBool(value)
	case "date":
		return time.Parse("2006-01-02", value)
	case "timestamp":
		return time.Parse(time.RFC3339, value)
	default:
		return value, nil
	}
}

// FormatValue formats a value as a string based on its type
func FormatValue(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case time.Time:
		return v.Format("2006-01-02")
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetTypeFromValue returns the data type of a value
func GetTypeFromValue(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch value.(type) {
	case int, int8, int16, int32, int64:
		return "integer"
	case float32, float64:
		return "float"
	case bool:
		return "boolean"
	case time.Time:
		return "date"
	case string:
		return "string"
	default:
		return "string"
	}
}

// IsNumeric checks if a string value is numeric
func IsNumeric(value string) bool {
	_, err1 := strconv.ParseInt(value, 10, 64)
	_, err2 := strconv.ParseFloat(value, 64)
	return err1 == nil || err2 == nil
}

// IsBoolean checks if a string value is a boolean
func IsBoolean(value string) bool {
	_, err := strconv.ParseBool(value)
	return err == nil
}

// IsDate checks if a string value is a date in YYYY-MM-DD format
func IsDate(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

// IsTimestamp checks if a string value is a timestamp in RFC3339 format
func IsTimestamp(value string) bool {
	_, err := time.Parse(time.RFC3339, value)
	return err == nil
}
