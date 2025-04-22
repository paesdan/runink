# Advanced Feature DSL Syntax – Runink

Runink uses `.feature` files with a human-friendly DSL (Domain-Specific Language) inspired by BDD and Gherkin. These files define declarative data pipelines using structured steps, annotations, and contract references.

This guide explores advanced syntax available for real-world data use cases including streaming, branching, conditionals, role enforcement, and metadata tagging.

---

## 📌 Anatomy of a `.feature` File

```gherkin
Feature: High-value customer segmentation

@module(layer="bronze", domain=)
Scenario: Normalize Bronze Layer Product Data
  @source("json://testdata/input/products_feed.json")
  @step("AddIngestMetadata")
  @step("TagMissingFields")
  @sink("json://bronze/products_bronze.json")

  Given the contract: contracts/products_raw.json

  When data is received from the vendor feed
  Then the ingestion timestamp should be added
  And missing required fields (sku, name, price) should be tagged
  Do log metadata for every record

---

## Scenario: Silver Layer - Normalize and Standardize

@module(layer="silver", domain=)
Scenario: Standardize Product Schema
  @source("json://bronze/products_bronze.json")
  @step("TrimProductNames")
  @step("StandardizeCurrency")
  @step("FixEmptyDescriptions")
  @sink("json://silver/products_silver.json")

  Given the contract: contracts/products_normalized.json

  When records contain inconsistent formatting
  Then product names should be trimmed
  And currencies should default to USD if missing
  And empty descriptions replaced with a default message
  Do emit standardized and validated output

---

## Scenario: Gold Layer - Curate for Analytics & Governance

@module(layer="gold", domain=)
Scenario: Enrich and Curate Product Catalog
  @source("json://silver/products_silver.json")
  @step("GroupVariantsByFamily")
  @step("EnrichWithCategoryLTV")
  @step("DetectDiscontinuedItems")
  @sink("json://gold/products_curated.json")

  Given the contract: contracts/products_curated.json

  When normalized product data is ready
  Then variants should be grouped by SKU family
  And categories should have a calculated LTV score
  And discontinued items should be flagged by description
  Do finalize product output with metadata for BI usage
```

---

## 🔁 Iterative Scenarios
Use multiple `Scenario` blocks per pipeline variant:
```gherkin
Scenario: Flag suspicious orders
Scenario: Apply loyalty points
```

---

## 🔄 Streaming with Windowing
```gherkin
@source("kafka://events.orders", @window(5m))
```

Windowing modes:
- `@window(5m)` — Tumbling
- `@window(sliding:10m, every:2m)` — Sliding
- `@window(session:15m)` — Session

---

## ⚠️ Conditional Routing
Send records to sinks based on dynamic conditions:
```gherkin
@sink("json://clean/users.json" when "user.active == true")
@sink("json://dlq/inactive.json" when "user.active == false")
```

---

## 🔐 Role-based Output
```gherkin
@sink("json://ops_view.json" with "role=ops")
@sink("json://finance_view.json" with "role=finance")
```
Paired with `access:"role"` in contracts.

---

## 🧪 Test Hooks
```gherkin
@golden("testdata/orders.golden.json")
@assert("records == 100")
```

---

## 📎 Source Metadata Tags
```gherkin
@source("csv://products.csv" @meta(region="us", vendor="x"))
```
This metadata is passed into transforms and lineage.

---

## 🧬 Reusing Steps (Registry)
Steps like `NormalizeEmail`, `FilterInactiveUsers`, `TagLTV` can be reused across scenarios with no code duplication.

---

## 🧰 Combining Sources
```gherkin
@source("csv://orders.csv")
@source("json://users.json")
@step("JoinOrdersAndUsers")
```
Each record is tagged with source path in the pipeline.

---

## 🧹 Side Effects and Emissions
```gherkin
@emits("alerts/vip_discovered")
@step("SendWebhook")
```
Triggers external systems while running pipelines.

---

## 🧠 LLM-based Annotations
Use AI-generated transform suggestions via:
```gherkin
@auto("summarize user activity")
```
Runink can then generate or recommend pipeline stages for the task.

---

## Summary
Runink’s `.feature` DSL lets you:
- Describe pipelines in natural, reusable syntax
- Build complex branching and streaming workflows
- Embed contracts, policies, and roles into your ETL

It’s not just testable — it’s self-documenting, composable, and production-ready.
