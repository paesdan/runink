package fdc3events

import (
	"strings"
	"time"
)

// ------------------------------
// @contract_version:"1.0.0"
// @herd:"finance"
// @feature:"trade_processing"
// @contract:"cdm_trade/fdc3events.go"
// @description:"Stream and Enrich CDM Trade Events into Snowflake with FDC3 Context Enrichment and Audit Failures."
// @metadata: owner="data.engineering@finance.example", domain="Trade", sla="P1"
// ------------------------------

// Step 0: Kafka Source Ingest
// @step:"KafkaConn"
type KafkaCDMTradeEvent struct {
	RawPayload string `json:"raw_payload" source_env:"CDM_SOURCE_KAFKA" window:"15m" format:"json" access:"ingest" step:"KafkaConn"`
}

// Step 1: DecodeCDMTrade
// @step:"DecodeCDMTrade" @requires_state:"false" @checkpoint_after:"true"
// @emit(lineage:"trade_processing/decoded_trade")
type TradeEvent struct {
	TradeID   string    `json:"trade_id" join_key:"trade_id" pii:"true" access:"trader"`
	Product   string    `json:"product"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	IsDecoded bool      `json:"is_decoded"`
	Step      string    `step:"DecodeCDMTrade" step_input:"KafkaCDMTradeEvent" step_output:"TradeEvent"`
}

// Step 2: ValidateTradeTypes
// @step:"ValidateTradeTypes" @requires_state:"true" @join_key:"trade_id"
// @affinity(group:"finance", colocate_with:"DecodeCDMTrade")
// @emit(lineage:"trade_processing/validated_trade")
type ValidatedTradeEvent struct {
	TradeEvent
	IsValid       bool   `json:"valid"`
	ValidationErr string `json:"validation_error,omitempty"`
	Step          string `step:"ValidateTradeTypes" step_input:"TradeEvent" step_output:"ValidatedTradeEvent"`
}

// Step 3: EnrichWithFDC3Context
// @step:"EnrichWithFDC3Context" @requires_state:"false"
// @emit(lineage:"trade_processing/enriched_with_fdc3")
type EnrichedTradeEvent struct {
	ValidatedTradeEvent
	FDC3Context string `json:"fdc3_context"`
	Step        string `step:"EnrichWithFDC3Context" step_input:"ValidatedTradeEvent" step_output:"EnrichedTradeEvent"`
}

// Step 4: RouteToAuditIfInvalid
// @step:"RouteToAuditIfInvalid"
// @emit(lineage:"trade_processing/audit_failed_trade")
type AuditRecord struct {
	ValidatedTradeEvent
	AuditReason string `json:"audit_reason,omitempty"`
	Step        string `step:"RouteToAuditIfInvalid" step_input:"ValidatedTradeEvent" step_output:"AuditRecord"`
}

// Step 5: LoadValidTrade
// @step:"LoadValidTrade"
type LoadedTradeEvent struct {
	EnrichedTradeEvent
	TargetSink string `sink_env:"CDM_FINAL_SINK"`
	Step       string `step:"LoadValidTrade" step_input:"EnrichedTradeEvent"`
}

// Step 6: LoadInvalidTrade
// @step:"LoadInvalidTrade"
type LoadedAuditRecord struct {
	AuditRecord
	TargetSink string `sink_env:"CDM_DLQ_SINK"`
	Step       string `step:"LoadInvalidTrade" step_input:"AuditRecord"`
}

// ------------------------------
// Transform Functions
// ------------------------------

// DecodeCDMTrade parses raw Kafka payload into a structured trade event.
func DecodeCDMTrade(input KafkaCDMTradeEvent) TradeEvent {
	return TradeEvent{
		TradeID:   "", // Placeholder for pipeline injection
		Product:   "",
		EventType: "",
		Timestamp: time.Now(),
		IsDecoded: true,
		Step:      "DecodeCDMTrade",
	}
}

// ValidateTradeTypes checks trade event for business rule violations.
func ValidateTradeTypes(event TradeEvent) ValidatedTradeEvent {
	isValid := event.EventType != "" && !strings.Contains(event.Product, "Invalid")
	var validationErr string
	if !isValid {
		validationErr = "Missing event type or invalid product category"
	}
	return ValidatedTradeEvent{
		TradeEvent:    event,
		IsValid:       isValid,
		ValidationErr: validationErr,
		Step:          "ValidateTradeTypes",
	}
}

// EnrichWithFDC3Context adds standard FDC3 metadata to valid trades.
func EnrichWithFDC3Context(event ValidatedTradeEvent) EnrichedTradeEvent {
	context := "fdc3.instrumentView:" + event.TradeID
	return EnrichedTradeEvent{
		ValidatedTradeEvent: event,
		FDC3Context:         context,
		Step:                "EnrichWithFDC3Context",
	}
}

// RouteToAuditIfInvalid creates an audit record for invalid trades.
func RouteToAuditIfInvalid(event ValidatedTradeEvent) (AuditRecord, bool) {
	if !event.IsValid {
		return AuditRecord{
			ValidatedTradeEvent: event,
			AuditReason:         event.ValidationErr,
			Step:                "RouteToAuditIfInvalid",
		}, true
	}
	return AuditRecord{}, false
}

// LoadValidTrade prepares validated trade for sink writing (Snowflake)
func LoadValidTrade(event EnrichedTradeEvent) LoadedTradeEvent {
	return LoadedTradeEvent{
		EnrichedTradeEvent: event,
		TargetSink:         "", // Resolved dynamically from CDM_FINAL_SINK
		Step:               "LoadValidTrade",
	}
}

// LoadInvalidTrade prepares invalid trade record for audit sink (DLQ)
func LoadInvalidTrade(audit AuditRecord) LoadedAuditRecord {
	return LoadedAuditRecord{
		AuditRecord: audit,
		TargetSink:  "", // Resolved dynamically from CDM_DLQ_SINK
		Step:        "LoadInvalidTrade",
	}
}
