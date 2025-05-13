package engine

import (
	"github.com/runink/runink/internal/parser"
)

// ParseFiles parses all the input files and returns the parsed data
func ParseFiles(config DAGConfig) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	// Parse contract file
	contractParser := parser.NewContractParser(config.ContractFile)
	contractData, err := contractParser.Parse()
	if err != nil {
		return nil, err
	}
	result["contract"] = contractData
	
	// Parse conf file
	confParser := parser.NewConfParser(config.ConfFile)
	confData, err := confParser.Parse()
	if err != nil {
		return nil, err
	}
	result["conf"] = confData
	
	// Parse DSL file
	dslParser := parser.NewDSLParser(config.DSLFile)
	dslData, err := dslParser.Parse()
	if err != nil {
		return nil, err
	}
	result["dsl"] = dslData
	
	// Parse herd file if provided
	if config.HerdFile != "" {
		herdParser := parser.NewHerdParser(config.HerdFile)
		herdData, err := herdParser.Parse()
		if err != nil {
			return nil, err
		}
		result["herd"] = herdData
	}
	
	return result, nil
}
