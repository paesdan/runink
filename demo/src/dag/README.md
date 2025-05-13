# RunInk DAG Generator

This package provides functionality to build and manage a Directed Acyclic Graph (DAG) for RunInk pipeline execution.

## Overview

The DAG generator takes parsed DSL files and creates a graph structure that represents the execution flow of a RunInk pipeline. The graph consists of nodes (processing steps) and edges (dependencies between steps).

## Core Components

- **Node**: Represents a processing step in the pipeline with an ID, name, description, type, and configuration.
- **Edge**: Represents a dependency between two nodes, with a "from" node and a "to" node.
- **DAG**: The main structure that holds nodes and edges, with methods for building, validating, and traversing the graph.

## Usage

```go
import (
    "github.com/runink/runink/dag"
    "github.com/runink/runink/parser"
)

// Parse a DSL file
dslFile, err := parser.ParseDSL("path/to/file.dsl")
if err != nil {
    // Handle error
}

// Build a DAG from the parsed DSL
dag, err := dag.Build(dslFile)
if err != nil {
    // Handle error
}

// Validate the DAG
if err := dag.Validate(); err != nil {
    // Handle validation error
}

// Get execution order using topological sort
order, err := dag.TopologicalSort()
if err != nil {
    // Handle error
}

// Execute the DAG (implementation depends on your runtime)
for _, node := range order {
    // Execute node
}
```

## Step Dependencies

The DAG generator supports two types of dependency specifications in the DSL:

1. **Implicit dependencies**: Steps are connected in the order they appear in the DSL file.
2. **Explicit dependencies**: Steps can specify dependencies using:
   - `after: stepName` - Indicates that this step should run after the named step
   - `depends_on: stepName` - Alternative syntax for specifying dependencies

## Example

```go
// Sample DSL file
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
dag, err := dag.Build(dsl)
```

This will create a DAG with the following structure:
- Source node
- Step nodes for each processing step
- Sink node
- Edges connecting the nodes based on dependencies

## Validation

The DAG generator performs several validations:
- Checks for cycles in the graph
- Ensures all referenced dependencies exist
- Verifies that all nodes in edges exist in the graph
