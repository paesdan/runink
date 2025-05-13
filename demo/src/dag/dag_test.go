
package dag

import (
        "testing"

        "github.com/runink/runink/parser"
)

// TestBuildSimpleDAG tests building a simple DAG from a DSL file
func TestBuildSimpleDAG(t *testing.T) {
        // Create a simple DSL file
        dsl := parser.DSLFile{
                Feature:  "Test Feature",
                Scenario: "Test Scenario",
                Source:   "test://source",
                Steps: []string{
                        "transform step1 (param1: value1)",
                        "filter step2 (param2: value2, after: step1)",
                        "aggregate step3 (depends_on: step2)",
                },
        }

        // Build the DAG
        dag, err := Build(dsl)
        if err != nil {
                t.Fatalf("Failed to build DAG: %v", err)
        }

        // Validate the DAG
        if err := dag.Validate(); err != nil {
                t.Fatalf("DAG validation failed: %v", err)
        }

        // Check number of nodes (source + 3 steps + sink)
        if len(dag.Nodes) != 5 {
                t.Errorf("Expected 5 nodes, got %d", len(dag.Nodes))
        }

        // Check number of edges (source->step1, step1->step2, step2->step3, step3->sink)
        if len(dag.Edges) != 4 {
                t.Errorf("Expected 4 edges, got %d", len(dag.Edges))
        }

        // Check topological sort
        order, err := dag.TopologicalSort()
        if err != nil {
                t.Fatalf("Topological sort failed: %v", err)
        }

        // Verify the order contains all nodes
        if len(order) != 5 {
                t.Fatalf("Expected 5 nodes in topological order, got %d", len(order))
        }

        // Check that the order is valid (source should be first, sink should be last)
        if order[0].ID != "source" {
                t.Errorf("Expected source to be first in topological order, got %s", order[0].ID)
        }
        
        if order[len(order)-1].ID != "sink" {
                t.Errorf("Expected sink to be last in topological order, got %s", order[len(order)-1].ID)
        }
        
        // Check that step1 comes before step2 and step2 comes before step3
        step1Pos, step2Pos, step3Pos := -1, -1, -1
        for i, node := range order {
                if node.ID == "step_0" { // step1
                        step1Pos = i
                } else if node.ID == "step_1" { // step2
                        step2Pos = i
                } else if node.ID == "step_2" { // step3
                        step3Pos = i
                }
        }
        
        if step1Pos == -1 || step2Pos == -1 || step3Pos == -1 {
                t.Errorf("Not all steps found in topological order")
        } else {
                if step1Pos >= step2Pos {
                        t.Errorf("step1 should come before step2 in topological order")
                }
                if step2Pos >= step3Pos {
                        t.Errorf("step2 should come before step3 in topological order")
                }
        }
}

// TestCyclicDAG tests that a cyclic DAG is detected and rejected
func TestCyclicDAG(t *testing.T) {
        // Create a DAG with a cycle
        dag := NewDAG()

        // Add nodes
        dag.AddNode(&Node{ID: "A", Name: "Node A"})
        dag.AddNode(&Node{ID: "B", Name: "Node B"})
        dag.AddNode(&Node{ID: "C", Name: "Node C"})

        // Add edges to create a cycle: A -> B -> C -> A
        dag.AddEdge("A", "B")
        dag.AddEdge("B", "C")
        dag.AddEdge("C", "A")

        // Validate should fail
        err := dag.Validate()
        if err == nil {
                t.Error("Expected validation to fail for cyclic DAG, but it passed")
        }

        // Topological sort should also fail
        _, err = dag.TopologicalSort()
        if err == nil {
                t.Error("Expected topological sort to fail for cyclic DAG, but it passed")
        }
}

// TestEmptyDAG tests building a DAG with no steps
func TestEmptyDAG(t *testing.T) {
        // Create a DSL file with no steps
        dsl := parser.DSLFile{
                Feature:  "Empty Feature",
                Scenario: "Empty Scenario",
                Source:   "test://source",
                Steps:    []string{},
        }

        // Build the DAG
        dag, err := Build(dsl)
        if err != nil {
                t.Fatalf("Failed to build DAG: %v", err)
        }

        // Validate the DAG
        if err := dag.Validate(); err != nil {
                t.Fatalf("DAG validation failed: %v", err)
        }

        // Check number of nodes (source + sink)
        if len(dag.Nodes) != 2 {
                t.Errorf("Expected 2 nodes, got %d", len(dag.Nodes))
        }

        // Check number of edges (source->sink)
        if len(dag.Edges) != 1 {
                t.Errorf("Expected 1 edge, got %d", len(dag.Edges))
        }
}
