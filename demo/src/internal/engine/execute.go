// Package engine provides the core DAG execution functionality for RunInk.
package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ExecutionMode defines how the DAG should be executed
type ExecutionMode int

const (
	// BatchMode executes the DAG in batch mode, where each node processes all inputs at once
	BatchMode ExecutionMode = iota
	
	// StreamingMode executes the DAG in streaming mode, where nodes process data as it arrives
	StreamingMode
)

// String returns a string representation of the execution mode
func (m ExecutionMode) String() string {
	switch m {
	case BatchMode:
		return "Batch"
	case StreamingMode:
		return "Streaming"
	default:
		return "Unknown"
	}
}

// ErrorHandler defines the interface for handling errors during DAG execution
type ErrorHandler interface {
	// Handle processes an error that occurred in a specific node
	// Returns true if execution should continue, false if it should stop
	Handle(nodeID string, err error) bool
	
	// GetErrors returns all errors that occurred during execution
	GetErrors() map[string]error
}

// DefaultErrorHandler provides a basic implementation of ErrorHandler
type DefaultErrorHandler struct {
	mu     sync.Mutex
	errors map[string]error
	fatal  bool
}

// NewDefaultErrorHandler creates a new DefaultErrorHandler
func NewDefaultErrorHandler(fatal bool) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		errors: make(map[string]error),
		fatal:  fatal,
	}
}

// Handle implements ErrorHandler.Handle
func (h *DefaultErrorHandler) Handle(nodeID string, err error) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.errors[nodeID] = err
	return !h.fatal
}

// GetErrors implements ErrorHandler.GetErrors
func (h *DefaultErrorHandler) GetErrors() map[string]error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Return a copy to avoid concurrent map access
	result := make(map[string]error, len(h.errors))
	for k, v := range h.errors {
		result[k] = v
	}
	return result
}

// Monitor defines the interface for monitoring DAG execution
type Monitor interface {
	// OnStart is called when a node starts execution
	OnStart(nodeID string)
	
	// OnSuccess is called when a node completes successfully
	OnSuccess(nodeID string, duration time.Duration)
	
	// OnError is called when a node encounters an error
	OnError(nodeID string, err error, duration time.Duration)
	
	// OnFinish is called when the entire DAG execution completes
	OnFinish(success bool, duration time.Duration)
}

// DefaultMonitor provides a basic implementation of Monitor
type DefaultMonitor struct{}

// NewDefaultMonitor creates a new DefaultMonitor
func NewDefaultMonitor() *DefaultMonitor {
	return &DefaultMonitor{}
}

// OnStart implements Monitor.OnStart
func (m *DefaultMonitor) OnStart(nodeID string) {
	fmt.Printf("Starting node: %s\n", nodeID)
}

// OnSuccess implements Monitor.OnSuccess
func (m *DefaultMonitor) OnSuccess(nodeID string, duration time.Duration) {
	fmt.Printf("Node %s completed successfully in %v\n", nodeID, duration)
}

// OnError implements Monitor.OnError
func (m *DefaultMonitor) OnError(nodeID string, err error, duration time.Duration) {
	fmt.Printf("Node %s failed after %v: %v\n", nodeID, duration, err)
}

// OnFinish implements Monitor.OnFinish
func (m *DefaultMonitor) OnFinish(success bool, duration time.Duration) {
	status := "successfully"
	if !success {
		status = "with errors"
	}
	fmt.Printf("DAG execution completed %s in %v\n", status, duration)
}

// Node represents a node in the DAG
type Node struct {
	ID           string
	Function     func(context.Context, <-chan any) (any, error)
	Dependencies []string
	Dependents   []string
}

// DAG represents a directed acyclic graph of nodes
type DAG struct {
	Nodes       map[string]*Node
	TopOrder    []string
	Isolate     bool // Whether to use Linux isolation features
	IsolationID string // Identifier for isolation namespace
}

// Execute runs the DAG with the specified execution mode, monitor, and error handler
func Execute(ctx context.Context, dag *DAG, mode ExecutionMode, monitor Monitor, errHandler ErrorHandler) error {
	if monitor == nil {
		monitor = NewDefaultMonitor()
	}
	
	if errHandler == nil {
		errHandler = NewDefaultErrorHandler(true)
	}
	
	startTime := time.Now()
	
	// Create a cancellable context for the execution
	execCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	
	// Create channels for data flow between nodes
	channels := make(map[string]chan any)
	for _, node := range dag.Nodes {
		for _, dep := range node.Dependents {
			channels[node.ID+"->"+dep] = make(chan any, 1)
		}
	}
	
	// Track remaining dependencies for each node
	remainingDeps := make(map[string]int)
	for _, node := range dag.Nodes {
		remainingDeps[node.ID] = len(node.Dependencies)
	}
	
	// Find nodes with no dependencies to start with
	var readyNodes []string
	for id, count := range remainingDeps {
		if count == 0 {
			readyNodes = append(readyNodes, id)
		}
	}
	
	// WaitGroup to track all running nodes
	var wg sync.WaitGroup
	
	// Mutex to protect access to shared state
	var mu sync.Mutex
	
	// Track execution success
	success := true
	
	// Function to execute a single node
	executeNode := func(nodeID string) {
		defer wg.Done()
		
		node := dag.Nodes[nodeID]
		
		// Setup isolation if required
		if dag.Isolate {
			// In a real implementation, this would set up Linux namespaces and cgroups
			// For this demo, we'll just log that isolation would happen
			fmt.Printf("Node %s would run in isolated environment %s\n", nodeID, dag.IsolationID)
		}
		
		// Notify monitor that node is starting
		nodeStartTime := time.Now()
		monitor.OnStart(nodeID)
		
		// Collect input channels
		var inputChan chan any
		if len(node.Dependencies) > 0 {
			// In a real implementation with multiple dependencies, we would merge channels
			// For simplicity, we'll assume just one input channel per node
			if len(node.Dependencies) == 1 {
				depID := node.Dependencies[0]
				inputChan = channels[depID+"->"+nodeID]
			} else {
				// Create a merged channel for multiple inputs
				inputChan = make(chan any)
				
				// Close the merged channel when all input channels are closed
				go func() {
					defer close(inputChan)
					
					// Wait for all input channels to close
					var inputWg sync.WaitGroup
					for _, depID := range node.Dependencies {
						depChan := channels[depID+"->"+nodeID]
						inputWg.Add(1)
						
						go func(ch <-chan any) {
							defer inputWg.Done()
							for data := range ch {
								// In streaming mode, forward each item
								// In batch mode, collect items (simplified here)
								if mode == StreamingMode {
									select {
									case inputChan <- data:
									case <-execCtx.Done():
										return
									}
								}
							}
						}(depChan)
					}
					
					inputWg.Wait()
				}()
			}
		} else {
			// For source nodes with no dependencies, create an empty channel
			inputChan = make(chan any)
			close(inputChan)
		}
		
		// Execute the node function with panic recovery
		var result any
		var err error
		
		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic in node %s: %v", nodeID, r)
				}
			}()
			
			// Execute the node function
			result, err = node.Function(execCtx, inputChan)
		}()
		
		duration := time.Since(nodeStartTime)
		
		if err != nil {
			// Handle error
			monitor.OnError(nodeID, err, duration)
			continueExecution := errHandler.Handle(nodeID, err)
			
			mu.Lock()
			success = false
			mu.Unlock()
			
			if !continueExecution {
				cancel() // Cancel context to stop other nodes
			}
		} else {
			// Node completed successfully
			monitor.OnSuccess(nodeID, duration)
			
			// Send result to dependent nodes
			for _, depID := range node.Dependents {
				outChan := channels[nodeID+"->"+depID]
				
				// In batch mode, send the entire result at once
				// In streaming mode, the node function should have handled streaming
				if mode == BatchMode {
					select {
					case outChan <- result:
					case <-execCtx.Done():
						// Context cancelled, stop sending
					}
				}
				
				// Close output channels to signal completion
				close(outChan)
			}
			
			// Mark dependent nodes as ready if all their dependencies are satisfied
			mu.Lock()
			for _, depID := range node.Dependents {
				remainingDeps[depID]--
				if remainingDeps[depID] == 0 {
					// This node is now ready to execute
					go func(id string) {
						wg.Add(1)
						executeNode(id)
					}(depID)
				}
			}
			mu.Unlock()
		}
	}
	
	// Start execution with nodes that have no dependencies
	for _, nodeID := range readyNodes {
		wg.Add(1)
		go executeNode(nodeID)
	}
	
	// Wait for all nodes to complete
	wg.Wait()
	
	// Notify monitor of completion
	totalDuration := time.Since(startTime)
	monitor.OnFinish(success, totalDuration)
	
	// Return error if execution was not successful
	if !success {
		return errors.New("DAG execution failed, check error handler for details")
	}
	
	return nil
}

// BuildDAG creates a new DAG from a list of nodes and their dependencies
// This is a helper function to simplify DAG creation
func BuildDAG(nodes []*Node, isolate bool, isolationID string) *DAG {
	nodeMap := make(map[string]*Node)
	
	// Add nodes to the map
	for _, node := range nodes {
		nodeMap[node.ID] = node
	}
	
	// Compute topological order
	var topOrder []string
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	
	var visit func(string)
	visit = func(nodeID string) {
		if temp[nodeID] {
			// Cycle detected
			panic(fmt.Sprintf("cycle detected in DAG at node %s", nodeID))
		}
		
		if !visited[nodeID] {
			temp[nodeID] = true
			
			node := nodeMap[nodeID]
			for _, depID := range node.Dependents {
				visit(depID)
			}
			
			visited[nodeID] = true
			temp[nodeID] = false
			
			// Prepend to topOrder (will be reversed later)
			topOrder = append([]string{nodeID}, topOrder...)
		}
	}
	
	// Visit all nodes
	for nodeID := range nodeMap {
		if !visited[nodeID] {
			visit(nodeID)
		}
	}
	
	return &DAG{
		Nodes:       nodeMap,
		TopOrder:    topOrder,
		Isolate:     isolate,
		IsolationID: isolationID,
	}
}
