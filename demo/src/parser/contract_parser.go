
package parser

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ParseContract parses a contract file and returns a ContractFile struct
func ParseContract(path string) (ContractFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return ContractFile{}, err
	}
	defer file.Close()

	contract := ContractFile{}
	
	// Current section being parsed
	var currentSection string
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Check if this is a section header
		sectionMatch := regexp.MustCompile(`^\[([a-zA-Z_]+)\]$`).FindStringSubmatch(line)
		if len(sectionMatch) > 1 {
			currentSection = sectionMatch[1]
			continue
		}
		
		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		
		// Process based on current section
		switch currentSection {
		case "contract":
			parseContractSection(&contract, key, value)
		case "compliance":
			parseComplianceSection(&contract, key, value)
		case "execution":
			parseExecutionSection(&contract, key, value)
		case "sources":
			parseSourcesSection(&contract, key, value)
		case "sinks":
			parseSinksSection(&contract, key, value)
		case "golden":
			parseGoldenSection(&contract, key, value)
		case "alerts":
			parseAlertsSection(&contract, key, value)
		case "retention":
			parseRetentionSection(&contract, key, value)
		case "audit":
			parseAuditSection(&contract, key, value)
		}
	}
	
	if err := scanner.Err(); err != nil {
		return ContractFile{}, err
	}
	
	return contract, nil
}

// Helper functions to parse each section

func parseContractSection(contract *ContractFile, key, value string) {
	switch key {
	case "name":
		contract.Contract.Name = value
	case "version":
		contract.Contract.Version = value
	case "schema_hash":
		contract.Contract.SchemaHash = value
	}
}

func parseComplianceSection(contract *ContractFile, key, value string) {
	switch key {
	case "level":
		// Parse array format: ["SOX", "GDPR", "PCI-DSS"]
		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			value = strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
			elements := strings.Split(value, ",")
			for i, e := range elements {
				elements[i] = strings.Trim(strings.TrimSpace(e), "\"")
			}
			contract.Compliance.Level = elements
		}
	case "classification":
		contract.Compliance.Classification = value
	case "governance_tier":
		contract.Compliance.GovernanceTier = value
	}
}

func parseExecutionSection(contract *ContractFile, key, value string) {
	switch key {
	case "herd":
		contract.Execution.Herd = value
	case "module_layer":
		contract.Execution.ModuleLayer = value
	case "slo_target":
		contract.Execution.SLOTarget = value
	case "default_masking_policy":
		contract.Execution.DefaultMaskingPolicy = value
	case "metrics_namespace":
		contract.Execution.MetricsNamespace = value
	case "logs_tag":
		contract.Execution.LogsTag = value
	case "tracing_sample_rate":
		rate, err := strconv.ParseFloat(value, 64)
		if err == nil {
			contract.Execution.TracingSampleRate = rate
		}
	}
}

func parseSourcesSection(contract *ContractFile, key, value string) {
	switch key {
	case "source_name":
		contract.Sources.SourceName = value
	case "source_uri":
		contract.Sources.SourceURI = value
	}
}

func parseSinksSection(contract *ContractFile, key, value string) {
	switch key {
	case "valid_sink_name":
		contract.Sinks.ValidSinkName = value
	case "valid_sink_uri":
		contract.Sinks.ValidSinkURI = value
	case "invalid_sink_name":
		contract.Sinks.InvalidSinkName = value
	case "invalid_sink_uri":
		contract.Sinks.InvalidSinkURI = value
	}
}

func parseGoldenSection(contract *ContractFile, key, value string) {
	switch key {
	case "input":
		contract.Golden.Input = value
	case "output":
		contract.Golden.Output = value
	}
}

func parseAlertsSection(contract *ContractFile, key, value string) {
	switch key {
	case "on_schema_drift":
		contract.Alerts.OnSchemaDrift = value
	case "on_masking_failure":
		contract.Alerts.OnMaskingFailure = value
	}
}

func parseRetentionSection(contract *ContractFile, key, value string) {
	switch key {
	case "lineage_retention_days":
		days, err := strconv.Atoi(value)
		if err == nil {
			contract.Retention.LineageRetentionDays = days
		}
	case "log_retention_days":
		days, err := strconv.Atoi(value)
		if err == nil {
			contract.Retention.LogRetentionDays = days
		}
	case "snapshot_retention_days":
		days, err := strconv.Atoi(value)
		if err == nil {
			contract.Retention.SnapshotRetentionDays = days
		}
	}
}

func parseAuditSection(contract *ContractFile, key, value string) {
	switch key {
	case "critical_event_alerting":
		contract.Audit.CriticalEventAlerting = value == "true"
	case "storage_backend":
		contract.Audit.StorageBackend = value
	}
}
