# Feature: Bronze to Gold Product Catalog ETL (Medallion)

This feature file describes the transformation of raw vendor product data into a business-curated catalog across the Bronze → Silver → Gold layers using clearly defined pipeline steps and Gherkin-style semantics.

---

## Scenario: Bronze Layer - Ingest Raw Product Data

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
