// Package dag provides functionality to build and manage a Directed Acyclic Graph (DAG)
package dag

import (
	"fmt"

	"github.com/runink/runink/parser"
)

// Example demonstrates how to use the DAG package
func Example() {
	// Create a sample DSL file (normally this would be parsed from a file)
	dsl := parser.DSLFile{
		Feature:  "Example Feature",
		Scenario: "Example Scenario",
		Source:   "kafka://example-topic",
		Steps: []string{
			"filter removeNulls (field: \"value\")",
			"transform normalizeData (format: \"json\")",
			"aggregate countByKey (key: \"user_id\", window: \"5m\")",
			"enrich addMetadata (source: \"metadata-service\", depends_on: normalizeData)",
			"transform finalTransform (after: countByKey, after: addMetadata)",
		},
	}

	// Build the DAG
	dag, err := Build(dsl)
	if err != nil {
		fmt.Printf("Error building DAG: %v\n", err)
		return
	}

	// Validate the DAG
	if err := dag.Validate(); err != nil {
		fmt.Printf("DAG validation failed: %v\n", err)
		return
	}

	// Print the DAG structure
	fmt.Println(dag.String())

	// Get topological sort order
	order, err := dag.TopologicalSort()
	if err != nil {
		fmt.Printf("Topological sort failed: %v\n", err)
		return
	}

	// Print execution order
	fmt.Println("Execution order:")
	for i, node := range order {
		fmt.Printf("%d. %s (%s)\n", i+1, node.Name, node.Type)
	}
}
