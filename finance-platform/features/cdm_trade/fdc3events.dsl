Feature: Trade Processing

@herd("finance")
@feature("trade_processing")
@scenario("stream_and_enrich")
@contract("trade_cdm_multi.go")
Scenario: Stream and Enrich CDM Trades

  When CDM trade events arrive in real time from Kafka
  Then decode the payload into structured trade events
  And validate product types and event transitions
  And route invalid records to the audit control table
  And tag valid records with FDC3 context
  Do write enriched valid trades to Snowflake for downstream analysis
