# Feature: Product Catalog Ingestion into Medallion Architecture

## 🔶 Scenario: Bronze – Ingest Raw Products
@module("bronze")
Scenario: IngestRawProductData
  @source("csv://bronze/products_feed.csv")
  @step("AddIngestMetadata")
  @step("TagMissingFields")
  @sink("json://bronze/products_bronze.json")
  Given the contract: contracts/products_raw.json

## 🔷 Scenario: Silver – Normalize Products
@module("silver")
Scenario: NormalizeProductDetails
  @source("json://bronze/products_bronze.json")
  @step("TrimProductNames")
  @step("StandardizeCurrency")
  @step("FixEmptyDescriptions")
  @sink("json://silver/products_silver.json")
  Given the contract: contracts/products_normalized.json

## 🟡 Scenario: Gold – Curated Product Catalog
@module("gold")
Scenario: CurateBusinessReadyProducts
  @source("json://silver/products_silver.json")
  @step("GroupVariantsByFamily")
  @step("EnrichWithCategoryLTV")
  @step("DetectDiscontinuedItems")
  @sink("json://gold/products_curated.json")
  Given the contract: contracts/products_curated.json