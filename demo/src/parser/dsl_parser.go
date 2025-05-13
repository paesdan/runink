
package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// ParseDSL parses a DSL file and returns a DSLFile struct
func ParseDSL(path string) (DSLFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return DSLFile{}, err
	}
	defer file.Close()

	dsl := DSLFile{
		Metadata:      make(map[string]interface{}),
		Steps:         []string{},
		Assertions:    []string{},
		Notifications: []Notification{},
	}

	scanner := bufio.NewScanner(file)
	
	// State tracking
	inMetadata := false
	inSteps := false
	inAssertions := false
	inGoldenTest := false
	inNotifications := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Feature detection
		if strings.HasPrefix(line, "Feature:") {
			dsl.Feature = strings.TrimSpace(strings.TrimPrefix(line, "Feature:"))
			continue
		}

		// Scenario detection
		if strings.HasPrefix(line, "Scenario:") {
			dsl.Scenario = strings.TrimSpace(strings.TrimPrefix(line, "Scenario:"))
			continue
		}

		// Metadata section
		if line == "Metadata:" {
			inMetadata = true
			inSteps = false
			inAssertions = false
			inGoldenTest = false
			inNotifications = false
			continue
		}

		// Source detection
		if strings.HasPrefix(line, "Given source") {
			sourceRegex := regexp.MustCompile(`"([^"]+)".*"([^"]+)"`)
			matches := sourceRegex.FindStringSubmatch(line)
			if len(matches) >= 3 {
				dsl.Source = matches[2] // Extract the URI
			}
			continue
		}

		// Steps section
		if line == "When events are received" {
			continue
		}

		if line == "Then:" {
			inMetadata = false
			inSteps = true
			inAssertions = false
			inGoldenTest = false
			inNotifications = false
			continue
		}

		// Assertions section
		if line == "Assertions:" {
			inMetadata = false
			inSteps = false
			inAssertions = true
			inGoldenTest = false
			inNotifications = false
			continue
		}

		// GoldenTest section
		if line == "GoldenTest:" {
			inMetadata = false
			inSteps = false
			inAssertions = false
			inGoldenTest = true
			inNotifications = false
			continue
		}

		// Notifications section
		if line == "Notifications:" {
			inMetadata = false
			inSteps = false
			inAssertions = false
			inGoldenTest = false
			inNotifications = true
			continue
		}

		// Process based on current section
		if inMetadata {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				
				// Handle array values
				if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
					arrayStr := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
					elements := strings.Split(arrayStr, ",")
					for i, e := range elements {
						elements[i] = strings.Trim(strings.TrimSpace(e), "\"")
					}
					dsl.Metadata[key] = elements
				} else {
					// Handle string values
					dsl.Metadata[key] = strings.Trim(value, "\"")
				}
			}
		} else if inSteps && strings.HasPrefix(line, "-") {
			step := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			dsl.Steps = append(dsl.Steps, step)
		} else if inAssertions && strings.HasPrefix(line, "-") {
			assertion := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			dsl.Assertions = append(dsl.Assertions, assertion)
		} else if inGoldenTest {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				
				switch key {
				case "input":
					dsl.GoldenTest.Input = value
				case "output":
					dsl.GoldenTest.Output = value
				case "validation":
					dsl.GoldenTest.Validation = value
				}
			}
		} else if inNotifications && strings.HasPrefix(line, "-") {
			notificationLine := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			parts := strings.SplitN(notificationLine, ",", 2)
			
			if len(parts) == 2 {
				conditionPart := strings.TrimSpace(parts[0])
				alertPart := strings.TrimSpace(parts[1])
				
				condition := strings.TrimSpace(strings.TrimPrefix(conditionPart, "On"))
				alert := ""
				
				alertMatch := regexp.MustCompile(`emit alert to "([^"]+)"`).FindStringSubmatch(alertPart)
				if len(alertMatch) >= 2 {
					alert = alertMatch[1]
				}
				
				dsl.Notifications = append(dsl.Notifications, Notification{
					Condition: condition,
					Alert:     alert,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return DSLFile{}, err
	}

	return dsl, nil
}
