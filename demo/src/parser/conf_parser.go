
package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// ParseConf parses a conf file and returns a ConfFile struct
func ParseConf(path string) (ConfFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return ConfFile{}, err
	}
	defer file.Close()

	conf := ConfFile{
		Config: make(map[string]interface{}),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		
		// Try to convert to appropriate type
		if intVal, err := strconv.Atoi(value); err == nil {
			conf.Config[key] = intVal
		} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			conf.Config[key] = floatVal
		} else if value == "true" || value == "false" {
			conf.Config[key] = value == "true"
		} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			// Parse array
			arrayStr := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
			elements := strings.Split(arrayStr, ",")
			for i, e := range elements {
				elements[i] = strings.Trim(strings.TrimSpace(e), "\"")
			}
			conf.Config[key] = elements
		} else {
			conf.Config[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return ConfFile{}, err
	}

	return conf, nil
}
