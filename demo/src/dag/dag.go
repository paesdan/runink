
// Package dag provides functionality to build and manage a Directed Acyclic Graph (DAG)
// for RunInk pipeline execution.
package dag

import (
        "fmt"
        "strings"

        "github.com/runink/runink/parser"
)

// Node represents a node in the DAG
type Node struct {
        ID          string
        Name        string
        Description string
        Type        string
        Config      map[string]interface{}
        Function    func() error
}

// Edge represents a directed edge between two nodes in the DAG
type Edge struct {
        From string
        To   string
}

// DAG represents a Directed Acyclic Graph with nodes and edges
type DAG struct {
        Nodes map[string]*Node
        Edges []Edge
}

// NewDAG creates a new empty DAG
func NewDAG() *DAG {
        return &DAG{
                Nodes: make(map[string]*Node),
                Edges: []Edge{},
        }
}

// AddNode adds a node to the DAG
func (d *DAG) AddNode(node *Node) {
        d.Nodes[node.ID] = node
}

// AddEdge adds an edge between two nodes in the DAG
func (d *DAG) AddEdge(from, to string) {
        d.Edges = append(d.Edges, Edge{From: from, To: to})
}

// GetNode returns a node by ID
func (d *DAG) GetNode(id string) *Node {
        return d.Nodes[id]
}

// GetPredecessors returns all nodes that are direct predecessors of the given node
func (d *DAG) GetPredecessors(nodeID string) []*Node {
        var predecessors []*Node
        for _, edge := range d.Edges {
                if edge.To == nodeID {
                        predecessors = append(predecessors, d.Nodes[edge.From])
                }
        }
        return predecessors
}

// GetSuccessors returns all nodes that are direct successors of the given node
func (d *DAG) GetSuccessors(nodeID string) []*Node {
        var successors []*Node
        for _, edge := range d.Edges {
                if edge.From == nodeID {
                        successors = append(successors, d.Nodes[edge.To])
                }
        }
        return successors
}

// Build constructs a DAG from the parsed DSL file
func Build(dsl parser.DSLFile) (*DAG, error) {
        dag := NewDAG()
        
        // Create a source node (entry point)
        sourceNode := &Node{
                ID:          "source",
                Name:        "Source",
                Description: "Data source entry point",
                Type:        "source",
                Config: map[string]interface{}{
                        "uri": dsl.Source,
                },
        }
        dag.AddNode(sourceNode)
        
        // Track dependencies for each step
        dependencies := make(map[string][]string)
        
        // First pass: create nodes for each step
        for i, step := range dsl.Steps {
                nodeID := fmt.Sprintf("step_%d", i)
                
                // Parse step to extract name, type, and dependencies
                stepParts := parseStep(step)
                
                node := &Node{
                        ID:          nodeID,
                        Name:        stepParts.name,
                        Description: step,
                        Type:        stepParts.stepType,
                        Config:      stepParts.config,
                }
                
                dag.AddNode(node)
                
                // Store dependencies for the second pass
                dependencies[nodeID] = stepParts.dependencies
        }
        
        // Create a sink node (exit point)
        sinkNode := &Node{
                ID:          "sink",
                Name:        "Sink",
                Description: "Data sink exit point",
                Type:        "sink",
        }
        dag.AddNode(sinkNode)
        
        // Second pass: create edges based on dependencies
        lastNodeID := "source"
        for i := range dsl.Steps {
                nodeID := fmt.Sprintf("step_%d", i)
                
                // If no explicit dependencies, connect to the previous node
                if len(dependencies[nodeID]) == 0 {
                        dag.AddEdge(lastNodeID, nodeID)
                } else {
                        // Add edges for each dependency
                        for _, dep := range dependencies[nodeID] {
                                // Find the node ID for the named dependency
                                depNodeID := findNodeIDByName(dag, dep)
                                if depNodeID != "" {
                                        dag.AddEdge(depNodeID, nodeID)
                                } else {
                                        return nil, fmt.Errorf("dependency '%s' not found for step '%s'", dep, nodeID)
                                }
                        }
                }
                
                lastNodeID = nodeID
        }
        
        // Connect the last step to the sink
        dag.AddEdge(lastNodeID, "sink")
        
        return dag, nil
}

// stepInfo holds parsed information about a step
type stepInfo struct {
        name         string
        stepType     string
        dependencies []string
        config       map[string]interface{}
}

// parseStep extracts information from a step string
func parseStep(step string) stepInfo {
        info := stepInfo{
                config: make(map[string]interface{}),
        }
        
        // Extract step type and name
        parts := strings.SplitN(step, " ", 2)
        if len(parts) > 0 {
                info.stepType = parts[0]
        }
        
        if len(parts) > 1 {
                // Extract name and parameters
                paramStart := strings.Index(parts[1], "(")
                if paramStart != -1 {
                        info.name = strings.TrimSpace(parts[1][:paramStart])
                        
                        // Extract parameters
                        paramStr := parts[1][paramStart+1:]
                        paramEnd := strings.LastIndex(paramStr, ")")
                        if paramEnd != -1 {
                                paramStr = paramStr[:paramEnd]
                                
                                // Parse parameters
                                params := strings.Split(paramStr, ",")
                                for _, param := range params {
                                        param = strings.TrimSpace(param)
                                        
                                        // Check for dependencies
                                        if strings.HasPrefix(param, "after:") || strings.HasPrefix(param, "depends_on:") {
                                                dep := strings.TrimSpace(strings.SplitN(param, ":", 2)[1])
                                                info.dependencies = append(info.dependencies, dep)
                                        } else if strings.Contains(param, ":") {
                                                // Regular config parameter
                                                kv := strings.SplitN(param, ":", 2)
                                                if len(kv) == 2 {
                                                        key := strings.TrimSpace(kv[0])
                                                        value := strings.TrimSpace(kv[1])
                                                        
                                                        // Remove quotes if present
                                                        value = strings.Trim(value, "\"'")
                                                        
                                                        info.config[key] = value
                                                }
                                        }
                                }
                        }
                } else {
                        info.name = strings.TrimSpace(parts[1])
                }
        }
        
        return info
}

// findNodeIDByName finds a node ID by its name
func findNodeIDByName(dag *DAG, name string) string {
        for id, node := range dag.Nodes {
                if node.Name == name {
                        return id
                }
        }
        return ""
}

// TopologicalSort returns nodes in topological order (from source to sink)
func (d *DAG) TopologicalSort() ([]*Node, error) {
        // Create a map to track visited nodes
        visited := make(map[string]bool)
        temp := make(map[string]bool)
        order := []*Node{}
        
        // Define a recursive visit function
        var visit func(string) error
        visit = func(nodeID string) error {
                // Check for cycles
                if temp[nodeID] {
                        return fmt.Errorf("cycle detected in DAG")
                }
                
                // Skip if already visited
                if visited[nodeID] {
                        return nil
                }
                
                temp[nodeID] = true
                
                // Visit all predecessors first
                for _, edge := range d.Edges {
                        if edge.To == nodeID {
                                if err := visit(edge.From); err != nil {
                                        return err
                                }
                        }
                }
                
                temp[nodeID] = false
                visited[nodeID] = true
                
                // Add to order
                order = append(order, d.Nodes[nodeID])
                
                return nil
        }
        
        // Start with sink node to ensure we get a complete path
        if _, exists := d.Nodes["sink"]; exists {
                if err := visit("sink"); err != nil {
                        return nil, err
                }
        }
        
        // Visit any remaining nodes
        for nodeID := range d.Nodes {
                if !visited[nodeID] {
                        if err := visit(nodeID); err != nil {
                                return nil, err
                        }
                }
        }
        
        return order, nil
}

// Validate checks if the DAG is valid (no cycles, all dependencies exist)
func (d *DAG) Validate() error {
        // Check if all nodes in edges exist
        for _, edge := range d.Edges {
                if _, exists := d.Nodes[edge.From]; !exists {
                        return fmt.Errorf("node '%s' referenced in edge does not exist", edge.From)
                }
                if _, exists := d.Nodes[edge.To]; !exists {
                        return fmt.Errorf("node '%s' referenced in edge does not exist", edge.To)
                }
        }
        
        // Check for cycles using topological sort
        _, err := d.TopologicalSort()
        return err
}

// String returns a string representation of the DAG
func (d *DAG) String() string {
        var sb strings.Builder
        
        sb.WriteString("DAG:\n")
        sb.WriteString("Nodes:\n")
        for id, node := range d.Nodes {
                sb.WriteString(fmt.Sprintf("  %s: %s (%s)\n", id, node.Name, node.Type))
        }
        
        sb.WriteString("Edges:\n")
        for _, edge := range d.Edges {
                fromNode := d.Nodes[edge.From]
                toNode := d.Nodes[edge.To]
                sb.WriteString(fmt.Sprintf("  %s (%s) -> %s (%s)\n", 
                        fromNode.Name, edge.From, 
                        toNode.Name, edge.To))
        }
        
        return sb.String()
}
