package dag

import (
	"fmt"

	"github.com/runink/runink/parser"
)

// BuildFromParsedFiles builds a DAG from the parsed files
func BuildFromParsedFiles(dsl parser.DSLFile, contract parser.ContractFile, conf parser.ConfFile, herd parser.HerdFile) (*DAG, error) {
	// Build the basic DAG structure from the DSL
	dag, err := Build(dsl)
	if err != nil {
		return nil, fmt.Errorf("failed to build DAG: %w", err)
	}

	// Enhance the DAG with information from other files
	enhanceWithContract(dag, contract)
	enhanceWithConf(dag, conf)
	enhanceWithHerd(dag, herd)

	// Validate the enhanced DAG
	if err := dag.Validate(); err != nil {
		return nil, fmt.Errorf("DAG validation failed: %w", err)
	}

	return dag, nil
}

// enhanceWithContract adds contract-specific information to the DAG
func enhanceWithContract(dag *DAG, contract parser.ContractFile) {
	// Add contract information to the source node
	if sourceNode, exists := dag.Nodes["source"]; exists {
		if sourceNode.Config == nil {
			sourceNode.Config = make(map[string]interface{})
		}
		
		// Add source information from contract
		if contract.Sources.SourceName != "" {
			sourceNode.Config["name"] = contract.Sources.SourceName
		}
		if contract.Sources.SourceURI != "" {
			sourceNode.Config["uri"] = contract.Sources.SourceURI
		}
	}

	// Add sink information to the sink node
	if sinkNode, exists := dag.Nodes["sink"]; exists {
		if sinkNode.Config == nil {
			sinkNode.Config = make(map[string]interface{})
		}
		
		// Add sink information from contract
		if contract.Sinks.ValidSinkName != "" {
			sinkNode.Config["name"] = contract.Sinks.ValidSinkName
		}
		if contract.Sinks.ValidSinkURI != "" {
			sinkNode.Config["uri"] = contract.Sinks.ValidSinkURI
		}
	}

	// Add contract metadata to all nodes
	for _, node := range dag.Nodes {
		if node.Config == nil {
			node.Config = make(map[string]interface{})
		}
		
		// Add contract information to all nodes
		node.Config["contract_name"] = contract.Contract.Name
		node.Config["contract_version"] = contract.Contract.Version
	}
}

// enhanceWithConf adds configuration-specific information to the DAG
func enhanceWithConf(dag *DAG, conf parser.ConfFile) {
	// Apply configuration to all nodes
	for _, node := range dag.Nodes {
		if node.Config == nil {
			node.Config = make(map[string]interface{})
		}
		
		// Add any relevant configuration from conf file
		// This is a simplified example - in a real implementation,
		// you would match configuration keys to specific nodes
		for key, value := range conf.Config {
			// Only apply configuration if it's relevant to this node type
			if isConfigRelevantToNode(key, node) {
				node.Config[key] = value
			}
		}
	}
}

// enhanceWithHerd adds herd-specific information to the DAG
func enhanceWithHerd(dag *DAG, herd parser.HerdFile) {
	// Apply herd information to all nodes
	for _, node := range dag.Nodes {
		if node.Config == nil {
			node.Config = make(map[string]interface{})
		}
		
		// Add herd information to all nodes
		node.Config["herd_id"] = herd.Herd.ID
		node.Config["herd_domain"] = herd.Herd.Domain
		
		// Add resource constraints
		node.Config["cpu_limit"] = herd.ResourceQuotas.CPULimit
		node.Config["memory_limit"] = herd.ResourceQuotas.MemoryLimit
	}
}

// isConfigRelevantToNode determines if a configuration key is relevant to a specific node
func isConfigRelevantToNode(key string, node *Node) bool {
	// This is a simplified implementation
	// In a real system, you would have more sophisticated logic to determine
	// which configuration keys apply to which node types
	
	// Example: Apply "timeout" config to all nodes
	if key == "timeout" {
		return true
	}
	
	// Example: Apply "batch_size" only to transform nodes
	if key == "batch_size" && node.Type == "transform" {
		return true
	}
	
	// Example: Apply "filter_condition" only to filter nodes
	if key == "filter_condition" && node.Type == "filter" {
		return true
	}
	
	return false
}
