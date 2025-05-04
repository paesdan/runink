Feature: Legal CSA Review and Risk Classification  
Scenario: Extract, Enrich, and Validate Legacy CSA Documents
  Metadata:
    purpose: "Use LLM to extract and classify legacy Credit Support Annexes for compliance and downstream automation"
    module_layer: "Silver"
    herd: "Finance"
    classification: "restricted"
    contract: "legal/csa_review.contract"
    contract_version: "1.1.0"
    slo: "99.9%"
    lineage_tracking: true

  Given source "CSAs from Vault" from "s3://legal/csas/2025/"
  When documents are received
    Then:
      - Parse CSA metadata from JSON
      - Enrich each CSA using LLM for clause summarization
      - Validate clause count >= 5
      - Mask legal reviewer information
      - Tag compliance risk score based on clause content
      - Detect schema drift against approved contract version
      - Route enriched CSAs to "Reviewed CSA Table" (sink: snowflake://LEGAL_DB.REVIEWED_CSAS)

  Assertions:
    - All records must contain `risk_score`
    - No reviewer name may appear unmasked
    - Schema version must match `1.1.0`

  GoldenTest:
    input: "legal/csas/sample.input"
    output: "legal/csas/sample.reviewed.golden"
    validation: strict

  Notifications:
    - On schema drift, alert "alerts/legal_schema_drift"
    - On unmasked reviewer name, alert "alerts/policy_violation"
