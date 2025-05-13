
// Package parser provides functionality to parse RunInk DSL, contract, conf, and herd files.
package parser

// This file serves as the main entry point for the parser package.
// It re-exports the main parsing functions for easy access.

// ParseFiles parses all the given files and returns their parsed representations
func ParseFiles(dslPath, contractPath, confPath, herdPath string) (DSLFile, ContractFile, ConfFile, HerdFile, error) {
	dsl, err := ParseDSL(dslPath)
	if err != nil {
		return DSLFile{}, ContractFile{}, ConfFile{}, HerdFile{}, err
	}

	contract, err := ParseContract(contractPath)
	if err != nil {
		return DSLFile{}, ContractFile{}, ConfFile{}, HerdFile{}, err
	}

	conf, err := ParseConf(confPath)
	if err != nil {
		return DSLFile{}, ContractFile{}, ConfFile{}, HerdFile{}, err
	}

	herd, err := ParseHerd(herdPath)
	if err != nil {
		return DSLFile{}, ContractFile{}, ConfFile{}, HerdFile{}, err
	}

	return dsl, contract, conf, herd, nil
}
