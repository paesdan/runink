Feature: Financial Trade Event Validation & Compliance  
Scenario: Validate and Enrich Incoming Trade Events
  Metadata:
    purpose: "Ensure all trade events meet mandatory compliance and security requirements before downstream analytics."
    module_layer: "Silver"
    herd: "Finance"
    lineage_tracking: true
    slo: "99.95%"
    classification: "restricted"
    compliance: ["SOX", "GDPR", "PCI-DSS"]
    contract: "cdm_trade/fdc3events.contract"
    contract_version: "1.0.0"

  Given source stream "Trade Events Stream" from "kafka://topics.trade_events"

  When events are received
    Then:
      - Decode each trade event according to CDM format
      - Validate mandatory fields (trade_id, symbol, price, timestamp)
      - Apply field masking to sensitive fields (SSN, bank accounts, emails)
      - Tag compliance metadata (classification, region, compliance tags)
      - Detect schema drift against approved contract version
      - Route valid records to "Validated Trades Table" (sink: snowflake://TRADE_PIPELINE.VALIDATED_TRADES)
      - Route invalid records to "DLQ for Invalid Trades" (sink: snowflake://DLQ.INVALID_TRADES)

  Assertions:
    - Processed record count must be >= 1000
    - No schema drift must be detected
    - All masked fields must pass redaction/tokenization checks

  GoldenTest:
    input: "cdm_trade/fdc3events.input"
    output: "cdm_trade/data/fdc3events.validated.golden"
    validation: strict

  Notifications:
    - On schema validation failure, emit alert to "alerts/finance_data_validation"
    - On masking failure, emit alert to "alerts/finance_security_breach"
