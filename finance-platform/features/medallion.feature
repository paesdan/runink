@feature(name="CDM Trade Processing", herd="cdm_trade/finance")
@contract(path="cdm_trade/trade_cdm_multi.go")

---

Scenario: Bronze Layer - Ingest CDM Trade Events
  @module(layer="bronze", herd="cdm_trade/finance")
  @source("kafka:topics.trade_events" @window(5m))
  @step("DecodeCDMEvents")
  @sink("sf:bronze.trade_cdm_decoded")

  When raw trade messages arrive
  Then decode the payload into structured CDM trade events
  Do store decoded payloads to audit and downstream consumption

---

Scenario: Silver Layer - Validate Lifecycle Rules
  @module(layer="silver", herd="cdm_trade/finance")
  @source("sf:bronze.trade_cdm_decoded")
  @step("ValidateLifecycle")
  @sink("sf:silver.trade_cdm_validated")
  @sink("sf:control.trade_invalid_events" when "!record.valid")

  Given the contract: trade_cdm_validated

  When decoded trade events are ingested
  Then validate schema and lifecycle transitions
  And route invalid messages to DLQ
  Do store valid events for enrichment

---

Scenario: Gold Layer - Enrich with FDC3
  @module(layer="gold", herd="finance")
  @source("sf://silver.trade_cdm_validated")
  @step("TagWithFDC3Context")
  @sink("sf://gold.trade_enriched_fdc3")

  @emits("alerts.trade_fdc3_tags")

  Given the contract: trade_cdm_tagged

  When validated trades are available
  Then enrich them with FDC3 context
  Do emit enrichment metadata and write to Snowflake
