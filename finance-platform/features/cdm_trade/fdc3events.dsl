feature: 
  name:
    - trade_processing
  herd: 
    - finance
  contract: 
    name: "cdm_trade/fdc3events"
    version: "1.0.0"
    path: "contracts/cdm_trade/fdc3events.go"
    
  scenario: 
    name: cdm_validation_fdc3
    description: "Stream and Enrich CDM Trades"
    layer: "bronze"

  steps:
    ingest:
      name: KafkaConn
      source_env: CDM_SOURCE_KAFKA
      topic_env: CDM_KAFKA_TOPIC
      group_id_env: CDM_KAFKA_GROUP_ID
      window: 15m

    when: 
      - "cdm_trade_events arrive from Kafka in real time, aggregated every 15 minutes."

    then:
      - action: "decode the incoming message into a trade_event"
        step: DecodeCDMTrade
        requires_state: false
        checkpoint_after: true

      - action: "validate the trade_event’s product_type and event transitions"
        step: ValidateTradeTypes
        requires_state: true
        join_key: trade_id

      - action: "tag valid trade_events with FDC3 context"
        step: EnrichWithFDC3Context
        requires_state: false

      - action: "route invalid trade_events to the audit sink"
        step: RouteToAuditIfInvalid
        requires_state: false

    do:
      - sink_env: CDM_FINAL_SINK
        condition: trade.valid == true
        checkpoint_after: true

      - sink_env: CDM_DLQ_SINK
        condition: trade.valid == false