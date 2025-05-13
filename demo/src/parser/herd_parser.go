
package parser

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ParseHerd parses a herd file and returns a HerdFile struct
func ParseHerd(path string) (HerdFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return HerdFile{}, err
	}
	defer file.Close()

	herd := HerdFile{
		Labels: make(map[string]interface{}),
		RBACPolicies: []RBACPolicy{},
		MaskingPolicies: MaskingPolicies{
			FieldLevelOverrides: []FieldLevelOverride{},
		},
	}

	// Current section being parsed
	var currentSection string
	var currentSubSection string
	var inRBACPolicy bool
	var currentRBACPolicy RBACPolicy
	var inFieldLevelOverride bool
	var currentFieldLevelOverride FieldLevelOverride

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Check if this is a section header
		sectionMatch := regexp.MustCompile(`^\[([a-zA-Z_.]+)\]$`).FindStringSubmatch(line)
		if len(sectionMatch) > 1 {
			section := sectionMatch[1]
			
			// Check if this is a subsection
			if strings.Contains(section, ".") {
				parts := strings.SplitN(section, ".", 2)
				currentSection = parts[0]
				currentSubSection = parts[1]
			} else {
				currentSection = section
				currentSubSection = ""
			}
			
			inRBACPolicy = false
			inFieldLevelOverride = false
			continue
		}
		
		// Check if this is an RBAC policy
		if line == "[[herd.rbac_policies]]" {
			inRBACPolicy = true
			currentRBACPolicy = RBACPolicy{}
			continue
		}
		
		// Check if this is a field level override
		if line == "[[herd.masking_policies.field_level_overrides]]" {
			inFieldLevelOverride = true
			currentFieldLevelOverride = FieldLevelOverride{}
			continue
		}
		
		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		
		// Process based on current section and subsection
		if inRBACPolicy {
			parseRBACPolicy(&currentRBACPolicy, key, value)
			
			// If we've parsed all fields, add to the list
			if currentRBACPolicy.Role != "" && len(currentRBACPolicy.Actions) > 0 {
				herd.RBACPolicies = append(herd.RBACPolicies, currentRBACPolicy)
				inRBACPolicy = false
			}
		} else if inFieldLevelOverride {
			parseFieldLevelOverride(&currentFieldLevelOverride, key, value)
			
			// If we've parsed all fields, add to the list
			if currentFieldLevelOverride.Field != "" && currentFieldLevelOverride.MaskType != "" {
				herd.MaskingPolicies.FieldLevelOverrides = append(
					herd.MaskingPolicies.FieldLevelOverrides, 
					currentFieldLevelOverride,
				)
				inFieldLevelOverride = false
			}
		} else {
			switch currentSection {
			case "herd":
				if currentSubSection == "" {
					parseHerdSection(&herd, key, value)
				} else if currentSubSection == "labels" {
					parseLabelsSection(&herd, key, value)
				} else if currentSubSection == "resource_quotas" {
					parseResourceQuotasSection(&herd, key, value)
				} else if currentSubSection == "secrets_scope" {
					parseSecretsScope(&herd, key, value)
				} else if currentSubSection == "compliance_requirements" {
					parseComplianceRequirements(&herd, key, value)
				} else if currentSubSection == "observability_hooks" {
					parseObservabilityHooks(&herd, key, value)
				} else if currentSubSection == "default_contract_policy" {
					parseDefaultContractPolicy(&herd, key, value)
				} else if currentSubSection == "masking_policies" {
					parseMaskingPolicies(&herd, key, value)
				} else if currentSubSection == "runtime_isolation" {
					parseRuntimeIsolation(&herd, key, value)
				} else if currentSubSection == "retention_policy" {
					parseHerdRetentionPolicy(&herd, key, value)
				} else if currentSubSection == "audit_policy" {
					parseHerdAuditPolicy(&herd, key, value)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return HerdFile{}, err
	}

	return herd, nil
}

// Helper functions to parse each section

func parseHerdSection(herd *HerdFile, key, value string) {
	switch key {
	case "id":
		herd.Herd.ID = value
	case "domain":
		herd.Herd.Domain = value
	case "description":
		herd.Herd.Description = value
	}
}

func parseLabelsSection(herd *HerdFile, key, value string) {
	// Parse array values
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		arrayStr := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
		elements := strings.Split(arrayStr, ",")
		for i, e := range elements {
			elements[i] = strings.Trim(strings.TrimSpace(e), "\"")
		}
		herd.Labels[key] = elements
	} else if value == "true" || value == "false" {
		// Parse boolean values
		herd.Labels[key] = value == "true"
	} else {
		// Parse string values
		herd.Labels[key] = value
	}
}

func parseResourceQuotasSection(herd *HerdFile, key, value string) {
	switch key {
	case "slices_max":
		if intVal, err := strconv.Atoi(value); err == nil {
			herd.ResourceQuotas.SlicesMax = intVal
		}
	case "cpu_limit":
		herd.ResourceQuotas.CPULimit = value
	case "memory_limit":
		herd.ResourceQuotas.MemoryLimit = value
	case "ephemeral_storage":
		herd.ResourceQuotas.EphemeralStorage = value
	case "gpu_limit":
		if intVal, err := strconv.Atoi(value); err == nil {
			herd.ResourceQuotas.GPULimit = intVal
		}
	case "slice_cpu_min":
		herd.ResourceQuotas.SliceCPUMin = value
	case "slice_memory_min":
		herd.ResourceQuotas.SliceMemoryMin = value
	}
}

func parseRBACPolicy(policy *RBACPolicy, key, value string) {
	switch key {
	case "role":
		policy.Role = value
	case "actions":
		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			arrayStr := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
			elements := strings.Split(arrayStr, ",")
			for i, e := range elements {
				elements[i] = strings.Trim(strings.TrimSpace(e), "\"")
			}
			policy.Actions = elements
		}
	}
}

func parseSecretsScope(herd *HerdFile, key, value string) {
	switch key {
	case "allow_cross_herd":
		herd.SecretsScope.AllowCrossHerd = value == "true"
	case "rotation_policy":
		herd.SecretsScope.RotationPolicy = value
	case "encryption":
		herd.SecretsScope.Encryption = value
	case "minimum_tls_version":
		herd.SecretsScope.MinimumTLSVersion = value
	case "approved_key_signers":
		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			arrayStr := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
			elements := strings.Split(arrayStr, ",")
			for i, e := range elements {
				elements[i] = strings.Trim(strings.TrimSpace(e), "\"")
			}
			herd.SecretsScope.ApprovedKeySigners = elements
		}
	case "token_lifetime_seconds":
		if intVal, err := strconv.Atoi(value); err == nil {
			herd.SecretsScope.TokenLifetimeSeconds = intVal
		}
	}
}

func parseComplianceRequirements(herd *HerdFile, key, value string) {
	switch key {
	case "encryption_at_rest":
		herd.ComplianceRequirements.EncryptionAtRest = value == "true"
	case "encryption_in_transit":
		herd.ComplianceRequirements.EncryptionInTransit = value == "true"
	case "audit_logging_enabled":
		herd.ComplianceRequirements.AuditLoggingEnabled = value == "true"
	case "external_auditor_ready":
		herd.ComplianceRequirements.ExternalAuditorReady = value == "true"
	case "automatic_snapshot_on_contract_change":
		herd.ComplianceRequirements.AutomaticSnapshotOnContractChange = value == "true"
	case "approved_storage_backends":
		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			arrayStr := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
			elements := strings.Split(arrayStr, ",")
			for i, e := range elements {
				elements[i] = strings.Trim(strings.TrimSpace(e), "\"")
			}
			herd.ComplianceRequirements.ApprovedStorageBackends = elements
		}
	}
}

func parseObservabilityHooks(herd *HerdFile, key, value string) {
	switch key {
	case "metrics_namespace":
		herd.ObservabilityHooks.MetricsNamespace = value
	case "logs_tag":
		herd.ObservabilityHooks.LogsTag = value
	case "tracing_sample_rate":
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			herd.ObservabilityHooks.TracingSampleRate = floatVal
		}
	}
}

func parseDefaultContractPolicy(herd *HerdFile, key, value string) {
	switch key {
	case "enforce_strict_mode":
		herd.DefaultContractPolicy.EnforceStrictMode = value == "true"
	case "allowed_drift_types":
		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			arrayStr := strings.TrimPrefix(strings.TrimSuffix(value, "]"), "[")
			elements := strings.Split(arrayStr, ",")
			for i, e := range elements {
				elements[i] = strings.Trim(strings.TrimSpace(e), "\"")
			}
			herd.DefaultContractPolicy.AllowedDriftTypes = elements
		}
	case "block_on_drift":
		herd.DefaultContractPolicy.BlockOnDrift = value == "true"
	case "auto_snapshot_on_deploy":
		herd.DefaultContractPolicy.AutoSnapshotOnDeploy = value == "true"
	case "golden_baseline_required":
		herd.DefaultContractPolicy.GoldenBaselineRequired = value == "true"
	case "schema_approval_flow_required":
		herd.DefaultContractPolicy.SchemaApprovalFlowRequired = value == "true"
	}
}

func parseMaskingPolicies(herd *HerdFile, key, value string) {
	switch key {
	case "default_masking":
		herd.MaskingPolicies.DefaultMasking = value
	}
}

func parseFieldLevelOverride(override *FieldLevelOverride, key, value string) {
	switch key {
	case "field":
		override.Field = value
	case "mask_type":
		override.MaskType = value
	}
}

func parseRuntimeIsolation(herd *HerdFile, key, value string) {
	switch key {
	case "ephemeral_user_namespace":
		herd.RuntimeIsolation.EphemeralUserNamespace = value == "true"
	case "pid_namespace_per_slice":
		herd.RuntimeIsolation.PIDNamespacePerSlice = value == "true"
	case "net_namespace_per_slice":
		herd.RuntimeIsolation.NetNamespacePerSlice = value == "true"
	case "mount_namespace_per_slice":
		herd.RuntimeIsolation.MountNamespacePerSlice = value == "true"
	}
}

func parseHerdRetentionPolicy(herd *HerdFile, key, value string) {
	switch key {
	case "lineage_retention_days":
		if intVal, err := strconv.Atoi(value); err == nil {
			herd.RetentionPolicy.LineageRetentionDays = intVal
		}
	case "log_retention_days":
		if intVal, err := strconv.Atoi(value); err == nil {
			herd.RetentionPolicy.LogRetentionDays = intVal
		}
	case "snapshot_retention_days":
		if intVal, err := strconv.Atoi(value); err == nil {
			herd.RetentionPolicy.SnapshotRetentionDays = intVal
		}
	}
}

func parseHerdAuditPolicy(herd *HerdFile, key, value string) {
	switch key {
	case "critical_event_alerting":
		herd.AuditPolicy.CriticalEventAlerting = value == "true"
	case "signed_audit_trails":
		herd.AuditPolicy.SignedAuditTrails = value == "true"
	case "audit_trail_storage":
		herd.AuditPolicy.AuditTrailStorage = value
	}
}
