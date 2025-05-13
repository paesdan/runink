package main

import (
	"fmt"
	"os"

	"github.com/runink/runink/dag"
	"github.com/runink/runink/parser"
)

func main() {
	// Create a sample DSL file for demonstration
	dsl := parser.DSLFile{
		Feature:  "Demo Feature",
		Scenario: "Demo Scenario",
		Source:   "kafka://demo-topic",
		Steps: []string{
			"filter removeNulls (field: \"value\")",
			"transform normalizeData (format: \"json\")",
			"aggregate countByKey (key: \"user_id\", window: \"5m\")",
			"enrich addMetadata (source: \"metadata-service\", depends_on: normalizeData)",
			"transform finalTransform (after: countByKey, after: addMetadata)",
		},
	}

	// Build the DAG
	dag, err := dag.Build(dsl)
	if err != nil {
		fmt.Printf("Error building DAG: %v\n", err)
		os.Exit(1)
	}

	// Validate the DAG
	if err := dag.Validate(); err != nil {
		fmt.Printf("DAG validation failed: %v\n", err)
		os.Exit(1)
	}

	// Print the DAG structure
	fmt.Println("DAG Structure:")
	fmt.Println(dag.String())

	// Get topological sort order
	order, err := dag.TopologicalSort()
	if err != nil {
		fmt.Printf("Topological sort failed: %v\n", err)
		os.Exit(1)
	}

	// Print execution order
	fmt.Println("\nExecution Order:")
	for i, node := range order {
		fmt.Printf("%d. %s (%s)\n", i+1, node.Name, node.Type)
	}
}
