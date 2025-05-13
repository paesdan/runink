
// Package parser provides functionality to parse RunInk DSL, contract, conf, and herd files.
package parser

// DSLFile represents the parsed content of a DSL file
type DSLFile struct {
        Feature     string
        Scenario    string
        Metadata    map[string]interface{}
        Source      string
        Steps       []string
        Assertions  []string
        GoldenTest  GoldenTest
        Notifications []Notification
}

// GoldenTest represents golden test configuration in DSL
type GoldenTest struct {
        Input      string
        Output     string
        Validation string
}

// Notification represents notification configuration in DSL
type Notification struct {
        Condition string
        Alert     string
}

// ContractFile represents the parsed content of a contract file
type ContractFile struct {
        Contract   ContractSection
        Compliance ComplianceSection
        Execution  ExecutionSection
        Sources    SourcesSection
        Sinks      SinksSection
        Golden     GoldenSection
        Alerts     AlertsSection
        Retention  RetentionSection
        Audit      AuditSection
}

// ContractSection represents the contract section in a contract file
type ContractSection struct {
        Name       string
        Version    string
        SchemaHash string
}

// ComplianceSection represents the compliance section in a contract file
type ComplianceSection struct {
        Level          []string
        Classification string
        GovernanceTier string
}

// ExecutionSection represents the execution section in a contract file
type ExecutionSection struct {
        Herd                string
        ModuleLayer         string
        SLOTarget           string
        DefaultMaskingPolicy string
        MetricsNamespace    string
        LogsTag             string
        TracingSampleRate   float64
}

// SourcesSection represents the sources section in a contract file
type SourcesSection struct {
        SourceName string
        SourceURI  string
}

// SinksSection represents the sinks section in a contract file
type SinksSection struct {
        ValidSinkName   string
        ValidSinkURI    string
        InvalidSinkName string
        InvalidSinkURI  string
}

// GoldenSection represents the golden section in a contract file
type GoldenSection struct {
        Input  string
        Output string
}

// AlertsSection represents the alerts section in a contract file
type AlertsSection struct {
        OnSchemaDrift    string
        OnMaskingFailure string
}

// RetentionSection represents the retention section in a contract file
type RetentionSection struct {
        LineageRetentionDays  int
        LogRetentionDays      int
        SnapshotRetentionDays int
}

// AuditSection represents the audit section in a contract file
type AuditSection struct {
        CriticalEventAlerting bool
        StorageBackend        string
}

// ConfFile represents the parsed content of a conf file
type ConfFile struct {
        // Conf files typically contain configuration for the runtime environment
        // This is a simplified representation
        Config map[string]interface{}
}

// HerdFile represents the parsed content of a herd file
type HerdFile struct {
        Herd                   HerdSection
        Labels                 map[string]interface{}
        ResourceQuotas         ResourceQuotasSection
        RBACPolicies           []RBACPolicy
        SecretsScope           SecretsScope
        ComplianceRequirements ComplianceRequirements
        ObservabilityHooks     ObservabilityHooks
        DefaultContractPolicy  DefaultContractPolicy
        MaskingPolicies        MaskingPolicies
        RuntimeIsolation       RuntimeIsolation
        RetentionPolicy        RetentionPolicy
        AuditPolicy            AuditPolicy
}

// HerdSection represents the herd section in a herd file
type HerdSection struct {
        ID          string
        Domain      string
        Description string
}

// ResourceQuotasSection represents the resource quotas section in a herd file
type ResourceQuotasSection struct {
        SlicesMax         int
        CPULimit          string
        MemoryLimit       string
        EphemeralStorage  string
        GPULimit          int
        SliceCPUMin       string
        SliceMemoryMin    string
}

// RBACPolicy represents an RBAC policy in a herd file
type RBACPolicy struct {
        Role    string
        Actions []string
}

// SecretsScope represents the secrets scope section in a herd file
type SecretsScope struct {
        AllowCrossHerd      bool
        RotationPolicy      string
        Encryption          string
        MinimumTLSVersion   string
        ApprovedKeySigners  []string
        TokenLifetimeSeconds int
}

// ComplianceRequirements represents the compliance requirements section in a herd file
type ComplianceRequirements struct {
        EncryptionAtRest               bool
        EncryptionInTransit            bool
        AuditLoggingEnabled            bool
        ExternalAuditorReady           bool
        AutomaticSnapshotOnContractChange bool
        ApprovedStorageBackends        []string
}

// ObservabilityHooks represents the observability hooks section in a herd file
type ObservabilityHooks struct {
        MetricsNamespace  string
        LogsTag           string
        TracingSampleRate float64
}

// DefaultContractPolicy represents the default contract policy section in a herd file
type DefaultContractPolicy struct {
        EnforceStrictMode        bool
        AllowedDriftTypes        []string
        BlockOnDrift             bool
        AutoSnapshotOnDeploy     bool
        GoldenBaselineRequired   bool
        SchemaApprovalFlowRequired bool
}

// MaskingPolicies represents the masking policies section in a herd file
type MaskingPolicies struct {
        DefaultMasking      string
        FieldLevelOverrides []FieldLevelOverride
}

// FieldLevelOverride represents a field level override in masking policies
type FieldLevelOverride struct {
        Field    string
        MaskType string
}

// RuntimeIsolation represents the runtime isolation section in a herd file
type RuntimeIsolation struct {
        EphemeralUserNamespace bool
        PIDNamespacePerSlice   bool
        NetNamespacePerSlice   bool
        MountNamespacePerSlice bool
}

// RetentionPolicy represents the retention policy section in a herd file
type RetentionPolicy struct {
        LineageRetentionDays  int
        LogRetentionDays      int
        SnapshotRetentionDays int
}

// AuditPolicy represents the audit policy section in a herd file
type AuditPolicy struct {
        CriticalEventAlerting bool
        SignedAuditTrails     bool
        AuditTrailStorage     string
}
