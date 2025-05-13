// Package nodes provides node implementations for the RunInk DAG execution engine.
package nodes

import (
	"fmt"
	"sync"

	"github.com/runink/runink/dag"
)

// NodeFactory is a function that creates a new node
type NodeFactory func(id string, config map[string]interface{}) (interface{}, error)

// NodeRegistry is a registry of node factories
type NodeRegistry struct {
	factories map[string]NodeFactory
	mu        sync.RWMutex
}

// NewNodeRegistry creates a new node registry
func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		factories: make(map[string]NodeFactory),
	}
}

// Register registers a node factory with the registry
func (r *NodeRegistry) Register(nodeType string, factory NodeFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[nodeType] = factory
}

// Get returns a node factory for the given node type
func (r *NodeRegistry) Get(nodeType string) (NodeFactory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	factory, ok := r.factories[nodeType]
	if !ok {
		return nil, fmt.Errorf("no factory registered for node type: %s", nodeType)
	}
	return factory, nil
}

// CreateNode creates a new node of the given type
func (r *NodeRegistry) CreateNode(nodeType, id string, config map[string]interface{}) (interface{}, error) {
	factory, err := r.Get(nodeType)
	if err != nil {
		return nil, err
	}
	return factory(id, config)
}

// DefaultRegistry is the default node registry
var DefaultRegistry = NewNodeRegistry()

// Register registers a node factory with the default registry
func Register(nodeType string, factory NodeFactory) {
	DefaultRegistry.Register(nodeType, factory)
}

// Initialize the registry with built-in node types
func init() {
	// Register CSV reader node
	Register("csv_reader", func(id string, config map[string]interface{}) (interface{}, error) {
		return NewCSVReaderNode(id, config)
	})

	// Register Parquet writer node
	Register("parquet_writer", func(id string, config map[string]interface{}) (interface{}, error) {
		return NewParquetWriterNode(id, config)
	})
}

// CreateDAGNode creates a DAG node from a node implementation
func CreateDAGNode(id string, nodeImpl interface{}, config map[string]interface{}) *dag.Node {
	// Create a DAG node that wraps our implementation
	return &dag.Node{
		ID:     id,
		Name:   fmt.Sprintf("%T", nodeImpl),
		Type:   fmt.Sprintf("%T", nodeImpl),
		Config: config,
		Function: func() error {
			// This would be implemented to call the Execute method of the node
			// For now, just return nil
			return nil
		},
	}
}
