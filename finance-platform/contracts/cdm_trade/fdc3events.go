package fdc3events

import (
	"os"
	"strings"
	"time"
)

// ------------------------------
// @herd:"finance"
// @feature:"trade_processing"
// @scenario:"fdc3events"
// ------------------------------

// Step 1: Ingest from Kafka
// @step:"Ingest"
type CDMTradeEvent struct {
	RawPayload string `json:"raw_payload" source:"${CDM_SOURCE_KAFKA}" window:"5m" access:"ingest" format:"json" golden:"${CDM_GOLDEN_INPUT}"`
}

// Step 2: Decode the payload
// @step:"DecodeCDMEvents" @step_input:"Ingest" @step_output:"Decoded"
type DecodedEvent struct {
	TradeID   string    `json:"trade_id" pii:"true" access:"trader"`
	Product   string    `json:"product"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	IsDecoded bool      `json:"is_decoded"`
	TempSink  string    `sink:"${CDM_TEMP_SINK}" purpose:"audit" golden:"${CDM_GOLDEN_DECODED}"`
}

// Step 3: Validate lifecycle & schema
// @step:"ValidateLifecycle" @step_input:"Decoded" @step_output:"Validated"
type ValidatedEvent struct {
	DecodedEvent
	IsValid       bool   `json:"valid"`
	ValidationErr string `json:"validation_error,omitempty"`
	DLQSink       string `sink:"${CDM_DLQ_SINK}" when:"!valid" access:"governance" purpose:"quarantine" golden:"${CDM_GOLDEN_VALIDATED}"`
}

// Step 4: Tag with FDC3 and route to Snowflake
// @step:"TagWithFDC3Context" @step_input:"Validated" @step_output:"Tagged"
type TaggedEvent struct {
	ValidatedEvent
	FDC3Context string `json:"fdc3_context"`
	FinalSink   string `sink:"${CDM_FINAL_SINK}" when:"valid" access:"analytics" pii:"false" golden:"${CDM_GOLDEN_TAGGED}"`
}

// DecodeCDMEvents transforms raw payloads into structured trade events.
func DecodeCDMEvents(input CDMTradeEvent) DecodedEvent {
	return DecodedEvent{
		TradeID:   "T1234", // placeholder parsing
		Product:   "Swap",
		EventType: "NewTrade",
		Timestamp: time.Now(),
		IsDecoded: true,
		TempSink:  os.Getenv("CDM_TEMP_SINK"),
	}
}

// ValidateLifecycle applies lifecycle rules to structured trade events.
func ValidateLifecycle(event DecodedEvent) ValidatedEvent {
	isValid := event.EventType != "" && !strings.Contains(event.Product, "Invalid")
	validationErr := ""
	if !isValid {
		validationErr = "Missing type or invalid product"
	}
	return ValidatedEvent{
		DecodedEvent:  event,
		IsValid:       isValid,
		ValidationErr: validationErr,
		DLQSink:       os.Getenv("CDM_DLQ_SINK"),
	}
}

// TagWithFDC3Context enriches valid trade events with an FDC3 metadata context.
func TagWithFDC3Context(event ValidatedEvent) TaggedEvent {
	context := "fdc3.instrumentView:" + event.TradeID
	return TaggedEvent{
		ValidatedEvent: event,
		FDC3Context:    context,
		FinalSink:      os.Getenv("CDM_FINAL_SINK"),
	}
}
