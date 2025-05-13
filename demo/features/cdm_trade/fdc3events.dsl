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

  Given source stream "Trade Events Stream" from "source/sample.csv"
  properties:
    description: "Schema for user data CSV file"
    fields:
      - name: id
        type: integer
        nullable: false
        description: "Unique identifier for the user"
      - name: name
        type: string
        nullable: false
        description: "Full name of the user"
      - name: age
        type: integer
        nullable: true
        description: "Age of the user in years"
      - name: email
        type: string
        nullable: false
        description: "Email address of the user"
      - name: country
        type: string
        nullable: true
        description: "Country of residence"
      - name: signup_date
        type: date
        nullable: false
        description: "Date when the user signed up"
        format: "yyyy-MM-dd"
      - name: active
        type: boolean
        nullable: false
        description: "Indicates if the user account is currently active"  

  When events are received
    Then:
      - Decode each trade event according to CDM format
      - Validate mandatory fields (trade_id, symbol, price, timestamp)
      - Apply field masking to sensitive fields (SSN, bank accounts, emails)
      - Tag compliance metadata (classification, region, compliance tags)
      - Detect schema drift against approved contract version
      - Route valid records to "Validated Trades Table" (sink: target/validated.csv)
      - Route invalid records to "DLQ for Invalid Trades" (sink: target/dql.csv)

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
