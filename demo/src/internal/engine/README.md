# RunInk Execution Engine

The execution engine is responsible for running DAGs (Directed Acyclic Graphs) of nodes in the RunInk system. It provides a robust, flexible framework for executing computational workflows with advanced features for error handling, data passing, monitoring, and isolation.

## Key Features

### Advanced Error Handling and Recovery

The engine implements sophisticated error handling mechanisms:
- Per-node panic recovery to prevent cascading failures
- Configurable error handling policies (fatal vs. non-fatal errors)
- Detailed error collection and reporting
- Graceful shutdown on critical failures

### Typed Channel Data Passing

Data flows between nodes through typed channels:
- Support for both generic (`any`) and strongly-typed data
- Automatic channel management (creation, closing)
- Buffered channels to prevent deadlocks
- Support for multiple input dependencies

### Execution Modes

Two primary execution modes are supported:
- **Batch Mode** (`engine.BatchMode`): Processes all data at once, suitable for bulk operations
- **Streaming Mode** (`engine.StreamingMode`): Processes data as it arrives, ideal for real-time applications

### Monitoring and Observability

Built-in monitoring hooks provide insights into execution:
- Node start/success/error events
- Execution timing and duration tracking
- Customizable monitoring implementations
- Aggregated execution statistics

### Isolation Support

The engine can leverage Linux isolation features:
- Namespace isolation for security
- Resource limits via cgroups
- Configurable isolation policies
- Isolation ID for tracking and management

## Usage Example

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/runink/internal/engine"
)

func main() {
    // Create nodes
    nodeA := &engine.Node{
        ID: "A",
        Function: func(ctx context.Context, in <-chan any) (any, error) {
            return "output from A", nil
        },
        Dependencies: []string{},
        Dependents:   []string{"B", "C"},
    }
    
    nodeB := &engine.Node{
        ID: "B",
        Function: func(ctx context.Context, in <-chan any) (any, error) {
            data := <-in
            return fmt.Sprintf("B processed: %v", data), nil
        },
        Dependencies: []string{"A"},
        Dependents:   []string{"D"},
    }
    
    nodeC := &engine.Node{
        ID: "C",
        Function: func(ctx context.Context, in <-chan any) (any, error) {
            data := <-in
            return fmt.Sprintf("C processed: %v", data), nil
        },
        Dependencies: []string{"A"},
        Dependents:   []string{"D"},
    }
    
    nodeD := &engine.Node{
        ID: "D",
        Function: func(ctx context.Context, in <-chan any) (any, error) {
            // In a real implementation, this would merge multiple inputs
            // For simplicity, we're just reading one
            data := <-in
            return fmt.Sprintf("D finalized: %v", data), nil
        },
        Dependencies: []string{"B", "C"},
        Dependents:   []string{},
    }
    
    // Build the DAG
    dag := engine.BuildDAG([]*engine.Node{nodeA, nodeB, nodeC, nodeD}, true, "run-123")
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Create monitor and error handler
    monitor := engine.NewDefaultMonitor()
    errHandler := engine.NewDefaultErrorHandler(true)
    
    // Execute the DAG
    err := engine.Execute(ctx, dag, engine.BatchMode, monitor, errHandler)
    if err != nil {
        fmt.Printf("Execution failed: %v\n", err)
        for nodeID, nodeErr := range errHandler.GetErrors() {
            fmt.Printf("Node %s error: %v\n", nodeID, nodeErr)
        }
    }
}
```
